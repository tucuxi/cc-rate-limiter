package ratelimit

import (
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	mu       sync.Mutex
	size     int
	capacity int
	current  int
	window   int64
}

func NewFixedWindowLimiter(size, capacity int) Limiter {
	return &FixedWindowLimiter{
		size:     size,
		capacity: capacity,
		current:  0,
		window:   time.Now().UnixMilli() / 1000 / int64(size),
	}
}

func (l *FixedWindowLimiter) Take() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	w := time.Now().UnixMilli() / 1000 / int64(l.size)
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
