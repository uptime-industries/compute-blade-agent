package hal

import (
	"errors"
	"os"
	"sync"
	"syscall"
	"time"
)

const (
	bcm2711PeripheryBaseAddr = 0xFE000000
	bcm2711PwmAddr           = bcm2711PeripheryBaseAddr + 0x20C000
	bcm2711GpioAddr          = bcm2711PeripheryBaseAddr + 0x200000
	bcm2711ClkAddr           = bcm2711PeripheryBaseAddr + 0x101000
	bcm2711ClkManagerPwd     = (0x5A << 24) //(31 - 24) on CM_GP0CTL/CM_GP1CTL/CM_GP2CTL regs

	bcm2711FrontButtonPin = 20
	bcm2711StealthPin     = 21
	bcm2711PwmFanPin      = 12
	bcm2711PwmTachPin     = 13
)

type bcm2711hal struct {
	// Config options
	opts ComputeBladeHalOpts

	wrMutex sync.Mutex

	devmem   *os.File
	gpioMem8 []uint8
	gpioMem  []uint32
	pwmMem8  []uint8
	pwmMem   []uint32
	clkMem8  []uint8
	clkMem   []uint32
}

func NewBcm2711Hal(opts ComputeBladeHalOpts) (*bcm2711hal, error) {
	// /dev/gpiomem doesn't allow complex operations for PWM fan control or WS281x
	devmem, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Setup memory mappings
	gpioMem, gpioMem8, err := mmap(devmem, bcm2711GpioAddr, 4096)
	if err != nil {
		return nil, err
	}
	pwmMem, pwmMem8, err := mmap(devmem, bcm2711PwmAddr, 4096)
	if err != nil {
		return nil, err
	}
	clkMem, clkMem8, err := mmap(devmem, bcm2711ClkAddr, 4096)
	if err != nil {
		return nil, err
	}

	return &bcm2711hal{
		devmem:   devmem,
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
func (hal *bcm2711hal) Close() error {
	return errors.Join(
		syscall.Munmap(hal.gpioMem8),
		syscall.Munmap(hal.pwmMem8),
		syscall.Munmap(hal.clkMem8),
		hal.devmem.Close(),
	)
}

// Init initialises GPIOs and sets sane defaults
func (hal *bcm2711hal) Init() {
	hal.InitGPIO()
	hal.SetFanSpeed(hal.opts.DefaultFanSpeed)
	hal.SetStealthMode(hal.opts.DefaultStealthModeEnabled)
}

// InitGPIO initalises GPIO configuration
func (hal *bcm2711hal) InitGPIO() {
	// based on https://datasheets.raspberrypi.com/bcm2711/bcm2711-peripherals.pdf
	hal.wrMutex.Lock()
	defer hal.wrMutex.Unlock()

	// Blade Butten (GPIO 20)
	// -> GPFSEL2 2:0, input
	hal.gpioMem[2] = (hal.gpioMem[2] &^ (0b111 << 0)) | (0b000 << 0)

	// Stealth Mode Output (GPIO 21)
	// -> GPFSEL2 5:3, output
	hal.gpioMem[2] = (hal.gpioMem[2] &^ (0b111 << 3)) | (0b001 << 3)

	// FAN PWM output for standard fan unit (GPIO 12)
	if hal.opts.FanUnit == FAN_UNIT_STANDARD {
		// -> GPFSEL1 8:6, alt0
		hal.gpioMem[1] = (hal.gpioMem[1] &^ (0b111 << 6)) | (0b100 << 6)

		// Stop pwm for both channels; this is required to set the new configuration
		hal.pwmMem[0] &^= 1<<8 | 1

		// Stop clock w/o any changes, they cannot be made in the same step
		hal.clkMem[40] = bcm2711ClkManagerPwd | (hal.clkMem[40] &^ (1 << 4))

		// Wait for the clock to not be busy so we can perform the changes
		for hal.clkMem[40]&(1<<7) != 0 {
			time.Sleep(time.Microsecond * 20)
		}

		// passwd, mash, disabled, source (oscillator)
		hal.clkMem[40] = bcm2711ClkManagerPwd | (0 << 9) | (0 << 4) | (1 << 0)

		// set PWM freq; the BCM2711 has an oscillator freq of 52 Mhz
		// Noctua fans are expecting a 25khz signal, where duty cycle controls fan on/speed/off
		// -> we'll need to get ~2.5Mhz of signal resultion in order to incorporate a 0-100 range
		// The following settings setup ~2.571Mhz resultion, resulting in a ~25,71khz signal
		// lying within the specifications of Noctua fans.
		hal.clkMem[41] = bcm2711ClkManagerPwd | (20 << 12) | (3276 << 0)

		// wait for changes to take effect before enabling it.
		// Note: 10us seems sufficient on idle systems, but doesn't always work when
		time.Sleep(time.Microsecond * 50)

		// Start clock (passwd, mash, enable, source)
		hal.clkMem[40] = bcm2711ClkManagerPwd | (0 << 9) | (1 << 4) | (1 << 0)

		// Start pwm for both channels again
		hal.pwmMem[0] &= 1<<8 | 1
	}

	// FAN TACH input for standard fan unit (GPIO 13)
	if hal.opts.FanUnit == FAN_UNIT_STANDARD {
		// -> GPFSEL1 11:9, input
		hal.gpioMem[1] = (hal.gpioMem[1] &^ (0b111 << 9)) | (0b000 << 9)
	}

	// FIXME add pullup

	// FIXME add WS2812 GPIO 18
}

// SetFanSpeed sets the fanspeed of a blade in percent (standard fan unit)
func (hal *bcm2711hal) SetFanSpeed(speed uint8) {
	hal.setFanSpeedPWM(speed)
}

func (hal *bcm2711hal) setFanSpeedPWM(speed uint8) {
	hal.wrMutex.Lock()
	defer hal.wrMutex.Unlock()

	// set MSEN=0
	hal.pwmMem[0] = hal.pwmMem[0]&^(0xff) | (0 << 7) | (1 << 0)

	hal.pwmMem[5] = uint32(speed)
	hal.pwmMem[4] = 100
	time.Sleep(3 * time.Microsecond)
}

func (hal *bcm2711hal) SetStealthMode(enable bool) {
	hal.wrMutex.Lock()
	defer hal.wrMutex.Unlock()

	if enable {
		// set high (bcm2711StealthPin == 21)
		hal.gpioMem[7] = 1 << (bcm2711StealthPin)
	} else {
		// clear high state (bcm2711StealthPin == 21)
		hal.gpioMem[10] = 1 << (bcm2711StealthPin)
	}
}
