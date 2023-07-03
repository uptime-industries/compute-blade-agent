package main

import "github.com/xvzf/computeblade-agent/pkg/hal"

func main() {
	hal, err := hal.NewComputeBladeHAL(hal.ComputeBladeHalOpts{
		ComputeModuleType:         hal.COMPUTE_MODULE_TYPE_CM4,
		FanUnit:                   hal.FAN_UNIT_STANDARD,
		DefaultFanSpeed:           50,
		DefaultStealthModeEnabled: true,
	})
	if err != nil {
		panic(err)
	}
	defer hal.Close()
	hal.Init()
}
