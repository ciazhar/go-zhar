package user

import (
	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func (uc *UserController) UpsertUser(ctx *fiber.Ctx) error {
	var (
		reqCtx, span = otel.Tracer("controller").Start(ctx.UserContext(), "UserController.UpsertUser")
		deferFn      = func() { span.End() }
		body         = ctx.Locals("body").(request.UpsertUserBodyRequest)
		log          = logger.FromContext(reqCtx).With().Any("body", body).Logger()
	)
	defer deferFn()

	if err := uc.service.UpsertUserByID(reqCtx, body); err != nil {
		log.Err(err).Send()
		return ctx.Status(fiber.StatusBadRequest).
			JSON(response.NewErrorResponse(reqCtx, "failed to upsert user"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(response.NewBaseResponse("Upsert user success"))
}
