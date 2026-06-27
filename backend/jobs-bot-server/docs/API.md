# API Reference

## HTTP Endpoints

### POST /runs

Start a new job application run.

- Request: `{ "max_applications": 10, "platforms": ["greenhouse", "lever"] }`
- Response: `{ "run_id": "...", "state": "PENDING", ... }`
- Errors: See [Error Codes](OBSERVABILITY.md#error-codes)

### GET /runs/{run_id}

Get status of a run.

### POST /runs/{run_id}/stop

Stop a running job application.

### GET /health

Health check.

## WebSocket Events

- `job_queued`, `job_started`, `job_skipped`, `applied`, `failed`, `run_completed`, `run_failed`, `run_cancelled`, `error`, etc.
- All events include: `event`, `run_id`, `job_id`, `platform`, `state`, `error_code`, `details`

## Example Event

```json
{
  "event": "job_started",
  "run_id": "abc123",
  "job_id": "job456",
  "platform": "greenhouse",
  "state": "APPLYING",
  "error_code": null,
  "details": "Started applying to job."
}
```
