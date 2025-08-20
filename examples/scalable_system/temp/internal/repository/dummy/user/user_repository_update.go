package user

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/scalable_system/temp/internal/model/request"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r userRepository) UpdateUser(ctx context.Context, id string, req request.UpdateUserBodyRequest) error {
	var (
		log = logger.FromContext(ctx).With().Str("id", id).Any("req", req).Logger()
	)

	log.Info().Msg("updating user in DB")

	return nil // or return error if needed
}
