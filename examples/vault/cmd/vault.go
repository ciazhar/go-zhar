package main

import (
	"context"
	"github.com/hashicorp/vault-client-go"
	"log"
	"time"

	"github.com/hashicorp/vault-client-go/schema"
)

func InitVault(address string) *vault.Client {
	client, err := vault.New(
		vault.WithAddress(address),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func main() {
	ctx := context.Background()

	client := InitVault("http://127.0.0.1:8200")

	// authenticate with a root token (insecure)
	if err := client.SetToken("my-token"); err != nil {
		log.Fatal(err)
	}

	// write a secret
	_, err := client.Secrets.KvV2Write(ctx, "foo", schema.KvV2WriteRequest{
		Data: map[string]any{
			"password1": "abc123",
			"password2": "correct horse battery staple",
		}},
		vault.WithMountPath("secret"),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("secret written successfully")

	// read the secret
	s, err := client.Secrets.KvV2Read(ctx, "foo", vault.WithMountPath("secret"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("secret retrieved:", s.Data.Data)
}
