package user

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar-scalable-system-l0/internal/model/request"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r UserRepository) UpdateUser(ctx context.Context, id string, req request.UpdateUserBodyRequest) error {
	var (
		reqCtx, span = otel.Tracer("repository").Start(ctx, "UserRepository.UpdateUser")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Str("id", id).Any("req", req).Logger()
	)
	defer deferFn()

	cmdTag, err := r.pg.Exec(reqCtx, queryUpdateUser, req.Username, req.Email, req.FullName, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to update user")
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		log.Error().Msg("no rows updated")
		return fmt.Errorf("no rows updated")
	}

	return nil
}
