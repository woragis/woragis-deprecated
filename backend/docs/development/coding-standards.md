# Coding Standards

**Last Updated:** 2025-12-22  
**Purpose:** Code style guidelines and best practices for the Woragis backend

---

## Overview

This document outlines coding standards for all languages used in the Woragis backend:
- **Go** - Main server and most workers
- **Python** - AI service, creative service, docs service, resume worker
- **JavaScript/TypeScript** - Job application worker

---

## General Principles

1. **Readability First** - Code should be easy to read and understand
2. **Consistency** - Follow language-specific conventions
3. **Documentation** - Document complex logic and public APIs
4. **Testing** - Write tests for all new features
5. **Error Handling** - Handle errors explicitly and gracefully
6. **Security** - Never commit secrets, validate inputs, use parameterized queries

---

## Go Standards

### Code Formatting

- Use `gofmt` for formatting (automatically applied by most editors)
- Run `golangci-lint` for linting
- Maximum line length: 100 characters (soft limit)

### Naming Conventions

```go
// Package names: lowercase, single word
package database

// Exported functions: PascalCase
func GetUser(ctx context.Context, id string) (*User, error)

// Unexported functions: camelCase
func validateUser(user *User) error

// Constants: PascalCase for exported, camelCase for unexported
const DefaultTimeout = 30 * time.Second
const maxRetries = 3

// Variables: camelCase
var userCount int
var dbConnection *sql.DB

// Types: PascalCase
type User struct {
    ID   string
    Name string
}

// Interfaces: PascalCase, often end with "er" if single method
type UserRepository interface {
    GetUser(ctx context.Context, id string) (*User, error)
}
```

### File Organization

```
service/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── handler/            # HTTP handlers
│   ├── service/            # Business logic
│   ├── repository/         # Data access
│   └── config/             # Configuration
├── pkg/                    # Public packages (if any)
├── go.mod
├── go.sum
└── README.md
```

### Error Handling

```go
// Always handle errors explicitly
user, err := repo.GetUser(ctx, id)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// Use errors.Is and errors.As for error checking
if errors.Is(err, sql.ErrNoRows) {
    return nil, ErrUserNotFound
}

// Wrap errors with context
if err != nil {
    return fmt.Errorf("processing user %s: %w", userID, err)
}
```

### Context Usage

