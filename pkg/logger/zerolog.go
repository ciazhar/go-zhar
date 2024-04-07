package logger

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

type Logger struct {
	serviceLogger zerolog.Logger
}

func Init() Logger {

	if _, err := os.Stat("./log"); os.IsNotExist(err) {
		err := os.Mkdir("./log", 0755)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to create log directory")
		}
	}

	file, err := os.OpenFile("./log/logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to open log file")
	}

	multi := zerolog.MultiLevelWriter(file, zerolog.ConsoleWriter{Out: os.Stdout})

	return Logger{
		serviceLogger: zerolog.New(multi).With().Timestamp().Logger(),
	}
}

func (l Logger) GetServiceLogger() zerolog.Logger {
	return l.serviceLogger
}

func (l Logger) Errorf(format string, a ...interface{}) error {
	errs := fmt.Errorf(format, a...)
	sentry.CaptureException(errs)
	l.serviceLogger.Error().Caller().Msgf(errs.Error())
	return errs
}

func (l Logger) Infof(format string, a ...interface{}) {
	l.serviceLogger.Info().Msgf(format, a...)
}

func (l Logger) Info(format string) {
	l.serviceLogger.Info().Msg(format)
}

func (l Logger) Fatalf(format string, a ...interface{}) {
	l.serviceLogger.Fatal().Msgf(format, a...)
}

func (l Logger) Fatal(format string) {
	l.serviceLogger.Fatal().Msg(format)
}
