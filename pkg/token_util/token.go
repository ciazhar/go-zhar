package token_util

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func ExtractToken(ctx *fiber.Ctx) (string, error) {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no token provided")
	}
	if len(authHeader) < len("Bearer ") {
		return "", errors.New("invalid token format")
	}
	return authHeader[len("Bearer "):], nil
}
