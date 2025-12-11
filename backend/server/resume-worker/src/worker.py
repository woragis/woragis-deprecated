"""
Resume Generation Worker

Long-running worker that processes resume generation jobs from Redis queue.
"""

import os
import sys
import time
import logging
import logging.handlers
import json
from datetime import datetime
from typing import Optional, Dict, Any

# Add src directory to path for imports
sys.path.insert(0, os.path.dirname(__file__))

from resume_queue import Queue
from database import Database
from ai_service import AIService
from resume_generator import ResumeGenerator
from translation_helper import TranslationHelper

# Configure logging
def setup_logging():
    """Setup structured logging with rotation."""
    log_dir = os.getenv("LOG_DIR", "/app/logs")
    os.makedirs(log_dir, exist_ok=True)

    # Create JSON formatter
    class JSONFormatter(logging.Formatter):
        def format(self, record):
            log_data = {
                "timestamp": datetime.utcnow().strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z",
                "level": record.levelname,
                "service": "resume-worker",
                "message": record.getMessage(),
            }
            
            # Extract context from record attributes (extra dict)
            context = {}
            for key, value in record.__dict__.items():
                if key not in ["name", "msg", "args", "created", "filename", "funcName", 
                               "levelname", "levelno", "lineno", "module", "msecs", 
                               "message", "pathname", "process", "processName", "relativeCreated",
                               "thread", "threadName", "exc_info", "exc_text", "stack_info",
                               "taskName", "processName", "threadName"]:
                    try:
                        # Only include JSON-serializable values
                        json.dumps(value)
                        context[key] = value
                    except (TypeError, ValueError):
                        context[key] = str(value)
            
            if context:
                log_data["context"] = context
            
            # Add exception info if present
            if record.exc_info:
                log_data["exception"] = self.formatException(record.exc_info)
            
            return json.dumps(log_data, default=str)
    
    formatter = JSONFormatter()

    # Main log file (all levels)
    main_handler = logging.handlers.TimedRotatingFileHandler(
        filename=os.path.join(log_dir, "resume-worker.log"),
        when="midnight",
        interval=1,
        backupCount=30,
        encoding="utf-8",
    )
    main_handler.setLevel(logging.DEBUG)
    main_handler.setFormatter(formatter)

    # Error log file (errors only)
    error_handler = logging.handlers.TimedRotatingFileHandler(
        filename=os.path.join(log_dir, "resume-worker.error.log"),
        when="midnight",
        interval=1,
        backupCount=30,
        encoding="utf-8",
    )
    error_handler.setLevel(logging.ERROR)
    error_handler.setFormatter(formatter)

    # Console handler (for Docker logs)
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setLevel(logging.INFO)
    console_formatter = logging.Formatter(
        "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    )
    console_handler.setFormatter(console_formatter)

    # Configure root logger
    root_logger = logging.getLogger()
    root_logger.setLevel(logging.DEBUG)
    root_logger.addHandler(main_handler)
    root_logger.addHandler(error_handler)
    root_logger.addHandler(console_handler)

    return root_logger


logger = setup_logging()

# Configuration
DATABASE_URL = os.getenv(
    "DATABASE_URL", "postgres://postgres:postgres@database:5432/woragis?sslmode=disable"
)
AI_SERVICE_URL = os.getenv("AI_SERVICE_URL", "http://ai-service:8000")
OUTPUT_DIR = os.getenv("RESUME_OUTPUT_DIR", "/app/output")
RESULTS_LOG_DIR = os.getenv("RESULTS_LOG_DIR", "/app/results")
MAX_RETRIES = int(os.getenv("MAX_RETRIES", "3"))
QUEUE_TIMEOUT = int(os.getenv("QUEUE_TIMEOUT", "5"))


def classify_error(error: Exception) -> str:
    """
    Classify error as transient or permanent.

    Args:
        error: Exception to classify

    Returns:
        "transient" or "permanent"
    """
    error_type = type(error).__name__
    error_msg = str(error).lower()

    # Permanent errors
    if "authentication" in error_msg or "unauthorized" in error_msg:
        return "permanent"
    if "not found" in error_msg or "does not exist" in error_msg:
        return "permanent"
    if "permission denied" in error_msg:
        return "permanent"
    if "disk full" in error_msg or "no space" in error_msg:
        return "permanent"

    # Transient errors (default)
    return "transient"


def calculate_backoff(retry_count: int, max_backoff: int = 60) -> int:
    """
    Calculate exponential backoff in seconds.

    Args:
        retry_count: Current retry count
        max_backoff: Maximum backoff in seconds

    Returns:
        Backoff seconds
    """
    backoff = min(2 ** retry_count, max_backoff)
    return backoff


class Worker:
    """Resume generation worker."""

    def __init__(self):
        self.queue = Queue()
        self.db = None
        self.ai_service = None
        self.translation_helper = None
        self.generator = None
        self.running = False

    def start(self):
        """Start the worker."""
        logger.info("Starting resume generation worker")

        # Initialize connections
        self.queue.connect()
        self.db = Database(DATABASE_URL)
        self.db.connect()

        self.ai_service = AIService(AI_SERVICE_URL)
        self.translation_helper = TranslationHelper(DATABASE_URL)
        self.translation_helper.connect()

        self.generator = ResumeGenerator(
            self.db, self.ai_service, OUTPUT_DIR, self.translation_helper
        )

        self.running = True
        logger.info("Resume generation worker started")

        # Start processing loop
        self.process_loop()

    def stop(self):
        """Stop the worker."""
        logger.info("Stopping resume generation worker")
        self.running = False

        if self.translation_helper:
            self.translation_helper.close()
        if self.db:
            self.db.close()
        if self.queue:
            self.queue.disconnect()

    def process_loop(self):
        """Main processing loop."""
        while self.running:
            try:
                # Dequeue job with timeout
                job = self.queue.dequeue_job(QUEUE_TIMEOUT)

                if not job:
                    # No job available, continue polling
                    continue

                job_id = job.get("id")
                user_id = job.get("user_id")
                job_application_id = job.get("job_application_id")

                logger.info(
                    "Processing resume generation job",
                    extra={
                        "job_id": job_id,
                        "user_id": user_id,
                        "job_application_id": job_application_id,
                    },
                )

                # Update status to processing
                self.queue.update_job_status(job_id, "processing")

                # Process the job
                self.process_job(job)

            except KeyboardInterrupt:
                logger.info("Received interrupt signal")
                self.stop()
                break
            except Exception as e:
                logger.error(
                    "Error in process loop",
                    extra={"error": str(e), "error_type": type(e).__name__},
                    exc_info=True,
                )
                # Continue processing other jobs
                time.sleep(1)

    def process_job(self, job: Dict[str, Any]):
        """
        Process a resume generation job.

        Args:
            job: Job data dictionary
        """
        job_id = job.get("id")
        user_id = job.get("user_id")
        job_description = job.get("job_description", "")
        job_title = job.get("job_title", "Software Engineer")
        language = job.get("language", "en")
        retry_count = job.get("retry_count", 0)

        start_time = time.time()

        try:
            # Generate resume
            result = self.generator.generate_resume(
                user_id=user_id,
                job_description=job_description,
                job_title=job_title,
                output_filename=None,  # Let worker generate it
                language=language,
            )

            duration_ms = int((time.time() - start_time) * 1000)

            # Mark job as complete
            self.queue.mark_job_complete(
                job_id,
                {
                    "output_path": result.get("output_path"),
                    "file_name": result.get("file_name"),
                    "file_size": result.get("file_size"),
                    "tags": result.get("tags", []),
                    "duration_ms": duration_ms,
                },
            )

            logger.info(
                "Resume generation completed",
                extra={
                    "job_id": job_id,
                    "user_id": user_id,
                    "duration_ms": duration_ms,
                    "file_size": result.get("file_size"),
                    "output_path": result.get("output_path"),
                },
            )

        except Exception as e:
            error_type = classify_error(e)
            error_msg = str(e)
            duration_ms = int((time.time() - start_time) * 1000)

            logger.error(
                "Resume generation failed",
                extra={
                    "job_id": job_id,
                    "user_id": user_id,
                    "error": error_msg,
                    "error_type": error_type,
                    "retry_count": retry_count,
                    "duration_ms": duration_ms,
                },
                exc_info=True,
            )

            # Handle retry logic
            if error_type == "permanent" or retry_count >= MAX_RETRIES:
                # Permanent error or max retries reached
                self.queue.mark_job_failed(
                    job_id, error_msg, error_type=error_type, retry_count=retry_count
                )

                if retry_count >= MAX_RETRIES:
                    # Move to dead letter queue
                    self.queue.move_to_dead_letter(
                        job_id, f"Max retries ({MAX_RETRIES}) reached"
                    )
            else:
                # Transient error, will retry
                retry_count += 1
                backoff_seconds = calculate_backoff(retry_count)

                self.queue.mark_job_retrying(
                    job_id, error_msg, retry_count, backoff_seconds
                )

                # Re-enqueue with backoff
                job["retry_count"] = retry_count
                job["status"] = "pending"
                job["scheduled_at"] = (
                    datetime.utcnow().timestamp() + backoff_seconds
                )

                # Store updated job
                job_key = f"resumes:job:{job_id}"
                job_data = json.dumps(job, default=str)
                self.queue.client.setex(job_key, 7 * 24 * 3600, job_data)

                # Re-enqueue after backoff
                time.sleep(backoff_seconds)
                self.queue.client.lpush("resumes:queue", job_id)

                logger.warning(
                    "Job will retry",
                    extra={
                        "job_id": job_id,
                        "retry_count": retry_count,
                        "backoff_seconds": backoff_seconds,
                    },
                )


def main():
    """Main entry point."""
    worker = Worker()

    try:
        worker.start()
    except KeyboardInterrupt:
        logger.info("Received interrupt signal")
    except Exception as e:
        logger.error("Worker failed to start", extra={"error": str(e)}, exc_info=True)
        sys.exit(1)
    finally:
        worker.stop()


if __name__ == "__main__":
    main()

