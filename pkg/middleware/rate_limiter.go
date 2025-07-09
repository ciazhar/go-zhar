package middleware

import (
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/rate_limiter"
	"github.com/ciazhar/go-start-small/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func RateLimitMiddleware(limiter rate_limiter.RateLimiter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := rate_limiter.GetKey(c, limiter.GetKeyType())

		allow, err := limiter.Allow(key)
		if err != nil {
			return response.HandleWarnings(c, logger.FromContext(c.UserContext()), fiber.StatusInternalServerError, "failed to check rate limit", []response.ValidationError{{Message: err.Error()}})
		}

		if !allow {
			return response.HandleWarnings(c, logger.FromContext(c.UserContext()), fiber.StatusTooManyRequests, "rate limit exceeded", []response.ValidationError{{Message: "rate limit exceeded"}})
		}

		return c.Next()
	}
}
