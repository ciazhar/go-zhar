package controller

import (
	"github.com/ciazhar/go-zhar/examples/mongodb/location/internal/location/model"
	"github.com/ciazhar/go-zhar/examples/mongodb/location/internal/location/service"
	"github.com/gofiber/fiber/v2"
)

type LocationController interface {
	Insert(ctx *fiber.Ctx) error
	Nearest(ctx *fiber.Ctx) error
	Top5Nearest(ctx *fiber.Ctx) error
}

type locationController struct {
	l service.LocationService
}

func (l locationController) Top5Nearest(ctx *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (l locationController) Insert(ctx *fiber.Ctx) error {
	var location model.InsertLocationForm
	if err := ctx.BodyParser(&location); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Error on parse request body",
				"data":    err,
			},
		)
	}

	if err := l.l.Insert(location); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Error on insert data",
				"data":    err,
			},
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Data inserted successfully",
			"data":    location,
		},
	)
}

func (l locationController) Nearest(ctx *fiber.Ctx) error {

	var form model.NearestLocationForm
	err := ctx.QueryParser(&form)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Error on parse query",
				"data":    err,
			},
		)
	}

	nearest, err := l.l.Nearest(form.Longitude, form.Latitude, form.MaxDistance, form.Limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Error on get nearest location",
				"data":    err,
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Get nearest location successfully",
			"data":    nearest,
		},
	)
}

func NewLocationController(l service.LocationService) LocationController {
	return &locationController{l}
}
