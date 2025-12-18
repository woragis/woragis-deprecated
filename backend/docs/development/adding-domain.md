# Adding a New Domain

## Overview

This guide explains how to add a new domain to the Woragis server. A domain represents a business entity (e.g., Projects, Resumes, Certifications) and follows a consistent structure.

## Domain Structure

Each domain follows this structure:

```
domains/{domain-name}/
├── model.go          # Database models (GORM)
├── repository.go     # Database operations
├── service.go        # Business logic
├── handler.go        # HTTP handlers
└── routes.go         # Route definitions
```

## Step-by-Step Guide

### 1. Create Domain Directory

```bash
mkdir -p server/app/internal/domains/{domain-name}
cd server/app/internal/domains/{domain-name}
```

### 2. Define Model

Create `model.go`:

```go
package {domainname}

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type {Entity} struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
    Name      string    `gorm:"not null"`
    // Add other fields...
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name
func ({Entity}) TableName() string {
    return "{entity}s"
}
```

### 3. Create Repository

Create `repository.go`:

```go
package {domainname}

import (
    "context"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Repository interface {
    Create(ctx context.Context, entity *{Entity}) error
    GetByID(ctx context.Context, id uuid.UUID) (*{Entity}, error)
    ListByUserID(ctx context.Context, userID uuid.UUID) ([]*{Entity}, error)
    Update(ctx context.Context, entity *{Entity}) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, entity *{Entity}) error {
    return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*{Entity}, error) {
    var entity {Entity}
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error
    if err != nil {
        return nil, err
    }
    return &entity, nil
}

// Implement other methods...
```

### 4. Create Service

Create `service.go`:

```go
package {domainname}

import (
    "context"
    "github.com/google/uuid"
    "log/slog"
)

type Service interface {
    Create{Entity}(ctx context.Context, userID uuid.UUID, input Create{Entity}Input) (*{Entity}, error)
    Get{Entity}(ctx context.Context, id uuid.UUID) (*{Entity}, error)
    List{Entity}s(ctx context.Context, userID uuid.UUID) ([]*{Entity}, error)
    Update{Entity}(ctx context.Context, id uuid.UUID, input Update{Entity}Input) (*{Entity}, error)
    Delete{Entity}(ctx context.Context, id uuid.UUID) error
}

type Create{Entity}Input struct {
    Name string `json:"name" validate:"required"`
    // Add other fields...
}

type Update{Entity}Input struct {
    Name *string `json:"name"`
    // Add other fields...
}

type service struct {
    repo   Repository
    logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
    return &service{
        repo:   repo,
        logger: logger,
    }
}

func (s *service) Create{Entity}(ctx context.Context, userID uuid.UUID, input Create{Entity}Input) (*{Entity}, error) {
    entity := &{Entity}{
        ID:     uuid.New(),
        UserID: userID,
        Name:   input.Name,
    }
    
    if err := s.repo.Create(ctx, entity); err != nil {
        s.logger.Error("failed to create entity", "error", err)
        return nil, err
    }
    
    return entity, nil
}

// Implement other methods...
```

### 5. Create Handler

Create `handler.go`:

```go
package {domainname}

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "log/slog"
    "github.com/woragis/backend/server/app/pkg/response"
)

type Handler interface {
    Create{Entity}(c *fiber.Ctx) error
    Get{Entity}(c *fiber.Ctx) error
    List{Entity}s(c *fiber.Ctx) error
    Update{Entity}(c *fiber.Ctx) error
    Delete{Entity}(c *fiber.Ctx) error
}

type handler struct {
    service Service
    logger  *slog.Logger
}

func NewHandler(service Service, logger *slog.Logger) Handler {
    return &handler{
        service: service,
        logger:  logger,
    }
}

func (h *handler) Create{Entity}(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)
    
    var input Create{Entity}Input
    if err := c.BodyParser(&input); err != nil {
        return response.BadRequest(c, "invalid request body", err)
    }
    
    entity, err := h.service.Create{Entity}(c.Context(), userID, input)
    if err != nil {
        h.logger.Error("failed to create entity", "error", err)
        return response.InternalError(c, "failed to create entity", err)
    }
    
    return response.Created(c, entity)
}

// Implement other handlers...
```

### 6. Create Routes

Create `routes.go`:

```go
package {domainname}

import "github.com/gofiber/fiber/v2"

func SetupRoutes(api fiber.Router, handler Handler) {
    api.Post("/", handler.Create{Entity})
    api.Get("/", handler.List{Entity}s)
    api.Get("/:id", handler.Get{Entity})
    api.Patch("/:id", handler.Update{Entity})
    api.Delete("/:id", handler.Delete{Entity})
}
```

### 7. Register Domain in Main

In `server/app/cmd/server/main.go`:

```go
// Import domain
{domainname}domain "github.com/woragis/backend/server/app/internal/domains/{domainname}"

// Initialize domain
{domainname}Repo := {domainname}domain.NewRepository(db)
{domainname}Service := {domainname}domain.NewService({domainname}Repo, slogLogger)
{domainname}Handler := {domainname}domain.NewHandler({domainname}Service, slogLogger)

// Register routes
api := app.Group("/api/{domain-name}")
{domainname}domain.SetupRoutes(api, {domainname}Handler)
```

### 8. Run Migrations

The domain will be auto-migrated on server startup (development). For production, create a migration:

```bash
# Create migration file
# migrations/YYYYMMDDHHMMSS_create_{entity}s.sql
```

## Testing

### Unit Tests

Create `repository_test.go`, `service_test.go`, `handler_test.go`:

```go
package {domainname}

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestCreate{Entity}(t *testing.T) {
    // Test implementation
}
```

### Integration Tests

Create integration tests in `tests/integration/`:

```go
package integration

import (
    "testing"
    // Test implementation
)
```

## Best Practices

1. **Follow Naming Conventions**:
   - Domain package: lowercase, singular (e.g., `projects`)
   - Entity: PascalCase, singular (e.g., `Project`)
   - Table name: lowercase, plural (e.g., `projects`)

2. **Use Interfaces**:
   - Define interfaces for Repository, Service, Handler
   - Makes testing easier (can mock interfaces)

3. **Error Handling**:
   - Return domain-specific errors
   - Log errors with context
   - Use response helpers for HTTP responses

4. **Validation**:
   - Validate input in handlers
   - Use struct tags for validation
   - Return clear validation errors

5. **Logging**:
   - Log important operations
   - Include context (user ID, entity ID, etc.)
   - Use appropriate log levels

6. **Database**:
   - Use transactions for multi-step operations
   - Handle soft deletes (if needed)
   - Add indexes for frequently queried fields

## Example: Complete Domain

See existing domains for complete examples:
- `domains/projects/` - Complex domain with multiple features
- `domains/skills/` - Simpler domain
- `domains/certifications/` - Domain with relationships

## Related Documentation

- [Testing Patterns](./testing-patterns.md) - How to test domains
- [Error Handling Patterns](./error-handling.md) - Error handling best practices
- [Logging Conventions](./logging-conventions.md) - Logging guidelines
