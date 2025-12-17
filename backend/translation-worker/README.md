# Translation Worker

A standalone Go-based worker service that processes translation jobs from RabbitMQ and writes results directly to the database.

## Overview

The Translation Worker follows the hybrid approach:
- **Go-based** for performance and consistency with other workers
- **HTTP-based translation APIs** (Google Translate, DeepL, LibreTranslate)
- **Direct database writes** for simplicity and performance

## Features

- ✅ RabbitMQ message consumption
- ✅ Multiple translation provider support (Google, DeepL, LibreTranslate)
- ✅ Automatic retry logic with configurable attempts
- ✅ Direct database writes (no round-trip to server)
- ✅ Health check endpoint (`/healthz`)
- ✅ Structured logging
- ✅ Graceful shutdown
- ✅ Source text auto-fetching from database when not provided

## Architecture

```
Server → RabbitMQ → Translation Worker → Translation API → Database
```

1. Server publishes translation jobs to RabbitMQ
2. Translation Worker consumes jobs from queue
3. Worker calls translation API (Google/DeepL/LibreTranslate)
4. Worker writes translated results directly to database

## Configuration

### Environment Variables

#### Required
- `DATABASE_URL` - PostgreSQL connection string
- `RABBITMQ_URL` - RabbitMQ connection URL (or use individual components)

#### Translation Provider
- `TRANSLATION_PROVIDER` - Provider to use: `google`, `deepl`, or `libre` (default: `google`)

#### Google Translate
- `GOOGLE_TRANSLATE_API_KEY` - Google Cloud Translation API key
- `GOOGLE_CLOUD_PROJECT_ID` - Google Cloud project ID (optional)

#### DeepL
- `DEEPL_API_KEY` - DeepL API key

#### LibreTranslate
- `LIBRE_TRANSLATE_API_URL` - LibreTranslate API URL (default: `https://libretranslate.com/translate`)
- `LIBRE_TRANSLATE_API_KEY` - LibreTranslate API key (optional for public instance)

#### RabbitMQ (if not using RABBITMQ_URL)
- `RABBITMQ_USER` - RabbitMQ username (default: `woragis`)
- `RABBITMQ_PASSWORD` - RabbitMQ password (default: `woragis`)
- `RABBITMQ_HOST` - RabbitMQ host (default: `rabbitmq`)
- `RABBITMQ_PORT` - RabbitMQ port (default: `5672`)
- `RABBITMQ_VHOST` - RabbitMQ vhost (default: `woragis`)

#### Worker Configuration
- `TRANSLATION_QUEUE_NAME` - Queue name (default: `translations.queue`)
- `TRANSLATION_EXCHANGE` - Exchange name (default: `woragis.tasks`)
- `TRANSLATION_ROUTING_KEY` - Routing key (default: `translations.process`)
- `TRANSLATION_PREFETCH_COUNT` - Prefetch count (default: `1`)

#### Translation API Settings
- `TRANSLATION_TIMEOUT` - Request timeout in seconds (default: `30`)
- `TRANSLATION_MAX_RETRIES` - Maximum retry attempts (default: `3`)
- `TRANSLATION_RETRY_DELAY` - Retry delay in milliseconds (default: `1000`)

#### Logging
- `ENV` - Environment: `development` or `production` (default: `development`)
- `LOG_TO_FILE` - Enable file logging in development (default: `false`)
- `LOG_DIR` - Log directory (default: `logs`)

## Message Format

The worker expects `TranslationJob` messages from RabbitMQ:

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

If `sourceText` is empty, the worker will automatically fetch it from the database based on `entityType` and `entityId`.

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

## Building

```bash
# Install dependencies
make install

# Build binary
make build

# Build Docker image
make docker-build
```

## Running

### Local Development

```bash
# Set environment variables
export DATABASE_URL="host=localhost port=5432 user=woragis password=woragis dbname=woragis sslmode=disable"
export RABBITMQ_URL="amqp://woragis:woragis@localhost:5672/woragis"
export TRANSLATION_PROVIDER=google
export GOOGLE_TRANSLATE_API_KEY="your-api-key"

# Run worker
./bin/translation-worker
```

### Docker

```bash
docker run --rm \
  -e DATABASE_URL="host=db port=5432 user=woragis password=woragis dbname=woragis sslmode=disable" \
  -e RABBITMQ_URL="amqp://woragis:woragis@rabbitmq:5672/woragis" \
  -e TRANSLATION_PROVIDER=google \
  -e GOOGLE_TRANSLATE_API_KEY="your-api-key" \
  translation-worker:latest
```

## Health Check

The worker exposes a health check endpoint at `GET /healthz` on port 8080.

```bash
curl http://localhost:8080/healthz
```

See [HEALTH_CHECK.md](./HEALTH_CHECK.md) for details.

## Logging

The worker uses structured logging. See [LOGGING.md](./LOGGING.md) for details.

## Testing

```bash
# Run all tests
make test

# Run with coverage
make test-cov

# Check coverage threshold (70%)
make test-cov-check
```

## Integration with Server

The server should publish translation jobs to the RabbitMQ queue:

```go
// In server code
job := &TranslationJob{
    ID:         uuid.New().String(),
    EntityType: "project",
    EntityID:   projectID.String(),
    Language:   "pt-BR",
    Fields:     []string{"name", "description"},
    SourceText: map[string]string{
        "name":        project.Name,
        "description": project.Description,
    },
}
translationQueue.EnqueueJob(ctx, job)
```

## Performance Considerations

- **Concurrency:** Set `TRANSLATION_PREFETCH_COUNT` to control concurrent job processing
- **Retries:** Configure `TRANSLATION_MAX_RETRIES` and `TRANSLATION_RETRY_DELAY` based on API reliability
- **Timeout:** Adjust `TRANSLATION_TIMEOUT` based on API response times
- **Scaling:** Run multiple worker instances for horizontal scaling

## Error Handling

- **Invalid messages:** Rejected without requeue
- **Translation failures:** Requeued for retry (up to max retries)
- **Database errors:** Requeued for retry
- **API errors:** Retried with exponential backoff

## Migration from Server

The translation worker has been moved from `server/app/cmd/translation-worker` to `backend/translation-worker` as a standalone service. The server should now only publish jobs to RabbitMQ, not process them directly.
