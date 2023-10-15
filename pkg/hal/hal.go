package hal

import (
	"context"
)

type FanUnit uint8
type ComputeModule uint8
type PowerStatus uint8

func (p PowerStatus) String() string {
	switch p {
	case PowerPoe802at:
		return "poe+"
	case PowerPoeOrUsbC:
		return "poeOrUsbC"
	default:
		return "undefined"
	}
}

const (
	FanUnitStandard = iota
	FanUnitSmart
)

const (
	PowerPoeOrUsbC = iota
	PowerPoe802at
)

const (
	SmartFanUnitBaudrate = 115200
)


const (
	LedTop = iota
	LedEdge
)

type LedColor struct {
	Red   uint8 `mapstructure:"red"`
	Green uint8 `mapstructure:"green"`
	Blue  uint8 `mapstructure:"blue"`
}

type ComputeBladeHalOpts struct {
	FanUnit FanUnit
}

// ComputeBladeHal abstracts hardware details of the Compute Blade and provides a simple interface
type ComputeBladeHal interface {
	Close() error
	// SetFanSpeed sets the fan speed in percent
	SetFanSpeed(speed uint8) error
	// GetFanSpeed returns the current fan speed in percent (based on moving average)
	GetFanRPM() (float64, error)
	// SetStealthMode enables/disables stealth mode of the blade (turning on/off the LEDs)
	SetStealthMode(enabled bool) error
	// SetLEDs sets the color of the LEDs
	SetLed(idx uint, color LedColor) error
	// GetPowerStatus returns the current power status of the blade
	GetPowerStatus() (PowerStatus, error)
	// GetTemperature returns the current temperature of the SoC in °C
	GetTemperature() (float64, error)
	// GetEdgeButtonPressChan returns a channel emitting edge button press events
	WaitForEdgeButtonPress(ctx context.Context) error
}
