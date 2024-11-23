package service

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-start-small/examples/kafka_producer_consumer_graceful_shutdown/internal/event/model"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

type EventService struct {
}

func (e *EventService) ProcessEvent(msg *sarama.ConsumerMessage) error {

	event := model.EmailEvent{}
	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return err
	}

	logger.LogInfo(context.Background(), "Processing event", map[string]interface{}{
		"key":   msg.Key,
		"value": event,
	})

	return nil
}

func NewEventService() *EventService {
	return &EventService{}

}
