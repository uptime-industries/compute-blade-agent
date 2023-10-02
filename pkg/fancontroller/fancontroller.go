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
	Percent uint8 `mapstructure:"speed"`
}

type FanControllerStep struct {
	// Temperature is the temperature to react to
	Temperature float64 `mapstructure:"temperature"`
	// Percent is the fan speed in percent
	Percent uint8 `mapstructure:"percent"`
}

// FanController configures a fan controller for the computeblade
type FanControllerConfig struct {
	// Steps defines the temperature/speed steps for the fan controller
	Steps []FanControllerStep `mapstructure:"steps"`
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
	if config.Steps[0].Percent > config.Steps[1].Percent {
		return nil, fmt.Errorf("step 1 speed must be lower than step 2 speed")
	}
	if config.Steps[0].Percent > 100 || config.Steps[1].Percent > 100 {
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
		return f.overrideOpts.Percent
	}

	if temperature <= f.config.Steps[0].Temperature {
		return f.config.Steps[0].Percent
	}
	if temperature >= f.config.Steps[1].Temperature {
		return f.config.Steps[1].Percent
	}

	// Calculate slope
	slope := float64(f.config.Steps[1].Percent-f.config.Steps[0].Percent) / (f.config.Steps[1].Temperature - f.config.Steps[0].Temperature)

	// Calculate speed
	speed := float64(f.config.Steps[0].Percent) + slope*(temperature-f.config.Steps[0].Temperature)

	return uint8(speed)
}
