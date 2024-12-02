run:
	go run cmd/main.go

send-order-status-notification:
	curl --location 'localhost:3000/notification/order-status' \
    --header 'Content-Type: application/json' \
    --data '{ \
        "order_id": "12345", \
        "status": "Order Shipped" \
    }'

publish-ttl:
	curl --location 'localhost:3000/notification/payment-reminder' \
    --header 'Content-Type: application/json' \
    --data '{ \
        "order_id": "12345", \
        "reminder": "Please complete your payment within 10 minutes." \
    }'