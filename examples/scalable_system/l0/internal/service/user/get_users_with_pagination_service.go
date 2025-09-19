package user

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar-scalable-system-l0/internal/model/response"

	"github.com/rs/zerolog/log"
)

func (u userService) GetUsersWithPagination(ctx context.Context, page, limit int) ([]response.User, int64, error) {
	var (
		reqCtx, span = otel.Tracer("service").Start(ctx, "UserService.GetUsersWithPagination")
		deferFn      = func() { span.End() }
		logger       = log.Ctx(reqCtx).With().Int("page", page).Int("limit", limit).Logger()
	)
	defer deferFn()

	users, total, err := u.repo.GetUsersWithPagination(reqCtx, page, limit)
	if err != nil {
		logger.Err(err).Msg("failed to fetch users")
		return nil, 0, err
	}

	return users, total, nil
}
