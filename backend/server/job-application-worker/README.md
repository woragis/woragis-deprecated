# Job Application Worker (Node.js)

This is a Node.js worker that processes job applications using Playwright for browser automation.

## Architecture

- **Queue**: Redis-based job queue (same as Go backend)
- **Database**: PostgreSQL (same database as Go backend)
- **Playwright**: Native Node.js browser automation
- **AI Service**: Calls Go backend's AI service for cover letter generation

## Setup

```bash
cd job-application-worker
npm install
```

## Environment Variables

```bash
DATABASE_URL=postgres://user:pass@host:5432/dbname
REDIS_URL=redis://localhost:6379/0
AI_SERVICE_URL=http://ai-service:8000
PLAYWRIGHT_HEADLESS=true
PLAYWRIGHT_SLOW_MO=100
PLAYWRIGHT_TIMEOUT=30000
```

## Development

```bash
npm run dev
```

## Production

The worker runs in a Docker container. See `Dockerfile.job-application-worker` in the parent directory.

## Communication with Go Backend

- **Redis Queue**: Go backend enqueues jobs, Node.js worker dequeues them
- **Database**: Both read/write to the same PostgreSQL database
- **AI Service**: Node.js worker calls Go backend's AI service via HTTP

## Implementation Status

- ✅ Queue system
- ✅ Database integration
- ✅ Rate limiting orchestrator
- ✅ Cover letter generation
- ✅ **AI-Powered Self-Healing Scrapers** - Automatically finds selectors when they break
- ✅ Selector caching in Redis (7-day TTL)
- ⏳ Website-specific scrapers (LinkedIn, Glassdoor, etc.) - Basic structure ready, needs refinement

## AI Self-Healing System

The scraper uses AI to automatically find HTML selectors when cached ones fail:

1. **Selector Cache**: Stores selectors in Redis with 7-day TTL
2. **AI Recovery**: When selectors fail, AI analyzes the page (HTML or screenshot) to find new ones
3. **Multiple Strategies**: Tries CSS selectors, XPath, text matching, etc.
4. **Auto-Update**: New selectors are automatically cached for future use

### How It Works

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

