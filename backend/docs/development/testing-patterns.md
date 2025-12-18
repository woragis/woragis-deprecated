# Testing Patterns

## Overview

This guide covers testing patterns and best practices for the Woragis backend. Each component (server, workers, services) has specific testing requirements.

## Testing Strategy

### Test Types

1. **Unit Tests**: Test individual functions/methods in isolation
2. **Integration Tests**: Test component interactions (database, queues, external APIs)
3. **End-to-End Tests**: Test complete workflows

### Coverage Goals

- **Unit Tests**: 70%+ coverage
- **Integration Tests**: Critical paths covered
- **End-to-End Tests**: Main user flows covered

## Go Testing

### Unit Tests

**File**: `{package}_test.go`

**Pattern**:
```go
package {package}

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestFunction(t *testing.T) {
    // Arrange
    input := "test"
    expected := "result"
    
    // Act
    result := Function(input)
    
    // Assert
    assert.Equal(t, expected, result)
}
```

### Table-Driven Tests

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "success case",
            input:    "test",
            expected: "result",
            wantErr:  false,
        },
        {
            name:     "error case",
            input:    "",
            expected: "",
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Function(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}
```

### Mocking

**Using testify/mock**:
```go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*Entity, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*Entity), args.Error(1)
}

func TestService(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("GetByID", mock.Anything, mock.Anything).Return(&Entity{}, nil)
    
    service := NewService(mockRepo, logger)
    // Test service...
}
```

### Integration Tests

**File**: `{package}_integration_test.go`

**Pattern**:
```go
// +build integration

package {package}

import (
    "testing"
    "github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    // Test with real database
    repo := NewRepository(db)
    // Test...
}
```

**Running Integration Tests**:
```bash
go test -tags=integration ./...
```

## Python Testing

### Unit Tests

**File**: `tests/unit/test_{module}.py`

**Pattern**:
```python
import pytest
from unittest.mock import Mock, patch

def test_function():
    # Arrange
    input_value = "test"
    expected = "result"
    
    # Act
    result = function(input_value)
    
    # Assert
    assert result == expected
```

### Pytest Fixtures

```python
import pytest
from app.database import get_db

@pytest.fixture
def db_session():
    # Setup test database
    db = get_db()
    yield db
    # Teardown
    db.close()

def test_with_db(db_session):
    # Use db_session
    pass
```

### Mocking

```python
from unittest.mock import Mock, patch

def test_with_mock():
    mock_service = Mock()
    mock_service.get_data.return_value = {"key": "value"}
    
    result = function_under_test(mock_service)
    assert result == expected
```

### Integration Tests

**File**: `tests/integration/test_{module}.py`

**Pattern**:
```python
import pytest

@pytest.mark.integration
def test_integration():
    # Test with real dependencies
    pass
```

**Running Integration Tests**:
```bash
pytest -m integration
```

## Node.js Testing

### Unit Tests

**File**: `tests/unit/{module}.test.js`

**Pattern**:
```javascript
const { describe, it, expect, jest } = require('@jest/globals');

describe('Function', () => {
    it('should return expected result', () => {
        // Arrange
        const input = 'test';
        const expected = 'result';
        
        // Act
        const result = function(input);
        
        // Assert
        expect(result).toBe(expected);
    });
});
```

### Mocking

```javascript
jest.mock('./dependency', () => ({
    dependencyFunction: jest.fn(() => 'mocked value')
}));
```

### Integration Tests

**File**: `tests/integration/{module}.test.js`

**Pattern**:
```javascript
describe('Integration Test', () => {
    it('should work with real dependencies', async () => {
        // Test with real dependencies
    });
});
```

## Testing Patterns by Component

### Server (Go)

#### Repository Tests

```go
func TestRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := NewRepository(db)
    
    entity := &Entity{Name: "Test"}
    err := repo.Create(context.Background(), entity)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, entity.ID)
}
```

#### Service Tests

```go
func TestService_CreateEntity(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    service := NewService(mockRepo, logger)
    result, err := service.CreateEntity(ctx, userID, input)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

#### Handler Tests

```go
func TestHandler_CreateEntity(t *testing.T) {
    mockService := new(MockService)
    mockService.On("CreateEntity", mock.Anything, mock.Anything, mock.Anything).Return(&Entity{}, nil)
    
    handler := NewHandler(mockService, logger)
    
    app := fiber.New()
    app.Post("/", handler.CreateEntity)
    
    req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"Test"}`))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    
    assert.Equal(t, 201, resp.StatusCode)
}
```

### Workers

#### Job Processing Tests

```go
func TestProcessJob(t *testing.T) {
    job := &Job{ID: uuid.New(), Data: "test"}
    
    err := processJob(job)
    
    assert.NoError(t, err)
    // Verify job was processed correctly
}
```

#### Queue Consumption Tests

```go
func TestConsumeQueue(t *testing.T) {
    // Setup test RabbitMQ connection
    conn := setupTestRabbitMQ(t)
    defer conn.Close()
    
    // Publish test message
    publishTestMessage(t, conn, queueName, testMessage)
    
    // Consume and verify
    // ...
}
```

### Services (Python/FastAPI)

#### Endpoint Tests

```python
from fastapi.testclient import TestClient
from app.main import app

client = TestClient(app)

def test_endpoint():
    response = client.post("/v1/endpoint", json={"field": "value"})
    assert response.status_code == 200
    assert "result" in response.json()
```

#### Service Logic Tests

```python
def test_service_logic():
    mock_provider = Mock()
    mock_provider.get_data.return_value = {"key": "value"}
    
    service = Service(mock_provider)
    result = service.process("input")
    
    assert result == expected
```

## Test Data Management

### Fixtures

**Go**:
```go
func createTestEntity(t *testing.T) *Entity {
    return &Entity{
        ID:     uuid.New(),
        Name:   "Test Entity",
        UserID: uuid.New(),
    }
}
```

**Python**:
```python
@pytest.fixture
def test_entity():
    return Entity(
        id=uuid.uuid4(),
        name="Test Entity",
        user_id=uuid.uuid4()
    )
```

### Test Database

**Setup**:
```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    
    // Run migrations
    db.AutoMigrate(&Entity{})
    
    return db
}
```

## Best Practices

1. **Test Naming**:
   - Use descriptive test names
   - Follow pattern: `Test{Function}_{Scenario}`

2. **Test Structure**:
   - Arrange-Act-Assert pattern
   - One assertion per test (when possible)

3. **Test Isolation**:
   - Tests should not depend on each other
   - Clean up test data after each test

4. **Mocking**:
   - Mock external dependencies
   - Don't mock code you own (unless necessary)

5. **Coverage**:
   - Aim for 70%+ coverage
   - Focus on critical paths
   - Don't obsess over 100% coverage

6. **Performance**:
   - Keep tests fast
   - Use test databases (in-memory when possible)
   - Parallelize tests when possible

## Running Tests

### Go

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Python

```bash
# Run all tests
pytest

# Run with coverage
pytest --cov=app --cov-report=html

# Run specific test
pytest tests/unit/test_module.py::test_function
```

### Node.js

```bash
# Run all tests
npm test

# Run with coverage
npm test -- --coverage

# Run specific test
npm test -- tests/unit/test_module.test.js
```

## Related Documentation

- [Testing Overview](../architecture/TESTING_OVERVIEW.md) - Comprehensive testing guide
- [Adding a Domain](./adding-domain.md) - Domain testing examples
- [Adding a Worker](./adding-worker.md) - Worker testing examples
- [Adding a Service](./adding-service.md) - Service testing examples
