package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ciazhar/go-start-small/pkg/context_util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig holds the configuration for the logger
type LogConfig struct {
	LogLevel      string
	LogFile       string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Compress      bool
	ConsoleOutput bool
}

// InitLogger initializes the logger with custom settings and log rotation
func InitLogger(config LogConfig) {
	// Set up lumberjack for log rotation
	logRotator := &lumberjack.Logger{
		Filename:   config.LogFile,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	// Determine the writers
	var writers []io.Writer
	writers = append(writers, logRotator)
	if config.ConsoleOutput {
		writers = append(writers, os.Stdout)
	}
	multiWriter := io.MultiWriter(writers...)

	// Custom function to format the caller to show only the file name and line number
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	// Set log level
	if err := SetLogLevel(config.LogLevel); err != nil {
		log.Fatal().Err(err).Msg("Failed to set log level")
	}

	// Configure zerolog
	log.Logger = zerolog.New(multiWriter).With().Timestamp().Caller().Logger().Output(zerolog.SyncWriter(multiWriter))
}

func SetLogLevel(level string) error {
	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(parsedLevel)
	return nil
}

// logEvent is a helper function to create a log event with common fields
func logEvent(ctx context.Context, event *zerolog.Event, fields map[string]interface{}) *zerolog.Event {
	if requestID, ok := ctx.Value(context_util.RequestIDKey).(string); ok {
		event = event.Str("request_id", requestID)
	}

	for key, value := range fields {
		switch v := value.(type) {
		case error:
			event = event.Err(v)
		case map[string]interface{}:
			event = event.Fields(v)
		default:
			event = event.Interface(key, value)
		}
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
