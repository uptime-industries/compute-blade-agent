package smartfanunit

import (
	"errors"
	"runtime/internal/atomic"

	"github.com/xvzf/computeblade-agent/pkg/hal"
	"github.com/xvzf/computeblade-agent/pkg/smartfanunit/proto"
)

const (
	// Blade -> FanUnit
	CmdSetFanSpeedPercent proto.Command = 0x01
	CmdSetLED             proto.Command = 0x02

	// FanUnit -> Blade, sent in regular intervals
	NotifyButtonPress        proto.Command = 0xa1
	NotifyAirFlowTemperature proto.Command = 0xa2
	NotifyFanSpeedRPM        proto.Command = 0xa3
)

var ErrInvalidCommand = errors.New("invalid command")


type PacketGenerator interface {
	Packet() proto.Packet
}

// SetFanSpeedPercentPacket is sent from the blade to the fan unit to set the fan speed in percent.
type SetFanSpeedPercentPacket struct {
	Percent uint8
}

func (p *SetFanSpeedPercentPacket) Packet() proto.Packet {
	return proto.Packet{
		Command: CmdSetFanSpeedPercent,
		Data:    proto.Data{p.Percent, 0, 0},
	}
}

func (p *SetFanSpeedPercentPacket) FromPacket(packet proto.Packet) error {
	if packet.Command != CmdSetFanSpeedPercent {
		return ErrInvalidCommand
	}
	p.Percent = packet.Data[0]
	return nil
}

// SetLEDPacket is sent from the blade to the fan unit to set the LED color.
type SetLEDPacket struct {
	Color hal.LedColor
}

func (p *SetLEDPacket) Packet() proto.Packet {
	return proto.Packet{
		Command: CmdSetLED,
		Data:    proto.Data{p.Color.Blue, p.Color.Green, p.Color.Red},
	}
}

func (p *SetLEDPacket) FromPacket(packet proto.Packet) error {
	if packet.Command != CmdSetLED {
		return ErrInvalidCommand
	}
	p.Color = hal.LedColor{
		Blue:  packet.Data[0],
		Green: packet.Data[1],
		Red:   packet.Data[2],
	}
	return nil
}

// ButtonPressPacket is sent from the fan unit to the blade when the button is pressed.
type ButtonPressPacket struct{}

func (p *ButtonPressPacket) Packet() proto.Packet {
	return proto.Packet{
		Command: NotifyButtonPress,
		Data:    proto.Data{},
	}
}

func (p *ButtonPressPacket) FromPacket(packet proto.Packet) error {
	if packet.Command != NotifyButtonPress {
		return ErrInvalidCommand
	}
	return nil
}

// AirFlowTemperaturePacket is sent from the fan unit to the blade to report the current air flow temperature.
type AirFlowTemperaturePacket struct {
	Temperature float32
}

func (p *AirFlowTemperaturePacket) Packet() proto.Packet {
	return proto.Packet{
		Command: NotifyAirFlowTemperature,
		Data: proto.Data(float32To24Bit(p.Temperature)),
	}
}

func (p *AirFlowTemperaturePacket) FromPacket(packet proto.Packet) error {
	if packet.Command != NotifyAirFlowTemperature {
		return ErrInvalidCommand
	}
	p.Temperature = float32From24Bit(packet.Data)
	return nil
}


// FanSpeedRPMPacket is sent from the fan unit to the blade to report the current fan speed in RPM.
type FanSpeedRPMPacket struct {
	RPM float32
}
func (p *FanSpeedRPMPacket) Packet() proto.Packet {
	return proto.Packet{
		Command: NotifyFanSpeedRPM,
		Data: float32To24Bit(p.RPM),
	}
}
func (p *FanSpeedRPMPacket) FromPacket(packet proto.Packet) error {
	if packet.Command != NotifyFanSpeedRPM {
		return ErrInvalidCommand
	}
	p.RPM = float32From24Bit(packet.Data)
	return nil
}

