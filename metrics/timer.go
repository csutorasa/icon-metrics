package metrics

import "time"

type Timer interface {
	End() time.Duration
}

type timeMeter struct {
	startTime time.Time
}

func NewTimer() Timer {
	return &timeMeter{
		startTime: time.Now(),
	}
}

func (this *timeMeter) End() time.Duration {
	return time.Since(this.startTime)
}
