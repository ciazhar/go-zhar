# Go parameters
GO_CMD=go
GOTEST=$(GO_CMD) test
GOFMT=$(GO_CMD) fmt
GOVET=$(GO_CMD) vet
GO_RUN=$(GO_CMD) run
BINARY_NAME=app
SRC_DIR=./

# Default target: Run the application
.PHONY: run
run:
	$(GO_RUN) cmd/analytic_report/main.go

# Build the Go application
.PHONY: build
build:
	$(GO_CMD) build -o $(BINARY_NAME) $(SRC_DIR)

# Test the Go application
.PHONY: test
test:
	$(GOTEST) ./...

# Clean up binaries
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

# Format the Go code
.PHONY: fmt
fmt:
	$(GOFMT) ./...

# Vet the Go code for errors
.PHONY: vet
vet:
	$(GOVET) ./...

# API Tests using curl, printing responses to stdout

# Test: Aggregate Total Amount for the Last Month
.PHONY: api-test-aggregate
api-test-aggregate:
	@echo "\nTest: Aggregate Total Amount for the Last Month"
	@curl -w "\nResponse time: %{time_total}s\n" -o /dev/null -s -X GET http://localhost:3000/aggregate-amount

# Test: Count Transactions Per Type in the Last Month
.PHONY: api-test-count-transactions
api-test-count-transactions:
	@echo "\nTest: Count Transactions Per Type in the Last Month"
	@curl -w "\nResponse time: %{time_total}s\n" -o /dev/null -s -X GET http://localhost:3000/count-transactions

# Test: Average Transaction Value for Each User
.PHONY: api-test-avg-transaction
api-test-avg-transaction:
	@echo "\nTest: Average Transaction Value for Each User"
	@curl -w "\nResponse time: %{time_total}s\n" -o /dev/null -s -X GET http://localhost:3000/average-transaction-value

# Test: Top 10 Users by Total Transaction Value
.PHONY: api-test-top-users
api-test-top-users:
	@echo "\nTest: Top 10 Users by Total Transaction Value"
	@curl -w "\nResponse time: %{time_total}s\n" -o /dev/null -s -X GET http://localhost:3000/top-users

# Test: Daily Transaction Volume for the Last 30 Days
.PHONY: api-test-daily-volume
api-test-daily-volume:
	@echo "\nTest: Daily Transaction Volume for the Last 30 Days"
	@curl -w "\nResponse time: %{time_total}s\n" -o /dev/null -s -X GET http://localhost:3000/daily-volume

# Test: Users with More than X Transactions in the Last Month
.PHONY: api-test-users-x-transactions
api-test-users-x-transactions:
	@echo "\nTest: Users with More than X Transactions in the Last Month"
	@curl -w "\nResponse time: %{time_total}s\n" -o /dev/null -s -X GET http://localhost:3000/users-x-transactions

# Test: Sum of Transaction Amounts Per Day for a Specific Type (e.g. Purchase)
.PHONY: api-test-purchase-sum
api-test-purchase-sum:
	@echo "\nTest: Sum of Transaction Amounts Per Day for 'Purchase' Type"
	@curl -w "\nResponse time: %{time_total}s\n" -o /dev/null -s -X GET http://localhost:3000/purchase-sum

# Test: Total Refunds Processed in the Last Month
.PHONY: api-test-total-refunds
api-test-total-refunds:
	@echo "\nTest: Total Refunds Processed in the Last Month"
	@curl -w "\nResponse time: %{time_total}s\n" -o /dev/null -s -X GET http://localhost:3000/total-refunds

# Test: Peak Hour for Transactions
.PHONY: api-test-peak-hour
api-test-peak-hour:
	@echo "\nTest: Peak Hour for Transactions"
	@curl -w "\nResponse time: %{time_total}s\n" -o /dev/null -s -X GET http://localhost:3000/peak-hour

# Test: Users Who Had Both Purchases and Refunds in the Last Month
.PHONY: api-test-both-purchases-refunds
api-test-both-purchases-refunds:
	@echo "\nTest: Users Who Had Both Purchases and Refunds in the Last Month"
	@curl -w "\nResponse time: %{time_total}s\n" -o /dev/null -s -X GET http://localhost:3000/both-purchases-refunds

# Full API test sequence
.PHONY: api-test-all
api-test-all: api-test-aggregate api-test-count-transactions api-test-avg-transaction api-test-top-users \
              api-test-daily-volume api-test-users-x-transactions api-test-purchase-sum api-test-total-refunds \
              api-test-peak-hour api-test-both-purchases-refunds

# Full process: Clean, format, vet, build, and run
.PHONY: all
all: clean fmt vet build run

# Run local server and test specific API
.PHONY: run-local
run-local:
	CLICKHOUSE_ADDR=localhost:9000 \
	CLICKHOUSE_USER=default \
	CLICKHOUSE_PASSWORD='' \
	$(GO_RUN) main.go
