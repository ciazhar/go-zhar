package controller

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/redis_idempotance_post/internal/model"
	"github.com/ciazhar/go-start-small/examples/redis_idempotance_post/internal/repository"
	"github.com/ciazhar/go-start-small/examples/redis_idempotance_post/pkg"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func CreateProductHandler(c *fiber.Ctx) error {
	// Ambil idempotency key dari header
	idempotencyKey := c.Get("Idempotency-Key")
	if idempotencyKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing Idempotency-Key header",
		})
	}

	// Cek ke Redis apakah key sudah ada
	exists, err := pkg.RedisClient.Exists(context.Background(), idempotencyKey).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check idempotency",
		})
	}

	if exists > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Duplicate request detected",
		})
	}

	// Parse request body
	var req model.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Simulasi insert ke DB (cek SKU unik)
	if _, err := repository.DB.FindProductBySKU(req.SKU); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "SKU already exists",
		})
	}

	// Insert product ke DB
	if err := repository.DB.InsertProduct(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert product",
		})
	}

	// Simpan Idempotency Key ke Redis untuk mencegah duplikat
	err = pkg.RedisClient.Set(context.Background(), idempotencyKey, "used", 5*time.Minute).Err()
	if err != nil {
		log.Println("Failed to store idempotency key:", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
	})
}
