package controller

import (
	"github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/purchase/service"
	"github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/transaction/model"
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
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    err,
			},
		)
	}

	err = p.p.Purchase(ctx.Context(), &transaction)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Could not purchase book",
				"data":    err,
			},
		)
	}

	return ctx.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Book purchased successfully",
			"data":    transaction,
		},
	)
}

func NewPurchaseController(p *service.PurchaseService) *PurchaseController {
	return &PurchaseController{p: p}
}
