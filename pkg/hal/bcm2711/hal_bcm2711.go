package bcm2711

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"github.com/xvzf/computeblade-agent/pkg/hal"
)

const (
	bcm2711PeripheryBaseAddr = 0xFE000000
	bcm2711RegPwmAddr        = bcm2711PeripheryBaseAddr + 0x20C000
	bcm2711GpioAddr          = bcm2711PeripheryBaseAddr + 0x200000
	bcm2711ClkAddr           = bcm2711PeripheryBaseAddr + 0x101000
	bcm2711ClkManagerPwd     = (0x5A << 24) //(31 - 24) on CM_GP0CTL/CM_GP1CTL/CM_GP2CTL regs
	bcm2711PageSize          = 4096         // theoretical page size

	bcm2711FrontButtonPin = 20
	bcm2711StealthPin     = 21
	bcm2711RegPwmTachPin  = 13

	bcm2711RegGpfsel1 = 0x01

	bcm2711RegPwmCtl  = 0x00
	bcm2711RegPwmRng1 = 0x04
	bcm2711RegPwmFif1 = 0x06

	bcm2711RegPwmCtlBitPwen2 = 8 // Enable (pwm2)
	bcm2711RegPwmCtlBitClrf1 = 6 // Clear FIFO
	bcm2711RegPwmCtlBitUsef1 = 5 // Use FIFO
	bcm2711RegPwmCtlBitSbit1 = 3 // Line level when not transmitting
	bcm2711RegPwmCtlBitRptl1 = 2 // Repeat last data when FIFO is empty
	bcm2711RegPwmCtlBitMode1 = 1 // Mode; 0: PWM, 1: Serializer
	bcm2711RegPwmCtlBitPwen1 = 0 // Enable (pwm1)

	bcm2711RegPwmclkCntrl          = 0x28
	bcm2711RegPwmclkDiv            = 0x29
	bcm2711RegPwmclkCntrlBitSrcOsc = 0
	bcm2711RegPwmclkCntrlBitEnable = 4
)

type bcm2711 struct {
	// Config options
	opts hal.ComputeBladeHalOpts

	wrMutex sync.Mutex

	// Keep track of the currently set fanspeed so it can later be restored after setting the ws281x LEDs
	currFanSpeed uint8

	devmem    *os.File
	gpioMem8  []uint8
	gpioMem   []uint32
	pwmMem8   []uint8
	pwmMem    []uint32
	clkMem8   []uint8
	clkMem    []uint32
	gpioChip0 *gpiod.Chip

	// Stealth mode output
	stealthModeLine *gpiod.Line

	// Edge button input
	edgeButtonLine *gpiod.Line

	// PoE detection input
	poeLine *gpiod.Line

	// Fan tach input
	fanEdgeLine      *gpiod.Line
	lastFanEdgeEvent *gpiod.LineEvent
	fanRpm           float64
}

func New(opts hal.ComputeBladeHalOpts) (*bcm2711, error) {
	// /dev/gpiomem doesn't allow complex operations for PWM fan control or WS281x
	devmem, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, os.ModePerm)
	if err != nil {
		return nil, err
	}

	gpioChip0, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return nil, err
	}

	// Setup memory mappings
	gpioMem, gpioMem8, err := mmap(devmem, bcm2711GpioAddr, bcm2711PageSize)
	if err != nil {
		return nil, err
	}
	pwmMem, pwmMem8, err := mmap(devmem, bcm2711RegPwmAddr, bcm2711PageSize)
	if err != nil {
		return nil, err
	}
	clkMem, clkMem8, err := mmap(devmem, bcm2711ClkAddr, bcm2711PageSize)
	if err != nil {
		return nil, err
	}

	bcm := &bcm2711{
		devmem:    devmem,
		gpioMem:   gpioMem,
		gpioMem8:  gpioMem8,
		pwmMem:    pwmMem,
		pwmMem8:   pwmMem8,
		clkMem:    clkMem,
		clkMem8:   clkMem8,
		gpioChip0: gpioChip0,
		opts:      opts,
	}

	return bcm, bcm.setup()
}

