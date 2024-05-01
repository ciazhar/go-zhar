# Define variables
HOST := localhost
PORT := 3000

# Define curl commands
GET := curl -sS
POST := curl -sS -X POST
DELETE := curl -sS -X DELETE

run:
	go run cmd/main.go

# Define targets for API endpoints
get:
	$(GET) http://$(HOST):$(PORT)/get

set:
	$(POST) -d "value=myvalue&expiration=1h" http://$(HOST):$(PORT)/set

gethash:
	$(GET) http://$(HOST):$(PORT)/gethash/myfield

sethash:
	$(POST) -d "field=myfield&value=myvalue" http://$(HOST):$(PORT)/sethash

sethashttl:
	$(POST) -d "field=myfield&value=myvalue&ttl=1h" http://$(HOST):$(PORT)/sethashttl

deletehash:
	$(DELETE) -d "field=myfield" http://$(HOST):$(PORT)/deletehash

test:
	go test ./internal/repository