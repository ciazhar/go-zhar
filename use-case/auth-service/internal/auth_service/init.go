package auth_service

import (
	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/auth_service/controller"
	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/auth_service/repository"
	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/auth_service/service"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Init(app fiber.Router, pg *pgxpool.Pool, redis *redis.Client) {
	authRedisRepo := repository.NewAuthRedisRepository(redis)
	authPostgresRepo := repository.NewUsersPostgresRepository(pg)
	authService := service.NewAuthService(authPostgresRepo, authRedisRepo)
	authController := controller.NewAuthController(authService)

	app.Post("/register", authController.RegisterUser)
	app.Post("/login", authController.Login)
	app.Post("/refresh", authController.RefreshToken)
	app.Get("/protected", authController.Protected)
	app.Post("/logout", authController.Logout)
	app.Post("/revoke", authController.Revoke)

}
