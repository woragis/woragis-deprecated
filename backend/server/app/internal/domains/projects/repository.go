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
	DeleteProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error
	GetProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*Project, error)
	GetProjectBySlug(ctx context.Context, slug string, userID uuid.UUID) (*Project, error)
	GetProjectBySlugPublic(ctx context.Context, slug string) (*Project, error)
	SearchProjectsBySlug(ctx context.Context, slug string, userID uuid.UUID) ([]Project, error)
	IsProjectSlugTaken(ctx context.Context, userID uuid.UUID, slug string, excludeID uuid.UUID) (bool, error)
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

	// Documentation operations
	CreateDocumentation(ctx context.Context, doc *ProjectDocumentation) error
	UpdateDocumentation(ctx context.Context, doc *ProjectDocumentation) error
	GetDocumentation(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*ProjectDocumentation, error)
	GetDocumentationByProjectID(ctx context.Context, projectID uuid.UUID) (*ProjectDocumentation, error)
	DeleteDocumentation(ctx context.Context, documentationID uuid.UUID, userID uuid.UUID) error

	// Documentation Section operations
	CreateDocumentationSection(ctx context.Context, section *DocumentationSection) error
	UpdateDocumentationSection(ctx context.Context, section *DocumentationSection) error
	DeleteDocumentationSection(ctx context.Context, sectionID uuid.UUID, userID uuid.UUID) error
	ListDocumentationSections(ctx context.Context, documentationID uuid.UUID, userID uuid.UUID) ([]DocumentationSection, error)
	GetDocumentationSection(ctx context.Context, sectionID uuid.UUID, userID uuid.UUID) (*DocumentationSection, error)
	BulkUpdateDocumentationSections(ctx context.Context, sections []*DocumentationSection) error

	// Technology operations
	CreateTechnology(ctx context.Context, tech *ProjectTechnology) error
	UpdateTechnology(ctx context.Context, tech *ProjectTechnology) error
	DeleteTechnology(ctx context.Context, techID uuid.UUID, userID uuid.UUID) error
	ListTechnologies(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectTechnology, error)
	GetTechnology(ctx context.Context, techID uuid.UUID, userID uuid.UUID) (*ProjectTechnology, error)
	BulkCreateTechnologies(ctx context.Context, technologies []*ProjectTechnology) error
	BulkUpdateTechnologies(ctx context.Context, technologies []*ProjectTechnology) error

	// File Structure operations
	CreateFileStructure(ctx context.Context, fs *ProjectFileStructure) error
	UpdateFileStructure(ctx context.Context, fs *ProjectFileStructure) error
	DeleteFileStructure(ctx context.Context, fsID uuid.UUID, userID uuid.UUID) error
	ListFileStructures(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectFileStructure, error)
	GetFileStructure(ctx context.Context, fsID uuid.UUID, userID uuid.UUID) (*ProjectFileStructure, error)
	BulkCreateFileStructures(ctx context.Context, structures []*ProjectFileStructure) error
	BulkUpdateFileStructures(ctx context.Context, structures []*ProjectFileStructure) error
	DeleteFileStructuresByProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error

	// Architecture Diagram operations
	CreateArchitectureDiagram(ctx context.Context, diagram *ProjectArchitectureDiagram) error
	UpdateArchitectureDiagram(ctx context.Context, diagram *ProjectArchitectureDiagram) error
	DeleteArchitectureDiagram(ctx context.Context, diagramID uuid.UUID, userID uuid.UUID) error
	ListArchitectureDiagrams(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectArchitectureDiagram, error)
	GetArchitectureDiagram(ctx context.Context, diagramID uuid.UUID, userID uuid.UUID) (*ProjectArchitectureDiagram, error)
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

