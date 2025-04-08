package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Ensure logs folder exists
	_ = os.MkdirAll("logdata", os.ModePerm)

	// Log file
	logFile, err := os.OpenFile("logdata/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Multi-writer: stdout + file
	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)

	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(multi).
		With().
		Timestamp().
		Str("app", "go-zerolog-app").
		Str("env", "dev").
		Caller().
		Logger()

	log.Logger = logger

	// Example periodic log
	go func() {
		for {
			log.Info().
				Str("event", "heartbeat").
				Msg("Application is alive")
			time.Sleep(5 * time.Second)
		}
	}()

	// Simple HTTP server for JSON response
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("path", r.URL.Path).
			Str("method", r.Method).
			Msg("Received /ping request")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","message":"pong"}`))
	})

	log.Info().Msg("Starting server on :8080")
	log.Fatal().Err(http.ListenAndServe(":8080", nil)).Msg("Server exited")
}
