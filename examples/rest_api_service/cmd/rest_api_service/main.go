package main

import (
	"context"
	"fmt"
	bootstrap2 "github.com/ciazhar/go-start-small/examples/rest_api_service/internal/bootstrap"
	"github.com/ciazhar/go-start-small/pkg/bootstrap"
	"github.com/ciazhar/go-start-small/pkg/bootstrap/server"
	"github.com/ciazhar/go-start-small/pkg/config"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/redis"
	"github.com/ciazhar/go-start-small/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"time"
)

func main() {
	ctx := context.Background()

	logger.InitLogger(logger.LogConfig{
		LogLevel:      "debug",
		LogFile:       "./logfile.log",
		MaxSize:       10,
		MaxBackups:    5,
		MaxAge:        30,
		Compress:      true,
		ConsoleOutput: true,
	})

	config.InitConfig(ctx, config.Config{
		Source: config.File,
		Type:   "json",
		File: config.FileConfig{
			FileName: "config.json",
			FilePath: "./configs",
		},
	})

	v := validator.New("id")

	// === INIT CLIENTS ===
	var clients []bootstrap.Service
	redisSvc, redisClient, err := redis.InitRedis(ctx, viper.GetString("redis.host"), viper.GetInt("redis.port"), viper.GetString("redis.password"))
	if err != nil {
		logger.LogFatal(err).Msg("failed to init redis client")
	}
	clients = append(clients, redisSvc)

	// === START ALL CLIENTS ===
	for _, svc := range clients {
		if err := svc.Start(); err != nil {
			logger.LogFatal(err).Msgf("🔥 %s failed", svc.Name())
		}
	}

	// === INIT SERVERS + WORKERS ===
	var serversAndWorkers []bootstrap.Service
	server1 := server.NewFiberServer(
		fmt.Sprintf("%s:%s", viper.GetString("application.name"), viper.GetString("application.version")),
		fmt.Sprintf(":%s", viper.GetString("application.port")),
		func(app *fiber.App) {
			bootstrap2.InitServer(app, v, redisClient)
		})
	serversAndWorkers = append(serversAndWorkers, server1)

	// === START ALL SERVERS + WORKERS ===
	for _, svc := range serversAndWorkers {
		go func(svc bootstrap.Service) {
			if err := svc.Start(); err != nil {
				logger.LogFatal(err).Msgf("🔥 %s failed", svc.Name())
			}
		}(svc)
	}

	// === GRACEFUL SHUTDOWN ===
	bootstrap.GracefulShutdown(ctx, 10*time.Second, clients, serversAndWorkers)
}
