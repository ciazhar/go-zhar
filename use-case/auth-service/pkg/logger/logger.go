package logger

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"strconv"
)

func InitLogger() {

	// Custom function to format the caller to show only the file name and line number
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	// Enable caller logging globally
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.With().Caller().Logger()

}

func LogAndReturnError(ctx context.Context, err error, msg string, fields map[string]string) error {
	event := log.Error().CallerSkipFrame(1).Err(err)

	// Optionally: you can add request IDs or context-related info if available
	requestID := ctx.Value("requestID")
	if requestID != nil {
		event.Str("requestID", requestID.(string))
	}

	// Add additional context fields to the log
	for key, value := range fields {
		event.Str(key, value)
	}

	event.Msg(msg)

	// Wrap the original error to preserve the stack trace
	return fmt.Errorf("%s: %w", msg, err)
}

func LogAndReturnWarning(ctx context.Context, err error, msg string, fields map[string]string) error {
	event := log.Warn().CallerSkipFrame(1).Err(err)

	// Optionally: you can add request IDs or context-related info if available
	requestID := ctx.Value("requestID")
	if requestID != nil {
		event.Str("requestID", requestID.(string))
	}

	// Add additional context fields to the log
	for key, value := range fields {
		event.Str(key, value)
	}

	event.Msg(msg)

	// Wrap the original error to preserve the stack trace
	return fmt.Errorf("%s: %w", msg, err)
}

func LogDebug(ctx context.Context, msg string, fields map[string]string) {
	event := log.Debug().CallerSkipFrame(1)

	// Optionally: you can add request IDs or context-related info if available
	requestID := ctx.Value("requestID")
	if requestID != nil {
		event.Str("requestID", requestID.(string))
	}

	// Add additional context fields to the log
	for key, value := range fields {
		event.Str(key, value)
	}

	event.Msg(msg)
}

func LogInfo(ctx context.Context, msg string, fields map[string]string) {
	event := log.Info().CallerSkipFrame(1)

	// Optionally: you can add request IDs or context-related info if available
	requestID := ctx.Value("requestID")
	if requestID != nil {
		event.Str("requestID", requestID.(string))
	}

	// Add additional context fields to the log
	for key, value := range fields {
		event.Str(key, value)
	}

	event.Msg(msg)
}
