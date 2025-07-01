package controller

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/redis_caching_for_every_data_structure/internal/repository"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type AgentController struct {
	repo  *repository.GeoTrackingDeliveryAgentLocationRepository
	repo2 *repository.SetActiveDeliveryAgentRepository
}

// NewAgentController initializes the controller with the repository
func NewAgentController(
	repo *repository.GeoTrackingDeliveryAgentLocationRepository,
	repo2 *repository.SetActiveDeliveryAgentRepository,
) *AgentController {
	return &AgentController{
		repo:  repo,
		repo2: repo2,
	}
}

// RegisterRoutes registers the routes for the GeoTracking controller
func (c *AgentController) RegisterRoutes(app *fiber.App) {
	app.Post("/add-location", c.AddAgentLocationHandler)
	app.Get("/nearby-agents", c.GetNearbyAgentsHandler)
	app.Post("/active/agents/:agentID", c.AddActiveAgentHandler)
	app.Get("/active/agents/:agentID", c.IsAgentActiveHandler)
}

// AddAgentLocationHandler handles adding an agent's location
func (c *AgentController) AddAgentLocationHandler(ctx *fiber.Ctx) error {
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
func (c *AgentController) GetNearbyAgentsHandler(ctx *fiber.Ctx) error {
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

// AddActiveAgentHandler adds a delivery agent to the active list
func (c *AgentController) AddActiveAgentHandler(ctx *fiber.Ctx) error {
	agentID := ctx.Params("agentID")
	if agentID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "agentID is required",
		})
	}

	err := c.repo2.AddActiveAgent(context.Background(), agentID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Agent added successfully",
	})
}

// IsAgentActiveHandler checks if a delivery agent is active
func (c *AgentController) IsAgentActiveHandler(ctx *fiber.Ctx) error {
	agentID := ctx.Params("agentID")
	if agentID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "agentID is required",
		})
	}

	isActive, err := c.repo2.IsAgentActive(context.Background(), agentID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"active": isActive,
	})
}
