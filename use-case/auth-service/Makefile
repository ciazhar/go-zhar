# Variables
APP_NAME=myapp
PORT=3000
REDIS_HOST=localhost:6379
POSTGRES_URL=postgres://user:password@localhost:5432/postgres
JWT_TOKEN=

# Run the application
run:
	@echo "Running the server..."
	REDIS_HOST=$(REDIS_HOST) POSTGRES_URL=$(POSTGRES_URL) go run main.go

# Build the application
build:
	@echo "Building the Go application..."
	@go build -o $(APP_NAME) main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(APP_NAME)

register:
	curl -X POST http://localhost:3000/register -H "Content-Type: application/json" -d '{"username":"testuser","password":"testpassword"}'

# Hit the login API (adjust username and password as necessary)
login:
	@echo "Logging in and retrieving JWT token..."
	@curl -X POST http://localhost:$(PORT)/login \
		-H "Content-Type: application/json" \
		-d '{"username":"testuser", "password":"testpassword"}' \
		| tee response.json | jq -r '.access_token' > jwt_token.txt
	@echo "Login complete. Token saved to jwt_token.txt"

# Hit the protected API (requires valid JWT token from login)
protected:
	@echo "Hitting protected route..."
	@JWT_TOKEN=$$(cat jwt_token.txt); \
	curl -X GET http://localhost:$(PORT)/protected \
		-H "Authorization: Bearer $$JWT_TOKEN"

# Logout by removing specific token
logout:
	@echo "Logging out..."
	@JWT_TOKEN=$$(cat jwt_token.txt); \
	curl -X POST http://localhost:$(PORT)/logout \
		-H "Authorization: Bearer $$JWT_TOKEN"
	@echo "Logout complete."

# Revoke all tokens for the user
revoke:
	@echo "Revoking all tokens for the user..."
	@JWT_TOKEN=$$(cat jwt_token.txt); \
	curl -X POST http://localhost:$(PORT)/revoke \
		-H "Authorization: Bearer $$JWT_TOKEN"
	@echo "All tokens revoked."

# Run the application and hit the APIs (for testing)
test: run login protected logout revoke

# Make help - displays available commands
help:
	@echo "Available commands:"
	@echo "  make run        - Run the server"
	@echo "  make build      - Build the Go app"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make login      - Test the login API"
	@echo "  make protected  - Test the protected route"
	@echo "  make logout     - Test the logout API"
	@echo "  make revoke     - Revoke all tokens for the user"
	@echo "  make test       - Run the app and test APIs"
