# Job Application Worker Component

## Overview

A Node.js-based worker service that automates job application processes. It uses Playwright for browser automation, web scraping (LinkedIn, job sites), and AI for cover letter generation.

## Architecture

- **Language**: Node.js 18+
- **Port**: 8080 (health check only)
- **Message Queue**: RabbitMQ
- **Database**: PostgreSQL
- **Browser Automation**: Playwright
- **AI Service**: Calls AI service for cover letter generation

## Responsibilities

1. **Consume Job Application Jobs**: Listen to `job-applications.queue` for application jobs
2. **Web Scraping**: Scrape job postings from websites (LinkedIn, Glassdoor, etc.)
3. **AI-Powered Self-Healing Scrapers**: Automatically finds selectors when they break
4. **Cover Letter Generation**: Generate personalized cover letters using AI
5. **Application Automation**: Automate job application process
6. **Update Database**: Record application status and results

## Health Check

**Endpoint**: `GET /healthz`

**Checks**:
- RabbitMQ connection - CRITICAL

**Response**:
```json
{
  "status": "healthy|unhealthy",
  "checks": [
    {"name": "rabbitmq", "status": "ok"}
  ]
}
```

## Metrics

**Endpoint**: `GET /metrics`

Exposes Prometheus metrics:
- Job processing rate (success/failed)
- Job processing duration
- Queue depth
- DLQ size
- Scraping success rate
- AI service call metrics

## Configuration

### Environment Variables

#### Required
- `DATABASE_URL` - PostgreSQL connection string
- `RABBITMQ_URL` - RabbitMQ connection URL
- `AI_SERVICE_URL` - AI service URL (default: `http://ai-service:8000`)

#### Optional
- `JOB_APPLICATION_QUEUE_NAME` - Queue name (default: `job-applications.queue`)
- `JOB_APPLICATION_EXCHANGE` - Exchange name (default: `woragis.tasks`)
- `JOB_APPLICATION_ROUTING_KEY` - Routing key (default: `job-applications.process`)
- `PLAYWRIGHT_HEADLESS` - Run browser in headless mode (default: `true`)
- `PLAYWRIGHT_SLOW_MO` - Slow down operations (ms) (default: `100`)
- `PLAYWRIGHT_TIMEOUT` - Operation timeout (ms) (default: `30000`)
- `ENV` - Environment (development/production)

## Message Format

**Queue**: `job-applications.queue`

**Message**:
```json
{
  "id": "job-uuid",
  "user_id": "user-uuid",
  "job_url": "https://linkedin.com/jobs/view/123456",
  "website": "linkedin",
  "application_data": {
    "resume_id": "resume-uuid",
    "cover_letter_template": "template-uuid"
  }
}
```

## Processing Flow

1. **Consume Message**: Worker consumes message from RabbitMQ
2. **Navigate to Job**: Open job URL in browser
3. **Scrape Job Details**: Extract job title, description, requirements
4. **Generate Cover Letter**: Use AI service to generate personalized cover letter
5. **Fill Application Form**: Automatically fill application form
6. **Submit Application**: Submit application
7. **Update Database**: Record application status
8. **Acknowledge**: Acknowledge message on success

## AI-Powered Self-Healing Scrapers

The worker uses AI to automatically find HTML selectors when cached ones fail:

### How It Works

1. **Selector Cache**: Stores selectors in Redis with 7-day TTL
2. **AI Recovery**: When selectors fail, AI analyzes the page (HTML or screenshot) to find new ones
3. **Multiple Strategies**: Tries CSS selectors, XPath, text matching, etc.
4. **Auto-Update**: New selectors are automatically cached for future use

### Selector Cache Structure

```json
{
  "primary": "button[data-testid='easy-apply']",
  "alternatives": [
    "button:has-text('Easy Apply')",
    ".jobs-apply-button"
  ],
  "xpath": "//button[contains(text(), 'Apply')]",
  "text": "Easy Apply",
  "explanation": "Button with data-testid attribute",
  "cachedAt": "2024-01-15T10:30:00Z",
  "website": "linkedin",
  "action": "easy-apply-button"
}
```

### Recovery Process

```javascript
// When a selector fails:
1. Check Redis cache for cached selectors
2. If cached selectors fail → Invalidate cache
3. Use AI to analyze page HTML/screenshot
4. AI returns new selectors (primary + alternatives)
5. Try new selectors
6. If successful → Cache new selectors in Redis
7. Continue with application
```

