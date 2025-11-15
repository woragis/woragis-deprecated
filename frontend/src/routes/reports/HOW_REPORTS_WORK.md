# How Reports Work in the Backend

## Overview

This document explains how reports are generated, what data they use, and how they're delivered in the Woragis backend.

## Report Generation Process

### Current Implementation

The report generation is handled by the `GenerateSummary` function in `backend/server/app/internal/domains/reports/service.go`. Here's how it works:

#### 1. **Data Aggregation (No AI, Uses SQL Queries)**

The system does **NOT** use AI to generate SQL queries. Instead, it uses **GORM (Go ORM)** which generates SQL queries automatically. The report generation process:

```go
func (s *Service) GenerateSummary(ctx context.Context, userID uuid.UUID) (Summary, error) {
    // 1. Fetch Ideas - uses repository method (generates SQL)
    ideas, err = s.ideasRepo.ListIdeas(ctx, userID)
    
    // 2. Fetch Projects - uses repository method (generates SQL)
    projects, err = s.projectsRepo.ListProjects(ctx, userID)
    
    // 3. Fetch Chats/Conversations - uses repository method (generates SQL)
    chats, err = s.chatsRepo.ListConversations(ctx, userID)
    
    // 4. Aggregate Finances - uses SQL aggregation queries
    finances, err = s.financeRepo.AggregateSummary(ctx, userID, from, to)
    
    // 5. Compile into Summary struct
    return Summary{
        IdeaCount: len(ideas),
        ProjectCount: len(projects),
        ConversationCount: len(chats),
        IncomeTotal: finances.IncomeTotal,
        ExpenseTotal: finances.ExpenseTotal,
        SavingsAllocation: finances.SavingsAllocation,
    }
}
```

#### 2. **SQL Queries (Example: Finance Aggregation)**

The finance aggregation uses **direct SQL aggregation queries** via GORM:

```go
// From finances/repository.go
func (r *gormRepository) AggregateSummary(ctx context.Context, userID uuid.UUID, from, to time.Time) (Summary, error) {
    query := r.db.WithContext(ctx).
        Model(&Transaction{}).
        Select("type, SUM(normalized_amount) as total").
        Where("user_id = ? AND is_archived = ?", userID, false).
        Group("type")
    
    if !from.IsZero() {
        query = query.Where("occurred_at >= ?", from)
    }
    if !to.IsZero() {
        query = query.Where("occurred_at <= ?", to)
    }
    
    // This generates SQL like:
    // SELECT type, SUM(normalized_amount) as total 
    // FROM transactions 
    // WHERE user_id = ? AND is_archived = false 
    //   AND occurred_at >= ? AND occurred_at <= ?
    // GROUP BY type
}
```

### How the System Knows Which Relations to Study

The system knows which relations to study because it's **hardcoded** in the `GenerateSummary` function. The report service has dependencies on:

1. **IdeaRepository** - Fetches all ideas for the user
2. **ProjectRepository** - Fetches all projects for the user
3. **ChatRepository** - Fetches all conversations for the user
4. **FinanceRepository** - Aggregates financial data (income, expenses, savings)

These repositories are injected into the service and called in a fixed order. The system doesn't dynamically determine which relations to query - it always queries all of them.

### ⚠️ Current Limitation: Report Definitions Not Used

**Important**: The `sections` and `filters` stored in `ReportDefinition` are **NOT currently used** during report generation. The `GenerateSummary` function ignores the definition and always generates the same summary structure.

The report definitions are stored in the database but serve more as **metadata/templates** for future use. The actual report generation is hardcoded to always include:
- Idea count
- Project count
- Conversation count
- Income total
- Expense total
- Savings allocation

## Report Delivery

### How Reports Are Sent

Reports are delivered through a **publisher pattern** using Redis channels:

1. **Report Generation**: `GenerateSummary` creates a `Summary` struct
2. **Formatting**: The summary is formatted into a text message using `formatSummary()`
3. **Publishing**: The formatted message is published to Redis channels:
   - `reports.email` - For email delivery
   - `reports.whatsapp` - For WhatsApp delivery
4. **Workers**: Background workers listen to these channels and send the actual messages

