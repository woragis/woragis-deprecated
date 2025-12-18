# Structured Logging Overview - Backend Architecture

## General Architecture

All backend components implement **structured logging** with consistent patterns but language-specific implementations. The system uses **JSON format in production** and **human-readable text in development** for easy debugging.

### Common Patterns

1. **Format**: JSON in production, text in development
2. **Service Name**: Automatically included in all logs (`service: "service-name"`)
3. **Trace ID**: Support for distributed tracing via context
4. **ISO 8601 Timestamps**: Standard timestamp format
5. **Structured Fields**: Key-value pairs (no string interpolation)
6. **Output**: stdout (for Kubernetes/log aggregation) or files (development only)

---

## Component Breakdown

### 1. **Server** (Main API Server)
**Location**: `backend/server/app/pkg/logger/logger.go`  
**Framework**: Go `log/slog`  
**Service Name**: `"server"`

#### Implementation:
- **Production**: JSON handler (`slog.NewJSONHandler`)
- **Development**: Text handler (`slog.NewTextHandler`)
- **Level**: INFO in production, DEBUG in development
- **Output**: stdout (production) or stdout+file (development if `LOG_TO_FILE=true`)

#### Features:
- ✅ Automatic service name injection
- ✅ Trace ID from context (`trace_id` key)
- ✅ ISO 8601 timestamp formatting
- ✅ Lowercase log levels
- ✅ Custom handler wrapper (`serviceHandler`)

#### Log Format Example:
**Production (JSON)**:
```json
{
  "timestamp": "2024-01-15T10:30:45.123456789Z",
  "level": "info",
  "service": "server",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Request processed",
  "method": "GET",
  "path": "/api/projects",
  "status": 200,
  "duration_ms": 45
}
```

**Development (Text)**:
```
2024-01-15T10:30:45.123456789Z INFO service=server trace_id=550e8400... message="Request processed" method=GET path=/api/projects status=200 duration_ms=45
```

#### Usage:
```go
logger := logger.New(os.Getenv("ENV"))
logger.Info("Request processed",
    slog.String("method", "GET"),
    slog.String("path", "/api/projects"),
    slog.Int("status", 200),
    slog.Int("duration_ms", 45))

// With trace ID
ctx := logger.WithTraceID(ctx, traceID)
logger.InfoContext(ctx, "Processing request", ...)
```

---

### 2. **Go Workers** (Email, WhatsApp, Translation)
**Locations**: 
- `backend/email-worker/pkg/logger/logger.go`
- `backend/whatsapp-worker/pkg/logger/logger.go`
- `backend/translation-worker/pkg/logger/logger.go`

**Framework**: Go `log/slog`  
**Service Names**: `"email-worker"`, `"whatsapp-worker"`, `"translation-worker"`

#### Implementation:
- **Identical Code**: All three workers share the exact same logger implementation (copy-paste pattern)
- **Production**: JSON handler
- **Development**: Text handler
- **Level**: INFO in production, DEBUG in development
- **Output**: stdout (production) or stdout+file (development if `LOG_TO_FILE=true`)

#### Features:
- ✅ Automatic service name injection (different per worker)
- ✅ Trace ID from context
- ✅ ISO 8601 timestamp formatting
- ✅ Lowercase log levels
- ✅ Custom handler wrapper

#### Log Format Example:
**Production (JSON)**:
```json
{
  "timestamp": "2024-01-15T10:30:45.123456789Z",
  "level": "info",
  "service": "email-worker",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Email sent successfully",
  "user_id": "6ad0d828-1234-5678-90ab-cdef12345678",
  "destination": "user@example.com"
}
```

**Development (Text)**:
```
2024-01-15T10:30:45.123456789Z INFO service=email-worker trace_id=550e8400... message="Email sent successfully" user_id=6ad0d828... destination=user@example.com
```

#### Usage:
```go
logger := logger.New(os.Getenv("ENV"))
logger.Info("Email sent successfully",
    slog.String("user_id", userID),
    slog.String("destination", destination))

// With trace ID
ctx := logger.WithTraceID(ctx, traceID)
logger.InfoContext(ctx, "Processing email", ...)
```

