run:
	go run cmd/main.go

create-customers:
	curl -X POST http://localhost:3000/customers \
    -H "Content-Type: application/json" \
    -d '{ \
      "name": "Jane Doe", \
      "email": "janedoe@example.com" \
    }'

create-products:
	curl -X POST http://localhost:3000/products \
	-H "Content-Type: application/json" \
	-d '{ \
	  "name": "Product A", \
	  "price": 10.00, \
	  "stock": 50 \
	}'
	curl -X POST http://localhost:3000/products \
	-H "Content-Type: application/json" \
	-d '{ \
	  "name": "Product B", \
	  "price": 20.00, \
	  "stock": 100 \
	}'

get-products:
	curl -X GET http://localhost:3000/products \
		 -H "Content-Type: application/json"

place-order:
	curl -X POST http://localhost:3000/orders \
    -H "Content-Type: application/json" \
    -d '{ \
      "customer_id": 1, \
      "items": [ \
        { \
          "product_id": 1, \
          "quantity": 2, \
          "price": 10.0 \
        }, \
        { \
          "product_id": 2, \
          "quantity": 1, \
          "price": 20.0 \
        } \
      ] \
    }'

process-payment:
	curl -X POST http://localhost:3000/orders/payment \
    -H "Content-Type: application/json" \
    -d '{ \
      "order_id": 1, \
      "method": "CreditCard", \
      "amount": 250.0 \
    }'

ship-order:
	curl -X POST http://localhost:3000/orders/ship \
    -H "Content-Type: application/json" \
    -d '{ \
      "order_id": 1, \
      "tracking_number": "TRACK123456", \
      "carrier": "FedEx" \
    }'

mark-order-as-delivered:
	curl -X PATCH http://localhost:3000/orders/1/delivered