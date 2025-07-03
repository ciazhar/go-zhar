package user

import (
	"context"
	"github.com/ciazhar/go-start-small/internal/base_project/model/request"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

func (u userService) CreateUser(ctx context.Context, req request.CreateUserBodyRequest) error {
	var (
		log = logger.FromContext(ctx).With().Any("req", req).Logger()
	)

	err := u.repo.CreateUser(ctx, req)
	if err != nil {
		log.Err(err).Msg("failed to insert user to DB")
		return err
	}

	return nil
}
