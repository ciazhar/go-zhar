bench-kafka: 
	go run cmd/main.go -brokers=localhost:9092 -topic=benchmark -message-size=1024 -message-count=1000000 -producers=4 -consumers=4 -duration=2m -redpanda=false

bench-redpanda:
	go run cmd/main.go -brokers=localhost:9093 -topic=benchmark -message-size=1024 -message-count=1000000 -producers=4 -consumers=4 -duration=2m -redpanda=true

run-kafka:
	cd ../../deployments/kafka && docker-compose up

run-redpanda:
	cd ../../deployments/redpanda && docker-compose up

stats:
	docker stats kafka redpanda