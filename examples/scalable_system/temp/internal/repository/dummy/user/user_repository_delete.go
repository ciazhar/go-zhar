package user

import (
	"context"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r userRepository) DeleteUser(ctx context.Context, id string) error {
	var (
		log = logger.FromContext(ctx).With().Str("id", id).Logger()
	)

	log.Info().Msg("deleting user in DB")

	return nil
}
