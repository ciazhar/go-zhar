package user

import (
	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func (uc *UserController) DeleteUser(ctx *fiber.Ctx) error {
	var (
		reqCtx, span = otel.Tracer("controller").Start(ctx.UserContext(), "UserController.DeleteUser")
		deferFn      = func() { span.End() }
		path         = ctx.Locals("path_param").(request.UserPathParam)
		log          = logger.FromContext(reqCtx).With().Any("path", path).Logger()
	)
	defer deferFn()

	if err := uc.service.SoftDeleteUser(reqCtx, path.ID); err != nil {
		log.Err(err).Send()
		return ctx.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(reqCtx, "failed to delete user"))
	}
	return ctx.Status(fiber.StatusOK).JSON(response.NewBaseResponse("Delete user success"))
}
