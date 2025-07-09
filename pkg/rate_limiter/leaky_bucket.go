package rate_limiter

import (
	"fmt"
	"time"
)

type leakyBucket struct {
	Tokens   float64   `json:"tokens"`
	LastLeak time.Time `json:"last_leak"`
}

type leakyBucketLimiter struct {
	store    RateLimitStore
	capacity int
	leakRate float64 // tokens per second
	window   time.Duration
	keyType  KeyType
}

func (l *leakyBucketLimiter) Allow(key string) (bool, error) {
	now := time.Now()

	var bucket leakyBucket
	found, err := l.store.Get(key, &bucket)
	if err != nil {
		return false, fmt.Errorf("failed to get leaky bucket data: %v", err)
	}
	if !found {
		bucket = leakyBucket{
			Tokens:   0,
			LastLeak: now,
		}
	}

	// Hitung token yang bocor sejak terakhir
	elapsed := now.Sub(bucket.LastLeak).Seconds()
	leaked := elapsed * l.leakRate
	bucket.Tokens -= leaked
	if bucket.Tokens < 0 {
		bucket.Tokens = 0
	}
	bucket.LastLeak = now

	if bucket.Tokens < float64(l.capacity) {
		bucket.Tokens += 1
		if err := l.store.Set(key, bucket, l.window); err != nil {
			return false, fmt.Errorf("failed to store leaky bucket: %v", err)
		}
		return true, nil
	}

	if err := l.store.Set(key, bucket, l.window); err != nil {
		return false, fmt.Errorf("failed to store leaky bucket: %v", err)
	}
	return false, nil
}

func (l *leakyBucketLimiter) GetKeyType() KeyType {
	return l.keyType
}

func NewLeakyBucketLimiter(cfg RateLimitConfig) RateLimiter {
	return &leakyBucketLimiter{
		store:    cfg.Store,
		capacity: cfg.Limit,
		leakRate: float64(cfg.Limit) / cfg.Window.Seconds(), // tokens per second
		window:   cfg.Window,
		keyType:  cfg.Key,
	}
}
