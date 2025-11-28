# Job Application Worker

## Quick Reference

- **Worker Type**: Always running, queue-based (like translation-worker)
- **Job Discovery**: Manual/API (Phase 1) → Automated scraping (Phase 2)
- **Activation**: Main server enqueues jobs via API
- **Technology**: Playwright + AI (cover letter generation)
- **Rate Limiting**: Per-website daily limits with automatic rotation
- **Documentation**: 
  - Main README: This file
  - Playwright & AI: [PLAYWRIGHT_AI_ENHANCEMENTS.md](./PLAYWRIGHT_AI_ENHANCEMENTS.md)

## Overview

The Job Application Worker is an automated system that applies to job postings across multiple employment websites (LinkedIn, Glassdoor, WeWorkRemotely, etc.) using Playwright for browser automation and AI for generating personalized cover letters.

## Architecture

### Worker Pattern

The worker follows the same pattern as other workers in the system (e.g., `translation-worker`):

- **Always Running**: The worker runs continuously, polling a Redis queue for job application tasks
- **Queue-Based**: Main server enqueues jobs via API endpoints
- **Separate Container**: Runs in its own Docker container with Playwright dependencies
- **Domain-Driven**: Follows the domain-driven design pattern used throughout the codebase

### Communication Flow

```
Main Server (API)
    ↓ (HTTP POST /api/job-applications)
Enqueue Job Application
    ↓ (Redis Queue: "job-applications:queue")
Job Application Worker
    ↓ (Polls queue continuously)
Check Rate Limits
    ↓
Launch Playwright Browser
    ↓
Navigate to Job Posting
    ↓
Fetch User Profile Data from DB
    ↓
Call AI Service (Generate Cover Letter)
    ↓
Fill Application Form / Submit
    ↓
Save to Database (JobApplication entity)
    ↓
Mark Job Complete in Queue
```

## Components

### Domain Structure

```
app/internal/domains/
  jobapplications/
    entity.go          # JobApplication model
    errors.go          # Domain-specific errors
    repository.go      # Database operations (GORM)
    service.go         # Business logic
    queue.go           # Redis queue interface
    handler.go         # HTTP handlers (for main server)
    routes.go          # API route definitions
  jobwebsites/
    entity.go          # JobWebsite model (rate limits, config)
    errors.go
    repository.go
    service.go
    handler.go
    routes.go
```

### Worker Structure

```
app/internal/workers/
  jobapplications/
    worker.go          # Main worker loop (polls Redis queue)
    orchestrator.go    # Manages rate limits & website rotation
```

### Services

```
app/internal/services/
  playwright/
    browser.go         # Browser instance management
    scraper.go         # Playwright operations (navigation, form filling)
    stealth.go         # Anti-detection measures
  ai/
    coverletter.go     # AI cover letter generation
```

### Command Entry Point

```
app/cmd/job-application-worker/
  main.go              # Worker entry point (similar to translation-worker)
```

## Database Schema

### JobApplication Entity

