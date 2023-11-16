package ratelimit

import (
	"sync"
	"time"
)

type BucketLimiter struct {
	mu       sync.Mutex
	capacity int
	rate     int
	current  int
	updated  time.Time
}

func NewBucketLimiter(capacity int, rate int) Limiter {
	return &BucketLimiter{
		capacity: capacity,
		rate:     rate,
		current:  capacity,
		updated:  time.Now(),
	}
}

func (l *BucketLimiter) Take() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refill()
	if l.current == 0 {
		return false
	}
	l.current--
	return true
}

func (l *BucketLimiter) refill() {
	now := time.Now()
	fill := int(now.Sub(l.updated).Seconds() * float64(l.rate))
	if fill > 0 {
		l.updated = now
		l.current += fill
		if l.current > l.capacity {
			l.current = l.capacity
		}
	}
}
