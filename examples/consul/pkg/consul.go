package pkg

import (
	"bytes"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"log"
)

func InitConsul(host string, port int, schema string) *consul.Client {
	config := consul.Config{
		Address: fmt.Sprintf("%s:%d", host, port),
		Scheme:  schema,
	}
	client, err := consul.NewClient(&config)
	if err != nil {
		log.Fatalf("Error initializing Consul client: %v\n", err)
		return nil
	}

	log.Println("Consul client initialized successfully")
	return client
}

func RetrieveConfiguration(consul *consul.Client, key string, configType string) {
	log.Printf("key: %s, configType: %s\n", key, configType)

	pair, _, err := consul.KV().Get(key, nil)
	if err != nil {
		log.Fatalf("Error retrieving key %s: %s", key, err)
	}

	if pair != nil {
		envData := pair.Value
		viper.SetConfigType(configType)
		if err := viper.ReadConfig(bytes.NewBuffer(envData)); err != nil {
			log.Fatalf("Error reading configuration for key %s: %s", key, err)
		}

		log.Printf("Configuration registered for key: %s\n", key)
	}
}

func RegisterService(c *consul.Client, id, name, host string, port int) {

	reg := consul.AgentServiceRegistration{
		ID:   id,
		Name: name,
		Port: port,
		Check: &consul.AgentServiceCheck{
			TCP:      fmt.Sprintf("%s:%d", host, port),
			Interval: "10s",
			Timeout:  "2s",
		},
	}

	log.Printf("Registering service with ID: %s, Name: %s\n", id, name)

	err := c.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
}

func DeregisterService(consul *consul.Client, id string) {
	err := consul.Agent().ServiceDeregister(id)
	if err != nil {
		log.Printf("Error deregister service: %v\n", err)
	}
}

func RetrieveServiceUrl(consul *consul.Client, id string) (string, error) {
	log.Printf("Retrieving service URL for id: %s\n", id)

	service, _, err := consul.Agent().Service(id, nil)
	if err != nil {
		log.Printf("Error retrieving service: %v\n", err)
		return "", err
	}

	url := fmt.Sprintf("http://%s:%v", service.Address, service.Port)
	log.Printf("Service URL: %s\n", url)

	return url, nil
}
