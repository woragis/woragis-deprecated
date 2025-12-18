# Email Worker Component

## Overview

A standalone Go-based worker service that processes email sending jobs from RabbitMQ. It handles asynchronous email delivery and updates the database with delivery status.

## Architecture

- **Language**: Go 1.21+
- **Port**: 8080 (health check only)
- **Message Queue**: RabbitMQ
- **Database**: PostgreSQL (via GORM)
- **SMTP**: External SMTP server

## Responsibilities

1. **Consume Email Jobs**: Listen to `emails.queue` for email sending jobs
2. **Send Emails**: Send emails via SMTP
3. **Update Database**: Record email delivery status
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
- `SMTP_HOST` - SMTP server host
- `SMTP_PORT` - SMTP server port
- `SMTP_USER` - SMTP username
- `SMTP_PASSWORD` - SMTP password
- `SMTP_FROM` - Default sender email

#### Optional
- `EMAIL_QUEUE_NAME` - Queue name (default: `emails.queue`)
- `EMAIL_EXCHANGE` - Exchange name (default: `woragis.tasks`)
- `EMAIL_ROUTING_KEY` - Routing key (default: `emails.send`)
- `EMAIL_PREFETCH_COUNT` - Prefetch count (default: `1`)
- `ENV` - Environment (development/production)

## Message Format

**Queue**: `emails.queue`

**Message**:
```json
{
  "id": "job-uuid",
  "to": "recipient@example.com",
  "subject": "Email subject",
  "body": "Email body (HTML or text)",
  "from": "sender@example.com",
  "reply_to": "reply@example.com",
  "user_id": "user-uuid"
}
```

## Processing Flow

1. **Consume Message**: Worker consumes message from RabbitMQ
2. **Send Email**: Send email via SMTP
3. **Update Database**: Record delivery status
4. **Acknowledge**: Acknowledge message on success
5. **Retry/DLQ**: Retry on transient errors, route to DLQ on permanent errors

## Error Handling

### Transient Errors
- Network errors
- SMTP server temporarily unavailable
- Timeout errors

**Action**: Retry with exponential backoff (up to 3 attempts)

### Permanent Errors
- Invalid email address
- Authentication failures
- Invalid message format

**Action**: Route to Dead Letter Queue (DLQ)

## Dead Letter Queue

**DLQ Exchange**: `woragis.dlx`
**DLQ Routing Key**: `emails.queue.failed`

Failed messages are automatically routed to DLQ after max retries.

## Logging

**Format**: Structured JSON (production), Text (development)

**Service Name**: `email-worker`

**Key Log Fields**:
- `job_id` - Job identifier
- `to` - Recipient email
- `status` - Success/failed
- `error` - Error message (if failed)

## Deployment

### Local Development

```bash
cd backend/email-worker
go run cmd/email-worker/main.go
```

### Docker

```bash
docker build -f Dockerfile.email-worker -t woragis/email-worker .
docker run --env-file .env woragis/email-worker
```

### Kubernetes

Deploy as a Deployment with:
- Health check probe on `/healthz`
- RabbitMQ connection required
- Database connection required

## Scaling

### Horizontal Scaling
- Multiple worker replicas consume from same queue
- RabbitMQ distributes messages (round-robin or fair dispatch)
- Each replica has its own database connection pool

### Resource Requirements
- **CPU**: 200m-500m (0.2-0.5 core)
- **Memory**: 256Mi-512Mi
- **Database Connections**: 5-10 per replica

## SMTP Configuration

### Common SMTP Providers

#### Gmail
- Host: `smtp.gmail.com`
- Port: `587` (TLS) or `465` (SSL)
- Requires app password (not regular password)

#### SendGrid
- Host: `smtp.sendgrid.net`
- Port: `587`
- Username: `apikey`
- Password: SendGrid API key

#### AWS SES
- Host: `email-smtp.{region}.amazonaws.com`
- Port: `587` or `465`
- Requires AWS credentials

## Monitoring

### Key Metrics
- Email sending rate (emails/second)
- Success rate
- Failure rate
- Queue depth
- DLQ size

### Alerts
- DLQ size > 100 messages
- Failure rate > 5%
- Queue depth > 1000 messages

## Troubleshooting

### Common Issues

#### Emails Not Sending
- Check SMTP configuration
- Check SMTP server status
- Check authentication credentials
- Check firewall/network connectivity

#### High Failure Rate
- Check SMTP server logs
- Verify email addresses are valid
- Check rate limits on SMTP provider

#### Queue Backlog
- Scale up worker replicas
- Check worker health
- Verify workers are consuming messages

## Related Documentation

- [Architecture Decision Records](../adr/) - Worker architecture decisions
- [Monitoring DLQ](../runbooks/monitoring-dlq.md) - DLQ monitoring procedures
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