func (r *gormRepository) DeleteProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	// Verify project exists and belongs to user
	_, err := r.GetProject(ctx, projectID, userID)
	if err != nil {
		return err
	}

	// Delete all related resources in a transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete milestones
		if err := tx.Where("project_id = ?", projectID).Delete(&Milestone{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
		}

		// Delete kanban cards (must be before columns)
		if err := tx.Where("project_id = ?", projectID).Delete(&KanbanCard{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
		}

		// Delete kanban columns
		if err := tx.Where("project_id = ?", projectID).Delete(&KanbanColumn{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
		}

		// Delete dependencies
		if err := tx.Where("project_id = ? OR depends_on = ?", projectID, projectID).Delete(&ProjectDependency{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
		}

		// Delete documentation sections
		var doc *ProjectDocumentation
		if err := tx.Where("project_id = ?", projectID).First(&doc).Error; err == nil {
			if err := tx.Where("documentation_id = ?", doc.ID).Delete(&DocumentationSection{}).Error; err != nil {
				return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
			}
		}

		// Delete documentation
		if err := tx.Where("project_id = ?", projectID).Delete(&ProjectDocumentation{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
		}

		// Delete technologies
		if err := tx.Where("project_id = ?", projectID).Delete(&ProjectTechnology{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
		}

		// Delete file structures (recursive delete handled by foreign key cascade or manual deletion)
		if err := tx.Where("project_id = ?", projectID).Delete(&ProjectFileStructure{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
		}

		// Delete architecture diagrams
		if err := tx.Where("project_id = ?", projectID).Delete(&ProjectArchitectureDiagram{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
		}

		// Delete project-skill relationships (table name: project_skills)
		if err := tx.Exec("DELETE FROM project_skills WHERE project_id = ?", projectID).Error; err != nil {
			// Log but don't fail - relationships might be managed elsewhere
			_ = err
		}

		// Finally delete the project
		if err := tx.Where("id = ? AND user_id = ?", projectID, userID).Delete(&Project{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
		}

		return nil
	})
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

func (r *gormRepository) GetProjectBySlugPublic(ctx context.Context, slug string) (*Project, error) {
	var project Project
	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
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

func (r *gormRepository) IsProjectSlugTaken(ctx context.Context, userID uuid.UUID, slug string, excludeID uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&Project{}).
		Where("user_id = ? AND slug = ?", userID, slug)
	if excludeID != uuid.Nil {
		query = query.Where("id <> ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return count > 0, nil
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

// Documentation operations

func (r *gormRepository) CreateDocumentation(ctx context.Context, doc *ProjectDocumentation) error {
	if err := doc.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(doc).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateDocumentation(ctx context.Context, doc *ProjectDocumentation) error {
	if err := doc.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Save(doc).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetDocumentation(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*ProjectDocumentation, error) {
	var doc ProjectDocumentation
	err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = project_documentations.project_id").
		Where("project_documentations.project_id = ? AND projects.user_id = ?", projectID, userID).
		First(&doc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrDocumentationNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &doc, nil
}

func (r *gormRepository) GetDocumentationByProjectID(ctx context.Context, projectID uuid.UUID) (*ProjectDocumentation, error) {
	var doc ProjectDocumentation
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		First(&doc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrDocumentationNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &doc, nil
}

func (r *gormRepository) DeleteDocumentation(ctx context.Context, documentationID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var doc ProjectDocumentation
		if err := tx.Joins("JOIN projects ON projects.id = project_documentations.project_id").
			Where("project_documentations.id = ? AND projects.user_id = ?", documentationID, userID).
			First(&doc).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NewDomainError(ErrCodeNotFound, ErrDocumentationNotFound)
			}
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
		}

		if err := tx.Where("documentation_id = ?", documentationID).Delete(&DocumentationSection{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}

		if err := tx.Delete(&doc).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}
		return nil
	})
}

// Documentation Section operations

func (r *gormRepository) CreateDocumentationSection(ctx context.Context, section *DocumentationSection) error {
	if err := section.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(section).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateDocumentationSection(ctx context.Context, section *DocumentationSection) error {
	if err := section.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Save(section).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) DeleteDocumentationSection(ctx context.Context, sectionID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var section DocumentationSection
		if err := tx.Joins("JOIN project_documentations ON project_documentations.id = documentation_sections.documentation_id").
			Joins("JOIN projects ON projects.id = project_documentations.project_id").
			Where("documentation_sections.id = ? AND projects.user_id = ?", sectionID, userID).
			First(&section).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NewDomainError(ErrCodeNotFound, ErrSectionNotFound)
			}
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
		}

		if err := tx.Delete(&section).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}
		return nil
	})
}

func (r *gormRepository) ListDocumentationSections(ctx context.Context, documentationID uuid.UUID, userID uuid.UUID) ([]DocumentationSection, error) {
	var sections []DocumentationSection
	if err := r.db.WithContext(ctx).
		Joins("JOIN project_documentations ON project_documentations.id = documentation_sections.documentation_id").
		Joins("JOIN projects ON projects.id = project_documentations.project_id").
		Where("documentation_sections.documentation_id = ? AND projects.user_id = ?", documentationID, userID).
		Order("documentation_sections.position asc, documentation_sections.created_at asc").
		Find(&sections).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return sections, nil
}

func (r *gormRepository) GetDocumentationSection(ctx context.Context, sectionID uuid.UUID, userID uuid.UUID) (*DocumentationSection, error) {
	var section DocumentationSection
	if err := r.db.WithContext(ctx).
		Joins("JOIN project_documentations ON project_documentations.id = documentation_sections.documentation_id").
		Joins("JOIN projects ON projects.id = project_documentations.project_id").
		Where("documentation_sections.id = ? AND projects.user_id = ?", sectionID, userID).
		First(&section).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrSectionNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &section, nil
}

func (r *gormRepository) BulkUpdateDocumentationSections(ctx context.Context, sections []*DocumentationSection) error {
	if len(sections) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, section := range sections {
			if err := section.Validate(); err != nil {
				return err
			}
			if err := tx.Save(section).Error; err != nil {
				return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
			}
		}
		return nil
	})
}

// Technology operations

func (r *gormRepository) CreateTechnology(ctx context.Context, tech *ProjectTechnology) error {
	if err := tech.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(tech).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateTechnology(ctx context.Context, tech *ProjectTechnology) error {
	if err := tech.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Save(tech).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) DeleteTechnology(ctx context.Context, techID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var tech ProjectTechnology
		if err := tx.Joins("JOIN projects ON projects.id = project_technologies.project_id").
			Where("project_technologies.id = ? AND projects.user_id = ?", techID, userID).
			First(&tech).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NewDomainError(ErrCodeNotFound, ErrTechnologyNotFound)
			}
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
		}

		if err := tx.Delete(&tech).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}
		return nil
	})
}

