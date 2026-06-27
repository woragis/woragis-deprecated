.PHONY: help build run dev test lint fmt migrate clean docker-build docker-up docker-down

help:
	@echo "Posts AI Service - Available commands:"
	@echo ""
	@echo "  make build          Build binary"
	@echo "  make run            Run server"
	@echo "  make dev            Run with hot reload (requires air)"
	@echo "  make test           Run tests"
	@echo "  make lint           Run linter (requires golangci-lint)"
	@echo "  make fmt            Format code"
	@echo "  make migrate        Run database migrations"
	@echo "  make clean          Clean build artifacts"
	@echo "  make docker-build   Build Docker image"
	@echo "  make docker-up      Start Docker Compose stack"
	@echo "  make docker-down    Stop Docker Compose stack"
	@echo ""

build:
	@echo "Building posts-ai-service..."
	go build -o dist/posts-ai-service cmd/server/main.go

run: build
	@echo "Running posts-ai-service..."
	./dist/posts-ai-service

dev:
	@echo "Running with hot reload (requires: go install github.com/cosmtrek/air@latest)..."
	air

test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	@echo "Running linter (requires: golangci-lint)..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
	go mod tidy

migrate:
	@echo "Running migrations..."
	# Migrations run automatically on service start
	# Use psql directly if needed:
	# psql -f db/migrations/001_init.sql $$DATABASE_URL

clean:
	@echo "Cleaning build artifacts..."
	rm -rf dist/
	rm -f coverage.out coverage.html
	rm -f *.log

docker-build:
	@echo "Building Docker image..."
	docker build -t posts-ai-service:latest .

docker-up:
	@echo "Starting Docker Compose stack..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker Compose stack..."
	docker-compose down

# Development helpers
db-create:
	createdb posts_ai || true

db-drop:
	dropdb posts_ai || true

db-reset: db-drop db-create

psql:
	psql -U woragis -d posts_ai

logs:
	@echo "Note: Add this to your Go code to see logs in JSON format"
	@echo "Logs are emitted via fiber logger middleware"
