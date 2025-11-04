.PHONY: test wire migrate-up migrate-down migrate-create migrate-goto migrate-version docker-up docker-down docker-restart docker-logs docker-build up down

test:
	@go test $(shell go list ./... | grep -v /tmp) -cover

wire:
	@go mod download && go mod tidy && wire ./cmd/restful

# Database migration targets
# Override with environment variable: make migrate-up DB_URL="postgres://..."
DB_URL ?= $(or $(POSTGRES_URL),postgres://postgres:postgres@postgresql:5432/postgres?sslmode=disable&x-multi-statement=true)
MIGRATION_PATH := /migrations

migrate-up:
	MSYS_NO_PATHCONV=1 docker compose run --rm postgres -path=$(MIGRATION_PATH) -database="$(DB_URL)" up

migrate-down:
	MSYS_NO_PATHCONV=1 docker compose run --rm postgres -path=$(MIGRATION_PATH) -database="$(DB_URL)" down

migrate-create:
	@read -p "Migration name: " name; \
	MSYS_NO_PATHCONV=1 docker compose run --rm -v $${PWD}/deployments/migrations/postgresql:$(MIGRATION_PATH) postgres \
		create -ext sql -dir $(MIGRATION_PATH) -seq $$name

migrate-goto:
	@read -p "Target version: " version; \
	MSYS_NO_PATHCONV=1 docker compose run --rm postgres -path=$(MIGRATION_PATH) -database="$(DB_URL)" goto $$version

migrate-version:
	MSYS_NO_PATHCONV=1 docker compose run --rm postgres -path=$(MIGRATION_PATH) -database="$(DB_URL)" version


# Docker targets
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-restart:
	docker compose restart restful

docker-logs:
	docker compose logs -f restful

docker-build:
	docker compose build --no-cache

up:
	@echo "Starting PostgreSQL..."
	docker compose up -d postgresql
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 10
	@echo "Running migrations..."
	$(MAKE) migrate-up
	@echo "Building and starting API..."
	docker compose up -d --build restful
	@echo "Startup complete!"
	@echo "API: http://localhost:8080"
	docker compose ps

down:
	@echo "Stopping all services..."
	docker compose down
	@echo "All services stopped."

# Development with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	@echo "Starting PostgreSQL..."
	@docker compose up -d postgresql
	@sleep 5
	@echo "Running migrations..."
	@$(MAKE) migrate-up
	@echo "Starting API with hot reload..."
	@echo "Install air if not exists: go install github.com/air-verse/air@latest"
	air || go run ./cmd/restful



