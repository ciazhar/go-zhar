package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

func GetValueFromJWT(jwtToken string, key string) interface{} {
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	return claims[key]
}
