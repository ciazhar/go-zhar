package rate_limiter

import (
	"fmt"
	"time"
)

type tokenBucket struct {
	Tokens     float64   `json:"tokens"`
	LastRefill time.Time `json:"last_refill"`
}

type tokenBucketLimiter struct {
	store      RateLimitStore
	capacity   int
	refillRate float64 // tokens per second
	window     time.Duration
	keyType    KeyType
}

func (l *tokenBucketLimiter) Allow(key string) (bool, error) {
	now := time.Now()

	var bucket tokenBucket
	found, err := l.store.Get(key, &bucket)
	if err != nil {
		return false, fmt.Errorf("failed to get token bucket: %v", err)
	}
	if !found {
		bucket = tokenBucket{
			Tokens:     float64(l.capacity),
			LastRefill: now,
		}
	}

	// Refill token
	elapsed := now.Sub(bucket.LastRefill).Seconds()
	bucket.Tokens += elapsed * l.refillRate
	if bucket.Tokens > float64(l.capacity) {
		bucket.Tokens = float64(l.capacity)
	}
	bucket.LastRefill = now

	if bucket.Tokens >= 1 {
		bucket.Tokens -= 1
		if err := l.store.Set(key, bucket, l.window); err != nil {
			return false, fmt.Errorf("failed to set token bucket: %v", err)
		}
		return true, nil
	}

	if err := l.store.Set(key, bucket, l.window); err != nil {
		return false, fmt.Errorf("failed to set token bucket: %v", err)
	}
	return false, nil
}

func (l *tokenBucketLimiter) GetKeyType() KeyType {
	return l.keyType
}

func NewTokenBucketLimiter(cfg RateLimitConfig) RateLimiter {
	return &tokenBucketLimiter{
		store:      cfg.Store,
		capacity:   cfg.Limit,
		refillRate: float64(cfg.Limit) / cfg.Window.Seconds(), // tokens per second
		window:     cfg.Window,
		keyType:    cfg.Key,
	}
}
