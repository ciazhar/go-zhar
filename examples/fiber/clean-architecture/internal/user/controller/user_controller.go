package controller

import (
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture/internal/user/model"
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture/internal/user/service"
	"github.com/gofiber/fiber/v2"
)

// UserController handles user-related requests
type UserController struct {
	UserService *service.UserService
}

// NewUserController creates a new UserController
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) AddUser(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString("Invalid request")
	}
	uc.UserService.AddUser(user)
	return c.SendStatus(201)
}

func (uc *UserController) GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := uc.UserService.GetUserByUsername(username)
	if err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.JSON(user)
}

func (uc *UserController) GetAllUsers(c *fiber.Ctx) error {
	users := uc.UserService.GetAllUsers()
	return c.JSON(users)
}

func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	username := c.Params("username")
	err := uc.UserService.DeleteUser(username)
	if err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.SendStatus(200)
}

func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString("Invalid request")
	}
	err := uc.UserService.UpdateUser(user)
	if err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.SendStatus(200)
}
