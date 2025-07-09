package rate_limiter

import (
	"time"
)

type tokenBucket struct {
	tokens     float64
	lastRefill time.Time
}

type tokenBucketLimiter struct {
	store      RateLimitStore
	capacity   int
	refillRate float64 // tokens per second
	window     time.Duration
	keyType    KeyType
}

func (l *tokenBucketLimiter) Allow(key string) bool {
	now := time.Now()

	val, ok := l.store.Get(key)
	var bucket *tokenBucket
	if ok {
		bucket = val.(*tokenBucket)
	} else {
		bucket = &tokenBucket{
			tokens:     float64(l.capacity),
			lastRefill: now,
		}
	}

	// Refill token berdasarkan waktu yang telah lewat
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * l.refillRate
	if bucket.tokens > float64(l.capacity) {
		bucket.tokens = float64(l.capacity)
	}
	bucket.lastRefill = now

	// Cek apakah bisa ambil 1 token
	if bucket.tokens >= 1 {
		bucket.tokens -= 1
		l.store.Set(key, bucket, l.window)
		return true
	}

	l.store.Set(key, bucket, l.window)
	return false
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
