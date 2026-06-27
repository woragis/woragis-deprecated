# Testing & Validation

## Running Tests

- Use `pytest` for unit and integration tests.
- Example: `pytest tests/`

## Coverage

- Input validation (API, LLM, job data, config)
- Error handling and propagation
- Observability: logs and WebSocket events

## Simulating Failures

- Provide invalid input to API and LLM
- Simulate platform/browser/LLM errors
- Verify logs, events, and DB records
