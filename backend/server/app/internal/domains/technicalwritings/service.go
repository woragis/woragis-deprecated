package technicalwritings

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates technical writing workflows.
type Service interface {
	CreateTechnicalWriting(ctx context.Context, userID uuid.UUID, req CreateTechnicalWritingRequest) (*TechnicalWriting, error)
	UpdateTechnicalWriting(ctx context.Context, userID, writingID uuid.UUID, req UpdateTechnicalWritingRequest) (*TechnicalWriting, error)
	GetTechnicalWriting(ctx context.Context, writingID uuid.UUID, userID uuid.UUID) (*TechnicalWriting, error)
	GetTechnicalWritingPublic(ctx context.Context, writingID uuid.UUID) (*TechnicalWriting, error)
	ListTechnicalWritings(ctx context.Context, filters ListTechnicalWritingFilters) ([]TechnicalWriting, error)
	ListFeaturedTechnicalWritings(ctx context.Context) ([]TechnicalWriting, error)
	GetWritingsByProject(ctx context.Context, projectID uuid.UUID) ([]TechnicalWriting, error)
	GetWritingsByType(ctx context.Context, writingType WritingType) ([]TechnicalWriting, error)
	GetWritingsByPlatform(ctx context.Context, platform PublicationPlatform) ([]TechnicalWriting, error)
	SearchTechnicalWritings(ctx context.Context, query string) ([]TechnicalWriting, error)
	DeleteTechnicalWriting(ctx context.Context, userID, writingID uuid.UUID) error
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

type CreateTechnicalWritingRequest struct {
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Type          WritingType        `json:"type"`
	Platform      PublicationPlatform `json:"platform"`
	URL           string             `json:"url"`
	Content       string             `json:"content,omitempty"`
	Excerpt       string             `json:"excerpt,omitempty"`
	CanonicalURL  string             `json:"canonicalUrl,omitempty"`
	CoverImageURL string             `json:"coverImageUrl,omitempty"`
	PublishedAt   *time.Time         `json:"publishedAt,omitempty"`
	ReadingTime   int                `json:"readingTime,omitempty"`
	Topics        []string           `json:"topics,omitempty"`
	Technologies  []string           `json:"technologies,omitempty"`
	Views         *int               `json:"views,omitempty"`
	Likes         *int               `json:"likes,omitempty"`
	Shares        *int               `json:"shares,omitempty"`
	Comments      *int               `json:"comments,omitempty"`
	ProjectID     *string            `json:"projectId,omitempty"`
	CaseStudyID   *string            `json:"caseStudyId,omitempty"`
	Featured      bool               `json:"featured,omitempty"`
	DisplayOrder  int                `json:"displayOrder,omitempty"`
}

type UpdateTechnicalWritingRequest struct {
	Title         *string             `json:"title,omitempty"`
	Description   *string             `json:"description,omitempty"`
	Type          *WritingType         `json:"type,omitempty"`
	Platform      *PublicationPlatform `json:"platform,omitempty"`
	URL           *string             `json:"url,omitempty"`
	Content       *string             `json:"content,omitempty"`
	Excerpt       *string             `json:"excerpt,omitempty"`
	CanonicalURL  *string             `json:"canonicalUrl,omitempty"`
	CoverImageURL *string             `json:"coverImageUrl,omitempty"`
	PublishedAt   *time.Time          `json:"publishedAt,omitempty"`
	ReadingTime   *int                `json:"readingTime,omitempty"`
	Topics        []string            `json:"topics,omitempty"`
	Technologies  []string            `json:"technologies,omitempty"`
	Views         *int                `json:"views,omitempty"`
	Likes         *int                `json:"likes,omitempty"`
	Shares        *int                `json:"shares,omitempty"`
	Comments      *int                `json:"comments,omitempty"`
	ProjectID     *string             `json:"projectId,omitempty"`
	CaseStudyID   *string             `json:"caseStudyId,omitempty"`
	Featured      *bool               `json:"featured,omitempty"`
	DisplayOrder  *int                `json:"displayOrder,omitempty"`
}

type ListTechnicalWritingFilters struct {
	UserID    *uuid.UUID
	Type      *WritingType
	Platform  *PublicationPlatform
	ProjectID *uuid.UUID
	Featured  *bool
	Limit     int
	Offset    int
	OrderBy   string
	Order     string
}

// Service methods

func (s *service) CreateTechnicalWriting(ctx context.Context, userID uuid.UUID, req CreateTechnicalWritingRequest) (*TechnicalWriting, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	writing, err := NewTechnicalWriting(userID, req.Title, req.Description, req.Type, req.Platform, req.URL)
	if err != nil {
		return nil, err
	}

	if req.Content != "" {
		writing.Content = req.Content
	}
	if req.Excerpt != "" {
		writing.Excerpt = req.Excerpt
	}
	if req.CanonicalURL != "" || req.CoverImageURL != "" {
		writing.SetURLs(req.URL, req.CanonicalURL, req.CoverImageURL)
	}
	if req.PublishedAt != nil || req.ReadingTime > 0 {
		writing.SetPublicationInfo(req.PublishedAt, req.ReadingTime)
	}
	if len(req.Topics) > 0 {
		writing.SetTopics(req.Topics)
	}
	if len(req.Technologies) > 0 {
		writing.SetTechnologies(req.Technologies)
	}
	if req.Views != nil || req.Likes != nil || req.Shares != nil || req.Comments != nil {
		writing.SetMetrics(req.Views, req.Likes, req.Shares, req.Comments)
	}
	if req.ProjectID != nil {
		projectID, err := uuid.Parse(*req.ProjectID)
		if err != nil {
			return nil, NewDomainError(ErrCodeInvalidPayload, "invalid project id format")
		}
		writing.SetProjectLink(projectID)
	}
	if req.CaseStudyID != nil {
		caseStudyID, err := uuid.Parse(*req.CaseStudyID)
		if err != nil {
			return nil, NewDomainError(ErrCodeInvalidPayload, "invalid case study id format")
		}
		writing.SetCaseStudyLink(caseStudyID)
	}
	writing.Featured = req.Featured
	writing.DisplayOrder = req.DisplayOrder

	if err := writing.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateTechnicalWriting(ctx, writing); err != nil {
		return nil, err
	}

	return writing, nil
}

func (s *service) UpdateTechnicalWriting(ctx context.Context, userID, writingID uuid.UUID, req UpdateTechnicalWritingRequest) (*TechnicalWriting, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if writingID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyWritingID)
	}

	// Get existing writing
	writing, err := s.repo.GetTechnicalWriting(ctx, writingID, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	title := writing.Title
	if req.Title != nil {
		title = *req.Title
	}
	description := writing.Description
	if req.Description != nil {
		description = *req.Description
	}
	content := writing.Content
	if req.Content != nil {
		content = *req.Content
	}
	excerpt := writing.Excerpt
	if req.Excerpt != nil {
		excerpt = *req.Excerpt
	}

	writing.UpdateDetails(title, description, content, excerpt)

	if req.Type != nil {
		writing.Type = *req.Type
	}
	if req.Platform != nil {
		writing.Platform = *req.Platform
	}
	if req.URL != nil || req.CanonicalURL != nil || req.CoverImageURL != nil {
		url := writing.URL
		canonicalURL := writing.CanonicalURL
		coverImageURL := writing.CoverImageURL
		if req.URL != nil {
			url = *req.URL
		}
		if req.CanonicalURL != nil {
			canonicalURL = *req.CanonicalURL
		}
		if req.CoverImageURL != nil {
			coverImageURL = *req.CoverImageURL
		}
		writing.SetURLs(url, canonicalURL, coverImageURL)
	}
	if req.PublishedAt != nil || req.ReadingTime != nil {
		publishedAt := writing.PublishedAt
		readingTime := writing.ReadingTime
		if req.PublishedAt != nil {
			publishedAt = req.PublishedAt
		}
		if req.ReadingTime != nil {
			readingTime = *req.ReadingTime
		}
		writing.SetPublicationInfo(publishedAt, readingTime)
	}
	if req.Topics != nil {
		writing.SetTopics(req.Topics)
	}
	if req.Technologies != nil {
		writing.SetTechnologies(req.Technologies)
	}
	if req.Views != nil || req.Likes != nil || req.Shares != nil || req.Comments != nil {
		views := writing.Views
		likes := writing.Likes
		shares := writing.Shares
		comments := writing.Comments
		if req.Views != nil {
			views = req.Views
		}
		if req.Likes != nil {
			likes = req.Likes
		}
		if req.Shares != nil {
			shares = req.Shares
		}
		if req.Comments != nil {
			comments = req.Comments
		}
		writing.SetMetrics(views, likes, shares, comments)
	}
	if req.ProjectID != nil {
		if *req.ProjectID == "" {
			writing.ClearProjectLink()
		} else {
			projectID, err := uuid.Parse(*req.ProjectID)
			if err != nil {
				return nil, NewDomainError(ErrCodeInvalidPayload, "invalid project id format")
			}
			writing.SetProjectLink(projectID)
		}
	}
	if req.CaseStudyID != nil {
		if *req.CaseStudyID == "" {
			writing.ClearCaseStudyLink()
		} else {
			caseStudyID, err := uuid.Parse(*req.CaseStudyID)
			if err != nil {
				return nil, NewDomainError(ErrCodeInvalidPayload, "invalid case study id format")
			}
			writing.SetCaseStudyLink(caseStudyID)
		}
	}
	if req.Featured != nil {
		writing.SetFeatured(*req.Featured)
	}
	if req.DisplayOrder != nil {
		writing.SetDisplayOrder(*req.DisplayOrder)
	}

	if err := writing.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateTechnicalWriting(ctx, writing); err != nil {
		return nil, err
	}

	return writing, nil
}

