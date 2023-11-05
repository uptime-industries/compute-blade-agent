//go:build tinygo

package main

import (
	"context"
	"machine"
	"time"

	"github.com/xvzf/computeblade-agent/pkg/eventbus"
	"github.com/xvzf/computeblade-agent/pkg/hal/led"
	"github.com/xvzf/computeblade-agent/pkg/smartfanunit"
	"github.com/xvzf/computeblade-agent/pkg/smartfanunit/emc2101"
	"github.com/xvzf/computeblade-agent/pkg/smartfanunit/proto"
	"golang.org/x/sync/errgroup"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/ws2812"
)

const (
	leftBladeTopicIn   = "left:in"
	leftBladeTopicOut  = "left:out"
	rightBladeTopicIn  = "right:in"
	rightBladeTopicOut = "right:out"
)

type Controller struct {
	DefaultFanSpeed uint8
	LEDs            ws2812.Device
	FanController   emc2101.EMC2101
	ButtonPin       machine.Pin

	LeftUART  drivers.UART
	RightUART drivers.UART

	eb               eventbus.EventBus
	leftLed          led.Color
	rightLed         led.Color
	leftReqFanSpeed  uint8
	rightReqFanSpeed uint8

	buttonPressed bool
}

func (c *Controller) Run(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	c.eb = eventbus.New()

	c.FanController.Init()
	c.FanController.SetFanPercent(c.DefaultFanSpeed)
	c.LEDs.Write([]byte{0, 0, 0, 0, 0, 0})

	group := errgroup.Group{}

	// LED Update events
	println("Starting LED update loop")
	group.Go(func() error {
		defer cancel()
		if err := c.updateLEDs(ctx); err != nil {
			return err
		}
		return nil
	})

	// Fan speed update events
	println("Starting fan update loop")
	group.Go(func() error {
		defer cancel()
		if err := c.updateFanSpeed(ctx); err != nil {
			return err
		}
		return nil
	})

	// Metric reporting events
	println("Starting metric reporting loop")
	group.Go(func() error {
		defer cancel()
		if err := c.metricReporter(ctx); err != nil {
			return err
		}
		return nil
	})

	// Left blade events
	println("Starting left blade event listener")
	group.Go(func() error {
		defer cancel()
		return c.listenEvents(ctx, c.LeftUART, leftBladeTopicIn)
	})
	println("Starting left blade event dispatcher")
	group.Go(func() error {
		defer cancel()
		return c.dispatchEvents(ctx, c.LeftUART, leftBladeTopicOut)
	})

	// right blade events
	println("Starting right blade event listener")
	group.Go(func() error {
		defer cancel()
		return c.listenEvents(ctx, c.RightUART, rightBladeTopicIn)
	})
	println("Starting right blade event dispatcher")
	group.Go(func() error {
		defer cancel()
		return c.dispatchEvents(ctx, c.RightUART, rightBladeTopicOut)
	})

	// Button Press events
	println("Starting button press event interrupt handler")
	c.ButtonPin.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		c.buttonPressed = true
	})

	group.Go(func() error {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
			}
			time.Sleep(10 * time.Millisecond)
			if c.buttonPressed {
				println("button pressed")
				btnPressed := smartfanunit.ButtonPressPacket{}
				c.eb.Publish(leftBladeTopicOut, btnPressed.Packet())
				c.eb.Publish(rightBladeTopicOut, btnPressed.Packet())
				c.buttonPressed = false
			}
		}
	})
	return group.Wait()
}

// listenEvents reads events from the UART interface and writes them to the eventbus
func (c *Controller) listenEvents(ctx context.Context, uart drivers.UART, targetTopic string) error {
	for {
		// Read packet from UART; blocks until packet is received
		pkt, err := proto.ReadPacket(ctx, uart)
		if err != nil {
			println(err.Error())
			continue
		}
		println("received packet from UART, publishing to topic", targetTopic)
		c.eb.Publish(targetTopic, pkt)
	}
}