```go
type JobApplication struct {
    ID              uuid.UUID
    UserID          uuid.UUID
    CompanyName     string
    Location        string
    JobTitle        string
    JobURL          string
    Website         string          // "linkedin", "glassdoor", "weworkremotely", etc.
    AppliedAt       time.Time
    CoverLetter     string          // AI-generated cover letter text
    LinkedInContact bool            // Whether employee was contacted via LinkedIn
    Status          ApplicationStatus // "pending", "applied", "contacted", "rejected", "accepted"
    ErrorMessage    string          // If application failed
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### JobWebsite Entity

```go
type JobWebsite struct {
    ID           uuid.UUID
    Name         string      // "linkedin", "glassdoor", etc.
    DisplayName string      // "LinkedIn", "Glassdoor"
    DailyLimit   int         // 50 for LinkedIn, etc.
    CurrentCount int        // Track today's applications
    LastReset    time.Time  // When counter was last reset
    Enabled      bool       // Whether this website is active
    BaseURL      string     // "https://www.linkedin.com"
    LoginURL     string     // "https://www.linkedin.com/login"
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

## Job Discovery

**Current Approach: Manual/API (Phase 1)**

- Jobs are discovered manually or via a separate process
- JobApplication records are created via API endpoint: `POST /api/job-applications`
- Worker only handles the application process (not discovery)

**Future Enhancement: Automated Discovery (Phase 2)**

- Worker can be extended to scrape job listings
- Filter by criteria (location, role, salary, etc.)
- Automatically create JobApplication records
- Then apply to discovered jobs

## Rate Limiting & Website Rotation

### Strategy

1. **Per-Website Limits**: Each website has a `DailyLimit` (e.g., LinkedIn: 50/day)
2. **Counter Tracking**: Track current count per website in Redis with TTL (or DB with daily reset)
3. **Rotation Logic**: Worker orchestrator rotates through websites:
   - Check current count for each enabled website
   - If under limit → process jobs for that website
   - If at limit → skip to next website
   - Reset counters at midnight (or use Redis TTL)

### Implementation

```go
// Pseudo-code for orchestrator
for {
    websites := getEnabledWebsites()
    for _, website := range websites {
        todayCount := getTodayCount(website.Name)
        if todayCount < website.DailyLimit {
            jobs := getJobsForWebsite(website.Name)
            for _, job := range jobs {
                processJob(job)
                incrementCount(website.Name)
                if getTodayCount(website.Name) >= website.DailyLimit {
                    break // Move to next website
                }
            }
        }
    }
    sleep(1 hour) // Check again later
}
```

## Playwright Integration

### Why Playwright?

- **JavaScript Support**: Handles modern SPAs (LinkedIn, Glassdoor are React-based)
- **Cross-Browser**: Supports Chromium, Firefox, WebKit
- **Human-Like Interactions**: Can simulate mouse movements, typing delays, scrolling
- **Stealth Capabilities**: Can be configured to avoid detection
- **Screenshot/Video**: Useful for debugging failed applications

### Limitations & AI Enhancements

Playwright has several limitations when automating job application websites. AI can significantly enhance its capabilities. See **[PLAYWRIGHT_AI_ENHANCEMENTS.md](./PLAYWRIGHT_AI_ENHANCEMENTS.md)** for detailed discussion of:

1. **Anti-Bot Detection** - AI-powered behavioral simulation
2. **Page Structure Changes** - Self-healing selectors with vision models
3. **Form Filling Complexity** - Intelligent field mapping and context-aware filling
4. **CAPTCHA Handling** - AI-powered solving services
5. **Session Management** - Smart authentication and 2FA handling
6. **Error Recovery** - Self-healing automation

### Basic Playwright Configuration

```go
// Stealth configuration
browser, err := playwright.Run()
page, err := browser.NewPage()

// Anti-detection measures
page.AddInitScript(`
    Object.defineProperty(navigator, 'webdriver', {get: () => undefined});
    window.chrome = {runtime: {}};
    Object.defineProperty(navigator, 'plugins', {get: () => [1, 2, 3]});
`)

// Human-like delays
page.SetExtraHTTPHeaders(map[string]string{
    "User-Agent": "Mozilla/5.0...",
})
```

## AI Cover Letter Generation

### Flow

1. **Job Details Extraction**: Worker extracts job description, company info, requirements
2. **User Profile Aggregation**: Service fetches from database:
   - Projects (from `projects` domain)
   - Posts (from `posts` domain)
   - Technical writings (from `technicalwritings` domain)
   - Skills, interests, certifications
3. **AI Service Call**: Sends to existing `ai-service`:
   ```json
   {
     "job_description": "...",
     "company_name": "...",
     "user_profile": {
       "projects": [...],
       "posts": [...],
       "technical_writings": [...],
       "skills": [...]
     }
   }
   ```
4. **Cover Letter Generation**: AI service returns personalized cover letter
5. **Application**: Worker uses generated text to fill application form

### Integration with Existing AI Service

- Use existing `ai-service` container (already in docker-compose.yml)
- Add new endpoint: `POST /api/cover-letters/generate`
- Or extend existing chat/completion endpoints

## API Endpoints (Main Server)

### Job Applications

- `POST /api/job-applications` - Create new job application (enqueues to worker)
- `GET /api/job-applications` - List all applications
- `GET /api/job-applications/:id` - Get specific application
- `PATCH /api/job-applications/:id/status` - Update status
- `DELETE /api/job-applications/:id` - Cancel/delete application

### Job Websites

- `GET /api/job-websites` - List all websites (with current counts)
- `PATCH /api/job-websites/:id` - Update website config (limits, enabled)
- `POST /api/job-websites/:id/reset-counter` - Manually reset daily counter

## Environment Variables

### Worker Container

```bash
# Required
DATABASE_URL=postgres://...
REDIS_URL=redis://...
AI_SERVICE_URL=http://ai-service:8000

# Optional
JOB_APPLICATION_ENABLED=true
PLAYWRIGHT_HEADLESS=true
PLAYWRIGHT_BROWSER_PATH=/usr/bin/chromium
PLAYWRIGHT_SLOW_MO=100  # Delay between actions (ms)
PLAYWRIGHT_TIMEOUT=30000  # Page timeout (ms)

# Per-website credentials (encrypted in DB or env)
LINKEDIN_USERNAME=...
LINKEDIN_PASSWORD=...
GLASSDOOR_USERNAME=...
GLASSDOOR_PASSWORD=...
```

## Deployment

### Docker Compose

Add to `docker-compose.yml`:

```yaml
job-application-worker:
  build:
    context: .
    dockerfile: Dockerfile.job-application-worker
  container_name: woragis-job-application-worker
  depends_on:
    database:
      condition: service_healthy
    redis:
      condition: service_healthy
    ai-service:
      condition: service_started
  environment:
    APP_ENV: ${APP_ENV:-development}
    DATABASE_URL: postgres://postgres:postgres@database:5432/woragis?sslmode=disable
    REDIS_URL: redis://redis:6379/0
    AI_SERVICE_URL: ${AI_SERVICE_URL:-http://ai-service:8000}
    JOB_APPLICATION_ENABLED: ${JOB_APPLICATION_ENABLED:-true}
    PLAYWRIGHT_HEADLESS: ${PLAYWRIGHT_HEADLESS:-true}
  command: ["job-application-worker"]
  restart: unless-stopped
```

## Future Enhancements

This section documents planned improvements to prevent knowledge loss over time.

### Phase 2: Automated Job Discovery

**Status**: Not Started  
**Priority**: Medium  
**Estimated Effort**: 2-3 weeks

- **Scraping Engine**: Extend worker to scrape job listings from websites
- **Job Filtering**: AI-powered job matching based on user preferences (location, salary, role, tech stack)
- **Duplicate Detection**: Prevent applying to same job twice (hash job URL + company name)
- **Job Quality Scoring**: AI rates job postings (salary, company reputation, fit score)
- **Job Queue Management**: Prioritize high-quality jobs

**Implementation Notes**:
- Start with one website (e.g., WeWorkRemotely - simpler structure)
- Use Playwright to scrape job listings
- Store discovered jobs in separate table before creating JobApplication records
- Add filtering criteria to JobWebsite entity

### Phase 3: Advanced AI Features

**Status**: Not Started  
**Priority**: High  
**Estimated Effort**: 3-4 weeks

- **Resume Tailoring**: AI adapts resume sections for each application
  - Reorder experience based on relevance
  - Emphasize matching skills
  - Generate ATS-friendly versions
- **Interview Prep**: AI generates interview questions based on job description
  - Technical questions based on required skills
  - Behavioral questions based on job requirements
- **Follow-up Automation**: AI sends follow-up messages after application
  - LinkedIn messages to recruiters
  - Email follow-ups (if email available)
  - Personalized based on application date
- **Response Analysis**: AI analyzes rejection emails for improvement
  - Extract feedback (if any)
  - Identify patterns in rejections
  - Suggest profile improvements

**Implementation Notes**:
- Extend AI service with new endpoints
- Store interview prep in database
- Schedule follow-ups via scheduler domain
- Create analysis dashboard

### Phase 4: Multi-Account Support

**Status**: Not Started  
**Priority**: Low  
**Estimated Effort**: 2 weeks

- **Account Rotation**: Support multiple LinkedIn/Glassdoor accounts
  - Rotate accounts to avoid rate limits
  - Balance load across accounts
- **Profile Optimization**: AI suggests profile improvements
  - Analyze profile completeness
  - Suggest keywords based on target roles
  - Optimize headline and summary
- **Network Building**: Automated connection requests with personalized messages
  - Find relevant connections (recruiters, engineers at target companies)
  - Generate personalized connection messages
  - Respect connection limits per account

**Implementation Notes**:
- Add Account entity (linked to User)
- Store credentials securely (encrypted)
- Implement account rotation logic in orchestrator
- Add connection tracking to prevent spam

### Phase 5: Analytics & Insights

**Status**: Not Started  
**Priority**: Medium  
**Estimated Effort**: 2-3 weeks

- **Application Analytics**: Success rates per website, company, role
  - Dashboard showing application statistics
  - Success rate by company size, industry
  - Response time analysis
- **Cover Letter Performance**: A/B testing different cover letter styles
  - Test different tones (formal vs. casual)
  - Test different structures
  - Track which styles get responses
- **Timing Optimization**: Best times to apply (ML-based)
  - Analyze when applications get responses
  - Optimize application timing
  - Consider timezone, day of week
- **Market Trends**: AI analyzes job market trends
  - Popular skills, technologies
  - Salary trends
  - Demand by location

**Implementation Notes**:
- Create analytics domain
- Store metrics in time-series database (or PostgreSQL with partitioning)
- Build dashboard endpoints
- Use existing monitoring infrastructure

### Phase 6: AI-Enhanced Playwright (Advanced)

**Status**: Not Started  
**Priority**: High (after basic implementation works)  
**Estimated Effort**: 4-6 weeks

See [PLAYWRIGHT_AI_ENHANCEMENTS.md](./PLAYWRIGHT_AI_ENHANCEMENTS.md) for detailed plan.

**Key Features**:
- Self-healing selectors using vision models
- Behavioral pattern learning
- Intelligent error recovery
- Adaptive anti-detection strategies

**Implementation Phases**:
1. Basic vision model integration for element finding
2. ML-based timing generation
3. Full self-healing automation
4. Reinforcement learning for strategy optimization

### Phase 7: Integration with Other Services

**Status**: Not Started  
**Priority**: Low  
**Estimated Effort**: 1-2 weeks

- **Calendar Integration**: Schedule interviews automatically
- **Email Integration**: Parse job-related emails
- **CRM Integration**: Track applications in external CRM
- **Slack/Discord Notifications**: Get notified of application status changes

### Phase 8: Mobile App Support

**Status**: Not Started  
**Priority**: Very Low  
**Estimated Effort**: 4+ weeks

- Some job sites have mobile apps with different flows
- Use Appium or similar for mobile automation
- Separate mobile scrapers/automation

## Upgrade Tracking

When implementing upgrades, update this section with:
- **Date**: When upgrade was completed
- **Version**: Version number or tag
- **Changes**: What was added/modified
- **Breaking Changes**: Any API or schema changes
- **Migration Notes**: How to upgrade from previous version

### Completed Upgrades

_None yet - initial implementation pending_

### In Progress

_None yet_

### Planned

See "Future Enhancements" section above.

## Security Considerations

- **Credential Storage**: Store website credentials encrypted in database
- **Session Management**: Secure cookie/session storage
- **Rate Limiting**: Respect website ToS and rate limits
- **Data Privacy**: Ensure user data is handled securely
- **CAPTCHA Handling**: May require manual intervention or paid services

## Monitoring & Logging

- **Application Success Rate**: Track successful vs failed applications
- **Rate Limit Usage**: Monitor daily limits per website
- **Error Tracking**: Log all failures with context
- **Performance Metrics**: Application processing time
- **AI Service Latency**: Cover letter generation time

## Troubleshooting

### Common Issues

1. **Worker Not Processing Jobs**
   - Check Redis connection
   - Verify queue name matches
   - Check worker logs

2. **Rate Limit Exceeded**
   - Check daily counters in database
   - Verify reset logic
   - Manually reset if needed

3. **Playwright Failures**
   - Check browser installation
   - Verify headless mode settings
   - Review screenshots/videos for debugging

4. **AI Service Errors**
   - Verify AI service is running
   - Check API key configuration
   - Review request/response logs

## Notes

- This worker is designed to be **always running** (like translation-worker)
- Jobs are **enqueued via API** from the main server
- **Job discovery is manual** initially (can be automated later)
- **Rate limiting** prevents exceeding website limits
- **AI integration** provides personalized cover letters
- **Playwright** handles browser automation with AI enhancements for robustness

