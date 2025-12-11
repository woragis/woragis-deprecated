"""
Redis Queue for Resume Generation Jobs

Handles job enqueueing, dequeueing, and status tracking in Redis.
"""

import os
import json
import logging
import redis
from typing import Optional, Dict, Any
from datetime import datetime, timedelta

logger = logging.getLogger(__name__)

QUEUE_KEY = "resumes:queue"
JOB_PREFIX = "resumes:job:"
DEAD_LETTER_QUEUE_KEY = "resumes:dead-letter:queue"
EVENTS_CHANNEL = "resumes:events"


class Queue:
    """Redis-based queue for resume generation jobs."""

    def __init__(self):
        self.client = None
        self.pubsub = None

    def connect(self):
        """Connect to Redis."""
        redis_url = os.getenv("REDIS_URL", "redis://localhost:6379/0")
        self.client = redis.from_url(redis_url, decode_responses=True)

        self.client.ping()
        logger.info("Connected to Redis", extra={"redis_url": redis_url})

    def disconnect(self):
        """Disconnect from Redis."""
        if self.pubsub:
            self.pubsub.close()
        if self.client:
            self.client.close()
            logger.info("Disconnected from Redis")

    def enqueue_job(self, job: Dict[str, Any]) -> str:
        """
        Enqueue a resume generation job.

        Args:
            job: Job data dictionary

        Returns:
            Job ID
        """
        if not job.get("id"):
            job["id"] = self._generate_id()

        job["status"] = "pending"
        job["created_at"] = datetime.utcnow().isoformat()
        job["updated_at"] = datetime.utcnow().isoformat()

        job_key = JOB_PREFIX + job["id"]
        job_data = json.dumps(job, default=str)

        # Store job data with 7 day TTL
        self.client.setex(job_key, 7 * 24 * 3600, job_data)

        # Add to queue
        self.client.lpush(QUEUE_KEY, job["id"])

        logger.info(
            "Job enqueued",
            extra={
                "job_id": job["id"],
                "user_id": job.get("user_id"),
                "job_application_id": job.get("job_application_id"),
            },
        )

        return job["id"]

    def dequeue_job(self, timeout: int = 5) -> Optional[Dict[str, Any]]:
        """
        Dequeue a job from the queue (blocking).

        Args:
            timeout: Timeout in seconds

        Returns:
            Job data or None if timeout
        """
        # Blocking pop from queue
        result = self.client.brpop(QUEUE_KEY, timeout)

        if not result:
            return None

        job_id = result[1]
        return self.get_job(job_id)

    def get_job(self, job_id: str) -> Optional[Dict[str, Any]]:
        """
        Get job data by ID.

        Args:
            job_id: Job ID

        Returns:
            Job data or None if not found
        """
        job_key = JOB_PREFIX + job_id
        job_data = self.client.get(job_key)

        if not job_data:
            return None

        return json.loads(job_data)

    def update_job_status(
        self,
        job_id: str,
        status: str,
        error: Optional[str] = None,
        error_type: Optional[str] = None,
        retry_count: Optional[int] = None,
        result: Optional[Dict[str, Any]] = None,
    ):
        """
        Update job status and metadata.

        Args:
            job_id: Job ID
            status: New status (pending, processing, completed, failed, retrying, dead_letter)
            error: Error message if failed
            error_type: Error type (transient, permanent)
            retry_count: Current retry count
            result: Result data if completed
        """
        job = self.get_job(job_id)
        if not job:
            logger.warning("Job not found for status update", extra={"job_id": job_id})
            return

        job["status"] = status
        job["updated_at"] = datetime.utcnow().isoformat()

        if error:
            job["last_error"] = error
            job["last_error_at"] = datetime.utcnow().isoformat()

        if error_type:
            job["last_error_type"] = error_type

        if retry_count is not None:
            job["retry_count"] = retry_count

        if result:
            job["result"] = result

        # Store error history
        if error and "errors_history" not in job:
            job["errors_history"] = []
        if error:
            job["errors_history"].append(
                {
                    "error": error,
                    "type": error_type,
                    "at": datetime.utcnow().isoformat(),
                }
            )
            # Keep only last 10 errors
            job["errors_history"] = job["errors_history"][-10:]

        job_key = JOB_PREFIX + job_id
        job_data = json.dumps(job, default=str)
        self.client.setex(job_key, 7 * 24 * 3600, job_data)

        # Publish event
        self._publish_event(job_id, status, error, result)

    def mark_job_complete(self, job_id: str, result: Dict[str, Any]):
        """Mark job as completed with result."""
        self.update_job_status(job_id, "completed", result=result)

    def mark_job_failed(
        self,
        job_id: str,
        error: str,
        error_type: str = "permanent",
        retry_count: int = 0,
    ):
        """Mark job as failed."""
        self.update_job_status(
            job_id, "failed", error=error, error_type=error_type, retry_count=retry_count
        )

    def mark_job_retrying(
        self, job_id: str, error: str, retry_count: int, backoff_seconds: int
    ):
        """Mark job as retrying."""
        self.update_job_status(
            job_id,
            "retrying",
            error=error,
            error_type="transient",
            retry_count=retry_count,
        )
        logger.warning(
            "Job will retry",
            extra={
                "job_id": job_id,
                "retry_count": retry_count,
                "backoff_seconds": backoff_seconds,
            },
        )

    def move_to_dead_letter(self, job_id: str, error: str):
        """Move job to dead letter queue."""
        job = self.get_job(job_id)
        if not job:
            return

        job["status"] = "dead_letter"
        job["updated_at"] = datetime.utcnow().isoformat()
        job["dead_letter_reason"] = error
        job["dead_lettered_at"] = datetime.utcnow().isoformat()

        job_key = JOB_PREFIX + job_id
        job_data = json.dumps(job, default=str)

        # Store in dead letter queue with 30 day TTL
        dead_letter_key = JOB_PREFIX + "dead-letter:" + job_id
        self.client.setex(dead_letter_key, 30 * 24 * 3600, job_data)

        # Add to dead letter queue list
        self.client.lpush(DEAD_LETTER_QUEUE_KEY, job_id)

        # Remove from main queue storage
        self.client.delete(job_key)

        self._publish_event(job_id, "dead_letter", error=error)

        logger.error(
            "Job moved to dead letter queue",
            extra={"job_id": job_id, "error": error},
        )

    def _publish_event(
        self,
        job_id: str,
        status: str,
        error: Optional[str] = None,
        result: Optional[Dict[str, Any]] = None,
    ):
        """Publish event to Redis channel."""
        event = {
            "job_id": job_id,
            "status": status,
            "timestamp": datetime.utcnow().isoformat(),
        }

        if error:
            event["error"] = error

        if result:
            event["result"] = result

        self.client.publish(EVENTS_CHANNEL, json.dumps(event))

    def _generate_id(self) -> str:
        """Generate unique job ID."""
        return f"{int(datetime.utcnow().timestamp() * 1000)}-{os.urandom(4).hex()}"

