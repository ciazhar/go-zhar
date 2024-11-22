run-orchestrator:
	go run cmd/orchestrator-service/main.go

run-stub:
	go run cmd/stub-service/main.go

run-worker:
	go run cmd/worker/main.go

install:
	go install github.com/swaggo/swag/cmd/swag@latest

generate:
	swag init --parseInternal --dir cmd/orchestrator-service,internal/controller --output=api/swagger

run:
	curl --location 'localhost:8080/onboarding' \
    --header 'Content-Type: application/json' \
    --data-raw '{ \
        "first_name": "Muhammad", \
        "last_name": "Hafidz", \
        "email": "ciazhar.id@gmail.com", \
        "city": "Kendal", \
        "account_type": "gtw", \
        "currency": "IDR" \
    }'

get-onboarding:
	curl --location 'localhost:8080/onboarding/{id}'

signature:
	curl --location 'localhost:8080/onboarding/dcf9500c-5b57-4e4a-a028-7209b372868e/signature' \
    --header 'Content-Type: application/json' \
    --data '{ \
        "agreement_id":"7772b295-147f-4e6e-bffc-27e391a69cf8", \
        "signature":"yes" \
    }'