package controller

import (
	"context"
	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/auth_service/model"
	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/auth_service/service"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/logger"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/response"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/token_util"
	"github.com/gofiber/fiber/v2"
	"time"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(
	authService *service.AuthService,
) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// RegisterUser Register User
func (c *AuthController) RegisterUser(ctx *fiber.Ctx) error {

	startTime := time.Now()
	newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	newCtx = context.WithValue(newCtx, "requestID", ctx.Locals("requestID").(string))
	defer cancel()

	var user model.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Response{
			Error: "Invalid input",
			Data: logger.LogAndReturnWarning(ctx.Context(), err,
				"Invalid input",
				map[string]string{
					"username": user.Username,
				},
			),
		})
	}

	logger.LogInfo(newCtx,
		"Register user request started",
		map[string]string{
			"username": user.Username,
		})

	err := c.authService.RegisterUser(newCtx, user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Could not register user",
			Data:  err.Error(),
		})
	}

	logger.LogInfo(newCtx,
		"Register user request completed",
		map[string]string{
			"username": user.Username,
			"duration": time.Since(startTime).String(),
		})

	return ctx.Status(fiber.StatusCreated).JSON(response.Response{
		RequestID: ctx.Locals("requestID").(string),
		Message:   "User registered successfully",
	})
}

// Login User
func (c *AuthController) Login(ctx *fiber.Ctx) error {

	startTime := time.Now()
	newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	newCtx = context.WithValue(newCtx, "requestID", ctx.Locals("requestID").(string))
	defer cancel()

	var body model.LoginRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Response{
			Error: "Invalid input",
			Data: logger.LogAndReturnWarning(ctx.Context(), err,
				"Invalid input",
				map[string]string{
					"username": body.Username,
					"password": body.Password,
				}),
		})
	}

	logger.LogInfo(newCtx,
		"Login request started",
		map[string]string{
			"username": body.Username,
		})

	login, err := c.authService.Login(newCtx, body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Could not login user",
			Data:  err.Error(),
		})
	}

	logger.LogInfo(newCtx,
		"Login request completed",
		map[string]string{
			"username":      body.Username,
			"access_token":  login.AccessToken,
			"refresh_token": login.RefreshToken,
			"duration":      time.Since(startTime).String(),
		})

	return ctx.JSON(response.Response{
		RequestID: ctx.Locals("requestID").(string),
		Message:   "User logged in successfully",
		Data:      login,
	})
}

// RefreshToken Refresh Token Handler
func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {

	startTime := time.Now()
	newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	newCtx = context.WithValue(newCtx, "requestID", ctx.Locals("requestID").(string))
	defer cancel()

	token, err := token_util.ExtractToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Error: "No token provided",
			Data:  logger.LogAndReturnWarning(ctx.Context(), err, "No token provided", nil),
		})
	}

	newToken, err := c.authService.RefreshToken(newCtx, token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Could not refresh token",
			Data:  err.Error(),
		})
	}

	logger.LogInfo(newCtx,
		"Refresh token request completed",
		map[string]string{
			"refresh_token": token,
			"access_token":  newToken.AccessToken,
			"duration":      time.Since(startTime).String(),
		})

	return ctx.JSON(response.Response{
		RequestID: ctx.Locals("requestID").(string),
		Message:   "Token refreshed successfully",
		Data:      newToken,
	})
}

// Protected route example
func (c *AuthController) Protected(ctx *fiber.Ctx) error {

	startTime := time.Now()
	newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	newCtx = context.WithValue(newCtx, "requestID", ctx.Locals("requestID").(string))
	defer cancel()

	token, err := token_util.ExtractToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Error: "No token provided",
			Data:  logger.LogAndReturnWarning(ctx.Context(), err, "No token provided", nil),
		})
	}

	logger.LogInfo(newCtx, "Protected route request started", nil)

	// Check if the token exists in Redis
	userId, err := c.authService.Protected(newCtx, token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Could not validate token",
			Data:  err.Error(),
		})
	}

	logger.LogInfo(newCtx, "Protected route request completed", map[string]string{
		"user_id":  userId,
		"duration": time.Since(startTime).String(),
	})

	return ctx.JSON(response.Response{
		RequestID: ctx.Locals("requestID").(string),
		Message:   "Token validated successfully",
		Data:      model.ProtectedResponse{UserId: userId},
	})
}

// Logout Handler
func (c *AuthController) Logout(ctx *fiber.Ctx) error {

	startTime := time.Now()
	newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	newCtx = context.WithValue(newCtx, "requestID", ctx.Locals("requestID").(string))
	defer cancel()

	token, err := token_util.ExtractToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Error: "No token provided",
			Data:  logger.LogAndReturnWarning(ctx.Context(), err, "No token provided", nil),
		})
	}

	logger.LogInfo(newCtx, "Logout request started", map[string]string{
		"token":    token,
		"duration": time.Since(startTime).String(),
	})

	if err := c.authService.Logout(newCtx, token); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Could not logout",
			Data:  err.Error(),
		})
	}

	logger.LogInfo(newCtx, "Logout request completed", map[string]string{
		"token":    token,
		"duration": time.Since(startTime).String(),
	})

	return ctx.JSON(response.Response{
		RequestID: ctx.Locals("requestID").(string),
		Message:   "User logged out successfully",
	})
}

// Revoke Handler
func (c *AuthController) Revoke(ctx *fiber.Ctx) error {

	startTime := time.Now()
	newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	newCtx = context.WithValue(newCtx, "requestID", ctx.Locals("requestID").(string))
	defer cancel()

	token, err := token_util.ExtractToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.Response{
			Error: "No token provided",
			Data:  logger.LogAndReturnWarning(ctx.Context(), err, "No token provided", nil),
		})
	}

	if err := c.authService.Revoke(newCtx, token); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Error: "Could not revoke tokens",
			Data:  err.Error(),
		})
	}

	logger.LogInfo(newCtx, "Revoke request completed", map[string]string{
		"token":    token,
		"duration": time.Since(startTime).String(),
	})

	return ctx.JSON(response.Response{
		RequestID: ctx.Locals("requestID").(string),
		Message:   "Tokens revoked successfully",
	})
}
