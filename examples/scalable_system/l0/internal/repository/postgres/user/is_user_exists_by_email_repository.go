package user

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r *UserRepository) IsUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var (
		reqCtx, span = otel.Tracer("repository").Start(ctx, "UserRepository.IsUserExistsByEmail")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Str("email", email).Logger()
		exists       = false
	)
	defer deferFn()

	err := r.pg.QueryRow(reqCtx, queryIsUserExistsByEmail, email).Scan(&exists)
	if err != nil {
		log.Error().Err(err).Msg("failed to check if user exists")
		return false, err
	}

	return exists, nil
}
