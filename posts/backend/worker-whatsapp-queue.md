# WhatsApp Worker Queue Strategy

## Overview
How the WhatsApp worker processes WhatsApp notifications via Redis Pub/Sub.

## Key Points

### Queue Architecture
- Redis Pub/Sub pattern (not traditional queue)
- Channel: `whatsapp` (whatsappChannel constant)
- Event-driven architecture
- Real-time notification delivery

### Pub/Sub Pattern
- Publisher sends messages to channel
- Worker subscribes to channel
- Messages broadcast to all subscribers
- No message persistence (fire-and-forget)

### Message Structure (ReportEnvelope)
- Destination: WhatsApp number
- TextMessage: Message text
- UserID: User identifier
- (Subject, HTMLMessage not used for WhatsApp)

### Processing Flow
1. Subscribe to WhatsApp channel
2. Receive message from channel
3. Parse JSON envelope
4. Extract destination and text message
5. Send via WhatsApp notifier
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
- WhatsApp service: Notifier interface
- Redis client for Pub/Sub
- Context for cancellation

## Potential Improvements
- Add message persistence (store in queue on failure)
- Implement retry mechanism for failed sends
- Add message acknowledgment
- Implement dead letter queue
- Add rate limiting per user/number
- Support message templates
- Add message delivery status tracking
- Implement message read receipts
- Support media attachments (images, documents)
- Add message priority levels
- Implement message scheduling
- Add message batching for efficiency
- Support message personalization
- Add message analytics (delivery, read rates)
- Support WhatsApp Business API features

