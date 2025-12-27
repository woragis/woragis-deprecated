package translations

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
)

// Worker processes translation jobs from Redis queue.
type Worker struct {
	queue    translationsdomain.Queue
	service  translationsdomain.Service
	logger   *slog.Logger
	stopChan chan struct{}
}

// NewWorker creates a new translation worker.
func NewWorker(
	queue translationsdomain.Queue,
	service translationsdomain.Service,
	logger *slog.Logger,
) *Worker {
	return &Worker{
		queue:    queue,
		service:  service,
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

// Start begins processing translation jobs.
func (w *Worker) Start(ctx context.Context) {
	w.logger.Info("translation worker started")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("translation worker stopping (context cancelled)")
			return
		case <-w.stopChan:
			w.logger.Info("translation worker stopping (stop signal)")
			return
		default:
			w.processJob(ctx)
		}
	}
}

// Stop signals the worker to stop.
func (w *Worker) Stop() {
	close(w.stopChan)
}

func (w *Worker) processJob(ctx context.Context) {
	// Dequeue job with 5 second timeout
	job, err := w.queue.DequeueJob(ctx, 5*time.Second)
	if err != nil {
		if err.Error() != "translations: job not found" {
			w.logger.Error("failed to dequeue job", slog.Any("error", err))
		}
		return
	}

	if job == nil {
		// No job available, continue polling
		return
	}

	w.logger.Info("processing translation job",
		slog.String("jobId", job.ID),
		slog.String("entityType", string(job.EntityType)),
		slog.String("entityId", job.EntityID),
		slog.String("language", string(job.Language)),
	)

	// Validate job data
	if err := ValidateTranslationJob(job); err != nil {
		w.logger.Error("invalid translation job",
			slog.String("jobId", job.ID),
			slog.Any("error", err),
		)
		_ = w.queue.MarkJobFailed(ctx, job.ID, fmt.Sprintf("validation failed: %v", err))
		return
	}

	// Validate job data again before processing (defense in depth)
	if err := ValidateTranslationJob(job); err != nil {
		w.logger.Error("translation job validation failed before processing",
			slog.String("jobId", job.ID),
			slog.Any("error", err),
		)
		_ = w.queue.MarkJobFailed(ctx, job.ID, fmt.Sprintf("validation failed: %v", err))
		return
	}

	// Process the translation
	if err := w.service.ProcessTranslationJob(ctx, job); err != nil {
		w.logger.Error("failed to process translation job",
			slog.String("jobId", job.ID),
			slog.Any("error", err),
		)
		// Mark job as failed
		_ = w.queue.MarkJobFailed(ctx, job.ID, err.Error())
		return
	}

	// Mark job as complete
	if err := w.queue.MarkJobComplete(ctx, job.ID); err != nil {
		w.logger.Warn("failed to mark job as complete",
			slog.String("jobId", job.ID),
			slog.Any("error", err),
		)
	}

	w.logger.Info("translation job completed",
		slog.String("jobId", job.ID),
	)
}

