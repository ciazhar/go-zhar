package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type BitmapDailyLoginTrackingRepository struct {
	redis *redis.Client
}

func NewBitmapDailyLoginTrackingRepository(redisClient *redis.Client) *BitmapDailyLoginTrackingRepository {
	return &BitmapDailyLoginTrackingRepository{
		redis: redisClient,
	}
}

// MarkUserLogin Mark a user as logged in for the day
func (r *BitmapDailyLoginTrackingRepository) MarkUserLogin(ctx context.Context, date string, userID int64) error {
	key := fmt.Sprintf("daily:login:%s", date)
	return r.redis.SetBit(ctx, key, userID, 1).Err()
}

// CheckUserLogin Check if a user logged in on a specific day
func (r *BitmapDailyLoginTrackingRepository) CheckUserLogin(ctx context.Context, date string, userID int64) (int64, error) {
	key := fmt.Sprintf("daily:login:%s", date)
	return r.redis.GetBit(ctx, key, userID).Result()
}
