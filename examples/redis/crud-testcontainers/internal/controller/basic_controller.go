package controller

import (
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/service"
	"github.com/gofiber/fiber/v2"
	"time"
)

type BasicController interface {
	Get(ctx *fiber.Ctx) error
	Set(ctx *fiber.Ctx) error
	GetHash(ctx *fiber.Ctx) error
	SetHash(ctx *fiber.Ctx) error
	SetHashTTL(ctx *fiber.Ctx) error
	DeleteHash(ctx *fiber.Ctx) error
}

type basicController struct {
	aService service.BasicService
}

func (b basicController) Get(ctx *fiber.Ctx) error {
	val, err := b.aService.Get()
	if err != nil {
		return err
	}
	return ctx.JSON(val)
}

func (b basicController) Set(ctx *fiber.Ctx) error {
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

func (b basicController) GetHash(ctx *fiber.Ctx) error {
	field := ctx.Params("field")
	val, err := b.aService.GetHash(field)
	if err != nil {
		return err
	}
	return ctx.JSON(val)
}

func (b basicController) SetHash(ctx *fiber.Ctx) error {
	field := ctx.FormValue("field")
	value := ctx.FormValue("value")
	if err := b.aService.SetHash(field, value); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (b basicController) SetHashTTL(ctx *fiber.Ctx) error {
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

func (b basicController) DeleteHash(ctx *fiber.Ctx) error {
	field := ctx.FormValue("field")
	if err := b.aService.DeleteHash(field); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func NewBasicController(aService service.BasicService) BasicController {
	return basicController{
		aService: aService,
	}
}
