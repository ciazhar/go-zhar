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
	  -d @deployments/sink-connector.json
