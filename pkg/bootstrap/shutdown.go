package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

// GracefulShutdown waits for termination syscalls and doing clean up operations after received it.
func GracefulShutdown(ctx context.Context, timeout time.Duration, clients []Service, serverAndWorkers []Service) {

	var log = logger.FromContext(ctx).With().Logger()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-s
	log.Info().Msg("🚦 Signal received. Shutting down...")

	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var wg sync.WaitGroup
	for _, client := range clients {
		wg.Add(1)
		go func(client Service) {
			defer wg.Done()
			log.Info().Msgf("Shutting down %s...", client.Name())
			if err := client.Shutdown(ctxTimeout); err != nil {
				log.Error().Err(err).Msgf("Error shutting down %s:\n", client.Name())
			}
		}(client)
	}
	for _, svc := range serverAndWorkers {
		wg.Add(1)
		go func(svc Service) {
			defer wg.Done()
			log.Info().Msgf("Shutting down %s...", svc.Name())
			if err := svc.Shutdown(ctxTimeout); err != nil {
				log.Error().Err(err).Msgf("Error shutting down %s:\n", svc.Name())
			}
		}(svc)
	}

	wg.Wait()
	log.Info().Msg("All services shut down gracefully.")
}
