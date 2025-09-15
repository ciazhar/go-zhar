package user

import (
	"github.com/ciazhar/go-zhar/examples/scalable_system/temp/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func (uc *UserController) GetUsers(ctx *fiber.Ctx) error {
	var (
		query = ctx.Locals("query_param").(request.GetUsersQueryParam)
		log   = logger.FromContext(ctx.UserContext()).With().Interface("query", query).Logger()
	)

	users, total, err := uc.service.GetUsers(ctx.UserContext(), query)
	if err != nil {
		log.Err(err).Send()
		return ctx.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse("failed to get users from DB"))
	}

	return ctx.Status(fiber.StatusOK).JSON(response.NewPageResponse("Get users success", users, total))
}
