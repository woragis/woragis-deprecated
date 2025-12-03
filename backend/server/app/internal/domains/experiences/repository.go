package experiences

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for experiences.
type Repository interface {
	CreateExperience(ctx context.Context, experience *Experience) error
	UpdateExperience(ctx context.Context, experience *Experience) error
	GetExperience(ctx context.Context, experienceID uuid.UUID) (*Experience, error)
	ListExperiences(ctx context.Context, filters ExperienceFilters) ([]Experience, error)
	DeleteExperience(ctx context.Context, experienceID uuid.UUID) error
	// Technology methods
	CreateExperienceTechnology(ctx context.Context, tech *ExperienceTechnology) error
	GetExperienceTechnologies(ctx context.Context, experienceID uuid.UUID) ([]ExperienceTechnology, error)
	DeleteExperienceTechnologies(ctx context.Context, experienceID uuid.UUID) error
	// Project methods
	CreateExperienceProject(ctx context.Context, project *ExperienceProject) error
	GetExperienceProjects(ctx context.Context, experienceID uuid.UUID) ([]ExperienceProject, error)
	DeleteExperienceProjects(ctx context.Context, experienceID uuid.UUID) error
	// Achievement methods
	CreateExperienceAchievement(ctx context.Context, achievement *ExperienceAchievement) error
	GetExperienceAchievements(ctx context.Context, experienceID uuid.UUID) ([]ExperienceAchievement, error)
	DeleteExperienceAchievements(ctx context.Context, experienceID uuid.UUID) error
}

// ExperienceFilters represents filtering options for listing experiences.
type ExperienceFilters struct {
	UserID   *uuid.UUID
	Type     *ExperienceType
	IsCurrent *bool
	Limit    int
	Offset   int
	OrderBy  string // "created_at", "updated_at", "display_order", "period_start"
	Order    string // "asc", "desc"
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateExperience(ctx context.Context, experience *Experience) error {
	if experience == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilExperience)
	}

	if err := experience.Validate(); err != nil {
		return err
	}

	now := time.Now()
	experience.CreatedAt = now
	experience.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(experience).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return NewDomainError(ErrCodeConflict, ErrExperienceAlreadyExists)
			}
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) UpdateExperience(ctx context.Context, experience *Experience) error {
	if experience == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilExperience)
	}

	if experience.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyExperienceID)
	}

	if err := experience.Validate(); err != nil {
		return err
	}

	experience.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Model(&Experience{}).
		Where("id = ?", experience.ID).
		Updates(experience)

	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrExperienceNotFound)
	}

	return nil
}

func (r *gormRepository) GetExperience(ctx context.Context, experienceID uuid.UUID) (*Experience, error) {
	if experienceID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyExperienceID)
	}

	var experience Experience
	if err := r.db.WithContext(ctx).Where("id = ?", experienceID).First(&experience).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrExperienceNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &experience, nil
}

func (r *gormRepository) ListExperiences(ctx context.Context, filters ExperienceFilters) ([]Experience, error) {
	query := r.db.WithContext(ctx).Model(&Experience{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.IsCurrent != nil {
		query = query.Where("is_current = ?", *filters.IsCurrent)
	}

	// Default ordering
	orderBy := normalizeOrderBy(filters.OrderBy)
	if orderBy == "" {
		orderBy = "display_order"
	}
	order := filters.Order
	if order == "" {
		order = "desc" // Most recent first by default
	}
	// Validate order direction
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	query = query.Order(orderBy + " " + order)

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var experiences []Experience
	if err := query.Find(&experiences).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return experiences, nil
}

func (r *gormRepository) DeleteExperience(ctx context.Context, experienceID uuid.UUID) error {
	if experienceID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyExperienceID)
	}

	// Delete related data first
	_ = r.DeleteExperienceTechnologies(ctx, experienceID)
	_ = r.DeleteExperienceProjects(ctx, experienceID)
	_ = r.DeleteExperienceAchievements(ctx, experienceID)

	result := r.db.WithContext(ctx).Where("id = ?", experienceID).Delete(&Experience{})
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrExperienceNotFound)
	}

	return nil
}

// Technology methods

func (r *gormRepository) CreateExperienceTechnology(ctx context.Context, tech *ExperienceTechnology) error {
	if tech == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilExperience)
	}

	now := time.Now()
	tech.CreatedAt = now
	tech.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(tech).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) GetExperienceTechnologies(ctx context.Context, experienceID uuid.UUID) ([]ExperienceTechnology, error) {
	if experienceID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyExperienceID)
	}

	var technologies []ExperienceTechnology
	if err := r.db.WithContext(ctx).Where("experience_id = ?", experienceID).Order("technology ASC").Find(&technologies).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return technologies, nil
}

func (r *gormRepository) DeleteExperienceTechnologies(ctx context.Context, experienceID uuid.UUID) error {
	if experienceID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyExperienceID)
	}

	result := r.db.WithContext(ctx).Where("experience_id = ?", experienceID).Delete(&ExperienceTechnology{})
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

// Project methods

func (r *gormRepository) CreateExperienceProject(ctx context.Context, project *ExperienceProject) error {
	if project == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilExperience)
	}

	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) GetExperienceProjects(ctx context.Context, experienceID uuid.UUID) ([]ExperienceProject, error) {
	if experienceID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyExperienceID)
	}

	var projects []ExperienceProject
	if err := r.db.WithContext(ctx).Where("experience_id = ?", experienceID).Order("name ASC").Find(&projects).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return projects, nil
}

func (r *gormRepository) DeleteExperienceProjects(ctx context.Context, experienceID uuid.UUID) error {
	if experienceID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyExperienceID)
	}

	result := r.db.WithContext(ctx).Where("experience_id = ?", experienceID).Delete(&ExperienceProject{})
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

// Achievement methods

func (r *gormRepository) CreateExperienceAchievement(ctx context.Context, achievement *ExperienceAchievement) error {
	if achievement == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilExperience)
	}

	now := time.Now()
	achievement.CreatedAt = now
	achievement.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(achievement).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) GetExperienceAchievements(ctx context.Context, experienceID uuid.UUID) ([]ExperienceAchievement, error) {
	if experienceID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyExperienceID)
	}

	var achievements []ExperienceAchievement
	if err := r.db.WithContext(ctx).Where("experience_id = ?", experienceID).Order("display_order ASC").Find(&achievements).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return achievements, nil
}

func (r *gormRepository) DeleteExperienceAchievements(ctx context.Context, experienceID uuid.UUID) error {
	if experienceID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyExperienceID)
	}

	result := r.db.WithContext(ctx).Where("experience_id = ?", experienceID).Delete(&ExperienceAchievement{})
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

// normalizeOrderBy converts camelCase orderBy values to snake_case database column names
func normalizeOrderBy(orderBy string) string {
	if orderBy == "" {
		return ""
	}

	allowedColumns := map[string]string{
		"createdAt":    "created_at",
		"updatedAt":    "updated_at",
		"displayOrder": "display_order",
		"periodStart":  "period_start",
		"created_at":   "created_at",
		"updated_at":   "updated_at",
		"display_order": "display_order",
		"period_start":  "period_start",
	}

	if normalized, ok := allowedColumns[orderBy]; ok {
		return normalized
	}

	return ""
}

