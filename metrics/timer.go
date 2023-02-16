package metrics

import "time"

// Helper struct to measure time spent.
type Timer interface {
	// Returns the duration since start.
	End() time.Duration
}

// Helper struct to measure time spent.
type timerImpl struct {
	startTime time.Time
}

// Creates and starts a new timer.
func NewTimer() Timer {
	return &timerImpl{
		startTime: time.Now(),
	}
}

// Returns the duration since start.
func (timer *timerImpl) End() time.Duration {
	return time.Since(timer.startTime)
}
