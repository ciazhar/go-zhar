package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ciazhar/go-start-small/pkg/context_util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestLogger(t *testing.T) {
	// Redirect log output to a buffer for testing
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).With().Timestamp().Caller().Logger()

	// Initialize logger with test configuration
	testConfig := LogConfig{
		ConsoleOutput: true,
	}
	InitLogger(testConfig)

	// Override the global logger to use our buffer
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(&buf)

	// Create a context with a request ID
	ctx := context.WithValue(context.Background(), context_util.RequestIDKey, "test-request-id")

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
			t.Log("output",output)

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


func TestLogRotation(t *testing.T) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Use a subdirectory in the current working directory for log files
	logDir := filepath.Join(currentDir, "test_logs")
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}
	defer os.RemoveAll(logDir) // Clean up after the test

	// Clear the directory before starting the test
	files, err := os.ReadDir(logDir)
	if err != nil {
		t.Fatalf("Failed to read log directory: %v", err)
	}
	for _, file := range files {
		os.Remove(filepath.Join(logDir, file.Name()))
	}

	logFile := filepath.Join(logDir, "test.log")

	config := LogConfig{
		LogLevel:      "debug",
		LogFile:       logFile,
		MaxSize:       1, // 1 MB
		MaxBackups:    3,
		MaxAge:        1,
		Compress:      true,
		ConsoleOutput: false,
	}

	InitLogger(config)

	ctx := context.WithValue(context.Background(), context_util.RequestIDKey, "test-request-id")

	// Write logs until rotation occurs
	for i := 0; i < 100000; i++ {
		LogInfo(ctx, fmt.Sprintf("Log message %d", i), nil)
	}

	// Check if log files were created
	files, err = os.ReadDir(logDir)
	if err != nil {
		t.Fatalf("Failed to read log directory: %v", err)
	}

	logFiles := 0
	gzFiles := 0
	for _, file := range files {
		if file.Name() == "test.log" {
			logFiles++
		} else if strings.HasPrefix(file.Name(), "test-") && strings.HasSuffix(file.Name(), ".log.gz") {
			gzFiles++
		}
	}

	if logFiles != 1 {
		t.Errorf("Expected 1 current log file, but found %d", logFiles)
	}

	if gzFiles != 3 {
		t.Errorf("Expected 3 compressed backup files, but found %d", gzFiles)
	}

	// Check if at least one compressed backup exists
	compressedBackupExists := false
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "test-") && strings.HasSuffix(file.Name(), ".log.gz") {
			compressedBackupExists = true
			break
		}
	}

	if !compressedBackupExists {
		t.Errorf("Expected at least one compressed backup file to exist")
	}
}
