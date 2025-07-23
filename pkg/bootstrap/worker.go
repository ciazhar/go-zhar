package bootstrap

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	name string
	stop chan struct{}
	wg   sync.WaitGroup
}

func NewWorker(name string) *Worker {
	return &Worker{
		name: name,
		stop: make(chan struct{}),
	}
}

func (w *Worker) Start() error {
	fmt.Printf("[%s] worker started\n", w.name)
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.stop:
				fmt.Printf("[%s] worker stopping...\n", w.name)
				return
			default:
				fmt.Printf("[%s] doing background job...\n", w.name)
				time.Sleep(2 * time.Second)
			}
		}
	}()
	return nil
}

func (w *Worker) Shutdown(ctx context.Context) error {
	close(w.stop)
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		fmt.Printf("[%s] shutdown complete\n", w.name)
	case <-ctx.Done():
		fmt.Printf("[%s] shutdown timeout\n", w.name)
	}
	return nil
}

func (w *Worker) Name() string {
	return w.name
}
