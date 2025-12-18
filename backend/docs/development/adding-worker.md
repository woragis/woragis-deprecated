# Adding a New Worker

## Overview

This guide explains how to create a new worker service for processing asynchronous jobs. Workers can be written in Go, Python, or Node.js depending on requirements.

## Worker Structure

Workers follow this structure:

```
{worker-name}/
├── cmd/{worker-name}/main.go  # Entry point (Go)
├── src/main.py                # Entry point (Python)
├── src/index.js               # Entry point (Node.js)
├── pkg/                       # Packages (Go)
│   ├── health/
│   ├── logger/
│   └── metrics/
├── internal/                  # Internal packages (Go)
├── src/                       # Source code (Python/Node.js)
├── tests/                     # Tests
├── Dockerfile
├── go.mod / requirements.txt / package.json
└── README.md
```

## Step-by-Step Guide

### 1. Choose Language

**Go**: Best for performance, consistency with other Go workers
**Python**: Best for AI/ML integration, rapid development
**Node.js**: Best for web scraping, browser automation

### 2. Create Worker Directory

```bash
mkdir -p backend/{worker-name}
cd backend/{worker-name}
```

### 3. Initialize Project

#### Go Worker

```bash
go mod init github.com/woragis/backend/{worker-name}
```

#### Python Worker

```bash
python -m venv venv
source venv/bin/activate  # or `venv\Scripts\activate` on Windows
pip install -r requirements.txt
```

#### Node.js Worker

```bash
npm init -y
npm install
```

### 4. Implement Core Components

#### Health Check

All workers expose `/healthz` endpoint:

**Go**:
```go
// pkg/health/health.go
func Handler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check RabbitMQ connection
        status := "healthy"
        if !isRabbitMQConnected() {
            status = "unhealthy"
            w.WriteHeader(503)
        }
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status": status,
            "checks": []map[string]string{
                {"name": "rabbitmq", "status": "ok"},
            },
        })
    }
}
```

**Python**:
```python
# src/health.py
def health_handler():
    # Check RabbitMQ connection
    status = "healthy"
    if not is_rabbitmq_connected():
        status = "unhealthy"
        return {"status": status, "checks": [...]}, 503
    return {"status": status, "checks": [...]}, 200
```

**Node.js**:
```javascript
// src/health.js
function healthHandler(req, res) {
    // Check RabbitMQ connection
    const status = isRabbitMQConnected() ? "healthy" : "unhealthy";
    res.status(status === "healthy" ? 200 : 503).json({
        status,
        checks: [{ name: "rabbitmq", status: "ok" }]
    });
}
```

#### Metrics

All workers expose `/metrics` endpoint:

**Go**:
```go
// pkg/metrics/metrics.go
var (
    JobsProcessed = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "worker_jobs_processed_total",
            Help: "Total jobs processed",
        },
        []string{"worker", "status"},
    )
)

// In main.go
http.Handle("/metrics", promhttp.Handler())
```

**Python**:
```python
# src/metrics.py
from prometheus_client import Counter, generate_latest

jobs_processed = Counter(
    'worker_jobs_processed_total',
    'Total jobs processed',
    ['worker', 'status']
)

def metrics_handler():
    return generate_latest()
```

**Node.js**:
```javascript
// src/metrics.js
const promClient = require('prom-client');

const jobsProcessed = new promClient.Counter({
    name: 'worker_jobs_processed_total',
    help: 'Total jobs processed',
    labelNames: ['worker', 'status']
});

function metricsHandler(req, res) {
    res.set('Content-Type', promClient.register.contentType);
    res.end(promClient.register.metrics());
}
```

#### Logging

All workers use structured logging:

**Go**:
```go
// pkg/logger/logger.go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger = logger.With("service", "worker-name")
logger.Info("Job processed", "job_id", jobID, "status", "success")
```

**Python**:
```python
# src/logger.py
import structlog
logger = structlog.get_logger().bind(service="worker-name")
logger.info("Job processed", job_id=job_id, status="success")
```

**Node.js**:
```javascript
// src/logger.js
const logger = createLogger({ service: 'worker-name' });
logger.info('Job processed', { job_id, status: 'success' });
```

### 5. Implement Message Consumption

#### Go Worker

