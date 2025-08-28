run:
	go run cmd/scalable_system/main.go

run-all:
	cd deployments && \
    	docker compose down && \
    	docker rmi scalable_backend || true && \
    	docker compose build --no-cache && \
    	docker compose up -d

wire:
	cd internal/bootstrap && wire