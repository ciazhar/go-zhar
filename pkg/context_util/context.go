package context_util

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func SetupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sigCh:
			fmt.Printf("Received signal: %v. Initiating shutdown...\n", sig)
			cancel() // Cancel the context to stop consumer goroutines
		}
	}()

	return ctx
}
