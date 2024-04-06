package controller

import (
	"github.com/ciazhar/go-zhar/examples/redis/clean-architecture/internal/basic/service"
	"github.com/gofiber/fiber/v2"
)

type BasicController interface {
	GetBasicHash(ctx *fiber.Ctx) error
	SetBasicHash(ctx *fiber.Ctx) error
}

type basicController struct {
	aService service.BasicService
}

func (b basicController) GetBasicHash(ctx *fiber.Ctx) error {
	key := ctx.Params("key")

	hash, err := b.aService.GetBasicHash(key)
	if err != nil {
		return err
	}

	return ctx.JSON(hash)
}

func (b basicController) SetBasicHash(ctx *fiber.Ctx) error {

	key := ctx.FormValue("key")
	value := ctx.FormValue("value")

	if err := b.aService.SetBasicHash(key, value); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func NewBasicController(aService service.BasicService) BasicController {
	return basicController{
		aService: aService,
	}
}
