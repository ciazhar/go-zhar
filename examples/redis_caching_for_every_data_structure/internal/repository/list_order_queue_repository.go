package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type ListOrderQueueRepository struct {
	redis *redis.Client
}

func NewListOrderQueueRepository(redisClient *redis.Client) *ListOrderQueueRepository {
	return &ListOrderQueueRepository{
		redis: redisClient,
	}
}

// AddOrderToQueue Add an order to the queue
func (r *ListOrderQueueRepository) AddOrderToQueue(ctx context.Context, orderID string) error {
	return r.redis.RPush(ctx, "order:queue", orderID).Err()
}

// ProcessNextOrder Process the next order in the queue
func (r *ListOrderQueueRepository) ProcessNextOrder(ctx context.Context) (string, error) {
	return r.redis.LPop(ctx, "order:queue").Result()
}
