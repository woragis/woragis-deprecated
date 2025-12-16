# Testing Strategy and Plans

## Overview

This document outlines comprehensive testing strategies for all backend components: server, services, and workers. Each component requires unit tests for isolated logic and integration tests for end-to-end workflows.

---

## Server (Go)

### Testing Framework
- **Unit Tests:** `testing` package (standard library) + `testify` for assertions
- **Integration Tests:** `testify` + Docker Compose for dependencies
- **HTTP Tests:** `net/http/httptest` for API endpoint testing
- **Mocking:** `gomock` or manual mocks for external dependencies

### Test Structure
```
server/
├── app/
│   ├── internal/
│   │   ├── domains/
│   │   │   └── {domain}/
│   │   │       ├── {domain}.go
│   │   │       ├── {domain}_test.go          # Unit tests
│   │   │       └── {domain}_integration_test.go # Integration tests
│   │   └── services/
│   │       └── {service}/
│   │           ├── {service}.go
│   │           └── {service}_test.go
│   └── pkg/
│       └── {package}/
│           ├── {package}.go
│           └── {package}_test.go
└── test/
    ├── fixtures/          # Test data
    ├── helpers/           # Test utilities
    └── integration/      # Integration test suites
```

### Unit Tests

#### Domain Layer
**Target Coverage:** 80%+

**Auth Domain:**
- [ ] User registration/login flows
- [ ] JWT token generation/validation
- [ ] OAuth provider integration (Google, GitHub, Microsoft)
- [ ] Password hashing/verification
- [ ] Session management
- [ ] Permission checks

**User Profiles Domain:**
- [ ] Profile creation/update
- [ ] Profile validation
- [ ] Profile retrieval
- [ ] Profile deletion

**Projects Domain:**
- [ ] Project CRUD operations
- [ ] Project case studies
- [ ] Project validation
- [ ] Project relationships (skills, technologies)

**Job Applications Domain:**
- [ ] Application creation/update
- [ ] Status transitions
- [ ] Interview stages management
- [ ] Application responses

**Other Domains:**
- Skills, Certifications, Experiences, Languages, etc.
- Each domain should have comprehensive CRUD tests
- Validation logic tests
- Business rule enforcement tests

#### Service Layer
**Target Coverage:** 70%+

**Email Service:**
- [ ] SMTP sender initialization
- [ ] Email template rendering
- [ ] Email sending logic
- [ ] Error handling and retries

**AI Service:**
- [ ] LangChain client initialization
- [ ] Chat completion handling
- [ ] Error handling

**Creative Service:**
- [ ] Client initialization
- [ ] API request handling
- [ ] Response parsing

**WhatsApp Service:**
- [ ] WhatsApp notifier initialization
- [ ] Message sending logic
- [ ] Connection management

#### Package Layer
**Target Coverage:** 90%+

**Logger:**
- [ ] Log level configuration
- [ ] JSON/text format switching
- [ ] Trace ID propagation
- [ ] File logging (development)

**Health:**
- [ ] Health check logic
- [ ] Dependency status checks
- [ ] Caching behavior
- [ ] Response formatting

**RabbitMQ:**
- [ ] Connection management
- [ ] Queue declaration
- [ ] Exchange declaration
- [ ] Message publishing/consuming

**Config:**
- [ ] Configuration loading
- [ ] Environment variable parsing
- [ ] Default values
- [ ] Validation

### Integration Tests

#### API Endpoint Tests
**Framework:** `net/http/httptest` + test database

**Test Infrastructure:**
- Docker Compose setup with PostgreSQL, Redis, RabbitMQ
- Test database with migrations
- Test fixtures for common data
- Test helpers for authentication

**Auth Endpoints:**
- [ ] POST `/api/auth/register` - User registration
- [ ] POST `/api/auth/login` - User login
- [ ] POST `/api/auth/logout` - User logout
- [ ] GET `/api/auth/me` - Current user info
- [ ] OAuth callback handlers

**User Profile Endpoints:**
- [ ] GET `/api/user-profiles/{id}` - Get profile
- [ ] PUT `/api/user-profiles/{id}` - Update profile
- [ ] POST `/api/user-profiles` - Create profile

**Projects Endpoints:**
- [ ] GET `/api/projects` - List projects
- [ ] POST `/api/projects` - Create project
- [ ] PUT `/api/projects/{id}` - Update project
- [ ] DELETE `/api/projects/{id}` - Delete project
- [ ] Project case studies endpoints

**Job Applications Endpoints:**
- [ ] POST `/api/job-applications` - Create application
- [ ] GET `/api/job-applications` - List applications
- [ ] PUT `/api/job-applications/{id}` - Update application
- [ ] Status transition endpoints