func (r *gormRepository) ListTechnologies(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectTechnology, error) {
	var technologies []ProjectTechnology
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = project_technologies.project_id").
		Where("project_technologies.project_id = ? AND projects.user_id = ?", projectID, userID).
		Order("project_technologies.category asc, project_technologies.name asc").
		Find(&technologies).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return technologies, nil
}

func (r *gormRepository) GetTechnology(ctx context.Context, techID uuid.UUID, userID uuid.UUID) (*ProjectTechnology, error) {
	var tech ProjectTechnology
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = project_technologies.project_id").
		Where("project_technologies.id = ? AND projects.user_id = ?", techID, userID).
		First(&tech).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrTechnologyNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &tech, nil
}

func (r *gormRepository) BulkCreateTechnologies(ctx context.Context, technologies []*ProjectTechnology) error {
	if len(technologies) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, tech := range technologies {
			if err := tech.Validate(); err != nil {
				return err
			}
		}
		if err := tx.Create(&technologies).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
		}
		return nil
	})
}

func (r *gormRepository) BulkUpdateTechnologies(ctx context.Context, technologies []*ProjectTechnology) error {
	if len(technologies) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, tech := range technologies {
			if err := tech.Validate(); err != nil {
				return err
			}
			if err := tx.Save(tech).Error; err != nil {
				return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
			}
		}
		return nil
	})
}

// File Structure operations

