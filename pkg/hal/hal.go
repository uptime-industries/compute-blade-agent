package hal

import "fmt"

type FanUnit uint8
type ComputeModule uint8

const (
	FAN_UNIT_STANDARD = iota
	FAN_UNIT_ADVANCED
)

const (
	COMPUTE_MODULE_TYPE_CM4 ComputeModule = iota
)

type ComputeBladeHalOpts struct {
	ComputeModuleType         ComputeModule
	FanUnit                   FanUnit
	DefaultFanSpeed           uint8
	DefaultStealthModeEnabled bool
}

// COmputeBladeHal abstracts hardware details of the Compute Blade and provides a simple interface
type ComputeBladeHal interface {
	Init()
	Close() error
	SetFanSpeed(speed uint8)
	SetStealthMode(enabled bool)
}

// NewComputeBladeHAL returns a new HAL for the Compute Blade and a given configuration
func NewComputeBladeHAL(opts ComputeBladeHalOpts) (ComputeBladeHal, error) {
	switch opts.ComputeModuleType {
	case COMPUTE_MODULE_TYPE_CM4:
		return NewBcm2711Hal(opts)
	default:
		return nil, fmt.Errorf("unsupported compute module type: %d", opts.ComputeModuleType)
	}
}
