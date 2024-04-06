package mongo

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init(host string, port int, username string, password string, database string, logger logger.Logger) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if username != "" && password != "" {

		url := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", username, password, host, port, database)
		connect, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
		if err != nil {
			logger.Fatalf("failed to connect to mongo: %v", err)
		}

		return connect.Database(database)
	} else {

		url := fmt.Sprintf("mongodb://%s:%d/%s", host, port, database)
		connect, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
		if err != nil {
			logger.Fatalf("failed to connect to mongo: %v", err)
		}

		return connect.Database(database)
	}
}
