package user

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/internal/base_project/model/request"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

func (u userRepository) CreateUser(ctx context.Context, req request.CreateUserBodyRequest) error {
	var (
		log = logger.FromContext(ctx).With().Any("req", req).Logger()
	)

	log.Info().Msg("inserting user to DB")
	// Simulasi insert ke database
	return fmt.Errorf("failed to insert user to DB")
}
