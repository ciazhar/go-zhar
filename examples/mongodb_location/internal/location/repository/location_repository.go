package repository

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/mongodb_location/internal/location/model"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type LocationRepository struct {
	c *mongo.Collection
}

func (l *LocationRepository) Insert(location model.InsertLocationForm) (err error) {
	_, err = l.c.InsertOne(context.Background(), location)
	return
}

func (l *LocationRepository) Nearest(long, lat float64, maxDistance int, limit int) (res []model.Location, err error) {

	// Define a point for which we want to find the nearest location
	targetPoint := []float64{long, lat}

	// Create a pipeline to calculate distances and find the nearest location
	pipeline := []bson.M{
		{
			"$geoNear": bson.M{
				"near": bson.M{
					"type":        "Point",
					"coordinates": targetPoint,
				},
				"distanceField": "distance",
				"spherical":     true,
				"maxDistance":   maxDistance, // Maximum distance in meters
			},
		},
		{"$limit": limit}, // Limit to the nearest location
		{"$project": bson.M{"name": 1, "coordinate": 1, "distance": 1}}, // Project only necessary fields
	}

	// Execute the aggregation pipeline
	cursor, err := l.c.Aggregate(context.Background(), pipeline)
	if err != nil {
		return
	}

	// Iterate over the cursor to get the results
	if err = cursor.All(context.Background(), &res); err != nil {
		return
	}

	return
}

func NewLocationRepository(conn *mongo.Database) *LocationRepository {

	// Get Collection
	coll := conn.Collection("locations")

	// Define the index model
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "coordinate", Value: "2dsphere"}},
		Options: options.Index().SetName("coordinate_2dsphere"),
	}

	// Create the index if not exists
	if _, err := coll.Indexes().CreateOne(
		context.Background(),
		indexModel,
		options.CreateIndexes().SetMaxTime(10*time.Second), // Optional: Set a timeout for the index creation
	); err != nil {
		// Check if the error is due to the index already existing
		if !strings.Contains(err.Error(), "index already exists") {
			logger.LogFatal(context.Background(), err, "Failed to create index", nil)
		}
	} else {
		logger.LogInfo(context.Background(), "2dsphere index on 'coordinate' field created successfully.", nil)
	}

	return &LocationRepository{
		c: coll,
	}
}
