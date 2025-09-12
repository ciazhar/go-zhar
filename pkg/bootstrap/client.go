package bootstrap

import (
	"context"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/rs/zerolog"
)

type ClientService struct {
	name string
	stop func() error
	log  zerolog.Logger
}

func (c *ClientService) Start() error {
	c.log.Info().Msgf("[%s] client initialized", c.name)
	return nil
}

func (c *ClientService) Shutdown(ctx context.Context) error {
	c.log.Info().Msgf("[%s] closing connection...", c.name)
	return c.stop()
}

func (c *ClientService) Name() string {
	return c.name
}

func NewClientService(ctx context.Context, name string, stop func() error) *ClientService {
	return &ClientService{
		name: name,
		stop: stop,
		log:  logger.FromContext(ctx),
	}
}
