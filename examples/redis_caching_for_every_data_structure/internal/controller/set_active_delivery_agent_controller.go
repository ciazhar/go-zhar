package controller

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/redis_caching_for_every_data_structure/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type DeliveryAgentController struct {
	repo *repository.SetActiveDeliveryAgentRepository
}

func NewDeliveryAgentController(repo *repository.SetActiveDeliveryAgentRepository) *DeliveryAgentController {
	return &DeliveryAgentController{
		repo: repo,
	}
}

// RegisterRoutes registers the routes for the DeliveryAgent controller
func (c *DeliveryAgentController) RegisterRoutes(app *fiber.App) {
	app.Post("/active/agents/:agentID", c.AddActiveAgentHandler)
	app.Get("/active/agents/:agentID", c.IsAgentActiveHandler)
}

// AddActiveAgentHandler adds a delivery agent to the active list
func (c *DeliveryAgentController) AddActiveAgentHandler(ctx *fiber.Ctx) error {
	agentID := ctx.Params("agentID")
	if agentID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "agentID is required",
		})
	}

	err := c.repo.AddActiveAgent(context.Background(), agentID)
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
func (c *DeliveryAgentController) IsAgentActiveHandler(ctx *fiber.Ctx) error {
	agentID := ctx.Params("agentID")
	if agentID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "agentID is required",
		})
	}

	isActive, err := c.repo.IsAgentActive(context.Background(), agentID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"active": isActive,
	})
}
