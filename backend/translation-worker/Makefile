.PHONY: test test-unit test-integration test-cov install clean build docker-build docker-run

# Install dependencies
install:
	go mod download
	go mod tidy

# Build the application
build:
	go build -o bin/translation-worker ./cmd/translation-worker

# Run all tests
test:
	go test ./... -v

# Run unit tests only
test-unit:
	go test ./internal/... ./pkg/... -v

# Run tests with coverage
test-cov:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

# Run tests with coverage threshold check
test-cov-check:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//' | awk '{if ($$1 < 70) {print "Coverage below 70%: " $$1 "%"; exit 1} else {print "Coverage: " $$1 "%"}}'

# Build Docker image
docker-build:
	docker build -t translation-worker:latest .

# Run Docker container
docker-run:
	docker run --rm \
		-e ENV=development \
		-e DATABASE_URL="host=localhost port=5432 user=woragis password=woragis dbname=woragis sslmode=disable" \
		-e RABBITMQ_URL="amqp://woragis:woragis@localhost:5672/woragis" \
		-e TRANSLATION_PROVIDER=google \
		-e GOOGLE_TRANSLATE_API_KEY="your-api-key" \
		translation-worker:latest

# Run tests in Docker
test-docker:
	docker build -f Dockerfile.test -t translation-worker-test .
	docker run --rm translation-worker-test

# Clean test artifacts
clean:
	rm -f coverage.out
	rm -f coverage.html
	rm -f bin/translation-worker
	go clean -testcache
