package config

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestInitConfig tests the InitConfig function for both file and Consul sources
func TestInitConfig(t *testing.T) {
	// Initialize logger for testing
	logger.InitLogger(logger.LogConfig{
		LogLevel:      "debug",
		ConsoleOutput: true,
	})

	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Create a sample configuration file for testing
	configFileName := "test_config.json"
	configFilePath := tempDir + "/" + configFileName
	configContent := `{
		"key1": "value1",
		"key2": "value2"
	}`
	if err := os.WriteFile(configFilePath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Test file-based configuration
	fileConfig := Config{
		Source: "file",
		File: FileConfig{
			FilePath: configFilePath,
			FileName: configFileName,
		},
		Type: "json",
	}

	// Set the config path to the temporary directory
	viper.AddConfigPath(tempDir) // Add the temporary directory to the config paths

	InitConfig(fileConfig)

	// Verify that the configuration was loaded correctly
	if viper.GetString("key1") != "value1" {
		t.Errorf("Expected key1 to be 'value1', got '%s'", viper.GetString("key1"))
	}
	if viper.GetString("key2") != "value2" {
		t.Errorf("Expected key2 to be 'value2', got '%s'", viper.GetString("key2"))
	}

}

func TestInitConfigWithConsul(t *testing.T) {
	// Start a Consul container
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "consul:1.14.0",
		ExposedPorts: []string{"8500/tcp"},
		WaitingFor:   wait.ForHTTP("/v1/status/leader").WithPort("8500/tcp").WithStatusCodeMatcher(func(status int) bool { return status == 200 }),
	}

	// Attempt to start the Consul container twice in case the first attempt fails
	for attempt := 0; attempt < 2; attempt++ {
		consulContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			if attempt == 1 {
				t.Fatal("Failed to start Consul container:", err)
			}
			continue // Try again if the first attempt fails
		}

		// Get the host and mapped port for Consul
		host, err := consulContainer.Host(ctx)
		if err != nil {
			t.Fatal("Failed to get Consul container host:", err)
		}

		port, err := consulContainer.MappedPort(ctx, "8500")
		if err != nil {
			t.Fatal("Failed to get Consul container port:", err)
		}

		consulEndpoint := fmt.Sprintf("%s:%s", host, port.Port())

		// Simulate adding key-value pairs to Consul using http.Client
		keyPath := "config/test"
		configData := `{
			"key1": "value1",
			"key2": "value2"
		}`

		// Use http.Client to send the PUT request to Consul
		httpClient := &http.Client{Timeout: 10 * time.Second}
		consulURL := fmt.Sprintf("http://%s/v1/kv/%s", consulEndpoint, keyPath)
		reqBody := bytes.NewReader([]byte(configData))

		request, err := http.NewRequest("PUT", consulURL, reqBody)
		if err != nil {
			t.Fatal("Failed to create PUT request:", err)
		}
		request.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(request)
		if err != nil {
			t.Fatal("Failed to write config to Consul:", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to write config to Consul: status code %d", resp.StatusCode)
		}

		// Initialize logger for the test
		logger.InitLogger(logger.LogConfig{
			LogLevel:      "debug",
			ConsoleOutput: true,
		})

		// Create a configuration struct for Consul
		consulViperConfig := Config{
			Source: "consul",
			Consul: ConsulConfig{
				Endpoint: consulEndpoint,
				Path:     keyPath,
			},
			Type: "json",
		}

		// Initialize Viper config from Consul
		InitConfig(consulViperConfig)

		// Retrieve all settings from Viper
		settings := viper.AllSettings()

		// Assertions to verify that the settings were loaded correctly
		assert.Equal(t, "value1", settings["key1"])
		assert.Equal(t, "value2", settings["key2"])

		err = resp.Body.Close()
		if err != nil {
			return
		}

		err = consulContainer.Terminate(ctx)
		if err != nil {
			return
		}

		return // Exit the loop if the test is successful
	}
}
