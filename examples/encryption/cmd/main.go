package main

import (
	"context"

	"github.com/ciazhar/go-start-small/pkg/encryption/aes"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

func main() {
	key := aes.GenerateKey()
	logger.LogInfo(context.Background(), "Key generated", map[string]interface{}{"key": key})

	plaintext := "Hello, AES in Go!"

	// Encrypt
	ciphertext := aes.Encrypt(plaintext, key)
	logger.LogInfo(context.Background(), "Encrypted text", map[string]interface{}{"ciphertext": ciphertext})

	// Decrypt
	decryptedText := aes.Decrypt(ciphertext, key)
	logger.LogInfo(context.Background(), "Decrypted text", map[string]interface{}{"decryptedText": decryptedText})
}
