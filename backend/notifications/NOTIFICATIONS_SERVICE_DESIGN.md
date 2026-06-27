# Notifications Service - Design Document

## Overview

The Notifications Service is responsible for managing and orchestrating all notification delivery across the Woragis platform. It acts as a central hub that receives notification requests from other microservices and dispatches them to appropriate worker services (Email Worker, WhatsApp Worker, etc.) via message queues.

## Architecture

### Service Responsibilities

1. **Notification Management**
   - Create, update, and track notification requests
   - Store notification history and delivery status
   - Manage notification templates
   - Handle notification preferences per user

2. **Queue Orchestration**
   - Publish notification jobs to RabbitMQ queues
   - Route notifications to appropriate workers based on channel (email, WhatsApp, push, etc.)
   - Handle retry logic and dead-letter queues
   - Monitor queue health

3. **Template Management**
   - Store and version notification templates
   - Support multi-language templates
   - Template variable substitution
   - Template rendering (HTML, plain text, etc.)

4. **Delivery Tracking**
   - Track notification status (pending, sent, delivered, failed)
   - Store delivery metadata (timestamps, provider responses, etc.)
   - Provide delivery analytics and reporting

## Integration with Existing Workers

### Email Worker
- **Location**: `Services-Workers/email-worker`
- **Queue**: `email_notifications` (or similar)
- **Integration**: Notifications service publishes email notification envelopes to the queue
- **Docker Image**: Use existing email-worker Docker image
- **Message Format**: Follow existing `ReportEnvelope` or similar structure

### WhatsApp Worker
- **Location**: `Services-Workers/whatsapp-worker`
- **Queue**: `whatsapp_notifications` (or similar)
- **Integration**: Notifications service publishes WhatsApp notification envelopes to the queue
- **Docker Image**: Use existing whatsapp-worker Docker image
- **Message Format**: Follow existing envelope structure

## Domain Model

### Core Entities

