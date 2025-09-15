package user

import (
	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func (uc *UserController) GetUserByID(ctx *fiber.Ctx) error {
	var (
		reqCtx, span = otel.Tracer("controller").Start(ctx.UserContext(), "UserController.GetUserByID")
		deferFn      = func() { span.End() }
		path         = ctx.Locals("path_param").(request.UserPathParam)
		log          = logger.FromContext(reqCtx).With().Any("path", path).Logger()
	)
	defer deferFn()

	user, err := uc.service.GetUserByID(reqCtx, path.ID)
	if err != nil {
		log.Err(err).Send()
		return ctx.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(reqCtx, "failed to get user by ID"))
	}
	return ctx.Status(fiber.StatusOK).JSON(response.NewDataResponse("Get user by ID success", user))
}
