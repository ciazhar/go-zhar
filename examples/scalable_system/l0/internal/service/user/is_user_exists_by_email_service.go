package user

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (u userService) IsUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var (
		reqCtx, span = otel.Tracer("service").Start(ctx, "UserService.IsUserExistsByEmail")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Str("email", email).Logger()
	)
	defer deferFn()

	exists, err := u.repo.IsUserExistsByEmail(reqCtx, email)
	if err != nil {
		log.Err(err).Msg("failed to check if user exists by email")
		return false, err
	}

	return exists, nil
}
