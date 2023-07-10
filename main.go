package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xvzf/computeblade-agent/pkg/hal"
	"github.com/xvzf/computeblade-agent/pkg/hal/bcm2711"
)

func main() {
	blade, err := bcm2711.New(hal.ComputeBladeHalOpts{
		FanUnit: hal.FanUnitStandard,
		Defaults: hal.ComputeBladeHalOptsDefault{
			StealthModeEnabled: false,
			FanSpeed:           20,
			TopLedColor:        hal.LedColor{},
			EdgeLedColor:       hal.LedColor{Red: 0, Green: 5, Blue: 5},
		},
	})
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	// setup signal handler
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		err := blade.Close()
		if err != nil {
			fmt.Printf("Error closing blade: %s\n", err)
		}
		fmt.Println("Exiting...")
		os.Exit(0)
	}()

	// just dummy print for now
	for {
		fmt.Printf("Power status: %s\n", blade.GetPowerStatus())
		fmt.Printf("Fan speed: %d\n", blade.GetFanSpeed())
		time.Sleep(1 * time.Second)
	}

}
