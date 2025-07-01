package user

import (
	"github.com/ciazhar/go-start-small/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	reqID := c.Locals(middleware.RequestIDKey).(string)

	// Passing request ID to service
	if err := uc.service.CreateUser(c.Context(), reqID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}
