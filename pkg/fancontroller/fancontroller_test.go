// fancontroller_test.go
package fancontroller_test

import (
	"testing"

	"github.com/uptime-induestries/compute-blade-agent/pkg/fancontroller"
)

func TestFanControllerLinear_GetFanSpeed(t *testing.T) {
	t.Parallel()

	config := fancontroller.FanControllerConfig{
		Steps: []fancontroller.FanControllerStep{
			{Temperature: 20, Percent: 30},
			{Temperature: 30, Percent: 60},
		},
	}

	controller, err := fancontroller.NewLinearFanController(config)
	if err != nil {
		t.Fatalf("Failed to create fan controller: %v", err)
	}

	testCases := []struct {
		temperature float64
		expected    uint8
	}{
		{15, 30}, // Should use the minimum speed
		{25, 45}, // Should calculate speed based on linear function
		{35, 60}, // Should use the maximum speed
	}

	for _, tc := range testCases {
		expected := tc.expected
		temperature := tc.temperature
		t.Run("", func(t *testing.T) {
			t.Parallel()
			speed := controller.GetFanSpeed(temperature)
			if speed != expected {
				t.Errorf("For temperature %.2f, expected speed %d but got %d", temperature, expected, speed)
			}
		})
	}
}

func TestFanControllerLinear_GetFanSpeedWithOverride(t *testing.T) {
	t.Parallel()

	config := fancontroller.FanControllerConfig{
		Steps: []fancontroller.FanControllerStep{
			{Temperature: 20, Percent: 30},
			{Temperature: 30, Percent: 60},
		},
	}

	controller, err := fancontroller.NewLinearFanController(config)
	if err != nil {
		t.Fatalf("Failed to create fan controller: %v", err)
	}
	controller.Override(&fancontroller.FanOverrideOpts{
		Percent: 99,
	})

	testCases := []struct {
		temperature float64
		expected    uint8
	}{
		{15, 99},
		{25, 99},
		{35, 99},
	}

	for _, tc := range testCases {
		expected := tc.expected
		temperature := tc.temperature
		t.Run("", func(t *testing.T) {
			t.Parallel()
			speed := controller.GetFanSpeed(temperature)
			if speed != expected {
				t.Errorf("For temperature %.2f, expected speed %d but got %d", temperature, expected, speed)
			}
		})
	}
}

func TestFanControllerLinear_ConstructionErrors(t *testing.T) {
	testCases := []struct {
		name   string
		config fancontroller.FanControllerConfig
		errMsg string
	}{
		{
			name: "InvalidStepCount",
			config: fancontroller.FanControllerConfig{
				Steps: []fancontroller.FanControllerStep{
					{Temperature: 20, Percent: 30},
				},
			},
			errMsg: "exactly two steps must be defined",
		},
		{
			name: "InvalidStepTemperatures",
			config: fancontroller.FanControllerConfig{
				Steps: []fancontroller.FanControllerStep{
					{Temperature: 30, Percent: 60},
					{Temperature: 20, Percent: 30},
				},
			},
			errMsg: "step 1 temperature must be lower than step 2 temperature",
		},
		{
			name: "InvalidStepSpeeds",
			config: fancontroller.FanControllerConfig{
				Steps: []fancontroller.FanControllerStep{
					{Temperature: 20, Percent: 60},
					{Temperature: 30, Percent: 30},
				},
			},
			errMsg: "step 1 speed must be lower than step 2 speed",
		},
		{
			name: "InvalidSpeedRange",
			config: fancontroller.FanControllerConfig{
				Steps: []fancontroller.FanControllerStep{
					{Temperature: 20, Percent: 10},
					{Temperature: 30, Percent: 200},
				},
			},
			errMsg: "speed must be between 0 and 100",
		},
	}

	for _, tc := range testCases {
		config := tc.config
		expectedErrMsg := tc.errMsg
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := fancontroller.NewLinearFanController(config)
			if err == nil {
				t.Errorf("Expected error with message '%s', but got no error", expectedErrMsg)
			} else if err.Error() != expectedErrMsg {
				t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
			}
		})
	}
}
