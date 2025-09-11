package user

import (
	"context"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (u userService) IsUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var (
		log = logger.FromContext(ctx).With().Str("email", email).Logger()
	)

	exists, err := u.repo.IsUserExistsByEmail(ctx, email)
	if err != nil {
		log.Err(err).Msg("failed to check if user exists by email")
		return false, err
	}

	return exists, nil
}
