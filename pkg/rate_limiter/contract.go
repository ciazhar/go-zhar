package rate_limiter

import "time"

type RateLimiter interface {
	Allow(key string) (bool, error)
	GetKeyType() KeyType
}

type RateLimitStore interface {
	Get(key string, out interface{}) (found bool, err error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string)
	Type() StorageType
}
