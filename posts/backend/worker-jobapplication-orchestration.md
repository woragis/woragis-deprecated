# Job Application Worker Orchestration Pattern

## Overview
How the job application worker orchestrates automated job applications including scraping, cover letter generation, and form submission.

## Key Points

### Worker Components
- Queue: Redis queue for job management
- Database: PostgreSQL connection for data storage/retrieval
- Orchestrator: Rate limiting and website management
- Scraper: Playwright-based web scraping
- CoverLetterService: AI-powered cover letter generation

### Initialization Flow
1. Connect to Redis queue
2. Connect to PostgreSQL database
3. Initialize Playwright scraper (browser automation)
4. Initialize orchestrator with database connection
5. Initialize cover letter service

### Processing Loop
- Continuous polling loop
- Dequeue with 5-second timeout
- Rate limit checking before processing
- Job re-enqueuing when rate limits reached
- Graceful error handling (continue after failures)

### Job Processing Flow
1. Dequeue job from Redis
2. Check if website rate limit reached (via Orchestrator)
3. If limit reached: re-enqueue job, wait 1 hour, continue
4. Find or create application record in database
5. Fetch comprehensive user profile data
6. Generate cover letter using AI service
7. Apply to job using Playwright scraper
8. Update application status (applied/failed)
9. Increment website count in orchestrator
10. Mark job as complete in queue

### Rate Limiting Orchestration
- Per-website daily limits (e.g., LinkedIn: 10/day)
- Automatic count reset at midnight
- Website enable/disable support
- Count tracking in database
- Prevents over-applications per website

### Scraping Orchestration
- Playwright browser automation
- Dynamic selector finding (cached selectors with fallback)
- Form filling automation
- Screenshot capture for debugging
- Error recovery (selector refresh, retry logic)

### Cover Letter Generation
- Fetches user profile: experiences, projects, skills, certifications
- Calls AI service with profile + job info
- Generates personalized cover letter
- Stores cover letter with application

### Integration Points
- **Database**: User profiles, job applications, websites, selectors cache
- **AI Service**: HTTP REST API for cover letter generation
- **Redis**: Job queue
- **Playwright**: Browser automation for scraping
- **Websites**: LinkedIn, Indeed, etc.

## Implementation Details

### Error Handling
- Scraping errors: caught and logged, application marked as failed
- AI service errors: handled gracefully
- Database errors: logged and propagated
- Network errors: retry logic where appropriate

### Resource Management
- Proper connection cleanup on shutdown
- Playwright browser cleanup
- Database connection management
- Queue disconnection

## Potential Improvements
- Add worker pool for parallel processing (multiple browsers)
- Implement job prioritization (urgent applications first)
- Add circuit breaker for AI service failures
- Implement health check endpoint
- Add job cancellation support
- Support job scheduling (apply at specific times)
- Implement job result caching (successful application patterns)
- Add metrics collection (throughput, latency, success rate per website)
- Support job dependency management (wait for resume generation)
- Implement A/B testing for cover letters
- Add browser pool management (reuse browsers)
- Implement proxy rotation for scraping
- Add CAPTCHA solving integration
- Support multiple cover letter templates
- Implement application tracking (follow up reminders)
- Add screenshot comparison for debugging
- Support video recording of application process

