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
│   │   └── registry.py
│   ├── providers/
│   │   ├── factory.py
│   │   └── cipher.py
│   └── main.py
└── tests/
    ├── unit/
    │   ├── test_agents.py        ✅ Implemented
    │   ├── test_providers.py     ✅ Implemented
    │   └── test_api.py           ✅ Implemented
    ├── integration/
    │   └── test_api.py           ✅ Implemented
    ├── conftest.py               ✅ Implemented
    └── fixtures/
        └── test_messages.json    ✅ Implemented
```

### Implementation Status
**Status:** ✅ **COMPLETE** | **Coverage:** 85% | **Tests:** 30+ passing

**Test Files:**
- `tests/unit/test_agents.py` - Agent registry and system message tests
- `tests/unit/test_providers.py` - Provider factory and CipherClient tests
- `tests/unit/test_api.py` - API endpoint tests (chat, streaming, images)
- `tests/integration/test_api.py` - Integration tests with real providers
- `tests/conftest.py` - Shared fixtures and mocks
- `pytest.ini` - Test configuration
- `Dockerfile.test` - Docker test image
- `Makefile` - Test commands

### Unit Tests

#### Agent Tests
**Target Coverage:** 80%+ | **Current Coverage:** 85%+

- [x] Agent registry initialization
- [x] Agent retrieval by name
- [x] Agent system message building
- [x] Agent response handling
- [x] Error handling for invalid agents

#### Provider Tests
**Target Coverage:** 80%+ | **Current Coverage:** 85%+

- [x] Provider factory initialization
- [x] Model creation (OpenAI, Anthropic, XAI, etc.)
- [x] Provider selection logic
- [x] API client initialization
- [x] Error handling and retries

#### API Endpoint Tests
**Target Coverage:** 70%+ | **Current Coverage:** 85%+

- [x] POST `/v1/chat` - Chat completion
- [x] POST `/v1/chat/stream` - Streaming chat
- [x] POST `/v1/images` - Image generation
- [x] GET `/v1/agents` - List agents
- [x] GET `/healthz` - Health check
- [x] Request validation
- [x] Response formatting
- [x] Error responses

### Integration Tests

#### API Integration Tests
- [x] End-to-end chat flow
- [x] Streaming response handling
- [x] Multiple agent interactions
- [x] Error scenarios (provider failures)
- [ ] Authentication (if added)

#### Provider Integration Tests
- [x] Real API calls (with test API keys)
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

# Using Docker
docker build -f Dockerfile.test -t ai-service-test .
docker run --rm ai-service-test

# Using Makefile
make test
make test-unit
make test-cov
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
│   │   ├── openai_image.py
│   │   ├── stable_diffusion.py
│   │   ├── cipher_image.py
│   │   ├── diagram_generator.py
│   │   └── video_generator.py
│   └── main.py
└── tests/
    ├── unit/
    │   ├── test_providers.py           ✅ Implemented
    │   ├── test_provider_implementations.py ✅ Implemented
    │   └── test_api.py                 ✅ Implemented
    ├── integration/
    │   └── test_api.py                 ✅ Implemented
    ├── conftest.py                     ✅ Implemented
    └── fixtures/
        └── test_requests.json          ✅ Implemented
```

### Implementation Status
**Status:** ✅ **COMPLETE** | **Coverage:** 84% | **Tests:** 46 passing

**Test Files:**
- `tests/unit/test_providers.py` - Provider factory tests (Image, Diagram, Video)
- `tests/unit/test_provider_implementations.py` - Provider implementation tests
- `tests/unit/test_api.py` - API endpoint tests (images, diagrams, videos)
- `tests/integration/test_api.py` - Integration tests
- `tests/conftest.py` - Shared fixtures and mocks
- `pytest.ini` - Test configuration
- `Dockerfile.test` - Docker test image
- `Makefile` - Test commands

### Unit Tests

#### Provider Tests
**Target Coverage:** 80%+ | **Current Coverage:** 92%+

- [x] Provider factory initialization
- [x] Image generation provider selection
- [x] Video generation provider selection
- [x] Diagram generation logic
- [x] Provider-specific implementations

#### Image Generation Tests
- [x] Prompt processing
- [x] Image generation requests
- [x] Response parsing
- [x] Error handling

#### Video Generation Tests
- [x] Video generation requests
- [x] Provider selection (Replicate, Runway)
- [x] Response handling

### Integration Tests

