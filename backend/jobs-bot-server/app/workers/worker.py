from app.utils.errors import JobsBotError
from app.utils.logging import log_event
from app.api.websocket import emit_ws_event
from app.workers.states import RunState, JobState
from app.orchestrator.queue import JobQueue


class Worker:
    def __init__(self, run_id, job_queue: JobQueue):
        try:
            self.run_id = run_id
            self.job_queue = job_queue
            self.state = RunState.PENDING
            self.jobs = []
            log_event("worker_initialized", run_id=run_id)
            # Optionally emit_ws_event("worker_initialized", run_id=run_id)
        except Exception as e:
            log_event("worker_init_failed", run_id=run_id,
                      error_code="JOBS001", details=str(e))
            raise JobsBotError("JOBS001", {"error": str(e)})

    def start(self):
        try:
            self.state = RunState.RUNNING
            log_event("worker_started", run_id=self.run_id)
            emit_ws_event("worker_started",
                          run_id=self.run_id, state=self.state)
            while not self.job_queue.is_empty():
                job = self.job_queue.get_job()
                self.jobs.append({"job": job, "state": JobState.QUEUED})
                log_event("job_queued", run_id=self.run_id,
                          job_id=job.get("id"))
                emit_ws_event("job_queued", run_id=self.run_id,
                              job_id=job.get("id"), state=JobState.QUEUED)
        except Exception as e:
            log_event("worker_start_failed", run_id=self.run_id,
                      error_code="JOBS002", details=str(e))
            emit_ws_event("worker_start_failed", run_id=self.run_id,
                          error_code="JOBS002", details=str(e))
            raise JobsBotError("JOBS002", {"error": str(e)})

    def stop(self):
        try:
            self.state = RunState.CANCELLED
            log_event("worker_stopped", run_id=self.run_id)
            emit_ws_event("worker_stopped",
                          run_id=self.run_id, state=self.state)
        except Exception as e:
            log_event("worker_stop_failed", run_id=self.run_id,
                      error_code="JOBS001", details=str(e))
            emit_ws_event("worker_stop_failed", run_id=self.run_id,
                          error_code="JOBS001", details=str(e))
            raise JobsBotError("JOBS001", {"error": str(e)})
