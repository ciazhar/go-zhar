install:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install github.com/google/wire/cmd/wire@latest


lint:
	golangci-lint run --output.text.print-issued-lines=false --fix --timeout=3m

.PHONY: coverage
coverage:
	go list ./... | grep -vE '\b(config)\b|route|log|bootstrap|asynq|infrastructure|requester|mock|mocks|healthcheck|route|grpc|middleware' | xargs go test -gcflags=all=-l -coverprofile=coverage.out && go tool cover -html=coverage.out

.PHONY: clean_mock
clean_mock:
	find . -type d -name "mock" -exec rm -rf {} +

.PHONY: mock_only
mock_only: ; $(info $(M) generating mock...) @
	@./scripts/mockgen.sh

.PHONY: mock
mock: clean_mock mock_only