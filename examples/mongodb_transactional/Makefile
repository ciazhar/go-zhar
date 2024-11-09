# Temporary file to store the last inserted ID
ID_FILE := .last_inserted_id

run-mongodb:
	cd ../../deployments/mongodb && docker-compose up

run:
	go run cmd/main.go

create-book:
	curl --location 'http://localhost:3000/book' \
    --header 'Content-Type: application/json' \
    --data '{\
        "title": "Sample Title",\
        "quantity": 100\
    }' | jq -r '.data.id' > $(ID_FILE) && \
    echo "Inserted ID: $$(cat $(ID_FILE))"

purchase-book:
	@if [ -f $(ID_FILE) ]; then \
		ID=$$(cat $(ID_FILE)); \
		curl --location 'http://localhost:3000/purchase' \
		--header 'Content-Type: application/json' \
		--data '{\
			"book_id": "'"$$ID"'",\
			"amount": 2\
		}'; \
	else \
		echo "ID file not found. Run 'make create-book' first."; \
	fi
