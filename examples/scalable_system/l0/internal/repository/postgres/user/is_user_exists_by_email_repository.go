package user

import (
	"context"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r *UserRepository) IsUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var (
		log    = logger.FromContext(ctx).With().Str("email", email).Logger()
		exists = false
	)

	err := r.pg.QueryRow(ctx, queryIsUserExistsByEmail, email).Scan(&exists)
	if err != nil {
		log.Error().Err(err).Msg("failed to check if user exists")
		return false, err
	}

	return exists, nil
}