**Other Endpoints:**
- Skills, Certifications, Experiences, etc.
- Test authentication/authorization
- Test validation errors
- Test pagination
- Test filtering/sorting

#### Database Integration Tests
- [ ] Migration tests
- [ ] Transaction handling
- [ ] Foreign key constraints
- [ ] Index usage
- [ ] Query performance

#### External Service Integration Tests
- [ ] RabbitMQ message publishing/consuming
- [ ] Redis operations
- [ ] External API mocking (AI, Creative services)

### Test Data Management
- **Fixtures:** JSON/YAML files for test data
- **Factories:** Go functions to generate test entities
- **Seeders:** Database seeders for integration tests
- **Cleanup:** Automatic cleanup after tests

### Running Tests
```bash
# Unit tests
go test ./app/internal/domains/... -v -cover

# Integration tests (requires Docker Compose)
docker-compose -f docker-compose.test.yml up -d
go test ./test/integration/... -v -tags=integration
docker-compose -f docker-compose.test.yml down

# All tests with coverage
go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## AI Service (Python)

### Testing Framework
- **Unit Tests:** `pytest` + `pytest-mock` for mocking
- **Integration Tests:** `pytest` + `httpx` for HTTP testing
- **Async Tests:** `pytest-asyncio` for async operations
- **Coverage:** `pytest-cov`

### Test Structure
```
ai-service/
├── app/
│   ├── agents/
│   │   ├── registry.py
│   │   └── test_registry.py
│   ├── providers/
│   │   ├── factory.py
│   │   └── test_factory.py
│   ├── main.py
│   └── test_main.py
└── tests/
    ├── unit/
    │   ├── test_agents.py
    │   ├── test_providers.py
    │   └── test_logger.py
    ├── integration/
    │   ├── test_api.py
    │   └── test_agents_integration.py
    └── fixtures/
        └── test_messages.json
```

### Unit Tests

#### Agent Tests
**Target Coverage:** 80%+

- [ ] Agent registry initialization
- [ ] Agent retrieval by name
- [ ] Agent system message building
- [ ] Agent response handling
- [ ] Error handling for invalid agents

#### Provider Tests
**Target Coverage:** 80%+

- [ ] Provider factory initialization
- [ ] Model creation (OpenAI, Anthropic, XAI, etc.)
- [ ] Provider selection logic
- [ ] API client initialization
- [ ] Error handling and retries

#### API Endpoint Tests
**Target Coverage:** 70%+

- [ ] POST `/chat` - Chat completion
- [ ] POST `/chat/stream` - Streaming chat
- [ ] Request validation
- [ ] Response formatting
- [ ] Error responses

### Integration Tests

#### API Integration Tests
- [ ] End-to-end chat flow
- [ ] Streaming response handling
- [ ] Multiple agent interactions
- [ ] Error scenarios (provider failures)
- [ ] Authentication (if added)

#### Provider Integration Tests
- [ ] Real API calls (with test API keys)
- [ ] Rate limiting handling
- [ ] Timeout handling
- [ ] Retry logic

### Test Configuration
```python
# pytest.ini
[pytest]
testpaths = tests
python_files = test_*.py
python_classes = Test*
python_functions = test_*
addopts = 
    -v
    --cov=app
    --cov-report=html
    --cov-report=term
```

### Running Tests
```bash
# Unit tests
pytest tests/unit/ -v

# Integration tests
pytest tests/integration/ -v

# With coverage
pytest --cov=app --cov-report=html

# Specific test
pytest tests/unit/test_agents.py::test_agent_registry -v
```

---

## Creative Service (Python)

### Testing Framework
- **Unit Tests:** `pytest` + `pytest-mock`
- **Integration Tests:** `pytest` + `httpx`
- **Coverage:** `pytest-cov`

### Test Structure
```
creative-service/
├── app/
│   ├── providers/
│   │   ├── factory.py
│   │   └── test_factory.py
│   ├── main.py
│   └── test_main.py
└── tests/
    ├── unit/
    │   ├── test_providers.py
    │   └── test_image_generation.py
    ├── integration/
    │   └── test_api.py
    └── fixtures/
        └── test_requests.json
