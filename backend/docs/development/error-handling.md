# Error Handling Patterns

## Overview

This guide covers error handling patterns and best practices across all backend components. Consistent error handling improves debugging and user experience.

## Error Handling Principles

1. **Fail Fast**: Detect errors early
2. **Fail Explicitly**: Return clear error messages
3. **Fail Gracefully**: Handle errors appropriately
4. **Log Errors**: Always log errors with context
5. **Propagate Appropriately**: Don't swallow errors

## Error Types

### Transient Errors
- Network errors
- Temporary service unavailability
- Timeout errors

**Action**: Retry with exponential backoff

### Permanent Errors
- Invalid input
- Authentication failures
- Resource not found

**Action**: Don't retry, return error immediately

### System Errors
- Database connection failures
- Configuration errors
- Critical service failures

**Action**: Log and alert, may require service restart

## Error Handling by Language

### Go

#### Error Creation

```go
import "errors"
import "fmt"

// Simple error
err := errors.New("entity not found")

// Formatted error
err := fmt.Errorf("failed to create entity: %w", originalErr)

// Custom error type
type NotFoundError struct {
    EntityType string
    EntityID   uuid.UUID
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with ID %s not found", e.EntityType, e.EntityID)
}
```

#### Error Wrapping

```go
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}
```

#### Error Checking

```go
result, err := operation()
if err != nil {
    // Handle error
    return err
}
// Use result
```

#### Error Propagation

```go
func (s *service) Process(ctx context.Context, input Input) error {
    entity, err := s.repo.GetByID(ctx, input.ID)
    if err != nil {
        return fmt.Errorf("failed to get entity: %w", err)
    }
    
    if err := s.validate(entity); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    return s.repo.Update(ctx, entity)
}
```

### Python

#### Error Creation

```python
# Simple error
raise ValueError("Invalid input")

# Custom error
class NotFoundError(Exception):
    def __init__(self, entity_type: str, entity_id: str):
        self.entity_type = entity_type
        self.entity_id = entity_id
        super().__init__(f"{entity_type} with ID {entity_id} not found")
```

#### Error Handling

```python
try:
    result = operation()
except SpecificError as e:
    logger.error("Operation failed", error=str(e))
    raise
except Exception as e:
    logger.exception("Unexpected error", exc_info=True)
    raise
```

#### Error Propagation

```python
def process(input_data):
    try:
        entity = repo.get_by_id(input_data.id)
    except NotFoundError:
        raise
    except Exception as e:
        raise ProcessingError(f"Failed to get entity: {e}") from e
    
    try:
        validate(entity)
    except ValidationError as e:
        raise ProcessingError(f"Validation failed: {e}") from e
    
    return repo.update(entity)
```

### Node.js

#### Error Creation

```javascript
// Simple error
throw new Error('Entity not found');

// Custom error
class NotFoundError extends Error {
    constructor(entityType, entityId) {
        super(`${entityType} with ID ${entityId} not found`);
        this.entityType = entityType;
        this.entityId = entityId;
        this.name = 'NotFoundError';
    }
}
```

#### Error Handling

```javascript
try {
    const result = await operation();
} catch (error) {
    if (error instanceof NotFoundError) {
        logger.error('Entity not found', { error: error.message });
        throw error;
    } else {
        logger.error('Unexpected error', { error: error.message, stack: error.stack });
        throw new Error('Operation failed');
    }
}
```

## Error Handling Patterns

### Repository Layer

**Pattern**: Return errors, don't handle them

```go
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Entity, error) {
    var entity Entity
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, &NotFoundError{EntityID: id}
        }
        return nil, fmt.Errorf("database error: %w", err)
    }
    return &entity, nil
}
```

### Service Layer

**Pattern**: Handle business logic errors, wrap database errors

```go
func (s *service) GetEntity(ctx context.Context, id uuid.UUID) (*Entity, error) {
    entity, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, &NotFoundError{}) {
            return nil, err // Propagate as-is
        }
        s.logger.Error("failed to get entity", "error", err, "id", id)
        return nil, fmt.Errorf("failed to get entity: %w", err)
    }
    
    // Business logic validation
    if !entity.IsActive() {
        return nil, &InactiveEntityError{EntityID: id}
    }
    
    return entity, nil
}
```

### Handler Layer

**Pattern**: Convert errors to HTTP responses

