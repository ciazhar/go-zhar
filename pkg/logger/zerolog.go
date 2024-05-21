package logger

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
)

// Config Configuration for logging
type Config struct {

	// Enable console logging
	ConsoleLoggingEnabled bool

	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool

	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool

	// Directory to log to to when filelogging is enabled
	Directory string

	// Filename is the name of the logfile which will be placed inside the directory
	Filename string

	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int

	// MaxBackups the max number of rolled files to keep
	MaxBackups int

	// MaxAge the max age in days to keep a logfile
	MaxAge int
}

type Logger struct {
	serviceLogger *zerolog.Logger
}

// Init sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func Init(config Config) *Logger {
	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	mw := io.MultiWriter(writers...)

	logger := zerolog.New(mw).With().Timestamp().Logger()

	logger.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("jsonLogOutput", config.EncodeLogsAsJson).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Msg("logging configured")

	return &Logger{
		serviceLogger: &logger,
	}
}

func newRollingFile(config Config) io.Writer {

	l := &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}

	return l
}

func (l *Logger) GetServiceLogger() *zerolog.Logger {
	return l.serviceLogger
}

func (l *Logger) Errorf(format string, a ...interface{}) error {
	errs := fmt.Errorf(format, a...)
	sentry.CaptureException(errs)
	l.serviceLogger.Error().Caller().Msgf(errs.Error())
	return errs
}

func (l *Logger) Infof(format string, a ...interface{}) {
	l.serviceLogger.Info().Msgf(format, a...)
}

func (l *Logger) Info(format string) {
	l.serviceLogger.Info().Msg(format)
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.serviceLogger.Fatal().Msgf(format, a...)
}

func (l *Logger) Fatal(format string) {
	l.serviceLogger.Fatal().Msg(format)
}
