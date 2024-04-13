package controller

import (
	"github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/book/model"
	"github.com/ciazhar/go-zhar/examples/mongodb/transactional/internal/book/service"
	"github.com/gofiber/fiber/v2"
)

type BookController interface {
	Insert(ctx *fiber.Ctx) error
}

type bookController struct {
	b service.BookService
}

func (b bookController) Insert(ctx *fiber.Ctx) error {

	var book model.Book
	err := ctx.BodyParser(&book)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"data":    err,
			},
		)
	}

	err = b.b.Insert(ctx.Context(), &book)
	if err != nil {
		return ctx.Status(500).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Could not insert book",
				"data":    err,
			},
		)
	}

	return ctx.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Book inserted successfully",
			"data":    book,
		},
	)
}

func NewBookController(b service.BookService) BookController {
	return &bookController{
		b: b,
	}
}
