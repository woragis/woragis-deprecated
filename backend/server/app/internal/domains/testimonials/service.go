package testimonials

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates testimonial workflows.
type Service interface {
	CreateTestimonial(ctx context.Context, userID uuid.UUID, req CreateTestimonialRequest) (*Testimonial, error)
	UpdateTestimonial(ctx context.Context, userID, testimonialID uuid.UUID, req UpdateTestimonialRequest) (*Testimonial, error)
	GetTestimonial(ctx context.Context, testimonialID uuid.UUID) (*Testimonial, error)
	ListTestimonials(ctx context.Context, filters ListTestimonialsFilters) ([]Testimonial, error)
	DeleteTestimonial(ctx context.Context, userID, testimonialID uuid.UUID) error
	ApproveTestimonial(ctx context.Context, userID, testimonialID uuid.UUID) (*Testimonial, error)
	RejectTestimonial(ctx context.Context, userID, testimonialID uuid.UUID) (*Testimonial, error)
	HideTestimonial(ctx context.Context, userID, testimonialID uuid.UUID) (*Testimonial, error)
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

type CreateTestimonialRequest struct {
	AuthorName    string `json:"authorName"`
	AuthorRole    string `json:"authorRole,omitempty"`
	AuthorCompany string `json:"authorCompany,omitempty"`
	AuthorPhoto   string `json:"authorPhoto,omitempty"`
	Content       string `json:"content"`
	Rating        *int   `json:"rating,omitempty"`
	LinkedInURL   string `json:"linkedinUrl,omitempty"`
	DisplayOrder  int    `json:"displayOrder,omitempty"`
}

type UpdateTestimonialRequest struct {
	AuthorName    *string `json:"authorName,omitempty"`
	AuthorRole    *string `json:"authorRole,omitempty"`
	AuthorCompany *string `json:"authorCompany,omitempty"`
	AuthorPhoto   *string `json:"authorPhoto,omitempty"`
	Content       *string `json:"content,omitempty"`
	Rating        *int    `json:"rating,omitempty"`
	LinkedInURL   *string `json:"linkedinUrl,omitempty"`
	Status        *TestimonialStatus `json:"status,omitempty"`
	DisplayOrder  *int    `json:"displayOrder,omitempty"`
}

type ListTestimonialsFilters struct {
	UserID *uuid.UUID
	Status *TestimonialStatus
	Limit  int
	Offset int
	OrderBy string
	Order   string
}

// Service methods

func (s *service) CreateTestimonial(ctx context.Context, userID uuid.UUID, req CreateTestimonialRequest) (*Testimonial, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	testimonial := NewTestimonial(userID, req.AuthorName, req.Content)
	
	if req.AuthorRole != "" {
		testimonial.AuthorRole = req.AuthorRole
	}
	if req.AuthorCompany != "" {
		testimonial.AuthorCompany = req.AuthorCompany
	}
	if req.AuthorPhoto != "" {
		testimonial.AuthorPhoto = req.AuthorPhoto
	}
	if req.Rating != nil {
		testimonial.Rating = req.Rating
	}
	if req.LinkedInURL != "" {
		testimonial.LinkedInURL = req.LinkedInURL
	}
	testimonial.DisplayOrder = req.DisplayOrder

	if err := testimonial.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateTestimonial(ctx, testimonial); err != nil {
		return nil, err
	}

	return testimonial, nil
}

func (s *service) UpdateTestimonial(ctx context.Context, userID, testimonialID uuid.UUID, req UpdateTestimonialRequest) (*Testimonial, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if testimonialID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	// Get existing testimonial
	testimonial, err := s.repo.GetTestimonial(ctx, testimonialID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if testimonial.UserID != userID {
		return nil, NewDomainError(ErrCodeUnauthorized, ErrUnauthorized)
	}

	// Update fields
	if req.AuthorName != nil {
		testimonial.AuthorName = *req.AuthorName
	}
	if req.AuthorRole != nil {
		testimonial.AuthorRole = *req.AuthorRole
	}
	if req.AuthorCompany != nil {
		testimonial.AuthorCompany = *req.AuthorCompany
	}
	if req.AuthorPhoto != nil {
		testimonial.AuthorPhoto = *req.AuthorPhoto
	}
	if req.Content != nil {
		testimonial.Content = *req.Content
	}
	if req.Rating != nil {
		testimonial.Rating = req.Rating
	}
	if req.LinkedInURL != nil {
		testimonial.LinkedInURL = *req.LinkedInURL
	}
	if req.Status != nil {
		testimonial.Status = *req.Status
	}
	if req.DisplayOrder != nil {
		testimonial.DisplayOrder = *req.DisplayOrder
	}

	testimonial.UpdatedAt = time.Now()

	if err := testimonial.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateTestimonial(ctx, testimonial); err != nil {
		return nil, err
	}

	return testimonial, nil
}

func (s *service) GetTestimonial(ctx context.Context, testimonialID uuid.UUID) (*Testimonial, error) {
	if testimonialID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	return s.repo.GetTestimonial(ctx, testimonialID)
}

func (s *service) ListTestimonials(ctx context.Context, filters ListTestimonialsFilters) ([]Testimonial, error) {
	repoFilters := TestimonialFilters{
		UserID: filters.UserID,
		Status: filters.Status,
		Limit:  filters.Limit,
		Offset: filters.Offset,
		OrderBy: filters.OrderBy,
		Order:   filters.Order,
	}

	return s.repo.ListTestimonials(ctx, repoFilters)
}

func (s *service) DeleteTestimonial(ctx context.Context, userID, testimonialID uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if testimonialID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	// Verify ownership
	testimonial, err := s.repo.GetTestimonial(ctx, testimonialID)
	if err != nil {
		return err
	}

	if testimonial.UserID != userID {
		return NewDomainError(ErrCodeUnauthorized, ErrUnauthorized)
	}

	return s.repo.DeleteTestimonial(ctx, testimonialID)
}

func (s *service) ApproveTestimonial(ctx context.Context, userID, testimonialID uuid.UUID) (*Testimonial, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if testimonialID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	// Verify ownership
	testimonial, err := s.repo.GetTestimonial(ctx, testimonialID)
	if err != nil {
		return nil, err
	}

	if testimonial.UserID != userID {
		return nil, NewDomainError(ErrCodeUnauthorized, ErrUnauthorized)
	}

	if err := s.repo.ApproveTestimonial(ctx, testimonialID); err != nil {
		return nil, err
	}

	return s.repo.GetTestimonial(ctx, testimonialID)
}

func (s *service) RejectTestimonial(ctx context.Context, userID, testimonialID uuid.UUID) (*Testimonial, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if testimonialID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	// Verify ownership
	testimonial, err := s.repo.GetTestimonial(ctx, testimonialID)
	if err != nil {
		return nil, err
	}

	if testimonial.UserID != userID {
		return nil, NewDomainError(ErrCodeUnauthorized, ErrUnauthorized)
	}

	if err := s.repo.RejectTestimonial(ctx, testimonialID); err != nil {
		return nil, err
	}

	return s.repo.GetTestimonial(ctx, testimonialID)
}

func (s *service) HideTestimonial(ctx context.Context, userID, testimonialID uuid.UUID) (*Testimonial, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if testimonialID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	// Verify ownership
	testimonial, err := s.repo.GetTestimonial(ctx, testimonialID)
	if err != nil {
		return nil, err
	}

	if testimonial.UserID != userID {
		return nil, NewDomainError(ErrCodeUnauthorized, ErrUnauthorized)
	}

	if err := s.repo.HideTestimonial(ctx, testimonialID); err != nil {
		return nil, err
	}

	return s.repo.GetTestimonial(ctx, testimonialID)
}

