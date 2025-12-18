# Testing Overview - Backend Architecture

## General Architecture

All backend components implement **automated testing** with language-specific frameworks. The system uses **consistent patterns** across similar components with **framework-specific implementations**.

### Common Patterns

1. **Test Structure**: Tests organized alongside source code or in dedicated `tests/` directory
2. **Coverage Target**: 70% minimum coverage for all components
3. **Test Types**: Unit tests (fast, isolated) and Integration tests (end-to-end)
4. **Docker Support**: `Dockerfile.test` for consistent test environments
5. **Makefile Commands**: Standardized `make test`, `make test-unit`, `make test-cov` commands
6. **Coverage Reports**: HTML and terminal coverage reports

---

## Component Breakdown

### 1. **Go Components** (Server, Email, WhatsApp, Translation Workers)
**Framework**: Go `testing` package (standard library)  
**Test Files**: `*_test.go` (alongside source files)  
**Assertions**: Manual assertions or `testify` (if used)

#### Test Structure:
```
component/
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go          ✅ Tests alongside source
│   ├── queue/
│   │   ├── queue.go
│   │   └── queue_test.go          ✅ Tests alongside source
│   └── sender/ (email-worker)
│       ├── smtp_sender.go
│       └── smtp_sender_test.go    ✅ Tests alongside source
├── pkg/
│   ├── health/
│   │   ├── health.go
│   │   └── health_test.go         ✅ Tests alongside source
│   └── logger/
│       ├── logger.go
│       └── logger_test.go         ✅ Tests alongside source
└── cmd/
    └── {worker}/
        ├── main.go
        └── main_test.go           ✅ Smoke tests
```

#### Test Patterns:
- **Table-driven tests**: Common pattern for multiple scenarios
- **Environment variable manipulation**: Save/restore pattern for config tests
- **Mock interfaces**: Manual mocks for dependencies
- **HTTP test servers**: `httptest` for HTTP endpoint testing
- **Context testing**: Context cancellation, timeouts

#### Example Test:
```go
func TestLoadConfig(t *testing.T) {
    tests := []struct {
        name string
        env  map[string]string
        want string
    }{
        {
            name: "valid config",
            env:  map[string]string{"KEY": "value"},
            want: "value",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Save original env
            originalEnv := os.Getenv("KEY")
            defer os.Setenv("KEY", originalEnv)
            
            // Set test env
            os.Setenv("KEY", tt.env["KEY"])
            
            // Test
            cfg := LoadConfig()
            if cfg.Key != tt.want {
                t.Errorf("Key = %v, want %v", cfg.Key, tt.want)
            }
        })
    }
}
```

#### Running Tests:
```bash
# All tests
go test ./... -v

# Unit tests only
go test ./internal/... ./pkg/... -v

# With coverage
go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out
go tool cover -func=coverage.out

# Coverage threshold check
make test-cov-check  # Fails if < 70%
```

#### Coverage Status:
- **Email Worker**: 49.7% (config: 100%, health: 94.6%, queue: 11.8%)
- **WhatsApp Worker**: 23.4% (config: 100%, health: 94.6%, queue: 9.8%)
- **Translation Worker**: 34.3% (config: 100%, health: 94.6%, translator: 32.8%)

#### Special Features:
- **CGO Handling**: WhatsApp worker tests skip CGO-dependent tests if CGO not enabled
- **Panic Recovery**: Some tests use `t.Skip()` for complex edge cases
- **Mock RabbitMQ**: Limited mocking (requires real RabbitMQ for full coverage)

---

### 2. **Python Services** (AI Service, Creative Service)
**Framework**: `pytest` + `pytest-mock` + `pytest-asyncio`  
**Test Files**: `test_*.py` in `tests/unit/` and `tests/integration/`  
**Fixtures**: `conftest.py` for shared fixtures

#### Test Structure:
```
service/
├── app/
│   ├── main.py
│   ├── agents/
│   └── providers/
└── tests/
    ├── unit/
    │   ├── test_agents.py          ✅ Unit tests
    │   ├── test_providers.py       ✅ Unit tests
    │   └── test_api.py             ✅ API endpoint tests
    ├── integration/
    │   └── test_api.py             ✅ Integration tests
    ├── conftest.py                 ✅ Shared fixtures
    └── fixtures/
        └── test_messages.json      ✅ Test data
```

