package skills

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service orchestrates skill workflows.
type Service interface {
	CreateSkill(ctx context.Context, req CreateSkillRequest) (*Skill, error)
	UpdateSkill(ctx context.Context, req UpdateSkillRequest) (*Skill, error)
	GetSkill(ctx context.Context, skillID uuid.UUID) (*Skill, error)
	GetSkillBySlug(ctx context.Context, slug string) (*Skill, error)
	ListSkills(ctx context.Context) ([]Skill, error)
	ListSkillsByCategory(ctx context.Context, category SkillCategory) ([]Skill, error)
	SearchSkills(ctx context.Context, query string) ([]Skill, error)
	GetAllSkillsWithProjectCounts(ctx context.Context) ([]SkillWithCount, error)

	// Project-Skill relationship operations
	AttachSkillToProject(ctx context.Context, projectID, skillID uuid.UUID) error
	DetachSkillFromProject(ctx context.Context, projectID, skillID uuid.UUID) error
	GetProjectSkills(ctx context.Context, projectID uuid.UUID) ([]Skill, error)
	GetProjectsBySkill(ctx context.Context, skillID uuid.UUID) ([]uuid.UUID, error)
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

type CreateSkillRequest struct {
	Name            string       `json:"name"`
	Description     string       `json:"description,omitempty"`
	Icon            string       `json:"icon,omitempty"`
	Color           string       `json:"color,omitempty"`
	BgGradient      string       `json:"bgGradient,omitempty"`
	BorderColor     string       `json:"borderColor,omitempty"`
	HoverBorderColor string      `json:"hoverBorderColor,omitempty"`
	ShadowColor     string       `json:"shadowColor,omitempty"`
	Category        SkillCategory `json:"category"`
}

type UpdateSkillRequest struct {
	SkillID         uuid.UUID    `json:"-"`
	Name            string        `json:"name,omitempty"`
	Description     string        `json:"description,omitempty"`
	Icon            string        `json:"icon,omitempty"`
	Color           string        `json:"color,omitempty"`
	BgGradient      string        `json:"bgGradient,omitempty"`
	BorderColor     string        `json:"borderColor,omitempty"`
	HoverBorderColor string       `json:"hoverBorderColor,omitempty"`
	ShadowColor     string        `json:"shadowColor,omitempty"`
	Category        SkillCategory `json:"category,omitempty"`
}

// Skill operations

func (s *service) CreateSkill(ctx context.Context, req CreateSkillRequest) (*Skill, error) {
	skill, err := NewSkill(req.Name, req.Description, req.Icon, req.Color, req.BgGradient, req.BorderColor, req.HoverBorderColor, req.ShadowColor, req.Category)
	if err != nil {
		return nil, err
	}

	// Check if skill with same name already exists
	existing, _ := s.repo.GetSkillByName(ctx, req.Name)
	if existing != nil {
		return nil, NewDomainError(ErrCodeConflict, ErrSkillAlreadyExists)
	}

	if err := s.repo.CreateSkill(ctx, skill); err != nil {
		return nil, err
	}

	return skill, nil
}

func (s *service) UpdateSkill(ctx context.Context, req UpdateSkillRequest) (*Skill, error) {
	skill, err := s.repo.GetSkill(ctx, req.SkillID)
	if err != nil {
		return nil, err
	}

	if err := skill.UpdateDetails(req.Name, req.Description, req.Icon, req.Color, req.BgGradient, req.BorderColor, req.HoverBorderColor, req.ShadowColor, req.Category); err != nil {
		return nil, err
	}

	// Check if new name conflicts with existing skill
	if req.Name != "" && req.Name != skill.Name {
		existing, _ := s.repo.GetSkillByName(ctx, req.Name)
		if existing != nil && existing.ID != skill.ID {
			return nil, NewDomainError(ErrCodeConflict, ErrSkillAlreadyExists)
		}
	}

	if err := s.repo.UpdateSkill(ctx, skill); err != nil {
		return nil, err
	}

	return skill, nil
}

func (s *service) GetSkill(ctx context.Context, skillID uuid.UUID) (*Skill, error) {
	return s.repo.GetSkill(ctx, skillID)
}

func (s *service) GetSkillBySlug(ctx context.Context, slug string) (*Skill, error) {
	return s.repo.GetSkillBySlug(ctx, slug)
}

func (s *service) ListSkills(ctx context.Context) ([]Skill, error) {
	return s.repo.ListSkills(ctx)
}

func (s *service) ListSkillsByCategory(ctx context.Context, category SkillCategory) ([]Skill, error) {
	return s.repo.ListSkillsByCategory(ctx, category)
}

func (s *service) SearchSkills(ctx context.Context, query string) ([]Skill, error) {
	return s.repo.SearchSkills(ctx, query)
}

func (s *service) GetAllSkillsWithProjectCounts(ctx context.Context) ([]SkillWithCount, error) {
	return s.repo.GetAllSkillsWithProjectCounts(ctx)
}

// Project-Skill relationship operations

func (s *service) AttachSkillToProject(ctx context.Context, projectID, skillID uuid.UUID) error {
	// Verify skill exists
	if _, err := s.repo.GetSkill(ctx, skillID); err != nil {
		return err
	}

	return s.repo.AttachSkillToProject(ctx, projectID, skillID)
}

func (s *service) DetachSkillFromProject(ctx context.Context, projectID, skillID uuid.UUID) error {
	return s.repo.DetachSkillFromProject(ctx, projectID, skillID)
}

func (s *service) GetProjectSkills(ctx context.Context, projectID uuid.UUID) ([]Skill, error) {
	return s.repo.GetProjectSkills(ctx, projectID)
}

func (s *service) GetProjectsBySkill(ctx context.Context, skillID uuid.UUID) ([]uuid.UUID, error) {
	return s.repo.GetProjectsBySkill(ctx, skillID)
}

