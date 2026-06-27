# Database Package

This package provides a unified interface for managing PostgreSQL and Redis database connections in the auth-service application.

## Features

- **PostgreSQL Integration**: GORM-based PostgreSQL connection with connection pooling
- **Redis Integration**: Redis client with proper error handling and health checks
- **Connection Management**: Centralized database manager for all connections
- **Health Checks**: Built-in health check functionality for monitoring
- **Configuration**: Environment-based configuration management
- **Graceful Shutdown**: Proper connection cleanup on application shutdown

## Usage

### Basic Setup

```go
package main

import (
    "log"
    "auth-service/internal/config"
    "auth-service/internal/database"
)

func main() {
    // Load configuration
    cfg := config.Load()

    // Create database manager
    dbManager, err := database.NewFromConfig(cfg)
    if err != nil {
        log.Fatalf("Failed to initialize database manager: %v", err)
    }
    defer dbManager.Close()

    // Use the connections
    postgresDB := dbManager.GetPostgres()
    redisClient := dbManager.GetRedis()
}
```

### Environment Variables

Create a `.env` file in your project root with the following variables:

```env
# Database
DATABASE_URL=postgres://username:password@localhost:5432/health_center?sslmode=disable
DATABASE_MAX_OPEN_CONNS=25
DATABASE_MAX_IDLE_CONNS=25
DATABASE_MAX_IDLE_TIME=15m

# Redis
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Security
JWT_SECRET=your-jwt-secret-key
AES_KEY=your-aes-key
HASH_SALT=your-hash-salt
BCRYPT_COST=12

# S3
S3_BUCKET_NAME=your-s3-bucket-name
```

### PostgreSQL Operations

```go
// Get PostgreSQL connection
db := dbManager.GetPostgres()

// Auto-migrate models
db.AutoMigrate(&YourModel{})

// Create a record
user := &User{Name: "John", Email: "john@example.com"}
db.Create(user)

// Find records
var users []User
db.Find(&users)
```

### Redis Operations

```go
// Get Redis client
redis := dbManager.GetRedis()

// Get context with timeout
ctx, cancel := database.GetRedisContext()
defer cancel()

// Set a value
err := redis.Set(ctx, "key", "value", 0).Err()

// Get a value
val, err := redis.Get(ctx, "key").Result()
```

### Health Checks

```go
// Perform health check on all connections
if err := dbManager.HealthCheck(); err != nil {
    log.Printf("Database health check failed: %v", err)
} else {
    log.Println("All database connections are healthy")
}
```

## File Structure

```
internal/database/
├── README.md           # This documentation
├── database.go         # Main package exports
├── manager.go          # Database manager implementation
├── postgres.go         # PostgreSQL connection handling
├── redis.go            # Redis connection handling
└── example.go          # Usage examples
```

## Configuration

The database package uses the configuration from `internal/config` package. Key configuration options:

- **PostgreSQL**: Connection string, pool settings, timeouts
- **Redis**: Connection URL, password, database selection
- **Security**: JWT secrets, encryption keys
- **Logging**: Log levels and formats

## Error Handling

All database operations include proper error handling:

- Connection failures are logged and returned as errors
- Health checks validate all connections
- Graceful shutdown ensures proper cleanup
- Context timeouts prevent hanging operations

## Dependencies

- `gorm.io/gorm` - PostgreSQL ORM
- `gorm.io/driver/postgres` - PostgreSQL driver
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/joho/godotenv` - Environment variable loading

## Best Practices

1. **Always use the database manager** instead of creating connections directly
2. **Perform health checks** before critical operations
3. **Use context timeouts** for Redis operations
4. **Handle errors gracefully** and log appropriately
5. **Close connections** properly on application shutdown
6. **Use connection pooling** settings appropriate for your load