#### API Integration Tests
- [x] POST `/v1/images/generate` - Image generation
- [x] POST `/v1/videos/generate` - Video generation
- [x] POST `/v1/diagrams/generate` - Diagram generation
- [x] GET `/v1/providers/*` - Provider listing
- [x] GET `/healthz` - Health check
- [x] Request validation
- [x] Error scenarios

### Running Tests
```bash
# All tests
pytest tests/ -v --cov=app

# Unit tests only
pytest tests/unit/ -v

# Integration tests
pytest tests/integration/ -v

# Using Docker
docker build -f Dockerfile.test -t creative-service-test .
docker run --rm creative-service-test

# Using Makefile
make test
make test-unit
make test-cov
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
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go           ✅ Implemented
│   ├── sender/
│   │   ├── smtp_sender.go
│   │   └── smtp_sender_test.go      ✅ Implemented
│   └── queue/
│       ├── queue.go
│       └── queue_test.go             ✅ Implemented
├── pkg/
│   ├── health/
│   │   ├── health.go
│   │   └── health_test.go            ✅ Implemented
│   └── logger/
│       ├── logger.go
│       └── logger_test.go            ✅ Implemented
└── tests/
    └── README.md                     ✅ Implemented
```

### Implementation Status
**Status:** ✅ **PARTIALLY COMPLETE** | **Coverage:** 49.7% overall | **Tests:** 30+ passing

**Coverage Breakdown:**
- `internal/config`: 100% ✅
- `pkg/health`: 94.6% ✅
- `internal/sender`: 65.0% ✅
- `pkg/logger`: 47.4% ✅
- `internal/queue`: 11.8% (requires RabbitMQ for full testing)

**Test Files:**
- `internal/config/config_test.go` - Configuration loading tests
- `internal/sender/smtp_sender_test.go` - SMTP sender tests
- `internal/queue/queue_test.go` - Queue structure tests
- `pkg/health/health_test.go` - Health check tests
- `pkg/logger/logger_test.go` - Logger tests
- `Makefile` - Test commands
- `Dockerfile.test` - Docker test image
- `tests/README.md` - Test documentation

**Note:** Queue package requires integration tests with real RabbitMQ instance to reach 70% coverage target.

### Unit Tests

#### Config Tests
**Target Coverage:** 90%+ | **Current Coverage:** 100% ✅

- [x] Email configuration loading
- [x] RabbitMQ configuration loading
- [x] Worker configuration loading
- [x] Environment variable parsing
- [x] Default values
- [x] Validation

#### SMTP Sender Tests
**Target Coverage:** 80%+ | **Current Coverage:** 65%

- [x] SMTP sender initialization
- [x] Email message formatting
- [x] MIME encoding
- [x] Connection handling
- [x] Error handling and retries
- [ ] Authentication (requires SMTP server)

#### Queue Tests
**Target Coverage:** 80%+ | **Current Coverage:** 11.8%

- [x] Message envelope serialization
- [x] Connection state checking
- [ ] Queue connection (requires RabbitMQ)
- [ ] Message consumption (requires RabbitMQ)
- [ ] Message acknowledgment (requires RabbitMQ)
- [ ] Message rejection/requeue (requires RabbitMQ)
- [ ] Error handling (requires RabbitMQ)

#### Health Check Tests
**Target Coverage:** 90%+ | **Current Coverage:** 94.6% ✅

- [x] Health check endpoint
- [x] RabbitMQ connection status
- [x] Response formatting
- [x] Caching behavior

#### Logger Tests
**Target Coverage:** 70%+ | **Current Coverage:** 47.4%

- [x] Logger initialization
- [x] Configuration handling
- [x] Trace ID context functions

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
go test ./internal/... ./pkg/... -v

# All tests with coverage
go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out

# Using Makefile
make test
make test-unit
make test-cov

# Using Docker
docker build -f Dockerfile.test -t email-worker-test .
docker run --rm email-worker-test
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
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go           ✅ Implemented
│   └── queue/
│       ├── queue.go
│       └── queue_test.go             ✅ Implemented
├── pkg/
│   ├── health/
│   │   ├── health.go
│   │   └── health_test.go            ✅ Implemented
│   └── logger/
│       ├── logger.go
│       └── (no tests yet)
└── tests/
    └── README.md                     ✅ Implemented
