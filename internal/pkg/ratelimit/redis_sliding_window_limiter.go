package ratelimit

import (
	"context"
	"log"
	"log/slog"
	"math"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisSlidingWindowLimiter is a rate limiter based on the sliding window
// algorithm that is backed by Redis.
type RedisSlidingWindowLimiter struct {
	rc   *redis.Client
	rate Rate
}

// NewRedisSlidingWindowLimiter returns a new RedisSlidingWindowLimiter with
// the given rate.
func NewRedisSlidingWindowLimiter(rc *redis.Client, rate Rate) *RedisSlidingWindowLimiter {
	return &RedisSlidingWindowLimiter{
		rc:   rc,
		rate: rate,
	}
}

// Take returns true as long as the request rate remains within the allowed limit.
// If Redis is unavailable, Take defaults to turning off rate limiting.
func (l *RedisSlidingWindowLimiter) Take() bool {
	ctx := context.TODO()

	now := time.Now().UnixMilli()
	d := l.rate.D.Milliseconds()
	w := now / d // Requires that clocks are synced.

	slog.Info("Take", "now", now, "d", d, "w", w)
	var np int

	nps, err := l.rc.Get(ctx, strconv.FormatInt(w-1, 16)).Result()
	if err != redis.Nil {
		if err != nil {
			log.Printf("Could not get previous window count from Redis: %v\n", err)
			return true // Turn off rate limiting when Redis is failing.
		}
		np, err = strconv.Atoi(nps)
		if err != nil {
			log.Printf("Got invalid previous window count %s: %v\n", nps, err)
			return false // This means that something messed with our Redis instance.
		}
		if np > l.rate.N {
			np = l.rate.N
		}
	}

	slog.Info("previous window", "nps", nps, "np", np)

	nc, err := l.rc.IncrBy(ctx, strconv.FormatInt(w, 16), 1).Result()
	if err != nil {
		log.Printf("Could not increment current window counter: %v\n", err)
		return true
	}

	slog.Info("current window", "nc", nc)

	p := 1 - float64(now-w*d)/float64(d)
	n := int(nc) + int(math.Round(p*float64(np)))

	slog.Info("result", "p", p, "n", n)
	return n <= l.rate.N
}
