# Backend Database Join Strategies

## Overview
How database joins are implemented across domains.

## Key Points

### Join Patterns Found

#### 1. Conversations with Messages
```go
Joins("JOIN conversations ON conversations.id = messages.conversation_id")
```
- Used in chats domain
- Retrieves messages with conversation context
- Filters by conversation properties

#### 2. Conversation Transcripts
```go
Joins("JOIN conversations ON conversations.id = conversation_transcripts.conversation_id")
```
- Links transcripts to conversations
- Enables transcript queries filtered by conversation

#### 3. Projects with Milestones
```go
Joins("JOIN projects ON projects.id = milestones.project_id")
```
- Links milestones to projects
- Filters milestones by project properties

#### 4. Projects with Kanban Columns
```go
Joins("JOIN projects ON projects.id = kanban_columns.project_id")
```
- Links kanban columns to projects
- Enables project-scoped kanban operations

#### 5. Resume Metrics (Subqueries)
- Subqueries for counting applications with completed interviews
- Subqueries for counting applications with offers
- Complex aggregations without explicit JOINs

### Join Usage Context

#### Performance Considerations
- Used when filtering by related table
- Used when need data from multiple tables
- GORM handles join optimization

#### Alternative Patterns
- Preload for eager loading relationships
- Separate queries + in-memory join (when appropriate)
- Subqueries for aggregations

## Potential Improvements
- Add explicit JOIN indexes
- Analyze join query performance
- Consider denormalization for frequently joined data
- Implement query result caching for expensive joins
- Add join query monitoring
- Use database views for complex joins
- Implement join query optimization hints
- Add join query pagination
- Support LEFT JOIN vs INNER JOIN where appropriate
- Add join query result size limits
- Implement join query batching for large datasets
- Add join query explain plans

