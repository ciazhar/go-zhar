package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// GenerateJWTToken Generate JWT access token
func GenerateJWTToken(userID string, ttl time.Duration) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := []byte(viper.GetString("jwt.secret"))

	return token.SignedString(jwtSecret)
}

// CustomClaims JWT Claims
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// ValidateJWT Validate JWT token
func ValidateJWT(tokenString string) (*CustomClaims, error) {

	jwtSecret := []byte(viper.GetString("jwt.secret"))

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(*CustomClaims), nil
}
