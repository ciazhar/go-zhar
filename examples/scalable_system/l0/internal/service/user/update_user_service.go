package user

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar-scalable-system-l0/internal/model/request"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (u userService) UpdateUser(ctx context.Context, id string, req request.UpdateUserBodyRequest) error {
	var (
		reqCtx, span = otel.Tracer("service").Start(ctx, "UserService.UpdateUser")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Str("id", id).Any("req", req).Logger()
	)
	defer deferFn()

	if err := u.repo.UpdateUser(reqCtx, id, req); err != nil {
		log.Err(err).Msg("failed to update user")
		return err
	}

	return nil
}
