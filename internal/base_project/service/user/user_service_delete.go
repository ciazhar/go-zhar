package user

import (
	"context"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

func (u userService) DeleteUser(ctx context.Context, id string) error {
	var (
		log = logger.FromContext(ctx).With().Str("id", id).Logger()
	)

	if err := u.repo.DeleteUser(ctx, id); err != nil {
		log.Err(err).Msg("failed to delete user")
		return err
	}

	return nil
}
