package metrics

import "time"

// Helper struct to measure time spent.
type Timer interface {
	// Returns the duration since start.
	End() time.Duration
	Reset()
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

func (timer *timerImpl) Reset() {
	timer.startTime = time.Now()
}

func (timer *timerImpl) End() time.Duration {
	return time.Since(timer.startTime)
}

func Timed(f func()) time.Duration {
	t := NewTimer()
	f()
	return t.End()
}
