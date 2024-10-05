package env

import (
	"github.com/joho/godotenv"
	"log"
)

func InitEnv() {

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
