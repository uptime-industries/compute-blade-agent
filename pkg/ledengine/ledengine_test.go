package ledengine_test

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xvzf/computeblade-agent/pkg/hal"
	"github.com/xvzf/computeblade-agent/pkg/ledengine"
	"github.com/xvzf/computeblade-agent/pkg/util"
)

func TestNewStaticPattern(t *testing.T) {
	t.Parallel()

	type args struct {
		color hal.LedColor
	}
	tests := []struct {
		name string
		args args
		want ledengine.BlinkPattern
	}{
		{
			"Green",
			args{hal.LedColor{Green: 255}},
			ledengine.BlinkPattern{
				BaseColor:   hal.LedColor{Green: 255},
				ActiveColor: hal.LedColor{Green: 255},
				Delays:      []time.Duration{time.Hour},
			},
		},
		{
			"Red",
			args{hal.LedColor{Red: 255}},
			ledengine.BlinkPattern{
				BaseColor:   hal.LedColor{Red: 255},
				ActiveColor: hal.LedColor{Red: 255},
				Delays:      []time.Duration{time.Hour},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ledengine.NewStaticPattern(tt.args.color); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStaticPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBurstPattern(t *testing.T) {
	t.Parallel()
	type args struct {
		baseColor  hal.LedColor
		burstColor hal.LedColor
	}
	tests := []struct {
		name string
		args args
		want ledengine.BlinkPattern
	}{
		{
			"Green <-> Red",
			args{
				baseColor:  hal.LedColor{Green: 255},
				burstColor: hal.LedColor{Red: 255},
			},
			ledengine.BlinkPattern{
				BaseColor:   hal.LedColor{Green: 255},
				ActiveColor: hal.LedColor{Red: 255},
				Delays: []time.Duration{
					750 * time.Millisecond,
					50 * time.Millisecond,
					50 * time.Millisecond,
					50 * time.Millisecond,
					50 * time.Millisecond,
					50 * time.Millisecond,
				},
			},
		},
		{
			"Green <-> Green (valid, but no visual effect)",
			args{
				baseColor:  hal.LedColor{Green: 255},
				burstColor: hal.LedColor{Green: 255},
			},
			ledengine.BlinkPattern{
				BaseColor:   hal.LedColor{Green: 255},
				ActiveColor: hal.LedColor{Green: 255},
				Delays: []time.Duration{
					750 * time.Millisecond,
					50 * time.Millisecond,
					50 * time.Millisecond,
					50 * time.Millisecond,
					50 * time.Millisecond,
					50 * time.Millisecond,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ledengine.NewBurstPattern(tt.args.baseColor, tt.args.burstColor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBurstPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSlowBlinkPattern(t *testing.T) {
	type args struct {
		baseColor   hal.LedColor
		activeColor hal.LedColor
	}
	tests := []struct {
		name string
		args args
		want ledengine.BlinkPattern
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ledengine.NewSlowBlinkPattern(tt.args.baseColor, tt.args.activeColor); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewSlowledengine.BlinkPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLedEngine(t *testing.T) {
	t.Parallel()
	engine := ledengine.LedEngineOpts{
		Clock:  util.RealClock{},
		LedIdx: 0,
		Hal:    &hal.ComputeBladeHalMock{},
	}
	assert.NotNil(t, engine)
}

func Test_LedEngine_SetPattern_WhileRunning(t *testing.T) {
	t.Parallel()

	clk := util.MockClock{}
	clkAfterChan := make(chan time.Time)
	clk.On("After", time.Hour).Times(2).Return(clkAfterChan)

	cbMock := hal.ComputeBladeHalMock{}
	cbMock.On("SetLed", uint(0), hal.LedColor{Green: 0, Blue: 0, Red: 0}).Once().Return(nil)
	cbMock.On("SetLed", uint(0), hal.LedColor{Green: 0, Blue: 0, Red: 255}).Once().Return(nil)

	opts := ledengine.LedEngineOpts{
		Hal:    &cbMock,
		Clock:  &clk,
		LedIdx: 0,
	}

	engine := ledengine.NewLedEngine(opts)

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log("LedEngine.Run() started")
		engine.Run(ctx)
		t.Log("LedEngine.Run() exited")
	}()

	// We want to change the pattern while the engine is running
	time.Sleep(5 * time.Millisecond)

	// Set pattern
	t.Log("Setting pattern")
	engine.SetPattern(ledengine.NewStaticPattern(hal.LedColor{Red: 255}))

	t.Log("Canceling context")
	cancel()
	wg.Wait()

	clk.AssertExpectations(t)
	cbMock.AssertExpectations(t)
}

func Test_LedEngine_SetPattern_BeforeRun(t *testing.T) {
	t.Parallel()

	clk := util.MockClock{}
	clkAfterChan := make(chan time.Time)
	clk.On("After", time.Hour).Once().Return(clkAfterChan)

	cbMock := hal.ComputeBladeHalMock{}
	cbMock.On("SetLed", uint(0), hal.LedColor{Green: 0, Blue: 0, Red: 255}).Once().Return(nil)

	opts := ledengine.LedEngineOpts{
		Hal:    &cbMock,
		Clock:  &clk,
		LedIdx: 0,
	}

	engine := ledengine.NewLedEngine(opts)
	// We want to change the pattern BEFORE the engine is started
	t.Log("Setting pattern")
	engine.SetPattern(ledengine.NewStaticPattern(hal.LedColor{Red: 255}))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log("LedEngine.Run() started")
		engine.Run(ctx)
		t.Log("LedEngine.Run() exited")
	}()

	t.Log("Waiting for context to timeout")
	wg.Wait()

	clk.AssertExpectations(t)
	cbMock.AssertExpectations(t)
}

func Test_LedEngine_SetPattern_SetLedFailureInPattern(t *testing.T) {
	t.Parallel()

	clk := util.MockClock{}
	clkAfterChan := make(chan time.Time)
	clk.On("After", time.Hour).Once().Return(clkAfterChan)

	cbMock := hal.ComputeBladeHalMock{}
	call0 := cbMock.On("SetLed", uint(0), hal.LedColor{Green: 0, Blue: 0, Red: 0}).Once().Return(nil)
	cbMock.On("SetLed", uint(0), hal.LedColor{Green: 0, Blue: 0, Red: 0}).Once().Return(errors.New("failure")).NotBefore(call0)

	opts := ledengine.LedEngineOpts{
		Hal:    &cbMock,
		Clock:  &clk,
		LedIdx: 0,
	}

	engine := ledengine.NewLedEngine(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log("LedEngine.Run() started")
		err := engine.Run(ctx)
		assert.Error(t, err)
		t.Log("LedEngine.Run() exited")
	}()
	time.Sleep(5 * time.Millisecond)

	// Time tick -> SetLed() fails
	clkAfterChan <- time.Now()

	t.Log("Waiting for context to timeout")
	wg.Wait()

	clk.AssertExpectations(t)
	cbMock.AssertExpectations(t)
}
