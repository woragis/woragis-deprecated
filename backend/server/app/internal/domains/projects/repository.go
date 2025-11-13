package projects

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for projects.
type Repository interface {
	CreateProject(ctx context.Context, project *Project) error
	UpdateProject(ctx context.Context, project *Project) error
	GetProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*Project, error)
	ListProjects(ctx context.Context, userID uuid.UUID) ([]Project, error)

	CreateMilestone(ctx context.Context, milestone *Milestone) error
	UpdateMilestone(ctx context.Context, milestone *Milestone) error
	ListMilestones(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]Milestone, error)
	GetMilestone(ctx context.Context, milestoneID uuid.UUID, userID uuid.UUID) (*Milestone, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateProject(ctx context.Context, project *Project) error {
	if err := project.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateProject(ctx context.Context, project *Project) error {
	if err := project.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(project).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*Project, error) {
	var project Project
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", projectID, userID).First(&project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrProjectNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &project, nil
}

func (r *gormRepository) ListProjects(ctx context.Context, userID uuid.UUID) ([]Project, error) {
	var projects []Project
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&projects).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return projects, nil
}

func (r *gormRepository) CreateMilestone(ctx context.Context, milestone *Milestone) error {
	if err := milestone.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(milestone).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateMilestone(ctx context.Context, milestone *Milestone) error {
	if err := milestone.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(milestone).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) ListMilestones(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]Milestone, error) {
	var milestones []Milestone
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = milestones.project_id").
		Where("milestones.project_id = ? AND projects.user_id = ?", projectID, userID).
		Order("milestones.due_date asc").
		Find(&milestones).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return milestones, nil
}

func (r *gormRepository) GetMilestone(ctx context.Context, milestoneID uuid.UUID, userID uuid.UUID) (*Milestone, error) {
	var milestone Milestone

	err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = milestones.project_id").
		Where("milestones.id = ? AND projects.user_id = ?", milestoneID, userID).
		First(&milestone).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrMilestoneNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &milestone, nil
}
