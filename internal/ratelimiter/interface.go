package ratelimiter

import (
	"sync"
	"time"
)

type Bucket struct {
	tokens     float64
	lastRefill time.Time
}

type RateLimiter struct {
	mu     sync.Mutex
	bucket *Bucket
	rate   float64 // limit/duration
	cap    float64 // bucket capacity
}

func NewRateLimiter(limit int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		bucket: &Bucket{},
		rate:   float64(limit) / duration.Seconds(),
		cap:    float64(limit),
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if rl.bucket == nil {
		bucket := &Bucket{tokens: rl.cap, lastRefill: now}
		rl.bucket = bucket
	}

	// refill tokens
	elapsed := now.Sub(rl.bucket.lastRefill).Seconds()
	rl.bucket.tokens += elapsed * rl.rate
	if rl.bucket.tokens > rl.cap {
		rl.bucket.tokens = rl.cap
	}
	rl.bucket.lastRefill = now

	// check availability
	if rl.bucket.tokens >= 1 {
		rl.bucket.tokens--
		return true
	}
	return false
}
