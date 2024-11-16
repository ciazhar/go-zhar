# Makefile for managing Kafka-Sink Connector with Docker and ClickHouse

# Define environment variables
CONNECTOR_NAME = clickhouse-sink-connector
CONNECT_HOST = localhost
CONNECT_PORT = 8083
KAFKA_TOPIC = your_kafka_topic
CLICKHOUSE_HOST = clickhouse
CLICKHOUSE_PORT = 8123
CLICKHOUSE_DB = default
CLICKHOUSE_TABLE = your_kafka_topic

# Run Docker Compose environment
run:
	@echo "Starting Docker Compose environment..."
	cd deployments && docker-compose up -d

# Create Kafka topic using Redpanda
create-topic:
	@echo "Creating Kafka topic $(KAFKA_TOPIC)..."
	docker exec -it deployments-redpanda-1 rpk topic create $(KAFKA_TOPIC) --brokers=redpanda:9092

# Create the ClickHouse Sink Connector
create-connector:
	@echo "Creating ClickHouse Sink Connector..."
	curl -X POST http://$(CONNECT_HOST):$(CONNECT_PORT)/connectors \
		-H "Content-Type: application/json" \
		--data '{"name":"$(CONNECTOR_NAME)", \
			"config": { \
				"connector.class": "com.clickhouse.kafka.connect.ClickHouseSinkConnector", \
				"topics": "$(KAFKA_TOPIC)", \
				"hostname": "$(KAFKA_TOPIC)", \
				"port": "$(CLICKHOUSE_PORT)", \
				"clickhouse.database": "$(CLICKHOUSE_DB)", \
				"clickhouse.table": "$(CLICKHOUSE_TABLE)", \
				"tasks.max": "1", \
				"key.converter": "org.apache.kafka.connect.storage.StringConverter", \
				"value.converter": "org.apache.kafka.connect.json.JsonConverter", \
				"value.converter.schemas.enable": false \
			} \
		}'

# Produce messages to Kafka topic
produce-message:
	@echo "Producing test messages to Kafka topic $(KAFKA_TOPIC)..."
	docker exec -it deployments-redpanda-1 rpk topic produce $(KAFKA_TOPIC) --brokers=redpanda:9092

# Check status of ClickHouse Sink Connector
check-connector:
	@echo "Checking ClickHouse Sink Connector status..."
	curl -X GET http://$(CONNECT_HOST):$(CONNECT_PORT)/connectors/$(CONNECTOR_NAME)

# Restart ClickHouse Sink Connector
restart-connector:
	@echo "Restarting ClickHouse Sink Connector..."
	curl -X POST http://$(CONNECT_HOST):$(CONNECT_PORT)/connectors/$(CONNECTOR_NAME)/restart

# Clean up the deployment environment (stop and remove containers)
cleanup:
	@echo "Cleaning up deployment environment..."
	cd deployments && docker-compose down
