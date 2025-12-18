package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// WorkerJobsProcessedTotal counts the total number of jobs processed
	WorkerJobsProcessedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_jobs_processed_total",
			Help: "Total number of jobs processed",
		},
		[]string{"worker", "status"}, // status: success, failed
	)

	// WorkerJobDuration tracks the duration of job processing in seconds
	WorkerJobDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "worker_job_duration_seconds",
			Help:    "Job processing duration in seconds",
			Buckets: []float64{.01, .05, .1, .25, .5, 1, 2.5, 5, 10, 30},
		},
		[]string{"worker"},
	)

	// WorkerJobsFailedTotal counts the total number of failed jobs
	WorkerJobsFailedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_jobs_failed_total",
			Help: "Total number of failed jobs",
		},
		[]string{"worker", "error_type"},
	)

	// WorkerJobsRetriedTotal counts the total number of retried jobs
	WorkerJobsRetriedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "worker_jobs_retried_total",
			Help: "Total number of retried jobs",
		},
		[]string{"worker"},
	)

	// QueueDepth tracks the current queue depth
	QueueDepth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "queue_depth",
			Help: "Current queue depth",
		},
		[]string{"queue_name"},
	)

	// QueueDLQSize tracks the dead letter queue size
	QueueDLQSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "queue_dlq_size",
			Help: "Dead letter queue size",
		},
		[]string{"queue_name"},
	)
)

// RecordJobProcessed records a processed job metric
func RecordJobProcessed(worker, status string, duration float64) {
	WorkerJobsProcessedTotal.WithLabelValues(worker, status).Inc()
	WorkerJobDuration.WithLabelValues(worker).Observe(duration)
}

// RecordJobFailed records a failed job metric
func RecordJobFailed(worker, errorType string) {
	WorkerJobsFailedTotal.WithLabelValues(worker, errorType).Inc()
}

// RecordJobRetried records a retried job metric
func RecordJobRetried(worker string) {
	WorkerJobsRetriedTotal.WithLabelValues(worker).Inc()
}

// SetQueueDepth sets the queue depth gauge
func SetQueueDepth(queueName string, depth float64) {
	QueueDepth.WithLabelValues(queueName).Set(depth)
}

// SetQueueDLQSize sets the DLQ size gauge
func SetQueueDLQSize(queueName string, size float64) {
	QueueDLQSize.WithLabelValues(queueName).Set(size)
}
