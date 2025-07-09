package rate_limiter

import (
	"fmt"
	"time"
)

type slidingWindowLimiter struct {
	store   RateLimitStore
	limit   int
	window  time.Duration
	keyType KeyType
}

type SlidingWindowData struct {
	Timestamps []time.Time `json:"timestamps"`
}

func (s *slidingWindowLimiter) Allow(key string) (bool, error) {
	now := time.Now()

	// Get or initialize the timestamps
	var data SlidingWindowData
	found, err := s.store.Get(key, &data)
	if err != nil {
		return false, fmt.Errorf("failed to get sliding window data: %v", err)
	}
	if !found {
		data = SlidingWindowData{}
	}

	// Filter out old timestamps
	filtered := make([]time.Time, 0, len(data.Timestamps))
	for _, ts := range data.Timestamps {
		if now.Sub(ts) < s.window {
			filtered = append(filtered, ts)
		}
	}

	if len(filtered) >= s.limit {
		return false, nil
	}

	filtered = append(filtered, now)
	data.Timestamps = filtered

	err = s.store.Set(key, data, s.window)
	if err != nil {
		return false, fmt.Errorf("failed to store sliding window data: %v", err)
	}
	return true, nil
}

func (s *slidingWindowLimiter) GetKeyType() KeyType {
	return s.keyType
}

func NewSlidingWindowLimiter(cfg RateLimitConfig) RateLimiter {
	return &slidingWindowLimiter{
		store:   cfg.Store,
		limit:   cfg.Limit,
		window:  cfg.Window,
		keyType: cfg.Key,
	}
}
