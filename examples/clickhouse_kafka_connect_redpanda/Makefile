# Makefile for managing Kafka-Sink Connector with Docker and ClickHouse

# Define environment variables
CONNECTOR_NAME = clickhouse-sink-connector-29
CONNECT_HOST = localhost
CONNECT_PORT = 8083
KAFKA_TOPIC = event_29
CLICKHOUSE_HOST = clickhouse
CLICKHOUSE_PORT = 8123
CLICKHOUSE_DB = default
CLICKHOUSE_TABLE = event_29

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
	@curl -X POST http://$(CONNECT_HOST):$(CONNECT_PORT)/connectors \
		-H "Content-Type: application/json" \
		--data '{ \
			"name": "$(CONNECTOR_NAME)", \
			"config": { \
				"connector.class": "com.clickhouse.kafka.connect.ClickHouseSinkConnector", \
				"format": "JSONEachRow", \
				"topics": "$(KAFKA_TOPIC)", \
				"hostname": "$(CLICKHOUSE_HOST)", \
				"port": "$(CLICKHOUSE_PORT)", \
				"clickhouse.database": "$(CLICKHOUSE_DB)", \
				"clickhouse.table": "$(CLICKHOUSE_TABLE)", \
				"tasks.max": "1", \
				"transforms": "RenameFields", \
				"transforms.RenameFields.type": "org.apache.kafka.connect.transforms.ReplaceField$$Value", \
				"transforms.RenameFields.renames": "eventId:event_id,companyId:company_id,type:type,injectionTime:injection_time,sendingDomain:sending_domain,recipientDomain:recipient_domain,subAccountId:sub_account_id,ipPool:ip_pool,reason:reason,bounceClassificationCode:bounce_classification_code,timestamp:timestamp,ampEnabled:amp_enabled,clickTracking:click_tracking,deliveryMethod:delivery_method,errorCode:error_code,friendlyFrom:friendly_from,geoIp:geo_ip,ipAddress:ip_address,mailboxProvider:mailbox_provider,mailboxProviderRegion:mailbox_provider_region,envelopeId:envelope_id,messageFrom:message_from,messageSize:message_size,numRetries:num_retries,openTracking:open_tracking,recipientMeta:recipient_meta,recipientTo:recipient_to,rawReason:raw_reason,rawRecipientTo:raw_recipient_to,recipientType:recipient_type,receivedMethod:received_method,routingDomain:routing_domain,scheduledTime:scheduled_time,sendingIp:sending_ip,subject:subject,targetLinkName:target_link_name,targetLinkUrl:target_link_url,templateId:template_id,templateVersion:template_version,transactional:transactional,transmissionId:transmission_id,userAgent:user_agent,userAgentParsed:user_agent_parsed", \
				"clickhouse.insert.column.definition.injection_time": "toDateTime64(injection_time / 1000, 3)", \
				"key.converter": "org.apache.kafka.connect.json.JsonConverter", \
				"value.converter": "org.apache.kafka.connect.json.JsonConverter", \
				"key.converter.schemas.enable": false, \
				"value.converter.schemas.enable": false, \
				"json.parse.mode": "ALLOW_EMPTY_OBJECTS", \
				"json.parse.with.nulls": "true" \
			} \
		}'

# Produce messages to Kafka topic
produce-message:
	@echo "Producing test messages to Kafka topic $(KAFKA_TOPIC)..."
	docker exec -it deployments-redpanda-1 rpk topic produce $(KAFKA_TOPIC) --brokers=redpanda:9092
	#docker exec -it deployments-redpanda-1 echo '{    "eventId": "9987a96c-7b0f-4b39-b9e8-f730d746d005",    "companyId": 513,    "transmissionId": "7efef88f-ac56-11ef-82ec-ce4148be6eef",    "subAccountId": 1256,    "ipAddress": "36.50.97.91",    "sendingDomain": "bpddiy.co.id",    "campaignId": null,    "customerId": null,    "sendingIp": "36.50.97.91",    "envelopeId": "7efff227-ac56-11ef-82ec-ce4148be6eef",    "ipPool": "shared pool",    "friendlyFrom": "digdaya@bpddiy.co.id",    "messageFrom": null,    "rawRecipientTo": null,    "recipientDomain": "gmail.com",    "recipientType": "TO",    "recipientMeta": {},    "recipientTags": null,    "recipientTo": "kalitirtoapprover@gmail.com",    "recipientHash": null,    "subject": "BUKTI TRANSAKSI - DIGDAYA BANK BPD DIY",    "messageSize": null,    "scheduledTime": null,    "injectionTime": 1732667296964,    "injectionTimeLong": 1732674633416,    "timestamp": 1732674633416,    "deliveryMethod": "SMTP",    "openTracking": false,    "clickTracking": false,    "initialPixel": null,    "mailboxProvider": "gmail-smtp-in.l.google.com.",    "mailboxProviderRegion": null,    "remoteAddress": "142.251.10.26",    "routingDomain": null,    "errorCode": 4,    "rawReason": "4.2.1 The user you are trying to contact is receiving mail at a rate that\n4.2.1 prevents additional messages from being delivered. Please resend your\n4.2.1 message at a later time. If the user is able to receive mail at that\n4.2.1 time, your message will be delivered. For more information, go to\n4.2.1  https://support.google.com/mail/?p=ReceivingRate 41be03b00d2f7-7fc84488e3asi1858835a12.815 - gsmtp",    "reason": "4.2.1 The user you are trying to contact is receiving mail at a rate that\n4.2.1 prevents additional messages from being delivered. Please resend your\n4.2.1 message at a later time. If the user is able to receive mail at that\n4.2.1 time, your message will be delivered. For more information, go to\n4.2.1  https://support.google.com/mail/?p=ReceivingRate 41be03b00d2f7-7fc84488e3asi1858835a12.815 - gsmtp",    "reportBy": null,    "reportTo": null,    "bounceClassificationCode": 70,    "receivedMethod": null,    "ampEnabled": null,    "numRetries": 4,    "outboundTls": null,    "queueTime": null,    "transactional": false,    "type": "delay"}' | rpk topic produce $(KAFKA_TOPIC) --brokers=redpanda:9092

# Check status of ClickHouse Sink Connector
check-connector:
	@echo "Checking ClickHouse Sink Connector status..."
	curl -X GET http://$(CONNECT_HOST):$(CONNECT_PORT)/connectors/$(CONNECTOR_NAME)/status

# Restart ClickHouse Sink Connector
restart-connector:
	@echo "Restarting ClickHouse Sink Connector..."
	curl -X POST http://$(CONNECT_HOST):$(CONNECT_PORT)/connectors/$(CONNECTOR_NAME)/restart

# Delete ClickHouse Sink Connector
restart-connector:
	@echo "Restarting ClickHouse Sink Connector..."
	curl -X DELETE http://$(CONNECT_HOST):$(CONNECT_PORT)/connectors/$(CONNECTOR_NAME)

# Clean up the deployment environment (stop and remove containers)
cleanup:
	@echo "Cleaning up deployment environment..."
	cd deployments && docker-compose down
