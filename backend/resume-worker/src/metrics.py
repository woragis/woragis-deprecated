"""
Prometheus metrics for Resume Worker.
"""
from prometheus_client import Counter, Histogram, Gauge, start_http_server
import threading

# Worker metrics
worker_jobs_processed_total = Counter(
    'worker_jobs_processed_total',
    'Total number of jobs processed',
    ['worker', 'status']  # status: success, failed
)

worker_job_duration_seconds = Histogram(
    'worker_job_duration_seconds',
    'Job processing duration in seconds',
    ['worker'],
    buckets=[.01, .05, .1, .25, .5, 1, 2.5, 5, 10, 30]
)

worker_jobs_failed_total = Counter(
    'worker_jobs_failed_total',
    'Total number of failed jobs',
    ['worker', 'error_type']
)

worker_jobs_retried_total = Counter(
    'worker_jobs_retried_total',
    'Total number of retried jobs',
    ['worker']
)

# Queue metrics
queue_depth = Gauge(
    'queue_depth',
    'Current queue depth',
    ['queue_name']
)

queue_dlq_size = Gauge(
    'queue_dlq_size',
    'Dead letter queue size',
    ['queue_name']
)

# Metrics server started flag
_metrics_server_started = False
_metrics_server_lock = threading.Lock()


def start_metrics_server(port=9090):
    """Start Prometheus metrics server on a separate port."""
    global _metrics_server_started
    
    with _metrics_server_lock:
        if _metrics_server_started:
            return
        
        # Start HTTP server for metrics (on different port than health)
        start_http_server(port)
        _metrics_server_started = True


def record_job_processed(worker, status, duration):
    """Record a processed job metric."""
    worker_jobs_processed_total.labels(worker=worker, status=status).inc()
    worker_job_duration_seconds.labels(worker=worker).observe(duration)


def record_job_failed(worker, error_type):
    """Record a failed job metric."""
    worker_jobs_failed_total.labels(worker=worker, error_type=error_type).inc()


def record_job_retried(worker):
    """Record a retried job metric."""
    worker_jobs_retried_total.labels(worker=worker).inc()


def set_queue_depth(queue_name, depth):
    """Set the queue depth gauge."""
    queue_depth.labels(queue_name=queue_name).set(depth)


def set_queue_dlq_size(queue_name, size):
    """Set the DLQ size gauge."""
    queue_dlq_size.labels(queue_name=queue_name).set(size)