```

### Implementation Status
**Status:** 🔄 **PARTIALLY COMPLETE** | **Coverage:** 23.4% | **Tests:** 20+ passing

**Coverage Breakdown:**
- `internal/config`: 100% ✅
- `pkg/health`: 94.6% ✅
- `internal/queue`: 9.8%
- `pkg/logger`: 0% (not yet tested)
- `internal/notifier`: 0% (not yet tested)
- `cmd/whatsapp-worker`: 0% (not yet tested)

**Test Files:**
- `internal/config/config_test.go` - Configuration loading tests
- `internal/queue/queue_test.go` - Queue structure tests
- `pkg/health/health_test.go` - Health check tests
- `Makefile` - Test commands
- `Dockerfile.test` - Docker test image
- `tests/README.md` - Test documentation

**Note:** Queue package requires integration tests with real RabbitMQ instance to reach 70% coverage target. Notifier package needs tests for WhatsApp connection and message sending.

### Unit Tests

#### Config Tests
**Target Coverage:** 90%+ | **Current Coverage:** 100% ✅

- [x] RabbitMQ configuration loading
- [x] Worker configuration loading
- [x] WhatsApp configuration loading
- [x] Environment variable parsing
- [x] Default values
- [x] Validation

#### Queue Tests
**Target Coverage:** 80%+ | **Current Coverage:** 9.8%

- [x] Message envelope serialization
- [x] Connection state checking
- [x] Connection closing
- [x] Context cancellation
- [x] Invalid message handling
- [ ] Queue connection (requires RabbitMQ)
- [ ] Message consumption (requires RabbitMQ)
- [ ] Message acknowledgment (requires RabbitMQ)
- [ ] Message rejection/requeue (requires RabbitMQ)

#### Health Check Tests
**Target Coverage:** 90%+ | **Current Coverage:** 94.6% ✅

- [x] Health check initialization
- [x] Health check endpoint
- [x] RabbitMQ connection status
- [x] Response formatting
- [x] Caching behavior
- [x] Unhealthy status handling

#### WhatsApp Notifier Tests
**Target Coverage:** 70%+ | **Current Coverage:** 0%

- [ ] Notifier initialization
- [ ] Connection handling
- [ ] Message sending logic
- [ ] QR code generation (if applicable)
- [ ] Error handling
- [ ] Disconnection handling

### Integration Tests

#### Worker Integration Tests
- [ ] End-to-end WhatsApp message sending
- [ ] RabbitMQ message consumption
- [ ] Connection recovery
- [ ] Error handling and retries

### Running Tests
```bash
# Unit tests
go test ./internal/... ./pkg/... -v

# All tests with coverage
go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out

# Using Makefile
make test
make test-unit
make test-cov

# Using Docker
docker build -f Dockerfile.test -t whatsapp-worker-test .
docker run --rm whatsapp-worker-test
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
│   ├── ai_service.py
│   ├── database.py
│   ├── logger.py
│   ├── main.py
│   ├── queue_consumer.py
│   ├── translation_helper.py
│   └── keyword_extractor.py
└── tests/
    ├── unit/
    │   ├── test_resume_generator.py    ✅ Implemented
    │   ├── test_translation_helper.py  ✅ Implemented
    │   ├── test_keyword_extractor.py   ✅ Implemented
    │   ├── test_ai_service.py          ✅ Implemented
    │   ├── test_database.py            ✅ Implemented
    │   ├── test_logger.py              ✅ Implemented
    │   ├── test_main.py                ✅ Implemented
    │   └── test_queue_consumer.py      ✅ Implemented
    ├── integration/
    │   └── test_worker.py              ✅ Implemented
    └── conftest.py                     ✅ Implemented