### Delivery Targets (Email/WhatsApp)

#### Current Implementation

Currently, the delivery targets (email address and phone number) are **passed in the request payload**, not automatically retrieved from user data:

```go
type DispatchOptions struct {
    SendEmail    bool
    EmailAddress string  // ← Must be provided in request
    SendWhatsApp bool
    PhoneNumber  string  // ← Must be provided in request
    AgentAlias   string
}
```

#### User Contact Information

The user's **email is available** in the authentication context:

```go
// From auth/context.go
type RequestUser struct {
    ID    uuid.UUID
    Email string  // ← Available from context
}
```

However, the current implementation requires the email/phone to be explicitly provided in the request. The system does **NOT automatically**:
- Use the user's email from the auth context
- Look up the user's phone number from a user profile

#### Recommendation for Improvement

To automatically use the user's email/phone:

1. **Email**: Can be retrieved from `authdomain.UserEmailFromContext(c)` 
2. **Phone Number**: Would need to be stored in the User entity or a separate profile table

Currently, when creating a delivery configuration in the frontend, you must manually specify the target (email address or phone number).

## Report Scheduling

### How Schedules Work

1. **Schedule Storage**: Schedules are stored with cron expressions in the `report_schedules` table
2. **Scheduler Service**: A separate scheduler service (in `domains/scheduler`) checks for due schedules
3. **Execution**: When a schedule is due:
   - The scheduler calls `reports.GenerateSummary()`
   - The summary is dispatched using the delivery options from the schedule
   - The next run time is calculated and stored

### Scheduler Integration

The scheduler service has its own `Schedule` entity that includes:
- Email address
- Phone number
- Agent alias
- Frequency, weekday, time of day
- RRule for complex schedules

When executing, it calls:
```go
summary, err := s.reports.GenerateSummary(ctx, schedule.UserID)
opts := reportsdomain.DispatchOptions{
    SendEmail:    schedule.Email != "",
    EmailAddress: schedule.Email,
    SendWhatsApp: schedule.PhoneNumber != "",
    PhoneNumber:  schedule.PhoneNumber,
    AgentAlias:   schedule.AgentAlias,
}
s.reports.DispatchSummary(ctx, summary, opts)
```

## Data Flow Summary

```
1. User creates ReportDefinition
   └─> Stored in database (sections/filters saved but not used yet)

2. User creates ReportSchedule
   └─> Cron expression stored, linked to definition

3. Scheduler checks for due schedules
   └─> Finds schedule that's due

4. GenerateSummary() called
   ├─> ideasRepo.ListIdeas() → SQL query
   ├─> projectsRepo.ListProjects() → SQL query
   ├─> chatsRepo.ListConversations() → SQL query
   └─> financeRepo.AggregateSummary() → SQL aggregation query

5. Summary compiled into struct
   └─> Contains counts and totals

6. formatSummary() formats the message
   └─> Uses agent profile for persona/tone

7. DispatchSummary() publishes to Redis
   ├─> PublishEmailReport() → "reports.email" channel
   └─> PublishWhatsAppReport() → "reports.whatsapp" channel

8. Background workers process
   ├─> Email worker sends via SMTP
   └─> WhatsApp worker sends via API
```

## Key Takeaways

1. **No AI for SQL**: Reports use GORM (ORM) which generates SQL queries automatically
2. **Hardcoded Relations**: The system always queries the same 4 domains (ideas, projects, chats, finances)
3. **Definitions Not Used**: Report definition sections/filters are stored but not used in generation
4. **Manual Delivery Targets**: Email/phone must be provided in requests, not auto-retrieved from user data
5. **User Email Available**: User email is in auth context but not automatically used for delivery
6. **Phone Number Missing**: No phone number stored in User entity - would need to be added

## Potential Improvements

1. **Use Report Definitions**: Implement logic to respect `sections` and `filters` from definitions
2. **Auto-retrieve Email**: Use `UserEmailFromContext()` for email delivery
3. **Add Phone to User**: Store phone number in User entity or profile
4. **Dynamic Relations**: Allow report definitions to specify which domains to include
5. **Custom Aggregations**: Allow definitions to specify custom SQL aggregations or filters

