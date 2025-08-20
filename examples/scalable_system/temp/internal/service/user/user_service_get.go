package user

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/scalable_system/temp/internal/model/request"
	"github.com/ciazhar/go-zhar/examples/scalable_system/temp/internal/model/response"

	"github.com/rs/zerolog/log"
)

func (u userService) GetUsers(ctx context.Context, req request.GetUsersQueryParam) ([]response.User, int64, error) {
	var (
		logger = log.Ctx(ctx).With().Interface("req", req).Logger()
	)

	users, total, err := u.repo.GetUsers(ctx, req.Page, req.Size)
	if err != nil {
		logger.Err(err).Msg("failed to fetch users")
		return nil, 0, err
	}

	return users, total, nil
}
