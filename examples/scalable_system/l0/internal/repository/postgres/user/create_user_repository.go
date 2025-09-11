package user

import (
	"context"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r UserRepository) CreateUser(ctx context.Context, req request.CreateUserBodyRequest) error {
	var (
		log = logger.FromContext(ctx).With().Any("req", req).Logger()
	)

	_, err := r.pg.Exec(ctx, queryCreateUser,
		req.Username,
		req.Email,
		req.Password,
		req.FullName,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert user")
		return err
	}

	return nil
}
