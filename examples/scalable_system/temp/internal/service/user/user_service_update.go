package user

import (
	"context"

	"github.com/ciazhar/go-zhar/examples/scalable_system/temp/internal/model/request"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (u userService) UpdateUser(ctx context.Context, id string, req request.UpdateUserBodyRequest) error {
	var (
		log = logger.FromContext(ctx).With().Str("id", id).Any("req", req).Logger()
	)

	if err := u.repo.UpdateUser(ctx, id, req); err != nil {
		log.Err(err).Msg("failed to update user")
		return err
	}

	return nil
}
