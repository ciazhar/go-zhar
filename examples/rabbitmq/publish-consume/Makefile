run:
	go run cmd/main.go

publish:
	curl --location 'localhost:3000/basic' \
    --form 'message="ahoy kapten"'

publish-ttl:
	curl --location 'localhost:3000/basic/ttl' \
    --form 'message="ahoy kapten"'