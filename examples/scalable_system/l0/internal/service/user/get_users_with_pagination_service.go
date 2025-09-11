package user

import (
	"context"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/response"

	"github.com/rs/zerolog/log"
)

func (u userService) GetUsersWithPagination(ctx context.Context, page, limit int) ([]response.User, int64, error) {
	var (
		logger = log.Ctx(ctx).With().Int("page", page).Int("limit", limit).Logger()
	)

	users, total, err := u.repo.GetUsersWithPagination(ctx, page, limit)
	if err != nil {
		logger.Err(err).Msg("failed to fetch users")
		return nil, 0, err
	}

	return users, total, nil
}
