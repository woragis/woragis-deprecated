# Projects Domain - Database Queries

## Overview
Database query patterns and joins used in the projects domain.

## Key Points

### Query Patterns

#### Project Queries
- User-scoped queries (filter by user_id)
- Get project by slug
- Search projects by slug (partial matching)
- List projects for user
- Project with all related entities (preload)

#### Join Patterns
- Projects with milestones: `JOIN projects ON projects.id = milestones.project_id`
- Projects with kanban columns: `JOIN projects ON projects.id = kanban_columns.project_id`
- Used for filtering related entities by project properties

#### Milestone Queries
- List milestones for project
- Filter milestones by project_id
- Join with projects for validation

#### Kanban Queries
- Get kanban board for project
- Load columns and cards together
- Order by position fields
- Transaction-based updates for board state

#### Technology Queries
- List technologies for project
- Filter by project_id
- Bulk operations (create/update multiple)

#### Documentation Queries
- Get documentation for project
- Get public documentation by slug
- List documentation sections (ordered)
- Reorder sections operation

### Complex Operations
- Project duplication (deep copy with all related entities)
- Bulk updates with transactions
- Position reordering operations

## Potential Improvements
- Add indexes for frequently queried fields (user_id, slug, project_id)
- Optimize join queries
- Add query result caching
- Implement query result pagination
- Add query logging for slow queries
- Optimize slug search queries
- Add full-text search for project names/descriptions
- Support advanced filtering (status, technologies, date range)
- Implement query result streaming for large datasets
- Add query explain plans for optimization
- Support project search across all fields
- Add query result aggregation (statistics)

