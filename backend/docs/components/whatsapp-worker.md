# WhatsApp Worker Component

## Overview

A standalone Go-based worker service that processes WhatsApp messaging jobs from RabbitMQ. It handles asynchronous WhatsApp message delivery and updates the database with delivery status.

## Architecture

- **Language**: Go 1.21+
- **Port**: 8080 (health check only)
- **Message Queue**: RabbitMQ
- **Database**: PostgreSQL (via GORM)
- **WhatsApp API**: External WhatsApp Business API

## Responsibilities

1. **Consume WhatsApp Jobs**: Listen to `whatsapp.queue` for messaging jobs
2. **Send Messages**: Send WhatsApp messages via API
3. **Update Database**: Record message delivery status
4. **Error Handling**: Retry transient failures, route permanent failures to DLQ

## Health Check

**Endpoint**: `GET /healthz`

**Checks**:
- RabbitMQ connection - CRITICAL

**Response**:
```json
{
  "status": "healthy|unhealthy",
  "checks": [
    {"name": "rabbitmq", "status": "ok"}
  ]
}
```

## Metrics

**Endpoint**: `GET /metrics`

Exposes Prometheus metrics:
- Job processing rate (success/failed)
- Job processing duration
- Queue depth
- DLQ size

## Configuration

### Environment Variables

#### Required
- `DATABASE_URL` - PostgreSQL connection string
- `RABBITMQ_URL` - RabbitMQ connection URL
- `WHATSAPP_API_URL` - WhatsApp Business API URL
- `WHATSAPP_API_TOKEN` - WhatsApp API authentication token
- `WHATSAPP_PHONE_NUMBER_ID` - WhatsApp Business phone number ID

#### Optional
- `WHATSAPP_QUEUE_NAME` - Queue name (default: `whatsapp.queue`)
- `WHATSAPP_EXCHANGE` - Exchange name (default: `woragis.tasks`)
- `WHATSAPP_ROUTING_KEY` - Routing key (default: `whatsapp.send`)
- `WHATSAPP_PREFETCH_COUNT` - Prefetch count (default: `1`)
- `ENV` - Environment (development/production)

## Message Format

**Queue**: `whatsapp.queue`

**Message**:
```json
{
  "id": "job-uuid",
  "to": "+1234567890",
  "message": "Message text",
  "media_url": "https://example.com/image.jpg",
  "user_id": "user-uuid"
}
```

## Processing Flow

1. **Consume Message**: Worker consumes message from RabbitMQ
2. **Send WhatsApp Message**: Send message via WhatsApp API
3. **Update Database**: Record delivery status
4. **Acknowledge**: Acknowledge message on success
5. **Retry/DLQ**: Retry on transient errors, route to DLQ on permanent errors

## Error Handling

### Transient Errors
- Network errors
- WhatsApp API temporarily unavailable
- Rate limit errors (with retry-after)

**Action**: Retry with exponential backoff (up to 3 attempts)

### Permanent Errors
- Invalid phone number
- Authentication failures
- Invalid message format

**Action**: Route to Dead Letter Queue (DLQ)

## Dead Letter Queue

**DLQ Exchange**: `woragis.dlx`
**DLQ Routing Key**: `whatsapp.queue.failed`

Failed messages are automatically routed to DLQ after max retries.

## Logging

**Format**: Structured JSON (production), Text (development)

**Service Name**: `whatsapp-worker`

**Key Log Fields**:
- `job_id` - Job identifier
- `to` - Recipient phone number
- `status` - Success/failed
- `error` - Error message (if failed)

## Deployment

### Local Development

```bash
cd backend/whatsapp-worker
go run cmd/whatsapp-worker/main.go
```

### Docker

```bash
docker build -f Dockerfile.whatsapp-worker -t woragis/whatsapp-worker .
docker run --env-file .env woragis/whatsapp-worker
```

### Kubernetes

Deploy as a Deployment or StatefulSet (if leader election needed):
- Health check probe on `/healthz`
- RabbitMQ connection required
- Database connection required
- See `server/k8s/whatsapp-worker/` for example

## Scaling

### Horizontal Scaling
- Multiple worker replicas consume from same queue
- RabbitMQ distributes messages (round-robin or fair dispatch)
- Each replica has its own database connection pool
- **Note**: WhatsApp API rate limits may limit scaling

### Resource Requirements
- **CPU**: 200m-500m (0.2-0.5 core)
- **Memory**: 256Mi-512Mi
- **Database Connections**: 5-10 per replica

## WhatsApp API Integration

### WhatsApp Business API

The worker integrates with WhatsApp Business API (Meta):
- REST API for sending messages
- Template messages for business communication
- Media support (images, documents, videos)

### Rate Limits

WhatsApp API has rate limits:
- Messages per second (varies by tier)
- Daily message limits
- Template message approval required

**Handling**:
- Respect rate limits
- Implement exponential backoff on rate limit errors
- Monitor rate limit usage

## Monitoring

### Key Metrics
- Message sending rate (messages/second)
- Success rate
- Failure rate
- Queue depth
- DLQ size
- Rate limit usage

### Alerts
- DLQ size > 100 messages
- Failure rate > 5%
- Queue depth > 1000 messages
- Rate limit approaching

## Troubleshooting

### Common Issues

#### Messages Not Sending
- Check WhatsApp API configuration
- Check API token validity
- Check phone number ID
- Verify phone numbers are registered

#### High Failure Rate
- Check WhatsApp API status
- Verify message templates are approved
- Check rate limits
- Verify phone number format

#### Queue Backlog
- Check rate limits (may be limiting throughput)
- Scale up worker replicas (if rate limits allow)
- Check worker health
- Verify workers are consuming messages

## Related Documentation

- [Architecture Decision Records](../adr/) - Worker architecture decisions
- [Monitoring DLQ](../runbooks/monitoring-dlq.md) - DLQ monitoring procedures
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
