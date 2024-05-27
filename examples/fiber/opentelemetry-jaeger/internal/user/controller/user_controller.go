package controller

import (
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/user/model"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/user/service"
	model2 "github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/pkg/model"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
	_, span := uc.tracer.Start(
		c.Context(),
		"UserController_AddUser")
	defer span.End()
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		span.RecordError(err)
		return c.Status(400).JSON(model2.Response{
			Message: "Invalid request",
			Error:   err.Error(),
			TraceID: span.SpanContext().TraceID().String(),
		})
	}
	uc.UserService.AddUser(c.Context(), user, span)
	return c.Status(201).JSON(model2.Response{
		Message: "User created",
		TraceID: span.SpanContext().TraceID().String(),
	})
}

// GetUserByUsername is the handler for getting a user by their username
func (uc *UserController) GetUserByUsername(c *fiber.Ctx) error {
	// Extract the context from the incoming request
	//ctx := otel.GetTextMapPropagator().Extract(c.Context(), propagation.HeaderCarrier(c.Request().Header))
	ctx := otel.GetTextMapPropagator().Extract(c.Context(), propagation.HeaderCarrier(c.GetReqHeaders()))

	_, span := uc.tracer.Start(
		ctx,
		"UserController_GetUserByUsername",
	)
	defer span.End()
	username := c.Params("username")
	user, err := uc.UserService.GetUserByUsername(c.Context(), username, span)
	if err != nil {
		span.RecordError(err)
		return c.Status(404).JSON(model2.Response{
			Message: "User not found",
			Error:   err.Error(),
			TraceID: span.SpanContext().TraceID().String(),
		})
	}
	return c.Status(200).JSON(model2.Response{
		Message: "User found",
		Data:    user,
		TraceID: span.SpanContext().TraceID().String(),
	})
}

// GetAllUsers is the handler for getting all users
func (uc *UserController) GetAllUsers(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(
		c.Context(),
		"UserController_GetAllUsers",
	)
	defer span.End()
	users := uc.UserService.GetAllUsers(c.Context(), span)
	return c.Status(200).JSON(model2.Response{
		Message: "Users found",
		Data:    users,
		TraceID: span.SpanContext().TraceID().String(),
	})
}

// DeleteUser is the handler for deleting a user
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(
		c.Context(),
		"UserController_DeleteUser",
	)
	defer span.End()
	username := c.Params("username")
	err := uc.UserService.DeleteUser(c.Context(), username, span)
	if err != nil {
		return c.Status(500).JSON(model2.Response{
			Message: "User not found",
			Error:   err.Error(),
			TraceID: span.SpanContext().TraceID().String(),
		})
	}
	return c.Status(200).JSON(model2.Response{
		Message: "User deleted",
		TraceID: span.SpanContext().TraceID().String(),
	})
}

// UpdateUser is the handler for updating a user's information
func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	_, span := uc.tracer.Start(
		c.Context(),
		"UserController_UpdateUser",
	)
	defer span.End()
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(model2.Response{
			Message: "Invalid request",
			Error:   err.Error(),
			TraceID: span.SpanContext().TraceID().String(),
		})
	}
	err := uc.UserService.UpdateUser(c.Context(), user, span)
	if err != nil {
		return c.Status(500).JSON(model2.Response{
			Message: "User not found",
			Error:   err.Error(),
			TraceID: span.SpanContext().TraceID().String(),
		})
	}
	return c.Status(200).JSON(model2.Response{
		Message: "User updated",
		TraceID: span.SpanContext().TraceID().String(),
	})
}
