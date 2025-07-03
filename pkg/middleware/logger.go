package middleware

import (
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID, _ := c.Locals("request_id").(string)

		logCtx := log.Logger.With().
			Str("request_id", reqID).
			Logger()

		ctx := logger.WithLogger(c.Context(), logCtx)
		c.SetUserContext(ctx)

		return c.Next()
	}
}
