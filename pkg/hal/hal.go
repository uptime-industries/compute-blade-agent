package hal

type FanUnit uint8
type ComputeModule uint8
type PowerStatus uint8

const (
	FanUnitStandard = iota
	FanUnitSmart
)

const (
	PoeOrUsbC = iota
	Poe802at
)

type LedColor struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

type ComputeBladeHalOptsDefault struct {
	StealthModeEnabled bool
	FanSpeed           uint8
	TopLedColor        LedColor
	EdgeLedColor       LedColor
}

type ComputeBladeHalOpts struct {
	FanUnit  FanUnit
	Defaults ComputeBladeHalOptsDefault
}

// COmputeBladeHal abstracts hardware details of the Compute Blade and provides a simple interface
type ComputeBladeHal interface {
	Init()
	Close() error
	SetFanSpeed(speed uint8)
	SetStealthMode(enabled bool)
	GetPowerStatus() PowerStatus
	SetLEDs(top, edge LedColor)
}
