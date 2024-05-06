package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// GenerateKey Function to generate a random 16-byte key as a base64-encoded string
func GenerateKey() string {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(key)
}

// Encrypt Function to perform AES encryption
func Encrypt(plainText string, key string) string {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		panic(err)
	}

	// Convert plaintext string to byte slice
	plainTextBytes := []byte(plainText)

	// Pad the plaintext
	plainTextBytes = PKCS7Pad(plainTextBytes)

	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainTextBytes)

	// Encode ciphertext to base64 string
	return base64.StdEncoding.EncodeToString(cipherText)
}

// Decrypt Function to perform AES decryption
func Decrypt(cipherText string, key string) string {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		panic(err)
	}

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		panic(err)
	}

	if len(cipherTextBytes) < aes.BlockSize {
		panic("cipherText too short")
	}

	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherTextBytes, cipherTextBytes)

	// Remove PKCS7 padding
	plainTextBytes := PKCS7Unpad(cipherTextBytes)

	// Convert byte slice to string
	return string(plainTextBytes)
}

// PKCS7Pad adds padding to the plaintext
func PKCS7Pad(plainText []byte) []byte {
	padding := aes.BlockSize - (len(plainText) % aes.BlockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plainText, padText...)
}

// PKCS7Unpad removes padding from the plaintext
func PKCS7Unpad(plainText []byte) []byte {
	padding := int(plainText[len(plainText)-1])
	return plainText[:len(plainText)-padding]
}
