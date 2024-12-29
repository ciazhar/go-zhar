package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type GeoTrackingDeliveryAgentLocationRepository struct {
	redis *redis.Client
}

func NewGeoTrackingDeliveryAgentLocationRepository(redisClient *redis.Client) *GeoTrackingDeliveryAgentLocationRepository {
	return &GeoTrackingDeliveryAgentLocationRepository{
		redis: redisClient,
	}
}

// AddAgentLocation Add agent location
func (r *GeoTrackingDeliveryAgentLocationRepository) AddAgentLocation(ctx context.Context, agentID string, lat, long float64) error {
	return r.redis.GeoAdd(ctx, "agent:locations", &redis.GeoLocation{
		Name:      agentID,
		Latitude:  lat,
		Longitude: long,
	}).Err()
}

// GetNearbyAgents Get nearby agents within a 5km radius
func (r *GeoTrackingDeliveryAgentLocationRepository) GetNearbyAgents(ctx context.Context, lat, long float64) ([]redis.GeoLocation, error) {
	return r.redis.GeoRadius(ctx, "agent:locations", long, lat, &redis.GeoRadiusQuery{
		Radius:    5,
		Unit:      "km",
		WithCoord: true,
	}).Result()
}