```go
// cmd/{worker-name}/main.go
func main() {
    conn, err := amqp.Dial(rabbitmqURL)
    // ... error handling
    
    ch, err := conn.Channel()
    // ... error handling
    
    msgs, err := ch.Consume(
        queueName,
        "",
        false, // manual ack
        false,
        false,
        false,
        nil,
    )
    
    for msg := range msgs {
        if err := processJob(msg.Body); err != nil {
            msg.Nack(false, true) // requeue
            continue
        }
        msg.Ack(false)
    }
}
```

#### Python Worker

```python
# src/main.py
import pika

connection = pika.BlockingConnection(pika.URLParameters(rabbitmq_url))
channel = connection.channel()

def callback(ch, method, properties, body):
    try:
        process_job(body)
        ch.basic_ack(delivery_tag=method.delivery_tag)
    except Exception as e:
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)

channel.basic_consume(queue=queue_name, on_message_callback=callback)
channel.start_consuming()
```

#### Node.js Worker

```javascript
// src/index.js
const amqp = require('amqplib');

const connection = await amqp.connect(rabbitmqURL);
const channel = await connection.createChannel();

await channel.consume(queueName, async (msg) => {
    try {
        await processJob(msg.content.toString());
        channel.ack(msg);
    } catch (error) {
        channel.nack(msg, false, true); // requeue
    }
});
```

### 6. Implement Job Processing

```go
// Example: Go worker
func processJob(body []byte) error {
    var job Job
    if err := json.Unmarshal(body, &job); err != nil {
        return err
    }
    
    // Process job
    result, err := doWork(job)
    if err != nil {
        return err
    }
    
    // Update database
    return updateDatabase(job.ID, result)
}
```

### 7. Add Dead Letter Queue

Configure queue with DLX:

**Go**:
```go
queueArgs := amqp.Table{
    "x-dead-letter-exchange":    "woragis.dlx",
    "x-dead-letter-routing-key": queueName + ".failed",
}
```

**Python**:
```python
queue_args = {
    'x-dead-letter-exchange': 'woragis.dlx',
    'x-dead-letter-routing-key': f'{queue_name}.failed'
}
```

**Node.js**:
```javascript
const queueArgs = {
    'x-dead-letter-exchange': 'woragis.dlx',
    'x-dead-letter-routing-key': `${queueName}.failed`
};
```

### 8. Create Dockerfile

**Go**:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/{worker-name} ./cmd/{worker-name}

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/{worker-name} /app/{worker-name}
EXPOSE 8080
CMD ["/app/{worker-name}"]
```

**Python**:
```dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
EXPOSE 8080
CMD ["python", "src/main.py"]
```

**Node.js**:
```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
EXPOSE 8080
CMD ["node", "src/index.js"]
```

### 9. Add to docker-compose.yml

```yaml
{worker-name}:
  build:
    context: .
    dockerfile: Dockerfile.{worker-name}
  environment:
    DATABASE_URL: ${DATABASE_URL}
    RABBITMQ_URL: ${RABBITMQ_URL}
    ENV: ${APP_ENV:-development}
  ports:
    - "8080:8080"
  restart: on-failure
```

### 10. Write Tests

Create unit and integration tests:

```go
// Go example
func TestProcessJob(t *testing.T) {
    // Test implementation
}
```

## Configuration

### Environment Variables

All workers need:
- `DATABASE_URL` - PostgreSQL connection string
- `RABBITMQ_URL` - RabbitMQ connection URL
- `ENV` - Environment (development/production)

Worker-specific:
- `{WORKER}_QUEUE_NAME` - Queue name
- `{WORKER}_EXCHANGE` - Exchange name
- `{WORKER}_ROUTING_KEY` - Routing key

## Best Practices

1. **Error Handling**:
   - Retry transient errors
   - Route permanent errors to DLQ
   - Log all errors with context

2. **Metrics**:
   - Record job processing rate
   - Record job duration
   - Record success/failure rates

3. **Logging**:
   - Log job start/end
   - Include job ID in logs
   - Use structured logging

4. **Resource Management**:
   - Close connections properly
   - Handle graceful shutdown
   - Clean up resources

5. **Testing**:
   - Unit tests for job processing logic
   - Integration tests for queue consumption
   - Mock external dependencies

## Example Workers

See existing workers for examples:
- `email-worker/` - Go worker example
- `resume-worker/` - Python worker example
- `job-application-worker/` - Node.js worker example

## Related Documentation

- [Component Documentation](../components/) - Worker component details
- [Testing Patterns](./testing-patterns.md) - Testing guidelines
- [Logging Conventions](./logging-conventions.md) - Logging guidelines
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
