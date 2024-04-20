package main

import (
	"github.com/ciazhar/go-zhar/examples/kafka/consumer-group/internal/event"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
)

func main() {

	log := logger.Init()

	env.Init("config.json", log)

	event.Init(log)
}
