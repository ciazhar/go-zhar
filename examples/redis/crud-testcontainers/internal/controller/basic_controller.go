package controller

import (
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/service"
	"github.com/gofiber/fiber/v2"
	"time"
)

type BasicController struct {
	aService *service.BasicService
}

func (b *BasicController) Get(ctx *fiber.Ctx) error {
	val, err := b.aService.Get()
	if err != nil {
		return err
	}
	return ctx.JSON(val)
}

func (b *BasicController) Set(ctx *fiber.Ctx) error {
	value := ctx.FormValue("value")
	expirationStr := ctx.FormValue("expiration")
	expiration, err := time.ParseDuration(expirationStr)
	if err != nil {
		// Handle error, perhaps by setting a default expiration
		expiration = time.Hour * 24 // Default expiration of 24 hours
	}
	if err := b.aService.Set(value, expiration); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (b *BasicController) Delete(ctx *fiber.Ctx) error {
	if err := b.aService.Delete(); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (b *BasicController) GetHash(ctx *fiber.Ctx) error {
	field := ctx.Params("field")
	val, err := b.aService.GetHash(field)
	if err != nil {
		return err
	}
	return ctx.JSON(val)
}

func (b *BasicController) SetHash(ctx *fiber.Ctx) error {
	field := ctx.FormValue("field")
	value := ctx.FormValue("value")
	if err := b.aService.SetHash(field, value); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (b *BasicController) SetHashTTL(ctx *fiber.Ctx) error {
	field := ctx.FormValue("field")
	value := ctx.FormValue("value")
	ttlStr := ctx.FormValue("ttl")
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		// Handle error, perhaps by setting a default TTL
		ttl = time.Hour * 24 // Default TTL of 24 hours
	}
	if err := b.aService.SetHashTTL(field, value, ttl); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (b *BasicController) DeleteHash(ctx *fiber.Ctx) error {
	field := ctx.FormValue("field")
	if err := b.aService.DeleteHash(field); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (b *BasicController) GetList(ctx *fiber.Ctx) error {
	val, err := b.aService.GetList()
	if err != nil {
		return err
	}
	return ctx.JSON(val)
}

func (b *BasicController) SetList(ctx *fiber.Ctx) error {
	listStr := ctx.FormValue("list")
	list := make([]string, 0)
	for _, val := range listStr {
		list = append(list, string(val))
	}
	if err := b.aService.SetList(list); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (b *BasicController) DeleteList(ctx *fiber.Ctx) error {
	value := ctx.FormValue("value")
	if err := b.aService.DeleteList(value); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func NewBasicController(aService *service.BasicService) *BasicController {
	return &BasicController{
		aService: aService,
	}
}
