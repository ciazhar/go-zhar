package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type SortedSetRealtimeDeliveryMetricRepository struct {
	redis *redis.Client
}

func NewSortedSetRealtimeDeliveryMetricRepository(redisClient *redis.Client) *SortedSetRealtimeDeliveryMetricRepository {
	return &SortedSetRealtimeDeliveryMetricRepository{
		redis: redisClient,
	}
}

// AddDeliveryMetric Add a delivery record with time taken
func (r *SortedSetRealtimeDeliveryMetricRepository) AddDeliveryMetric(ctx context.Context, orderID string, timeTaken float64) error {
	return r.redis.ZAdd(ctx, "delivery:metrics", &redis.Z{
		Score:  timeTaken,
		Member: orderID,
	}).Err()
}

// GetTopDeliveries Get top 5 fastest deliveries
func (r *SortedSetRealtimeDeliveryMetricRepository) GetTopDeliveries(ctx context.Context) ([]redis.Z, error) {
	return r.redis.ZRangeWithScores(ctx, "delivery:metrics", 0, 4).Result()
}
