install:
	go install github.com/swaggo/swag/cmd/swag@latest

generate:
	swag init --parseInternal --dir cmd/,internal/ --output=api/swagger

run:
	go run cmd/main.go