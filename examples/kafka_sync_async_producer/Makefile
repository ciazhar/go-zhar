run-kafka:
	cd ../../deployments/kafka && docker compose up

run-service:
	go run cmd/main.go

produce-sync-message:
	curl --location 'localhost:3000/sync'

produce-async-message:
	curl --location 'localhost:3000/async'