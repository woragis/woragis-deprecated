# Email Worker Queue Strategy

## Overview
How the email worker processes email notifications via Redis Pub/Sub.

## Key Points

### Queue Architecture
- Redis Pub/Sub pattern (not traditional queue)
- Channel: `email` (emailChannel constant)
- Event-driven architecture
- Real-time notification delivery

### Pub/Sub Pattern
- Publisher sends messages to channel
- Worker subscribes to channel
- Messages broadcast to all subscribers
- No message persistence (fire-and-forget)

### Message Structure (ReportEnvelope)
- Destination: Email address
- Subject: Email subject
- TextMessage: Plain text body
- HTMLMessage: HTML body
- UserID: User identifier

### Processing Flow
1. Subscribe to email channel
2. Receive message from channel
3. Parse JSON envelope
4. Convert to emailservice.Message
5. Send via email sender
6. Handle errors (log, continue)

### Error Handling
- Invalid payload: log error, continue processing
- Send failures: log error with user_id, continue
- Context cancellation: unsubscribe and exit

## Implementation Details

### Redis Operations
- `Subscribe`: Subscribe to channel
- `Channel`: Get message channel
- Message handling in goroutine

### Integration Points
- Email service: Sender interface
- Redis client for Pub/Sub
- Context for cancellation

## Potential Improvements
- Add message persistence (store in queue on failure)
- Implement retry mechanism for failed sends
- Add message acknowledgment
- Implement dead letter queue
- Add rate limiting per user/email
- Support email templates
- Add email delivery status tracking
- Implement email bounce handling
- Add email unsubscription support
- Support email attachments
- Add email priority levels
- Implement email scheduling
- Add email batching for efficiency
- Support email personalization
- Add email analytics (open rates, clicks)

