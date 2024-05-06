package main

import (
	"github.com/ciazhar/go-zhar/pkg/consul"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"time"
)

func main() {

	// Logger
	log := logger.Init()

	// Environment configuration
	env.Init("client.json", log)
	c := consul.Init(
		viper.GetString("consul.host"),
		viper.GetInt("consul.port"),
		viper.GetString("consul.scheme"),
		log,
	)

	myServiceUrl, err := c.RetrieveServiceUrl("my-service")
	if err != nil {
		return
	}

	log.Infof("Starting Client.")
	var client = &http.Client{
		Timeout: time.Second * 30,
	}
	callServerEvery(10*time.Second, client, myServiceUrl, log)
}

func hello(t time.Time, client *http.Client, url string, log logger.Logger) {
	response, err := client.Get(url)
	if err != nil {
		log.Infof("Error: %v", err)
		return
	}
	body, _ := io.ReadAll(response.Body)
	log.Infof("%s. Time is %v", body, t)
}

func callServerEvery(d time.Duration, client *http.Client, url string, log logger.Logger) {
	for x := range time.Tick(d) {
		hello(x, client, url, log)
	}
}
