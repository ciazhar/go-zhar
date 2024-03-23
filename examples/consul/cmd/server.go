package main

import (
	"fmt"
	"github.com/ciazhar/go-zhar/pkg/consul"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Environment configuration
	env.Init("server.json")
	c := consul.Init(
		viper.GetString("consul.host"),
		viper.GetInt("consul.port"),
		viper.GetString("consul.scheme"),
	)
	c.RetrieveConfiguration(viper.GetString("consul.key"), viper.GetString("consul.configType"))
	c.RegisterService(
		viper.GetString("application.name"),
		viper.GetString("application.name"),
		viper.GetString("application.host"),
		viper.GetInt("application.port"),
	)

	// Handle termination signals to deregister the service
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start the HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from my-service!")
	})

	go func() {
		fmt.Printf("Starting server on :%s...\n", viper.GetString("application.port"))
		if err := http.ListenAndServe(":"+viper.GetString("application.port"), nil); err != nil {
			log.Fatal(err)
		}
	}()

	<-sigCh
	fmt.Println("Deregistering service...")
	c.DeregisterService(viper.GetString("application.name"))
	os.Exit(1)
}
