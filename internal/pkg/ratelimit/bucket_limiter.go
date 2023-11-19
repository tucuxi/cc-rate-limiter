package ratelimit

import (
	"sync"
	"time"
)

type BucketLimiter struct {
	mu      sync.Mutex
	rate    Rate
	current int
	updated time.Time
}

func NewBucketLimiter(rate Rate) Limiter {
	return &BucketLimiter{
		rate:    rate,
		current: rate.N,
		updated: time.Now(),
	}
}

func (l *BucketLimiter) Take() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	fill := int(int64(l.rate.N) * now.Sub(l.updated).Milliseconds() / l.rate.D.Milliseconds())
	if fill > 0 {
		l.updated = now
		l.current = min(l.current+fill, l.rate.N)
	}

	take := l.current > 0
	if take {
		l.current--
	}
	return take
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
