# Chats Domain - Database Queries

## Overview
Database query patterns used in the chats domain.

## Key Points

### Query Patterns

#### Conversation Queries
- Get conversation by ID + user_id
- List conversations for user
- Search conversations (full-text search)
- Filter by linked entities (idea_id, project_id, job_application_id)
- Archive/restore operations

#### Message Queries
- Get messages for conversation
- Join queries: `JOIN conversations ON conversations.id = messages.conversation_id`
- Filter by conversation_id
- Order by created_at

#### Transcript Queries
- Get conversation transcripts
- Join queries: `JOIN conversations ON conversations.id = conversation_transcripts.conversation_id`
- TTL-based cleanup (7-day retention)
- Transcript search functionality

### Join Patterns
- Messages joined with conversations for filtering
- Transcripts joined with conversations
- Efficient retrieval of related data

### Search Functionality
- Full-text search across conversation titles and messages
- Include archived conversations option
- Limit and pagination support

## Potential Improvements
- Add indexes for frequently queried fields (user_id, conversation_id)
- Optimize join queries
- Add query result caching
- Implement query result pagination
- Add query logging for slow queries
- Optimize search queries (full-text indexes)
- Add advanced search filters (date range, entity type, etc.)
- Implement query result streaming for large datasets
- Add query explain plans for optimization
- Support search result highlighting
- Add search analytics

