package server

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"

	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	name      string
	app       *fiber.App
	addr      string
	serverErr error
	log       zerolog.Logger
}

func NewFiberServer(ctx context.Context, name, addr string, setup func(app *fiber.App)) *FiberServer {
	app := fiber.New()
	setup(app)
	return &FiberServer{
		name: name,
		app:  app,
		addr: addr,
		log:  logger.FromContext(ctx),
	}
}

func (f *FiberServer) Start() error {
	f.log.Info().Msgf("[%s] starting at %s", f.name, f.addr)
	err := f.app.Listen(f.addr)
	if err != nil {
		return fmt.Errorf("[%s] server error: %w", f.name, err)
	}
	return nil
}

func (f *FiberServer) Shutdown(ctx context.Context) error {
	f.log.Info().Msgf("[%s] shutting down...", f.name)
	return f.app.Shutdown()
}

func (f *FiberServer) Name() string {
	return f.name
}
