package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	leakRate        = time.Second * 1 // Kecepatan mengeluarkan request dari bucket
	rateLimitExpire = time.Minute * 5 // TTL untuk key Redis

	// Redis keys
	leakyBucketRateLimit     = "leaky_bucket_rate_limit"
	clientRateLimitThreshold = "client_rate_limit_threshold"
)

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{Addr: "localhost:6377"})
)

// Middleware Rate Limiting (Leaky Bucket)
func rateLimiter(defaultRateLimit int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Granular Limit: Ambil key berdasarkan IP atau User ID
		clientID := c.IP() // Bisa diganti dengan User ID jika ada autentikasi
		bucketKey := fmt.Sprintf("%s:%s", leakyBucketRateLimit, clientID)
		now := time.Now().Unix()

		// Hitung jumlah request dalam bucket
		reqCount, _ := rdb.LLen(ctx, bucketKey).Result()

		// Jika bucket penuh, tolak request dengan HTTP 429
		limit := getRateLimit(clientID, defaultRateLimit)
		c.Set("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(limit-int(reqCount)))

		// Deteksi anomali
		anomaly := detectAnomaly(clientID)
		if anomaly {
			setRateLimit(clientID, limit/2)
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "The system detected an unusual traffic spike. Rate limits are temporarily lowered.",
			})
		}

		// Jika request melebihi batas, tolak request
		if reqCount >= int64(limit) {
			retryAfter := leakRate.Seconds()
			log.Warn().Str("client", clientID).Msg("Rate limit exceeded")

			c.Set("Retry-After", strconv.Itoa(int(retryAfter)))
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Too many requests",
				"retry_after": retryAfter,
			})
		}

		// Tambahkan timestamp request ke bucket
		_, err := rdb.RPush(ctx, bucketKey, now).Result()
		if err != nil {
			log.Error().Err(err).Msg("Redis error")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}

		// Set expiration untuk bucket (agar tidak ada data expired di Redis)
		rdb.Expire(ctx, bucketKey, rateLimitExpire)

		// Lanjutkan ke handler berikutnya
		return c.Next()
	}
}

func detectAnomaly(clientID string) bool {
	now := time.Now().Unix()
	threshold := 100 // Batas anomali dalam 1 menit

	// Simpan timestamp request
	rdb.ZAdd(ctx, "traffic_monitor:"+clientID, redis.Z{
		Score:  float64(now),
		Member: now,
	})

	// Hapus data yang lebih lama dari 1 menit
	rdb.ZRemRangeByScore(ctx, "traffic_monitor:"+clientID, "0", strconv.FormatInt(now-60, 10))

	// Hitung total request dalam 1 menit
	count, _ := rdb.ZCard(ctx, "traffic_monitor:"+clientID).Result()

	return int(count) > threshold
}

func getRateLimit(clientID string, defaultLimit int) int {

	// Cek di Redis apakah ada limit yang disesuaikan
	limit, err := rdb.HGet(ctx, clientRateLimitThreshold, clientID).Int()
	if errors.Is(err, redis.Nil) {
		return defaultLimit // Gunakan default jika tidak ada custom limit
	} else if err != nil {
		log.Error().Err(err).Msg("Gagal mengambil rate limit dari Redis")
		return defaultLimit
	}
	return limit
}

func setRateLimit(clientID string, limit int) {
	// Simpan limit ke Redis
	err := rdb.HSet(ctx, clientRateLimitThreshold, clientID, limit).Err()
	if err != nil {
		log.Error().Err(err).Msg("Gagal menyimpan rate limit ke Redis")
	}
}

// Goroutine untuk membuang request dari bucket (Leaky Process)
func leakBucket() {
	for {
		time.Sleep(leakRate)

		// Ambil semua key yang berkaitan dengan rate limiting
		keys, err := rdb.Keys(ctx, leakyBucketRateLimit+":*").Result()
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch keys")
			continue
		}

		for _, key := range keys {
			// Hapus 1 request dari setiap bucket
			_, err := rdb.LPop(ctx, key).Result()
			if errors.Is(err, redis.Nil) {
				continue // Tidak ada request untuk dikeluarkan
			} else if err != nil {
				log.Error().Err(err).Str("key", key).Msg("Failed to leak request")
			}
		}
	}
}

func main() {

	// Inisialisasi Viper
	viper.SetDefault("rate_limit", 10)
	defaultRateLimit := viper.GetInt("rate_limit")

	// Inisialisasi logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Jalankan proses Leaky Bucket
	go leakBucket()

	// Inisialisasi Fiber
	app := fiber.New()

	// Terapkan global middleware rate limiting
	app.Use(rateLimiter(defaultRateLimit))

	// Route test
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Request allowed"})
	})

	// Route test dengan middleware
	//app.Get("/", rateLimiter(defaultRateLimit), func(c *fiber.Ctx) error {
	//	return c.JSON(fiber.Map{"message": "Request allowed"})
	//})

	type RateLimit struct {
		Limit    int    `json:"limit"`
		ClientId string `json:"client_id"`
	}

	app.Put("/rate_limit_threshold", func(c *fiber.Ctx) error {
		var limit RateLimit
		if err := c.BodyParser(&limit); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if limit.Limit == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Limit must be greater than 0"})
		}
		if limit.ClientId == "" {
			limit.ClientId = c.IP()
		}
		setRateLimit(limit.ClientId, limit.Limit)

		return c.JSON(fiber.Map{"message": "Rate limit updated"})
	})

	type DefaultRateLimit struct {
		Limit int `json:"limit"`
	}

	app.Put("/default_rate_limit", func(c *fiber.Ctx) error {

		var limit DefaultRateLimit
		if err := c.BodyParser(&limit); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if limit.Limit == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Limit must be greater than 0"})
		}

		viper.Set("rate_limit", limit.Limit)

		return c.JSON(fiber.Map{"message": "Default rate limit updated"})
	})

	// Jalankan server
	log.Info().Msg("Starting Fiber server on :3000")
	log.Fatal().Err(app.Listen(":3000"))
}
