package logger

import (
	"context"

	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/getsentry/sentry-go"
)

func InitSentry(dsn string) {

	err := sentry.Init(
		sentry.ClientOptions{
			Dsn:   dsn,
			Debug: true,
		})
	if err != nil {
		logger.LogFatal(context.Background(), err, "failed to initialize sentry", nil)
	}

	logger.LogInfo(context.Background(), "Sentry initialized", nil)
}
