package ratelimit

import (
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	mu      sync.Mutex
	rate    Rate
	current int
	window  int64
}

func NewFixedWindowLimiter(rate Rate) Limiter {
	return &FixedWindowLimiter{
		rate:   rate,
		window: time.Now().UnixMilli() / rate.D.Milliseconds(),
	}
}

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
