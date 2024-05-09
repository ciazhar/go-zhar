package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ciazhar/go-zhar/pkg/logger"
)

func Init(hosts string, database string, username string, password string, debug bool, logger logger.Logger) clickhouse.Conn {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{hosts},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		Debug: debug,
	})
	if err != nil {
		logger.Fatalf("Failed to initialize ClickHouse connection : %v", err)
	}
	logger.Info("ClickHouse connection initialized successfully")
	return conn
}
