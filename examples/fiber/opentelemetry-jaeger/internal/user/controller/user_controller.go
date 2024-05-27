package controller

import (
	"fmt"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/user/model"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/user/service"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

// UserController handles user-related requests
type UserController struct {
	UserService *service.UserService
	tracer      trace.Tracer
}

// NewUserController creates a new UserController
func NewUserController(
	userService *service.UserService,
	tracer trace.Tracer,
) *UserController {
	return &UserController{
		UserService: userService,
		tracer:      tracer,
	}
}

// AddUser @Summary Add a new user
func (uc *UserController) AddUser(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(c.Context(), "AddUser")
	defer span.End()
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString("Invalid request")
	}
	uc.UserService.AddUser(c.Context(), user, span)
	return c.SendStatus(201)
}

// GetUserByUsername is the handler for getting a user by their username
func (uc *UserController) GetUserByUsername(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(c.Context(), "GetUserByUsername")
	defer span.End()
	username := c.Params("username")
	user, err := uc.UserService.GetUserByUsername(c.Context(), username)
	if err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.JSON(user)
}

// GetAllUsers is the handler for getting all users
func (uc *UserController) GetAllUsers(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(c.Context(), "GetAllUsersController")
	defer span.End()
	users := uc.UserService.GetAllUsers(c.Context(), span)
	fmt.Println(span.SpanContext().TraceID())
	return c.JSON(users)
}

// DeleteUser is the handler for deleting a user
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(c.Context(), "DeleteUser")
	defer span.End()
	username := c.Params("username")
	err := uc.UserService.DeleteUser(c.Context(), username)
	if err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.SendStatus(200)
}

// UpdateUser is the handler for updating a user's information
func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(c.Context(), "UpdateUser")
	defer span.End()
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString("Invalid request")
	}
	err := uc.UserService.UpdateUser(c.Context(), user)
	if err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.SendStatus(200)
}
