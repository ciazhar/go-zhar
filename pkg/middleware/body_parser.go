package middleware

import (
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/ciazhar/go-zhar/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

func BodyParserMiddleware[T any](v validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			body T
			log  = logger.FromContext(c.UserContext())
		)

		if err := c.BodyParser(&body); err != nil {
			return response.HandleWarning(c, log, fiber.StatusBadRequest, "invalid request body", err)
		}

		errs, err := v.ValidateStruct(body)
		if err != nil {
			if err.Error() != "validation failed" {
				return response.HandleWarning(c, log, fiber.StatusBadRequest, "failed to validate request body", err)
			}
		}
		if len(errs) > 0 {
			return response.HandleWarnings(c, log, fiber.StatusBadRequest, "failed to validate request body", errs)
		}

		c.Locals("body", body)
		return c.Next()
	}
}
