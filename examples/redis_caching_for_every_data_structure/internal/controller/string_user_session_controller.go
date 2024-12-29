package controller

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/redis_caching_for_every_data_structure/internal/repository"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type UserSessionController struct {
	repo *repository.StringUserSessionRepository
}

// NewUserSessionController creates a new instance of UserSessionController
func NewUserSessionController(repo *repository.StringUserSessionRepository) *UserSessionController {
	return &UserSessionController{repo: repo}
}

func (c *UserSessionController) RegisterRoutes(app *fiber.App) {
	app.Post("/user/:userID/session", c.SetSession)
	app.Get("/user/:userID/session", c.GetSession)
}

// SetSession handles the request to create a user session
func (c *UserSessionController) SetSession(ctx *fiber.Ctx) error {
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

	if err := c.repo.SetUserSession(context.Background(), req.UserID, req.SessionToken); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to set user session",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Session created successfully",
	})
}

// GetSession handles the request to retrieve a user session
func (c *UserSessionController) GetSession(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")
	if userID == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "userID is required",
		})
	}

	sessionToken, err := c.repo.GetUserSession(context.Background(), userID)
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
