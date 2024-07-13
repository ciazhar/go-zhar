package mongo

import (
	"context"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init(host string, username string, password string, database string, logger *logger.Logger) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := generateUrl(username, password, host, database)
	connect, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.Fatalf("failed to connect to mongo: %v", err)
	}

	err = connect.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("MongoDB connection initialized successfully")

	return connect.Database(database)
}

func generateUrl(username string, password string, hosts string, database string) string {
	var sb strings.Builder

	// MongoDB scheme
	sb.WriteString("mongodb://")

	// Username and password
	if username != "" && password != "" {
		sb.WriteString(url.QueryEscape(username))
		sb.WriteString(":")
		sb.WriteString(url.QueryEscape(password))
		sb.WriteString("@")
	}

	// Hosts
	hostsArr := strings.Split(hosts, ",")
	sb.WriteString(hostsArr[0])

	// Database
	if database != "" {
		sb.WriteString("/")
		sb.WriteString(database)
	}

	q := sb.String()

	// Determine if it's a cluster
	if len(hostsArr) > 1 {
		q += "?connect=direct"
	}

	return q
}
