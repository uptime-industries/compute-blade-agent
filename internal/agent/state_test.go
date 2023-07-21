package agent_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xvzf/computeblade-agent/internal/agent"
)

func TestNewComputeBladeState(t *testing.T) {
	t.Parallel()

	state := agent.NewComputeBladeState()
	assert.NotNil(t, state)
}

func TestComputeBladeState_RegisterEventIdentify(t *testing.T) {
	t.Parallel()

	state := agent.NewComputeBladeState()

	// Identify event
	state.RegisterEvent(agent.IdentifyEvent)
	assert.True(t, state.IdentifyActive())
	state.RegisterEvent(agent.IdentifyConfirmEvent)
	assert.False(t, state.IdentifyActive())
}

func TestComputeBladeState_RegisterEventCritical(t *testing.T) {
	t.Parallel()

	state := agent.NewComputeBladeState()

	// critical event
	state.RegisterEvent(agent.CriticalEvent)
	assert.True(t, state.CriticalActive())
	state.RegisterEvent(agent.CriticalResetEvent)
	assert.False(t, state.CriticalActive())
}

func TestComputeBladeState_RegisterEventMixed(t *testing.T) {
	t.Parallel()

	state := agent.NewComputeBladeState()

	// Send a bunch of events
	state.RegisterEvent(agent.CriticalEvent)
	state.RegisterEvent(agent.CriticalResetEvent)
	state.RegisterEvent(agent.NoopEvent)
	state.RegisterEvent(agent.CriticalEvent)
	state.RegisterEvent(agent.NoopEvent)
	state.RegisterEvent(agent.IdentifyEvent)
	state.RegisterEvent(agent.IdentifyEvent)
	state.RegisterEvent(agent.CriticalResetEvent)
	state.RegisterEvent(agent.IdentifyEvent)

	assert.False(t, state.CriticalActive())
	assert.True(t, state.IdentifyActive())
}

func TestComputeBladeState_WaitForIdentifyConfirm_NoTimeout(t *testing.T) {
	t.Parallel()

	state := agent.NewComputeBladeState()

	// send identify event
	t.Log("Setting identify event")
	state.RegisterEvent(agent.IdentifyEvent)
	assert.True(t, state.IdentifyActive())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := context.Background()

		// Block until identify status is cleared
		t.Log("Waiting for identify confirm")
		err := state.WaitForIdentifyConfirm(ctx)
		assert.NoError(t, err)
	}()

	// Give goroutine time to start
	time.Sleep(50 * time.Millisecond)

	// confirm event
	state.RegisterEvent(agent.IdentifyConfirmEvent)
	t.Log("Identify event confirmed")

	wg.Wait()
}

func TestComputeBladeState_WaitForIdentifyConfirm_Timeout(t *testing.T) {
	t.Parallel()

	state := agent.NewComputeBladeState()

	// send identify event
	t.Log("Setting identify event")
	state.RegisterEvent(agent.IdentifyEvent)
	assert.True(t, state.IdentifyActive())

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Block until identify status is cleared
		t.Log("Waiting for identify confirm")
		err := state.WaitForIdentifyConfirm(ctx)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	}()

	// Give goroutine time to start.
	time.Sleep(50 * time.Millisecond)

	// confirm event
	state.RegisterEvent(agent.IdentifyConfirmEvent)
	t.Log("Identify event confirmed")

	wg.Wait()
}

func TestComputeBladeState_WaitForCriticalClear_NoTimeout(t *testing.T) {
	t.Parallel()

	state := agent.NewComputeBladeState()

	// send critical event
	t.Log("Setting critical event")
	state.RegisterEvent(agent.CriticalEvent)
	assert.True(t, state.CriticalActive())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := context.Background()

		// Block until critical status is cleared
		t.Log("Waiting for critical confirm")
		err := state.WaitForCriticalClear(ctx)
		assert.NoError(t, err)
	}()

	// Give goroutine time to start
	time.Sleep(50 * time.Millisecond)

	// confirm event
	state.RegisterEvent(agent.CriticalResetEvent)
	t.Log("critical event confirmed")

	wg.Wait()
}

func TestComputeBladeState_WaitForCriticalClear_Timeout(t *testing.T) {
	t.Parallel()

	state := agent.NewComputeBladeState()

	// send critical event
	t.Log("Setting critical event")
	state.RegisterEvent(agent.CriticalEvent)
	assert.True(t, state.CriticalActive())

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Block until critical status is cleared
		t.Log("Waiting for critical confirm")
		err := state.WaitForCriticalClear(ctx)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	}()

	// Give goroutine time to start.
	time.Sleep(50 * time.Millisecond)

	// confirm event
	state.RegisterEvent(agent.CriticalResetEvent)
	t.Log("critical event confirmed")

	wg.Wait()
}
