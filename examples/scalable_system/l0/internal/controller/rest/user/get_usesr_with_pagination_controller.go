package user

import (
	"github.com/ciazhar/go-zhar-scalable-system-l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func (uc *UserController) GetUsers(ctx *fiber.Ctx) error {
	var (
		reqCtx, span = otel.Tracer("controller").Start(ctx.UserContext(), "UserController.GetUsers")
		deferFn      = func() { span.End() }
		query        = ctx.Locals("query_param").(request.GetUsersQueryParam)
		log          = logger.FromContext(reqCtx).With().Interface("query", query).Logger()
	)
	defer deferFn()

	users, total, err := uc.service.GetUsersWithPagination(reqCtx, query.Page, query.Size)
	if err != nil {
		log.Err(err).Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.NewErrorResponse(reqCtx, "failed to get users from DB"))
	}

	return ctx.Status(fiber.StatusOK).JSON(response.NewPageResponse("Get users success", users, total))
}
