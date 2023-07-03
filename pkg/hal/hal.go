package hal

type FanUnit uint8
type ComputeModule uint8

const (
	FAN_UNIT_STANDARD = iota
	FAN_UNIT_ADVANCED
)

type ComputeBladeHalOpts struct {
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
