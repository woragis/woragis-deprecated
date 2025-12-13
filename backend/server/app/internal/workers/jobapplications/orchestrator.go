package jobapplications

import (
	"context"
	"log/slog"

	jobwebsitesdomain "github.com/woragis/backend/server/app/internal/domains/jobwebsites"
)

// Orchestrator manages rate limits and website rotation.
type Orchestrator struct {
	websiteService jobwebsitesdomain.Service
	logger         *slog.Logger
}

// NewOrchestrator creates a new orchestrator.
func NewOrchestrator(websiteService jobwebsitesdomain.Service, logger *slog.Logger) *Orchestrator {
	return &Orchestrator{
		websiteService: websiteService,
		logger:         logger,
	}
}

// GetAvailableWebsites returns websites that are enabled and haven't reached their daily limit.
func (o *Orchestrator) GetAvailableWebsites(ctx context.Context) ([]*jobwebsitesdomain.JobWebsite, error) {
	websites, err := o.websiteService.ListJobWebsites(ctx, true) // enabled only
	if err != nil {
		return nil, err
	}

	var available []*jobwebsitesdomain.JobWebsite
	for i := range websites {
		website := &websites[i]

		// Check if should reset (new day)
		if website.ShouldReset() {
			if err := o.websiteService.ResetCount(ctx, website.ID); err != nil {
				o.logger.Warn("failed to reset website count",
					slog.String("website", website.Name),
					slog.Any("error", err),
				)
				continue
			}
			// Reload website to get updated count
			website, err = o.websiteService.GetJobWebsite(ctx, website.ID)
			if err != nil {
				continue
			}
		}

		// Check if limit reached
		if !website.IsLimitReached() {
			available = append(available, website)
		} else {
			o.logger.Info("website limit reached",
				slog.String("website", website.Name),
				slog.Int("currentCount", website.CurrentCount),
				slog.Int("dailyLimit", website.DailyLimit),
			)
		}
	}

	return available, nil
}

// IncrementWebsiteCount increments the count for a website.
func (o *Orchestrator) IncrementWebsiteCount(ctx context.Context, websiteName string) error {
	return o.websiteService.IncrementCount(ctx, websiteName)
}

// ShouldProcessWebsite checks if we should process jobs for a specific website.
func (o *Orchestrator) ShouldProcessWebsite(ctx context.Context, websiteName string) (bool, error) {
	website, err := o.websiteService.GetJobWebsiteByName(ctx, websiteName)
	if err != nil {
		return false, err
	}

	if !website.Enabled {
		return false, nil
	}

	// Check if should reset
	if website.ShouldReset() {
		if err := o.websiteService.ResetCount(ctx, website.ID); err != nil {
			return false, err
		}
		// Reload to get updated count
		website, err = o.websiteService.GetJobWebsiteByName(ctx, websiteName)
		if err != nil {
			return false, err
		}
	}

	return !website.IsLimitReached(), nil
}