#### Test Patterns:
- **Pytest fixtures**: `@pytest.fixture` for test dependencies
- **Mocking**: `@patch` decorator for external dependencies
- **FastAPI TestClient**: `TestClient` for API endpoint testing
- **Async testing**: `pytest-asyncio` for async operations
- **Markers**: `@pytest.mark.unit`, `@pytest.mark.integration`

#### Example Test:
```python
class TestHealthCheck:
    """Tests for GET /healthz endpoint."""
    
    def test_health_check(self, client):
        """Test health check endpoint."""
        response = client.get("/healthz")
        assert response.status_code == 200
        data = response.json()
        assert "status" in data
        assert "checks" in data

class TestChatEndpoint:
    """Tests for POST /v1/chat endpoint."""
    
    @patch('app.main.make_model')
    @patch('app.main.build_agent')
    def test_chat_success(self, mock_build_agent, mock_make_model, client, mock_chat_model):
        """Test successful chat completion."""
        mock_make_model.return_value = mock_chat_model
        # ... test logic
```

#### Running Tests:
```bash
# All tests
pytest

# Unit tests only
pytest tests/unit/ -v

# Integration tests
pytest tests/integration/ -v

# With coverage
pytest --cov=app --cov-report=html --cov-report=term

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

#### Coverage Status:
- **AI Service**: ✅ **85%** (30+ tests passing)
- **Creative Service**: ✅ **84%** (46 tests passing)

#### Special Features:
- **FastAPI TestClient**: Built-in test client for API testing
- **AsyncMock**: For async function mocking
- **Shared Fixtures**: `conftest.py` provides reusable test setup
- **Test Markers**: Mark tests for selective execution

---

### 3. **Resume Worker** (Python)
**Framework**: `pytest` + `pytest-mock`  
**Test Files**: `test_*.py` in `tests/unit/` and `tests/integration/`  
**Fixtures**: `conftest.py` for shared fixtures

#### Test Structure:
```
resume-worker/
├── src/
│   ├── resume_generator.py
│   ├── ai_service.py
│   ├── database.py
│   └── ...
└── tests/
    ├── unit/
    │   ├── test_resume_generator.py    ✅ 83% coverage
    │   ├── test_ai_service.py          ✅ 56% coverage
    │   ├── test_database.py            ✅ 78% coverage
    │   ├── test_keyword_extractor.py   ✅ 93% coverage
    │   ├── test_translation_helper.py  ✅ 61% coverage
    │   ├── test_logger.py              ✅ 100% coverage
    │   ├── test_main.py                ✅ 87% coverage
    │   └── test_queue_consumer.py     ✅ 74% coverage
    ├── integration/
    │   └── test_worker.py             ✅ Integration tests
    └── conftest.py                     ✅ Shared fixtures
```

#### Test Patterns:
- **Mocking external services**: AI service, database, RabbitMQ
- **Fixture-based setup**: Database fixtures, mock services
- **Integration tests**: Full worker flow with test containers

#### Example Test:
```python
class TestResumeGenerator:
    """Tests for ResumeGenerator class."""
    
    @patch('resume_generator.AIService')
    def test_generate_resume(self, mock_ai_service):
        """Test resume generation."""
        mock_ai_service.return_value.generate_resume_section.return_value = "<p>Profile</p>"
        generator = ResumeGenerator(mock_ai_service.return_value)
        result = generator.generate(user_id, job_description)
        assert result is not None
```

#### Running Tests:
```bash
# All tests
pytest

# Unit tests only
pytest tests/unit/ -v --cov=src

# Integration tests
pytest tests/integration/ -v

# With coverage
pytest --cov=src --cov-report=html --cov-report=term

# Using Docker
docker build -f Dockerfile.test -t resume-worker-test .
docker run --rm resume-worker-test

