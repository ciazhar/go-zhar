package user

import (
	"context"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (u userService) SoftDeleteUser(ctx context.Context, id string) error {
	var (
		log = logger.FromContext(ctx).With().Str("id", id).Logger()
	)

	if err := u.repo.SoftDeleteUser(ctx, id); err != nil {
		log.Err(err).Msg("failed to delete user")
		return err
	}

	return nil
}
