# Structured Logging Implementation

## Current Status

**Logs are currently written to:**
- **Development:** stdout (default) or files (if enabled)
- **Production:** stdout (for Kubernetes/log aggregation)

## Log Storage Architecture

### Development Environment

**Option 1: stdout (default)**
- Logs go to console/terminal
- Simple, no file management
- Good for local development

**Option 2: File logging (optional)**
```go
logger := logger.NewWithConfig(logger.LogConfig{
    Env:       "development",
    LogDir:    "logs",        // Directory for log files
    LogToFile: true,          // Enable file logging
})
```
- Logs written to `logs/server.log`
- Also outputs to stdout (dual output)
- Useful for debugging and log persistence

### Production Environment

**Recommended Architecture:**
```
Application (stdout) 
    ↓
Kubernetes Log Collection (automatic)
    ↓
Log Aggregator (Fluentd/Fluent Bit)
    ↓
Centralized Storage
    ├─→ Database (PostgreSQL/MySQL) - for structured queries
    ├─→ ELK Stack (Elasticsearch) - for full-text search
    ├─→ Loki - for log aggregation
    └─→ CloudWatch/Cloud Logging - cloud-native
```

**Why stdout in production?**
- ✅ Kubernetes automatically collects stdout/stderr
- ✅ No file management in containers
- ✅ Log aggregators read from stdout
- ✅ Better performance (no file I/O)
- ✅ Standard practice for containerized apps

**Why NOT direct database writes?**
- ❌ High write volume impacts app performance
- ❌ Databases not optimized for log streams
- ❌ Harder to query/search at scale
- ❌ Better to use log aggregators that buffer/batch writes

## Log Format

### Production (JSON)
```json
{
  "timestamp": "2024-01-15T10:30:45.123456789Z",
  "level": "info",
  "service": "server",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "http request",
  "method": "GET",
  "path": "/api/projects",
  "status": 200,
  "duration": "45ms"
}
```

### Development (Text)
```
2024-01-15T10:30:45.123456789Z INFO service=server trace_id=550e8400... message="http request" method=GET path=/api/projects status=200
```

## Configuration

### Environment Variables

```bash
# Environment (affects log format and level)
ENV=development  # or "production"

# Optional: Enable file logging in development
LOG_TO_FILE=true
LOG_DIR=logs
```

### Code Usage

```go
// Default: stdout in dev, stdout in prod
logger := logger.New(os.Getenv("ENV"))

// Custom: file logging in development
logger := logger.NewWithConfig(logger.LogConfig{
    Env:       os.Getenv("ENV"),
    LogDir:    os.Getenv("LOG_DIR"),
    LogToFile: os.Getenv("LOG_TO_FILE") == "true",
})
```

## Production Log Aggregation Setup

### Option 1: Database Storage (via Log Aggregator)

**Architecture:**
1. Application writes to stdout (JSON format)
2. Kubernetes collects logs
3. Fluentd/Fluent Bit forwards to database
4. One table per service: `server_logs`, `email_worker_logs`, etc.

**Fluentd Configuration Example:**
```ruby
<match woragis.server>
  @type postgres
  host postgres-logging.example.com
  database logs
  table server_logs
  <buffer>
    flush_interval 10s
  </buffer>
</match>
```

**Database Schema:**
```sql
CREATE TABLE server_logs (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL,
    level VARCHAR(10) NOT NULL,
    service VARCHAR(50) NOT NULL,
    trace_id UUID,
    message TEXT,
    data JSONB,  -- All other structured fields
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_server_logs_timestamp ON server_logs(timestamp);
CREATE INDEX idx_server_logs_trace_id ON server_logs(trace_id);
CREATE INDEX idx_server_logs_level ON server_logs(level);
```

### Option 2: ELK Stack (Elasticsearch)

- Better for full-text search
- Handles high volume better
- More complex setup

### Option 3: Loki (Grafana)

- Lightweight, designed for logs
- Integrates with Grafana
- Good for Kubernetes

## Log Retention

- **Development files:** Manual cleanup or log rotation
- **Production:** Configure retention in log aggregator/storage
  - Database: Partition tables by date, archive old partitions
  - ELK: Index lifecycle management
  - Loki: Retention policies

## Best Practices

1. **Never log sensitive data** (passwords, tokens, PII)
2. **Use appropriate log levels** (debug/info/warn/error)
3. **Include trace_id** for distributed tracing
4. **Structured fields** instead of string interpolation
5. **Production: stdout only** (let Kubernetes handle collection)
6. **Development: stdout or files** (your choice)

## Migration Path

1. ✅ **Current:** stdout in all environments
2. **Next:** Add file logging option for development
3. **Production:** Set up log aggregator (Fluentd/Fluent Bit)
4. **Storage:** Configure database/ELK/Loki for centralized storage