## Supported Websites

- **LinkedIn** - Job applications via LinkedIn Easy Apply
- **Glassdoor** - Job applications (basic support)
- **Generic** - Generic job sites (with AI selector finding)

## Rate Limiting

The worker implements rate limiting to avoid:
- Being blocked by websites
- Triggering anti-bot measures
- Overwhelming target websites

**Orchestrator**: Manages rate limiting across multiple workers

## Error Handling

### Transient Errors
- Website temporarily unavailable
- Network errors
- Selector failures (auto-recovered via AI)

**Action**: Retry with exponential backoff (up to 3 attempts)

### Permanent Errors
- Invalid job URL
- Website structure changed (AI recovery failed)
- Authentication failures

**Action**: Route to Dead Letter Queue (DLQ)

## Dead Letter Queue

**DLQ Exchange**: `woragis.dlx`
**DLQ Routing Key**: `job-applications.queue.failed`

Failed messages are automatically routed to DLQ after max retries.

## Logging

**Format**: Structured JSON (production), Text (development)

**Service Name**: `job-application-worker`

**Key Log Fields**:
- `job_id` - Job identifier
- `user_id` - User ID
- `job_url` - Job URL
- `website` - Website name
- `status` - Success/failed
- `error` - Error message (if failed)
- `selector_cache_hit` - Whether selector cache was used

## Deployment

### Local Development

```bash
cd backend/job-application-worker
npm install
npm run dev
```

### Docker

```bash
docker build -f Dockerfile.job-application-worker -t woragis/job-application-worker .
docker run --env-file .env --shm-size=2gb woragis/job-application-worker
```

**Note**: Playwright requires shared memory (`--shm-size=2gb`)

### Kubernetes

Deploy as a Deployment:
- Health check probe on `/healthz`
- RabbitMQ connection required
- Database connection required
- AI service connection required
- Redis connection required (for selector cache)
- Shared memory requirement for Playwright

## Scaling

### Horizontal Scaling
- Multiple worker replicas consume from same queue
- RabbitMQ distributes messages (round-robin or fair dispatch)
- Each replica has its own browser instances
- **Note**: Browser instances are memory-intensive, limit replicas based on memory

### Resource Requirements
- **CPU**: 500m-1000m (0.5-1 core) - Browser automation is CPU-intensive
- **Memory**: 1Gi-2Gi (browser instances are memory-intensive)
- **Database Connections**: 5-10 per replica
- **Shared Memory**: 2GB (for Playwright)

## AI Service Integration

The worker calls the AI service for:
- Cover letter generation
- Selector finding (when selectors fail)
- Job description analysis

**AI Service Endpoints Used**:
- `POST /v1/chat` - Chat completion for cover letter generation
- `POST /v1/chat` - Selector finding assistance

## Monitoring

### Key Metrics
- Job application rate (applications/second)
- Success rate
- Failure rate
- Queue depth
- DLQ size
- Scraping success rate
- Selector cache hit rate
- AI service call duration

### Alerts
- DLQ size > 100 messages
- Failure rate > 10% (higher threshold due to website changes)
- Queue depth > 1000 messages
- Selector cache hit rate < 50% (may indicate website changes)

## Troubleshooting

### Common Issues

#### Applications Not Processing
- Check Playwright installation
- Check browser dependencies
- Check database connection
- Verify message format

#### High Failure Rate
- Check website structure changes (selectors may need updating)
- Verify AI service is available (for selector recovery)
- Check rate limiting (may be blocked)
- Verify job URLs are valid

#### Selector Failures
- Check selector cache in Redis
- Verify AI service is available (for recovery)
- Check website structure changes
- Review selector cache hit rate

#### Queue Backlog
- Check browser instances (may be memory-limited)
- Scale up worker replicas (if memory allows)
- Check worker health
- Verify workers are consuming messages

#### Browser Crashes
- Increase memory limits
- Check shared memory (`--shm-size`)
- Reduce concurrent browser instances
- Check for memory leaks

## Related Documentation

- [Architecture Decision Records](../adr/) - Worker architecture decisions
- [Monitoring DLQ](../runbooks/monitoring-dlq.md) - DLQ monitoring procedures
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