```

### Implementation Status
**Status:** ✅ **COMPLETE** | **Coverage:** 76.78% | **Tests:** 93 passing

**Coverage Breakdown:**
- `logger.py`: 100% ✅
- `resume_generator.py`: 83% ✅
- `keyword_extractor.py`: 93% ✅
- `main.py`: 87% ✅
- `database.py`: 78% ✅
- `queue_consumer.py`: 74% ✅
- `health.py`: 74% ✅
- `ai_service.py`: 56%
- `translation_helper.py`: 61%

**Test Files:**
- `tests/unit/test_resume_generator.py` - Resume generation, HTML sanitization, template rendering
- `tests/unit/test_translation_helper.py` - Translation operations
- `tests/unit/test_keyword_extractor.py` - Keyword extraction from job descriptions
- `tests/unit/test_ai_service.py` - AI service client tests
- `tests/unit/test_database.py` - Database operations (projects, certifications, experiences, publications)
- `tests/unit/test_logger.py` - Logger configuration and trace ID management
- `tests/unit/test_main.py` - Main entry point (CLI mode, queue mode)
- `tests/unit/test_queue_consumer.py` - RabbitMQ queue consumer
- `tests/integration/test_worker.py` - Integration tests
- `pytest.ini` - Test configuration
- `Dockerfile.test` - Docker test image
- `Makefile` - Test commands

### Unit Tests

#### Resume Generator Tests
**Target Coverage:** 80%+ | **Current Coverage:** 83% ✅

- [x] Resume template rendering
- [x] HTML generation
- [x] HTML sanitization
- [x] PDF generation (mocked)
- [x] Data formatting
- [x] Template variable substitution
- [x] Language-specific templates

#### Translation Helper Tests
**Target Coverage:** 80%+ | **Current Coverage:** 61%

- [x] Translation logic
- [x] Language code mapping
- [x] Field translation (projects, certifications, experiences, posts)
- [x] Error handling
- [x] Connection management

#### Keyword Extractor Tests
**Target Coverage:** 70%+ | **Current Coverage:** 93% ✅

- [x] Keyword extraction logic
- [x] Skill matching
- [x] Tech category matching
- [x] Certification category matching
- [x] Case-insensitive matching
- [x] Multiple keyword variations (Golang, go programming, go language)

#### AI Service Tests
**Target Coverage:** 70%+ | **Current Coverage:** 56%

- [x] AI service client initialization
- [x] API request handling
- [x] Response parsing
- [x] Error handling
- [x] Section generation (profile, about, experience, skills)

#### Database Tests
**Target Coverage:** 70%+ | **Current Coverage:** 78% ✅

- [x] Database connection
- [x] User projects retrieval (with filters)
- [x] User certifications retrieval
- [x] User publications retrieval
- [x] User experiences retrieval
- [x] User info retrieval

#### Logger Tests
**Target Coverage:** 70%+ | **Current Coverage:** 100% ✅

- [x] Logger configuration (development/production)
- [x] File logging setup
- [x] Trace ID management
- [x] Service name binding

#### Main Entry Point Tests
**Target Coverage:** 70%+ | **Current Coverage:** 87% ✅

- [x] CLI mode execution
- [x] Queue mode execution
- [x] Job processing
- [x] Result saving
- [x] Error handling

#### Queue Consumer Tests
**Target Coverage:** 70%+ | **Current Coverage:** 74% ✅

- [x] Queue connection
- [x] Message consumption
- [x] Message acknowledgment
- [x] Message rejection/requeue
- [x] Error handling
- [x] Health check integration

### Integration Tests

#### Worker Integration Tests
- [x] Health check endpoint
- [ ] End-to-end resume generation flow (requires full stack)
- [ ] RabbitMQ message consumption (requires RabbitMQ)
- [ ] Database operations (requires database)
- [ ] AI service integration (requires AI service)
- [ ] Translation workflow (requires database)

### Running Tests
```bash
# Unit tests
pytest tests/unit/ -v --cov=src

# Integration tests
pytest tests/integration/ -v

# All tests with coverage
pytest --cov=src --cov-report=html --cov-report=term

# Using Docker
docker build -f Dockerfile.test -t resume-worker-test .
docker run --rm resume-worker-test

# Using Makefile
make test
make test-unit
make test-cov
```

---

## Job Application Worker (Node.js)

### Testing Framework
- **Unit Tests:** Jest
- **Integration Tests:** Jest + test RabbitMQ container
- **Mocking:** Jest mocks (`jest.fn()`, `jest.unstable_mockModule()`)
- **Coverage:** Jest built-in coverage

### Test Structure
```
job-application-worker/
├── src/
│   ├── worker.js
│   ├── orchestrator.js
│   ├── coverLetter.js
│   ├── health.js
│   └── ...
└── tests/
    ├── __tests__/
    │   ├── health.test.js          ✅ Implemented
    │   ├── orchestrator.test.js    ✅ Implemented
    │   └── coverLetter.test.js     ✅ Implemented
    ├── setup.js                    ✅ Implemented
    └── README.md                   ✅ Implemented
