package schedulerworker

import (
	"context"
	"log/slog"
	"time"

	schedulerdomain "github.com/woragis/backend/server/app/internal/domains/scheduler"
)

// Runner periodically checks due schedules and executes them.
type Runner struct {
	service  *schedulerdomain.Service
	logger   *slog.Logger
	interval time.Duration
}

// NewRunner builds a new Runner.
func NewRunner(service *schedulerdomain.Service, logger *slog.Logger, interval time.Duration) *Runner {
	if interval <= 0 {
		interval = time.Minute
	}
	return &Runner{
		service:  service,
		logger:   logger,
		interval: interval,
	}
}

// Start begins processing schedules until the context is cancelled.
func (r *Runner) Start(ctx context.Context) {
	if r.service == nil {
		return
	}
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	r.process(ctx)

	for {
		select {
		case <-ticker.C:
			r.process(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (r *Runner) process(ctx context.Context) {
	now := time.Now().UTC()
	schedules, err := r.service.ListDue(ctx, now)
	if err != nil {
		if r.logger != nil {
			r.logger.Error("scheduler worker: list due failed", slog.Any("error", err))
		}
		return
	}

	for _, schedule := range schedules {
		s := schedule
		if err := r.service.Execute(ctx, &s); err != nil && r.logger != nil {
			r.logger.Error("scheduler worker: execute failed", slog.String("schedule_id", s.ID.String()), slog.Any("error", err))
		}
	}
}
