package controller

import (
	"github.com/ciazhar/go-zhar/examples/mongodb/crud-testcontainers/internal/person/model"
	"github.com/ciazhar/go-zhar/examples/mongodb/crud-testcontainers/internal/person/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

type PersonController interface {
	Insert(ctx *fiber.Ctx) error
	InsertBatch(ctx *fiber.Ctx) error
	FindOne(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindCountry(ctx *fiber.Ctx) error
	FindAgeRange(ctx *fiber.Ctx) error
	FindHobby(ctx *fiber.Ctx) error
	FindMinified(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type personController struct {
	p service.PersonService
}

func (p personController) Insert(ctx *fiber.Ctx) error {
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

	err = p.p.Insert(&person)
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

func (p personController) InsertBatch(ctx *fiber.Ctx) error {
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

	err = p.p.InsertBatch(&persons)
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

func (p personController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	one, err := p.p.FindOne(id)
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

func (p personController) FindAll(ctx *fiber.Ctx) error {
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

	all, err := p.p.FindAll(name, email, age)
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

func (p personController) FindCountry(ctx *fiber.Ctx) error {
	country := ctx.Query("country")

	all, err := p.p.FindCountry(country)
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

func (p personController) FindAgeRange(ctx *fiber.Ctx) error {
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

	all, err := p.p.FindAgeRange(start, end)
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

func (p personController) FindHobby(ctx *fiber.Ctx) error {
	hobby := ctx.Query("hobby")
	hobbys := strings.Split(hobby, ",")

	findHobby, err := p.p.FindHobby(hobbys)
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

func (p personController) FindMinified(ctx *fiber.Ctx) error {

	all, err := p.p.FindMinified()
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

func (p personController) Update(ctx *fiber.Ctx) error {

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

	err = p.p.Update(id, person)
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

func (p personController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := p.p.Delete(id)
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

func NewPersonController(p service.PersonService) PersonController {
	return &personController{p}
}
