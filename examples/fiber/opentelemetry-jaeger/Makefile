run-user:
	go run cmd/user/main.go

run-order:
	go run cmd/order/main.go

curl-new-user:
	curl -X POST -H "Content-Type: application/json" -d '{"username":"user1","email":"user1@example.com","age":25}' -H "Content-Type: application/json" http://localhost:3000/users

curl-get-user-by-username:
	curl -X GET http://localhost:3000/users/user1

curl-get-all-user:
	curl -X GET http://localhost:3000/users

curl-delete-user:
	curl -X DELETE http://localhost:3000/users/user1

curl-update-user:
	curl -X PUT -H "Content-Type: application/json" -d '{"username":"user1","email":"user1@example.com","age":26}' -H "Content-Type: application/json" http://localhost:3000/users

curl-new-order:
	curl -X POST -H "Content-Type: application/json" \
         -d '{"order_id":"order123","order_date":"2023-05-27","username":"user1"}' \
         http://localhost:3001/orders

curl-get-order-by-id:
	curl -X GET http://localhost:3001/orders/order123

curl-get-all-order:
	curl -X GET http://localhost:3001/orders

curl-delete-order:
	curl -X DELETE http://localhost:3001/orders/order123

curl-update-order:
	curl -X PUT -H "Content-Type: application/json" \
         -d '{"order_id":"order123","order_date":"2023-05-28","username":"user456"}' \
         http://localhost:3001/orders

run-sequence:
	make curl-new-user
	make curl-new-order
	make curl-get-all-order