```go
// Always accept context as first parameter
func GetUser(ctx context.Context, id string) (*User, error)

// Pass context through call chain
func (s *Service) ProcessUser(ctx context.Context, id string) error {
    user, err := s.repo.GetUser(ctx, id)
    // ...
}

// Use context for cancellation and timeouts
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

### Struct Tags

```go
type User struct {
    ID        string    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Email     string    `json:"email" db:"email"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

### Comments

```go
// Package comment
// Package database provides database connection and query utilities.
package database

// Function comment for exported functions
// GetUser retrieves a user by ID from the database.
// It returns an error if the user is not found or if the database query fails.
func GetUser(ctx context.Context, id string) (*User, error) {
    // Implementation
}

// Inline comments for complex logic
// Calculate the exponential backoff delay: baseDelay * 2^attempt
delay := baseDelay * time.Duration(1<<attempt)
```

### Testing

```go
// Test file: user_test.go
func TestGetUser(t *testing.T) {
    tests := []struct {
        name    string
        userID  string
        want    *User
        wantErr bool
    }{
        {
            name:    "valid user",
            userID:  "123",
            want:    &User{ID: "123", Name: "John"},
            wantErr: false,
        },
        {
            name:    "not found",
            userID:  "999",
            want:    nil,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := GetUser(context.Background(), tt.userID)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetUser() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## Python Standards

### Code Formatting

- Use `black` for formatting (line length: 100)
- Use `isort` for import sorting
- Use `pylint` or `flake8` for linting
- Follow [PEP 8](https://pep8.org/) style guide

### Naming Conventions

```python
# Module names: lowercase, underscores
import user_service
import database_connection

# Class names: PascalCase
class UserService:
    pass

class DatabaseConnection:
    pass

# Function names: snake_case
def get_user(user_id: str) -> User:
    pass

def process_translation(text: str) -> str:
    pass

# Constants: UPPER_SNAKE_CASE
DEFAULT_TIMEOUT = 30
MAX_RETRIES = 3

# Variables: snake_case
user_count = 0
db_connection = None

# Private functions/variables: leading underscore
def _validate_user(user: User) -> bool:
    pass

_private_variable = "hidden"
```

### Type Hints

```python
from typing import Optional, List, Dict

def get_user(user_id: str) -> Optional[User]:
    """Get user by ID."""
    pass

def process_users(users: List[User]) -> Dict[str, int]:
    """Process users and return statistics."""
    pass

# Use Optional for nullable types
def find_user(name: str) -> Optional[User]:
    pass
```

### Docstrings

```python
def get_user(user_id: str) -> Optional[User]:
    """
    Retrieve a user by their unique identifier.
    
    Args:
        user_id: The unique identifier for the user
        
    Returns:
        User object if found, None otherwise
        
    Raises:
        ValueError: If user_id is empty or invalid
        DatabaseError: If database query fails
    """
    pass

class UserService:
    """Service for user-related operations."""
    
    def __init__(self, db: Database):
        """
        Initialize UserService.
        
        Args:
            db: Database connection instance
        """
        self.db = db
```

### Error Handling

```python
# Use specific exceptions
try:
    user = db.get_user(user_id)
except DatabaseError as e:
    logger.error(f"Database error: {e}")
    raise
except ValueError as e:
    logger.warning(f"Invalid user_id: {e}")
    raise

# Create custom exceptions
class UserNotFoundError(Exception):
    """Raised when user is not found."""
    pass

# Use context managers for resources
with open(file_path, 'r') as f:
    content = f.read()
```

### Imports

```python
# Standard library imports
import os
import sys
from typing import Optional, List

# Third-party imports
import fastapi
from pydantic import BaseModel

# Local imports
from app.models import User
from app.services import user_service
```

### Testing

```python
import pytest
from unittest.mock import Mock, patch

def test_get_user():
    """Test retrieving a valid user."""
    user = get_user("123")
    assert user.id == "123"
    assert user.name is not None

def test_get_user_not_found():
    """Test retrieving a non-existent user."""
    with pytest.raises(UserNotFoundError):
        get_user("999")

@pytest.fixture
def mock_db():
    """Fixture for mock database."""
    db = Mock()
    db.get_user.return_value = User(id="123", name="John")
    return db

def test_get_user_with_fixture(mock_db):
    """Test with mock database."""
    user = get_user("123", db=mock_db)
    assert user.name == "John"
```

---

## JavaScript/TypeScript Standards

### Code Formatting

- Use `prettier` for formatting
- Use `eslint` for linting
- Follow [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)

### Naming Conventions

```javascript
// Variables and functions: camelCase
const userName = "John";
function getUserById(id) {}

// Constants: UPPER_SNAKE_CASE
const DEFAULT_TIMEOUT = 30000;
const MAX_RETRIES = 3;

// Classes: PascalCase
class UserService {
    constructor() {}
}

// Private methods: leading underscore
class UserService {
    _validateUser(user) {}
}
```

### TypeScript Types

```typescript
// Use interfaces for object shapes
interface User {
    id: string;
    name: string;
    email: string;
}

// Use types for unions, intersections, etc.
type UserId = string;
type UserStatus = 'active' | 'inactive' | 'pending';

// Function types
type GetUser = (id: string) => Promise<User>;

// Generic types
interface Repository<T> {
    findById(id: string): Promise<T | null>;
}
```

### Error Handling

```typescript
// Use try-catch for async operations
async function getUser(id: string): Promise<User> {
    try {
        const user = await db.getUser(id);
        if (!user) {
            throw new UserNotFoundError(`User ${id} not found`);
        }
        return user;
    } catch (error) {
        logger.error(`Failed to get user ${id}:`, error);
        throw error;
    }
}

// Custom error classes
class UserNotFoundError extends Error {
    constructor(message: string) {
        super(message);
        this.name = 'UserNotFoundError';
    }
}
```

### Async/Await

```typescript
// Prefer async/await over promises
async function processUsers(userIds: string[]): Promise<User[]> {
    const users = await Promise.all(
        userIds.map(id => getUser(id))
    );
    return users;
}

// Handle errors properly
async function processUser(id: string): Promise<User | null> {
    try {
        return await getUser(id);
    } catch (error) {
        logger.error(`Error processing user ${id}:`, error);
        return null;
    }
}
```

### Testing

```typescript
import { describe, it, expect, beforeEach } from 'vitest';

describe('UserService', () => {
    let userService: UserService;
    
    beforeEach(() => {
        userService = new UserService(mockDb);
    });
    
    it('should get user by id', async () => {
        const user = await userService.getUser('123');
        expect(user).toBeDefined();
        expect(user.id).toBe('123');
    });
    
    it('should throw error for invalid id', async () => {
        await expect(userService.getUser('')).rejects.toThrow();
    });
});
```

---

## Common Patterns

### Logging

```go
// Go
logger.Info("Processing user", "user_id", userID, "action", "update")
logger.Error("Failed to process user", "error", err, "user_id", userID)
```

```python
# Python
logger.info("Processing user", extra={"user_id": user_id, "action": "update"})
logger.error("Failed to process user", exc_info=True, extra={"user_id": user_id})
```

```typescript
// TypeScript
logger.info('Processing user', { userId, action: 'update' });
logger.error('Failed to process user', { error, userId });
```

### Configuration

```go
// Go - Use environment variables
port := os.Getenv("APP_PORT")
if port == "" {
    port = "8080"
}
```

```python
# Python - Use environment variables
import os
port = os.getenv("APP_PORT", "8080")
```

```typescript
// TypeScript
const port = process.env.APP_PORT || '8080';
```

### Database Queries

```go
// Go - Use parameterized queries
query := "SELECT * FROM users WHERE id = $1"
row := db.QueryRow(ctx, query, userID)
```

```python
# Python - Use parameterized queries
query = "SELECT * FROM users WHERE id = %s"
cursor.execute(query, (user_id,))
```

```typescript
// TypeScript - Use parameterized queries
const query = 'SELECT * FROM users WHERE id = $1';
const result = await db.query(query, [userId]);
```

---

## Security Best Practices

1. **Never commit secrets** - Use environment variables or secrets management
2. **Validate all inputs** - Check user input before processing
3. **Use parameterized queries** - Prevent SQL injection
4. **Sanitize output** - Escape HTML/JavaScript in user-generated content
5. **Use HTTPS** - Encrypt data in transit
6. **Implement rate limiting** - Prevent abuse
7. **Keep dependencies updated** - Patch security vulnerabilities
8. **Use secure defaults** - Secure by default

---

## Code Review Checklist

- [ ] Code follows language-specific style guide
- [ ] Functions are well-documented
- [ ] Error handling is appropriate
- [ ] Tests are included and passing
- [ ] No secrets or sensitive data committed
- [ ] Input validation is present
- [ ] Logging is appropriate (not too verbose, not too sparse)
- [ ] Performance considerations addressed
- [ ] Security best practices followed

---

## Tools

### Go
- `gofmt` - Formatting
- `golangci-lint` - Linting
- `go vet` - Static analysis
- `go test` - Testing

### Python
- `black` - Formatting
- `isort` - Import sorting
- `pylint` or `flake8` - Linting
- `pytest` - Testing
- `mypy` - Type checking

### JavaScript/TypeScript
- `prettier` - Formatting
- `eslint` - Linting
- `vitest` or `jest` - Testing
- `typescript` - Type checking

---

## Related Documentation

- [Development Setup Guide](./setup-guide.md)
- [Contributing Guide](./contributing.md)
- [Testing Patterns](./testing-patterns.md)
- [Error Handling](./error-handling.md)
- [Logging Conventions](./logging-conventions.md)

---

**Last Updated:** 2025-12-22
