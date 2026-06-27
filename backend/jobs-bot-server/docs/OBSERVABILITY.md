# Observability & Debugging

## Logging

- All major events, errors, and state changes are logged with context: `run_id`, `job_id`, `platform`, `state`, `error_code`, `details`, `timestamp`.
- Logs are persisted to both file and database (`logs/` and `logs` table).

## WebSocket Events

- Real-time events for all workflow steps and errors.
- See [API.md](API.md#websocket-events) for event types and fields.

## Database Records

- All run/job state transitions and errors are persisted in the `logs` table.
- Schema: `id`, `run_id`, `event`, `timestamp`, `details`

## Error Codes

| Code       | Meaning                    |
| ---------- | -------------------------- |
| API001     | General API error          |
| API002     | Run not found              |
| API003     | Run already in progress    |
| JOBS001    | Job queue error            |
| JOBS002    | Job application failed     |
| LLM001     | LLM fit check failed       |
| LLM002     | LLM response parsing error |
| BROWSER001 | Browser session crashed    |
| ...        | ...                        |

## Tracing & Debugging

- Use `run_id` and `job_id` to trace all related logs, events, and DB records.
- All errors include a unique code and context for fast debugging.
