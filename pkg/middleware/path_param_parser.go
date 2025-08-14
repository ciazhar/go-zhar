package middleware

import (
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/ciazhar/go-zhar/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

func PathParamParserMiddleware[T any](v validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			pathParam T
			log       = logger.FromContext(c.UserContext())
		)

		if err := c.ParamsParser(&pathParam); err != nil {
			return response.HandleWarning(c, log, fiber.StatusBadRequest, "invalid request path param", err)
		}

		errs, err := v.ValidateStruct(pathParam)
		if err != nil {
			if err.Error() != "validation failed" {
				return response.HandleWarning(c, log, fiber.StatusBadRequest, "failed to validate request path param", err)
			}
		}
		if len(errs) > 0 {
			return response.HandleWarnings(c, log, fiber.StatusBadRequest, "failed to validate request path param", errs)
		}

		c.Locals("path_param", pathParam)
		return c.Next()
	}
}
