package user

import (
	"github.com/ciazhar/go-zhar/examples/scalable_system/temp/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func (uc *UserController) GetUserByID(ctx *fiber.Ctx) error {
	var (
		path = ctx.Locals("path_param").(request.UserPathParam)
		log  = logger.FromContext(ctx.UserContext()).With().Any("path", path).Logger()
	)

	user, err := uc.service.GetUserByID(ctx.UserContext(), path.ID)
	if err != nil {
		log.Err(err).Send()
		return ctx.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse("failed to get user by ID", err))
	}
	return ctx.Status(fiber.StatusOK).JSON(response.NewDataResponse("Get user by ID success", user))
}
