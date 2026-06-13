
# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOTEST = $(GOCMD) test
GOMOD = $(GOCMD) mod
GOGENERATE = $(GOCMD) generate

# Binary name
BINARY_NAME = auth-svc

.PHONY: all build run dev test clean migrate-up migrate-down migrate-status migrate-create wire swagger help

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
