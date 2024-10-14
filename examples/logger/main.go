package main

import (
	"errors"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	errs := errors.New("ayeeyey")
	//log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	//log.Info().Msgf("hai anjeng %s, err %v", "satu", errs)

	//data := fmt.Sprintf("hai anjeng %s, err %v", "satu", errs)
	//println(data)

	Infof("hai anjeng %s, err %v", "satu", errs)
}

func Infof(format string, v ...interface{}) {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	//data := fmt.Sprintf(format, v)
	//println(data)
	//
	//log.Info().Msg(data)
	log.Info().Msgf(format, v...)
}
