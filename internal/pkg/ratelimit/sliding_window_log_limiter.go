package ratelimit

import (
	"sync"
	"time"
)

// SlidingWindowLogLimiter is a rate limiter based on the sliding window log algorithm.
type SlidingWindowLogLimiter struct {
	mu   sync.Mutex
	rate Rate
	log  []int64
}

// NewSlidingWindowLogLimiter returns a new SlidingWindowLogLimiter with the given rate.
func NewSlidingWindowLogLimiter(rate Rate) *SlidingWindowLogLimiter {
	return &SlidingWindowLogLimiter{
		rate: rate,
	}
}

// Take returns true as long as the request rate remains within the allowed limit.
func (l *SlidingWindowLogLimiter) Take() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UnixMilli()
	threshold := now - l.rate.D.Milliseconds()

	for len(l.log) > 0 && l.log[0] <= threshold {
		l.log = l.log[1:]
	}

	take := len(l.log) < l.rate.N
	if take {
		l.log = append(l.log, now)
	}
	return take
}
