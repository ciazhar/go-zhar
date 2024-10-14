package pkg

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

const ConfigPath = "./configs"

func InitEnv(name string) {
	log.Printf("Initializing configuration with name: %s\n", name)

	splitName := strings.Split(name, ".")
	if len(splitName) != 2 {
		log.Fatal("Failed to parse config name")
	}

	configName := splitName[0]
	configType := splitName[1]

	viper.SetConfigName(configName)
	viper.SetConfigType(configType) // REQUIRED if the configs file does not have the extension in the name
	viper.AddConfigPath(ConfigPath) // optionally look for configs in the working directory

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error configs file: %s", err)
	}
}
