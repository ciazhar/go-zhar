package main

import (
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/controller"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/antifraud-service/checks", controller.AntiFraudChecksHandler)
	app.Post("/user-service/users", controller.UsersHandler)
	app.Post("/agreement-service/agreements", controller.AgreementsHandler)
	app.Post("/signature-service/signatures", controller.SignaturesHandler)
	app.Post("/account-service/accounts", controller.AccountsHandler)
	app.Post("/card-service/cards", controller.CardsHandler)

	log.Fatal(app.Listen(":8081"))
}
