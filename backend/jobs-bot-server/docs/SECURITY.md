# Security & Hardening

## Sensitive Data

- Store secrets in environment variables
- Never log API keys or credentials

## Concurrency & Scheduling

- Only one run at a time (enforced in orchestrator)
- Use APScheduler for scheduled runs

## Recovery & Restart

- All state is persisted; system can be restarted safely
- See logs and DB for last known state
