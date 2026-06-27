# Backend Database Query Patterns

## Overview
Common query patterns and strategies used across domains.

## Key Points

### Query Patterns

#### 1. Simple CRUD Queries
- Standard GORM queries (Find, Create, Update, Delete)
- User-scoped queries (filter by user_id)
- Soft deletes support

#### 2. Join Queries
- Projects with milestones: `JOIN milestones ON milestones.project_id = projects.id`
- Conversations with messages: `JOIN conversations ON conversations.id = messages.conversation_id`
- Used for related data fetching

#### 3. Aggregation Queries
- Resume metrics (counts of applications, interviews, offers)
- Subqueries for filtering (e.g., completed interview stages)
- Complex aggregations in resume domain

#### 4. Filtered Queries
- Tag-based filtering (resumes by tags)
- Status-based filtering (applications by status)
- Date range filtering

#### 5. Preloading
- Eager loading of relationships (GORM Preload)
- Nested preloading for deep relationships
- Reduces N+1 query problems

### Domain-Specific Patterns

#### Resumes Domain
- Complex aggregations for metrics
- Subqueries for filtering applications
- Joins with job_applications table

#### Chats Domain
- Join queries for messages and conversations
- Join queries for transcripts
- Efficient message retrieval

#### Projects Domain
- Joins with milestones
- Joins with kanban columns
- Nested relationship loading

## Query Optimization Strategies

### Current Practices
- Use of indexes (implied by GORM model definitions)
- User-scoped queries (security + performance)
- Selective field loading

### Potential Optimizations
- Add explicit indexes for frequently queried fields
- Use query hints for complex joins
- Implement query result caching
- Add query logging for slow queries
- Use database views for complex aggregations
- Implement pagination for large result sets
- Add query timeout configuration
- Use prepared statements for repeated queries

## Potential Improvements
- Add query performance monitoring
- Implement slow query logging (> 100ms)
- Add query explain plans for optimization
- Implement query result pagination consistently
- Add database connection pooling optimization
- Implement read replicas for read-heavy operations
- Add query result caching (Redis)
- Support query batching for bulk operations
- Implement database sharding for scale
- Add query metrics collection
- Support database query rewriting for optimization
- Implement query result streaming for large datasets
- Add query sanitization for security
- Support database query versioning (migration-safe queries)

