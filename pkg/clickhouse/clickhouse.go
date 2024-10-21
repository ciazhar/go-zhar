package clickhouse

import (
	"context"
	"fmt"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

// [clickhouse][conn=1][127.0.0.1:9000][send query] compression="none"
// SELECT UserID, COUNT(*) AS transaction_count
// FROM transactions
// WHERE Timestamp >= now() - INTERVAL 30 DAY
// GROUP BY UserID
// HAVING transaction_count > 10
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
			if strings.Contains(format, "send query") {
				logger.LogDebug(context.Background(), fmt.Sprintf(format, v...), nil)
			}
		},
	})
	if err != nil {
		logger.LogFatal(context.Background(), err, "ClickHouse connection failed", nil)
	}
	logger.LogInfo(context.Background(), "ClickHouse connection initialized", nil)
	return conn
}
