run:
	go run cmd/main.go

put:
	curl --location 'localhost:3000/basic/hash' \
    --form 'key="fizz"' \
    --form 'value="buzz"'

get:
	curl 'localhost:3000/basic/hash/fizz'