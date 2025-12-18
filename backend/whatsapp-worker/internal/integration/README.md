# WhatsApp Worker Integration Tests

This directory contains integration tests for the WhatsApp worker service.

## Prerequisites

1. **Docker and Docker Compose** - Required to run test dependencies
2. **RabbitMQ** - Message queue for WhatsApp message processing

## Running Tests

### Using Docker Compose (Recommended)

From the `backend/server` directory:

```bash
# Start test dependencies
docker-compose -f docker-compose.test.yml up -d

# Run tests
cd ../../whatsapp-worker
go test -tags=integration ./internal/integration/...

# Cleanup
cd ../server
docker-compose -f docker-compose.test.yml down
```

### Manual Setup

1. Start RabbitMQ locally:
   ```bash
   docker run -d --name test-rabbitmq \
     -p 5673:5672 \
     -e RABBITMQ_DEFAULT_USER=test \
     -e RABBITMQ_DEFAULT_PASS=test \
     -e RABBITMQ_DEFAULT_VHOST=test \
     rabbitmq:3.13-management-alpine
   ```

2. Set environment variables:
   ```bash
   export RABBITMQ_URL="amqp://test:test@localhost:5673/test"
   ```

3. Run tests:
   ```bash
   go test -tags=integration ./internal/integration/...
   ```

## Test Coverage

### Queue Setup Tests
- ✅ `TestWhatsAppWorkerQueueSetup` - Tests queue and exchange creation
- ✅ `TestWhatsAppWorkerMessagePublish` - Tests publishing messages to queue

### Message Processing Tests
- ✅ `TestWhatsAppWorkerMessageConsume` - Tests consuming and processing messages
- ✅ `TestWhatsAppWorkerInvalidMessage` - Tests handling of invalid JSON messages
- ✅ `TestWhatsAppWorkerMissingDestination` - Tests validation of missing destination
- ✅ `TestWhatsAppWorkerRetryOnFailure` - Tests retry behavior on send failure
- ✅ `TestWhatsAppWorkerMultipleMessages` - Tests processing multiple messages

### Dead Letter Queue Tests
- ✅ `TestWhatsAppWorkerDeadLetterQueue` - Tests DLQ setup and binding

## Test Structure

Tests use a mock WhatsApp notifier to avoid actually sending WhatsApp messages during testing. The mock notifier:
- Tracks sent messages
- Can simulate failures for retry testing
- Validates message content and destination

## Environment Variables

- `RABBITMQ_URL` - RabbitMQ connection URL (default: `amqp://test:test@localhost:5673/test`)
