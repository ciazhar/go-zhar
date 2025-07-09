package rate_limiter

import (
	"fmt"
	"time"
)

type fixedWindowData struct {
	Count   int       `json:"count"`
	ResetAt time.Time `json:"reset_at"`
}

type fixedWindowLimiter struct {
	store   RateLimitStore
	limit   int
	window  time.Duration
	keyType KeyType
}

func (f *fixedWindowLimiter) Allow(key string) (bool, error) {
	now := time.Now()

	var data fixedWindowData
	found, err := f.store.Get(key, &data)
	if err != nil {
		return false, fmt.Errorf("failed to get fixed window data: %v", err)
	}
	if !found {
		data = fixedWindowData{
			Count:   0,
			ResetAt: now.Add(f.window),
		}
	}

	if now.After(data.ResetAt) {
		// Reset window
		data.Count = 0
		data.ResetAt = now.Add(f.window)
	}

	if data.Count < f.limit {
		data.Count++
		if err := f.store.Set(key, data, f.window); err != nil {
			return false, fmt.Errorf("failed to store fixed window data: %v", err)
		}
		return true, nil
	}

	// still store to maintain reset time
	if err := f.store.Set(key, data, f.window); err != nil {
		return false, fmt.Errorf("failed to store fixed window data: %v", err)
	}
	return false, nil
}

func (f *fixedWindowLimiter) GetKeyType() KeyType {
	return f.keyType
}

func NewFixedWindowLimiter(cfg RateLimitConfig) RateLimiter {
	return &fixedWindowLimiter{
		store:   cfg.Store,
		limit:   cfg.Limit,
		window:  cfg.Window,
		keyType: cfg.Key,
	}
}
