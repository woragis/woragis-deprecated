package projects

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for projects.
type Repository interface {
	CreateProject(ctx context.Context, project *Project) error
	UpdateProject(ctx context.Context, project *Project) error
	GetProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*Project, error)
	GetProjectBySlug(ctx context.Context, slug string, userID uuid.UUID) (*Project, error)
	SearchProjectsBySlug(ctx context.Context, slug string, userID uuid.UUID) ([]Project, error)
	ListProjects(ctx context.Context, userID uuid.UUID) ([]Project, error)

	CreateMilestone(ctx context.Context, milestone *Milestone) error
	UpdateMilestone(ctx context.Context, milestone *Milestone) error
	BulkUpdateMilestones(ctx context.Context, milestones []*Milestone) error
	CreateMilestones(ctx context.Context, milestones []*Milestone) error
	ListMilestones(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]Milestone, error)
	GetMilestone(ctx context.Context, milestoneID uuid.UUID, userID uuid.UUID) (*Milestone, error)

	CreateKanbanColumn(ctx context.Context, column *KanbanColumn) error
	UpdateKanbanColumn(ctx context.Context, column *KanbanColumn) error
	DeleteKanbanColumn(ctx context.Context, columnID uuid.UUID, userID uuid.UUID) error
	ListKanbanColumns(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]KanbanColumn, error)
	GetKanbanColumn(ctx context.Context, columnID uuid.UUID, userID uuid.UUID) (*KanbanColumn, error)

	CreateKanbanCard(ctx context.Context, card *KanbanCard) error
	UpdateKanbanCard(ctx context.Context, card *KanbanCard) error
	DeleteKanbanCard(ctx context.Context, cardID uuid.UUID, userID uuid.UUID) error
	ListKanbanCards(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]KanbanCard, error)
	GetKanbanCard(ctx context.Context, cardID uuid.UUID, userID uuid.UUID) (*KanbanCard, error)

	CreateDependency(ctx context.Context, dependency *ProjectDependency) error
	DeleteDependency(ctx context.Context, dependencyID uuid.UUID, userID uuid.UUID) error
	ListDependencies(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectDependency, error)
	GetDependency(ctx context.Context, dependencyID uuid.UUID, userID uuid.UUID) (*ProjectDependency, error)
	DependencyExists(ctx context.Context, projectID, dependsOn uuid.UUID) (bool, error)

	CreateProjectWithRelated(ctx context.Context, project *Project, columns []*KanbanColumn, cards []*KanbanCard, milestones []*Milestone) error
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

func (r *gormRepository) GetProjectBySlug(ctx context.Context, slug string, userID uuid.UUID) (*Project, error) {
	var project Project
	err := r.db.WithContext(ctx).
		Where("slug = ? AND user_id = ?", slug, userID).
		First(&project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrProjectNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &project, nil
}

func (r *gormRepository) SearchProjectsBySlug(ctx context.Context, slug string, userID uuid.UUID) ([]Project, error) {
	query := strings.TrimSpace(strings.ToLower(slug))
	if query == "" {
		return []Project{}, nil
	}

	var projects []Project
	pattern := "%" + query + "%"
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND LOWER(slug) LIKE ?", userID, pattern).
		Order("created_at desc").
		Find(&projects).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return projects, nil
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

func (r *gormRepository) BulkUpdateMilestones(ctx context.Context, milestones []*Milestone) error {
	if len(milestones) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, milestone := range milestones {
			if err := milestone.Validate(); err != nil {
				return err
			}

			if err := tx.Save(milestone).Error; err != nil {
				return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
			}
		}
		return nil
	})
}

func (r *gormRepository) CreateMilestones(ctx context.Context, milestones []*Milestone) error {
	if len(milestones) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, milestone := range milestones {
			if err := milestone.Validate(); err != nil {
				return err
			}
		}

		if err := tx.Create(&milestones).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
		}
		return nil
	})
}

func (r *gormRepository) ListMilestones(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]Milestone, error) {
	var milestones []Milestone
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = milestones.project_id").
		Where("milestones.project_id = ? AND projects.user_id = ?", projectID, userID).
		Order("milestones.due_date asc, milestones.created_at asc").
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

func (r *gormRepository) CreateKanbanColumn(ctx context.Context, column *KanbanColumn) error {
	if err := column.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(column).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateKanbanColumn(ctx context.Context, column *KanbanColumn) error {
	if err := column.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(column).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) DeleteKanbanColumn(ctx context.Context, columnID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var column KanbanColumn
		if err := tx.Joins("JOIN projects ON projects.id = kanban_columns.project_id").
			Where("kanban_columns.id = ? AND projects.user_id = ?", columnID, userID).
			First(&column).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NewDomainError(ErrCodeNotFound, ErrKanbanColumnNotFound)
			}
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
		}

		if err := tx.Where("column_id = ?", columnID).Delete(&KanbanCard{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}

		if err := tx.Delete(&column).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}
		return nil
	})
}

func (r *gormRepository) ListKanbanColumns(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]KanbanColumn, error) {
	var columns []KanbanColumn
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = kanban_columns.project_id").
		Where("kanban_columns.project_id = ? AND projects.user_id = ?", projectID, userID).
		Order("kanban_columns.position asc, kanban_columns.created_at asc").
		Find(&columns).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return columns, nil
}

