package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const ConfigPath = "./configs"

// Config holds the configuration details for loading settings
type Config struct {
	Source ConfigSource
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

type ConfigSource int

const (
	File ConfigSource = iota
	Consul
	Vault
)

func (r ConfigSource) String() string {
	return [...]string{"file", "consul", "vault"}[r]
}

func ParseConfigSource(s string) (ConfigSource, error) {
	switch strings.ToLower(s) {
	case "file":
		return File, nil
	case "consul":
		return Consul, nil
	case "vault":
		return Vault, nil
	default:
		return ConfigSource(-1), fmt.Errorf("invalid ConfigSource: %s", s)
	}
}

// InitConfig initializes the configuration based on the provided config struct and source
func InitConfig(ctx context.Context, config Config) {

	var (
		log = logger.FromContext(ctx).With().Any("config", config).Logger()
	)

	switch config.Source {
	case File:

		log.Info().Msg("Loading config from file")

		viper.SetConfigName(config.File.FileName)
		viper.SetConfigType(config.Type) // REQUIRED if the config file does not have the extension in the name
		if config.Consul.Path == "" {
			viper.AddConfigPath(ConfigPath) // optionally look for configs in the working directory
		} else {
			viper.AddConfigPath(config.File.FilePath)
		}

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal().Msg("Failed to read config from file")
		}

	case Consul:

		log.Info().Msg("Loading config from Consul")

		// Load configuration from Consul
		viper.SetConfigType(config.Type)
		err := viper.AddRemoteProvider(config.Source.String(), config.Consul.Endpoint, config.Consul.Path)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to add remote provider")
		}
		err = viper.ReadRemoteConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read config from Consul")
		}
	default:
		log.Fatal().Msg("Invalid config source")
	}

	log.Info().Msg("Config loaded successfully")
}
