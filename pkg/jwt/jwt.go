package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateToken(data map[string]interface{}, keyString string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	for key, value := range data {
		claims[key] = value
	}

	// Set expiration
	expirationTime := time.Now().Add(2 * time.Hour)
	claims["exp"] = expirationTime.Unix()

	// Sign the token with the key
	tokenString, err := token.SignedString([]byte(keyString))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, keyString string) (map[string]interface{}, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(keyString), nil
	})
	if err != nil {
		return nil, err
	}

	// Verify token
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	return claims, nil
}