// dispatchEvents reads events from the eventbus and writes them to the UART interface
func (c *Controller) dispatchEvents(ctx context.Context, uart drivers.UART, sourceTopic string) error {
	sub := c.eb.Subscribe(sourceTopic, 4, eventbus.MatchAll)
	defer sub.Unsubscribe()
	for {
		select {
		case msg := <-sub.C():
			println("dispatching event to UART from topic", sourceTopic)
			pkt := msg.(proto.Packet)
			err := proto.WritePacket(ctx, uart, pkt)
			if err != nil {
				println(err.Error())
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *Controller) metricReporter(ctx context.Context) error {
	var err error

	ticker := time.NewTicker(2 * time.Second)
	airFlowTempRight := smartfanunit.AirFlowTemperaturePacket{}
	airFlowTempLeft := smartfanunit.AirFlowTemperaturePacket{}
	fanRpm := smartfanunit.FanSpeedRPMPacket{}
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}

		airFlowTempLeft.Temperature, err = c.FanController.InternalTemperature()
		if err != nil {
			println(err.Error())
		}
		airFlowTempRight.Temperature, err = c.FanController.ExternalTemperature()
		if err != nil {
			println(err.Error())
		}
		fanRpm.RPM, err = c.FanController.FanRPM()
		if err != nil {
			println(err.Error())
		}

		// Publish metrics
		c.eb.Publish(leftBladeTopicOut, airFlowTempLeft.Packet())
		c.eb.Publish(rightBladeTopicOut, airFlowTempRight.Packet())
		c.eb.Publish(leftBladeTopicOut, fanRpm.Packet())
		c.eb.Publish(rightBladeTopicOut, fanRpm.Packet())
	}
}

func (c *Controller) updateFanSpeed(ctx context.Context) error {
	var pkt smartfanunit.SetFanSpeedPercentPacket

	subLeft := c.eb.Subscribe(leftBladeTopicIn, 1, eventbus.MatchAll)
	defer subLeft.Unsubscribe()
	subRight := c.eb.Subscribe(rightBladeTopicIn, 1, eventbus.MatchAll)
	defer subRight.Unsubscribe()

	for {
		// Update LED color depending on blade
		select {
		case msg := <-subLeft.C():
			pkt.FromPacket(msg.(proto.Packet))
			c.leftReqFanSpeed = pkt.Percent
		case msg := <-subRight.C():
			pkt.FromPacket(msg.(proto.Packet))
			c.rightReqFanSpeed = pkt.Percent
		case <-ctx.Done():
			return nil
		}

		// Update fan speed with the max requested speed
		if c.leftReqFanSpeed > c.rightReqFanSpeed {
			c.FanController.SetFanPercent(c.leftReqFanSpeed)
		} else {
			c.FanController.SetFanPercent(c.rightReqFanSpeed)
		}
	}
}

func (c *Controller) updateLEDs(ctx context.Context) error {
	subLeft := c.eb.Subscribe(leftBladeTopicIn, 1, smartfanunit.MatchCmd(smartfanunit.CmdSetLED))
	defer subLeft.Unsubscribe()
	subRight := c.eb.Subscribe(rightBladeTopicIn, 1, eventbus.MatchAll)
	defer subRight.Unsubscribe()

	var pkt smartfanunit.SetLEDPacket
	for {
		// Update LED color depending on blade
		select {
		case msg := <-subLeft.C():
			pkt.FromPacket(msg.(proto.Packet))
			c.leftLed = pkt.Color
		case msg := <-subRight.C():
			pkt.FromPacket(msg.(proto.Packet))
			c.rightLed = pkt.Color
		case <-ctx.Done():
			return nil
		}
		// Write to LEDs (they are in a chain -> we always have to update both)
		_, err := c.LEDs.Write([]byte{
			c.rightLed.Blue, c.rightLed.Green, c.rightLed.Red,
			c.leftLed.Blue, c.leftLed.Green, c.leftLed.Red,
		})
		if err != nil {
			println(err.Error())
			return err
		}
	}
}
