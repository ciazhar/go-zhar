package user

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/rest_api_service/internal/model/response"
	"github.com/rs/zerolog/log"
)

func (r userRepository) GetUsers(ctx context.Context, page, limit int) ([]response.User, int64, error) {
	var (
		logger = log.Ctx(ctx).With().Int("page", page).Int("limit", limit).Logger()
	)

	logger.Info().Msg("Fetching users from DB")

	// Dummy data
	users := []response.User{
		{ID: "4f372c6c-9272-4062-9261-dccda33ecbb9", Name: "Alice", Age: 25},
		{ID: "4dbcf43b-9471-44a0-aba7-1d26c2b27394", Name: "Bob", Age: 30},
		{ID: "ff8a8538-2c00-44a4-914a-87c47c65b21b", Name: "Charlie", Age: 35},
	}

	total := int64(len(users))

	// Simulate pagination (no real DB limit/offset here)
	start := (page - 1) * limit
	end := start + limit
	if start >= len(users) {
		return []response.User{}, total, nil
	}
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], total, nil
}
