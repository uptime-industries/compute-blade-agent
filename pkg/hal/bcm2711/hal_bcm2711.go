package bcm2711

import (
	"errors"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/xvzf/computeblade-agent/pkg/hal"
)

const (
	bcm2711PeripheryBaseAddr = 0xFE000000
	bcm2711PwmAddr           = bcm2711PeripheryBaseAddr + 0x20C000
	bcm2711GpioAddr          = bcm2711PeripheryBaseAddr + 0x200000
	bcm2711ClkAddr           = bcm2711PeripheryBaseAddr + 0x101000
	bcm2711ClkManagerPwd     = (0x5A << 24) //(31 - 24) on CM_GP0CTL/CM_GP1CTL/CM_GP2CTL regs
	bcm2711PageSize          = 4096         // theoretical page size

	bcm2711FrontButtonPin = 20
	bcm2711StealthPin     = 21
	bcm2711PwmFanPin      = 12
	bcm2711PwmTachPin     = 13

	GPFSEL0 = 0x00
	GPFSEL1 = 0x01
	GPFSEL2 = 0x02

	PWM_CTL  = 0x00
	PWM_RNG1 = 0x04
	PWM_DAT1 = 0x05

	PWMCLK_CNTL = 40
	PWMCLK_DIV  = 41
)

type bcm2711bcm struct {
	// Config options
	opts hal.ComputeBladeHalOpts

	wrMutex sync.Mutex

	// Keep track of the currently set fanspeed so it can later be restored after setting the ws281x LEDs
	currFanSpeed uint8

	devmem   *os.File
	mbox     *os.File
	gpioMem8 []uint8
	gpioMem  []uint32
	pwmMem8  []uint8
	pwmMem   []uint32
	clkMem8  []uint8
	clkMem   []uint32
}

