package controller

import (
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/purchase/service"
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/transaction/model"
	"github.com/ciazhar/go-start-small/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type PurchaseController struct {
	p *service.PurchaseService
}

func (p *PurchaseController) Purchase(ctx *fiber.Ctx) error {

	var transaction model.Transaction
	err := ctx.BodyParser(&transaction)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Review your input",
				Error:   err.Error(),
			},
		)
	}

	err = p.p.Purchase(ctx.Context(), &transaction)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Could not purchase book",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(200).JSON(
		response.Response{
			Message: "Book purchased successfully",
			Data:    transaction,
		},
	)
}

func NewPurchaseController(p *service.PurchaseService) *PurchaseController {
	return &PurchaseController{p: p}
}
