package middleware

import (
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/response"
	"github.com/ciazhar/go-start-small/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

func QueryParamParserMiddleware[T any](v validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			queryParam T
			log        = logger.FromContext(c.UserContext())
		)
		if err := c.QueryParser(&queryParam); err != nil {
			return response.HandleWarning(c, log, fiber.StatusBadRequest, "invalid request query param", err)
		}

		errs, err := v.ValidateStruct(queryParam)
		if err != nil {
			if err.Error() != "validation failed" {
				return response.HandleWarning(c, log, fiber.StatusBadRequest, "failed to validate request query param", err)
			}
		}
		if len(errs) > 0 {
			return response.HandleWarnings(c, log, fiber.StatusBadRequest, "failed to validate request query param", errs)
		}

		c.Locals("query_param", queryParam)
		return c.Next()
	}
}
