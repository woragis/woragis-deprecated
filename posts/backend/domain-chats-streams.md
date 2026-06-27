# Chats Domain - Stream Handling

## Overview
How AI response streaming is implemented for chat conversations.

## Key Points

### Streaming Architecture
- StreamHub integration for WebSocket broadcasting
- AI service streaming responses
- Incremental message updates
- Real-time token-by-token delivery

### Streaming Flow
1. User sends message via API
2. Message appended to conversation
3. AI service called with streaming enabled
4. Stream tokens received incrementally
5. Each token broadcast to WebSocket subscribers
6. Final message saved when stream completes

### Context Building
- ContextBuilder builds conversation context
- Includes user profile data (experiences, projects, skills)
- Includes conversation history
- Includes linked entities (ideas, projects, job applications)

### AI Integration
- LangChain client for AI interactions
- Multiple provider support (OpenAI, Anthropic, etc.)
- Model configuration per conversation
- Temperature and token limits configurable

### Message Management
- Messages stored in database
- Conversation transcripts with TTL (7 days)
- Message search functionality
- Message archiving support

## Potential Improvements
- Add streaming error recovery
- Implement streaming timeout handling
- Add streaming progress indicators
- Support multiple concurrent streams
- Add streaming cancellation
- Implement streaming queue (rate limiting)
- Add streaming metrics (tokens/sec)
- Support streaming for multiple AI providers
- Add streaming content filtering
- Implement streaming response caching
- Add streaming analytics

