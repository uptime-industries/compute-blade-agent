//go:build darwin

package hal

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// fails if SimulatedHal does not implement ComputeBladeHal
var _ ComputeBladeHal = &SimulatedHal{}

// ComputeBladeMock implements a mock for the ComputeBladeHal interface
type SimulatedHal struct {
	logger *zap.Logger
}

func NewCm4Hal(opts ComputeBladeHalOpts) (ComputeBladeHal, error) {
	logger := zap.L().Named("hal").Named("simulated-cm4")
	logger.Warn("Using simulated hal")

	return &SimulatedHal{
		logger: logger,
	}, nil
}

func (m *SimulatedHal) Close() error {
	return nil
}

func (m *SimulatedHal) SetFanSpeed(percent uint8) error {
	m.logger.Info("SetFanSpeed", zap.Uint8("percent", percent))
	return nil
}

func (m *SimulatedHal) GetFanRPM() (float64, error) {
	return 1337, nil
}

func (m *SimulatedHal) SetStealthMode(enabled bool) error {
	m.logger.Info("SetStealthMode", zap.Bool("enabled", enabled))
	return nil
}

func (m *SimulatedHal) GetPowerStatus() (PowerStatus, error) {
	m.logger.Info("GetPowerStatus")
	return PowerPoe802at, nil
}

func (m *SimulatedHal) WaitForEdgeButtonPress(ctx context.Context) error {
	m.logger.Info("WaitForEdgeButtonPress")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return nil
	}
}

func (m *SimulatedHal) SetLed(idx uint, color LedColor) error {
	m.logger.Info("SetLed", zap.Uint("idx", idx), zap.Any("color", color))
	return nil
}
