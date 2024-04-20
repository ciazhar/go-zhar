package service

import "github.com/ciazhar/go-zhar/pkg/logger"

type EventService interface {
	ProcessEvent(key, value string)
}

type eventService struct {
	logger logger.Logger
}

func (e eventService) ProcessEvent(key, value string) {

	e.logger.Infof("key: %s", key)
	e.logger.Infof("value: %s", value)
}

func NewEventService(logger logger.Logger) EventService {
	return &eventService{logger: logger}

}
