# Variables
APP_NAME=myapp
PORT=3000
USERNAME=testuser
PASSWORD=testpassword

install:
	@echo "Installing dependencies..."
	go install github.com/vektra/mockery/v2@latest

generate:
	@echo "Generating mocks..."
	mockery  --output ./internal/mocks --dir ./internal --all

# Run the application
run:
	@echo "Running the server..."
	go run ./cmd/main.go --source=file

# Build the application
build:
	@echo "Building the Go application..."
	@go build -o $(APP_NAME) main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(APP_NAME)

# Register new user
register:
	@echo "Registering new user..."
	curl -X POST http://localhost:$(PORT)/api/v1/register \
		-H "Content-Type: application/json" \
		-d '{"username":"$(USERNAME)","password":"$(PASSWORD)"}'

# Hit the login API (adjust username and password as necessary)
login:
	@echo "Logging in and retrieving JWT tokens..."
	@curl -X POST http://localhost:$(PORT)/api/v1/login \
		-H "Content-Type: application/json" \
		-d '{"username":"$(USERNAME)", "password":"$(PASSWORD)"}' \
		| tee response.json | jq -r '.data.access_token' > access_token.txt
	@cat response.json | jq -r '.data.refresh_token' > refresh_token.txt
	@echo "Login complete. Access token saved to access_token.txt and refresh token saved to refresh_token.txt"

# Hit the refresh token API
refresh_token:
	@echo "Refreshing token..."
	@OLD_REFRESH_TOKEN=$$(cat refresh_token.txt); \
	curl -X POST http://localhost:$(PORT)/api/v1/refresh \
		-H "Authorization: Bearer $$OLD_REFRESH_TOKEN" \
		-H "Content-Type: application/json" \
		-d '{"refresh_token":"$$OLD_REFRESH_TOKEN"}' \
		| tee response.json | jq -r '.data.access_token' > access_token.txt
	@echo "New JWT token saved to access_token.txt"

# Hit the protected API (requires valid JWT token from login)
protected:
	@echo "Hitting protected route..."
	@JWT_TOKEN=$$(cat access_token.txt); \
	curl -X GET http://localhost:$(PORT)/api/v1/protected \
		-H "Authorization: Bearer $$JWT_TOKEN"

# Logout by removing specific token
logout:
	@echo "Logging out..."
	@JWT_TOKEN=$$(cat access_token.txt); \
	curl -X POST http://localhost:$(PORT)/api/v1/logout \
		-H "Authorization: Bearer $$JWT_TOKEN"
	@echo "Logout complete."

# Revoke all tokens for the user
revoke:
	@echo "Revoking all tokens for the user..."
	@JWT_TOKEN=$$(cat access_token.txt); \
	curl -X POST http://localhost:$(PORT)/api/v1/revoke \
		-H "Authorization: Bearer $$JWT_TOKEN"
	@echo "All tokens revoked."

# Run the application and hit the APIs (for testing)
test: run register login protected refresh_token logout revoke

# Make help - displays available commands
help:
	@echo "Available commands:"
	@echo "  make run            - Run the server"
	@echo "  make build          - Build the Go app"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make register       - Register a new user"
	@echo "  make login          - Test the login API"
	@echo "  make refresh_token  - Test the refresh token API"
	@echo "  make protected      - Test the protected route"
	@echo "  make logout         - Test the logout API"
	@echo "  make revoke         - Revoke all tokens for the user"
	@echo "  make test           - Run the app and test APIs"
