# Astroneko Backend Makefile

.PHONY: help build run test test-unit test-integration test-coverage test-benchmark clean migrate migrate-status gen-mock

# Default target
help:
	@echo "Available commands:"
	@echo "  build            - Build the application"
	@echo "  run              - Run the application"
	@echo "  test             - Run all tests"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-integration - Run integration tests"
	@echo "  test-coverage    - Run tests with coverage"
	@echo "  test-benchmark   - Run benchmark tests"
	@echo "  clean            - Clean build artifacts"
	@echo "  migrate          - Run database migrations"
	@echo "  migrate-status   - Show migration status"
	@echo "  gen-mock         - Generate mock files"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/astroneko-backend cmd/api/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/api/main.go

# Run all tests
test:
	@echo "Running all tests..."
	go test -v ./...

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test -v -short ./...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmark tests
test-benchmark:
	@echo "Running benchmark tests..."
	go test -bench=. -benchmem ./...

# Generate mock files
gen-mock:
	@echo "Generating mock files..."
	@echo "Generating user repository mock..."
	mockgen -source=internal/core/ports/user/repository.go -package=mock_ports -destination=testings/mock_ports/user_repository.go
	@echo "Generating database interface mock..."
	mockgen -source=internal/core/ports/database.go -package=mock_ports -destination=testings/mock_ports/database.go
	@echo "Generating referral code repository mock..."
	mockgen -source=internal/core/ports/referral_code/repository.go -package=mock_ports -destination=testings/mock_ports/referral_code_repository.go
	@echo "Generating logger mock..."
	mockgen -source=pkg/logger/logger.go -package=mock_logger -destination=testings/mock_logger/logger.go
	@echo "Generating Firebase client mock..."
	mockgen -source=pkg/firebase/firebase.go -package=mock_firebase -destination=testings/mock_firebase/firebase.go

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run database migrations
migrate:
	@echo "Running database migrations..."
	go run cmd/migrate/main.go

# Show migration status
migrate-status:
	@echo "Migration status:"
	@echo "Applied migrations:"
	@go run -c "SELECT version, applied_at FROM schema_migrations ORDER BY applied_at;" cmd/migrate/main.go || echo "No migrations table found"

# Run test generation
gen-test: gen-mock test-unit
