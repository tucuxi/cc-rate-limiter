package ratelimit

import (
	"sync"
	"time"
)

type SlidingWindowLogLimiter struct {
	mu   sync.Mutex
	rate Rate
	log  []int64
}

func NewSlidingWindowLogLimiter(rate Rate) Limiter {
	return &SlidingWindowLogLimiter{
		rate: rate,
	}
}

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
