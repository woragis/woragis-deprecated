# Job Application Worker Tests

## Overview

This directory contains unit and integration tests for the Job Application Worker.

## Test Structure

```
tests/
├── __tests__/          # Unit tests
│   ├── health.test.js
│   ├── orchestrator.test.js
│   └── coverLetter.test.js
└── README.md
```

## Running Tests

### All Tests
```bash
npm run test:jest
```

### With Coverage
```bash
npm run test:jest:coverage
```

### Specific Test File
```bash
npm run test:jest -- tests/__tests__/health.test.js
```

## Coverage Goals

- **Unit Tests:** 70%+ coverage
- **Integration Tests:** Critical paths covered
- **Overall:** 70%+ coverage

## Test Requirements

### Unit Tests
- No external dependencies
- Fast execution (< 1 second each)
- Use mocks for external services (database, RabbitMQ, AI service, Playwright)

### Integration Tests
- May require database, RabbitMQ, and AI service
- Test with real external services
- Can be skipped if services unavailable

## Continuous Integration

Tests are run in CI/CD pipeline:
- All unit tests on every push
- Integration tests on merge to main
- Coverage reports generated
- Fail build if coverage < 70%

## Running Tests in Docker

```bash
# Build test image
docker build -f Dockerfile.test -t job-application-worker-test .

# Run tests
docker run --rm job-application-worker-test
```
