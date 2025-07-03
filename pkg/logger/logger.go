package logger

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"strconv"
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

// contextKey is used to prevent context key collisions
type contextKey string

const (
	LoggerKey contextKey = "logger"
)

// InitLogger initializes the logger with custom settings and log rotation
func InitLogger(config LogConfig) {
	multiWriter := setupWriters(config)

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	setLogLevel(config.LogLevel)

	log.Logger = zerolog.New(zerolog.SyncWriter(multiWriter)).With().Timestamp().Caller().Logger()
}

func setupWriters(config LogConfig) io.Writer {
	logRotator := &lumberjack.Logger{
		Filename:   config.LogFile,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	var writers []io.Writer
	writers = append(writers, logRotator)
	if config.ConsoleOutput {
		writers = append(writers, os.Stdout)
	}

	return io.MultiWriter(writers...)
}

func setLogLevel(level string) {
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// WithLogger returns a new context with the logger attached
func WithLogger(ctx context.Context, l zerolog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, l)
}

// FromContext retrieves the logger from the context
func FromContext(ctx context.Context) zerolog.Logger {
	if l, ok := ctx.Value(LoggerKey).(zerolog.Logger); ok {
		return l
	}
	return log.Logger
}

// LogFatal is a wrapper for zerolog.LogFatal
func LogFatal(err error) *zerolog.Event {
	return log.Fatal().Err(err)
}

// LogError is a wrapper for zerolog.LogError
func LogError(err error) *zerolog.Event {
	return log.Error().Err(err)
}

// LogWarn is a wrapper for zerolog.LogWarn
func LogWarn(msg error) *zerolog.Event {
	return log.Warn().Err(msg)
}
