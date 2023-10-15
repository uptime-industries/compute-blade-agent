package smartfanunit

import (
	"context"
	"io"

	"github.com/jacobsa/go-serial/serial"
	"github.com/xvzf/computeblade-agent/pkg/hal"
)

type FanUnitClient interface {
	// Run the client with event loop
	Run(context.Context) error

	// SetFanSpeedPercent sets the fan speed in percent.
	SetFanSpeedPercent(context.Context, uint8) error
	// SetLed sets the LED color.
	SetLed(context.Context, hal.LedColor) error

	// FanSpeedRPM returns the current fan speed in rotations per minute.
	FanSpeedRPM(context.Context) (uint8, error)
	// WaitForButtonPress blocks until the button is pressed.
	WaitForButtonPress(context.Context) error
	// AirFlowTemperature returns the temperature of the air flow.
	AirFlowTemperature(context.Context) (float32, error)
}

func NewFanUnitClient(portName string) (FanUnitClient, error) {
	// Open the serial port.
	_, err := serial.Open(serial.OpenOptions{
		PortName:          portName,
		BaudRate:          hal.SmartFanUnitBaudrate,
		DataBits:          8,
		StopBits:          1,
		MinimumReadSize:   1,
		RTSCTSFlowControl: false,
	})
	if err != nil {
		return nil, err
	}

	// return &fanUnitClient{rwc: rwc}, nil
	return nil, nil
}

type fanUnitClient struct {
	rwc io.ReadWriteCloser

}


// Run the client with event loop
func (fanunitclient *fanUnitClient) Run(_ context.Context) error {
	panic("not implemented") // TODO: Implement
}
// SetFanSpeedPercent sets the fan speed in percent.
func (fanunitclient *fanUnitClient) SetFanSpeedPercent(_ context.Context, _ uint8) error {
	panic("not implemented") // TODO: Implement
}
// SetLed sets the LED color.
func (fanunitclient *fanUnitClient) SetLed(_ context.Context, _ hal.LedColor) error {
	panic("not implemented") // TODO: Implement
}
// FanSpeedRPM returns the current fan speed in rotations per minute.
func (fanunitclient *fanUnitClient) FanSpeedRPM(_ context.Context) (uint8, error) {
	panic("not implemented") // TODO: Implement
}
// WaitForButtonPress blocks until the button is pressed.
func (fanunitclient *fanUnitClient) WaitForButtonPress(_ context.Context) error {
	panic("not implemented") // TODO: Implement
}
// AirFlowTemperature returns the temperature of the air flow.
func (fanunitclient *fanUnitClient) AirFlowTemperature(_ context.Context) (float32, error) {
	panic("not implemented") // TODO: Implement
}
