package user

import (
	"context"
	"fmt"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r UserRepository) UpdateUser(ctx context.Context, id string, req request.UpdateUserBodyRequest) error {
	var (
		log = logger.FromContext(ctx).With().Str("id", id).Any("req", req).Logger()
	)

	cmdTag, err := r.pg.Exec(ctx, queryUpdateUser, req.Username, req.Email, req.FullName, id)
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
