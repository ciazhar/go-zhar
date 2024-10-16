package main

import (
	"fmt"

	"github.com/ciazhar/go-start-small/pkg/encryption/aes"
)

func main() {
	key := aes.GenerateKey()
	fmt.Println("Key:", key)

	plaintext := "Hello, AES in Go!"

	// Encrypt
	ciphertext := aes.Encrypt(plaintext, key)
	fmt.Println("Encrypted:", ciphertext)

	// Decrypt
	decryptedText := aes.Decrypt(ciphertext, key)
	fmt.Println("Decrypted:", decryptedText)
}
