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
	PWM_STA  = 0x01
	PWM_DMAC = 0x02
	PWM_RNG1 = 0x04
	PWM_DAT1 = 0x05
	PWM_FIF1 = 0x06

	PWM_CTL_PWEN2 = 8 // Enable (pwm2)
	PWM_CTL_CLRF1 = 6 // Clear FIFO
	PWM_CTL_MSEN1 = 7 // Use M/S algorithm
	PWM_CTL_USEF1 = 5 // Use FIFO
	PWM_CTL_POLA1 = 4 // Invert polarity
	PWM_CTL_SBIT1 = 3 // Line level when not transmitting
	PWM_CTL_RPTL1 = 2 // Repeat last data when FIFO is empty
	PWM_CTL_MODE1 = 1 // Mode; 0: PWM, 1: Serializer
	PWM_CTL_PWEN1 = 0 // Enable (pwm1)

	PWM_STA_STA1  = 9 // Status
	PWM_STA_BERR  = 8 // Bus Error
	PWM_STA_GAPO1 = 4 // Gap detected
	PWM_STA_RERR1 = 3 // FIFO Read Error
	PWM_STA_WERR1 = 2 // FIFO Write Error
	PWM_STA_EMPT1 = 1 // FIFO Empty
	PWM_STA_FULL1 = 0 // FIFO Full

	PWMCLK_CNTL         = 0x28
	PWMCLK_CNTL_SRC_OSC = 0
	PWMCLK_CNTL_ENABLE  = 4
	PWMCLK_DIV          = 0x29
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
	// bcm.SetFanSpeed(bcm.opts.DefaultFanSpeed)
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

	// FAN PWM output for standard fan unit (GPIO 12)
	if bcm.opts.FanUnit == hal.FAN_UNIT_STANDARD {
		// -> GPFSEL1 8:6, alt0
		bcm.gpioMem[GPFSEL1] = (bcm.gpioMem[GPFSEL1] &^ (0b111 << 6)) | (0b100 << 6)
		bcm.setFanSpeedPWM(bcm.opts.DefaultFanSpeed)
	}

	// FAN TACH input for standard fan unit (GPIO 13)
	if bcm.opts.FanUnit == hal.FAN_UNIT_STANDARD {
		// -> GPFSEL1 11:9, input
		bcm.gpioMem[GPFSEL1] = (bcm.gpioMem[GPFSEL1] &^ (0b111 << 9)) | (0b000 << 9)
	}

	// Set WS2812 output (GPIO 18)
	// -> GPFSEL1 24:26, set as regular output by default. On-demand, it's mapped to pwm0
	bcm.gpioMem[GPFSEL1] = (bcm.gpioMem[GPFSEL1] &^ (0b111 << 24)) | (0b000 << 24)

	// FIXME add edge button
}

func (bcm *bcm2711bcm) setPwm0Freq(targetFrequency uint64) error {
	// Calculate PWM divisor based on target frequency
	divisor := 54000000 / targetFrequency
	realDivisor := divisor & 0xfff // 12 bits
	if divisor != realDivisor {
		return errors.New("invalid frequency, max divisor is 4095, calculated divisor is " + string(divisor))
	}

	// Stop pwm for both channels; this is required to set the new configuration
	bcm.pwmMem[PWM_CTL] &^= (1 << PWM_CTL_PWEN1) | (1 << PWM_CTL_PWEN2)
	time.Sleep(time.Microsecond * 10)

	// Stop clock w/o any changes, they cannot be made in the same step
	bcm.clkMem[PWMCLK_CNTL] = bcm2711ClkManagerPwd | (bcm.clkMem[PWMCLK_CNTL] &^ (1 << 4))
	time.Sleep(time.Microsecond * 10)

	// Wait for the clock to not be busy so we can perform the changes
	for bcm.clkMem[PWMCLK_CNTL]&(1<<7) != 0 {
		time.Sleep(time.Microsecond * 10)
	}

	// passwd, disabled, source (oscillator)
	bcm.clkMem[PWMCLK_CNTL] = bcm2711ClkManagerPwd | (0 << PWMCLK_CNTL_ENABLE) | (1 << PWMCLK_CNTL_SRC_OSC)
	time.Sleep(time.Microsecond * 10)

	bcm.clkMem[PWMCLK_DIV] = bcm2711ClkManagerPwd | (uint32(divisor) << 12)
	time.Sleep(time.Microsecond * 10)

	// Start clock (passwd, enable, source)
	bcm.clkMem[PWMCLK_CNTL] = bcm2711ClkManagerPwd | (1 << PWMCLK_CNTL_ENABLE) | (1 << PWMCLK_CNTL_SRC_OSC)
	time.Sleep(time.Microsecond * 10)

	// Start pwm for both channels again
	bcm.pwmMem[PWM_CTL] &= (1 << PWM_CTL_PWEN1)
	time.Sleep(time.Microsecond * 10)

	return nil
}

// SetFanSpeed sets the fanspeed of a blade in percent (standard fan unit)
func (bcm *bcm2711bcm) SetFanSpeed(speed uint8) {
	bcm.setFanSpeedPWM(speed)
}

