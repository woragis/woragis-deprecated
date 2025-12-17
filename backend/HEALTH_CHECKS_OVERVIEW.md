# Health Checks Overview - Backend Architecture

## General Architecture

All backend components implement health checks via a **`GET /healthz`** endpoint that returns JSON responses. The system uses a **consistent pattern** across all services with **component-specific variations**.

### Common Patterns

1. **Endpoint**: All use `GET /healthz` (no authentication required)
2. **Response Format**: JSON with `status` and `checks` array
3. **Caching**: Results cached for **5 seconds** to reduce load
4. **HTTP Status Codes**:
   - `200 OK`: Healthy or degraded (service operational)
   - `503 Service Unavailable`: Unhealthy (service cannot function)
5. **Timeout**: All checks use **2-second timeouts** to prevent hanging

---

## Component Breakdown

### 1. **Server** (Main API Server)
**Location**: `backend/server/app/pkg/health/health.go`  
**Port**: `8080` (same as main API)  
**Framework**: Fiber (Go)

#### Health Checks:
- ✅ **Database (PostgreSQL)** - **CRITICAL**
  - Checks: Connection ping with 2s timeout
  - Failure: Service cannot function → **503 Unhealthy**
  
- ✅ **Redis** - **CRITICAL**
  - Checks: Connection ping with 2s timeout
  - Failure: Service cannot function → **503 Unhealthy**
  
- ⚠️ **RabbitMQ** - **OPTIONAL** (non-critical)
  - Checks: Connection status via `IsConnected()` interface
  - Failure: Service degrades (falls back to Redis) → **200 Degraded**

#### Status Values:
- `healthy`: All critical dependencies OK
- `degraded`: RabbitMQ failed (non-critical)
- `unhealthy`: Database or Redis failed (critical)

#### Special Features:
- **Multi-dependency checking**: Most comprehensive health check
- **Degraded state**: Distinguishes between critical and non-critical failures
- **Dynamic RabbitMQ**: RabbitMQ checker set after connection (optional)
- **Fiber integration**: Uses Fiber handler (not standard HTTP)

#### Response Example:
```json
{
  "status": "healthy",
  "checks": [
    {"name": "database", "status": "ok"},
    {"name": "redis", "status": "ok"},
    {"name": "rabbitmq", "status": "ok"}
  ]
}
```

---

### 2. **Go Workers** (Email, WhatsApp, Translation)
**Locations**: 
- `backend/email-worker/pkg/health/health.go`
- `backend/whatsapp-worker/pkg/health/health.go`
- `backend/translation-worker/pkg/health/health.go`

**Port**: `8080` (dedicated HTTP server)  
**Framework**: Standard `net/http` (Go)

#### Health Checks:
- ✅ **RabbitMQ** - **CRITICAL**
  - Checks: Connection status via `IsClosed()` method
  - Failure: Worker cannot process jobs → **503 Unhealthy**

#### Status Values:
- `healthy`: RabbitMQ connected
- `unhealthy`: RabbitMQ disconnected or error

#### Special Features:
- **Dedicated HTTP server**: Runs on separate port 8080 alongside worker
- **Simple binary status**: Only healthy/unhealthy (no degraded state)
- **Standard HTTP**: Uses `net/http.HandlerFunc` (not Fiber)
- **Graceful shutdown**: Health server shuts down with worker

#### Response Example:
```json
{
  "status": "healthy",
  "checks": [
    {"name": "rabbitmq", "status": "ok"}
  ]
}
```

#### Differences Between Go Workers:
- **Email Worker**: Identical implementation
- **WhatsApp Worker**: Identical implementation  
- **Translation Worker**: Identical implementation
- All three share the **exact same code pattern** (copy-paste with different package names)

---

### 3. **Python Workers** (Resume, Job Application)
**Locations**:
- `backend/resume-worker/` (Python)
- `backend/job-application-worker/` (Node.js)

#### Resume Worker (Python)
**Port**: `8080`  
**Framework**: Python's `HTTPServer`

#### Health Checks:
- ✅ **Service** - **CRITICAL**
  - Checks: Service availability (always ok if endpoint reachable)
  
- ✅ **RabbitMQ** - **CRITICAL**
  - Checks: Connection status
  - Failure: Worker cannot consume jobs → **503 Unhealthy**

#### Special Features:
- **Daemon thread**: HTTP server runs in background thread
- **Python-specific**: Uses Python's built-in HTTP server
- **Service check**: Includes service availability check

#### Response Example:
```json
{
  "status": "healthy",
  "checks": [
    {"name": "service", "status": "ok"},
    {"name": "rabbitmq", "status": "ok"}
  ]
}
```

#### Job Application Worker (Node.js)
**Port**: `8080`  
**Framework**: Node.js `http.createServer`

#### Health Checks:
- ✅ **Service** - **CRITICAL**
  - Checks: Service availability
  
- ✅ **RabbitMQ** - **CRITICAL**
  - Checks: Connection status
  - Failure: Worker cannot consume jobs → **503 Unhealthy**

#### Special Features:
- **Node.js HTTP**: Uses `http.createServer`
- **Module-level caching**: Uses module variables for caching
- **Service check**: Includes service availability check

#### Response Example:
```json
{
  "status": "healthy",
  "checks": [
    {"name": "service", "status": "ok"},
    {"name": "rabbitmq", "status": "ok"}
  ]
}
```

---

### 4. **Python Services** (AI Service, Creative Service)
**Locations**:
- `backend/ai-service/` (FastAPI)
- `backend/creative-service/` (FastAPI)

