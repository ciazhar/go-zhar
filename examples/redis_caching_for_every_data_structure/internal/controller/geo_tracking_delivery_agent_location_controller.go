package controller

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/redis_caching_for_every_data_structure/internal/repository"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type GeoTrackingController struct {
	repo *repository.GeoTrackingDeliveryAgentLocationRepository
}

// NewGeoTrackingController initializes the controller with the repository
func NewGeoTrackingController(repo *repository.GeoTrackingDeliveryAgentLocationRepository) *GeoTrackingController {
	return &GeoTrackingController{
		repo: repo,
	}
}

// RegisterRoutes registers the routes for the GeoTracking controller
func (c *GeoTrackingController) RegisterRoutes(app *fiber.App) {
	app.Post("/add-location", c.AddAgentLocationHandler)
	app.Get("/nearby-agents", c.GetNearbyAgentsHandler)
}

// AddAgentLocationHandler handles adding an agent's location
func (c *GeoTrackingController) AddAgentLocationHandler(ctx *fiber.Ctx) error {
	agentID := ctx.Query("agentID")
	latStr := ctx.Query("lat")
	longStr := ctx.Query("long")

	if agentID == "" || latStr == "" || longStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required parameters: agentID, lat, long")
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid latitude value")
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid longitude value")
	}

	if err := c.repo.AddAgentLocation(context.Background(), agentID, lat, long); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to add agent location")
	}

	return ctx.Status(fiber.StatusOK).SendString("Agent location added successfully")
}

// GetNearbyAgentsHandler handles retrieving nearby agents
func (c *GeoTrackingController) GetNearbyAgentsHandler(ctx *fiber.Ctx) error {
	latStr := ctx.Query("lat")
	longStr := ctx.Query("long")

	if latStr == "" || longStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required parameters: lat, long")
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid latitude value")
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid longitude value")
	}

	agents, err := c.repo.GetNearbyAgents(context.Background(), lat, long)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve nearby agents")
	}

	response := make([]fiber.Map, len(agents))
	for i, agent := range agents {
		response[i] = fiber.Map{
			"agentID":   agent.Name,
			"latitude":  agent.Latitude,
			"longitude": agent.Longitude,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