// Close cleans all memory mappings
func (bcm *bcm2711) Close() error {
	errs := errors.Join(
		syscall.Munmap(bcm.gpioMem8),
		syscall.Munmap(bcm.pwmMem8),
		syscall.Munmap(bcm.clkMem8),
		bcm.devmem.Close(),
		bcm.gpioChip0.Close(),
		bcm.edgeButtonLine.Close(),
		bcm.poeLine.Close(),
		bcm.stealthModeLine.Close(),
	)

	if bcm.fanEdgeLine != nil {
		return errors.Join(errs, bcm.fanEdgeLine.Close())
	}
	return errs
}

// handleFanEdge handles an edge event on the fan tach input for the standard fan unite.
// Exponential moving average is used to smooth out the fan speed.
func (bcm *bcm2711) handleFanEdge(evt gpiod.LineEvent) {
	// Ensure we're always storing the last event
	defer func() {
		bcm.lastFanEdgeEvent = &evt
	}()

	// First event, we cannot extrapolate the fan speed yet
	if bcm.lastFanEdgeEvent == nil {
		return
	}

	// Calculate time delta between events
	delta := evt.Timestamp - bcm.lastFanEdgeEvent.Timestamp
	ticksPerSecond := 1000.0 / float64(delta.Milliseconds())
	rpm := (ticksPerSecond * 60.0) / 2.0 // 2 ticks per revolution

	// Simple moving average to smooth out the fan speed
	bcm.fanRpm = (rpm * 0.1) + (bcm.fanRpm * 0.9)
}

func (bcm *bcm2711) handleEdgeButtonEdge(evt gpiod.LineEvent) {
	fmt.Printf("button pressed\n")
}

// Init initialises GPIOs and sets sane defaults
func (bcm *bcm2711) setup() error {
	var err error = nil

	// Register edge event handler for edge button
	bcm.edgeButtonLine, err = bcm.gpioChip0.RequestLine(
		rpi.GPIO20, gpiod.WithEventHandler(bcm.handleEdgeButtonEdge),
		gpiod.WithFallingEdge, gpiod.WithPullUp, gpiod.WithDebounce(10*time.Millisecond))
	if err != nil {
		return err
	}

	// Register input for PoE detection
	bcm.poeLine, err = bcm.gpioChip0.RequestLine(rpi.GPIO23, gpiod.AsInput)
	if err != nil {
		return err
	}

	// Register output for stealth mode
	bcm.stealthModeLine, err = bcm.gpioChip0.RequestLine(rpi.GPIO21, gpiod.AsOutput(1))
	if err != nil {
		return err
	}

	// standard fan unit
	if bcm.opts.FanUnit == hal.FanUnitStandard {
		// FAN PWM output for standard fan unit (GPIO 12)
		// -> bcm2711RegGpfsel1 8:6, alt0
		bcm.gpioMem[bcm2711RegGpfsel1] = (bcm.gpioMem[bcm2711RegGpfsel1] &^ (0b111 << 6)) | (0b100 << 6)
		// Register edge event handler for fan tach input
		bcm.fanEdgeLine, err = bcm.gpioChip0.RequestLine(rpi.GPIO13, gpiod.WithEventHandler(bcm.handleFanEdge), gpiod.WithFallingEdge, gpiod.WithPullUp)
		if err != nil {
			return err
		}
	}

	// Setup defaults

	bcm.SetFanSpeed(bcm.opts.Defaults.FanSpeed)
	bcm.SetStealthMode(bcm.opts.Defaults.StealthModeEnabled)
	bcm.SetLEDs(bcm.opts.Defaults.TopLedColor, bcm.opts.Defaults.EdgeLedColor)

	return err
}

func (bcm2711 *bcm2711) GetFanSpeed() int {
	return int(bcm2711.fanRpm)
}

