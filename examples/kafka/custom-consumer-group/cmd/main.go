package main

import (
	"github.com/ciazhar/go-zhar/examples/kafka/custom-consumer-group/internal/event"
	"github.com/ciazhar/go-zhar/pkg/context_util"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/kafka"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/spf13/viper"
	"sync"
)

func main() {

	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})
	ctx := context_util.SetupSignalHandler()
	var wg sync.WaitGroup

	env.Init("config.json", log)

	producer := kafka.NewSyncProducer(viper.GetString("kafka.brokers"), log)

	event.Init(ctx, &wg, log, producer)
}
