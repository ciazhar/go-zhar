package env

import (
	"embed"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/spf13/viper"
	"strings"
)

const ConfigPath = "./configs"

const (
	ErrFailedToParseConfig = "Failed to parse config name"
	ErrFailedToReadConfig  = "Fatal error configs file: %s"
)

// Init initializes the configuration with the given name using the provided logger.
//
// Parameters:
// - name: the name of the configuration file.
// - logger: the logger used for logging.
//
// Returns:
// None.
func Init(name string, logger logger.Logger) {
	logger.Infof("Initializing configuration with name: %s", name)

	splitName := strings.Split(name, ".")
	if len(splitName) != 2 {
		logger.Fatal(ErrFailedToParseConfig)
	}

	configName := splitName[0]
	configType := splitName[1]

	viper.SetConfigName(configName)
	viper.SetConfigType(configType) // REQUIRED if the configs file does not have the extension in the name
	viper.AddConfigPath(ConfigPath) // optionally look for configs in the working directory

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatalf(ErrFailedToReadConfig, err)
	}
}

// InitFS initializes the configuration with the given name using the provided embed.FS and logger.
//
// Parameters:
// - name: the name of the configuration file.
// - embedFs: the embed.FS used to open the configuration file.
// - logger: the logger used for logging.
//
// Returns:
// None.
func InitFS(name string, embedFs embed.FS, logger logger.Logger) {
	logger.Infof("Initializing configuration with name: %s", name)

	splitName := strings.Split(name, ".")
	if len(splitName) != 2 {
		logger.Fatal(ErrFailedToParseConfig)
	}

	configType := splitName[1]
	viper.SetConfigType(configType) // REQUIRED if the configs file does not have the extension in the name

	open, err := embedFs.Open(name)
	if err != nil {
		logger.Fatalf(ErrFailedToReadConfig, err)
	}

	err = viper.ReadConfig(open)
	if err != nil {
		logger.Fatalf(ErrFailedToReadConfig, err)
	}
}
