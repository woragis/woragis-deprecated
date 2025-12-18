# ADR-002: Standalone Workers Architecture

## Status
Accepted

## Context
The backend requires asynchronous processing of various jobs (emails, translations, resume generation, etc.). We need to decide how to structure these workers: as part of the main server or as standalone services.

**Requirements:**
- Asynchronous job processing
- Independent scaling
- Fault isolation
- Language flexibility (Go, Python, Node.js)
- Easy deployment and maintenance

**Constraints:**
- Different workers have different requirements (memory, CPU, external APIs)
- Some workers need specific languages (Python for AI, Node.js for web scraping)
- Workers should be independently deployable
- Need to minimize coupling between server and workers

## Decision
We will implement **standalone worker services** that are completely independent from the main server.

**Architecture:**
- Each worker is a separate service/application
- Workers consume messages from RabbitMQ independently
- Workers connect directly to the database
- Workers expose their own health check endpoints
- Workers can be written in different languages (Go, Python, Node.js)
- Workers are deployed and scaled independently

**Worker Structure:**
- Each worker has its own repository/directory
- Each worker has its own Dockerfile
- Each worker has its own configuration
- Each worker has its own health check endpoint (`/healthz`)
- Each worker has its own metrics endpoint (`/metrics`)

## Consequences

### Positive
- ✅ **Independent Scaling**: Scale workers independently based on queue depth
- ✅ **Fault Isolation**: Failure in one worker doesn't affect others or the server
- ✅ **Language Flexibility**: Use the best language for each worker's needs
- ✅ **Independent Deployment**: Deploy workers without affecting the server
- ✅ **Resource Optimization**: Allocate resources per worker based on needs
- ✅ **Development Velocity**: Teams can work on workers independently

### Negative
- ⚠️ **More Services**: More services to manage, monitor, and deploy
- ⚠️ **Code Duplication**: Some shared code (logging, health checks) duplicated across workers
- ⚠️ **Configuration Management**: Each worker needs its own configuration
- ⚠️ **Deployment Complexity**: More services to deploy and coordinate

### Neutral
- Workers share the same database and message queue infrastructure
- Each worker follows similar patterns (health checks, logging, metrics)

## Implementation Details

### Worker Types

#### Go Workers (Email, WhatsApp, Translation)
- **Language**: Go
- **Structure**: `cmd/{worker-name}/main.go`, `pkg/`, `internal/`
- **Health Check**: Standard `net/http` server on port 8080
- **Metrics**: Prometheus client library
- **Logging**: `log/slog` with structured logging

#### Python Workers (Resume)
- **Language**: Python
- **Structure**: `src/main.py`, `src/`, `tests/`
- **Health Check**: Custom HTTP server on port 8080
- **Metrics**: `prometheus-client` library
- **Logging**: `structlog` with structured logging

#### Node.js Workers (Job Application)
- **Language**: Node.js/JavaScript
- **Structure**: `src/index.js`, `src/`, `tests/`
- **Health Check**: Custom HTTP server on port 8080
- **Metrics**: `prom-client` library
- **Logging**: Custom logger with structured logging

### Common Patterns

#### Health Checks
All workers expose `/healthz` endpoint:
- Checks RabbitMQ connection
- Returns JSON: `{"status": "healthy|unhealthy", "checks": [...]}`
- HTTP 200 for healthy, 503 for unhealthy

#### Metrics
All workers expose `/metrics` endpoint:
- Prometheus format
- Worker-specific metrics (job processing, failures, retries)
- Queue metrics (depth, DLQ size)

#### Logging
All workers use structured logging:
- JSON format in production
- Service name in logs
- Trace ID support
- Consistent log levels

#### Message Consumption
All workers:
- Connect to RabbitMQ
- Consume from specific queue
- Process messages asynchronously
- Acknowledge messages after successful processing
- Route failed messages to DLQ

## Alternatives Considered

### 1. Workers as Part of Server
- **Pros**: Simpler deployment, shared code
- **Cons**: Can't scale independently, fault coupling, language restrictions
- **Rejected**: Doesn't meet scaling and fault isolation requirements

### 2. Workers as Server Plugins/Modules
- **Pros**: Shared code, simpler deployment
- **Cons**: Still coupled to server, language restrictions
- **Rejected**: Doesn't meet language flexibility requirement

### 3. Workers as Kubernetes Jobs
- **Pros**: Kubernetes-native, scheduled execution
- **Cons**: Not suitable for long-running consumers, more complex
- **Rejected**: Workers need to run continuously, not as scheduled jobs

### 4. Workers as Serverless Functions
- **Pros**: Auto-scaling, pay-per-use
- **Cons**: Cold starts, vendor lock-in, complexity
- **Rejected**: Need long-running connections to RabbitMQ, not suitable for serverless

## Examples

### Email Worker Structure
```
email-worker/
├── cmd/
│   └── email-worker/
│       └── main.go
├── pkg/
│   ├── health/
│   ├── logger/
│   ├── metrics/
│   └── queue/
├── go.mod
├── Dockerfile
└── README.md
```

### Resume Worker Structure
```
resume-worker/
├── src/
│   ├── main.py
│   ├── metrics.py
│   └── ...
├── tests/
├── requirements.txt
├── Dockerfile
└── README.md
```

## Notes
- Workers are stateless and can run multiple replicas
- Each worker replica consumes from the same queue (RabbitMQ distributes messages)
- Workers write directly to the database (no round-trip through server)
- Future enhancement: Shared libraries for common functionality (if needed)
- Future enhancement: Worker registry/service discovery (if needed)

## Related ADRs
- [ADR-001: RabbitMQ with Redis Fallback](./001-rabbitmq-redis-fallback.md) - Workers use RabbitMQ directly
- [ADR-003: Structured Logging Implementation](./003-structured-logging.md) - Logging patterns in workers
- [ADR-005: Health Checks Implementation Strategy](./005-health-checks.md) - Health check patterns in workers
