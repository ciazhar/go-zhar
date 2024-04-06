package location

import (
	"github.com/ciazhar/go-zhar/examples/mongodb/location/internal/location/controller"
	"github.com/ciazhar/go-zhar/examples/mongodb/location/internal/location/repository"
	"github.com/ciazhar/go-zhar/examples/mongodb/location/internal/location/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init(router fiber.Router, conn *mongo.Database) {
	r := repository.NewLocationRepository(conn)
	s := service.NewLocationService(r)
	c := controller.NewLocationController(s)

	router.Post("/location", c.Insert)
	router.Get("/location", c.Nearest)
}
