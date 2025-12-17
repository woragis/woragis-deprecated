# Translation Worker Tests

## Overview

This directory contains test files and documentation for the translation worker.

## Test Structure

Tests are organized alongside the source code:

- `internal/config/config_test.go` - Configuration loading tests
- `internal/queue/queue_test.go` - Queue and message handling tests
- `internal/translator/translator_test.go` - Translation API client tests
- `internal/database/database_test.go` - Database operation tests
- `pkg/logger/logger_test.go` - Logger tests
- `pkg/health/health_test.go` - Health check tests
- `cmd/translation-worker/main_test.go` - Main entry point smoke tests

## Running Tests

### Local Testing

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run with coverage
make test-cov

# Check coverage threshold (70%)
make test-cov-check
```

### Docker Testing

```bash
# Build and run tests in Docker
make test-docker

# Or manually
docker build -f Dockerfile.test -t translation-worker-test .
docker run --rm translation-worker-test
```

## Test Coverage

Target coverage: **70%**

Current coverage can be checked with:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Test Categories

### Unit Tests
- Configuration loading from environment variables
- Logger creation and trace ID handling
- Health check functionality
- Queue message serialization/deserialization
- Translation entity field operations
- Language code mapping

### Integration Tests
- Full translation workflow (requires database and RabbitMQ)
- Translation API calls (requires API keys or mocks)
- Database operations (requires test database)

## Mocking

Tests use:
- HTTP test servers for translation API mocking
- Environment variable manipulation for configuration tests
- Mock RabbitMQ connections for health checks

## Notes

- Translation API tests use mock HTTP servers to avoid requiring real API keys
- Database tests verify entity methods but full repository tests require a test database
- Queue tests verify message handling but full integration requires RabbitMQ
