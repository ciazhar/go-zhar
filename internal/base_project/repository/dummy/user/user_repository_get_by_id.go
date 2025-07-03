package user

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/internal/base_project/model/response"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

func (r userRepository) GetUserByID(ctx context.Context, id string) (*response.User, error) {
	var (
		log = logger.FromContext(ctx).With().Str("id", id).Logger()
	)

	log.Info().Msg("retrieving user from DB")

	// Dummy simulation
	if id == "not-found" {
		return nil, fmt.Errorf("user not found")
	}

	return &response.User{ID: id, Name: "yabai", Age: 10}, nil
}
