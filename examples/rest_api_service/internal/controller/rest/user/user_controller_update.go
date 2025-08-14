package user

import (
	"github.com/ciazhar/go-zhar/examples/rest_api_service/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func (uc *UserController) UpdateUser(ctx *fiber.Ctx) error {

	var (
		body = ctx.Locals("body").(request.UpdateUserBodyRequest)
		path = ctx.Locals("path_param").(request.UserPathParam)
		log  = logger.FromContext(ctx.UserContext()).With().Any("body", body).Any("path", path).Logger()
	)

	if err := uc.service.UpdateUser(ctx.UserContext(), path.ID, body); err != nil {
		log.Err(err).Send()
		return ctx.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse("failed to update user", err))
	}
	return ctx.Status(fiber.StatusOK).JSON(response.NewBaseResponse("Update user success"))
}
