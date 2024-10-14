package logger

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"strconv"
)

// InitLogger initializes the logger with custom settings
func InitLogger() {

	// Custom function to format the caller to show only the file name and line number
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	// Enable caller logging globally
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.With().Caller().Logger()

}

// logEvent is a helper function to create a log event with common fields
func logEvent(ctx context.Context, event *zerolog.Event, fields map[string]interface{}) *zerolog.Event {
	if requestID, ok := ctx.Value("requestID").(string); ok {
		event = event.Str("requestID", requestID)
	}

	for key, value := range fields {
		event = event.Interface(key, value)
	}

	return event.CallerSkipFrame(1)
}

// LogAndReturnError logs an error and returns it
func LogAndReturnError(ctx context.Context, err error, msg string, fields map[string]interface{}) error {
	logEvent(ctx, log.Error().Err(err), fields).Msg(msg)
	return fmt.Errorf("%s: %w", msg, err)
}

// LogAndReturnWarning logs a warning and returns it as an error
func LogAndReturnWarning(ctx context.Context, err error, msg string, fields map[string]interface{}) error {
	logEvent(ctx, log.Warn().Err(err), fields).Msg(msg)
	return fmt.Errorf("%s: %w", msg, err)
}

// LogDebug logs a debug message
func LogDebug(ctx context.Context, msg string, fields map[string]interface{}) {
	logEvent(ctx, log.Debug(), fields).Msg(msg)
}

// LogInfo logs an info message
func LogInfo(ctx context.Context, msg string, fields map[string]interface{}) {
	logEvent(ctx, log.Info(), fields).Msg(msg)
}
