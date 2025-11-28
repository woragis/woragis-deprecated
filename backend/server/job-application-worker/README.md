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
- ⏳ Website-specific scrapers (LinkedIn, Glassdoor, etc.) - TODO

