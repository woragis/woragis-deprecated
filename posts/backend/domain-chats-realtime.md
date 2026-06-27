# Chats Domain - Real-time Communication

## Overview
How real-time chat communication is implemented using WebSockets.

## Key Points

### WebSocket Architecture
- StreamHub manages WebSocket connections
- Subscribers per conversation (map[conversationID]map[connection])
- Broadcast to all subscribers of a conversation
- Thread-safe operations (sync.RWMutex)

### Connection Management
- **Register**: Add connection to conversation's subscriber list
- **Unregister**: Remove connection from subscriber list
- **Broadcast**: Send event to all subscribers of conversation
- Automatic cleanup on connection close

### Real-time Events
- Message updates broadcast to all subscribers
- Conversation updates broadcast
- AI response streaming via WebSockets
- Status updates (typing indicators, etc.)

### StreamHub Implementation
- Concurrent-safe operations
- Connection pool management per conversation
- Automatic connection cleanup on errors
- JSON message format

## Potential Improvements
- Add connection heartbeat/ping-pong
- Implement connection authentication
- Add connection rate limiting
- Support room-based broadcasting
- Add presence tracking (who's online)
- Implement message delivery acknowledgments
- Add connection reconnection handling
- Support binary message formats
- Add connection load balancing
- Implement connection metrics
- Support private messaging
- Add typing indicators
- Implement read receipts
- Add message reactions

