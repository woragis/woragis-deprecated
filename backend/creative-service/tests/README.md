# Creative Service Tests

## Overview

This directory contains unit and integration tests for the Creative Service.

## Test Structure

```
tests/
├── unit/              # Unit tests (fast, isolated)
│   ├── test_providers.py
│   └── test_api.py
├── integration/       # Integration tests (require external services)
│   └── test_api.py
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
pytest --cov=app --cov-report=html
```

### Specific Test File
```bash
pytest tests/unit/test_api.py -v
```

## Test Markers

- `@pytest.mark.unit` - Unit tests (fast, isolated)
- `@pytest.mark.integration` - Integration tests (require external services)
- `@pytest.mark.requires_api_key` - Tests that require API keys
- `@pytest.mark.slow` - Slow running tests

## Running Marked Tests

```bash
# Only unit tests
pytest -m unit

# Only integration tests
pytest -m integration

# Skip tests requiring API keys
pytest -m "not requires_api_key"
```

## Coverage Goals

- **Unit Tests:** 80%+ coverage
- **Integration Tests:** Critical paths covered
- **Overall:** 70%+ coverage

## Test Requirements

### Unit Tests
- No external dependencies
- Fast execution (< 1 second each)
- Use mocks for external services

### Integration Tests
- May require API keys (set in environment)
- Test with real external services
- Marked with `@pytest.mark.integration`
- Can be skipped if services unavailable

## Environment Variables for Testing

Create a `.env.test` file for integration tests:

```env
OPENAI_API_KEY=test-key
ANTHROPIC_API_KEY=test-key
REPLICATE_API_KEY=test-key
CIPHER_API_KEY=test-key
RUNWAY_API_KEY=test-key
```

## Continuous Integration

Tests are run in CI/CD pipeline:
- All unit tests on every push
- Integration tests on merge to main
- Coverage reports generated
- Fail build if coverage < 70%
