package config

import (
	"context"

	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const ConfigPath = "./configs"

// Config holds the configuration details for loading settings
type Config struct {
	Source string
	Type   string       // Type of the configuration file (e.g., json, yaml)
	File   FileConfig   // File configuration details
	Consul ConsulConfig // Consul configuration details
}

// FileConfig holds the configuration details for a file
type FileConfig struct {
	FileName string
	FilePath string
}

// ConsulConfig holds the configuration details for Consul
type ConsulConfig struct {
	Endpoint string // Consul endpoint
	Path     string // Path to the configuration in Consul
}

// InitConfig initializes the configuration based on the provided config struct and source
func InitConfig(config Config) {
	if config.Source == "file" {

		logger.LogInfo(context.Background(), "Loading config from file", map[string]interface{}{
			"file_name": config.File.FileName,
			"file_path": config.File.FilePath,
		})

		viper.SetConfigName(config.File.FileName)
		viper.SetConfigType(config.Type) // REQUIRED if the config file does not have the extension in the name
		if config.Consul.Path == "" {
			viper.AddConfigPath(ConfigPath) // optionally look for configs in the working directory
		} else {
			viper.AddConfigPath(config.File.FilePath)
		}

		err := viper.ReadInConfig()
		if err != nil {
			logger.LogFatal(context.Background(), err, "Failed to read config from file", map[string]interface{}{
				"file_name": config.File.FileName,
				"file_path": config.File.FilePath,
			})
		}
	} else if config.Source == "consul" {

		logger.LogInfo(context.Background(), "Loading config from Consul", map[string]interface{}{
			"endpoint": config.Consul.Endpoint,
			"path":     config.Consul.Path,
		})

		// Load configuration from Consul
		viper.SetConfigType(config.Type)
		viper.AddRemoteProvider(config.Source, config.Consul.Endpoint, config.Consul.Path)
		err := viper.ReadRemoteConfig()
		if err != nil {
			logger.LogFatal(context.Background(), err, "Failed to read config from Consul", map[string]interface{}{
				"endpoint": config.Consul.Endpoint,
				"path":     config.Consul.Path,
			})
		}
		
	} else {
		logger.LogFatal(context.Background(), nil, "Invalid config source", map[string]interface{}{
			"source": config.Source,
		})
	}

	logger.LogInfo(context.Background(), "Config loaded successfully", nil)
}
