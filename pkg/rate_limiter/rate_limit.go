package rate_limiter

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

const (

	// Bucket types
	TokenBucketType   = "token"
	LeakyBucketType   = "leaky"
	FixedWindowType   = "fixed"
	SlidingWindowType = "sliding"

	// Key types
	IpAddress = "ip_address"
	ApiKey    = "api_key"
	UserId    = "user_id"

	// Storage
	InMemory = "in_memory"
	Redis    = "redis"
)

type RateLimitConfig struct {
	Type   string
	Key    string
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

func GetKey(c *fiber.Ctx, keyType string) string {
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
