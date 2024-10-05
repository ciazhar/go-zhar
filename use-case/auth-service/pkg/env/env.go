package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func InitEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	requiredEnvVars := []string{"POSTGRES_URL", "REDIS_HOST", "JWT_SECRET", "PORT"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Required environment variable %s not set", envVar)
		}
	}
}
