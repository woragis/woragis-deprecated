# Job Applications Domain - Application Workflow

## Overview
Complete workflow for job application management from creation to application.

## Key Points

### Application Creation Flow
1. User requests application creation (company, job title, URL, website)
2. Service creates application record with "pending" status
3. Job enqueued to Redis queue for processing
4. Response returned to user with application ID

### Application Processing Flow (Worker)
1. Worker dequeues job
2. Checks website rate limit (via Orchestrator)
3. Creates/finds application record
4. Fetches comprehensive user profile data
5. Generates cover letter via AI service
6. Applies to job using Playwright scraper
7. Updates application status (applied/failed)
8. Stores cover letter and application details

### Status Transitions
- **pending**: Initial state, queued for processing
- **processing**: Worker is processing application
- **applied**: Successfully applied to job
- **failed**: Application failed (scraping error, etc.)
- Additional statuses: interview_scheduled, offer_received, rejected, etc.

### Rate Limiting
- Per-website daily limits (e.g., LinkedIn: 10/day)
- Orchestrator manages rate limits
- Jobs re-enqueued if rate limit reached
- Automatic count reset at midnight

### Components Integration
- **Queue**: Redis queue for job management
- **Database**: Application records, websites, user profiles
- **Orchestrator**: Rate limiting and website management
- **Scraper**: Playwright-based web scraping
- **Cover Letter Service**: AI-powered cover letter generation

## Potential Improvements
- Add application status webhooks
- Implement application tracking (follow-up reminders)
- Add application analytics (success rate per website)
- Support bulk application creation
- Add application scheduling (apply at specific times)
- Implement application templates
- Add application notes and tags
- Support application status updates from websites
- Add application search and filtering
- Implement application export functionality

