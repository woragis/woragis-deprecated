# Job Applications Domain - Database Queries & Joins

## Overview
Database query patterns used in the job applications domain.

## Key Points

### Query Patterns

#### Simple Queries
- User-scoped queries (filter by user_id)
- Get application by ID
- List applications with filters (status, company, date range)
- Filter by resume_id (applications using specific resume)
- Filter by job_application_id (for related entities)

#### Filtering Patterns
- Status-based filtering
- Date range filtering (created_at, applied_at, deadline)
- Company name filtering
- Website filtering
- Resume-based filtering

#### Update Operations
- Status updates
- Field updates (salary, notes, tags, etc.)
- Cascade updates (unlink conversations on delete)

### Integration with Other Domains

#### Resumes Domain
- Applications linked to resumes
- Resume metrics calculation uses applications
- Applications filtered by resume_id

#### Chats Domain
- Applications can have linked conversations
- Conversations unlinked when application deleted

#### User Preferences
- Default language and currency fetched from preferences
- Used for application creation defaults

### Complex Queries
- Applications with completed interviews (for resume metrics)
- Applications with offers (for resume metrics)
- Applications by status aggregation

## Potential Improvements
- Add indexes for frequently queried fields (user_id, status, company, website)
- Add full-text search for company names, job titles
- Optimize filtering queries
- Add query result pagination
- Implement query result caching
- Add query logging for slow queries
- Support advanced filtering (salary range, location, etc.)
- Add query result sorting options
- Implement query result aggregation (statistics)
- Add query explain plans for optimization
- Support query result export

