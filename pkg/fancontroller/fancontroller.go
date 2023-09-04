package fancontroller

import (
	"fmt"
	"sync"
)

type FanController interface {
	Override(opts *FanOverrideOpts)
	GetFanSpeed(temperature float64) uint8
}

type FanOverrideOpts struct {
	Speed uint8
}

type FanControllerStep struct {
	// Temperature is the temperature to react to
	Temperature float64
	// Speed is the fan speed in percent
	Speed uint8
}

// FanController configures a fan controller for the computeblade
type FanControllerConfig struct {
	// Steps defines the temperature/speed steps for the fan controller
	Steps []FanControllerStep
}

// FanController is a simple fan controller that reacts to temperature changes with a linear function
type fanControllerLinear struct {
	mu       sync.Mutex
	overrideOpts *FanOverrideOpts
	config   FanControllerConfig
}

// NewFanControllerLinear creates a new FanControllerLinear
func NewLinearFanController(config FanControllerConfig) (FanController, error) {

	// Validate config for a very simple linear fan controller
	if len(config.Steps) != 2 {
		return nil, fmt.Errorf("exactly two steps must be defined")
	}
	if config.Steps[0].Temperature > config.Steps[1].Temperature {
		return nil, fmt.Errorf("step 1 temperature must be lower than step 2 temperature")
	}
	if config.Steps[0].Speed > config.Steps[1].Speed {
		return nil, fmt.Errorf("step 1 speed must be lower than step 2 speed")
	}
	if config.Steps[0].Speed > 100 || config.Steps[1].Speed > 100 {
		return nil, fmt.Errorf("speed must be between 0 and 100")
	}

	return &fanControllerLinear{
		config: config,
	}, nil
}

func (f *fanControllerLinear) Override(opts *FanOverrideOpts) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.overrideOpts = opts
}

// GetFanSpeed returns the fan speed in percent based on the current temperature
func (f *fanControllerLinear) GetFanSpeed(temperature float64) uint8 {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.overrideOpts != nil {
		return f.overrideOpts.Speed
	}

	if temperature <= f.config.Steps[0].Temperature {
		return f.config.Steps[0].Speed
	}
	if temperature >= f.config.Steps[1].Temperature {
		return f.config.Steps[1].Speed
	}

	// Calculate slope
	slope := float64(f.config.Steps[1].Speed-f.config.Steps[0].Speed) / (f.config.Steps[1].Temperature - f.config.Steps[0].Temperature)

	// Calculate speed
	speed := float64(f.config.Steps[0].Speed) + slope*(temperature-f.config.Steps[0].Temperature)

	return uint8(speed)
}
