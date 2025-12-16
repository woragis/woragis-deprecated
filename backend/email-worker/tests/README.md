# Email Worker Tests

## Overview

This directory contains unit and integration tests for the Email Worker.

## Test Structure

```
email-worker/
├── internal/
│   ├── config/
│   │   └── config_test.go
│   ├── sender/
│   │   └── smtp_sender_test.go
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
go test ./internal/sender -v
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
- Use mocks for external services (SMTP, RabbitMQ)

### Integration Tests
- May require RabbitMQ and SMTP servers
- Test with real external services
- Can be skipped if services unavailable

## Test Categories

### Config Tests
- Email configuration loading
- RabbitMQ configuration loading
- Worker configuration loading
- Environment variable parsing
- Default values
- Validation

### SMTP Sender Tests
- Sender initialization
- Email message formatting
- MIME encoding
- Connection handling
- Error handling

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
docker build -f Dockerfile.test -t email-worker-test .

# Run tests
docker run --rm email-worker-test
```

## Mocking

For unit tests, we use:
- Manual mocks for RabbitMQ connections
- Manual mocks for SMTP connections
- Interface-based design for testability

## Example Test

```go
func TestLoadEmailConfig(t *testing.T) {
    os.Setenv("SMTP_HOST", "smtp.example.com")
    os.Setenv("SMTP_FROM", "noreply@example.com")
    defer os.Unsetenv("SMTP_HOST")
    defer os.Unsetenv("SMTP_FROM")

    cfg, err := config.LoadEmailConfig()
    if err != nil {
        t.Fatalf("LoadEmailConfig() error = %v", err)
    }
    if !cfg.Enabled() {
        t.Error("EmailConfig should be enabled")
    }
}
```
