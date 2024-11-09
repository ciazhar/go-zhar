package controller

import (
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/book/model"
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/book/service"
	"github.com/ciazhar/go-start-small/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type BookController struct {
	b *service.BookService
}

func (b *BookController) Insert(ctx *fiber.Ctx) error {

	var book model.Book
	err := ctx.BodyParser(&book)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Review your input",
				Error:   err.Error(),
			},
		)
	}

	err = b.b.Insert(ctx.Context(), &book)
	if err != nil {
		return ctx.Status(500).JSON(
			response.Response{
				Message: "Could not insert book",
				Error:   err.Error(),
			},
		)
	}

	return ctx.Status(200).JSON(
		response.Response{
			Message: "Book inserted successfully",
			Data:    book,
		},
	)
}

func NewBookController(b *service.BookService) *BookController {
	return &BookController{
		b: b,
	}
}