func (bcm *bcm2711bcm) setFanSpeedPWM(speed uint8) {
	// Noctua fans are expecting a 25khz signal, where duty cycle controls fan on/speed/off
	// With the usage of the FIFO, we can alter the duty cycle by the number of bits set in the FIFO, maximum of 32.
	// We therefore need a frequency of 32*25khz = 800khz, which is a divisor of 67.5 (thus we'll use 68).
	// This results in an actual period frequency of 24.8khz, which is within the specifications of Noctua fans.
	err := bcm.setPwm0Freq(800000)
	if err != nil {
		// we know it produces a valid divisor, so this should never happen
		panic(err)
	}

	// Using hardware ticks would offer a better resultion, but this works for now.
	var targetvalue uint32 = 0
	if speed == 0 {
		targetvalue = 0
	} else if speed <= 100 {
		for i := 0; i <= int((float64(speed)/100.0)*32.0); i++ {
			targetvalue |= (1 << i)
		}
	} else {
		targetvalue = ^(uint32(0))
	}

	// Use fifo, repeat, ...
	bcm.pwmMem[PWM_CTL] = (1 << PWM_CTL_PWEN1) | (1 << PWM_CTL_MODE1) | (1 << PWM_CTL_RPTL1) | (1 << PWM_CTL_USEF1)
	time.Sleep(10 * time.Microsecond)
	bcm.pwmMem[PWM_RNG1] = 32
	time.Sleep(10 * time.Microsecond)
	bcm.pwmMem[PWM_FIF1] = targetvalue

	// Store fan speed for later use
	bcm.currFanSpeed = speed
}

type LedColor struct {
	Red   uint8
	Green uint8
	Blue  uint8
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

// serializePwmDataFrame converts a byte to a 24 bit PWM data frame for WS281x LEDs
func serializePwmDataFrame(data uint8) uint32 {
	var result uint32 = 0
	for i := 7; i >= 0; i-- {
		if i != 7 {
			result <<= 3
		}
		if (uint32(data)&(1<<i))>>i == 0 {
			result |= 0b100 // -__
		} else {
			result |= 0b110 // --_
		}
	}
	return result
}

// SetLEDs sets the color of the WS281x LEDs
func (bcm *bcm2711bcm) SetLEDs(top LedColor, edge LedColor) {
	bcm.wrMutex.Lock()
	defer bcm.wrMutex.Unlock()

	// Set frequency to 3*800khz.
	// we'll bit-bang the data, so we'll need to send 3 bits per bit of data.
	bcm.setPwm0Freq(3 * 800000)
	time.Sleep(10 * time.Microsecond)

	// WS281x Output (GPIO 18)
	// -> GPFSEL1 24:26, regular output; it's configured as alt5 whenever pixel data is sent.
	// This is not optimal but required as the pwm0 peripheral is shared between fan and data line for the LEDs.
	time.Sleep(10 * time.Microsecond)
	bcm.gpioMem[GPFSEL1] = (bcm.gpioMem[GPFSEL1] &^ (0b111 << 24)) | (0b010 << 24)
	time.Sleep(10 * time.Microsecond)
	defer func() {
		// Set to regular output again so the PWM signal doesn't confuse the WS2812
		bcm.gpioMem[GPFSEL1] = (bcm.gpioMem[GPFSEL1] &^ (0b111 << 24)) | (0b000 << 24)
		bcm.setFanSpeedPWM(bcm.currFanSpeed)
	}()

	bcm.pwmMem[PWM_CTL] = (1 << PWM_CTL_MODE1) | (1 << PWM_CTL_RPTL1) | (0 << PWM_CTL_SBIT1) | (1 << PWM_CTL_USEF1) | (1 << PWM_CTL_CLRF1)
	time.Sleep(10 * time.Microsecond)
	bcm.pwmMem[PWM_RNG1] = 32
	time.Sleep(10 * time.Microsecond)

	// Add sufficient padding to clear
	bcm.pwmMem[PWM_FIF1] = 0
	bcm.pwmMem[PWM_FIF1] = 0
	bcm.pwmMem[PWM_FIF1] = 0
	// Write top LED data
	bcm.pwmMem[PWM_FIF1] = serializePwmDataFrame(top.Red)
	bcm.pwmMem[PWM_FIF1] = serializePwmDataFrame(top.Green)
	bcm.pwmMem[PWM_FIF1] = serializePwmDataFrame(top.Blue)
	// Write edge LED data
	bcm.pwmMem[PWM_FIF1] = serializePwmDataFrame(edge.Red)
	bcm.pwmMem[PWM_FIF1] = serializePwmDataFrame(edge.Green)
	bcm.pwmMem[PWM_FIF1] = serializePwmDataFrame(edge.Blue)
	// make sure there's >50us of silence
	bcm.pwmMem[PWM_FIF1] = 0
	bcm.pwmMem[PWM_FIF1] = 0
	bcm.pwmMem[PWM_FIF1] = 0

	bcm.pwmMem[PWM_CTL] = (1 << PWM_CTL_PWEN1) | (1 << PWM_CTL_MODE1) | (1 << PWM_CTL_RPTL1) | (0 << PWM_CTL_SBIT1) | (1 << PWM_CTL_USEF1)
	// sleep for 4*50us to ensure the data is sent.
	time.Sleep(200 * time.Microsecond)
}
