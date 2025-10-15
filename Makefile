# Makefile for the clinicplus project

# Variables
APP_NAME = clinicplus-api
BUILD_DIR = ./cmd/server
BINARY = main
MIGRATION_DIR = ./migrations

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building the application..."
	go build -o $(BINARY) $(BUILD_DIR)

# Run the application
.PHONY: run
run: build
	@echo "Running the application..."
	./$(BINARY)

# Run database migrations
.PHONY: migrate
migrate:
	@echo "Running migrations..."
	goose up

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY)

# Start observability stack
.PHONY: observability-up
observability-up:
	@echo "Starting observability stack..."
	docker-compose -f docker-compose.observability.yml up -d

# Stop observability stack
.PHONY: observability-down
observability-down:
	@echo "Stopping observability stack..."
	docker-compose -f docker-compose.observability.yml down

# View observability logs
.PHONY: observability-logs
observability-logs:
	@echo "Viewing observability logs..."
	docker-compose -f docker-compose.observability.yml logs -f

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Run load test (Go version)
.PHONY: load-test
load-test:
	@echo "Running Go load test..."
	go run load_test.go

# Run load test (Bash version)
.PHONY: load-test-bash
load-test-bash:
	@echo "Running bash load test..."
	./load_test.sh

# Quick load test on health endpoint
.PHONY: load-test-health
load-test-health:
	@echo "Running quick load test on health endpoint..."
	@for i in $$(seq 1 100); do \
		curl -s -o /dev/null -w "Request $$i: %{http_code} %{time_total}s\n" http://localhost:8080/health & \
		if [ $$((i % 10)) -eq 0 ]; then wait; fi \
	done; wait