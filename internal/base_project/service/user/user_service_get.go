package user

import (
	"context"
	"github.com/ciazhar/go-start-small/internal/base_project/model/request"
	"github.com/ciazhar/go-start-small/internal/base_project/model/response"
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
