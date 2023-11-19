package ratelimit

import "time"

type Rate struct {
	N int
	D time.Duration
}

func PerSecond(n int) Rate {
	return Rate{n, time.Second}
}

func PerMinute(n int) Rate {
	return Rate{n, time.Minute}
}

type Limiter interface {
	Take() bool
}
