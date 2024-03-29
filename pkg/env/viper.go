package env

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

const ConfigPath = "./configs"

const (
	ErrFailedToParseConfig = "Failed to parse config name"
	ErrFailedToReadConfig  = "Fatal error configs file: %s"
)

// Init initializes the configuration with the given name.
//
// It takes basic name string as basic parameter and does not return anything.
func Init(name string) {
	log.Printf("Initializing configuration with name: %s", name)

	splitName := strings.Split(name, ".")
	if len(splitName) != 2 {
		log.Fatalln(ErrFailedToParseConfig)
	}

	configName := splitName[0]
	configType := splitName[1]

	viper.SetConfigName(configName)
	viper.SetConfigType(configType) // REQUIRED if the configs file does not have the extension in the name
	viper.AddConfigPath(ConfigPath) // optionally look for configs in the working directory

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf(ErrFailedToReadConfig, err)
	}
}
