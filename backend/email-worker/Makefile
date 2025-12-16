.PHONY: test test-unit test-integration test-cov install clean

# Install dependencies
install:
	go mod download
	go mod tidy

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

# Clean test artifacts
clean:
	rm -f coverage.out
	rm -f coverage.html
	go clean -testcache
