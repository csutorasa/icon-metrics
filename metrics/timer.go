package metrics

import "time"

// Helper struct to measure time spent.
type Timer struct {
	startTime time.Time
}

// Creates and starts a new timer.
func NewTimer() *Timer {
	return &Timer{
		startTime: time.Now(),
	}
}

// Returns the duration since start.
func (timer *Timer) End() time.Duration {
	return time.Since(timer.startTime)
}
