package user

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/rest_api_service/internal/model/response"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

func (u userService) GetUserByID(ctx context.Context, id string) (*response.User, error) {
	var (
		log = logger.FromContext(ctx).With().Str("id", id).Logger()
	)

	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		log.Err(err).Msg("failed to get user by ID")
		return nil, err
	}

	return user, nil
}
