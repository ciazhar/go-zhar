package main

import (
	"log"
	"os"
	"time"
)

func main() {
	// Create log directory if it doesn't exist
	_ = os.MkdirAll("/logdata", 0755)

	f, err := os.OpenFile("/logdata/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags|log.Lshortfile)

	for {
		logger.Println("INFO: Hello from Go App")
		time.Sleep(5 * time.Second)
	}
}
