package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/db_pg_csv_zip_http/internal/model"
	"github.com/ciazhar/go-start-small/examples/db_pg_csv_zip_http/internal/service"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/response"
	"github.com/gofiber/fiber/v2"
	"os"
)

type Controller interface {
	SendEmail(ctx *fiber.Ctx) error
	ExportUnoptimized(ctx *fiber.Ctx) error
	ExportOptimized(ctx *fiber.Ctx) error
}

type EmailController struct {
	s service.Service
}

func (e EmailController) ExportUnoptimized(ctx *fiber.Ctx) error {
	err := e.s.ExportAndSendUnoptimized(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Failed to export data",
			Data: logger.LogAndReturnWarning(ctx.Context(), err,
				"Failed to export data",
				map[string]interface{}{}),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Response{
		Error: "",
		Data: map[string]interface{}{
			"message": "Success",
		},
	})
}

func (e EmailController) ExportOptimized(ctx *fiber.Ctx) error {
	err := e.s.ExportAndSendOptimizedCopy(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Failed to export data",
			Data: logger.LogAndReturnWarning(ctx.Context(), err,
				"Failed to export data",
				map[string]interface{}{}),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Response{
		Error: "",
		Data: map[string]interface{}{
			"message": "Success",
		},
	})
}

func (e EmailController) SendEmail(ctx *fiber.Ctx) error {

	var body model.JsonBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Response{
			Error: "Invalid input",
			Data: logger.LogAndReturnWarning(ctx.Context(), err,
				"Invalid input",
				map[string]interface{}{}),
		})
	}
	fmt.Println("Received email..." + body.FileName)

	decoded, err := base64.StdEncoding.DecodeString(body.Base64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Failed to decode base64",
			Data: logger.LogAndReturnWarning(ctx.Context(), err,
				"Failed to decode base64",
				map[string]interface{}{}),
		})
	}

	filePath := fmt.Sprintf("./%s", body.FileName)
	file, err := os.Create("/Users/ciazhar/GolandProjects/go-start-small/datasets/diabetes_clinical_100k/" + filePath)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Failed to create file",
			Data: logger.LogAndReturnWarning(ctx.Context(), err,
				"Failed to create file",
				map[string]interface{}{}),
		})
	}
	defer file.Close()

	if _, err := file.Write(decoded); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Failed to write to file",
			Data: logger.LogAndReturnWarning(ctx.Context(), err,
				"Failed to write to file",
				map[string]interface{}{}),
		})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func NewEmailController(s service.Service) Controller {
	return &EmailController{
		s: s,
	}
}