```

### Implementation Status
**Status:** 🔄 **PARTIALLY COMPLETE** | **Coverage:** 15.46% | **Tests:** 23 passing

**Coverage Breakdown:**
- `orchestrator.js`: 100% ✅
- `health.js`: 94.73% ✅
- `coverLetter.js`: 41.33%
- `logger.js` (utils): 63.26%
- Other modules: 0% (not yet tested)

**Test Files:**
- `tests/__tests__/health.test.js` - Health check tests
- `tests/__tests__/orchestrator.test.js` - Website rate limiting and count tracking
- `tests/__tests__/coverLetter.test.js` - Cover letter generation
- `jest.config.js` - Jest configuration for ES modules
- `tests/setup.js` - Jest setup for global `jest` availability
- `Dockerfile.test` - Docker test image
- `Makefile` - Test commands
- `tests/README.md` - Test documentation

**Note:** Coverage is low because only a few modules are tested. To reach 70%, tests are needed for:
- `worker.js` (0%)
- `scraper.js` (0%)
- `database.js` (0%)
- `queue.js` / `queue_rabbitmq.js` (0%)
- `aiSelectorFinder.js` (0%)
- `selectorCache.js` (0%)
- `selfHealingScraper.js` (0%)

### Unit Tests

#### Health Check Tests
**Target Coverage:** 90%+ | **Current Coverage:** 94.73% ✅

- [x] Health check endpoint
- [x] RabbitMQ connection status
- [x] Response formatting

#### Orchestrator Tests
**Target Coverage:** 80%+ | **Current Coverage:** 100% ✅

- [x] Website rate limiting
- [x] Website count tracking
- [x] Should process logic
- [x] Count increment logic
- [x] New day reset logic

#### Cover Letter Service Tests
**Target Coverage:** 70%+ | **Current Coverage:** 41.33%

- [x] Cover letter generation
- [x] Profile data processing
- [x] Job info processing
- [x] Error handling
- [ ] Edge cases and retries

#### Worker Tests
**Target Coverage:** 70%+ | **Current Coverage:** 0%

- [ ] Worker initialization
- [ ] Job processing logic
- [ ] Error handling
- [ ] Retry logic
- [ ] Graceful shutdown

#### Scraper Tests
**Target Coverage:** 80%+ | **Current Coverage:** 0%

- [ ] Playwright browser initialization
- [ ] Selector finding logic
- [ ] Form filling logic
- [ ] Job application submission
- [ ] Error handling
- [ ] Selector caching

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
  coverageThreshold: {
    global: {
      statements: 70,
      branches: 70,
      functions: 70,
      lines: 70
    }
  },
  setupFilesAfterEnv: ['<rootDir>/tests/setup.js']
};
```

### Running Tests
```bash
# Unit tests
npm run test:jest

# With coverage
npm run test:jest:coverage

# Using Docker
docker build -f Dockerfile.test -t job-application-worker-test .
docker run --rm job-application-worker-test

# Using Makefile
make test
make test-unit
make test-cov
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

### Phase 1: Foundation (Week 1-2) ✅ COMPLETE
- [x] Set up test infrastructure (Docker Compose, test helpers)
- [x] Create example unit tests for services (AI Service, Creative Service)
- [x] Create example integration tests for service endpoints
- [x] Set up coverage reporting

### Phase 2: Core Components (Week 3-4) ✅ IN PROGRESS
- [ ] Server: Auth, User Profiles, Projects domains
- [ ] Server: Critical API endpoints
- [x] Workers: Queue and message processing tests (Email Worker)
- [x] Services: Provider and agent tests (AI Service, Creative Service)

### Phase 3: Expansion (Week 5-6) 🔄 IN PROGRESS
- [ ] Server: Remaining domains
- [ ] Workers: End-to-end integration tests
- [x] Services: API integration tests (AI Service, Creative Service)
- [x] Edge cases and error scenarios (AI Service, Creative Service)

### Phase 4: Optimization (Week 7-8)
- [x] Improve coverage to target levels (AI Service: 85%, Creative Service: 84%)
- [ ] Improve coverage for Email Worker (currently 49.7%, needs integration tests)
- [ ] Optimize test execution time
- [ ] Add performance tests (if needed)
- [x] Document test patterns and best practices

## Current Implementation Status

### ✅ Completed
- **AI Service:** 85% coverage, 30+ tests passing
- **Creative Service:** 84% coverage, 46 tests passing
- **Email Worker:** 49.7% coverage, 30+ tests passing (config: 100%, health: 94.6%)
- **Resume Worker:** 76.78% coverage, 93 tests passing ✅

### 🔄 In Progress
- **Email Worker:** Integration tests needed for queue package to reach 70% target
- **Job Application Worker:** 15.46% coverage, 23 tests passing (needs more test coverage)
- **WhatsApp Worker:** 23.4% coverage, 20+ tests passing (config: 100%, health: 94.6%, needs notifier and queue integration tests)

### ⏳ Pending
- **Server:** Unit and integration tests
- **Translation Worker:** Unit and integration tests

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
