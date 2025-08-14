CONNECT_URL ?= http://localhost:8083/connectors

.PHONY: all source sink

all: source sink

source:
	curl -X POST $(CONNECT_URL) \
	  -H "Content-Type: application/json" \
	  -d @deployments/source-connector.json

sink:
	curl -X POST $(CONNECT_URL) \
	  -H "Content-Type: application/json" \
	  -d @deployments/sink-clickhouse-connector.json

source-status:
	curl $(CONNECT_URL)/source-postgres-connector/status | jq

sink-status:
	curl $(CONNECT_URL)/clickhouse-sink-connector/status | jq

delete-source:
	curl -X DELETE $(CONNECT_URL)/source-postgres-connector

delete-sink:
	curl -X DELETE $(CONNECT_URL)/clickhouse-sink-connector

rebuild-docker:
	docker compose --file deployments/docker-compose.yml down && docker compose --file deployments/docker-compose.yml up -d --build
