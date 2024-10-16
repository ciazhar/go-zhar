package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/ciazhar/go-start-small/pkg/config"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/spf13/viper"
)

func main() {

	// Configuration using flags for source, type, and other details
	var logLevel string
	var consoleOutput bool
	var source, configType, fileName, filePath, consulEndpoint, consulPath string

	// Parse command-line flags
	flag.StringVar(&logLevel, "log-level", "debug", "Log level")
	flag.BoolVar(&consoleOutput, "console-output", true, "Console output")
	flag.StringVar(&source, "source", "consul", "Configuration source (file or consul)")
	flag.StringVar(&fileName, "file-name", "config.json", "Name of the configuration file")
	flag.StringVar(&filePath, "file-path", "configs", "Path to the configuration file")
	flag.StringVar(&configType, "config-type", "json", "Configuration file type")
	flag.StringVar(&consulEndpoint, "consul-endpoint", "localhost:8500", "Consul endpoint")
	flag.StringVar(&consulPath, "consul-path", "path/to/config", "Path to the configuration in Consul")
	flag.Parse()

	// Initialize logger with parsed configuration
	logger.InitLogger(logger.LogConfig{
		LogLevel:      logLevel,
		ConsoleOutput: consoleOutput,
	})

	// Configuration using flags for source, type, and other details
	fileConfig := config.Config{
		Source: source,
		Type:   configType,
		File: config.FileConfig{
			FileName: fileName,
			FilePath: filePath,
		},
		Consul: config.ConsulConfig{
			Endpoint: consulEndpoint,
			Path:     consulPath,
		},
	}

	config.InitConfig(fileConfig)

	// Verify that the configuration was loaded correctly
	if viper.GetString("key1") != "value1" {
		logger.LogFatal(context.Background(), nil, fmt.Sprintf("Expected key1 to be 'value1', got '%s'", viper.GetString("key1")), nil)
	}
	if viper.GetString("key2") != "value2" {
		logger.LogFatal(context.Background(), nil, fmt.Sprintf("Expected key2 to be 'value2', got '%s'", viper.GetString("key2")), nil)
	}
}
