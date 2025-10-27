package metrics

import "time"

var startTime = time.Now()

// Gets the uptime duration.
func Uptime() time.Duration {
	return time.Since(startTime)
}
