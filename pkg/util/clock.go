package util

import (
	"time"
)

// Clock is an interface for the time package
type Clock interface {
	// Now returns the current time.
	Now() time.Time
	// After waits for the duration to elapse and then sends the current time
	After(d time.Duration) <-chan time.Time
}

// RealClock implements the Clock interface using the real time package.
type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

func (RealClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}
