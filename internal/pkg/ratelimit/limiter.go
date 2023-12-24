package ratelimit

import "time"

// Rate represents a number of requests per time unit.
type Rate struct {
	N int
	D time.Duration
}

// PerSecond gives a rate of n requests per second.
func PerSecond(n int) Rate {
	return Rate{n, time.Second}
}

// PerMinute gives a rate of n requests per minute.
func PerMinute(n int) Rate {
	return Rate{n, time.Minute}
}

// Limiter is the interface that rate limiters must implement.
type Limiter interface {
	Take() bool
}
