package clickhouse

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

func Init(hosts string, database string, username string, password string, debug bool) clickhouse.Conn {

	logger.LogInfo(context.Background(), "ClickHouse connection initializing", nil)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{hosts},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		Debug: debug,
		Debugf: func(format string, v ...any) {
			logger.LogDebug(context.Background(), format, toMap(v))
		},
	})
	if err != nil {
		logger.LogFatal(context.Background(), err, "ClickHouse connection failed", nil)
	}
	logger.LogInfo(context.Background(), "ClickHouse connection initialized", nil)
	return conn
}

func toMap(v ...any) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(v); i += 2 {
		if i+1 < len(v) {
			key, ok := v[i].(string)
			if ok {
				m[key] = v[i+1]
			}
		}
	}
	return m
}
