package rate_limiter

import (
	"time"
)

type slidingWindowLimiter struct {
	store   RateLimitStore
	limit   int
	window  time.Duration
	keyType KeyType
}

func (s *slidingWindowLimiter) Allow(key string) bool {
	now := time.Now()

	// Get or initialize the timestamps
	val, ok := s.store.Get(key)
	var timestamps []time.Time
	if ok {
		timestamps = val.([]time.Time)
	}

	// Filter out old timestamps
	filtered := make([]time.Time, 0, len(timestamps))
	for _, ts := range timestamps {
		if now.Sub(ts) < s.window {
			filtered = append(filtered, ts)
		}
	}

	if len(filtered) >= s.limit {
		return false
	}

	filtered = append(filtered, now)
	s.store.Set(key, filtered, s.window)
	return true
}

func (l *slidingWindowLimiter) GetKeyType() KeyType {
	return l.keyType
}

func NewSlidingWindowLimiter(cfg RateLimitConfig) RateLimiter {
	return &slidingWindowLimiter{
		store:   cfg.Store,
		limit:   cfg.Limit,
		window:  cfg.Window,
		keyType: cfg.Key,
	}
}
