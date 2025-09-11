package user

import (
	"context"
	"fmt"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r UserRepository) SoftDeleteUser(ctx context.Context, id string) error {
	var (
		log = logger.FromContext(ctx).With().Str("id", id).Logger()
	)

	cmdTag, err := r.pg.Exec(ctx, querySoftDeleteUser, id)
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
