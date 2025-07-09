package rate_limiter

import "time"

type leakyBucket struct {
	tokens   float64
	lastLeak time.Time
}

type leakyBucketLimiter struct {
	store    RateLimitStore
	capacity int
	leakRate float64 // tokens per second
	window   time.Duration
	keyType  string
}

func (l *leakyBucketLimiter) Allow(key string) bool {
	now := time.Now()

	val, ok := l.store.Get(key)
	var bucket *leakyBucket
	if ok {
		bucket = val.(*leakyBucket)
	} else {
		bucket = &leakyBucket{
			tokens:   0,
			lastLeak: now,
		}
	}

	// Hitung token yang bisa keluar (leak) sejak terakhir
	elapsed := now.Sub(bucket.lastLeak).Seconds()
	leaked := elapsed * l.leakRate
	bucket.tokens -= leaked
	if bucket.tokens < 0 {
		bucket.tokens = 0
	}
	bucket.lastLeak = now

	// Cek apakah ember masih muat
	if bucket.tokens < float64(l.capacity) {
		bucket.tokens += 1
		l.store.Set(key, bucket, l.window)
		return true
	}

	// Ember penuh
	l.store.Set(key, bucket, l.window)
	return false
}

func (l *leakyBucketLimiter) GetKeyType() string {
	return l.keyType
}

func NewLeakyBucketLimiter(cfg RateLimitConfig) RateLimiter {
	return &leakyBucketLimiter{
		store:    cfg.Store,
		capacity: cfg.Limit,
		leakRate: float64(cfg.Limit) / cfg.Window.Seconds(), // 1 request per X seconds
		window:   cfg.Window,
		keyType:  cfg.Key,
	}
}