# Using Makefile
make test
make test-unit
make test-cov
```

#### Coverage Status:
- **Resume Worker**: ✅ **76.78%** (93 tests passing)
- **Coverage Breakdown**:
  - `logger.py`: 100% ✅
  - `keyword_extractor.py`: 93% ✅
  - `main.py`: 87% ✅
  - `resume_generator.py`: 83% ✅
  - `database.py`: 78% ✅
  - `queue_consumer.py`: 74% ✅
  - `ai_service.py`: 56%
  - `translation_helper.py`: 61%

#### Special Features:
- **Comprehensive coverage**: Most complete test suite
- **Mock-heavy**: Extensive use of mocks for external dependencies
- **Integration tests**: Full worker integration tests

---

### 4. **Job Application Worker** (Node.js)
**Framework**: Jest (`@jest/globals`)  
**Test Files**: `*.test.js` in `tests/__tests__/`  
**Configuration**: `jest.config.js` for ES modules

#### Test Structure:
```
job-application-worker/
├── src/
│   ├── worker.js
│   ├── orchestrator.js
│   ├── coverLetter.js
│   └── ...
└── tests/
    ├── __tests__/
    │   ├── health.test.js          ✅ 94.73% coverage
    │   ├── orchestrator.test.js    ✅ 100% coverage
    │   └── coverLetter.test.js    ✅ 41.33% coverage
    ├── setup.js                    ✅ Jest setup
    └── README.md
```

#### Test Patterns:
- **Jest mocks**: `jest.fn()`, `jest.unstable_mockModule()` for ES modules
- **Describe/Test blocks**: Organized test suites
- **Mock modules**: Module-level mocking for dependencies
- **Async testing**: Native async/await support

#### Example Test:
```javascript
import { describe, test, expect, beforeEach, jest } from '@jest/globals';
import { Orchestrator } from '../../src/orchestrator.js';

describe('Orchestrator', () => {
  let orchestrator;
  let mockDb;

  beforeEach(() => {
    mockDb = {
      getWebsiteByName: jest.fn(),
      resetWebsiteCount: jest.fn(),
      updateWebsiteCount: jest.fn(),
    };
    orchestrator = new Orchestrator(mockDb);
  });

  test('should return false when website is not found', async () => {
    mockDb.getWebsiteByName.mockResolvedValue(null);
    const result = await orchestrator.shouldProcessWebsite('unknown-site');
    expect(result).toBe(false);
  });
});
```

#### Running Tests:
```bash
# All tests
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

#### Coverage Status:
- **Job Application Worker**: 🔄 **15.46%** (23 tests passing)
- **Coverage Breakdown**:
  - `orchestrator.js`: 100% ✅
  - `health.js`: 94.73% ✅
  - `coverLetter.js`: 41.33%
  - `logger.js`: 63.26%
  - Other modules: 0% (not yet tested)

#### Special Features:
- **ES Module Support**: Jest configured for ES modules (`moduleNameMapper`)
- **Coverage Threshold**: 70% threshold in `jest.config.js`
- **Setup File**: `tests/setup.js` for global Jest availability

---

## Comparison Table

| Component | Language | Framework | Test Files | Coverage | Tests | Status |
|-----------|----------|-----------|------------|----------|-------|--------|
| **Server** | Go | `testing` | `*_test.go` | ⏳ Not tested | 0 | Pending |
| **Email Worker** | Go | `testing` | `*_test.go` | 49.7% | 30+ | Partial |
| **WhatsApp Worker** | Go | `testing` | `*_test.go` | 23.4% | 20+ | Partial |
| **Translation Worker** | Go | `testing` | `*_test.go` | 34.3% | 20+ | Partial |
| **Resume Worker** | Python | `pytest` | `test_*.py` | 76.78% | 93 | ✅ Complete |
| **Job App Worker** | Node.js | Jest | `*.test.js` | 15.46% | 23 | Partial |
| **AI Service** | Python | `pytest` | `test_*.py` | 85% | 30+ | ✅ Complete |
| **Creative Service** | Python | `pytest` | `test_*.py` | 84% | 46 | ✅ Complete |

---

## Key Differences

### 1. **Test Organization**
- **Go**: Tests alongside source (`*_test.go` next to `*.go`)
- **Python**: Tests in dedicated `tests/` directory (`test_*.py`)
- **Node.js**: Tests in `tests/__tests__/` directory (`*.test.js`)

