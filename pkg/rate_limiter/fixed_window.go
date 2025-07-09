package rate_limiter

import (
	"time"
)

type fixedWindowData struct {
	Count   int
	ResetAt time.Time
}

type fixedWindowLimiter struct {
	store   RateLimitStore
	limit   int
	window  time.Duration
	keyType KeyType
}

func (f *fixedWindowLimiter) Allow(key string) bool {
	now := time.Now()

	val, ok := f.store.Get(key)
	var data *fixedWindowData
	if ok {
		data = val.(*fixedWindowData)
	} else {
		data = &fixedWindowData{
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
		f.store.Set(key, data, f.window)
		return true
	}

	f.store.Set(key, data, f.window)
	return false
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
