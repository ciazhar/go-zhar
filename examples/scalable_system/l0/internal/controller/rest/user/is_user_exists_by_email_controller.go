package user

import (
	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func (uc *UserController) IsUserExistByEmail(ctx *fiber.Ctx) error {
	var (
		query = ctx.Locals("query_param").(request.UserEmailQueryParam)
		log   = logger.FromContext(ctx.UserContext()).With().Interface("query", query).Logger()
	)

	exists, err := uc.service.IsUserExistsByEmail(ctx.UserContext(), query.Email)
	if err != nil {
		log.Err(err).Send()
		return ctx.Status(fiber.StatusBadRequest).
			JSON(response.NewErrorResponse("failed to check user existence", err))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(response.NewDataResponse("Check user existence success", map[string]bool{"exists": exists}))
}
