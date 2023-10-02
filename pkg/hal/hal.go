package hal

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type FanUnit uint8
type ComputeModule uint8
type PowerStatus uint8

var (
	fanSpeedTargetPercent = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "computeblade",
		Name:      "fan_speed_target_percent",
		Help:      "Target fanspeed in percent",
	})
	fanSpeed = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "computeblade",
		Name:      "fan_speed",
		Help:      "Fan speed in RPM",
	})
	socTemperature = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "computeblade",
		Name:      "soc_temperature",
		Help:      "SoC temperature in °C",
	})
	computeModule = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "computeblade",
		Name:      "compute_modul_present",
		Help:      "Compute module type",
	}, []string{"type"})
	ledColorChangeEventCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "computeblade",
		Name:      "led_color_change_event_count",
		Help:      "Led color change event_count",
	})
	powerStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "computeblade",
		Name:      "power_status",
		Help:      "Power status of the blade",
	}, []string{"type"})
	stealthModeEnabled = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "computeblade",
		Name:      "stealth_mode_enabled",
		Help:      "Stealth mode enabled",
	})
	fanUnit = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "computeblade",
		Name:      "fan_unit",
		Help:      "Fan unit",
	}, []string{"type"})
	edgeButtonEventCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "computeblade",
		Name:      "edge_button_event_count",
		Help:      "Number of edge button presses",
	})
)

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
