package experiences

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates experience workflows.
type Service interface {
	CreateExperience(ctx context.Context, userID uuid.UUID, req CreateExperienceRequest) (*Experience, error)
	UpdateExperience(ctx context.Context, userID, experienceID uuid.UUID, req UpdateExperienceRequest) (*Experience, error)
	GetExperience(ctx context.Context, experienceID uuid.UUID) (*ExperienceWithDetails, error)
	ListExperiences(ctx context.Context, filters ListExperiencesFilters) ([]ExperienceWithDetails, error)
	DeleteExperience(ctx context.Context, userID, experienceID uuid.UUID) error
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

// ExperienceWithDetails includes related entities (technologies, projects, achievements).
type ExperienceWithDetails struct {
	Experience
	Technologies []ExperienceTechnology `json:"technologies"`
	Projects     []ExperienceProject    `json:"projects"`
	Achievements  []ExperienceAchievement `json:"achievements"`
}

// Request payloads

type CreateExperienceRequest struct {
	Company      string          `json:"company"`
	Position     string          `json:"position"`
	PeriodStart  *time.Time      `json:"periodStart,omitempty"`
	PeriodEnd    *time.Time      `json:"periodEnd,omitempty"`
	PeriodText   string          `json:"periodText,omitempty"`
	Location     string          `json:"location,omitempty"`
	Description  string          `json:"description,omitempty"`
	Type         ExperienceType  `json:"type,omitempty"`
	CompanyURL   string          `json:"companyUrl,omitempty"`
	LinkedInURL  string          `json:"linkedinUrl,omitempty"`
	DisplayOrder int             `json:"displayOrder,omitempty"`
	IsCurrent    bool            `json:"isCurrent,omitempty"`
	Technologies []string        `json:"technologies,omitempty"`
	Projects     []ProjectInput  `json:"projects,omitempty"`
	Achievements []AchievementInput `json:"achievements,omitempty"`
}

type UpdateExperienceRequest struct {
	Company      *string         `json:"company,omitempty"`
	Position     *string         `json:"position,omitempty"`
	PeriodStart  *time.Time      `json:"periodStart,omitempty"`
	PeriodEnd    *time.Time      `json:"periodEnd,omitempty"`
	PeriodText   *string         `json:"periodText,omitempty"`
	Location     *string         `json:"location,omitempty"`
	Description  *string         `json:"description,omitempty"`
	Type         *ExperienceType `json:"type,omitempty"`
	CompanyURL   *string         `json:"companyUrl,omitempty"`
	LinkedInURL  *string         `json:"linkedinUrl,omitempty"`
	DisplayOrder *int            `json:"displayOrder,omitempty"`
	IsCurrent    *bool           `json:"isCurrent,omitempty"`
	Technologies []string        `json:"technologies,omitempty"`
	Projects     []ProjectInput  `json:"projects,omitempty"`
	Achievements []AchievementInput `json:"achievements,omitempty"`
}

type ProjectInput struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type AchievementInput struct {
	Metric       string `json:"metric"`
	Description  string `json:"description"`
	Icon         string `json:"icon,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

type ListExperiencesFilters struct {
	UserID    *uuid.UUID
	Type      *ExperienceType
	IsCurrent *bool
	Limit     int
	Offset    int
	OrderBy   string
	Order     string
}

// Service methods

func (s *service) CreateExperience(ctx context.Context, userID uuid.UUID, req CreateExperienceRequest) (*Experience, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	experience := NewExperience(userID, req.Company, req.Position)

	if req.PeriodStart != nil {
		experience.PeriodStart = req.PeriodStart
	}
	if req.PeriodEnd != nil {
		experience.PeriodEnd = req.PeriodEnd
	}
	if req.PeriodText != "" {
		experience.SetPeriodText(req.PeriodText)
	}
	if req.Location != "" {
		experience.Location = req.Location
	}
	if req.Description != "" {
		experience.SetDescription(req.Description)
	}
	if req.Type != "" {
		if err := experience.SetType(req.Type); err != nil {
			return nil, err
		}
	}
	if req.CompanyURL != "" {
		experience.CompanyURL = req.CompanyURL
	}
	if req.LinkedInURL != "" {
		experience.LinkedInURL = req.LinkedInURL
	}
	experience.DisplayOrder = req.DisplayOrder
	experience.SetIsCurrent(req.IsCurrent)

	if err := experience.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateExperience(ctx, experience); err != nil {
		return nil, err
	}

	// Create technologies
	for _, techName := range req.Technologies {
		tech := &ExperienceTechnology{
			ID:           uuid.New(),
			ExperienceID: experience.ID,
			Technology:   techName,
		}
		if err := s.repo.CreateExperienceTechnology(ctx, tech); err != nil {
			return nil, err
		}
	}

	// Create projects
	for _, projectInput := range req.Projects {
		project := &ExperienceProject{
			ID:           uuid.New(),
			ExperienceID: experience.ID,
			Name:         projectInput.Name,
			URL:          projectInput.URL,
		}
		if err := s.repo.CreateExperienceProject(ctx, project); err != nil {
			return nil, err
		}
	}

	// Create achievements
	for _, achievementInput := range req.Achievements {
		achievement := &ExperienceAchievement{
			ID:           uuid.New(),
			ExperienceID: experience.ID,
			Metric:       achievementInput.Metric,
			Description:  achievementInput.Description,
			Icon:         achievementInput.Icon,
			DisplayOrder: achievementInput.DisplayOrder,
		}
		if err := s.repo.CreateExperienceAchievement(ctx, achievement); err != nil {
			return nil, err
		}
	}

	return experience, nil
}

func (s *service) UpdateExperience(ctx context.Context, userID, experienceID uuid.UUID, req UpdateExperienceRequest) (*Experience, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if experienceID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyExperienceID)
	}

	// Get existing experience
	experience, err := s.repo.GetExperience(ctx, experienceID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if experience.UserID != userID {
		return nil, NewDomainError(ErrCodeUnauthorized, ErrUnauthorized)
	}

	// Update fields
	if req.Company != nil {
		experience.Company = *req.Company
	}
	if req.Position != nil {
		experience.Position = *req.Position
	}
	if req.PeriodStart != nil {
		experience.PeriodStart = req.PeriodStart
	}
	if req.PeriodEnd != nil {
		experience.PeriodEnd = req.PeriodEnd
	}
	if req.PeriodText != nil {
		experience.SetPeriodText(*req.PeriodText)
	}
	if req.Location != nil {
		experience.Location = *req.Location
	}
	if req.Description != nil {
		experience.SetDescription(*req.Description)
	}
	if req.Type != nil {
		if err := experience.SetType(*req.Type); err != nil {
			return nil, err
		}
	}
	if req.CompanyURL != nil {
		experience.CompanyURL = *req.CompanyURL
	}
	if req.LinkedInURL != nil {
		experience.LinkedInURL = *req.LinkedInURL
	}
	if req.DisplayOrder != nil {
		experience.DisplayOrder = *req.DisplayOrder
	}
	if req.IsCurrent != nil {
		experience.SetIsCurrent(*req.IsCurrent)
	}

	experience.UpdatedAt = time.Now()

	if err := experience.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateExperience(ctx, experience); err != nil {
		return nil, err
	}

	// Update related entities if provided
	if req.Technologies != nil {
		// Delete existing technologies
		if err := s.repo.DeleteExperienceTechnologies(ctx, experienceID); err != nil {
			return nil, err
		}
		// Create new technologies
		for _, techName := range req.Technologies {
			tech := &ExperienceTechnology{
				ID:           uuid.New(),
				ExperienceID: experienceID,
				Technology:   techName,
			}
			if err := s.repo.CreateExperienceTechnology(ctx, tech); err != nil {
				return nil, err
			}
		}
	}

	if req.Projects != nil {
		// Delete existing projects
		if err := s.repo.DeleteExperienceProjects(ctx, experienceID); err != nil {
			return nil, err
		}
		// Create new projects
		for _, projectInput := range req.Projects {
			project := &ExperienceProject{
				ID:           uuid.New(),
				ExperienceID: experienceID,
				Name:         projectInput.Name,
				URL:          projectInput.URL,
			}
			if err := s.repo.CreateExperienceProject(ctx, project); err != nil {
				return nil, err
			}
		}
	}

	if req.Achievements != nil {
		// Delete existing achievements
		if err := s.repo.DeleteExperienceAchievements(ctx, experienceID); err != nil {
			return nil, err
		}
		// Create new achievements
		for _, achievementInput := range req.Achievements {
			achievement := &ExperienceAchievement{
				ID:           uuid.New(),
				ExperienceID: experienceID,
				Metric:       achievementInput.Metric,
				Description:  achievementInput.Description,
				Icon:         achievementInput.Icon,
				DisplayOrder: achievementInput.DisplayOrder,
			}
			if err := s.repo.CreateExperienceAchievement(ctx, achievement); err != nil {
				return nil, err
			}
		}
	}

	return experience, nil
}

func (s *service) GetExperience(ctx context.Context, experienceID uuid.UUID) (*ExperienceWithDetails, error) {
	if experienceID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyExperienceID)
	}

	experience, err := s.repo.GetExperience(ctx, experienceID)
	if err != nil {
		return nil, err
	}

	technologies, _ := s.repo.GetExperienceTechnologies(ctx, experienceID)
	projects, _ := s.repo.GetExperienceProjects(ctx, experienceID)
	achievements, _ := s.repo.GetExperienceAchievements(ctx, experienceID)

	return &ExperienceWithDetails{
		Experience:   *experience,
		Technologies: technologies,
		Projects:     projects,
		Achievements:  achievements,
	}, nil
}

func (s *service) ListExperiences(ctx context.Context, filters ListExperiencesFilters) ([]ExperienceWithDetails, error) {
	repoFilters := ExperienceFilters{
		UserID:    filters.UserID,
		Type:      filters.Type,
		IsCurrent: filters.IsCurrent,
		Limit:     filters.Limit,
		Offset:    filters.Offset,
		OrderBy:   filters.OrderBy,
		Order:     filters.Order,
	}

	experiences, err := s.repo.ListExperiences(ctx, repoFilters)
	if err != nil {
		return nil, err
	}

	result := make([]ExperienceWithDetails, len(experiences))
	for i, exp := range experiences {
		technologies, _ := s.repo.GetExperienceTechnologies(ctx, exp.ID)
		projects, _ := s.repo.GetExperienceProjects(ctx, exp.ID)
		achievements, _ := s.repo.GetExperienceAchievements(ctx, exp.ID)

		result[i] = ExperienceWithDetails{
			Experience:   exp,
			Technologies: technologies,
			Projects:     projects,
			Achievements:  achievements,
		}
	}

	return result, nil
}

func (s *service) DeleteExperience(ctx context.Context, userID, experienceID uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if experienceID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyExperienceID)
	}

	// Verify ownership
	experience, err := s.repo.GetExperience(ctx, experienceID)
	if err != nil {
		return err
	}

	if experience.UserID != userID {
		return NewDomainError(ErrCodeUnauthorized, ErrUnauthorized)
	}

	return s.repo.DeleteExperience(ctx, experienceID)
}