func (r *gormRepository) GetKanbanColumn(ctx context.Context, columnID uuid.UUID, userID uuid.UUID) (*KanbanColumn, error) {
	var column KanbanColumn
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = kanban_columns.project_id").
		Where("kanban_columns.id = ? AND projects.user_id = ?", columnID, userID).
		First(&column).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrKanbanColumnNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &column, nil
}

func (r *gormRepository) CreateKanbanCard(ctx context.Context, card *KanbanCard) error {
	if err := card.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(card).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateKanbanCard(ctx context.Context, card *KanbanCard) error {
	if err := card.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(card).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) DeleteKanbanCard(ctx context.Context, cardID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var card KanbanCard
		if err := tx.Joins("JOIN projects ON projects.id = kanban_cards.project_id").
			Where("kanban_cards.id = ? AND projects.user_id = ?", cardID, userID).
			First(&card).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NewDomainError(ErrCodeNotFound, ErrKanbanCardNotFound)
			}
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
		}

		if err := tx.Delete(&card).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}
		return nil
	})
}

func (r *gormRepository) ListKanbanCards(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]KanbanCard, error) {
	var cards []KanbanCard
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = kanban_cards.project_id").
		Where("kanban_cards.project_id = ? AND projects.user_id = ?", projectID, userID).
		Order("kanban_cards.column_id asc, kanban_cards.position asc, kanban_cards.created_at asc").
		Find(&cards).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return cards, nil
}

func (r *gormRepository) GetKanbanCard(ctx context.Context, cardID uuid.UUID, userID uuid.UUID) (*KanbanCard, error) {
	var card KanbanCard
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = kanban_cards.project_id").
		Where("kanban_cards.id = ? AND projects.user_id = ?", cardID, userID).
		First(&card).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrKanbanCardNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &card, nil
}

func (r *gormRepository) CreateDependency(ctx context.Context, dependency *ProjectDependency) error {
	if err := dependency.Validate(); err != nil {
		return err
	}

	exists, err := r.DependencyExists(ctx, dependency.ProjectID, dependency.DependsOnProjectID)
	if err != nil {
		return err
	}
	if exists {
		return NewDomainError(ErrCodeConflict, ErrDependencyAlreadyExists)
	}

	if err := r.db.WithContext(ctx).Create(dependency).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) DeleteDependency(ctx context.Context, dependencyID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var dep ProjectDependency
		if err := tx.Table("project_dependencies").
			Select("project_dependencies.*").
			Joins("JOIN projects p1 ON p1.id = project_dependencies.project_id").
			Where("project_dependencies.id = ? AND p1.user_id = ?", dependencyID, userID).
			First(&dep).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NewDomainError(ErrCodeNotFound, ErrDependencyNotFound)
			}
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
		}

		if err := tx.Delete(&dep).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}
		return nil
	})
}

func (r *gormRepository) ListDependencies(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectDependency, error) {
	var deps []ProjectDependency
	if err := r.db.WithContext(ctx).
		Table("project_dependencies").
		Select("project_dependencies.*").
		Joins("JOIN projects p1 ON p1.id = project_dependencies.project_id").
		Joins("JOIN projects p2 ON p2.id = project_dependencies.depends_on_project_id").
		Where("project_dependencies.project_id = ? AND p1.user_id = ?", projectID, userID).
		Find(&deps).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return deps, nil
}

func (r *gormRepository) GetDependency(ctx context.Context, dependencyID uuid.UUID, userID uuid.UUID) (*ProjectDependency, error) {
	var dep ProjectDependency
	if err := r.db.WithContext(ctx).
		Table("project_dependencies").
		Select("project_dependencies.*").
		Joins("JOIN projects p1 ON p1.id = project_dependencies.project_id").
		Where("project_dependencies.id = ? AND p1.user_id = ?", dependencyID, userID).
		First(&dep).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrDependencyNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &dep, nil
}

func (r *gormRepository) DependencyExists(ctx context.Context, projectID, dependsOn uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&ProjectDependency{}).
		Where("project_id = ? AND depends_on_project_id = ?", projectID, dependsOn).
		Count(&count).Error; err != nil {
		return false, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return count > 0, nil
}

func (r *gormRepository) CreateProjectWithRelated(ctx context.Context, project *Project, columns []*KanbanColumn, cards []*KanbanCard, milestones []*Milestone) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := project.Validate(); err != nil {
			return err
		}

		if err := tx.Create(project).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
		}

		if len(columns) > 0 {
			for _, column := range columns {
				if err := column.Validate(); err != nil {
					return err
				}
			}
			if err := tx.Create(&columns).Error; err != nil {
				return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
			}
		}

		if len(cards) > 0 {
			for _, card := range cards {
				if err := card.Validate(); err != nil {
					return err
				}
			}
			if err := tx.Create(&cards).Error; err != nil {
				return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
			}
		}

		if len(milestones) > 0 {
			for _, milestone := range milestones {
				if err := milestone.Validate(); err != nil {
					return err
				}
			}
			if err := tx.Create(&milestones).Error; err != nil {
				return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
			}
		}

		return nil
	})
}
