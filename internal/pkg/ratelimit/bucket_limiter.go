package ratelimit

import (
	"sync"
	"time"
)

// BucketLimiter is a rate limiter based on the token bucket algorithm.
type BucketLimiter struct {
	mu      sync.Mutex
	rate    Rate
	current int
	updated time.Time
}

// NewBucketLimiter returns a new BucketLimiter with the given rate.
func NewBucketLimiter(rate Rate) *BucketLimiter {
	return &BucketLimiter{
		rate:    rate,
		current: rate.N,
		updated: time.Now(),
	}
}

// Take returns true as long as the request rate remains within the allowed limit.
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
