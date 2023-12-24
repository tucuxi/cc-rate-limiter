package ratelimit

import (
	"sync"
	"time"
)

// FixedWindowLimiter is a rate limiter based on the fixed window algorithm.
type FixedWindowLimiter struct {
	mu      sync.Mutex
	rate    Rate
	current int
	window  int64
}

// NewFixedWindowLimiter returns a new FixedWindowLimiter with the given rate.
func NewFixedWindowLimiter(rate Rate) Limiter {
	return &FixedWindowLimiter{
		rate:   rate,
		window: time.Now().UnixMilli() / rate.D.Milliseconds(),
	}
}

// Take returns true as long as the request rate remains within the allowed limit.
func (l *FixedWindowLimiter) Take() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if w := time.Now().UnixMilli() / l.rate.D.Milliseconds(); w > l.window {
		l.window = w
		l.current = 0
	}

	take := l.current < l.rate.N
	if take {
		l.current++
	}
	return take
}