#### Differences Between Go Workers:
- **Service Name**: Only difference (email-worker, whatsapp-worker, translation-worker)
- **Code**: Identical implementation (copy-paste pattern)

---

### 3. **Python Services** (AI Service, Creative Service)
**Locations**:
- `backend/ai-service/app/logger.py`
- `backend/creative-service/app/logger.py`

**Framework**: Python `structlog`  
**Service Names**: `"ai-service"`, `"creative-service"`

#### Implementation:
- **Library**: `structlog` (structured logging library)
- **Production**: JSON renderer (`structlog.processors.JSONRenderer`)
- **Development**: Console renderer (`structlog.dev.ConsoleRenderer`)
- **Level**: INFO in production, DEBUG in development
- **Output**: stdout (production) or stdout+file (development if `LOG_TO_FILE=true`)

#### Features:
- ✅ Automatic service name binding (`logger.bind(service="ai-service")`)
- ✅ Trace ID via context variables (`ContextVar`)
- ✅ ISO 8601 timestamp (`TimeStamper(fmt="iso")`)
- ✅ Stack trace support
- ✅ Exception formatting
- ✅ Thread-safe context variables

#### Log Format Example:
**Production (JSON)**:
```json
{
  "timestamp": "2024-01-15T10:30:45.123456789Z",
  "level": "info",
  "logger": "woragis.ai-service",
  "service": "ai-service",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Chat request completed",
  "request_id": "abc123",
  "provider": "openai",
  "duration_ms": 1234
}
```

**Development (Text)**:
```
2024-01-15T10:30:45.123456789Z [info     ] Chat request completed     request_id=abc123 provider=openai duration_ms=1234 service=ai-service trace_id=550e8400... logger=woragis.ai-service
```

#### Usage:
```python
from app.logger import get_logger, set_trace_id

logger = get_logger()
logger.info("Chat request completed",
    request_id=request_id,
    provider=provider,
    duration_ms=duration)

# With trace ID
set_trace_id(trace_id)
logger.info("Processing request", ...)
```

#### Differences:
- **Service Name**: Only difference (ai-service, creative-service)
- **Code**: Identical implementation (copy-paste pattern)

---

### 4. **Resume Worker** (Python)
**Location**: `backend/resume-worker/src/logger.py`  
**Framework**: Python `structlog`  
**Service Name**: `"resume-worker"`

#### Implementation:
- **Library**: `structlog` (same as Python services)
- **Production**: JSON renderer
- **Development**: Console renderer
- **Level**: INFO in production, DEBUG in development
- **Output**: stdout (production) or stdout+file (development if `LOG_TO_FILE=true`)

#### Features:
- ✅ Automatic service name binding
- ✅ Trace ID via context variables
- ✅ ISO 8601 timestamp
- ✅ Stack trace support
- ✅ Exception formatting
- ✅ Thread-safe context variables

#### Log Format Example:
**Production (JSON)**:
```json
{
  "timestamp": "2024-01-15T10:30:45.123456789Z",
  "level": "info",
  "logger": "woragis.resume-worker",
  "service": "resume-worker",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Resume successfully generated",
  "user_id": "6ad0d828-1234-5678-90ab-cdef12345678",
  "file_size": 123456,
  "duration_ms": 4523
}
```

**Development (Text)**:
```
2024-01-15T10:30:45.123456789Z [info     ] Resume successfully generated     user_id=6ad0d828... file_size=123456 duration_ms=4523 service=resume-worker trace_id=550e8400... logger=woragis.resume-worker
```

#### Usage:
```python
from logger import get_logger, set_trace_id

logger = get_logger()
logger.info("Resume successfully generated",
    user_id=user_id,
    file_size=file_size,
    duration_ms=duration)

# With trace ID
set_trace_id(trace_id)
logger.info("Processing resume job", ...)
```

#### Differences from Python Services:
- **Location**: Different directory structure
- **Logger Name**: `"woragis.resume-worker"` vs `"woragis.ai-service"`
- **Implementation**: Same `structlog` pattern

---

