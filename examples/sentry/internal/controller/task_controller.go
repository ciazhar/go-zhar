package controller

import (
	"github.com/ciazhar/go-start-small/examples/sentry/internal/model"
	"github.com/ciazhar/go-start-small/examples/sentry/internal/repository"
	"github.com/gofiber/fiber/v2"
)

// RegisterTaskRoutes sets up routes related to tasks
func RegisterTaskRoutes(app *fiber.App) {
	app.Post("/tasks/create", createTask)
	app.Get("/tasks", getAllTasks)
}

func createTask(c *fiber.Ctx) error {
	var task model.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := repository.CreateTask(task); err != nil {
		if err.Error() == "task with this title already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Task created successfully"})
}

func getAllTasks(c *fiber.Ctx) error {
	tasks, err := repository.GetAllTasks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tasks"})
	}
	return c.JSON(tasks)
}
