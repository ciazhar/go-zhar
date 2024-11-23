package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ciazhar/go-start-small/examples/socket_io/internal/model"
	"github.com/ciazhar/go-start-small/examples/socket_io/internal/service"
	"github.com/ciazhar/go-start-small/pkg/config"
	"github.com/ciazhar/go-start-small/pkg/socketio"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
)

func main() {
	// Load configuration
	config.InitConfig(
		config.Config{
			Source: "file",
			Type:   "json",
			File: config.FileConfig{
				FileName: "config.json",
				FilePath: "./configs",
			},
		},
	)

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Create service
	bloodService := service.NewBloodService(viper.GetInt("queue_size"))

	// Routes
	r.Post("/blood", func(w http.ResponseWriter, r *http.Request) {

		var blood model.BloodAvailability
		err := json.NewDecoder(r.Body).Decode(&blood)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Set timestamps
		now := time.Now().Format(time.RFC3339)
		blood.ApplicationID = "digisar"
		blood.CreatedAt = now
		blood.UpdatedAt = now

		if err := bloodService.BroadcastLatestBlood(r.Context(), &blood); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/socket.io/", socketio.Init(model.RoomName, bloodService.ListenLatestBlood).ServeHTTP)

	// Create server
	srv := &http.Server{
		Addr:    ":" + viper.GetString("port"),
		Handler: r,
	}

	// Start server
	go func() {
		log.Printf("Starting server on port %s", viper.GetString("port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