func (bcm *bcm2711) GetPowerStatus() hal.PowerStatus {

	// GPIO 23 is used for PoE detection
	if val, err := bcm.poeLine.Value(); err != nil && val == 1 {
		return hal.PowerPoe802at
	}
	return hal.PowerPoeOrUsbC
}

func (bcm *bcm2711) setPwm0Freq(targetFrequency uint64) error {
	// Calculate PWM divisor based on target frequency
	divisor := 54000000 / targetFrequency
	realDivisor := divisor & 0xfff // 12 bits
	if divisor != realDivisor {
		return fmt.Errorf("invalid frequency, max divisor is 4095, calculated divisor is %d", divisor)
	}

	// Stop pwm for both channels; this is required to set the new configuration
	bcm.pwmMem[bcm2711RegPwmCtl] &^= (1 << bcm2711RegPwmCtlBitPwen1) | (1 << bcm2711RegPwmCtlBitPwen2)
	time.Sleep(time.Microsecond * 10)

	// Stop clock w/o any changes, they cannot be made in the same step
	bcm.clkMem[bcm2711RegPwmclkCntrl] = bcm2711ClkManagerPwd | (bcm.clkMem[bcm2711RegPwmclkCntrl] &^ (1 << 4))
	time.Sleep(time.Microsecond * 10)

	// Wait for the clock to not be busy so we can perform the changes
	for bcm.clkMem[bcm2711RegPwmclkCntrl]&(1<<7) != 0 {
		time.Sleep(time.Microsecond * 10)
	}

	// passwd, disabled, source (oscillator)
	bcm.clkMem[bcm2711RegPwmclkCntrl] = bcm2711ClkManagerPwd | (0 << bcm2711RegPwmclkCntrlBitEnable) | (1 << bcm2711RegPwmclkCntrlBitSrcOsc)
	time.Sleep(time.Microsecond * 10)

	bcm.clkMem[bcm2711RegPwmclkDiv] = bcm2711ClkManagerPwd | (uint32(divisor) << 12)
	time.Sleep(time.Microsecond * 10)

	// Start clock (passwd, enable, source)
	bcm.clkMem[bcm2711RegPwmclkCntrl] = bcm2711ClkManagerPwd | (1 << bcm2711RegPwmclkCntrlBitEnable) | (1 << bcm2711RegPwmclkCntrlBitSrcOsc)
	time.Sleep(time.Microsecond * 10)

	// Start pwm for both channels again
	bcm.pwmMem[bcm2711RegPwmCtl] &= (1 << bcm2711RegPwmCtlBitPwen1)
	time.Sleep(time.Microsecond * 10)

	return nil
}

// SetFanSpeed sets the fanspeed of a blade in percent (standard fan unit)
func (bcm *bcm2711) SetFanSpeed(speed uint8) {
	bcm.setFanSpeedPWM(speed)
}

func (bcm *bcm2711) setFanSpeedPWM(speed uint8) {
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
	bcm.pwmMem[bcm2711RegPwmCtl] = (1 << bcm2711RegPwmCtlBitPwen1) | (1 << bcm2711RegPwmCtlBitMode1) | (1 << bcm2711RegPwmCtlBitRptl1) | (1 << bcm2711RegPwmCtlBitUsef1)
	time.Sleep(10 * time.Microsecond)
	bcm.pwmMem[bcm2711RegPwmRng1] = 32
	time.Sleep(10 * time.Microsecond)
	bcm.pwmMem[bcm2711RegPwmFif1] = targetvalue

	// Store fan speed for later use
	bcm.currFanSpeed = speed
}

