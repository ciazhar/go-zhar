package main

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
	jwtSecret       = []byte("your-secret")
	rdb             = redis.NewClient(&redis.Options{
		Addr: "localhost:6377",
	})
	ctx = context.Background()
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
}

func generateJWT(userID string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func saveRefreshToken(userID, token string) error {
	return rdb.Set(ctx, "refresh:"+token, userID, refreshTokenTTL).Err()
}

func deleteRefreshToken(token string) error {
	return rdb.Del(ctx, "refresh:"+token).Err()
}

func loginHandler(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if req.Username != "john_doe" || req.Password != "123456" {
		return fiber.ErrUnauthorized
	}

	userID := "u123"
	accessToken, _ := generateJWT(userID, accessTokenTTL)
	refreshToken := uuid.NewString()
	saveRefreshToken(userID, refreshToken)

	return c.JSON(TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
	})
}

func profileHandler(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	return c.JSON(fiber.Map{
		"id":       userID,
		"username": "john_doe",
		"email":    "john@mail.com",
	})
}

func refreshHandler(c *fiber.Ctx) error {
	type Req struct {
		RefreshToken string `json:"refresh_token"`
	}
	var req Req
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	userID, err := rdb.Get(ctx, "refresh:"+req.RefreshToken).Result()
	if errors.Is(err, redis.Nil) {
		return fiber.ErrUnauthorized
	} else if err != nil {
		return fiber.ErrInternalServerError
	}

	// Optional: Token Rotation
	_ = deleteRefreshToken(req.RefreshToken)
	newRefresh := uuid.NewString()
	_ = saveRefreshToken(userID, newRefresh)

	accessToken, _ := generateJWT(userID, accessTokenTTL)
	return c.JSON(TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefresh,
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
	})
}

func logoutHandler(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.ErrUnauthorized
	}

	// Extract token (simplified)
	token := authHeader[len("Bearer "):]
	_ = deleteRefreshToken(token) // token di sini dianggap refresh token

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func jwtMiddleware(c *fiber.Ctx) error {
	tokenStr := c.Get("Authorization")
	if tokenStr == "" {
		return fiber.ErrUnauthorized
	}
	tokenStr = tokenStr[len("Bearer "):]

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return fiber.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Locals("userID", claims["sub"])
	return c.Next()
}

func main() {
	app := fiber.New()

	app.Post("/login", loginHandler)
	app.Post("/refresh-token", refreshHandler)
	app.Post("/logout", logoutHandler)

	app.Get("/user/profile", jwtMiddleware, profileHandler)

	log.Fatal(app.Listen(":3000"))
}
