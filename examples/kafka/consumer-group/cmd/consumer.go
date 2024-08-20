package main

import (
	"github.com/ciazhar/go-zhar/examples/kafka/consumer-group/internal/event"
	"github.com/ciazhar/go-zhar/pkg/context_util"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"sync"
)

func main() {

	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})
	ctx := context_util.SetupSignalHandler()
	var wg sync.WaitGroup

	env.Init("config.json", log)

	event.Init(log, ctx, &wg)

	// Wait for all consumer goroutines to complete before exiting
	wg.Wait()
}
