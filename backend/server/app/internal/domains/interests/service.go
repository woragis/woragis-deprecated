package interests

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service orchestrates interest workflows.
type Service interface {
	CreateInterest(ctx context.Context, req CreateInterestRequest) (*Interest, error)
	UpdateInterest(ctx context.Context, req UpdateInterestRequest) (*Interest, error)
	GetInterest(ctx context.Context, interestID uuid.UUID) (*Interest, error)
	GetInterestBySlug(ctx context.Context, slug string) (*Interest, error)
	ListInterests(ctx context.Context) ([]Interest, error)
	ListFeaturedInterests(ctx context.Context) ([]Interest, error)
	SearchInterests(ctx context.Context, query string) ([]Interest, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

var _ Service = (*service)(nil)

// NewService constructs a Service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// Request payloads

type CreateInterestRequest struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	Icon             string `json:"icon,omitempty"`
	Color            string `json:"color,omitempty"`
	BgGradient       string `json:"bgGradient,omitempty"`
	BorderColor      string `json:"borderColor,omitempty"`
	HoverBorderColor string `json:"hoverBorderColor,omitempty"`
	ShadowColor      string `json:"shadowColor,omitempty"`
	FullWidth        bool   `json:"fullWidth"`
	Featured         bool   `json:"featured"`
}

type UpdateInterestRequest struct {
	InterestID       uuid.UUID `json:"-"`
	Title            string    `json:"title,omitempty"`
	Description      string    `json:"description,omitempty"`
	Icon             string    `json:"icon,omitempty"`
	Color            string    `json:"color,omitempty"`
	BgGradient       string    `json:"bgGradient,omitempty"`
	BorderColor      string    `json:"borderColor,omitempty"`
	HoverBorderColor string    `json:"hoverBorderColor,omitempty"`
	ShadowColor      string    `json:"shadowColor,omitempty"`
	FullWidth        *bool     `json:"fullWidth,omitempty"`
	Featured         *bool     `json:"featured,omitempty"`
}

// Interest operations

func (s *service) CreateInterest(ctx context.Context, req CreateInterestRequest) (*Interest, error) {
	interest, err := NewInterest(req.Title, req.Description, req.Icon, req.Color, req.BgGradient, req.BorderColor, req.HoverBorderColor, req.ShadowColor, req.FullWidth, req.Featured)
	if err != nil {
		return nil, err
	}

	// Check if interest with same title already exists
	existing, _ := s.repo.GetInterestByTitle(ctx, req.Title)
	if existing != nil {
		return nil, NewDomainError(ErrCodeConflict, ErrInterestAlreadyExists)
	}

	if err := s.repo.CreateInterest(ctx, interest); err != nil {
		return nil, err
	}

	return interest, nil
}

func (s *service) UpdateInterest(ctx context.Context, req UpdateInterestRequest) (*Interest, error) {
	interest, err := s.repo.GetInterest(ctx, req.InterestID)
	if err != nil {
		return nil, err
	}

	if err := interest.UpdateDetails(req.Title, req.Description, req.Icon, req.Color, req.BgGradient, req.BorderColor, req.HoverBorderColor, req.ShadowColor, req.FullWidth, req.Featured); err != nil {
		return nil, err
	}

	// Check if new title conflicts with existing interest
	if req.Title != "" && req.Title != interest.Title {
		existing, _ := s.repo.GetInterestByTitle(ctx, req.Title)
		if existing != nil && existing.ID != interest.ID {
			return nil, NewDomainError(ErrCodeConflict, ErrInterestAlreadyExists)
		}
	}

	if err := s.repo.UpdateInterest(ctx, interest); err != nil {
		return nil, err
	}

	return interest, nil
}

func (s *service) GetInterest(ctx context.Context, interestID uuid.UUID) (*Interest, error) {
	return s.repo.GetInterest(ctx, interestID)
}

func (s *service) GetInterestBySlug(ctx context.Context, slug string) (*Interest, error) {
	return s.repo.GetInterestBySlug(ctx, slug)
}

func (s *service) ListInterests(ctx context.Context) ([]Interest, error) {
	return s.repo.ListInterests(ctx)
}

func (s *service) ListFeaturedInterests(ctx context.Context) ([]Interest, error) {
	return s.repo.ListFeaturedInterests(ctx)
}

func (s *service) SearchInterests(ctx context.Context, query string) ([]Interest, error) {
	return s.repo.SearchInterests(ctx, query)
}