func (s *service) GetTechnicalWriting(ctx context.Context, writingID uuid.UUID, userID uuid.UUID) (*TechnicalWriting, error) {
	if writingID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyWritingID)
	}

	return s.repo.GetTechnicalWriting(ctx, writingID, userID)
}

func (s *service) GetTechnicalWritingPublic(ctx context.Context, writingID uuid.UUID) (*TechnicalWriting, error) {
	if writingID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyWritingID)
	}

	return s.repo.GetTechnicalWritingPublic(ctx, writingID)
}

func (s *service) ListTechnicalWritings(ctx context.Context, filters ListTechnicalWritingFilters) ([]TechnicalWriting, error) {
	repoFilters := TechnicalWritingFilters{
		UserID:    filters.UserID,
		Type:      filters.Type,
		Platform:  filters.Platform,
		ProjectID: filters.ProjectID,
		Featured:  filters.Featured,
		Limit:     filters.Limit,
		Offset:    filters.Offset,
		OrderBy:   filters.OrderBy,
		Order:     filters.Order,
	}

	return s.repo.ListTechnicalWritings(ctx, repoFilters)
}

func (s *service) ListFeaturedTechnicalWritings(ctx context.Context) ([]TechnicalWriting, error) {
	return s.repo.ListFeaturedTechnicalWritings(ctx)
}

func (s *service) GetWritingsByProject(ctx context.Context, projectID uuid.UUID) ([]TechnicalWriting, error) {
	if projectID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, "project id cannot be empty")
	}

	return s.repo.GetWritingsByProject(ctx, projectID)
}

func (s *service) GetWritingsByType(ctx context.Context, writingType WritingType) ([]TechnicalWriting, error) {
	return s.repo.GetWritingsByType(ctx, writingType)
}

func (s *service) GetWritingsByPlatform(ctx context.Context, platform PublicationPlatform) ([]TechnicalWriting, error) {
	return s.repo.GetWritingsByPlatform(ctx, platform)
}

func (s *service) SearchTechnicalWritings(ctx context.Context, query string) ([]TechnicalWriting, error) {
	if query == "" {
		return nil, NewDomainError(ErrCodeInvalidPayload, "search query cannot be empty")
	}

	return s.repo.SearchTechnicalWritings(ctx, query)
}

func (s *service) DeleteTechnicalWriting(ctx context.Context, userID, writingID uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if writingID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyWritingID)
	}

	return s.repo.DeleteTechnicalWriting(ctx, writingID, userID)
}

