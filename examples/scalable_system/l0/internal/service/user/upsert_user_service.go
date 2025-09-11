package user

import (
	"context"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (u userService) UpsertUserByID(ctx context.Context, req request.UpsertUserBodyRequest) error {
	var (
		log = logger.FromContext(ctx).With().Any("req", req).Logger()
	)

	if err := u.repo.UpsertUserByID(ctx, req); err != nil {
		log.Err(err).Msg("failed to upsert user by ID")
		return err
	}

	return nil
}
