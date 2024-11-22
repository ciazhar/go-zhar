package controller

import (
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/model"
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
)

type OnboardingController struct {
	service *service.OnboardingService
}

func NewOnboardingController(service *service.OnboardingService) *OnboardingController {
	return &OnboardingController{
		service: service,
	}
}

// CreateOnboarding godoc
// @Summary Create onboarding
// @Description Create onboarding
// @Tags onboarding
// @Accept json
// @Produce json
// @Param request body http.OnboardingRequest true "Onboarding request"
// @Success 200 {object} http.OnboardingStatus
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /onboarding [post]
func (c *OnboardingController) CreateOnboarding(ctx *fiber.Ctx) error {
	var request model.OnboardingRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := c.service.CreateOnboarding(ctx.Context(), &request)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// GetOnboarding godoc
// @Summary Get onboarding
// @Description Get onboarding
// @Tags onboarding
// @Accept json
// @Produce json
// @Param id path string true "Onboarding ID"
// @Success 200 {object} http.OnboardingStatus
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /onboarding/{id} [get]
func (c *OnboardingController) GetOnboarding(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := c.service.GetOnboarding(ctx.Context(), id)
	if err != nil {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(response)
}

// SignAgreement godoc
// @Summary Sign agreement
// @Description Sign agreement
// @Tags onboarding
// @Accept json
// @Produce json
// @Param id path string true "Onboarding ID"
// @Param request body http.SignatureRequest true "Signature request"
// @Success 200 {object} http.OnboardingStatus
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /onboarding/{id}/signature [post]
func (c *OnboardingController) SignAgreement(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var request model.SignatureRequest
	err = ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := c.service.SignAgreement(ctx.Context(), id, &request)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(response)
}
