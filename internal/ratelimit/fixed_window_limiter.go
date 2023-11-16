package ratelimit

import (
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	mu       sync.Mutex
	capacity int
	current  int
	window   int64
}

const WINDOW_SIZE = 60000 /* ms */

func NewFixedWindowLimiter(rate int) Limiter {
	return &FixedWindowLimiter{
		capacity: rate,
		current:  0,
		window:   time.Now().UnixMilli() / WINDOW_SIZE,
	}
}

func (l *FixedWindowLimiter) Take() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	w := time.Now().UnixMilli() / WINDOW_SIZE
	if w > l.window {
		l.window = w
		l.current = 0
	}
	if l.current == l.capacity {
		return false
	}
	l.current++
	return true
}
