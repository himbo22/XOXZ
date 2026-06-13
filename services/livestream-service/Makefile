# Auth Service Makefile
# Commands for database migrations using Goose

# Load environment variables from .env file
include .env
export

# Database connection string
DB_STRING ?= "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSL_MODE)"

# Goose migration directory
MIGRATIONS_DIR = ./migrations

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOTEST = $(GOCMD) test
GOMOD = $(GOCMD) mod
GOGENERATE = $(GOCMD) generate

# Binary name
BINARY_NAME = account-svc

.PHONY: all build run dev test clean migrate-up migrate-down migrate-status migrate-create wire swagger help

## help: Display this help message
help:
	@echo "Auth Service - Available Commands:"
	@echo ""
	@echo "Build & Run:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run with hot reload (requires air)"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "Database Migrations (Goose):"
	@echo "  make migrate-up     - Run all pending migrations"
	@echo "  make migrate-down   - Rollback last migration"
	@echo "  make migrate-down-all - Rollback all migrations"
	@echo "  make migrate-status - Show migration status"
	@echo "  make migrate-reset  - Rollback all and re-run migrations"
	@echo "  make migrate-create name=<name> - Create new migration"
	@echo ""
	@echo "Code Generation:"
	@echo "  make wire           - Generate Wire dependency injection"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make generate       - Run all code generators"
	@echo ""
	@echo "Dependencies:"
	@echo "  make deps           - Download dependencies"
	@echo "  make tidy           - Tidy go.mod"
	@echo "  make install-tools  - Install required tools"
	@echo ""

## all: Build and run the application
all: build

## build: Build the application binary
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) ./cmd/main.go

## run: Run the application
run:
	$(GORUN) ./cmd/main.go

## dev: Run with hot reload (requires air)
dev:
	@command -v air > /dev/null 2>&1 || { echo "Installing air..."; go install github.com/air-verse/air@latest; }
	air

d-run:
	doppler run -- $(GORUN) ./cmd/main.go

d-dev:
	doppler run --	air

seed:
	$(GORUN) ./cmd/seeder/main.go

d-seed:
	doppler run -- $(GORUN) ./cmd/seeder/main.go

## test: Run tests
test:
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	$(GOTEST) -v -cover -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## clean: Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html


# ==================== RSA KEY GENERATION =====================
rsa-key:
	openssl genrsa -out private.pem 2048 && openssl rsa -in private.pem -pubout -out public.pem
# ==================== Database Migrations ====================

## migrate-up: Run all pending migrations
migrate-up:
	@echo "Running migrations..."
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) up

## migrate-up-one: Run one migration
migrate-up-one:
	@echo "Running one migration..."
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) up-by-one

## migrate-down: Rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) down

## migrate-down-all: Rollback all migrations
migrate-down-all:
	@echo "Rolling back all migrations..."
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) reset

## migrate-status: Show migration status
migrate-status:
	@echo "Migration status:"
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) status

## migrate-reset: Rollback all and re-run migrations
migrate-reset: migrate-down-all migrate-up
	@echo "Database reset complete!"

## migrate-create: Create new migration (usage: make migrate-create name=create_table)
migrate-create:
ifndef name
	$(error name is required. Usage: make migrate-create name=migration_name)
endif
	@echo "Creating migration: $(name)"
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

## migrate-version: Show current migration version
migrate-version:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) version

## migrate-validate: Validate migrations
migrate-validate:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) validate

# ==================== Code Generation ====================

## wire: Generate Wire dependency injection code
wire:
	@echo "Generating Wire code..."
	@command -v wire > /dev/null 2>&1 || { echo "Installing wire..."; go install github.com/google/wire/cmd/wire@latest; }
	cd internal/di && wire

## swagger: Generate Swagger documentation
swagger:
	@echo "Generating Swagger docs..."
	@command -v swag > /dev/null 2>&1 || { echo "Installing swag..."; go install github.com/swaggo/swag/cmd/swag@latest; }
	swag init -g cmd/main.go -o docs

## generate: Run all code generators
generate: wire swagger
	$(GOGENERATE) ./...

# ==================== Dependencies ====================

## deps: Download dependencies
deps:
	$(GOMOD) download

## tidy: Tidy go.mod
tidy:
	$(GOMOD) tidy

## install-tools: Install required development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest
	@echo "All tools installed!"

# ==================== Docker ====================

## docker-build: Build Docker image
docker-build:
	docker build -t $(BINARY_NAME):latest .

## docker-run: Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

## docker-compose-up: Start with docker-compose
docker-compose-up:
	docker-compose up -d

## docker-compose-down: Stop docker-compose
docker-compose-down:
	docker-compose down

# ==================== Lint ====================

## lint: Run linter
lint:
	@command -v golangci-lint > /dev/null 2>&1 || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	golangci-lint run ./...

## fmt: Format code
fmt:
	$(GOCMD) fmt ./...
