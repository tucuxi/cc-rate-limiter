package ratelimit

import (
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	mu    sync.Mutex
	rate  Rate
	index int
	w     [2]int64
	n     [2]int
}

func NewSlidingWindowLimiter(rate Rate) Limiter {
	return &SlidingWindowLimiter{
		rate: rate,
	}
}

func (l *SlidingWindowLimiter) Take() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UnixMilli()
	d := l.rate.D.Milliseconds()
	w := now / d

	if w > l.w[l.index] {
		l.index = 1 - l.index
		l.w[l.index] = w
		l.n[l.index] = 0
	}

	var n int

	if l.w[1-l.index] == w-1 {
		p := float64(now-w*d) / float64(d)
		n = int((1-p)*float64(l.n[1-l.index]) + p*float64(l.n[l.index]))
	} else {
		n = l.n[l.index]
	}

	take := n < l.rate.N
	if take {
		l.n[l.index]++
	}
	return take
}
