# WhatsApp Worker Tests

## Overview

This directory contains unit and integration tests for the WhatsApp Worker.

## Test Structure

```
whatsapp-worker/
├── internal/
│   ├── config/
│   │   └── config_test.go
│   └── queue/
│       └── queue_test.go
├── pkg/
│   └── health/
│       └── health_test.go
└── tests/
    └── README.md
```

## Running Tests

### All Tests
```bash
go test ./... -v
```

### Unit Tests Only
```bash
go test ./internal/... ./pkg/... -v
```

### With Coverage
```bash
go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Specific Package
```bash
go test ./internal/config -v
go test ./internal/queue -v
go test ./pkg/health -v
```

## Test Coverage Goals

- **Unit Tests:** 70%+ coverage
- **Integration Tests:** Critical paths covered
- **Overall:** 70%+ coverage

## Test Requirements

### Unit Tests
- No external dependencies
- Fast execution (< 1 second each)
- Use mocks for external services (RabbitMQ, WhatsApp)

### Integration Tests
- May require RabbitMQ and WhatsApp services
- Test with real external services
- Can be skipped if services unavailable

## Test Categories

### Config Tests
- RabbitMQ configuration loading
- Worker configuration loading
- WhatsApp configuration loading
- Environment variable parsing
- Default values
- Validation

### Queue Tests
- Queue connection
- Message consumption
- Message acknowledgment
- Message rejection/requeue
- Error handling

### Health Check Tests
- Health check endpoint
- RabbitMQ connection status
- Response formatting
- Caching behavior

## Continuous Integration

Tests are run in CI/CD pipeline:
- All unit tests on every push
- Integration tests on merge to main
- Coverage reports generated
- Fail build if coverage < 70%

## Running Tests in Docker

```bash
# Build test image
docker build -f Dockerfile.test -t whatsapp-worker-test .

# Run tests
docker run --rm whatsapp-worker-test
```

## Mocking

For unit tests, we use:
- Manual mocks for RabbitMQ connections
- Manual mocks for WhatsApp connections
- Interface-based design for testability

## Example Test

```go
func TestLoadWorkerConfig(t *testing.T) {
    os.Setenv("WHATSAPP_QUEUE_NAME", "test.queue")
    defer os.Unsetenv("WHATSAPP_QUEUE_NAME")

    cfg := config.LoadWorkerConfig()
    if cfg.QueueName != "test.queue" {
        t.Errorf("QueueName = %v, want test.queue", cfg.QueueName)
    }
}
```
