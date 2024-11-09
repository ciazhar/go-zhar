package mongo

import (
	"context"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB(host string, username string, password string, database string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := generateUrl(username, password, host, database)
	connect, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.LogFatal(ctx, err, "failed to connect to mongo", nil)
	}

	err = connect.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.LogFatal(ctx, err, "failed to ping mongo", nil)
	}

	logger.LogInfo(ctx, "connected to mongo", nil)

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
