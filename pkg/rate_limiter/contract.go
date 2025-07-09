package rate_limiter

import "time"

type RateLimiter interface {
	Allow(key string) bool
	GetKeyType() string
}

type RateLimitStore interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
}
