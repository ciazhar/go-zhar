package postgres

import (
	"context"
	"fmt"

	"github.com/ciazhar/go-zhar/pkg/bootstrap"
	"github.com/ciazhar/go-zhar/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

// ZerologAdapter adapts zerolog to pgx logger interface
type ZerologAdapter struct {
	zerolog.Logger
}

func (l *ZerologAdapter) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	zerologLevel := zerolog.DebugLevel

	switch level {
	case tracelog.LogLevelTrace:
		zerologLevel = zerolog.TraceLevel
	case tracelog.LogLevelDebug:
		zerologLevel = zerolog.DebugLevel
	case tracelog.LogLevelInfo:
		zerologLevel = zerolog.InfoLevel
	case tracelog.LogLevelWarn:
		zerologLevel = zerolog.WarnLevel
	case tracelog.LogLevelError:
		zerologLevel = zerolog.ErrorLevel
	}

	l.WithLevel(zerologLevel).Fields(data).Msg(msg)
}

func InitPostgres(ctx context.Context, host string, port int, dbName string, username string, password string, logLevel string) (*bootstrap.ClientService, *pgxpool.Pool) {
	var (
		log = logger.FromContext(ctx).With().
			Str("host", host).
			Int("port", port).
			Str("db_name", dbName).
			Logger()
	)

	// Connection string
	connString := ""
	if username == "" && password == "" {
		connString = fmt.Sprintf("postgresql://%s:%d/%s?sslmode=disable", host, port, dbName)
	} else {
		connString = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbName)
	}

	// Configure the connection pool with debug logging
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse connection string")
	}

	// Configure the connection pool with debug logging
	if logLevel == "debug" {
		pgxLogger := &ZerologAdapter{Logger: logger.GetLogger()}
		config.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   pgxLogger,
			LogLevel: tracelog.LogLevelDebug,
		}
	}

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create connection pool")
	}

	return bootstrap.NewClientService(ctx, "postgres", func() error {
		pool.Close()
		return nil
	}), pool
}
