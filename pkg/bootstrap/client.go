package bootstrap

import (
	"context"
	"github.com/ciazhar/go-zhar/pkg/logger"
)

type ClientService struct {
	name string
	stop func() error
}

func (c *ClientService) Start() error {
	logger.LogInfof("[%s] client initialized", c.name)
	return nil
}

func (c *ClientService) Shutdown(ctx context.Context) error {
	logger.LogInfof("[%s] closing connection...", c.name)
	return c.stop()
}

func (c *ClientService) Name() string {
	return c.name
}

func NewClientService(name string, stop func() error) *ClientService {
	return &ClientService{
		name: name,
		stop: stop,
	}
}
