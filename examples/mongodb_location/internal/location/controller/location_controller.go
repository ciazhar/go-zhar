package controller

import (
	"github.com/ciazhar/go-start-small/examples/mongodb_location/internal/location/model"
	"github.com/ciazhar/go-start-small/examples/mongodb_location/internal/location/service"
	"github.com/ciazhar/go-start-small/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type LocationController struct {
	l *service.LocationService
}

func (l *LocationController) Top5Nearest(ctx *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (l *LocationController) Insert(ctx *fiber.Ctx) error {
	var location model.InsertLocationForm
	if err := ctx.BodyParser(&location); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Response{
				Message: "Error on parse request body",
				Error:   err.Error(),
			},
		)
	}

	if err := l.l.Insert(location); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			response.Response{
				Message: "Error on insert data",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		response.Response{
			Message: "Data inserted successfully",
			Data:    location,
		},
	)
}

func (l *LocationController) Nearest(ctx *fiber.Ctx) error {

	var form model.NearestLocationForm
	err := ctx.QueryParser(&form)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Response{
				Message: "Error on parse query",
				Error:   err.Error(),
			},
		)
	}

	nearest, err := l.l.Nearest(form.Longitude, form.Latitude, form.MaxDistance, form.Limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			response.Response{
				Message: "Error on get nearest location",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		response.Response{
			Message: "Get nearest location successfully",
			Data:    nearest,
		},
	)
}

func NewLocationController(l *service.LocationService) *LocationController {
	return &LocationController{l}
}
