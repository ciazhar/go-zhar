package controller

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/redis_caching_for_every_data_structure/internal/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DailyLoginTrackingController struct {
	repo *repository.BitmapDailyLoginTrackingRepository
}

// NewDailyLoginTrackingController initializes the controller with the repository
func NewDailyLoginTrackingController(repo *repository.BitmapDailyLoginTrackingRepository) *DailyLoginTrackingController {
	return &DailyLoginTrackingController{
		repo: repo,
	}
}

// RegisterRoutes registers the controller routes with the Fiber app
func (c *DailyLoginTrackingController) RegisterRoutes(app *fiber.App) {
	app.Post("/mark-login", c.MarkUserLoginHandler)
	app.Get("/check-login", c.CheckUserLoginHandler)
}

// MarkUserLoginHandler handles marking a user as logged in
func (c *DailyLoginTrackingController) MarkUserLoginHandler(ctx *fiber.Ctx) error {
	// Parse query parameters
	date := ctx.Query("date")
	userIDStr := ctx.Query("userID")

	if date == "" {
		date = time.Now().Format("2006-01-02") // Default to today's date
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid userID parameter")
	}

	// Mark the user login
	if err := c.repo.MarkUserLogin(context.Background(), date, userID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to mark user login: %v", err))
	}

	return ctx.Status(fiber.StatusOK).SendString("User login marked successfully")
}

// CheckUserLoginHandler handles checking if a user logged in on a specific day
func (c *DailyLoginTrackingController) CheckUserLoginHandler(ctx *fiber.Ctx) error {
	// Parse query parameters
	date := ctx.Query("date")
	userIDStr := ctx.Query("userID")

	if date == "" {
		date = time.Now().Format("2006-01-02") // Default to today's date
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid userID parameter")
	}

	// Check if the user logged in
	isLoggedIn, err := c.repo.CheckUserLogin(context.Background(), date, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to check user login: %v", err))
	}

	// Respond with the result
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"date":       date,
		"userID":     userID,
		"isLoggedIn": isLoggedIn == 1,
	})
}
