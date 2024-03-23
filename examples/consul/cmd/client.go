package main

import (
	"fmt"
	"github.com/ciazhar/go-zhar/pkg/consul"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"time"
)

func main() {

	// Environment configuration
	env.Init("client.json")
	c := consul.Init(
		viper.GetString("consul.host"),
		viper.GetInt("consul.port"),
		viper.GetString("consul.scheme"),
	)

	myServiceUrl, err := c.RetrieveServiceUrl("my-service")
	if err != nil {
		return
	}

	fmt.Println("Starting Client.")
	var client = &http.Client{
		Timeout: time.Second * 30,
	}
	callServerEvery(10*time.Second, client, myServiceUrl)
}

func hello(t time.Time, client *http.Client, url string) {
	response, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := io.ReadAll(response.Body)
	fmt.Printf("%s. Time is %v\n", body, t)
}

func callServerEvery(d time.Duration, client *http.Client, url string) {
	for x := range time.Tick(d) {
		hello(x, client, url)
	}
}
