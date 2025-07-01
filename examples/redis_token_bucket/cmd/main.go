package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

const (
	defaultBucketCapacity = 10 // Kapasitas maksimum token default
	tokenRefillRate       = 1  // Token diisi setiap detik
	tokenRefillTime       = 1 * time.Second
)

var redisClient *redis.Client

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

func getRateLimit(clientID string, defaultLimit int) int {
	ctx := context.Background()
	limit, err := redisClient.HGet(ctx, "client_rate_limit", clientID).Int()
	if errors.Is(err, redis.Nil) {
		return defaultLimit
	} else if err != nil {
		log.Println("Gagal mengambil rate limit dari Redis:", err)
		return defaultLimit
	}
	return limit
}

func tokenBucketMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()

		// Granular Limit: Ambil key berdasarkan IP atau User ID
		clientID := c.IP() // Bisa diganti dengan User ID jika ada autentikasi

		bucketCapacity := getRateLimit(clientID, defaultBucketCapacity)
		key := "token_bucket:" + clientID
		now := time.Now().Unix()

		// Ambil token dan timestamp terakhir dari Redis
		lastRefillTimeStr, _ := redisClient.HGet(ctx, key, "last_refill").Result()
		currentTokensStr, _ := redisClient.HGet(ctx, key, "tokens").Result()

		var lastRefillTime int64
		var currentTokens int

		if lastRefillTimeStr != "" {
			lastRefillTime, _ = strconv.ParseInt(lastRefillTimeStr, 10, 64)
		} else {
			lastRefillTime = now
		}

		if currentTokensStr != "" {
			currentTokens, _ = strconv.Atoi(currentTokensStr)
		} else {
			currentTokens = bucketCapacity
		}

		// Hitung token yang perlu diisi ulang
		elapsedTime := now - lastRefillTime
		newTokens := int(elapsedTime) * tokenRefillRate
		currentTokens = min(bucketCapacity, currentTokens+newTokens)

		// Simpan token terbaru dan timestamp refill di Redis
		redisClient.HSet(ctx, key, "tokens", currentTokens)
		redisClient.HSet(ctx, key, "last_refill", now)

		// Periksa apakah masih ada token yang tersedia
		if currentTokens > 0 {
			redisClient.HIncrBy(ctx, key, "tokens", -1)
			return c.Next()
		}

		// Jika token habis, kirim respons 429 Too Many Requests
		retryAfter := tokenRefillTime.Seconds()
		c.Set("Retry-After", fmt.Sprintf("%d", int(retryAfter)))
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error":       "Rate limit exceeded. Please try again later.",
			"retry_after": retryAfter,
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	initRedis()
	r := fiber.New()
	r.Use(tokenBucketMiddleware())

	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome! You are within the rate limit.")
	})

	log.Fatal(r.Listen(":3000"))
}
