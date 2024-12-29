package controller

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/redis_caching_for_every_data_structure/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type UserProfileController struct {
	repo *repository.HashUserProfileRepository
}

// NewUserProfileController initializes the controller with the repository
func NewUserProfileController(repo *repository.HashUserProfileRepository) *UserProfileController {
	return &UserProfileController{
		repo: repo,
	}
}

// RegisterRoutes registers the routes for the UserProfile controller
func (c *UserProfileController) RegisterRoutes(app *fiber.App) {
	app.Post("/user/profile", c.SetUserProfileHandler)
	app.Get("/user/profile/:userID", c.GetUserProfileHandler)
}

// SetUserProfileHandler handles setting user profile details
func (c *UserProfileController) SetUserProfileHandler(ctx *fiber.Ctx) error {
	userID := ctx.Query("userID")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required parameter: userID")
	}

	var profile map[string]interface{}
	if err := ctx.BodyParser(&profile); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := c.repo.SetUserProfile(context.Background(), userID, profile); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to set user profile")
	}

	return ctx.Status(fiber.StatusOK).SendString("User profile set successfully")
}

// GetUserProfileHandler handles retrieving user profile details
func (c *UserProfileController) GetUserProfileHandler(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required parameter: userID")
	}

	profile, err := c.repo.GetUserProfile(context.Background(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user profile")
	}

	return ctx.Status(fiber.StatusOK).JSON(profile)
}
