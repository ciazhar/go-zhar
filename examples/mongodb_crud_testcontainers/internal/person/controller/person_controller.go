package controller

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/mongodb_crud_testcontainers/internal/person/model"
	"github.com/ciazhar/go-start-small/examples/mongodb_crud_testcontainers/internal/person/service"
	"github.com/ciazhar/go-start-small/pkg/response"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

type PersonController struct {
	p *service.PersonService
}

func (p *PersonController) Insert(ctx *fiber.Ctx) error {
	var person model.Person
	err := ctx.BodyParser(&person)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Review your input",
				Error:   err.Error(),
			})
	}

	err = p.p.Insert(context.Background(), &person)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Could not insert person",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(200).JSON(
		response.Response{
			Message: "Person inserted successfully",
			Data:    person,
		},
	)
}

func (p *PersonController) InsertBatch(ctx *fiber.Ctx) error {
	var persons []model.Person
	err := ctx.BodyParser(&persons)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Review your input",
				Error:   err.Error(),
			},
		)
	}

	err = p.p.InsertBatch(context.Background(), &persons)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Could not insert persons",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(200).JSON(
		response.Response{
			Message: "Persons inserted successfully",
			Data:    persons,
		},
	)
}

func (p *PersonController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	one, err := p.p.FindOne(context.Background(), id)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Could not find person",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(200).JSON(
		response.Response{
			Message: "Person found successfully",
			Data:    one,
		},
	)
}

func (p *PersonController) FindAllPageSize(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	size := ctx.QueryInt("size", 10)
	sort := ctx.Query("sort")
	name := ctx.Query("name")
	email := ctx.Query("email")
	age := ctx.QueryInt("age")
	all, err := p.p.FindAllPageSize(context.Background(), page, size, sort, name, email, age)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Could not find persons",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(200).JSON(
		response.Response{
			Message: "Persons found successfully",
			Data:    all,
		},
	)
}

func (p *PersonController) FindCountry(ctx *fiber.Ctx) error {
	country := ctx.Query("country")

	all, err := p.p.FindCountry(context.Background(), country)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Could not find persons",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(200).JSON(
		response.Response{
			Message: "Persons found successfully",
			Data:    all,
		},
	)
}

func (p *PersonController) FindAgeRange(ctx *fiber.Ctx) error {
	startAge := ctx.Query("startAge")
	endAge := ctx.Query("endAge")

	start, err := strconv.Atoi(startAge)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Review your input",
				Error:   err.Error(),
			},
		)
	}

	end, err := strconv.Atoi(endAge)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Review your input",
				Error:   err.Error(),
			},
		)
	}

	all, err := p.p.FindAgeRange(context.Background(), start, end)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Could not find persons",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(200).JSON(
		response.Response{
			Message: "Persons found successfully",
			Data:    all,
		},
	)
}

func (p *PersonController) FindHobby(ctx *fiber.Ctx) error {
	hobby := ctx.Query("hobby")
	hobbys := strings.Split(hobby, ",")

	findHobby, err := p.p.FindHobby(context.Background(), hobbys)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Could not find persons",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(200).JSON(
		response.Response{
			Message: "Persons found successfully",
			Data:    findHobby,
		},
	)
}

func (p *PersonController) FindMinified(ctx *fiber.Ctx) error {

	all, err := p.p.FindMinified(context.Background())
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Could not find persons",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(200).JSON(
		response.Response{
			Message: "Persons found successfully",
			Data:    all,
		},
	)

}

func (p *PersonController) Update(ctx *fiber.Ctx) error {

	id := ctx.Params("id")

	var person model.UpdatePersonForm
	err := ctx.BodyParser(&person)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Review your input",
				Error:   err.Error(),
			},
		)
	}

	err = p.p.Update(context.Background(), id, person)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Could not update person",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(200).JSON(
		response.Response{
			Message: "Person updated successfully",
			Data:    person,
		},
	)
}

func (p *PersonController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := p.p.Delete(context.Background(), id)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Could not delete person",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(200).JSON(
		response.Response{
			Message: "Person deleted successfully",
		},
	)
}

func NewPersonController(p *service.PersonService) *PersonController {
	return &PersonController{p}
}
