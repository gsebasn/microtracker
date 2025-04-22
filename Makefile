.PHONY: run test clean build docker-build docker-run lint swagger

# Development
run:
	@echo "Starting development server..."
	@GIN_MODE=debug go run main.go

run-watch:
	@echo "Starting development server with file watcher..."
	@air

# Testing
test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

test-bench:
	@echo "Running benchmarks..."
	@go test -bench=. ./...

# Build
build:
	@echo "Building application..."
	@go build -o microtracker main.go

build-linux:
	@echo "Building Linux binary..."
	@GOOS=linux GOARCH=amd64 go build -o microtracker-linux main.go

build-all:
	@echo "Building for all platforms..."
	@GOOS=linux GOARCH=amd64 go build -o bin/linux/microtracker main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/darwin/microtracker main.go
	@GOOS=windows GOARCH=amd64 go build -o bin/windows/microtracker.exe main.go

# Clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f microtracker
	@rm -f microtracker-linux
	@rm -f coverage.out
	@rm -f coverage.txt
	@rm -f coverage.html
	@rm -rf bin/

# Docker
docker-build:
	@echo "Building Docker image..."
	@docker build -t microtracker .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env.development microtracker

docker-compose-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

docker-compose-down:
	@echo "Stopping services with Docker Compose..."
	@docker-compose down

# Database
db-start:
	@echo "Starting MongoDB container..."
	@docker run -d -p 27017:27017 --name mongodb mongo:latest

db-stop:
	@echo "Stopping MongoDB container..."
	@docker stop mongodb
	@docker rm mongodb

db-shell:
	@echo "Connecting to MongoDB shell..."
	@docker exec -it mongodb mongosh

# Environment
env-dev:
	@echo "Setting up development environment..."
	@cp .env.example .env.development

env-prod:
	@echo "Setting up production environment..."
	@cp .env.example .env.production

# Documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g main.go -o docs

# Linting
lint:
	@echo "Running linter..."
	@golangci-lint run

lint-fix:
	@echo "Running linter with auto-fix..."
	@golangci-lint run --fix

# Dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# Help
help:
	@echo "Available commands:"
	@echo "  make run              - Start development server"
	@echo "  make run-watch        - Start development server with file watcher"
	@echo "  make test             - Run tests"
	@echo "  make test-coverage    - Run tests with coverage report"
	@echo "  make test-bench       - Run benchmarks"
	@echo "  make build            - Build application"
	@echo "  make build-linux      - Build Linux binary"
	@echo "  make build-all        - Build for all platforms"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-run       - Run Docker container"
	@echo "  make docker-compose-up - Start services with Docker Compose"
	@echo "  make docker-compose-down - Stop services with Docker Compose"
	@echo "  make db-start         - Start MongoDB container"
	@echo "  make db-stop          - Stop MongoDB container"
	@echo "  make db-shell         - Connect to MongoDB shell"
	@echo "  make env-dev          - Setup development environment"
	@echo "  make env-prod         - Setup production environment"
	@echo "  make swagger          - Generate Swagger documentation"
	@echo "  make lint             - Run linter"
	@echo "  make lint-fix         - Run linter with auto-fix"
	@echo "  make deps             - Install dependencies"
	@echo "  make deps-update      - Update dependencies"
	@echo "  make help             - Show this help message" 