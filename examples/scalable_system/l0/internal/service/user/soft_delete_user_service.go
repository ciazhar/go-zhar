package user

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (u userService) SoftDeleteUser(ctx context.Context, id string) error {
	var (
		reqCtx, span = otel.Tracer("service").Start(ctx, "UserService.SoftDeleteUser")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Str("id", id).Logger()
	)
	defer deferFn()

	if err := u.repo.SoftDeleteUser(reqCtx, id); err != nil {
		log.Err(err).Msg("failed to delete user")
		return err
	}

	return nil
}
