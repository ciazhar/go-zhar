package user

import (
	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func (uc *UserController) CreateUser(ctx *fiber.Ctx) error {
	var (
		reqCtx, span = otel.Tracer("controller").Start(ctx.UserContext(), "UserController.CreateUser")
		deferFn      = func() { span.End() }
		body         = ctx.Locals("body").(request.CreateUserBodyRequest)
		log          = logger.FromContext(reqCtx).With().Any("body", body).Logger()
	)
	defer deferFn()

	if err := uc.service.CreateUser(reqCtx, body); err != nil {
		log.Err(err).Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.NewErrorResponse(reqCtx, "failed to insert user to DB"))
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.NewBaseResponse("Create user success"))
}
