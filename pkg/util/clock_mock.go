package util

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// MockClock implements the Clock interface using the testify mock package.
type MockClock struct {
	mock.Mock
}

// Now returns the current time.
func (mc *MockClock) Now() time.Time {
	args := mc.Called()
	return args.Get(0).(time.Time)
}

// After waits for the duration to elapse and then sends the current time
func (mc *MockClock) After(d time.Duration) <-chan time.Time {
	args := mc.Called(d)
	return args.Get(0).(chan time.Time)
}
