//go:build !tinygo

package hal_test

import (
	"context"
	"log"

	"github.com/xvzf/computeblade-agent/pkg/hal"
	"github.com/xvzf/computeblade-agent/pkg/hal/led"
)

func ExampleNewSmartFanUnitClient() {
	ctx := context.Background()

	client, err := hal.NewSmartFanUnitClient("/dev/tty.usbmodem11102")
	if err != nil {
		panic(err)
	}
	go func() {
		err := client.Run(ctx)
		if err != nil {
			panic(err)
		}
	}()

	// Set LED color for the blade to red
	client.SetLed(ctx, led.Color{Red: 100, Green: 0, Blue: 0})

	// Set fanspeed to 20%
	client.SetFanSpeedPercent(ctx, 20)

	tmp, err := client.AirFlowTemperature(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("AirflowTemp", tmp)
	rpm, err := client.FanSpeedRPM(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("RPM", rpm)
}
