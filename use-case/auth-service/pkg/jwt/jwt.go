package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var (
	jwtSecret       = []byte(os.Getenv("JWT_SECRET")) // Change this to a strong secret
	AccessTokenTTL  = 15 * time.Minute                // Access token TTL
	RefreshTokenTTL = 24 * time.Hour
)

// GenerateJWTToken Generate JWT access token
func GenerateJWTToken(userID string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// CustomClaims JWT Claims
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// ValidateJWT Validate JWT token
func ValidateJWT(tokenString string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(*CustomClaims), nil
}
