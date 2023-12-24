package ratelimit

import (
	"math"
	"sync"
	"time"
)

// SlidingWindowLimiter is a rate limiter based on the sliding window algorithm.
type SlidingWindowLimiter struct {
	mu   sync.Mutex
	rate Rate
	wp   int64
	np   int
	w    int64
	n    int
}

// NewSlidingWindowLimiter returns a new SlidingWindowLimiter with the given rate.
func NewSlidingWindowLimiter(rate Rate) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		rate: rate,
	}
}

// Take returns true as long as the request rate remains within the allowed limit.
func (l *SlidingWindowLimiter) Take() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UnixMilli()
	d := l.rate.D.Milliseconds()
	w := now / d

	if w > l.w {
		if l.wp+1 == w {
			l.np = l.n
		} else {
			l.np = 0
		}
		l.wp = w - 1
		l.w = w
		l.n = 0
	}

	l.n++
	p := 1 - float64(now-w*d)/float64(d)
	n := l.n + int(math.Round(p*float64(l.np)))

	return n <= l.rate.N
}