func (bcm *bcm2711) SetStealthMode(enable bool) {
	if enable {
		_ = bcm.stealthModeLine.SetValue(1)
	} else {
		_ = bcm.stealthModeLine.SetValue(0)
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
func (bcm *bcm2711) SetLEDs(top hal.LedColor, edge hal.LedColor) {
	bcm.wrMutex.Lock()
	defer bcm.wrMutex.Unlock()

	// Set frequency to 3*800khz.
	// we'll bit-bang the data, so we'll need to send 3 bits per bit of data.
	bcm.setPwm0Freq(3 * 800000)
	time.Sleep(10 * time.Microsecond)

	// WS281x Output (GPIO 18)
	// -> bcm2711RegGpfsel1 24:26, regular output; it's configured as alt5 whenever pixel data is sent.
	// This is not optimal but required as the pwm0 peripheral is shared between fan and data line for the LEDs.
	time.Sleep(10 * time.Microsecond)
	bcm.gpioMem[bcm2711RegGpfsel1] = (bcm.gpioMem[bcm2711RegGpfsel1] &^ (0b111 << 24)) | (0b010 << 24)
	time.Sleep(10 * time.Microsecond)
	defer func() {
		// Set to regular output again so the PWM signal doesn't confuse the WS2812
		bcm.gpioMem[bcm2711RegGpfsel1] = (bcm.gpioMem[bcm2711RegGpfsel1] &^ (0b111 << 24)) | (0b001 << 24)
		bcm.setFanSpeedPWM(bcm.currFanSpeed)
	}()

	bcm.pwmMem[bcm2711RegPwmCtl] = (1 << bcm2711RegPwmCtlBitMode1) | (1 << bcm2711RegPwmCtlBitRptl1) | (0 << bcm2711RegPwmCtlBitSbit1) | (1 << bcm2711RegPwmCtlBitUsef1) | (1 << bcm2711RegPwmCtlBitClrf1)
	time.Sleep(10 * time.Microsecond)
	// bcm.pwmMem[bcm2711RegPwmRng1] = 32
	bcm.pwmMem[bcm2711RegPwmRng1] = 24 // we only need 24 bits per LED
	time.Sleep(10 * time.Microsecond)

	// Add sufficient padding to clear 50us of silence with ~412.5ns per bit -> at least 121 bits -> let's be safe and send 6*24=144 bits of silence
	bcm.pwmMem[bcm2711RegPwmFif1] = 0
	bcm.pwmMem[bcm2711RegPwmFif1] = 0
	bcm.pwmMem[bcm2711RegPwmFif1] = 0
	bcm.pwmMem[bcm2711RegPwmFif1] = 0
	bcm.pwmMem[bcm2711RegPwmFif1] = 0
	bcm.pwmMem[bcm2711RegPwmFif1] = 0
	// Write top LED data
	bcm.pwmMem[bcm2711RegPwmFif1] = serializePwmDataFrame(top.Red) << 8
	bcm.pwmMem[bcm2711RegPwmFif1] = serializePwmDataFrame(top.Green) << 8
	bcm.pwmMem[bcm2711RegPwmFif1] = serializePwmDataFrame(top.Blue) << 8
	// Write edge LED data
	bcm.pwmMem[bcm2711RegPwmFif1] = serializePwmDataFrame(edge.Red) << 8
	bcm.pwmMem[bcm2711RegPwmFif1] = serializePwmDataFrame(edge.Green) << 8
	bcm.pwmMem[bcm2711RegPwmFif1] = serializePwmDataFrame(edge.Blue) << 8
	// make sure there's >50us of silence
	bcm.pwmMem[bcm2711RegPwmFif1] = 0 // auto-repeated, so no need to feed the FIFO further.

	bcm.pwmMem[bcm2711RegPwmCtl] = (1 << bcm2711RegPwmCtlBitPwen1) | (1 << bcm2711RegPwmCtlBitMode1) | (1 << bcm2711RegPwmCtlBitRptl1) | (0 << bcm2711RegPwmCtlBitSbit1) | (1 << bcm2711RegPwmCtlBitUsef1)
	// sleep for 4*50us to ensure the data is sent. This is probably a bit too gracious but does not have a significant impact, so let's be safe data gets out.
	time.Sleep(200 * time.Microsecond)
}
