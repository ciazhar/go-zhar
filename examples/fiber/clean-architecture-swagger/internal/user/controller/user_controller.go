package controller

import (
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture-swagger/internal/user/model"
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture-swagger/internal/user/service"
	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	AddUser(c *fiber.Ctx) error
	GetUserByUsername(c *fiber.Ctx) error
	GetAllUsers(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
}

// UserController handles user-related requests
type userController struct {
	UserService service.UserService
}

// NewUserController creates a new UserController
func NewUserController(userService service.UserService) UserController {
	return &userController{
		UserService: userService,
	}
}

// AddUser @Summary Add a new user
// @Description Add a new user to the system
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.User true "User"
// @Success 201
// @Router /users [post]
func (uc *userController) AddUser(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString("Invalid request")
	}
	uc.UserService.AddUser(user)
	return c.SendStatus(201)
}

// GetUserByUsername is the handler for getting a user by their username
// @Summary Get a user by username
// @Description Get a user by their username
// @Tags User
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} model.User
// @Failure 404
// @Router /users/{username} [get]
func (uc *userController) GetUserByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := uc.UserService.GetUserByUsername(username)
	if err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.JSON(user)
}

// GetAllUsers is the handler for getting all users
// @Summary Get all users
// @Description Retrieve all users
// @Tags User
// @Produce json
// @Success 200 {object} map[string]model.User
// @Router /users [get]
func (uc *userController) GetAllUsers(c *fiber.Ctx) error {
	users := uc.UserService.GetAllUsers()
	return c.JSON(users)
}

// DeleteUser is the handler for deleting a user
// @Summary Delete a user
// @Description Delete a user by their username
// @Tags User
// @Param username path string true "Username"
// @Success 200
// @Failure 404
// @Router /users/{username} [delete]
func (uc *userController) DeleteUser(c *fiber.Ctx) error {
	username := c.Params("username")
	err := uc.UserService.DeleteUser(username)
	if err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.SendStatus(200)
}

// UpdateUser is the handler for updating a user's information
// @Summary Update a user
// @Description Update a user's details
// @Tags User
// @Param user body model.User true "User"
// @Success 200
// @Failure 404
// @Router /users [put]
func (uc *userController) UpdateUser(c *fiber.Ctx) error {
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
