proto:
	protoc --go_out=. user.proto

bench:
	go test -bench=.

consume:
	go run cmd/consumer/main.go

produce:
	go run cmd/producer/main.go