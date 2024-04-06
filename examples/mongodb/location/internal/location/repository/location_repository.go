package repository

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/mongodb/location/internal/location/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"time"
)

type LocationRepository interface {
	Insert(location model.InsertLocationForm) error
	Nearest(long, lat float64, maxDistance int, limit int) ([]model.Location, error)
}

type locationRepository struct {
	c *mongo.Collection
}

func (l locationRepository) Insert(location model.InsertLocationForm) error {
	_, err := l.c.InsertOne(context.Background(), location)
	return err
}

func (l locationRepository) Nearest(long, lat float64, maxDistance int, limit int) ([]model.Location, error) {

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
		return nil, err
	}

	// Iterate over the cursor to get the results
	var results []model.Location
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}

	return results, nil
}

func NewLocationRepository(conn *mongo.Database) LocationRepository {

	// Get Collection
	coll := conn.Collection("locations")

	// Define the index model
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"coordinate", "2dsphere"}},
		Options: options.Index().SetName("coordinate_2dsphere").SetBackground(true),
	}

	// Create the index if not exists
	_, err := coll.Indexes().CreateOne(
		context.Background(),
		indexModel,
		options.CreateIndexes().SetMaxTime(10*time.Second), // Optional: Set a timeout for the index creation
	)
	if err != nil {
		// Check if the error is due to the index already existing
		if !strings.Contains(err.Error(), "index already exists") {
			log.Fatal(err)
		}
		log.Println("Index already exists. Skipping index creation.")
	} else {
		log.Println("2dsphere index on 'coordinate' field created successfully.")
	}

	return &locationRepository{
		c: coll,
	}
}
