package main

import (
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"pkg"
	"time"
)

func main() {

	// Environment configuration
	pkg.InitEnv("server.json")

	c := pkg.InitConsul(
		viper.GetString("consul.host"),
		viper.GetInt("consul.port"),
		viper.GetString("consul.scheme"),
	)

	serverServiceUrl, err := pkg.RetrieveServiceUrl(c, "server-service")
	if err != nil {
		return
	}

	log.Println("Starting Client.")
	var client = &http.Client{
		Timeout: time.Second * 30,
	}
	callServerEvery(10*time.Second, client, serverServiceUrl)
}

func hello(t time.Time, client *http.Client, url string) {
	response, err := client.Get(url)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	body, _ := io.ReadAll(response.Body)
	log.Printf("%s. Time is %v\n", body, t)
}

func callServerEvery(d time.Duration, client *http.Client, url string) {
	for x := range time.Tick(d) {
		hello(x, client, url)
	}
}
