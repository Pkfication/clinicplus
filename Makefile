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