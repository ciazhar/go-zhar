package postgres

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Init(username string, password string, host string, port int, database string, schema string, debug bool, logger *logger.Logger) *pgxpool.Pool {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?search_path=%s", username, password, host, port, database, schema)

	c, err := pgxpool.ParseConfig(url)
	if err != nil {
		logger.Fatalf("Failed to parse postgres connection string: %v", err)
	}

	if debug {
		c.ConnConfig.Logger = zerologadapter.NewLogger(*logger.GetServiceLogger())
	}

	conn, err := pgxpool.ConnectConfig(context.Background(), c)
	if err != nil {
		logger.Fatalf("Failed to connect to postgres: %v", err)
	}

	return conn
}
