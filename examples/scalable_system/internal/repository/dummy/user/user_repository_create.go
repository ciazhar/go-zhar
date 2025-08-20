package user

import (
	"context"

	"github.com/ciazhar/go-zhar/examples/scalable_system/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (u userRepository) CreateUser(ctx context.Context, req request.CreateUserBodyRequest) error {
	var (
		log = logger.FromContext(ctx).With().Any("req", req).Logger()
	)

	log.Info().Msg("inserting user to DB")

	return nil
}
