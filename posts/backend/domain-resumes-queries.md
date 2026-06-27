# Resumes Domain - Database Queries & Joins

## Overview
Database query patterns and joins used in the resumes domain.

## Key Points

### Query Patterns

#### Simple Queries
- User-scoped queries (filter by user_id)
- Get resume by ID + user_id
- List resumes for user (ordered by: is_main DESC, is_featured DESC, created_at DESC)
- List resumes by tags (must have at least one matching tag)

#### Complex Aggregations
- CalculateResumeMetrics: Complex aggregations using subqueries
- ApplicationsUsed: Count applications linked to resume
- InterviewCount: Count applications with completed interviews (subquery)
- OfferCount: Count applications with offers or status='accepted' (subquery)
- InterviewRate, OfferRate: Calculated percentages

#### Subquery Patterns
- `SELECT DISTINCT job_application_id FROM job_application_stages WHERE completed_date IS NOT NULL`
- `SELECT DISTINCT job_application_id FROM job_application_responses WHERE response_type = 'offer'`
- Used for filtering applications with specific characteristics

### Metric Calculations
- Resume metrics cached in database (UpdateResumeMetrics)
- Recalculation triggered on demand
- Metrics include: application counts, interview rates, offer rates

### Tag Filtering
- ListResumesByTags: Resume must have at least one matching tag
- Uses JSON array operations (GORM)

## Potential Improvements
- Add indexes for frequently queried fields (user_id, is_main, is_featured)
- Optimize subquery performance
- Add query result caching
- Implement pagination for large result sets
- Add query logging for slow queries
- Optimize tag filtering queries
- Add full-text search for resume titles
- Support advanced filtering (date ranges, file size, etc.)
- Implement query result streaming for large datasets
- Add query explain plans for optimization

