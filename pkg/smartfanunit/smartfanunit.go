package smartfanunit

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/jacobsa/go-serial/serial"
	"github.com/xvzf/computeblade-agent/pkg/eventbus"
	"github.com/xvzf/computeblade-agent/pkg/hal"
	"github.com/xvzf/computeblade-agent/pkg/smartfanunit/proto"
	"golang.org/x/sync/errgroup"
)

type FanUnitClient interface {
	// Run the client with event loop
	Run(context.Context) error

	// SetFanSpeedPercent sets the fan speed in percent.
	SetFanSpeedPercent(context.Context, uint8) error
	// SetLed sets the LED color.
	SetLed(context.Context, hal.LedColor) error

	// FanSpeedRPM returns the current fan speed in rotations per minute.
	FanSpeedRPM(context.Context) (float32, error)
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

var ErrCommunicationFailed = errors.New("communication failed")

const (
	inboundTopic  = "smartfanunit:inbound"
	outboundTopic = "smartfanunit:outbound"
)

// FIXME move to internal & log err
type fanUnitClient struct {
	rwc io.ReadWriteCloser
	mu  sync.Mutex // write mutex

	speed FanSpeedRPMPacket
	airflow AirFlowTemperaturePacket

	eb eventbus.EventBus
}

func matchCmd(cmd proto.Command) func(any) bool {
	return func(pktAny any) bool {
		pkt, ok := pktAny.(proto.Packet)
		if !ok {
			return false
		}
		if pkt.Command == cmd {
			return true
		}
		return false
	}
}

// Run the client with event loop
func (fuc *fanUnitClient) Run(parentCtx context.Context) error {

	ctx, cancel := context.WithCancelCause(parentCtx)
	defer cancel(nil)

	wg := errgroup.Group{}

	// Start read loop
	wg.Go(func() error {
		errCounter := 0
		for {
			pkt, err := proto.ReadPacket(ctx, fuc.rwc)
			if err != nil {
				errCounter++
				continue
			}
			fuc.eb.Publish(inboundTopic, pkt)
		}
	})

	// Subscribe to fan speed updates
	wg.Go(func() error {
		sub := fuc.eb.Subscribe(inboundTopic, matchCmd(NotifyFanSpeedRPM))
		defer sub.Unsubscribe()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case pktAny := <-sub.C():
				rawPkt := pktAny.(proto.Packet)
				if err := fuc.speed.FromPacket(rawPkt); err != nil {
					return err
				}
			}
		}
	})


	// Subscribe to air flow temperature updates
	wg.Go(func() error {
		sub := fuc.eb.Subscribe(inboundTopic, matchCmd(NotifyAirFlowTemperature))
		defer sub.Unsubscribe()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case pktAny := <-sub.C():
				rawPkt := pktAny.(proto.Packet)
				if err := fuc.airflow.FromPacket(rawPkt); err != nil {
					cancel(err)
					return err
				}
			}
		}
	})

	return wg.Wait()
}

func (fuc *fanUnitClient) write(ctx context.Context, pktGen PacketGenerator) error {
	fuc.mu.Lock()
	defer fuc.mu.Unlock()
	return proto.WritePacket(ctx, fuc.rwc, pktGen.Packet())
}

// SetFanSpeedPercent sets the fan speed in percent.
func (fuc *fanUnitClient) SetFanSpeedPercent(ctx context.Context, percent uint8) error {
	return fuc.write(ctx, &SetFanSpeedPercentPacket{Percent: percent})
}

// SetLed sets the LED color.
func (fuc *fanUnitClient) SetLed(ctx context.Context, color hal.LedColor) error {
	return fuc.write(ctx, &SetLEDPacket{Color: color})
}

// FanSpeedRPM returns the current fan speed in rotations per minute.
func (fuc *fanUnitClient) FanSpeedRPM(_ context.Context) (float32, error) {
	return fuc.speed.RPM, nil
}

// WaitForButtonPress blocks until the button is pressed.
func (fuc *fanUnitClient) WaitForButtonPress(ctx context.Context) error {

	sub := fuc.eb.Subscribe(inboundTopic, matchCmd(NotifyButtonPress))
	defer sub.Unsubscribe()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case pktAny := <-sub.C():
		rawPkt := pktAny.(proto.Packet)
		if rawPkt.Command != NotifyButtonPress {
			return errors.New("unexpected packet")
		}
	}

	return nil
}

// AirFlowTemperature returns the temperature of the air flow.
func (fuc *fanUnitClient) AirFlowTemperature(_ context.Context) (float32, error) {
	return fuc.airflow.Temperature, nil
}
