package user

import (
	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func (uc *UserController) CreateUser(ctx *fiber.Ctx) error {
	var (
		body = ctx.Locals("body").(request.CreateUserBodyRequest)
		log  = logger.FromContext(ctx.UserContext()).With().Any("body", body).Logger()
	)

	if err := uc.service.CreateUser(ctx.UserContext(), body); err != nil {
		log.Err(err).Send()
		return ctx.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse("failed to insert user to DB", err))
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.NewBaseResponse("Create user success"))
}
