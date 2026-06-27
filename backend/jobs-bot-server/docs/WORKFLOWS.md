# Workflows & State Machines

## Run Lifecycle

- States: PENDING → RUNNING → COMPLETED / FAILED / CANCELLED
- Events: `run_started`, `run_completed`, `run_failed`, `run_cancelled`

## Job Lifecycle

- States: QUEUED → SKIPPED / APPLYING → APPLIED / FAILED
- Events: `job_queued`, `job_started`, `job_skipped`, `applied`, `failed`

## Platform Integration Flow

1. Login (verify session)
2. Collect jobs
3. For each job:
   - LLM fit check
   - LLM text generation
   - Apply (browser automation)
   - Emit events and log
