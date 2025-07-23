package bootstrap

import (
	"context"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// GracefulShutdown waits for termination syscalls and doing clean up operations after received it.
func GracefulShutdown(ctx context.Context, timeout time.Duration, clients []Service, serverAndWorkers []Service) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-s
	logger.LogInfo("🚦 Signal received. Shutting down...")

	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var wg sync.WaitGroup
	for _, client := range clients {
		wg.Add(1)
		go func(client Service) {
			defer wg.Done()
			logger.LogInfof("🧹 Shutting down: %s", client.Name())
			if err := client.Shutdown(ctxTimeout); err != nil {
				logger.LogErrorf(err, "⚠️ Error shutting down %s:\n", client.Name())
			}
		}(client)
	}
	for _, svc := range serverAndWorkers {
		wg.Add(1)
		go func(svc Service) {
			defer wg.Done()
			logger.LogInfof("🧹 Shutting down: %s", svc.Name())
			if err := svc.Shutdown(ctxTimeout); err != nil {
				logger.LogErrorf(err, "⚠️ Error shutting down %s:\n", svc.Name())
			}
		}(svc)
	}

	wg.Wait()
	logger.LogInfo("✅ All services shut down gracefully.")
}
