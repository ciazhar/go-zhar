package person

import (
	"github.com/ciazhar/go-start-small/examples/mongodb_crud_testcontainers/internal/person/controller"
	"github.com/ciazhar/go-start-small/examples/mongodb_crud_testcontainers/internal/person/repository"
	"github.com/ciazhar/go-start-small/examples/mongodb_crud_testcontainers/internal/person/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init(router fiber.Router, conn *mongo.Database) {
	r := repository.NewPersonRepository(conn)
	s := service.NewPersonService(r)
	c := controller.NewPersonController(s)

	router.Post("/person", c.Insert)
	router.Post("/person/batch", c.InsertBatch)
	router.Get("/person", c.FindAllPageSize)
	router.Get("/person/country", c.FindCountry)
	router.Get("/person/age-range", c.FindAgeRange)
	router.Get("/person/hobby", c.FindHobby)
	router.Get("/person/minified", c.FindMinified)
	router.Get("/person/:id", c.FindOne)
	router.Put("/person/:id", c.Update)
	router.Delete("/person/:id", c.Delete)
}
