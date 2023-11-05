//go:build tinygo

package main

import (
	"context"
	"machine"
	"time"

	"github.com/xvzf/computeblade-agent/pkg/smartfanunit"
	"github.com/xvzf/computeblade-agent/pkg/smartfanunit/emc2101"
	"tinygo.org/x/drivers/ws2812"
)

func main() {
	// Configure status LED
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Configure UARTs
	machine.UART0.Configure(machine.UARTConfig{TX: machine.UART0_TX_PIN, RX: machine.UART0_RX_PIN})
	machine.UART0.SetBaudRate(smartfanunit.Baudrate)
	machine.UART1.Configure(machine.UARTConfig{TX: machine.UART1_TX_PIN, RX: machine.UART1_RX_PIN})
	machine.UART1.SetBaudRate(smartfanunit.Baudrate)

	// Enables fan, DO NOT CHANGE
	machine.GP16.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.GP16.Set(true)

	// WS2812 LEDs
	machine.GP15.Configure(machine.PinConfig{Mode: machine.PinOutput})
	brgLeds := ws2812.New(machine.GP15)

	// Configure button
	machine.GP12.Configure(machine.PinConfig{Mode: machine.PinInput})

	// Setup emc2101
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 100 * machine.KHz,
		SDA:       machine.I2C0_SDA_PIN,
		SCL:       machine.I2C0_SCL_PIN,
	})
	emc := emc2101.New(machine.I2C0)
	err := emc.Init()
	if err != nil {
		println("Failed to initialize emc2101")
	}

	println("IO initialized, starting controller...")

	// Run controller
	controller := &Controller{
		DefaultFanSpeed: 50,
		LEDs:            brgLeds,
		FanController:   emc,
		ButtonPin:       machine.GP12,
		LeftUART:        machine.UART0,
		RightUART:       machine.UART1,
	}

	// Solid green indicates the fan unit is running
	machine.LED.Set(true)
	err = controller.Run(context.Background())

	// Blinking -> something went wrong
	ledState := false
	for {
		ledState = !ledState
		machine.LED.Set(ledState)
		println("Controller exited with error:", err)
		time.Sleep(500 * time.Millisecond)
	}
}
