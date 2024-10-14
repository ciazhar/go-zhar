package logger

import (
	"github.com/getsentry/sentry-go"
)

func InitSentry(dsn string, logger Logger) {

	err := sentry.Init(
		sentry.ClientOptions{
			Dsn:   dsn,
			Debug: true,
		})
	if err != nil {
		logger.Fatalf("Failed to initialize Sentry: %v", err)
	}

	logger.Info("Sentry initialized successfully")
}
