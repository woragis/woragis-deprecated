# Finances Domain - Database Queries

## Overview
Database query patterns used in the finances domain.

## Key Points

### Query Patterns

#### Transaction Queries
- User-scoped queries (filter by user_id)
- ListTransactions: Date range queries
- QueryTransactions: Advanced filtering (category, type, tags, date range)
- Get transaction by ID
- Filter by archived status

#### Summary Queries
- GetSummary: Aggregation queries
- Group by category, type, date
- Currency conversion calculations
- Income vs expense calculations
- Balance calculations

#### Template Queries
- ListTemplates: User's recurring templates
- Get template by ID
- Filter by active/inactive

#### Cashflow Projection Queries
- Project future transactions from templates
- Aggregate scheduled transactions
- Date range projections

#### Bulk Operations
- Bulk updates with transactions
- Bulk deletes with validation
- Efficient batch processing

### Aggregation Patterns
- SUM aggregations for amounts
- COUNT aggregations for counts
- GROUP BY for categorizations
- Date range filtering
- Currency grouping

### Filtering Patterns
- By date range (from, to)
- By category
- By type (income, expense, transfer)
- By tags (array filtering)
- By archived status
- By recurring flag
- By essential flag

## Potential Improvements
- Add indexes for frequently queried fields (user_id, occurred_at, category, type)
- Optimize aggregation queries
- Add query result caching
- Implement query result pagination
- Add query logging for slow queries
- Optimize date range queries
- Add full-text search for descriptions
- Support advanced filtering (amount ranges, currency, etc.)
- Implement query result streaming
- Add query explain plans
- Support query result export
- Add query analytics
- Implement materialized views for summaries

