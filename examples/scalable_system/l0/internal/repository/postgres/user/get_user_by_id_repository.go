package user

import (
	"context"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/response"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r UserRepository) GetUserByID(ctx context.Context, id string) (*response.User, error) {
	var (
		log  = logger.FromContext(ctx).With().Str("id", id).Logger()
		resp = new(response.User)
	)

	err := r.pg.QueryRow(ctx, queryGetUserByID, id).
		Scan(&resp.ID, &resp.Username, &resp.Email, &resp.FullName, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user by ID")
		return nil, err
	}

	return resp, nil
}
