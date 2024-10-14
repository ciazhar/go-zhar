package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestLogger(t *testing.T) {
	// Redirect log output to a buffer for testing
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf)

	// Initialize logger
	InitLogger()

	// Create a context with a request ID
	ctx := context.WithValue(context.Background(), "requestID", "test-request-id")

	// Test cases
	tests := []struct {
		name     string
		logFunc  func()
		expected []string
		level    zerolog.Level
	}{
		{
			name: "LogAndReturnError",
			logFunc: func() {
				err := LogAndReturnError(ctx, errors.New("test error"), "Error occurred", map[string]interface{}{"key": "value"})
				if err == nil {
					t.Error("Expected an error, got nil")
				}
			},
			expected: []string{"error", "test error", "Error occurred", "test-request-id", "key", "value"},
		},
		{
			name: "LogAndReturnWarning",
			logFunc: func() {
				err := LogAndReturnWarning(ctx, errors.New("test warning"), "Warning occurred", map[string]interface{}{"key": "value"})
				if err == nil {
					t.Error("Expected an error, got nil")
				}
			},
			expected: []string{"warn", "test warning", "Warning occurred", "test-request-id", "key", "value"},
		},
		{
            name: "LogDebug",
            logFunc: func() {
                LogDebug(ctx, "Debug message", map[string]interface{}{"key": "value"})
            },
            expected: []string{"debug", "Debug message", "test-request-id", "key", "value"},
            level:    zerolog.DebugLevel,
        },
		{
			name: "LogInfo",
			logFunc: func() {
				LogInfo(ctx, "Info message", map[string]interface{}{"key": "value"})
			},
			expected: []string{"info", "Info message", "test-request-id", "key", "value"},
		},
	}

	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            buf.Reset()
            
            // Set the log level for this test
            zerolog.SetGlobalLevel(tt.level)
            
            tt.logFunc()

            output := buf.String()
            if output == "" {
                t.Fatalf("No log output generated. Check if the log level is set correctly.")
            }

            var logEntry map[string]interface{}
            err := json.Unmarshal([]byte(output), &logEntry)
            if err != nil {
                t.Fatalf("Failed to parse log output: %v\nRaw output: %s", err, output)
            }

            for _, exp := range tt.expected {
                if !strings.Contains(output, exp) {
                    t.Errorf("Expected log to contain '%s', but it didn't. Log: %s", exp, output)
                }
            }

            if caller, ok := logEntry["caller"].(string); !ok || caller == "" {
                t.Error("Expected caller information, but it was missing or empty")
            }
        })
    }
}