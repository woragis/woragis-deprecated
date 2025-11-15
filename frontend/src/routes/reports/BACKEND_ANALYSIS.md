# Reports Backend Analysis

## Overview

The reports system in the backend is built using Go (Fiber framework) and follows a domain-driven architecture. Here's how it works:

## Architecture

### Core Components

1. **ReportDefinition** - Defines what data to include in a report
2. **ReportSchedule** - Defines when to automatically generate reports (using cron expressions)
3. **ReportDelivery** - Defines how to send/distribute reports
4. **ReportRun** - Tracks individual report generation executions

## Database Schema

### ReportDefinition
- `id` (UUID) - Primary key
- `user_id` (UUID) - Owner of the report
- `name` (string) - Report name
- `description` (string) - Optional description
- `sections` (JSONB) - Defines what sections to include (e.g., `{"overview": []}`)
- `filters` (JSONB) - Defines date ranges and filters (e.g., `{"date_range": "last_30_days"}`)
- `is_favorite` (boolean) - Favorite flag
- `archived_at` (timestamp) - Soft archive support
- `created_at`, `updated_at`, `deleted_at` - Standard timestamps

### ReportSchedule
- `id` (UUID) - Primary key
- `report_id` (UUID) - Links to ReportDefinition
- `cron` (string) - Cron expression (e.g., `"0 8 * * *"` for daily at 8 AM)
- `frequency` (string) - Human-readable frequency (e.g., "daily", "weekly", "monthly")
- `timezone` (string) - Timezone for scheduling (e.g., "America/New_York")
- `next_run` (timestamp) - When the next execution should occur
- `last_run_at` (timestamp) - Last execution time
- `enabled` (boolean) - Whether the schedule is active
- `meta` (JSONB) - Additional metadata

### ReportDelivery
- `id` (UUID) - Primary key
- `report_id` (UUID) - Links to ReportDefinition
- `channel` (string) - Delivery method (e.g., "email", "whatsapp")
- `target` (string) - Recipient (e.g., email address)
- `template` (JSONB) - Custom template for formatting
- `enabled` (boolean) - Whether delivery is active

### ReportRun
- `id` (UUID) - Primary key
- `report_id` (UUID) - Links to ReportDefinition
- `status` (string) - "pending", "running", "completed", "failed"
- `output_location` (string) - Where the generated report is stored
- `error_message` (string) - Error details if failed
- `metadata` (JSONB) - Additional run metadata
- `created_at`, `updated_at` - Timestamps

## API Endpoints

### Definitions
- `POST /reports` - Create a new report definition
- `GET /reports` - List report definitions (with filters)
- `GET /reports/:id` - Get definition detail (includes schedules and deliveries)
- `PUT /reports/:id` - Update a definition
- `POST /reports/archive` - Archive definitions (bulk)
- `POST /reports/restore` - Restore archived definitions (bulk)
- `POST /reports/delete` - Delete definitions (bulk)
- `POST /reports/favorite` - Toggle favorite status

### Schedules
- `POST /reports/:id/schedules` - Create a schedule for a report
- `GET /reports/:id/schedules` - List schedules for a report
- `PUT /reports/schedules/:scheduleID` - Update a schedule
- `POST /reports/schedules/:scheduleID/toggle` - Enable/disable a schedule
- `DELETE /reports/schedules/:scheduleID` - Delete a schedule

### Deliveries
- `POST /reports/:id/deliveries` - Create a delivery configuration
- `GET /reports/:id/deliveries` - List deliveries for a report
- `PUT /reports/deliveries/:deliveryID` - Update a delivery
- `POST /reports/deliveries/:deliveryID/toggle` - Enable/disable a delivery
- `DELETE /reports/deliveries/:deliveryID` - Delete a delivery

### Runs
- `POST /reports/runs/bulk` - Queue report runs for multiple definitions
- `GET /reports/:id/runs` - List run history for a report

## Cron Expression Format

The backend uses standard cron expressions with 5 fields:
```
minute hour day-of-month month day-of-week
```

**Examples:**
- `0 8 * * *` - Daily at 8:00 AM
- `0 9 * * 1` - Every Monday at 9:00 AM (1 = Monday, 0 = Sunday)
- `0 6 1 * *` - 1st of every month at 6:00 AM
- `0 */7 * * *` - Every 7 days at midnight (approximation)

**Day of week values:**
- 0 = Sunday
- 1 = Monday
- 2 = Tuesday
- 3 = Wednesday
- 4 = Thursday
- 5 = Friday
- 6 = Saturday

## How Scheduling Works

1. **Schedule Creation**: When a schedule is created, the cron expression is stored along with timezone and frequency metadata.

2. **Next Run Calculation**: The `next_run` field can be set manually or calculated based on the cron expression and timezone.

3. **Execution**: The backend stores the cron expressions, but actual execution would typically be handled by:
   - A background worker/daemon that periodically checks for due schedules
   - An external cron scheduler
   - A job queue system (like the scheduler domain in the codebase)

4. **Report Generation**: When a schedule triggers:
   - A `ReportRun` is created with status "pending"
   - The report is generated based on the definition's sections and filters
   - The status is updated to "running", then "completed" or "failed"
   - Deliveries are processed based on enabled delivery configurations

## Data Flow

1. **User creates a ReportDefinition** with sections and filters
2. **User adds a ReportSchedule** with cron expression and timezone
3. **User adds ReportDelivery** configurations (email, etc.)
4. **Scheduler/Worker** checks for due schedules (based on `next_run` or cron evaluation)
5. **Report is generated** using the definition's sections and filters
6. **ReportRun** is created and tracked
7. **Deliveries are executed** based on enabled delivery channels
8. **Next run is calculated** and schedule is updated

## Key Features

- **Soft Deletes**: Definitions, schedules, and deliveries use soft deletes (`deleted_at`)
- **Archiving**: Definitions can be archived without deletion
- **Favorites**: Users can mark definitions as favorites
- **Bulk Operations**: Archive, restore, delete, and queue runs support bulk operations
- **Metadata**: Schedules and runs support custom metadata (JSONB)
- **Timezone Support**: Schedules respect timezone settings
- **Enable/Disable**: Schedules and deliveries can be toggled without deletion

## Integration Points

The reports service integrates with:
- **Ideas Domain**: Can include ideas in reports
- **Projects Domain**: Can include projects in reports
- **Finances Domain**: Can aggregate financial summaries
- **Chats Domain**: Can include conversation data
- **Notification System**: Publishes reports via email/WhatsApp

## Notes

- The cron expressions are stored as strings and validated on creation
- The `frequency` field is metadata for UI purposes and doesn't affect execution
- The backend doesn't include a cron parser/executor - this would be handled by a separate scheduler service or worker
- Report generation aggregates data from multiple domains based on the definition's sections and filters

