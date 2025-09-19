package user

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar-scalable-system-l0/internal/model/response"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (u userService) GetUserByID(ctx context.Context, id string) (*response.User, error) {
	var (
		reqCtx, span = otel.Tracer("service").Start(ctx, "UserService.GetUserByID")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Str("id", id).Logger()
	)
	defer deferFn()

	user, err := u.repo.GetUserByID(reqCtx, id)
	if err != nil {
		log.Err(err).Msg("failed to get user by ID")
		return nil, err
	}

	return user, nil
}