func New(opts hal.ComputeBladeHalOpts) (*bcm2711bcm, error) {
	// /dev/gpiomem doesn't allow complex operations for PWM fan control or WS281x
	devmem, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// /dev/vcio for ioctl with VC mailbox
	mbox, err := os.OpenFile("/dev/vcio", os.O_RDWR|os.O_SYNC, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Setup memory mappings
	gpioMem, gpioMem8, err := mmap(devmem, bcm2711GpioAddr, bcm2711PageSize)
	if err != nil {
		return nil, err
	}
	pwmMem, pwmMem8, err := mmap(devmem, bcm2711PwmAddr, bcm2711PageSize)
	if err != nil {
		return nil, err
	}
	clkMem, clkMem8, err := mmap(devmem, bcm2711ClkAddr, bcm2711PageSize)
	if err != nil {
		return nil, err
	}

	return &bcm2711bcm{
		devmem:   devmem,
		mbox:     mbox,
		gpioMem:  gpioMem,
		gpioMem8: gpioMem8,
		pwmMem:   pwmMem,
		pwmMem8:  pwmMem8,
		clkMem:   clkMem,
		clkMem8:  clkMem8,
		opts:     opts,
	}, nil
}

// Close cleans all memory mappings
func (bcm *bcm2711bcm) Close() error {
	return errors.Join(
		syscall.Munmap(bcm.gpioMem8),
		syscall.Munmap(bcm.pwmMem8),
		syscall.Munmap(bcm.clkMem8),
		bcm.devmem.Close(),
		bcm.mbox.Close(),
	)
}

// Init initialises GPIOs and sets sane defaults
func (bcm *bcm2711bcm) Init() {
	bcm.InitGPIO()
	bcm.SetFanSpeed(bcm.opts.DefaultFanSpeed)
	bcm.SetStealthMode(bcm.opts.DefaultStealthModeEnabled)
}

// InitGPIO initalises GPIO configuration
func (bcm *bcm2711bcm) InitGPIO() {
	// based on https://datasheets.raspberrypi.com/bcm2711/bcm2711-peripherals.pdf
	bcm.wrMutex.Lock()
	defer bcm.wrMutex.Unlock()

	// Blade Butten (GPIO 20)
	// -> GPFSEL2 2:0, input
	bcm.gpioMem[GPFSEL2] = (bcm.gpioMem[GPFSEL2] &^ (0b111 << 0)) | (0b000 << 0)

	// Stealth Mode Output (GPIO 21)
	// -> GPFSEL2 5:3, output
	bcm.gpioMem[GPFSEL2] = (bcm.gpioMem[GPFSEL2] &^ (0b111 << 3)) | (0b001 << 3)

	// WS281x Output (GPIO 18)
	// -> GPFSEL1 24:26, regular output; it's configured as alt5 whenever pixel data is sent.
	// This is not performant but required as the pwm0 peripheral is shared between fan and data line for the LEDs.
	bcm.gpioMem[GPFSEL1] = (bcm.gpioMem[GPFSEL1] &^ (0b111 << 24)) | (0b000 << 24)

	// FAN PWM output for standard fan unit (GPIO 12)
	if bcm.opts.FanUnit == hal.FAN_UNIT_STANDARD {
		// -> GPFSEL1 8:6, alt0
		bcm.gpioMem[GPFSEL1] = (bcm.gpioMem[GPFSEL1] &^ (0b111 << 6)) | (0b100 << 6)
		bcm.setupFanPwm0()
	}

	// FAN TACH input for standard fan unit (GPIO 13)
	if bcm.opts.FanUnit == hal.FAN_UNIT_STANDARD {
		// -> GPFSEL1 11:9, input
		bcm.gpioMem[GPFSEL1] = (bcm.gpioMem[GPFSEL1] &^ (0b111 << 9)) | (0b000 << 9)
	}

	// FIXME add pullup

	// FIXME add WS2812 GPIO 18
}

func (bcm *bcm2711bcm) setupFanPwm0() {

	// Stop pwm for both channels; this is required to set the new configuration
	bcm.pwmMem[PWM_CTL] &^= 1<<8 | 1

	// Stop clock w/o any changes, they cannot be made in the same step
	bcm.clkMem[PWMCLK_CNTL] = bcm2711ClkManagerPwd | (bcm.clkMem[PWMCLK_CNTL] &^ (1 << 4))

	// Wait for the clock to not be busy so we can perform the changes
	for bcm.clkMem[PWMCLK_CNTL]&(1<<7) != 0 {
		time.Sleep(time.Microsecond * 20)
	}

	// passwd, mash, disabled, source (oscillator)
	bcm.clkMem[PWMCLK_CNTL] = bcm2711ClkManagerPwd | (0 << 9) | (0 << 4) | (1 << 0)

	// set PWM freq; the BCM2711 has an oscillator freq of 52 Mhz
	// Noctua fans are expecting a 25khz signal, where duty cycle controls fan on/speed/off
	// -> we'll need to get ~2.5Mhz of signal resultion in order to incorporate a 0-100 range
	// The following settings setup ~2.571Mhz resultion, resulting in a ~25,71khz signal
	// lying within the specifications of Noctua fans.
	bcm.clkMem[PWMCLK_DIV] = bcm2711ClkManagerPwd | (20 << 12) | (3276 << 0)

	// wait for changes to take effect before enabling it.
	// Note: 10us seems sufficient on idle systems, but doesn't always work when
	time.Sleep(time.Microsecond * 50)

	// Start clock (passwd, mash, enable, source)
	bcm.clkMem[PWMCLK_CNTL] = bcm2711ClkManagerPwd | (0 << 9) | (1 << 4) | (1 << 0)

	// Start pwm for both channels again
	bcm.pwmMem[PWM_CTL] &= 1<<8 | 1
}

// SetFanSpeed sets the fanspeed of a blade in percent (standard fan unit)
func (bcm *bcm2711bcm) SetFanSpeed(speed uint8) {
	bcm.setFanSpeedPWM(speed)
}

func (bcm *bcm2711bcm) setFanSpeedPWM(speed uint8) {
	bcm.wrMutex.Lock()
	defer bcm.wrMutex.Unlock()

	// set MSEN=0
	bcm.pwmMem[PWM_CTL] = bcm.pwmMem[PWM_CTL]&^(0xff) | (0 << 7) | (1 << 0)

	bcm.pwmMem[PWM_DAT1] = uint32(speed)
	bcm.pwmMem[PWM_RNG1] = 100
	time.Sleep(3 * time.Microsecond)

	// Store fan speed for later use
	bcm.currFanSpeed = speed
}

type LedColor struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// SetLEDs sets the color of the WS281x LEDs
func (bcm *bcm2711bcm) SetLEDs(top LedColor, edge LedColor) {
	bcm.wrMutex.Lock()
	defer bcm.wrMutex.Unlock()
	// Restore fan PWM after setting the LEDs & set the GPIO18 as a regular output
	defer func() {
		bcm.setupFanPwm0()
		bcm.setFanSpeedPWM(bcm.currFanSpeed)
	}()

	// Datarate for WS281x LEDs is 800kHz
	// Every bit transmitted takes 3 bits on of buffer (-../--.) thus we need (3*3*8) = 72 bits per LED and therefore 144 bits in total
	// ws2812 reset expects 55us of low signal, which is 132bits in the buffer (44 logocal bits) -> will take 55us to transmit -> reset signal
	// -> we need 144 + 132 = 276 bits in total, which when rounded up to the next multiple of 8 is 280 bits or 35 bytes
	const bufferSize = 35

	// Get DMA buffer

	// Stop pwm for both channels; this is required to set the new configuration
	bcm.pwmMem[PWM_CTL] &^= 1<<8 | 1

	// Set GPIO18 to alt5 (PWM0_0)
	bcm.gpioMem[GPFSEL1] = (bcm.gpioMem[GPFSEL1] &^ (0b111 << 24)) | (0b010 << 24)

	// Start pwm for both channels again
	bcm.pwmMem[PWM_CTL] &= 1<<8 | 1
}

func (bcm *bcm2711bcm) SetStealthMode(enable bool) {
	bcm.wrMutex.Lock()
	defer bcm.wrMutex.Unlock()

	if enable {
		// set high (bcm2711StealthPin == 21)
		bcm.gpioMem[7] = 1 << (bcm2711StealthPin)
	} else {
		// clear high state (bcm2711StealthPin == 21)
		bcm.gpioMem[10] = 1 << (bcm2711StealthPin)
	}
}
