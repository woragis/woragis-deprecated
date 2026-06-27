# Backend Caching - Redis Integration

## Overview
How Redis is integrated and used for caching across the backend.

## Key Points

### Current Redis Usage
- **Queues**: Job queues for workers (translations, resumes, job-applications)
- **Pub/Sub**: Real-time notifications (email, WhatsApp)
- **Job Storage**: Temporary job data with TTL
- **Not extensively used for caching yet** (per caching strategies doc)

### Redis Configuration
- Redis client integration in main server
- Connection pooling
- Context-aware operations
- Error handling for Redis unavailability

### Queue Patterns
- LPush/BRPop for FIFO queues
- Job data stored separately from queue (key-value)
- TTL-based expiration (24 hours, 7 days)
- Pub/Sub for events and notifications

### Potential Caching Patterns

#### Cache-Aside Pattern
- Check Redis cache
- If miss: fetch from database, store in cache
- Return data

#### Write-Through Pattern
- Write to database
- Update cache
- Return success

#### Cache Keys
- Translation: `cache:translation:{entity_type}:{entity_id}:{language}`
- User Profile: `cache:user:{user_id}:profile`
- Resume: `cache:resume:{user_id}:{job_hash}`
- API Response: `cache:api:{endpoint}:{params_hash}`

## Potential Improvements
- Implement caching layer abstraction
- Add cache warming strategies
- Implement cache invalidation patterns
- Add cache hit/miss metrics
- Support cache versioning
- Implement cache compression
- Add cache TTL strategies
- Support cache tags (invalidate by tag)
- Implement cache locking (prevent thundering herd)
- Add cache monitoring and alerting
- Support distributed caching (Redis Cluster)
- Implement cache partitioning
- Add cache analytics
- Support multi-level caching
- Add cache encryption for sensitive data

