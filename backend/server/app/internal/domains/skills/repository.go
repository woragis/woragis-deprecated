package skills

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for skills.
type Repository interface {
	CreateSkill(ctx context.Context, skill *Skill) error
	UpdateSkill(ctx context.Context, skill *Skill) error
	DeleteSkill(ctx context.Context, skillID uuid.UUID) error
	GetSkill(ctx context.Context, skillID uuid.UUID) (*Skill, error)
	GetSkillBySlug(ctx context.Context, slug string) (*Skill, error)
	GetSkillByName(ctx context.Context, name string) (*Skill, error)
	ListSkills(ctx context.Context) ([]Skill, error)
	ListSkillsByCategory(ctx context.Context, category SkillCategory) ([]Skill, error)
	SearchSkills(ctx context.Context, query string) ([]Skill, error)

	// Project-Skill relationship operations
	AttachSkillToProject(ctx context.Context, projectID, skillID uuid.UUID) error
	DetachSkillFromProject(ctx context.Context, projectID, skillID uuid.UUID) error
	GetProjectSkills(ctx context.Context, projectID uuid.UUID) ([]Skill, error)
	GetProjectsBySkill(ctx context.Context, skillID uuid.UUID) ([]uuid.UUID, error)
	GetSkillProjectCount(ctx context.Context, skillID uuid.UUID) (int64, error)
	GetAllSkillsWithProjectCounts(ctx context.Context) ([]SkillWithCount, error)
	ProjectHasSkill(ctx context.Context, projectID, skillID uuid.UUID) (bool, error)
	
	// Timeline operations
	GetSkillsTimeline(ctx context.Context) ([]Skill, error) // Ordered by firstUsedDate
}

// SkillWithCount represents a skill with its project count.
type SkillWithCount struct {
	Skill
	ProjectCount int64 `json:"projectCount"`
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateSkill(ctx context.Context, skill *Skill) error {
	if err := skill.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(skill).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrSkillAlreadyExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateSkill(ctx context.Context, skill *Skill) error {
	if err := skill.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(skill).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrSkillAlreadyExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) DeleteSkill(ctx context.Context, skillID uuid.UUID) error {
	// First, delete all project-skill relationships
	if err := r.db.WithContext(ctx).Where("skill_id = ?", skillID).Delete(&ProjectSkill{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
	}

	// Then delete the skill
	if err := r.db.WithContext(ctx).Where("id = ?", skillID).Delete(&Skill{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
	}
	return nil
}

func (r *gormRepository) GetSkill(ctx context.Context, skillID uuid.UUID) (*Skill, error) {
	var skill Skill
	err := r.db.WithContext(ctx).Where("id = ?", skillID).First(&skill).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrSkillNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &skill, nil
}

func (r *gormRepository) GetSkillBySlug(ctx context.Context, slug string) (*Skill, error) {
	var skill Skill
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&skill).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrSkillNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &skill, nil
}

func (r *gormRepository) GetSkillByName(ctx context.Context, name string) (*Skill, error) {
	var skill Skill
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&skill).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrSkillNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &skill, nil
}

func (r *gormRepository) ListSkills(ctx context.Context) ([]Skill, error) {
	var skills []Skill
	if err := r.db.WithContext(ctx).
		Order("name asc").
		Find(&skills).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return skills, nil
}

func (r *gormRepository) ListSkillsByCategory(ctx context.Context, category SkillCategory) ([]Skill, error) {
	var skills []Skill
	if err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Order("name asc").
		Find(&skills).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return skills, nil
}

func (r *gormRepository) SearchSkills(ctx context.Context, query string) ([]Skill, error) {
	var skills []Skill
	pattern := "%" + query + "%"
	if err := r.db.WithContext(ctx).
		Where("LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", pattern, pattern).
		Order("name asc").
		Find(&skills).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return skills, nil
}

// Project-Skill relationship operations

func (r *gormRepository) AttachSkillToProject(ctx context.Context, projectID, skillID uuid.UUID) error {
	projectSkill := ProjectSkill{
		ProjectID: projectID,
		SkillID:   skillID,
		CreatedAt: time.Now().UTC(),
	}

	if err := r.db.WithContext(ctx).Create(&projectSkill).Error; err != nil {
		if isUniqueConstraintError(err) {
			// Already attached, no error
			return nil
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) DetachSkillFromProject(ctx context.Context, projectID, skillID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("project_id = ? AND skill_id = ?", projectID, skillID).
		Delete(&ProjectSkill{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetProjectSkills(ctx context.Context, projectID uuid.UUID) ([]Skill, error) {
	var skills []Skill
	if err := r.db.WithContext(ctx).
		Table("skills").
		Joins("JOIN project_skills ON project_skills.skill_id = skills.id").
		Where("project_skills.project_id = ?", projectID).
		Order("skills.name asc").
		Find(&skills).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return skills, nil
}

func (r *gormRepository) GetProjectsBySkill(ctx context.Context, skillID uuid.UUID) ([]uuid.UUID, error) {
	var projectIDs []uuid.UUID
	if err := r.db.WithContext(ctx).
		Table("project_skills").
		Where("skill_id = ?", skillID).
		Pluck("project_id", &projectIDs).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return projectIDs, nil
}

func (r *gormRepository) GetSkillProjectCount(ctx context.Context, skillID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&ProjectSkill{}).
		Where("skill_id = ?", skillID).
		Count(&count).Error; err != nil {
		return 0, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return count, nil
}

func (r *gormRepository) GetAllSkillsWithProjectCounts(ctx context.Context) ([]SkillWithCount, error) {
	var results []SkillWithCount
	if err := r.db.WithContext(ctx).
		Table("skills").
		Select("skills.*, COUNT(project_skills.project_id) as project_count").
		Joins("LEFT JOIN project_skills ON project_skills.skill_id = skills.id").
		Group("skills.id").
		Order("project_count desc, skills.name asc").
		Scan(&results).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return results, nil
}

func (r *gormRepository) ProjectHasSkill(ctx context.Context, projectID, skillID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&ProjectSkill{}).
		Where("project_id = ? AND skill_id = ?", projectID, skillID).
		Count(&count).Error; err != nil {
		return false, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return count > 0, nil
}

func (r *gormRepository) GetSkillsTimeline(ctx context.Context) ([]Skill, error) {
	var skills []Skill
	if err := r.db.WithContext(ctx).
		Where("first_used_date IS NOT NULL").
		Order("first_used_date asc, name asc").
		Find(&skills).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return skills, nil
}

// isUniqueConstraintError checks if the error is a unique constraint violation.
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// Check for PostgreSQL unique constraint violation (code 23505)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}

	// Check for SQLite unique constraint violation
	if err.Error() == "UNIQUE constraint failed" {
		return true
	}

	return false
}

