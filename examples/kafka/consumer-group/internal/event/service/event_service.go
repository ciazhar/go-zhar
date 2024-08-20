package service

import (
	"encoding/json"
	"github.com/ciazhar/go-zhar/examples/kafka/consumer-group/internal/event/model"
	"github.com/ciazhar/go-zhar/pkg/logger"
)

type EventService struct {
	logger *logger.Logger
}

func (e *EventService) ProcessEvent(key, value string) {

	event := model.EmailEvent{}
	err := json.Unmarshal([]byte(value), &event)
	if err != nil {
		return
	}

	e.logger.Infof("key: %s", key)
	e.logger.Infof("value: %v", event)
}

func NewEventService(logger *logger.Logger) *EventService {
	return &EventService{logger: logger}

}
