package user

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (u userService) UpsertUserByID(ctx context.Context, req request.UpsertUserBodyRequest) error {
	var (
		reqCtx, span = otel.Tracer("service").Start(ctx, "UserService.UpsertUserByID")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Any("req", req).Logger()
	)
	defer deferFn()

	if err := u.repo.UpsertUserByID(reqCtx, req); err != nil {
		log.Err(err).Msg("failed to upsert user by ID")
		return err
	}

	return nil
}
