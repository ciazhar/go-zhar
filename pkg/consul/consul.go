package consul

import (
	"bytes"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"log"
)

type Consul struct {
	client *consul.Client
}

func Init(host string, port int, schema string) *Consul {
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
	return &Consul{
		client: client,
	}
}

func (c Consul) RetrieveConfiguration(key string, configType string) {
	pair, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		log.Fatalf("Error retrieving key %s: %s", key, err)
	}

	if pair != nil {
		envData := pair.Value
		viper.SetConfigType(configType)
		if err := viper.ReadConfig(bytes.NewBuffer(envData)); err != nil {
			log.Fatalf("Error reading configuration for key %s: %s", key, err)
		}

		log.Printf("Configuration registered for key: %s", key)
	}
}

func (c Consul) RegisterService(id, name, host string, port int) {

	reg := &consul.AgentServiceRegistration{
		ID:   id,
		Name: name,
		Port: port,
		Check: &consul.AgentServiceCheck{
			TCP:      fmt.Sprintf("%s:%d", host, port),
			Interval: "10s",
			Timeout:  "2s",
		},
	}

	log.Printf("Registering service with ID: %s, Name: %s", id, name)

	err := c.client.Agent().ServiceRegister(reg)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
}

func (c Consul) DeregisterService(id string) {
	err := c.client.Agent().ServiceDeregister(id)
	if err != nil {
		log.Printf("Error deregistering service: %v", err)
	}
}

func (c Consul) RetrieveServiceUrl(id string) (string, error) {
	log.Printf("Retrieving service URL for id: %s", id)

	service, _, err := c.client.Agent().Service(id, nil)
	if err != nil {
		log.Printf("Error retrieving service: %v", err)
		return "", err
	}

	url := fmt.Sprintf("http://%s:%v", service.Address, service.Port)
	log.Printf("Service URL: %s", url)

	return url, nil
}
