# Define the Go command
GO := go
APP_NAME := jwt-auth-app

# Define the .env file
ENV_FILE := .env

# Define Redis address and port
REDIS_HOST := localhost
REDIS_PORT := 6379

# Variables
HOST := http://localhost:3000
USERNAME := admin
PASSWORD := password
TOKEN := 

# Run the application
run: check-env redis-server
	$(GO) run main.go

# Check if the .env file exists
check-env:
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "Error: $(ENV_FILE) file not found."; \
		exit 1; \
	fi

# Start Redis server
redis-server:
	@echo "Checking if Redis is running..."
	@redis-cli -h $(REDIS_HOST) -p $(REDIS_PORT) ping > /dev/null || \
		(redis-server --daemonize yes && echo "Started Redis server")

# Login and save token
login:
	@echo "Logging in..."
	@TOKEN=$$(curl -X POST -d "username=$(USERNAME)&password=$(PASSWORD)" $(HOST)/login | jq -r '.token'); \
	echo "Token: $$TOKEN"

# Access protected route
protected:
	@echo "Accessing protected route..."
	@TOKEN=$$(curl -X POST -d "username=$(USERNAME)&password=$(PASSWORD)" $(HOST)/login | jq -r '.token'); \
	curl -H "Authorization: Bearer $$TOKEN" $(HOST)/protected

# Logout
logout:
	@echo "Logging out..."
	@TOKEN=$$(curl -X POST -d "username=$(USERNAME)&password=$(PASSWORD)" $(HOST)/login | jq -r '.token'); \
	curl -X GET -H "Authorization: Bearer $$TOKEN" $(HOST)/logout

.PHONY: login protected logout
