# =============================================================================
# Makefile for Contact Form API
# =============================================================================

.PHONY: help build run test clean docker-build docker-run docker-stop docker-logs

# Variables
APP_NAME := contact-form-api
VERSION := 1.0.0
DOCKER_IMAGE := $(APP_NAME):$(VERSION)

# Default target
help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application locally"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run with Docker Compose"
	@echo "  make docker-stop  - Stop Docker containers"
	@echo "  make docker-logs  - View Docker logs"
	@echo "  make lint         - Run linter"

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -ldflags="-w -s -X main.Version=$(VERSION)" -o bin/$(APP_NAME) ./cmd/api
	@echo "Build complete: bin/$(APP_NAME)"

# Run the application locally
run:
	@echo "Running $(APP_NAME)..."
	@go run ./cmd/api

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f *.exe
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) --build-arg VERSION=$(VERSION) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

# Run with Docker Compose
docker-run:
	@echo "Starting containers..."
	@docker-compose up -d
	@echo "Containers started. View logs with: make docker-logs"

# Stop Docker containers
docker-stop:
	@echo "Stopping containers..."
	@docker-compose down
	@echo "Containers stopped"

# View Docker logs
docker-logs:
	@docker-compose logs -f

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Generate mocks (for testing)
mocks:
	@echo "Generating mocks..."
	@mockgen -source=internal/domain/repository/email_repository.go -destination=internal/domain/repository/mocks/email_repository_mock.go
	@mockgen -source=internal/usecase/contact/interface.go -destination=internal/usecase/contact/mocks/contact_usecase_mock.go