func (r *gormRepository) CreateFileStructure(ctx context.Context, fs *ProjectFileStructure) error {
	if err := fs.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(fs).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateFileStructure(ctx context.Context, fs *ProjectFileStructure) error {
	if err := fs.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Save(fs).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) DeleteFileStructure(ctx context.Context, fsID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var fs ProjectFileStructure
		if err := tx.Joins("JOIN projects ON projects.id = project_file_structures.project_id").
			Where("project_file_structures.id = ? AND projects.user_id = ?", fsID, userID).
			First(&fs).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NewDomainError(ErrCodeNotFound, ErrFileStructureNotFound)
			}
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
		}

		if err := tx.Where("parent_id = ?", fsID).Delete(&ProjectFileStructure{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}

		if err := tx.Delete(&fs).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}
		return nil
	})
}

func (r *gormRepository) ListFileStructures(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectFileStructure, error) {
	var structures []ProjectFileStructure
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = project_file_structures.project_id").
		Where("project_file_structures.project_id = ? AND projects.user_id = ?", projectID, userID).
		Order("project_file_structures.position asc, project_file_structures.path asc").
		Find(&structures).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return structures, nil
}

func (r *gormRepository) GetFileStructure(ctx context.Context, fsID uuid.UUID, userID uuid.UUID) (*ProjectFileStructure, error) {
	var fs ProjectFileStructure
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = project_file_structures.project_id").
		Where("project_file_structures.id = ? AND projects.user_id = ?", fsID, userID).
		First(&fs).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrFileStructureNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &fs, nil
}

func (r *gormRepository) BulkCreateFileStructures(ctx context.Context, structures []*ProjectFileStructure) error {
	if len(structures) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, fs := range structures {
			if err := fs.Validate(); err != nil {
				return err
			}
		}
		if err := tx.Create(&structures).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
		}
		return nil
	})
}

func (r *gormRepository) BulkUpdateFileStructures(ctx context.Context, structures []*ProjectFileStructure) error {
	if len(structures) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, fs := range structures {
			if err := fs.Validate(); err != nil {
				return err
			}
			if err := tx.Save(fs).Error; err != nil {
				return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
			}
		}
		return nil
	})
}

func (r *gormRepository) DeleteFileStructuresByProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN projects ON projects.id = project_file_structures.project_id").
			Where("project_file_structures.project_id = ? AND projects.user_id = ?", projectID, userID).
			Delete(&ProjectFileStructure{}).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}
		return nil
	})
}

// Architecture Diagram operations

func (r *gormRepository) CreateArchitectureDiagram(ctx context.Context, diagram *ProjectArchitectureDiagram) error {
	if err := diagram.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(diagram).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateArchitectureDiagram(ctx context.Context, diagram *ProjectArchitectureDiagram) error {
	if err := diagram.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Save(diagram).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) DeleteArchitectureDiagram(ctx context.Context, diagramID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var diagram ProjectArchitectureDiagram
		if err := tx.Joins("JOIN projects ON projects.id = project_architecture_diagrams.project_id").
			Where("project_architecture_diagrams.id = ? AND projects.user_id = ?", diagramID, userID).
			First(&diagram).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NewDomainError(ErrCodeNotFound, ErrDiagramNotFound)
			}
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
		}

		if err := tx.Delete(&diagram).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}
		return nil
	})
}

func (r *gormRepository) ListArchitectureDiagrams(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectArchitectureDiagram, error) {
	var diagrams []ProjectArchitectureDiagram
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = project_architecture_diagrams.project_id").
		Where("project_architecture_diagrams.project_id = ? AND projects.user_id = ?", projectID, userID).
		Order("project_architecture_diagrams.created_at desc").
		Find(&diagrams).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return diagrams, nil
}

func (r *gormRepository) GetArchitectureDiagram(ctx context.Context, diagramID uuid.UUID, userID uuid.UUID) (*ProjectArchitectureDiagram, error) {
	var diagram ProjectArchitectureDiagram
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = project_architecture_diagrams.project_id").
		Where("project_architecture_diagrams.id = ? AND projects.user_id = ?", diagramID, userID).
		First(&diagram).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrDiagramNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &diagram, nil
}
