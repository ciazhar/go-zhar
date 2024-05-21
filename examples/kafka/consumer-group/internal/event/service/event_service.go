package service

import "github.com/ciazhar/go-zhar/pkg/logger"

type EventService struct {
	logger *logger.Logger
}

func (e *EventService) ProcessEvent(key, value string) {

	e.logger.Infof("key: %s", key)
	e.logger.Infof("value: %s", value)
}

func NewEventService(logger *logger.Logger) *EventService {
	return &EventService{logger: logger}

}
