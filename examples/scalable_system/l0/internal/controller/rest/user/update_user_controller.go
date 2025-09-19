package user

import (
	"github.com/ciazhar/go-zhar-scalable-system-l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func (uc *UserController) UpdateUser(ctx *fiber.Ctx) error {
	var (
		reqCtx, span = otel.Tracer("controller").Start(ctx.UserContext(), "UserController.UpdateUser")
		deferFn      = func() { span.End() }
		body         = ctx.Locals("body").(request.UpdateUserBodyRequest)
		path         = ctx.Locals("path_param").(request.UserPathParam)
		log          = logger.FromContext(reqCtx).With().Any("body", body).Any("path", path).Logger()
	)
	defer deferFn()

	if err := uc.service.UpdateUser(reqCtx, path.ID, body); err != nil {
		log.Err(err).Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.NewErrorResponse(reqCtx, "failed to update user"))
	}
	return ctx.Status(fiber.StatusOK).JSON(response.NewBaseResponse("Update user success"))
}