**Port**: `8000` (same as main API)  
**Framework**: FastAPI

#### Health Checks:
- ✅ **Service** - **CRITICAL**
  - Checks: Service availability (always ok if endpoint reachable)
  - Failure: Service cannot function → **503 Unhealthy**

#### Status Values:
- `healthy`: Service is running
- `unhealthy`: Service unavailable

#### Special Features:
- **FastAPI integration**: Uses FastAPI route decorators
- **Simplest check**: Only checks if service is reachable
- **No external dependencies**: Doesn't check database/Redis/RabbitMQ
- **Service-level**: Focuses on service availability only

#### Response Example:
```json
{
  "status": "healthy",
  "checks": [
    {"name": "service", "status": "ok"}
  ]
}
```

#### Differences:
- **AI Service**: FastAPI with `/healthz` endpoint
- **Creative Service**: FastAPI with `/healthz` endpoint
- Both have **identical implementations** (simple service check)

---

## Comparison Table

| Component | Port | Framework | Checks | Status Values | Special Features |
|-----------|------|-----------|--------|---------------|------------------|
| **Server** | 8080 | Fiber | DB, Redis, RabbitMQ (opt) | healthy/degraded/unhealthy | Multi-dependency, degraded state |
| **Email Worker** | 8080 | net/http | RabbitMQ | healthy/unhealthy | Dedicated HTTP server |
| **WhatsApp Worker** | 8080 | net/http | RabbitMQ | healthy/unhealthy | Dedicated HTTP server |
| **Translation Worker** | 8080 | net/http | RabbitMQ | healthy/unhealthy | Dedicated HTTP server |
| **Resume Worker** | 8080 | HTTPServer | Service, RabbitMQ | healthy/unhealthy | Python daemon thread |
| **Job App Worker** | 8080 | http.createServer | Service, RabbitMQ | healthy/unhealthy | Node.js module caching |
| **AI Service** | 8000 | FastAPI | Service | healthy/unhealthy | Simplest, no deps |
| **Creative Service** | 8000 | FastAPI | Service | healthy/unhealthy | Simplest, no deps |

---

## Key Differences

### 1. **Status Granularity**
- **Server**: 3 states (healthy/degraded/unhealthy) - most sophisticated
- **All Others**: 2 states (healthy/unhealthy) - binary

### 2. **Dependency Checking**
- **Server**: Checks 3 dependencies (DB, Redis, RabbitMQ)
- **Go Workers**: Check 1 dependency (RabbitMQ)
- **Python/Node Workers**: Check 2 dependencies (Service, RabbitMQ)
- **Python Services**: Check 0 dependencies (just service availability)

### 3. **HTTP Framework**
- **Server**: Fiber (Go web framework)
- **Go Workers**: Standard `net/http`
- **Python Services**: FastAPI
- **Python Workers**: Python's `HTTPServer`
- **Node Workers**: Node.js `http.createServer`

### 4. **Port Strategy**
- **Server**: Port 8080 (same as main API)
- **Go Workers**: Port 8080 (dedicated HTTP server)
- **Python/Node Workers**: Port 8080 (dedicated HTTP server)
- **Python Services**: Port 8000 (same as main API)

### 5. **Caching Implementation**
- **Go Components**: `sync.RWMutex` with time-based cache
- **Python Services**: FastAPI dependency injection caching
- **Node Workers**: Module-level variables

### 6. **Criticality**
- **Server**: Database and Redis are critical; RabbitMQ is optional
- **Workers**: RabbitMQ is always critical (cannot function without it)
- **Services**: Only service availability matters

---

## Common Implementation Details

### Caching Strategy
All components cache health check results for **5 seconds** to:
- Reduce load on dependencies (DB, Redis, RabbitMQ)
- Improve response time for frequent health checks
- Prevent excessive connection attempts

### Timeout Strategy
All dependency checks use **2-second timeouts** to:
- Prevent hanging on slow/unresponsive dependencies
- Ensure health checks complete quickly
- Fail fast rather than wait indefinitely

### Thread Safety
- **Go Components**: Use `sync.RWMutex` for concurrent access
- **Python Components**: Use thread-safe patterns or daemon threads
- **Node Components**: Single-threaded event loop (inherently safe)

---

## Kubernetes Integration

All components support Kubernetes probes with consistent configuration:

### Liveness Probe
```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080  # or 8000 for services
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

### Readiness Probe
```yaml
readinessProbe:
  httpGet:
    path: /healthz
    port: 8080  # or 8000 for services
  initialDelaySeconds: 10
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

---

## Usage Examples

### Check Server Health
```bash
curl http://localhost:8080/healthz
```

### Check Worker Health
```bash
curl http://localhost:8080/healthz  # Email/WhatsApp/Translation workers
```

### Check Service Health
```bash
curl http://localhost:8000/healthz  # AI/Creative services
```

### Monitor Health Status
```bash
# Watch for unhealthy status
watch -n 1 'curl -s http://localhost:8080/healthz | jq .status'
```

---

## Design Philosophy

1. **Consistency**: All components use `/healthz` endpoint
2. **Simplicity**: Workers check only what they need (RabbitMQ)
3. **Comprehensiveness**: Server checks all dependencies
4. **Performance**: Caching prevents excessive load
5. **Reliability**: Timeouts prevent hanging
6. **Observability**: JSON responses provide detailed status

---

## Future Improvements

Potential enhancements:
- Add metrics/telemetry to health checks
- Implement health check aggregation endpoint
- Add dependency health history/trending
- Support for custom health check plugins
- Health check webhooks/notifications