### 2. **Test Framework**
- **Go**: Standard library `testing` package
- **Python**: `pytest` (third-party, feature-rich)
- **Node.js**: Jest (third-party, JavaScript-focused)

### 3. **Mocking Strategy**
- **Go**: Manual mocks (interfaces) or `testify/mock` (if used)
- **Python**: `unittest.mock` (`@patch` decorator)
- **Node.js**: Jest mocks (`jest.fn()`, `jest.unstable_mockModule()`)

### 4. **Test Execution**
- **Go**: `go test ./...` (parallel by default)
- **Python**: `pytest` (parallel with `pytest-xdist`)
- **Node.js**: `jest` (parallel by default)

### 5. **Coverage Tools**
- **Go**: Built-in `go tool cover`
- **Python**: `pytest-cov` (wrapper around `coverage.py`)
- **Node.js**: Jest built-in coverage

### 6. **Test Types**
- **Go**: Unit tests (integration tests use build tags)
- **Python**: Unit + Integration (marked with `@pytest.mark`)
- **Node.js**: Unit tests (integration tests in separate files)

---

## Common Test Patterns

### 1. **Table-Driven Tests** (Go)
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name string
        input string
        want string
    }{
        {"case1", "input1", "output1"},
        {"case2", "input2", "output2"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Function(tt.input)
            if got != tt.want {
                t.Errorf("Function() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 2. **Fixture-Based Tests** (Python)
```python
@pytest.fixture
def mock_service():
    """Mock service for testing."""
    return Mock()

def test_function(mock_service):
    """Test using fixture."""
    result = function(mock_service)
    assert result is not None
```

### 3. **Mock-Based Tests** (Node.js)
```javascript
describe('Function', () => {
  let mockService;
  
  beforeEach(() => {
    mockService = {
      method: jest.fn().mockResolvedValue('result'),
    };
  });
  
  test('should call service', async () => {
    await function(mockService);
    expect(mockService.method).toHaveBeenCalled();
  });
});
```

---

## Test Infrastructure

### Docker Test Images

All components have `Dockerfile.test` for consistent test environments:

**Go Components**:
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["go", "test", "./...", "-v", "-coverprofile=coverage.out"]
```

**Python Components**:
```dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt ./
RUN pip install -r requirements.txt
COPY src tests ./
CMD ["pytest", "-v", "--cov=src", "--cov-report=term"]
```

**Node.js Components**:
```dockerfile
FROM node:20-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY src tests ./
CMD ["npm", "run", "test:jest:coverage"]
```

### Makefile Commands

All components have standardized Makefile commands:

**Go**:
```makefile
test: go test ./... -v
test-unit: go test ./internal/... ./pkg/... -v
test-cov: go test ./... -v -coverprofile=coverage.out
test-cov-check: # Fails if coverage < 70%
```

**Python**:
```makefile
test: pytest -v
test-unit: pytest tests/unit/ -v
test-cov: pytest --cov=app --cov-report=html
```

**Node.js**:
```makefile
test: npm run test:jest
test-unit: npm run test:jest
test-cov: npm run test:jest:coverage
```

---

## Coverage Goals and Status

### Coverage Targets
- **Overall**: 70% minimum
- **Critical Components**: 90%+ (auth, payments, etc.)
- **Package/Utility Code**: 90%+ (logger, health, config)
- **Domain Logic**: 80%+
- **Integration**: Critical paths covered

### Current Status

| Component | Target | Current | Status |
|-----------|--------|---------|--------|
| **AI Service** | 70% | 85% | ✅ Exceeds |
| **Creative Service** | 70% | 84% | ✅ Exceeds |
| **Resume Worker** | 70% | 76.78% | ✅ Exceeds |
| **Email Worker** | 70% | 49.7% | 🔄 Needs work |
| **Translation Worker** | 70% | 34.3% | 🔄 Needs work |
| **WhatsApp Worker** | 70% | 23.4% | 🔄 Needs work |
| **Job App Worker** | 70% | 15.46% | 🔄 Needs work |
| **Server** | 70% | 0% | ⏳ Not started |

---

## Test Categories

### Unit Tests
**Purpose**: Test isolated functions/classes  
**Speed**: Fast (< 1 second each)  
**Dependencies**: Mocked  
**Examples**:
- Configuration loading
- Data transformation
- Business logic
- Utility functions

### Integration Tests
**Purpose**: Test component interactions  
**Speed**: Slower (< 30 seconds each)  
**Dependencies**: Real services (test containers)  
**Examples**:
- API endpoint tests
- Database operations
- RabbitMQ message flow
- Service-to-service communication

---

## Best Practices

### 1. **Test Organization**
- One test file per source file
- Group related tests (subtests in Go, classes in Python/Node.js)
- Use descriptive test names
- Keep tests isolated (no shared state)

### 2. **Test Data**
- Use fixtures for complex data
- Create factories for generating test entities
- Clean up after tests
- Use realistic data matching production

### 3. **Mocking**
- Mock external dependencies (databases, APIs, services)
- Use real implementations in integration tests
- Verify mock interactions
- Keep mocks simple

### 4. **Error Testing**
- Test error cases (failures, timeouts, invalid input)
- Test edge cases (empty data, null values, boundaries)
- Test retry logic
- Test graceful degradation

### 5. **Performance**
- Fast unit tests (< 1 second)
- Reasonable integration tests (< 30 seconds)
- Run tests in parallel when possible
- Skip slow tests in development

---

## Running Tests

### Local Development

**Go**:
```bash
cd backend/email-worker
make test              # All tests
make test-unit         # Unit tests only
make test-cov          # With coverage
make test-cov-check    # Check 70% threshold
```

**Python**:
```bash
cd backend/ai-service
make test              # All tests
make test-unit         # Unit tests only
make test-cov          # With coverage
```

**Node.js**:
```bash
cd backend/job-application-worker
make test              # All tests
make test-unit         # Unit tests only
make test-cov          # With coverage
```

### Docker

**All Components**:
```bash
# Build test image
docker build -f Dockerfile.test -t {component}-test .

# Run tests
docker run --rm {component}-test
```

### CI/CD Integration

Tests should run:
- On every push/PR
- Fail builds on test failures
- Generate coverage reports
- Track coverage trends
- Run integration tests in CI (with test containers)

---

## Test Infrastructure Needs

### Missing Infrastructure

1. **Server Tests**: No tests implemented yet
2. **Integration Test Containers**: Docker Compose for test dependencies
3. **Test Helpers**: Shared test utilities across components
4. **CI/CD Pipeline**: Automated test execution
5. **Coverage Tracking**: Historical coverage trends

### Recommended Additions

1. **Test Containers**: Use testcontainers for integration tests
2. **Shared Test Utilities**: Common helpers for database, RabbitMQ, etc.
3. **Performance Tests**: Load testing for critical endpoints
4. **E2E Tests**: Full system end-to-end tests
5. **Mutation Testing**: Verify test quality

---

## Summary

### Consistency
- ✅ All components have test infrastructure
- ✅ All use Docker for test environments
- ✅ All have Makefile commands
- ✅ All target 70% coverage minimum
- ✅ All support coverage reporting

### Differences
- **Framework**: Go (`testing`) vs Python (`pytest`) vs Node.js (`Jest`)
- **Organization**: Alongside source (Go) vs dedicated directory (Python/Node.js)
- **Mocking**: Manual (Go) vs decorators (Python) vs Jest (Node.js)
- **Coverage**: Built-in (Go) vs third-party (Python/Node.js)

### Status
- **Complete**: AI Service (85%), Creative Service (84%), Resume Worker (76.78%)
- **Partial**: Email Worker (49.7%), Translation Worker (34.3%), WhatsApp Worker (23.4%), Job App Worker (15.46%)
- **Pending**: Server (0%)

### Design Philosophy
1. **Fast Feedback**: Unit tests run quickly
2. **Comprehensive Coverage**: Target 70%+ coverage
3. **Isolation**: Tests don't depend on each other
4. **Realistic**: Test data matches production scenarios
5. **Maintainable**: Clear, descriptive test names
