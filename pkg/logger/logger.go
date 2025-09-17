package logger

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"io"
	"os"
	"path/filepath"
	"strconv"

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

// FromContext retrieves the logger from the context
func FromContext(ctx context.Context) zerolog.Logger {
	span := trace.SpanFromContext(ctx)
	spanCtx := span.SpanContext()
	if spanCtx.HasTraceID() {
		return log.Logger.With().
			Str("trace_id", spanCtx.TraceID().String()).
			Str("span_id", spanCtx.SpanID().String()).
			Logger()
	}
	return log.Logger
}

// LogFatal is a wrapper for zerolog.LogFatal
func LogFatal(err error) *zerolog.Event {
	return log.Fatal().Err(err)
}

func GetLogger() zerolog.Logger {
	return log.Logger
}
