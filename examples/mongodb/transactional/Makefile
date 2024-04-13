run:
	go run cmd/main.go

create-book:
	curl --location 'http://localhost:3000/book' \
    --header 'Content-Type: application/json' \
    --data '{\
        "title": "Sample Title",\
        "quantity": 100\
    }'

purchase-book:
	curl --location 'http://localhost:3000/purchase' \
	--header 'Content-Type: application/json' \
	--data '{\
		"book_id": "6617b2e9db70728460a5d3e4",\
		"amount": "2"\
	}'

docker:
	cd ../../../deployments/mongodb && docker-compose up