```

### Unit Tests

#### Provider Tests
**Target Coverage:** 80%+

- [ ] Provider factory initialization
- [ ] Image generation provider selection
- [ ] Video generation provider selection
- [ ] Diagram generation logic
- [ ] Provider-specific implementations

#### Image Generation Tests
- [ ] Prompt processing
- [ ] Image generation requests
- [ ] Response parsing
- [ ] Error handling

#### Video Generation Tests
- [ ] Video generation requests
- [ ] Provider selection (Replicate, Runway)
- [ ] Response handling

### Integration Tests

#### API Integration Tests
- [ ] POST `/generate/image` - Image generation
- [ ] POST `/generate/video` - Video generation
- [ ] POST `/generate/diagram` - Diagram generation
- [ ] Request validation
- [ ] Error scenarios

### Running Tests
```bash
# All tests
pytest tests/ -v --cov=app

# Unit tests only
pytest tests/unit/ -v

# Integration tests
pytest tests/integration/ -v
```

---

## Email Worker (Go)

### Testing Framework
- **Unit Tests:** `testing` + `testify`
- **Integration Tests:** `testify` + test RabbitMQ container
- **Mocking:** Manual mocks for SMTP and RabbitMQ

### Test Structure
```
email-worker/
├── internal/
│   ├── sender/
│   │   ├── smtp_sender.go
│   │   └── smtp_sender_test.go
│   └── queue/
│       ├── queue.go
│       └── queue_test.go
└── test/
    └── integration/
        └── worker_test.go
```

### Unit Tests

#### SMTP Sender Tests
**Target Coverage:** 80%+

- [ ] SMTP sender initialization
- [ ] Email message formatting
- [ ] Template rendering
- [ ] Connection handling
- [ ] Error handling and retries
- [ ] Authentication

#### Queue Tests
**Target Coverage:** 80%+

- [ ] Queue connection
- [ ] Message consumption
- [ ] Message acknowledgment
- [ ] Message rejection/requeue
- [ ] Error handling

### Integration Tests

#### Worker Integration Tests
- [ ] End-to-end email sending flow
- [ ] RabbitMQ message consumption
- [ ] SMTP email delivery
- [ ] Error handling and retries
- [ ] Graceful shutdown

#### Test Infrastructure
- Docker Compose with RabbitMQ
- Mock SMTP server (MailHog or similar)
- Test message fixtures

### Running Tests
```bash
# Unit tests
go test ./internal/... -v -cover

# Integration tests
docker-compose -f docker-compose.test.yml up -d
go test ./test/integration/... -v -tags=integration
docker-compose -f docker-compose.test.yml down
```

---

## WhatsApp Worker (Go)

### Testing Framework
- **Unit Tests:** `testing` + `testify`
- **Integration Tests:** `testify` + test RabbitMQ container
- **Mocking:** Manual mocks for WhatsApp and RabbitMQ

### Test Structure
```
whatsapp-worker/
├── internal/
│   ├── notifier/
│   │   ├── whatsmeow_notifier.go
│   │   └── whatsmeow_notifier_test.go
│   └── queue/
│       ├── queue.go
│       └── queue_test.go
└── test/
    └── integration/
        └── worker_test.go
```

### Unit Tests

#### WhatsApp Notifier Tests
**Target Coverage:** 70%+

- [ ] Notifier initialization
- [ ] Connection handling
- [ ] Message sending logic
- [ ] QR code generation (if applicable)
- [ ] Error handling
- [ ] Disconnection handling

#### Queue Tests
**Target Coverage:** 80%+

- [ ] Queue connection
- [ ] Message consumption
- [ ] Message processing
- [ ] Error handling

### Integration Tests

#### Worker Integration Tests
- [ ] End-to-end WhatsApp message sending
- [ ] RabbitMQ message consumption
- [ ] Connection recovery
- [ ] Error handling and retries

### Running Tests
```bash
# Unit tests
go test ./internal/... -v -cover

# Integration tests (requires test WhatsApp account or mocks)
go test ./test/integration/... -v -tags=integration
```

---

## Resume Worker (Python)

### Testing Framework
- **Unit Tests:** `pytest` + `pytest-mock`
- **Integration Tests:** `pytest` + test RabbitMQ container
- **Coverage:** `pytest-cov`

### Test Structure
```
resume-worker/
├── src/
│   ├── resume_generator.py
│   ├── test_resume_generator.py
│   ├── ai_service.py
│   ├── test_ai_service.py
│   └── queue_consumer.py
└── tests/
    ├── unit/
    │   ├── test_resume_generation.py
    │   ├── test_translation.py
    │   └── test_keyword_extraction.py
    ├── integration/
    │   └── test_worker.py
    └── fixtures/
        └── test_resume_data.json
