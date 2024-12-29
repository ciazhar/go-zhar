package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type SetActiveDeliveryAgentRepository struct {
	redis *redis.Client
}

func NewSetActiveDeliveryAgentRepository(redisClient *redis.Client) *SetActiveDeliveryAgentRepository {
	return &SetActiveDeliveryAgentRepository{
		redis: redisClient,
	}
}

// AddActiveAgent Add a delivery agent to the active list
func (r *SetActiveDeliveryAgentRepository) AddActiveAgent(ctx context.Context, agentID string) error {
	return r.redis.SAdd(ctx, "active:agents", agentID).Err()
}

// IsAgentActive Check if an agent is active
func (r *SetActiveDeliveryAgentRepository) IsAgentActive(ctx context.Context, agentID string) (bool, error) {
	return r.redis.SIsMember(ctx, "active:agents", agentID).Result()
}
