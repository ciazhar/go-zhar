run-redis:
	cd ../../deployments/redis && docker-compose up

run-go:
	go run cmd/main.go

bitmap-daily-login-mark-login:
	curl --location --request POST 'http://localhost:3000/mark-login?userID=12345&date=2024-12-28'

bitmap-daily-login-check-login:
	curl --location 'http://localhost:3000/check-login?userID=12345&date=2024-12-28'

geo-add-agent-location:
	curl --location 'http://localhost:3000/nearby-agents?lat=40.7128&long=-74.0060'

geo-get-nearby-agents:
	curl --location 'http://localhost:3000/nearby-agents?lat=40.7128&long=-74.0060'

hash-create-profile:
	curl --location 'http://localhost:3000/user/profile?userID=123' \
    --header 'Content-Type: application/json' \
    --data-raw '{ \
      "name": "John Doe", \
      "email": "john.doe@example.com", \
      "age": 30 \
    }'

hash-get-profile-update:
	curl --location 'http://localhost:3000/user/profile/123'

list-add-order-to-queue:
	curl --location 'http://localhost:3000/orders/queue' \
    --header 'Content-Type: application/json' \
    --data '{ \
      "order_id": "order123" \
    }'

list-process-next-order:
	curl --location --request POST 'http://localhost:3000/orders/next'

set-add-active-agent:
	curl --location --request POST 'http://localhost:3000/active/agents/1'

set-is-agent-active:
	curl --location 'http://localhost:3000/active/agents/1'

sorted-set-add-delivery-metric:
	curl --location 'http://localhost:3000/add-delivery-metric' \
    --header 'Content-Type: application/json' \
    --data '{ \
      "order_id": "order123", \
      "time_taken": 12.5 \
    }'

sorted-set-get-top-deliveries:
	curl --location 'http://localhost:3000/top-deliveries'

string-set-user-session:
	curl --location 'http://localhost:3000/user/user123/session' \
    --header 'Content-Type: application/json' \
    --data '{ \
      "user_id": "user123", \
      "session_token": "sessionToken123" \
    }'

string-get-user-session:
	curl -X GET http://localhost:3000/user/user123/session