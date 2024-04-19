run:
	go run cmd/main.go

docker:
	cd ../../../deployments/kafka && docker compose up

sync:
	curl --location 'localhost:3000/sync'

async:
	curl --location 'localhost:3000/async'