package main

import (
	"context"
	"flag"

	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional/internal"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional/db/migrations"
	"github.com/ciazhar/go-start-small/pkg/config"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/postgres"
	"github.com/spf13/viper"
)

func main() {

	// Configuration using flags for source, type, and other details
	var logLevel string
	var consoleOutput bool
	var source, configType, fileName, filePath, consulEndpoint, consulPath string

	// Parse command-line flags
	flag.StringVar(&logLevel, "log-level", "debug", "Log level (default: debug)")
	flag.BoolVar(&consoleOutput, "console-output", true, "Console output (default: true)")
	flag.StringVar(&source, "source", "file", "Configuration source (file or consul)")
	flag.StringVar(&fileName, "file-name", "config.json", "Name of the configuration file")
	flag.StringVar(&filePath, "file-path", "configs", "Path to the configuration file")
	flag.StringVar(&configType, "config-type", "json", "Configuration file type")
	flag.StringVar(&consulEndpoint, "consul-endpoint", "localhost:8500", "Consul endpoint")
	flag.StringVar(&consulPath, "consul-path", "path/to/config", "Path to the configuration in Consul")
	flag.Parse()

	// Initialize logger with parsed configuration
	logger.InitLogger(logger.LogConfig{
		LogLevel:      logLevel,
		ConsoleOutput: consoleOutput,
	})

	// Configuration using flags for source, type, and other details
	fileConfig := config.Config{
		Source: source,
		Type:   configType,
		File: config.FileConfig{
			FileName: fileName,
			FilePath: filePath,
		},
		Consul: config.ConsulConfig{
			Endpoint: consulEndpoint,
			Path:     consulPath,
		},
	}

	config.InitConfig(fileConfig)

	// Initialize connection pool
	pool := postgres.InitPostgres(
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.dbname"),
		viper.GetString("postgres.username"),
		viper.GetString("postgres.password"),
		logLevel,
	)
	defer pool.Close()

	postgres.InitDBMigration(
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.dbname"),
		viper.GetString("postgres.username"),
		viper.GetString("postgres.password"),
		migrations.MigrationsFS,
	)

	// // Example query
	// var greeting string
	// err := pool.QueryRow(context.Background(), "SELECT 'Hello, pgxpool with zerolog!'").Scan(&greeting)
	// if err != nil {
	// 	logger.LogFatal(context.Background(), err, "QueryRow failed", nil)
	// }

	// // Example with query arguments
	// var name string
	// var age int
	// err = pool.QueryRow(context.Background(), "SELECT $1::text, $2::int", "Alice", 30).Scan(&name, &age)
	// if err != nil {
	// 	logger.LogFatal(context.Background(), err, "QueryRow with arguments failed", nil)
	// }

	// Create repository
	userRepo := internal.NewPgxUserRepository(pool)

	// Example: Create a user
	newUser := &internal.User{Name: "John Doe", Email: "johndoe@example.com"}
	id, err := userRepo.Create(context.Background(), newUser)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to create user", nil)
	}
	logger.LogInfo(context.Background(), "User created successfully", map[string]interface{}{"id": id})

	// Example: Fetch the user
	user, err := userRepo.GetByID(context.Background(), id)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to fetch user", map[string]interface{}{"id": id})
	}
	logger.LogInfo(context.Background(), "User fetched successfully", map[string]interface{}{"user": user})

	// Example: Batch create users with transaction
	users := []internal.User{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
	}
	err = userRepo.BatchCreateWithTrx(context.Background(), users)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to batch create users", nil)
	}
	logger.LogInfo(context.Background(), "Users created successfully", map[string]interface{}{"users": users})

	// Example: Fetch all users
	users, err = userRepo.GetAll(context.Background())
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to fetch users", nil)
	}
	logger.LogInfo(context.Background(), "Users fetched successfully", map[string]interface{}{"users": users})

	// Example: Update user
	user.Name = "Jane Doe"
	err = userRepo.Update(context.Background(), user)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to update user", map[string]interface{}{"user": user})
	}
	logger.LogInfo(context.Background(), "User updated successfully", map[string]interface{}{"user": user})

	// Example: Delete user
	err = userRepo.Delete(context.Background(), id)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to delete user", map[string]interface{}{"id": id})
	}
	logger.LogInfo(context.Background(), "User deleted successfully", map[string]interface{}{"id": id})
}
