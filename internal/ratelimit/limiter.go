package ratelimit

type Limiter interface {
	Take() bool
}
