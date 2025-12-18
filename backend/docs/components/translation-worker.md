# Translation Worker Component

## Overview

A standalone Go-based worker service that processes translation jobs from RabbitMQ. It translates content using external translation APIs (Google Translate, DeepL, LibreTranslate) and writes results directly to the database.

## Architecture

- **Language**: Go 1.21+
- **Port**: 8080 (health check only)
- **Message Queue**: RabbitMQ
- **Database**: PostgreSQL (via GORM)
- **Translation APIs**: Google Translate, DeepL, LibreTranslate

## Responsibilities

1. **Consume Translation Jobs**: Listen to `translations.queue` for translation jobs
2. **Translate Content**: Call translation APIs to translate text
3. **Write to Database**: Write translations directly to database (no round-trip through server)
4. **Error Handling**: Retry transient failures with exponential backoff

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
- Translation API call metrics

## Configuration

### Environment Variables

#### Required
- `DATABASE_URL` - PostgreSQL connection string
- `RABBITMQ_URL` - RabbitMQ connection URL

#### Translation Provider
- `TRANSLATION_PROVIDER` - Provider: `google`, `deepl`, or `libre` (default: `google`)

#### Google Translate
- `GOOGLE_TRANSLATE_API_KEY` - Google Cloud Translation API key
- `GOOGLE_CLOUD_PROJECT_ID` - Google Cloud project ID (optional)

#### DeepL
- `DEEPL_API_KEY` - DeepL API key

#### LibreTranslate
- `LIBRE_TRANSLATE_API_URL` - LibreTranslate API URL (default: `https://libretranslate.com/translate`)
- `LIBRE_TRANSLATE_API_KEY` - LibreTranslate API key (optional)

#### Queue Configuration
- `TRANSLATION_QUEUE_NAME` - Queue name (default: `translations.queue`)
- `TRANSLATION_EXCHANGE` - Exchange name (default: `woragis.tasks`)
- `TRANSLATION_ROUTING_KEY` - Routing key (default: `translations.process`)
- `TRANSLATION_PREFETCH_COUNT` - Prefetch count (default: `1`)

#### Translation Settings
- `TRANSLATION_TIMEOUT` - Request timeout in seconds (default: `30`)
- `TRANSLATION_MAX_RETRIES` - Maximum retry attempts (default: `3`)
- `TRANSLATION_RETRY_DELAY` - Retry delay in milliseconds (default: `1000`)

## Message Format

**Queue**: `translations.queue`

**Message**:
```json
{
  "id": "job-uuid",
  "entityType": "project",
  "entityId": "entity-uuid",
  "language": "pt-BR",
  "fields": ["name", "description"],
  "sourceText": {
    "name": "My Project",
    "description": "Project description"
  }
}
```

**Note**: If `sourceText` is empty, worker automatically fetches from database.

## Processing Flow

1. **Consume Message**: Worker consumes message from RabbitMQ
2. **Fetch Source Text** (if not provided): Fetch from database based on `entityType` and `entityId`
3. **Translate**: Call translation API for each field
4. **Write to Database**: Write translations directly to database
5. **Acknowledge**: Acknowledge message on success
6. **Retry/DLQ**: Retry on transient errors, route to DLQ on permanent errors

## Supported Entity Types

- `testimonial`
- `project`
- `certification`
- `skill`
- (More can be added in `internal/database/database.go`)

## Supported Languages

- `en` - English
- `pt-BR` - Portuguese (Brazil)
- `pt` - Portuguese
- `es` - Spanish
- `fr` - French
- `de` - German
- `ru` - Russian
- `ja` - Japanese
- `ko` - Korean
- `zh-CN` - Chinese (Simplified)
- `el` - Greek
- `la` - Latin

## Retry Logic

- **Max Retries**: 3 attempts (configurable)
- **Retry Delay**: Exponential backoff (1s, 2s, 4s)
- **Retry Conditions**: Transient errors (network, timeouts)
- **No Retry**: Permanent errors (invalid input, auth failures)

## Error Handling

### Transient Errors
- Network errors
- Translation API temporarily unavailable
- Timeout errors

**Action**: Retry with exponential backoff (up to 3 attempts)

### Permanent Errors
- Invalid language code
- Authentication failures
- Invalid entity type/ID

**Action**: Route to Dead Letter Queue (DLQ)

## Dead Letter Queue

**DLQ Exchange**: `woragis.dlx`
**DLQ Routing Key**: `translations.queue.failed`

Failed messages are automatically routed to DLQ after max retries.

## Logging

**Format**: Structured JSON (production), Text (development)

**Service Name**: `translation-worker`

**Key Log Fields**:
- `job_id` - Job identifier
- `entity_type` - Entity type
- `entity_id` - Entity ID
- `language` - Target language
- `status` - Success/failed
- `error` - Error message (if failed)

## Deployment

### Local Development

```bash
cd backend/translation-worker
go run cmd/translation-worker/main.go
```

### Docker

```bash
docker build -f Dockerfile.translation-worker -t woragis/translation-worker .
docker run --env-file .env woragis/translation-worker
```

### Kubernetes

Deploy as a Deployment:
- Health check probe on `/healthz`
- RabbitMQ connection required
- Database connection required

## Scaling

### Horizontal Scaling
- Multiple worker replicas consume from same queue
- RabbitMQ distributes messages (round-robin or fair dispatch)
- Each replica has its own database connection pool
- **Note**: Translation API rate limits may limit scaling

### Resource Requirements
- **CPU**: 200m-500m (0.2-0.5 core)
- **Memory**: 256Mi-512Mi
- **Database Connections**: 5-10 per replica

## Translation API Providers

### Google Translate
- **API**: Google Cloud Translation API
- **Rate Limits**: Per project quotas
- **Cost**: Pay-per-character
- **Quality**: High

### DeepL
- **API**: DeepL API v2
- **Rate Limits**: Per subscription tier
- **Cost**: Subscription-based
- **Quality**: Very high (best for European languages)

### LibreTranslate
- **API**: LibreTranslate API (self-hosted or public)
- **Rate Limits**: Self-hosted (unlimited) or public (limited)
- **Cost**: Free (self-hosted) or usage-based (public)
- **Quality**: Good (open-source)

## Monitoring

### Key Metrics
- Translation job processing rate (jobs/second)
- Success rate
- Failure rate
- Queue depth
- DLQ size
- Translation API call duration
- Translation API error rate

### Alerts
- DLQ size > 100 messages
- Failure rate > 5%
- Queue depth > 1000 messages
- Translation API error rate > 10%

## Troubleshooting

### Common Issues

#### Translations Not Processing
- Check translation API configuration
- Check API keys validity
- Check database connection
- Verify message format

#### High Failure Rate
- Check translation API status
- Verify language codes are valid
- Check rate limits
- Verify entity types/IDs exist in database

#### Queue Backlog
- Check rate limits (may be limiting throughput)
- Scale up worker replicas (if rate limits allow)
- Check worker health
- Verify workers are consuming messages

## Related Documentation

- [ADR-004: Translation Worker Architecture](../adr/004-translation-worker.md) - Architecture decision
- [Architecture Decision Records](../adr/) - Other architectural decisions
- [Monitoring DLQ](../runbooks/monitoring-dlq.md) - DLQ monitoring procedures
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
