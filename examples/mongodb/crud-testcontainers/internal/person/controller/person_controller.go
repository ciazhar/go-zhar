package controller

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/mongodb/crud-testcontainers/internal/person/model"
	"github.com/ciazhar/go-zhar/examples/mongodb/crud-testcontainers/internal/person/service"
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
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    err,
			},
		)
	}

	err = p.p.Insert(context.Background(), &person)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Could not insert person",
				"data":    err,
			},
		)
	}

	return ctx.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Person inserted successfully",
			"data":    person,
		},
	)
}

func (p *PersonController) InsertBatch(ctx *fiber.Ctx) error {
	var persons []model.Person
	err := ctx.BodyParser(&persons)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    err,
			},
		)
	}

	err = p.p.InsertBatch(context.Background(), &persons)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Could not insert persons",
				"data":    err,
			},
		)
	}

	return ctx.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Persons inserted successfully",
			"data":    persons,
		},
	)

}

func (p *PersonController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	one, err := p.p.FindOne(context.Background(), id)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Could not find person",
				"data":    err,
			},
		)
	}

	return ctx.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Person found successfully",
			"data":    one,
		},
	)
}

func (p *PersonController) FindAllPageSize(ctx *fiber.Ctx) error {
	page := ctx.Query("page")
	size := ctx.Query("size")
	sort := ctx.Query("sort")
	pageI, err := strconv.Atoi(page)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    err,
			},
		)
	}
	sizeI, err := strconv.Atoi(size)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    err,
			},
		)
	}
	name := ctx.Query("name")
	email := ctx.Query("email")
	ageS := ctx.Query("age")

	age := 0
	if ageS != "" {
		ageI, err := strconv.Atoi(ageS)
		if err != nil {
			return ctx.Status(500).JSON(
				fiber.Map{
					"status":  "error",
					"message": "Review your input",
					"data":    err,
				},
			)
		}
		age = ageI
	}

	all, err := p.p.FindAllPageSize(context.Background(), pageI, sizeI, sort, name, email, age)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Could not find persons",
				"data":    err,
			},
		)
	}

	return ctx.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Persons found successfully",
			"data":    all,
		},
	)
}

func (p *PersonController) FindCountry(ctx *fiber.Ctx) error {
	country := ctx.Query("country")

	all, err := p.p.FindCountry(context.Background(), country)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Could not find persons",
				"data":    err,
			},
		)
	}

	return ctx.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Persons found successfully",
			"data":    all,
		},
	)
}

func (p *PersonController) FindAgeRange(ctx *fiber.Ctx) error {
	startAge := ctx.Query("startAge")
	endAge := ctx.Query("endAge")

	start, err := strconv.Atoi(startAge)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    err,
			},
		)
	}

	end, err := strconv.Atoi(endAge)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    err,
			},
		)
	}

	all, err := p.p.FindAgeRange(context.Background(), start, end)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Could not find persons",
				"data":    err,
			},
		)
	}

	return ctx.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Persons found successfully",
			"data":    all,
		},
	)
}

func (p *PersonController) FindHobby(ctx *fiber.Ctx) error {
	hobby := ctx.Query("hobby")
	hobbys := strings.Split(hobby, ",")

	findHobby, err := p.p.FindHobby(context.Background(), hobbys)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Could not find persons",
				"data":    err,
			},
		)
	}

	return ctx.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Persons found successfully",
			"data":    findHobby,
		},
	)
}

func (p *PersonController) FindMinified(ctx *fiber.Ctx) error {

	all, err := p.p.FindMinified(context.Background())
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Could not find persons",
				"data":    err,
			},
		)
	}

	return ctx.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Persons found successfully",
			"data":    all,
		},
	)

}

func (p *PersonController) Update(ctx *fiber.Ctx) error {

	id := ctx.Params("id")

	var person model.UpdatePersonForm
	err := ctx.BodyParser(&person)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    err,
			},
		)
	}

	err = p.p.Update(context.Background(), id, person)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Could not update person",
				"data":    err,
			},
		)
	}

	return ctx.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Person updated successfully",
			"data":    person,
		},
	)
}

func (p *PersonController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := p.p.Delete(context.Background(), id)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Could not delete person",
				"data":    err,
			},
		)
	}

	return ctx.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Person deleted successfully",
		},
	)
}

func NewPersonController(p *service.PersonService) *PersonController {
	return &PersonController{p}
}
