package book

import (
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/book/controller"
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/book/repository"
	"github.com/ciazhar/go-start-small/examples/mongodb_transactional/internal/book/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init(fiberApp fiber.Router, db *mongo.Database) {
	r := repository.NewBookRepository(db)
	s := service.NewBookService(r)
	c := controller.NewBookController(s)

	fiberApp.Post("/book", c.Insert)
}
