package rate_limiter

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

type BucketType int

const (
	SlidingWindowType BucketType = iota
	TokenBucketType
	LeakyBucketType
	FixedWindowType
)

func (r BucketType) String() string {
	return [...]string{"sliding_window", "token_bucket", "leaky_bucket", "fixed_window"}[r]
}

type KeyType int

const (
	IpAddress KeyType = iota
	ApiKey
	UserId
)

func (r KeyType) String() string {
	return [...]string{"ip_address", "api_key", "user_id"}[r]
}

type RateLimitConfig struct {
	Type   BucketType
	Key    KeyType
	Store  RateLimitStore
	Limit  int
	Window time.Duration
}

func NewRateLimiter(cfg RateLimitConfig) RateLimiter {
	switch cfg.Type {
	case TokenBucketType:
		return NewTokenBucketLimiter(cfg)
	case LeakyBucketType:
		return NewLeakyBucketLimiter(cfg)
	case FixedWindowType:
		return NewFixedWindowLimiter(cfg)
	case SlidingWindowType:
		return NewSlidingWindowLimiter(cfg)
	default:
		panic("invalid rate limiter type")
	}
}

func GetKey(c *fiber.Ctx, keyType KeyType) string {
	switch keyType {
	case IpAddress:
		return c.IP()
	case ApiKey:
		return c.Get("X-Api-Key")
	case UserId:
		return c.Locals("user_id").(string)
	default:
		panic("invalid key type")
	}
}
