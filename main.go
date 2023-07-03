package main

import (
	"github.com/xvzf/computeblade-agent/pkg/hal"
	"github.com/xvzf/computeblade-agent/pkg/hal/bcm2711"
)

func main() {
	blade, err := bcm2711.New(hal.ComputeBladeHalOpts{
		FanUnit:                   hal.FAN_UNIT_STANDARD,
		DefaultFanSpeed:           40,
		DefaultStealthModeEnabled: false,
	})
	if err != nil {
		panic(err)
	}
	defer blade.Close()
	blade.Init()
}
