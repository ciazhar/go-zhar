package user

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar-scalable-system-l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r *UserRepository) UpsertUserByID(ctx context.Context, req request.UpsertUserBodyRequest) error {
	var (
		reqCtx, span = otel.Tracer("repository").Start(ctx, "UserRepository.UpsertUserByID")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Any("req", req).Logger()
	)
	defer deferFn()

	if _, err := r.pg.Exec(reqCtx, queryUpsertUser, req.Id, req.Username, req.Email, req.Password, req.FullName); err != nil {
		log.Error().Err(err).Msg("failed to upsert user")
		return err
	}

	return nil
}