### 5. **Job Application Worker** (Node.js)
**Location**: `backend/job-application-worker/src/utils/logger.js`  
**Framework**: Custom JSON logger (Node.js)  
**Service Name**: `"job-application-worker"`

#### Implementation:
- **Custom Logger**: Built with Node.js `console` and file streams
- **Format**: Always JSON (no text format)
- **Level**: All levels in development, info/warn/error in production
- **Output**: stdout (always) + file (development if `LOG_TO_FILE=true`)

#### Features:
- ✅ Automatic service name in all logs
- ✅ Trace ID via module-level variable (simple, not thread-safe)
- ✅ ISO 8601 timestamp
- ✅ JSON.stringify for consistent format
- ✅ File logging support (development only)
- ⚠️ **No text format** (always JSON)

#### Log Format Example:
**All Environments (JSON)**:
```json
{
  "timestamp": "2024-01-15T10:30:45.123Z",
  "level": "info",
  "service": "job-application-worker",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Job application completed",
  "jobId": "123e4567-e89b-12d3-a456-426614174000",
  "user_id": "6ad0d828-1234-5678-90ab-cdef12345678",
  "website": "linkedin.com"
}
```

#### Usage:
```javascript
import { logger, setTraceId } from './utils/logger.js';

logger.info('Job application completed', {
  jobId: job.id,
  user_id: job.userId,
  website: job.website,
  status: 'success'
});

// With trace ID
setTraceId(traceId);
logger.info('Processing job', { jobId });
```

#### Special Features:
- **Always JSON**: Unlike other components, always outputs JSON
- **Module-level trace ID**: Simple variable (not AsyncLocalStorage)
- **File stream**: Uses `fs.createWriteStream` for file logging
- **Close method**: `logger.close()` to close file stream on shutdown

---

## Comparison Table

| Component | Language | Library | Format (Prod) | Format (Dev) | Trace ID | File Logging |
|-----------|----------|---------|---------------|--------------|----------|--------------|
| **Server** | Go | `log/slog` | JSON | Text | Context | Optional |
| **Email Worker** | Go | `log/slog` | JSON | Text | Context | Optional |
| **WhatsApp Worker** | Go | `log/slog` | JSON | Text | Context | Optional |
| **Translation Worker** | Go | `log/slog` | JSON | Text | Context | Optional |
| **Resume Worker** | Python | `structlog` | JSON | Text | ContextVar | Optional |
| **Job App Worker** | Node.js | Custom | JSON | JSON | Variable | Optional |
| **AI Service** | Python | `structlog` | JSON | Text | ContextVar | Optional |
| **Creative Service** | Python | `structlog` | JSON | Text | ContextVar | Optional |

---

## Key Differences

### 1. **Language & Library**
- **Go Components**: Use `log/slog` (Go 1.21+ standard library)
- **Python Components**: Use `structlog` (third-party library)
- **Node.js Component**: Custom implementation (no library)

### 2. **Format Strategy**
- **Go & Python**: JSON in production, text in development
- **Node.js**: Always JSON (no text format)

### 3. **Trace ID Implementation**
- **Go**: Context values (`context.WithValue`)
- **Python**: Context variables (`ContextVar` - thread-safe)
- **Node.js**: Module-level variable (simple, not thread-safe)

### 4. **Service Name Injection**
- **Go**: Custom handler wrapper (`serviceHandler`)
- **Python**: Logger binding (`logger.bind(service="...")`)
- **Node.js**: Included in log entry creation

### 5. **Timestamp Format**
- **All**: ISO 8601 format (RFC3339Nano for Go, ISO for Python/Node.js)
- **Go**: `time.RFC3339Nano`
- **Python**: `TimeStamper(fmt="iso")`
- **Node.js**: `new Date().toISOString()`

### 6. **File Logging**
- **Go**: `io.MultiWriter` (file + stdout)
- **Python**: `logging.FileHandler` added to root logger
- **Node.js**: `fs.createWriteStream` + `console.log`

### 7. **Log Level Control**
- **Go**: Handler options (`slog.LevelInfo` vs `slog.LevelDebug`)
- **Python**: Standard logging level (`logging.INFO` vs `logging.DEBUG`)
- **Node.js**: Conditional debug logging (`if (!isProduction)`)

