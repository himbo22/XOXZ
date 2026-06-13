# Artist Service Makefile

include .env
export

DB_STRING ?= "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSL_MODE)"
MIGRATIONS_DIR = ./migrations

GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOTEST = $(GOCMD) test
GOMOD = $(GOCMD) mod

BINARY_NAME = artist-svc

.PHONY: build run test clean migrate-up migrate-down migrate-status migrate-create wire swagger deps tidy fmt

build:
	$(GOBUILD) -o bin/$(BINARY_NAME) ./cmd/main.go

run:
	$(GORUN) ./cmd/main.go

test:
	$(GOTEST) -v ./...

clean:
	rm -rf bin/

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) status

migrate-create:
ifndef name
	$(error name is required. Usage: make migrate-create name=migration_name)
endif
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

wire:
	cd internal/di && wire

swagger:
	swag init -g cmd/main.go -o docs

deps:
	$(GOMOD) download

tidy:
	$(GOMOD) tidy

fmt:
	$(GOCMD) fmt ./...
