package user

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r UserRepository) SoftDeleteUser(ctx context.Context, id string) error {
	var (
		reqCtx, span = otel.Tracer("repository").Start(ctx, "UserRepository.SoftDeleteUser")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Str("id", id).Logger()
	)
	defer deferFn()

	cmdTag, err := r.pg.Exec(reqCtx, querySoftDeleteUser, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete user")
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		log.Error().Msg("no rows deleted")
		return fmt.Errorf("no rows deleted")
	}

	return nil
}