---

## Common Features

### 1. **Structured Fields**
All components use key-value pairs instead of string interpolation:

**Go**:
```go
logger.Info("Email sent", 
    slog.String("user_id", userID),
    slog.Int("duration_ms", duration))
```

**Python**:
```python
logger.info("Email sent",
    user_id=user_id,
    duration_ms=duration)
```

**Node.js**:
```javascript
logger.info('Email sent', {
  user_id: userID,
  duration_ms: duration
});
```

### 2. **Service Name**
All logs automatically include service name:
- Go: Injected by `serviceHandler`
- Python: Bound via `logger.bind(service="...")`
- Node.js: Added in `createLogEntry`

### 3. **Trace ID Support**
All components support trace ID for distributed tracing:
- **Go**: `logger.WithTraceID(ctx, traceID)` → `logger.InfoContext(ctx, ...)`
- **Python**: `set_trace_id(trace_id)` → `logger.info(...)`
- **Node.js**: `setTraceId(traceId)` → `logger.info(...)`

### 4. **Environment-Based Configuration**
All components respect `ENV` variable:
- `development` or `dev`: Text format (except Node.js), DEBUG level
- `production` or `prod`: JSON format, INFO level

### 5. **File Logging (Development Only)**
All components support optional file logging in development:
- **Environment Variable**: `LOG_TO_FILE=true`
- **Directory**: `LOG_DIR=logs` (default)
- **Behavior**: Writes to both file and stdout (dual output)

---

## Configuration

### Environment Variables

**Common to All**:
```bash
ENV=development  # or "production"
LOG_TO_FILE=true  # Optional: Enable file logging in development
LOG_DIR=logs  # Optional: Log directory (default: "logs")
```

**Component-Specific**:
- None - all use the same environment variables

---

## Log Levels

### Standard Levels (All Components)

1. **DEBUG**: Detailed diagnostic information
   - Only in development
   - Internal state, detailed processing steps
   - Example: "Processing step 3 of 5", "Cache hit for key X"

2. **INFO**: General operational messages
   - Normal business events
   - Successful operations
   - Example: "Email sent", "Job completed", "Request processed"

3. **WARN**: Warnings and recoverable issues
   - Retry attempts
   - Degraded performance
   - Non-critical failures
   - Example: "Retrying after error", "Rate limit approaching", "Fallback to Redis"

4. **ERROR**: Errors and failures
   - System failures
   - Unrecoverable errors
   - Critical issues
   - Example: "Database connection failed", "Translation API error", "Job processing failed"

---

## Trace ID Propagation

### How It Works

1. **Server/Service**: Generate or extract trace ID from request headers
2. **Set in Context**: Store trace ID in context (Go) or context variable (Python) or variable (Node.js)
3. **Automatic Inclusion**: All logs in that context automatically include trace ID
4. **Cross-Service**: Pass trace ID in HTTP headers or message metadata

### Implementation Patterns

**Go (Context)**:
```go
// Set trace ID
ctx := logger.WithTraceID(ctx, traceID)

// Use in logs (automatic)
logger.InfoContext(ctx, "Processing request", ...)
```

**Python (ContextVar)**:
```python
# Set trace ID
set_trace_id(trace_id)

# Use in logs (automatic)
logger.info("Processing request", ...)
```

**Node.js (Variable)**:
```javascript
// Set trace ID
setTraceId(traceId);

// Use in logs (automatic)
logger.info('Processing request', { ... });
```

---

## Best Practices

### 1. **Always Use Structured Fields**
```go
// ✅ Good
logger.Info("Email sent", 
    slog.String("user_id", userID),
    slog.Int("duration_ms", duration))

// ❌ Bad
logger.Info(fmt.Sprintf("Email sent for user %s in %dms", userID, duration))
```

### 2. **Include Relevant Context**
```go
// ✅ Good - includes all relevant context
logger.Error("Translation failed",
    slog.String("job_id", jobID),
    slog.String("entity_type", entityType),
    slog.String("language", language),
    slog.Any("error", err))

// ❌ Bad - missing context
logger.Error("Translation failed", slog.Any("error", err))
```