```go
func (h *handler) GetEntity(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.BadRequest(c, "invalid entity ID", err)
    }
    
    entity, err := h.service.GetEntity(c.Context(), id)
    if err != nil {
        switch {
        case errors.Is(err, &NotFoundError{}):
            return response.NotFound(c, "entity not found", err)
        case errors.Is(err, &InactiveEntityError{}):
            return response.Forbidden(c, "entity is inactive", err)
        default:
            h.logger.Error("failed to get entity", "error", err, "id", id)
            return response.InternalError(c, "failed to get entity", err)
        }
    }
    
    return response.OK(c, entity)
}
```

## HTTP Error Responses

### Standard Error Format

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {
    "field": "Additional details"
  }
}
```

### Status Code Mapping

- `400 Bad Request` - Invalid input, validation errors
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server errors
- `503 Service Unavailable` - Service unavailable

### Error Response Helpers

**Go (Fiber)**:
```go
package response

func BadRequest(c *fiber.Ctx, message string, err error) error {
    return c.Status(400).JSON(fiber.Map{
        "error": message,
        "code": "BAD_REQUEST",
        "details": err.Error(),
    })
}

func NotFound(c *fiber.Ctx, message string, err error) error {
    return c.Status(404).JSON(fiber.Map{
        "error": message,
        "code": "NOT_FOUND",
    })
}

func InternalError(c *fiber.Ctx, message string, err error) error {
    return c.Status(500).JSON(fiber.Map{
        "error": message,
        "code": "INTERNAL_ERROR",
    })
}
```

**Python (FastAPI)**:
```python
from fastapi import HTTPException

def bad_request(message: str, details: str = None):
    raise HTTPException(
        status_code=400,
        detail={
            "error": message,
            "code": "BAD_REQUEST",
            "details": details
        }
    )

def not_found(message: str):
    raise HTTPException(
        status_code=404,
        detail={
            "error": message,
            "code": "NOT_FOUND"
        }
    )
```

## Retry Logic

### Exponential Backoff

**Go**:
```go
func retryWithBackoff(ctx context.Context, maxRetries int, operation func() error) error {
    var lastErr error
    delay := time.Second
    
    for i := 0; i < maxRetries; i++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        lastErr = err
        if !isTransientError(err) {
            return err // Don't retry permanent errors
        }
        
        if i < maxRetries-1 {
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(delay):
                delay *= 2 // Exponential backoff
            }
        }
    }
    
    return fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

**Python**:
```python
import time
from tenacity import retry, stop_after_attempt, wait_exponential

@retry(
    stop=stop_after_attempt(3),
    wait=wait_exponential(multiplier=1, min=1, max=10),
    retry=retry_if_transient_error
)
def operation_with_retry():
    # Operation that may fail
    pass
```

## Error Logging

### Always Log Errors

```go
logger.Error("Operation failed",
    slog.String("operation", "create_entity"),
    slog.String("user_id", userID.String()),
    slog.String("error", err.Error()))
```

### Include Context

```go
logger.Error("Failed to process job",
    slog.String("job_id", jobID.String()),
    slog.String("job_type", jobType),
    slog.Int("retry_count", retryCount),
    slog.String("error", err.Error()))
```

### Log Stack Traces (when available)

**Python**:
```python
logger.exception("Operation failed",
    operation="create_entity",
    entity_id=entity_id,
    exc_info=True)
```

## Best Practices

1. **Error Messages**:
   - Be clear and specific
   - Don't expose internal details to users
   - Include actionable information

2. **Error Codes**:
   - Use consistent error codes
   - Map errors to appropriate HTTP status codes
   - Document error codes

3. **Error Propagation**:
   - Wrap errors with context
   - Don't swallow errors
   - Return errors, don't panic (unless fatal)

4. **Error Handling**:
   - Handle errors at the right level
   - Don't handle errors you can't handle
   - Log errors before returning

5. **Retry Logic**:
   - Only retry transient errors
   - Use exponential backoff
   - Limit retry count

6. **Validation**:
   - Validate input early
   - Return clear validation errors
   - Don't proceed with invalid data

## Related Documentation

- [Adding a Domain](./adding-domain.md) - Domain error handling examples
- [Adding a Worker](./adding-worker.md) - Worker error handling examples
- [Adding a Service](./adding-service.md) - Service error handling examples
- [Resilience Overview](../architecture/RESILIENCE_OVERVIEW.md) - Resilience patterns
