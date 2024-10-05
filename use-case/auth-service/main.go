package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var (
	redisClient *redis.Client
	dbPool      *pgxpool.Pool
	ctx         = context.Background()
	jwtSecret   = []byte(os.Getenv("JWT_SECRET")) // Change this to a strong secret
	tokenTTL    = 15 * time.Minute                // Access token TTL
)

// User structure
type User struct {
	ID       string
	Username string
	Password string // This should be a hashed password
}

// Initialize Redis and PostgreSQL
func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
	})

	// Initialize PostgreSQL
	var err error
	dbPool, err = pgxpool.Connect(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	runMigrations()
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash checks if the password is correct
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate JWT access token
func generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func generateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Longer-lived token
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWT Claims
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Validate JWT token
func validateJWT(tokenString string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(*CustomClaims), nil
}

// Register User
func registerUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Hash the password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}
	user.Password = hashedPassword

	// Store user in the database
	_, err = dbPool.Exec(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

// Login User
func login(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Retrieve user from the database
	var user User
	err := dbPool.QueryRow(ctx, "SELECT id, username, password FROM users WHERE username=$1", body.Username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Check password
	if !CheckPasswordHash(body.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate tokens
	accessToken, err := generateAccessToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate access token"})
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate refresh token"})
	}

	// Store JWT in Redis (allow multiple tokens per user)
	err = redisClient.SAdd(ctx, "user:"+user.ID, accessToken).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save token to Redis"})
	}

	// Store refresh token in Redis with expiration
	err = redisClient.Set(ctx, "refresh:"+user.ID, refreshToken, tokenTTL*24).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save refresh token"})
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Refresh Token Handler
func refreshToken(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validate refresh token
	claims, err := validateJWT(body.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	// Check if refresh token exists in Redis
	val, err := redisClient.Get(ctx, "refresh:"+claims.UserID).Result()
	if err != nil || val != body.RefreshToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Refresh token not found"})
	}

	// Generate new access token
	newAccessToken, err := generateAccessToken(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate new access token"})
	}

	// Store JWT in Redis (allow multiple tokens per user)
	err = redisClient.SAdd(ctx, "user:"+claims.UserID, newAccessToken).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save token to Redis"})
	}

	return c.JSON(fiber.Map{
		"access_token": newAccessToken,
	})
}

// Protected route example
func protected(c *fiber.Ctx) error {
	// Get token from Authorization header
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No token provided"})
	}

	// Validate token
	claims, err := validateJWT(token[len("Bearer "):])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Check if the token exists in Redis
	exists := redisClient.SIsMember(ctx, "user:"+claims.UserID, token[len("Bearer "):]).Val()
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not found"})
	}

	return c.JSON(fiber.Map{"message": "Protected data", "user_id": claims.UserID})
}

// Logout Handler
func logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No token provided"})
	}

	// Validate token
	claims, err := validateJWT(token[len("Bearer "):])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Remove the specific token from Redis
	err = redisClient.SRem(ctx, "user:"+claims.UserID, token[len("Bearer "):]).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not remove token"})
	}

	// Delete refresh token from Redis
	err = redisClient.Del(ctx, "refresh:"+claims.UserID).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not logout"})
	}

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

// Revoke Handler
func revoke(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No token provided"})
	}

	// Validate token
	claims, err := validateJWT(token[len("Bearer "):])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Remove all tokens for the user from Redis
	err = redisClient.Del(ctx, "user:"+claims.UserID).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not revoke tokens"})
	}

	// Delete all tokens for the user
	err = redisClient.Del(ctx, "refresh:"+claims.UserID).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not revoke tokens"})
	}

	return c.JSON(fiber.Map{"message": "All tokens revoked"})
}

func runMigrations() {

	// Create a new migrate instance
	m, err := migrate.New("file://db/migrations", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	// Run the migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No changes to apply in migrations.")
		} else {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	}

	log.Println("Migrations applied successfully!")
}

// Main function
func main() {
	app := fiber.New()

	app.Post("/register", registerUser)
	app.Post("/login", login)
	app.Post("/refresh", refreshToken)
	app.Get("/protected", protected)
	app.Post("/logout", logout)
	app.Post("/revoke", revoke)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
