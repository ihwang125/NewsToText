# Variables
BINARY_NAME=news-to-text
MAIN_PATH=cmd/server/main.go

# Build the application
.PHONY: build
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
.PHONY: run
run:
	go run $(MAIN_PATH)

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Generate swagger docs
.PHONY: swagger
swagger:
	swag init -g $(MAIN_PATH) -o docs/

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Docker commands
.PHONY: docker-build
docker-build:
	docker build -t $(BINARY_NAME) .

.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: docker-logs
docker-logs:
	docker-compose logs -f

# Database migration
.PHONY: migrate-up
migrate-up:
	mysql -h localhost -u root -p < migrations/001_initial_schema.sql

# Development setup
.PHONY: dev-setup
dev-setup: deps
	@echo "Setting up development environment..."
	@echo "Please ensure MySQL and Redis are running"
	@echo "Copy .env.example to .env and configure your environment variables"

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  swagger       - Generate swagger documentation"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-up     - Start Docker containers"
	@echo "  docker-down   - Stop Docker containers"
	@echo "  docker-logs   - View Docker logs"
	@echo "  migrate-up    - Run database migrations"
	@echo "  dev-setup     - Setup development environment"
	@echo "  help          - Show this help message"