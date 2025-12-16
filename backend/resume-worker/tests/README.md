# Resume Worker Tests

## Overview

This directory contains unit and integration tests for the Resume Worker.

## Test Structure

```
tests/
├── unit/              # Unit tests (fast, isolated)
│   ├── test_keyword_extractor.py
│   ├── test_ai_service.py
│   └── test_translation_helper.py
├── integration/       # Integration tests (require external services)
│   └── test_worker.py
└── conftest.py        # Shared fixtures and configuration
```

## Running Tests

### All Tests
```bash
pytest
```

### Unit Tests Only
```bash
pytest tests/unit/ -v
```

### Integration Tests Only
```bash
pytest tests/integration/ -m integration
```

### With Coverage
```bash
pytest --cov=src --cov-report=html
```

### Specific Test File
```bash
pytest tests/unit/test_keyword_extractor.py -v
```

## Test Markers

- `@pytest.mark.unit` - Unit tests (fast, isolated)
- `@pytest.mark.integration` - Integration tests (require external services)
- `@pytest.mark.requires_api_key` - Tests that require API keys
- `@pytest.mark.requires_database` - Tests that require database
- `@pytest.mark.requires_rabbitmq` - Tests that require RabbitMQ
- `@pytest.mark.slow` - Slow running tests

## Running Marked Tests

```bash
# Only unit tests
pytest -m unit

# Only integration tests
pytest -m integration

# Skip tests requiring database
pytest -m "not requires_database"
```

## Coverage Goals

- **Unit Tests:** 80%+ coverage
- **Integration Tests:** Critical paths covered
- **Overall:** 70%+ coverage

## Test Requirements

### Unit Tests
- No external dependencies
- Fast execution (< 1 second each)
- Use mocks for external services (database, AI service, RabbitMQ)

### Integration Tests
- May require database, RabbitMQ, and AI service
- Test with real external services
- Marked with `@pytest.mark.integration`
- Can be skipped if services unavailable

## Environment Variables for Testing

Create a `.env.test` file for integration tests:

```env
DATABASE_URL=postgresql://test:test@localhost/test
AI_SERVICE_URL=http://ai-service:8000
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

## Continuous Integration

Tests are run in CI/CD pipeline:
- All unit tests on every push
- Integration tests on merge to main
- Coverage reports generated
- Fail build if coverage < 70%

## Running Tests in Docker

```bash
# Build test image
docker build -f Dockerfile.test -t resume-worker-test .

# Run tests
docker run --rm resume-worker-test
```
