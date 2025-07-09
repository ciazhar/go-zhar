package middleware

import (
	"github.com/ciazhar/go-start-small/pkg/rate_limiter"
	"github.com/gofiber/fiber/v2"
)

func RateLimitMiddleware(limiter rate_limiter.RateLimiter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := rate_limiter.GetKey(c, limiter.GetKeyType())

		if !limiter.Allow(key) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "rate limit exceeded",
			})
		}

		return c.Next()
	}
}