#### Notification
```go
type Notification struct {
    ID          uuid.UUID
    UserID      uuid.UUID
    Type        NotificationType  // email, whatsapp, push, sms
    Status      NotificationStatus // pending, queued, sent, delivered, failed
    Priority    int
    Subject     string  // For email
    Content     string  // Rendered content
    TemplateID  *uuid.UUID
    Variables   map[string]interface{} // Template variables
    Metadata    map[string]interface{} // Additional metadata
    ScheduledAt *time.Time
    SentAt      *time.Time
    DeliveredAt *time.Time
    FailedAt    *time.Time
    ErrorMessage string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

#### NotificationTemplate
```go
type NotificationTemplate struct {
    ID          uuid.UUID
    Name        string
    Type        NotificationType
    Subject     string  // For email
    Body        string  // Template body
    Variables   []string // Required variables
    Language    string  // ISO 639-1
    IsActive    bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

#### NotificationPreference
```go
type NotificationPreference struct {
    ID              uuid.UUID
    UserID          uuid.UUID
    NotificationType string
    Channel         string  // email, whatsapp, push
    Enabled         bool
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

## API Endpoints

### Notification Management
- `POST /api/v1/notifications` - Create and send notification
- `GET /api/v1/notifications` - List notifications (with filters)
- `GET /api/v1/notifications/:id` - Get notification details
- `GET /api/v1/notifications/:id/status` - Get delivery status
- `POST /api/v1/notifications/bulk` - Send bulk notifications

### Template Management
- `POST /api/v1/templates` - Create template
- `GET /api/v1/templates` - List templates
- `GET /api/v1/templates/:id` - Get template
- `PATCH /api/v1/templates/:id` - Update template
- `DELETE /api/v1/templates/:id` - Delete template
- `POST /api/v1/templates/:id/render` - Preview rendered template

### Preferences
- `GET /api/v1/preferences` - Get user notification preferences
- `PATCH /api/v1/preferences` - Update preferences
- `GET /api/v1/preferences/channels` - Get available channels

### Webhooks (for worker callbacks)
- `POST /api/v1/webhooks/email/delivery` - Email worker delivery callback
- `POST /api/v1/webhooks/whatsapp/delivery` - WhatsApp worker delivery callback
- `POST /api/v1/webhooks/email/bounce` - Email bounce callback
- `POST /api/v1/webhooks/whatsapp/status` - WhatsApp status callback

## Message Queue Integration

### RabbitMQ Queues

#### Email Queue
- **Queue Name**: `email_notifications`
- **Exchange**: `notifications`
- **Routing Key**: `email`
- **Message Format**: 
  ```json
  {
    "notification_id": "uuid",
    "user_id": "uuid",
    "to": "email@example.com",
    "subject": "string",
    "body": "string",
    "html_body": "string (optional)",
    "metadata": {}
  }
  ```

#### WhatsApp Queue
- **Queue Name**: `whatsapp_notifications`
- **Exchange**: `notifications`
- **Routing Key**: `whatsapp`
- **Message Format**:
  ```json
  {
    "notification_id": "uuid",
    "user_id": "uuid",
    "to": "phone_number",
    "message": "string",
    "template_id": "string (optional)",
    "metadata": {}
  }
  ```

### Queue Configuration
- Use existing RabbitMQ instance
- Configure dead-letter queues for failed notifications
- Set appropriate TTL and retry policies
- Monitor queue depth and processing rates

## Service Dependencies

### Required Services
1. **Auth Service** - For user authentication and JWT validation
2. **RabbitMQ** - Message queue for worker communication
3. **PostgreSQL** - For notification and template storage
4. **Redis** - For caching templates and rate limiting

### External Workers (via Docker)
1. **Email Worker** - `email-worker` Docker image
2. **WhatsApp Worker** - `whatsapp-worker` Docker image

## Implementation Considerations

### 1. Template Rendering
- Support multiple template engines (Go templates, Handlebars, etc.)
- Variable validation before rendering
- Multi-language template support
- Template caching in Redis

### 2. Rate Limiting
- Per-user rate limits
- Per-channel rate limits
- Global rate limits
- Respect user preferences (opt-out)

### 3. Retry Logic
- Exponential backoff for failed notifications
- Maximum retry attempts
- Dead-letter queue handling
- Manual retry endpoint

### 4. Delivery Tracking
- Webhook endpoints for worker callbacks
- Status updates from workers
- Delivery confirmation tracking
- Bounce and error handling

### 5. Security
- Validate notification requests
- Sanitize template content
- Rate limiting to prevent abuse
- Webhook signature verification

### 6. Scalability
- Horizontal scaling support
- Queue-based architecture for async processing
- Database connection pooling
- Caching layer for templates

## Database Schema

### Tables
- `notifications` - Notification records
- `notification_templates` - Template storage
- `notification_preferences` - User preferences
- `notification_delivery_logs` - Delivery history and status

### Indexes
- `notifications.user_id` - For user notification queries
- `notifications.status` - For status-based queries
- `notifications.created_at` - For time-based queries
- `notification_templates.name` - For template lookup
- `notification_preferences.user_id` - For preference queries

## Environment Variables

```env
# Service Configuration
NOTIFICATIONS_SERVICE_PORT=3006
ENV=development

# Database
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=woragis
POSTGRES_PASSWORD=password
POSTGRES_DB=notifications_db

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
RABBITMQ_EXCHANGE=notifications
RABBITMQ_EMAIL_QUEUE=email_notifications
RABBITMQ_WHATSAPP_QUEUE=whatsapp_notifications

# Auth Service
AUTH_SERVICE_URL=http://auth-service:3000

# Worker URLs (for webhooks)
EMAIL_WORKER_URL=http://email-worker:3000
WHATSAPP_WORKER_URL=http://whatsapp-worker:3000

# Rate Limiting
RATE_LIMIT_PER_USER=100  # notifications per hour
RATE_LIMIT_PER_CHANNEL=50  # per channel per hour
```

## Docker Compose Integration

### Service Definition
```yaml
notifications-service:
  build: ./notifications/server
  ports:
    - "3006:3006"
  environment:
    - AUTH_SERVICE_URL=http://auth-service:3000
    - RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
  depends_on:
    - postgres
    - redis
    - rabbitmq
    - auth-service

email-worker:
  image: woragis/email-worker:latest  # Use existing image
  environment:
    - RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
  depends_on:
    - rabbitmq
    - notifications-service

whatsapp-worker:
  image: woragis/whatsapp-worker:latest  # Use existing image
  environment:
    - RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
  depends_on:
    - rabbitmq
    - notifications-service
```

## Migration Strategy

### Phase 1: Service Setup
1. Create notifications service structure (copy from auth/jobs pattern)
2. Set up database schema
3. Implement basic notification creation and storage
4. Set up RabbitMQ integration

### Phase 2: Template System
1. Implement template CRUD operations
2. Add template rendering engine
3. Implement variable substitution
4. Add template caching

### Phase 3: Queue Integration
1. Integrate with existing email-worker
2. Integrate with existing whatsapp-worker
3. Implement webhook endpoints for delivery callbacks
4. Add retry and error handling

### Phase 4: Advanced Features
1. Notification preferences
2. Delivery tracking and analytics
3. Rate limiting
4. Scheduled notifications

## Testing Considerations

### Unit Tests
- Template rendering
- Variable substitution
- Queue message formatting
- Validation logic

### Integration Tests
- RabbitMQ message publishing
- Worker communication
- Webhook callbacks
- Database operations

### E2E Tests
- Full notification flow (create → queue → worker → callback)
- Multi-channel notifications
- Error scenarios

## Monitoring and Observability

### Metrics
- Notification creation rate
- Queue depth
- Delivery success rate
- Average delivery time
- Error rates by channel

### Logging
- Notification creation events
- Queue publishing events
- Delivery status updates
- Error logs with context

### Tracing
- End-to-end notification flow
- Worker communication traces
- Template rendering performance

## Future Enhancements

1. **Push Notifications** - Add push notification support (FCM, APNS)
2. **SMS Notifications** - Add SMS channel support
3. **Notification Scheduling** - Advanced scheduling with timezone support
4. **A/B Testing** - Template A/B testing capabilities
5. **Analytics Dashboard** - Real-time notification analytics
6. **Multi-tenant Support** - Organization-level notification settings

## Notes

- This service should be implemented carefully to avoid duplicating existing worker logic
- Reuse existing worker Docker images rather than reimplementing
- Focus on orchestration and management, not delivery implementation
- Ensure webhook security with signature verification
- Consider notification batching for high-volume scenarios

