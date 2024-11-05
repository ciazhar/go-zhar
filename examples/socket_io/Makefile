run-frontend:
	go run cmd/frontend/main.go

run-backend:
	go run cmd/backend/main.go

insert-blood:
	@echo "Sending POST request to /blood..."
	@curl -X POST http://localhost:8091/blood \
		-H "Content-Type: application/json" \
		-d '{ \
			"ApplicationID": "digisar", \
			"ID": 1, \
			"A": 5, \
			"B": 3, \
			"AB": 2, \
			"O": 4 \
		}'