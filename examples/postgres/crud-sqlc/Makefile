generate:
	cd configs && sqlc generate

run:
	cd configs && sqlc generate
	go run cmd/main.go


# Variables
BASE_URL=http://localhost:3000/api/products
PRODUCT_ID=1

# Create a new product
create-product:
	@curl -X POST $(BASE_URL) \
		-H "Content-Type: application/json" \
		-d '{"name": "Product 1", "price": 9.99}'

# Get products with filtering, sorting, and pagination
get-products:
	@curl -X GET "$(BASE_URL)?name=Product&price=10&sortBy=price&page=1&size=10"

# Get products with cursor-based pagination
get-products-cursor:
	@curl -X GET "$(BASE_URL)/cursor?name=next,&price=10&cursor=10&size=10"

# Update product price
update-product-price:
	@curl -X PUT $(BASE_URL)/$(PRODUCT_ID)/price \
		-H "Content-Type: application/json" \
		-d '{"name": "Product 1", "price": 14.99}'

# Delete product
delete-product:
	@curl -X DELETE $(BASE_URL)/$(PRODUCT_ID)