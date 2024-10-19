package internal

import (
	"github.com/ciazhar/go-start-small/examples/jwt_authorization_service/internal/controller"
	"github.com/ciazhar/go-start-small/examples/jwt_authorization_service/internal/repository"
	"github.com/ciazhar/go-start-small/examples/jwt_authorization_service/internal/service"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
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