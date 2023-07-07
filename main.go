package main

import (
	"time"

	"github.com/xvzf/computeblade-agent/pkg/hal"
	"github.com/xvzf/computeblade-agent/pkg/hal/bcm2711"
)

func main() {
	blade, err := bcm2711.New(hal.ComputeBladeHalOpts{
		FanUnit:                   hal.FAN_UNIT_STANDARD,
		DefaultFanSpeed:           10,
		DefaultStealthModeEnabled: false,
	})
	if err != nil {
		panic(err)
	}
	defer blade.Close()
	blade.Init()

	ledColorToggle := []bcm2711.LedColor{
		{
			Green: 50,
			Blue:  0,
			Red:   0,
		},
		{
			Green: 0,
			Blue:  50,
			Red:   0,
		},
		{
			Green: 0,
			Blue:  0,
			Red:   50,
		},
	}

	// just randomly go through colors!
	idxTop := 1
	idxEdge := 0
	for {
		<-time.After(time.Second * 1)
		idxTop++
		if idxTop > 2 {
			idxTop = 0
		}
		idxEdge++
		if idxEdge > 2 {
			idxEdge = 0
		}
		blade.SetLEDs(ledColorToggle[idxTop], ledColorToggle[idxEdge])
		blade.SetFanSpeed(30 * uint8(idxTop))
	}

}
