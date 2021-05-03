package metrics

import "time"

type Timer struct {
	startTime time.Time
}

func NewTimer() *Timer {
	return &Timer{
		startTime: time.Now(),
	}
}

func (this *Timer) End() time.Duration {
	return time.Since(this.startTime)
}
