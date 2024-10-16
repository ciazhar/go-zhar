package aes

import (
	"encoding/base64"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	key := GenerateKey()
	if len(key) == 0 {
		t.Error("Generated key is empty")
	}
	if _, err := base64.StdEncoding.DecodeString(key); err != nil {
		t.Errorf("Generated key is not valid base64: %v", err)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := GenerateKey()
	plainText := "Hello, World!"

	cipherText := Encrypt(plainText, key)
	if cipherText == "" {
		t.Error("Encryption failed; ciphertext is empty")
	}

	decryptedText := Decrypt(cipherText, key)
	if decryptedText != plainText {
		t.Errorf("Decrypted text does not match original. Got: %s, Want: %s", decryptedText, plainText)
	}
}

func TestDecryptInvalidCipherText(t *testing.T) {
	key := GenerateKey()
	invalidCipherText := "this_is_not_a_valid_base64_string"

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Decrypt did not panic on invalid cipher text")
		}
	}()

	Decrypt(invalidCipherText, key)
}

func TestDecryptInvalidKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Decrypt did not panic on invalid key")
		}
	}()

	Decrypt("validCipherText", "this_is_not_a_valid_base64_string")
}