```

### Unit Tests

#### Resume Generator Tests
**Target Coverage:** 80%+

- [ ] Resume template rendering
- [ ] HTML generation
- [ ] PDF generation (if applicable)
- [ ] Data formatting
- [ ] Template variable substitution

#### Translation Helper Tests
**Target Coverage:** 80%+

- [ ] Translation logic
- [ ] Language detection
- [ ] Field translation
- [ ] Error handling

#### Keyword Extractor Tests
**Target Coverage:** 70%+

- [ ] Keyword extraction logic
- [ ] Skill matching
- [ ] Experience parsing

#### AI Service Tests
**Target Coverage:** 70%+

- [ ] AI service client initialization
- [ ] API request handling
- [ ] Response parsing
- [ ] Error handling

### Integration Tests

#### Worker Integration Tests
- [ ] End-to-end resume generation flow
- [ ] RabbitMQ message consumption
- [ ] Database operations
- [ ] AI service integration
- [ ] Translation workflow

### Running Tests
```bash
# Unit tests
pytest tests/unit/ -v --cov=src

# Integration tests
docker-compose -f docker-compose.test.yml up -d
pytest tests/integration/ -v
docker-compose -f docker-compose.test.yml down
```

---

## Job Application Worker (Node.js)

### Testing Framework
- **Unit Tests:** Jest or Mocha + Chai
- **Integration Tests:** Jest + test RabbitMQ container
- **Mocking:** Jest mocks or Sinon
- **Coverage:** `nyc` or `c8`

### Test Structure
```
job-application-worker/
├── src/
│   ├── worker.js
│   ├── worker.test.js
│   ├── scraper.js
│   ├── scraper.test.js
│   └── orchestrator.js
└── tests/
    ├── unit/
    │   ├── worker.test.js
    │   ├── scraper.test.js
    │   └── orchestrator.test.js
    ├── integration/
    │   └── worker.integration.test.js
    └── fixtures/
        └── test_jobs.json
```

### Unit Tests

#### Worker Tests
**Target Coverage:** 70%+

- [ ] Worker initialization
- [ ] Job processing logic
- [ ] Error handling
- [ ] Retry logic
- [ ] Graceful shutdown

#### Scraper Tests
**Target Coverage:** 80%+

- [ ] Playwright browser initialization
- [ ] Selector finding logic
- [ ] Form filling logic
- [ ] Job application submission
- [ ] Error handling
- [ ] Selector caching

#### Orchestrator Tests
**Target Coverage:** 80%+

- [ ] Website rate limiting
- [ ] Website count tracking
- [ ] Should process logic
- [ ] Count increment logic

#### Cover Letter Service Tests
**Target Coverage:** 70%+

- [ ] Cover letter generation
- [ ] Profile data processing
- [ ] Job info processing

### Integration Tests

#### Worker Integration Tests
- [ ] End-to-end job application flow
- [ ] RabbitMQ message consumption
- [ ] Database operations
- [ ] Scraper integration
- [ ] Error handling and retries

### Test Configuration
```javascript
// jest.config.js
module.exports = {
  testEnvironment: 'node',
  coverageDirectory: 'coverage',
  collectCoverageFrom: ['src/**/*.js'],
  testMatch: ['**/__tests__/**/*.js', '**/*.test.js'],
  setupFilesAfterEnv: ['<rootDir>/tests/setup.js']
};
```

### Running Tests
```bash
# Unit tests
npm test

# Integration tests
docker-compose -f docker-compose.test.yml up -d
npm run test:integration
docker-compose -f docker-compose.test.yml down

# With coverage
npm run test:coverage
```

---

## Translation Worker (Go)

### Testing Framework
- **Unit Tests:** `testing` + `testify`
- **Integration Tests:** `testify` + test RabbitMQ container
- **Mocking:** Manual mocks for AI service and RabbitMQ

### Test Structure
```
server/app/cmd/translation-worker/
├── pkg/
│   └── health/
│       ├── health.go
│       └── health_test.go
└── test/
    └── integration/
        └── worker_test.go
```

### Unit Tests

#### Translation Service Tests
**Target Coverage:** 80%+

- [ ] Translation logic
- [ ] Language detection
- [ ] Content enrichment
- [ ] AI service integration
- [ ] Error handling

#### Queue Tests
**Target Coverage:** 80%+

- [ ] Queue connection
- [ ] Message consumption
- [ ] Message processing
- [ ] Error handling

### Integration Tests

#### Worker Integration Tests
- [ ] End-to-end translation flow
- [ ] RabbitMQ message consumption
- [ ] Database operations
- [ ] AI service integration
- [ ] Error handling and retries

### Running Tests
```bash
# Unit tests
go test ./... -v -cover

