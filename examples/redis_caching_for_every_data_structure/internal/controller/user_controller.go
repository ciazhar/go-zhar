package controller

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/redis_caching_for_every_data_structure/internal/repository"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
	"time"
)

type UserController struct {
	repo  *repository.BitmapDailyLoginTrackingRepository
	repo2 *repository.HashUserProfileRepository
	repo3 *repository.StringUserSessionRepository
}

// NewUserController initializes the controller with the repository
func NewUserController(repo *repository.BitmapDailyLoginTrackingRepository, repo2 *repository.HashUserProfileRepository, repo3 *repository.StringUserSessionRepository) *UserController {
	return &UserController{
		repo:  repo,
		repo2: repo2,
		repo3: repo3,
	}
}

// RegisterRoutes registers the controller routes with the Fiber app
func (c *UserController) RegisterRoutes(app *fiber.App) {

	user := app.Group("/user")
	user.Post("/profile", c.SetUserProfileHandler)
	user.Get("/profile/:userID", c.GetUserProfileHandler)
	user.Post("/mark-login", c.MarkUserLoginHandler)
	user.Get("/check-login", c.CheckUserLoginHandler)
	user.Post("/:userID/session", c.SetSession)
	user.Get("/:userID/session", c.GetSession)
}

// SetUserProfileHandler handles setting user profile details
func (c *UserController) SetUserProfileHandler(ctx *fiber.Ctx) error {
	userID := ctx.Query("userID")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required parameter: userID")
	}

	var profile map[string]interface{}
	if err := ctx.BodyParser(&profile); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := c.repo2.SetUserProfile(context.Background(), userID, profile); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to set user profile")
	}

	return ctx.Status(fiber.StatusOK).SendString("User profile set successfully")
}

// GetUserProfileHandler handles retrieving user profile details
func (c *UserController) GetUserProfileHandler(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required parameter: userID")
	}

	profile, err := c.repo2.GetUserProfile(context.Background(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user profile")
	}

	return ctx.Status(fiber.StatusOK).JSON(profile)
}

// MarkUserLoginHandler handles marking a user as logged in
func (c *UserController) MarkUserLoginHandler(ctx *fiber.Ctx) error {
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
func (c *UserController) CheckUserLoginHandler(ctx *fiber.Ctx) error {
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

// SetSession handles the request to create a user session
func (c *UserController) SetSession(ctx *fiber.Ctx) error {
	type request struct {
		UserID       string `json:"user_id"`
		SessionToken string `json:"session_token"`
	}

	var req request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if req.UserID == "" || req.SessionToken == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id and session_token are required",
		})
	}

	if err := c.repo3.SetUserSession(context.Background(), req.UserID, req.SessionToken); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to set user session",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Session created successfully",
	})
}

// GetSession handles the request to retrieve a user session
func (c *UserController) GetSession(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")
	if userID == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "userID is required",
		})
	}

	sessionToken, err := c.repo3.GetUserSession(context.Background(), userID)
	if err != nil {
		if err.Error() == "redis: nil" {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Session not found",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user session",
		})
	}

	return ctx.JSON(fiber.Map{
		"user_id":       userID,
		"session_token": sessionToken,
	})
}
