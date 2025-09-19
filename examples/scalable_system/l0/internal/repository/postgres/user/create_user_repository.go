package user

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar-scalable-system-l0/internal/model/request"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r UserRepository) CreateUser(ctx context.Context, req request.CreateUserBodyRequest) error {
	var (
		reqCtx, span = otel.Tracer("repository").Start(ctx, "UserRepository.CreateUser")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Any("req", req).Logger()
	)
	defer deferFn()

	_, err := r.pg.Exec(reqCtx, queryCreateUser,
		req.Username,
		req.Email,
		req.Password,
		req.FullName,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert user")
		return err
	}

	return nil
}