# Integration tests
docker-compose -f docker-compose.test.yml up -d
go test ./test/integration/... -v -tags=integration
docker-compose -f docker-compose.test.yml down
```

---

## Test Infrastructure

### Docker Compose for Testing
```yaml
# docker-compose.test.yml
version: '3.8'
services:
  postgres-test:
    image: postgres:15
    environment:
      POSTGRES_DB: woragis_test
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
    ports:
      - "5433:5432"
  
  redis-test:
    image: redis:7-alpine
    ports:
      - "6380:6379"
  
  rabbitmq-test:
    image: rabbitmq:3-management-alpine
    ports:
      - "5673:5672"
      - "15673:15672"
    environment:
      RABBITMQ_DEFAULT_USER: test
      RABBITMQ_DEFAULT_PASS: test
```

### Test Helpers and Utilities

#### Go Test Helpers
```go
// test/helpers/db.go
func SetupTestDB(t *testing.T) *gorm.DB {
    // Create test database connection
}

func CleanupTestDB(t *testing.T, db *gorm.DB) {
    // Clean up test data
}

// test/helpers/auth.go
func CreateTestUser(t *testing.T, db *gorm.DB) *User {
    // Create test user
}

func GetTestToken(t *testing.T, user *User) string {
    // Generate JWT token for test user
}
```

#### Python Test Fixtures
```python
# tests/fixtures/conftest.py
@pytest.fixture
def test_db():
    # Setup test database
    yield db
    # Cleanup

@pytest.fixture
def mock_rabbitmq():
    # Mock RabbitMQ connection
    pass
```

#### Node.js Test Helpers
```javascript
// tests/helpers/db.js
export async function setupTestDB() {
  // Setup test database
}

export async function cleanupTestDB() {
  // Cleanup test data
}
```

---

## Coverage Goals

### Overall Targets
- **Unit Tests:** 70-80% coverage minimum
- **Integration Tests:** Critical paths covered
- **Critical Components:** 90%+ coverage (auth, payments, etc.)

### Coverage by Component
- **Server:** 75% unit, 50% integration
- **Services:** 70% unit, 40% integration
- **Workers:** 70% unit, 50% integration

---

## Best Practices

### Test Organization
1. **One test file per source file** (`*_test.go`, `test_*.py`, `*.test.js`)
2. **Group related tests** using subtests or test suites
3. **Use descriptive test names** that explain what is being tested
4. **Keep tests isolated** - no shared state between tests

### Test Data
1. **Use fixtures** for complex test data
2. **Create factories** for generating test entities
3. **Clean up after tests** - always reset state
4. **Use realistic data** that matches production scenarios

### Test Execution
1. **Fast unit tests** - should complete in < 1 second each
2. **Integration tests** - should complete in < 30 seconds each
3. **Run tests in parallel** when possible
4. **Skip slow tests** in development, run in CI

### Mocking
1. **Mock external dependencies** (databases, APIs, services)
2. **Use real implementations** in integration tests
3. **Verify mock interactions** when testing behavior
4. **Keep mocks simple** - don't over-mock

### Error Testing
1. **Test error cases** - failures, timeouts, invalid input
2. **Test edge cases** - empty data, null values, boundaries
3. **Test retry logic** - ensure retries work correctly
4. **Test graceful degradation** - partial failures

### CI Integration
1. **Run tests on every push/PR**
2. **Fail builds on test failures**
3. **Generate coverage reports**
4. **Track coverage trends over time**
5. **Run integration tests in CI** (with test containers)

---

## Implementation Roadmap

### Phase 1: Foundation (Week 1-2)
- [ ] Set up test infrastructure (Docker Compose, test helpers)
- [ ] Create example unit tests for one domain (e.g., Auth)
- [ ] Create example integration tests for one endpoint
- [ ] Set up coverage reporting

### Phase 2: Core Components (Week 3-4)
- [ ] Server: Auth, User Profiles, Projects domains
- [ ] Server: Critical API endpoints
- [ ] Workers: Queue and message processing tests
- [ ] Services: Provider and agent tests

### Phase 3: Expansion (Week 5-6)
- [ ] Server: Remaining domains
- [ ] Workers: End-to-end integration tests
- [ ] Services: API integration tests
- [ ] Edge cases and error scenarios

### Phase 4: Optimization (Week 7-8)
- [ ] Improve coverage to target levels
- [ ] Optimize test execution time
- [ ] Add performance tests (if needed)
- [ ] Document test patterns and best practices

---

## Continuous Improvement

### Regular Reviews
- Review test coverage monthly
- Identify gaps in test coverage
- Refactor tests for maintainability
- Update test strategies based on learnings

### Metrics to Track
- Test coverage percentage
- Test execution time
- Flaky test rate
- Test failure rate
- Time to fix failing tests
