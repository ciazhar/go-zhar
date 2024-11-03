run-dashboard:
	go run cmd/dashboard_service/main.go

run-order:
	go run cmd/order_service/main.go

run-user:
	go run cmd/user_service/main.go

run-product:
	go run cmd/product_service/main.go

run:
	curl -w "\nResponse time: %{time_total}s\n" -s -X GET http://localhost:3000/dashboard?userID=1