package workerpool

import (
	"bytes"
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// Result represents the outcome of an async function execution
type Result struct {
	Index int
	Err   error
	Time  time.Duration
}

// Options configures the behavior of RunAsync
type Options struct {
	// Timeout for the entire operation. Zero means no timeout
	Timeout time.Duration
	// MaxConcurrent limits the number of concurrent goroutines. Zero means no limit
	MaxConcurrent int
	// ContinueOnError determines whether to continue executing remaining functions after an error
	ContinueOnError bool
}

// DefaultOptions provides sensible default settings
var DefaultOptions = Options{
	Timeout:        30 * time.Second,
	MaxConcurrent:  0, // no limit
	ContinueOnError: true,
}

// GenericFunction represents a function that takes no args and returns an error
type GenericFunction func(ctx context.Context) error

// RunAsync executes functions concurrently and returns their results
func RunAsync(ctx context.Context, opts Options, functions ...GenericFunction) []Result {
	if len(functions) == 0 {
		return nil
	}

	// Create context with timeout if specified
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// Initialize result slice
	results := make([]Result, len(functions))

	// Create worker pool if MaxConcurrent is specified
	var semaphore chan struct{}
	if opts.MaxConcurrent > 0 {
		semaphore = make(chan struct{}, opts.MaxConcurrent)
	}

	// Create wait group for synchronization
	var wg sync.WaitGroup
	wg.Add(len(functions))

	// Execute functions concurrently
	for i, fn := range functions {
		i, fn := i, fn // Create new variables for goroutine closure

		go func() {
			defer wg.Done()

			// Acquire semaphore if using worker pool
			if semaphore != nil {
				select {
				case semaphore <- struct{}{}:
					defer func() { <-semaphore }()
				case <-ctx.Done():
					results[i] = Result{Index: i, Err: ctx.Err()}
					return
				}
			}

			start := time.Now()
			results[i] = executeFunction(ctx, i, fn)
			results[i].Time = time.Since(start)

			// Check if we should continue after error
			if !opts.ContinueOnError && results[i].Err != nil {
				ctx.Done()
			}
		}()
	}

	// Wait for completion or context cancellation
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return results
	case <-ctx.Done():
		return results
	}
}

// executeFunction runs a single function with panic recovery
func executeFunction(ctx context.Context, index int, fn GenericFunction) Result {
	result := Result{Index: index}

	// Handle panics
	defer func() {
		if r := recover(); r != nil {
			result.Err = fmt.Errorf("panic in async function %d: %v\n%s", index, r, formatStack(4))
		}
	}()

	// Execute function with context
	if err := fn(ctx); err != nil {
		result.Err = fmt.Errorf("error in async function %d: %w", index, err)
	}

	return result
}

// formatStack returns a formatted stack trace, skipping the specified number of frames
func formatStack(skip int) string {
	stack := debug.Stack()
	lines := bytes.Split(bytes.TrimSpace(stack), []byte("\n"))
	
	// Skip the specified number of frame pairs (each frame has 2 lines)
	if skip*2 < len(lines) {
		lines = lines[1+skip*2:]
	}
	
	return string(bytes.Join(lines, []byte("\n")))
}