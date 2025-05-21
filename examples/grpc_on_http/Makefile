generate-proto:
	protoc \
      --proto_path=api/proto \
      --go_out=internal/product/controller/grpc \
      --go-grpc_out=internal/product/controller/grpc \
      api/proto/*.proto