### 3. **Never Log Sensitive Data**
```go
// ❌ Bad
logger.Info("User login", slog.String("password", password))

// ✅ Good
logger.Info("User login", slog.String("user_id", userID))
```

### 4. **Use Appropriate Log Levels**
- **ERROR**: System failures, unrecoverable errors
- **WARN**: Recoverable issues, retries, degraded performance
- **INFO**: Normal operations, successful completions
- **DEBUG**: Detailed diagnostic info (development only)

### 5. **Include Trace ID for Distributed Tracing**
```go
// Always set trace ID at request/job start
ctx := logger.WithTraceID(ctx, traceID)

// All subsequent logs in this context will include trace_id
logger.InfoContext(ctx, "Processing", ...)
```

---

## Log Aggregation

### Production Strategy

All components write to **stdout** in production, which is collected by:
- **Kubernetes**: Captures stdout/stderr from containers
- **Log Aggregator**: Fluentd/Fluent Bit collects logs
- **Storage**: ELK Stack, Loki, CloudWatch, etc.

### Development Strategy

- **Default**: stdout (human-readable text for Go/Python, JSON for Node.js)
- **Optional**: File logging (`LOG_TO_FILE=true`) writes to both file and stdout
- **Files**: `logs/{service-name}.log`

---

## Code Patterns

### Go Pattern (Server, Workers)
```go
package logger

const ServiceName = "service-name"
const TraceIDKey = "trace_id"

func New(env string) *slog.Logger {
    // Configure handler based on environment
    // JSON in production, text in development
    // Wrap with serviceHandler to inject service name and trace_id
}

type serviceHandler struct {
    slog.Handler
    service string
}

func (h *serviceHandler) Handle(ctx context.Context, r slog.Record) error {
    r.AddAttrs(slog.String("service", h.service))
    if traceID := ctx.Value(TraceIDKey); traceID != nil {
        r.AddAttrs(slog.String("trace_id", traceID.(string)))
    }
    return h.Handler.Handle(ctx, r)
}
```

### Python Pattern (Services, Resume Worker)
```python
import structlog
from contextvars import ContextVar

trace_id_var: ContextVar[Optional[str]] = ContextVar("trace_id", default=None)

def configure_logging(env: str, log_to_file: bool, log_dir: str):
    # Configure structlog processors
    # JSON in production, console in development
    # Add file handler if log_to_file is True

def get_logger(name: str):
    logger = structlog.get_logger(name)
    logger = logger.bind(service="service-name")
    if trace_id := trace_id_var.get():
        logger = logger.bind(trace_id=trace_id)
    return logger
```

### Node.js Pattern (Job Application Worker)
```javascript
let traceId = null;

export const setTraceId = (id) => { traceId = id; };

const createLogEntry = (level, message, meta = {}) => {
  const entry = {
    timestamp: new Date().toISOString(),
    level,
    service: 'job-application-worker',
    message,
    ...meta,
  };
  if (traceId) entry.trace_id = traceId;
  return entry;
};

export const logger = {
  info: (message, meta = {}) => {
    const entry = createLogEntry('info', message, meta);
    console.log(JSON.stringify(entry));
  },
  // ... other levels
};
```

---

## Summary

### Consistency
- ✅ All use structured logging (key-value pairs)
- ✅ All include service name automatically
- ✅ All support trace ID for distributed tracing
- ✅ All use ISO 8601 timestamps
- ✅ All output JSON in production
- ✅ All support file logging in development

### Differences
- **Language**: Go (`slog`) vs Python (`structlog`) vs Node.js (custom)
- **Format**: Text in dev (Go/Python) vs Always JSON (Node.js)
- **Trace ID**: Context (Go) vs ContextVar (Python) vs Variable (Node.js)
- **Implementation**: Standard library (Go) vs Third-party (Python) vs Custom (Node.js)

### Design Philosophy
1. **Consistency**: Same patterns across similar components
2. **Simplicity**: Workers check only what they need
3. **Observability**: Structured logs enable easy parsing and analysis
4. **Performance**: Non-blocking, efficient logging
5. **Flexibility**: Environment-based configuration
