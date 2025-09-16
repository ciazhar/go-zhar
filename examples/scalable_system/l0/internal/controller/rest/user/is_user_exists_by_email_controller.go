package user

import (
	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func (uc *UserController) IsUserExistByEmail(ctx *fiber.Ctx) error {
	var (
		reqCtx, span = otel.Tracer("controller").Start(ctx.UserContext(), "UserController.IsUserExistByEmail")
		deferFn      = func() { span.End() }
		query        = ctx.Locals("query_param").(request.UserEmailQueryParam)
		log          = logger.FromContext(reqCtx).With().Interface("query", query).Logger()
	)
	defer deferFn()

	exists, err := uc.service.IsUserExistsByEmail(reqCtx, query.Email)
	if err != nil {
		log.Err(err).Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.NewErrorResponse(reqCtx, "failed to check user existence"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(response.NewDataResponse("Check user existence success", map[string]bool{"exists": exists}))
}
