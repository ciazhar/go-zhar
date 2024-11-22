package main

import (
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/controller"
	"github.com/ciazhar/go-start-small/examples/temporal_banking_onboarding/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.temporal.io/sdk/client"
	"log"
	"log/slog"
)

func main() {

	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		slog.Error("Unable to create Temporal client:", slog.String("error", err.Error()))
		return
	}
	defer temporalClient.Close()

	s := service.NewOnboardingService(temporalClient)
	c := controller.NewOnboardingController(s)

	app := fiber.New()
	app.Post("/onboarding", c.CreateOnboarding)
	app.Post("/onboarding/:id/signature", c.SignAgreement)
	app.Get("/onboarding/:id", c.GetOnboarding)
	////docs.SwaggerInfo.BasePath = "/"
	//app.Get("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	log.Fatal(app.Listen(":8080"))
}
