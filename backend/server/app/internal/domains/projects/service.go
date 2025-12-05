package projects

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates project workflows.
type Service interface {
	CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error)
	UpdateProjectStatus(ctx context.Context, req UpdateStatusRequest) (*Project, error)
	UpdateProjectMetrics(ctx context.Context, req UpdateMetricsRequest) (*Project, error)
	DeleteProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error
	ListProjects(ctx context.Context, userID uuid.UUID) ([]Project, error)
	GetProjectBySlug(ctx context.Context, userID uuid.UUID, slug string) (*Project, error)
	SearchProjectsBySlug(ctx context.Context, userID uuid.UUID, slug string) ([]Project, error)

	AddMilestone(ctx context.Context, req AddMilestoneRequest) (*Milestone, error)
	ToggleMilestone(ctx context.Context, req ToggleMilestoneRequest) (*Milestone, error)
	ListMilestones(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]Milestone, error)
	BulkUpdateMilestones(ctx context.Context, req BulkUpdateMilestonesRequest) ([]*Milestone, error)

	CreateKanbanColumn(ctx context.Context, req CreateKanbanColumnRequest) (KanbanBoard, error)
	UpdateKanbanColumn(ctx context.Context, req UpdateKanbanColumnRequest) (KanbanBoard, error)
	ReorderKanbanColumns(ctx context.Context, req ReorderKanbanColumnsRequest) (KanbanBoard, error)
	DeleteKanbanColumn(ctx context.Context, req DeleteKanbanColumnRequest) (KanbanBoard, error)

	CreateKanbanCard(ctx context.Context, req CreateKanbanCardRequest) (KanbanBoard, error)
	UpdateKanbanCard(ctx context.Context, req UpdateKanbanCardRequest) (KanbanBoard, error)
	MoveKanbanCard(ctx context.Context, req MoveKanbanCardRequest) (KanbanBoard, error)
	DeleteKanbanCard(ctx context.Context, req DeleteKanbanCardRequest) (KanbanBoard, error)
	GetKanbanBoard(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (KanbanBoard, error)

	CreateDependency(ctx context.Context, req CreateDependencyRequest) (*ProjectDependency, error)
	DeleteDependency(ctx context.Context, req DeleteDependencyRequest) error
	ListDependencies(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectDependency, error)

	DuplicateProject(ctx context.Context, req DuplicateProjectRequest) (*Project, error)

	// Documentation operations
	CreateDocumentation(ctx context.Context, req CreateDocumentationRequest) (*ProjectDocumentation, error)
	UpdateDocumentationVisibility(ctx context.Context, req UpdateDocumentationVisibilityRequest) (*ProjectDocumentation, error)
	GetDocumentation(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*ProjectDocumentation, error)
	GetPublicDocumentation(ctx context.Context, projectSlug string) (*ProjectDocumentation, error)
	DeleteDocumentation(ctx context.Context, req DeleteDocumentationRequest) error

	// Documentation Section operations
	CreateDocumentationSection(ctx context.Context, req CreateDocumentationSectionRequest) (*DocumentationSection, error)
	UpdateDocumentationSection(ctx context.Context, req UpdateDocumentationSectionRequest) (*DocumentationSection, error)
	DeleteDocumentationSection(ctx context.Context, req DeleteDocumentationSectionRequest) error
	ListDocumentationSections(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]DocumentationSection, error)
	ReorderDocumentationSections(ctx context.Context, req ReorderDocumentationSectionsRequest) ([]DocumentationSection, error)

	// Technology operations
	CreateTechnology(ctx context.Context, req CreateTechnologyRequest) (*ProjectTechnology, error)
	UpdateTechnology(ctx context.Context, req UpdateTechnologyRequest) (*ProjectTechnology, error)
	DeleteTechnology(ctx context.Context, req DeleteTechnologyRequest) error
	ListTechnologies(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectTechnology, error)
	BulkCreateTechnologies(ctx context.Context, req BulkCreateTechnologiesRequest) ([]ProjectTechnology, error)
	BulkUpdateTechnologies(ctx context.Context, req BulkUpdateTechnologiesRequest) ([]ProjectTechnology, error)

	// File Structure operations
	CreateFileStructure(ctx context.Context, req CreateFileStructureRequest) (*ProjectFileStructure, error)
	UpdateFileStructure(ctx context.Context, req UpdateFileStructureRequest) (*ProjectFileStructure, error)
	DeleteFileStructure(ctx context.Context, req DeleteFileStructureRequest) error
	ListFileStructures(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectFileStructure, error)
	BulkCreateFileStructures(ctx context.Context, req BulkCreateFileStructuresRequest) ([]ProjectFileStructure, error)
	BulkUpdateFileStructures(ctx context.Context, req BulkUpdateFileStructuresRequest) ([]ProjectFileStructure, error)

	// Architecture Diagram operations
	CreateArchitectureDiagram(ctx context.Context, req CreateArchitectureDiagramRequest) (*ProjectArchitectureDiagram, error)
	UpdateArchitectureDiagram(ctx context.Context, req UpdateArchitectureDiagramRequest) (*ProjectArchitectureDiagram, error)
	DeleteArchitectureDiagram(ctx context.Context, req DeleteArchitectureDiagramRequest) error
	ListArchitectureDiagrams(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectArchitectureDiagram, error)
	GetArchitectureDiagram(ctx context.Context, diagramID uuid.UUID, userID uuid.UUID) (*ProjectArchitectureDiagram, error)
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

type CreateProjectRequest struct {
	UserID      uuid.UUID
	Name        string
	Description string
	Status      ProjectStatus
	HealthScore int
	MRR         float64
	CAC         float64
	LTV         float64
	ChurnRate   float64
}

type UpdateStatusRequest struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Status    ProjectStatus
}

type UpdateMetricsRequest struct {
	ProjectID   uuid.UUID
	UserID      uuid.UUID
	HealthScore int
	MRR         float64
	CAC         float64
	LTV         float64
	ChurnRate   float64
}

type AddMilestoneRequest struct {
	ProjectID   uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	DueDate     time.Time
}

type ToggleMilestoneRequest struct {
	MilestoneID uuid.UUID
	UserID      uuid.UUID
	Completed   bool
}

type MilestoneUpdate struct {
	MilestoneID uuid.UUID
	Title       *string
	Description *string
	DueDate     *time.Time
	Completed   *bool
}

type BulkUpdateMilestonesRequest struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Updates   []MilestoneUpdate
}

type CreateKanbanColumnRequest struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Name      string
	WIPLimit  *int
	Position  *int
}

type UpdateKanbanColumnRequest struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	ColumnID  uuid.UUID
	Name      *string
	WIPLimit  *int
}

type ReorderKanbanColumnsRequest struct {
	ProjectID   uuid.UUID
	UserID      uuid.UUID
	ColumnOrder []uuid.UUID
}

type DeleteKanbanColumnRequest struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	ColumnID  uuid.UUID
}

type CreateKanbanCardRequest struct {
	ProjectID   uuid.UUID
	UserID      uuid.UUID
	ColumnID    uuid.UUID
	Title       string
	Description string
	DueDate     *time.Time
	MilestoneID *uuid.UUID
	Position    *int
}

type UpdateKanbanCardRequest struct {
	ProjectID   uuid.UUID
	UserID      uuid.UUID
	CardID      uuid.UUID
	Title       *string
	Description *string
	DueDate     *time.Time
	MilestoneID *uuid.UUID
}

type MoveKanbanCardRequest struct {
	ProjectID      uuid.UUID
	UserID         uuid.UUID
	CardID         uuid.UUID
	TargetColumnID uuid.UUID
	TargetPosition int
}

type DeleteKanbanCardRequest struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	CardID    uuid.UUID
}

type CreateDependencyRequest struct {
	ProjectID          uuid.UUID
	UserID             uuid.UUID
	DependsOnProjectID uuid.UUID
	Type               DependencyType
}

type DeleteDependencyRequest struct {
	DependencyID uuid.UUID
	UserID       uuid.UUID
}

type DuplicateProjectRequest struct {
	TemplateProjectID uuid.UUID
	UserID            uuid.UUID
	Name              string
	Description       string
	Status            *ProjectStatus
	HealthScore       *int
	MRR               *float64
	CAC               *float64
	LTV               *float64
	ChurnRate         *float64
	CopyBoard         bool
	CopyMilestones    bool
	CopyDependencies  bool
}

// Documentation request types

type CreateDocumentationRequest struct {
	ProjectID  uuid.UUID
	UserID     uuid.UUID
	Visibility DocumentationVisibility
}

type UpdateDocumentationVisibilityRequest struct {
	ProjectID  uuid.UUID
	UserID     uuid.UUID
	Visibility DocumentationVisibility
}

type DeleteDocumentationRequest struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
}

type CreateDocumentationSectionRequest struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Type      DocumentationSectionType
	Title     string
	Content   string
	Position  *int
}

type UpdateDocumentationSectionRequest struct {
	SectionID uuid.UUID
	UserID    uuid.UUID
	Title     *string
	Content   *string
	Position  *int
}

type DeleteDocumentationSectionRequest struct {
	SectionID uuid.UUID
	UserID    uuid.UUID
}

type ReorderDocumentationSectionsRequest struct {
	ProjectID    uuid.UUID
	UserID       uuid.UUID
	SectionOrder []uuid.UUID
}

// Technology request types

type CreateTechnologyRequest struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Name      string
	Version   string
	Category  TechnologyCategory
	Purpose   string
	Link      string
}

type UpdateTechnologyRequest struct {
	TechID   uuid.UUID
	UserID   uuid.UUID
	Name     *string
	Version  *string
	Category *TechnologyCategory
	Purpose  *string
	Link     *string
}

type DeleteTechnologyRequest struct {
	TechID uuid.UUID
	UserID uuid.UUID
}

type BulkCreateTechnologiesRequest struct {
	ProjectID    uuid.UUID
	UserID       uuid.UUID
	Technologies []CreateTechnologyRequest
}

type BulkUpdateTechnologiesRequest struct {
	ProjectID    uuid.UUID
	UserID       uuid.UUID
	Technologies []UpdateTechnologyRequest
}

// File Structure request types

type CreateFileStructureRequest struct {
	ProjectID   uuid.UUID
	UserID      uuid.UUID
	Path        string
	Name        string
	IsDirectory bool
	ParentID    *uuid.UUID
	Language    string
	LineCount   int
	Purpose     string
	Position    *int
}

type UpdateFileStructureRequest struct {
	FileStructureID uuid.UUID
	UserID          uuid.UUID
	Purpose         *string
	LineCount       *int
	Language        *string
	Position        *int
}

type DeleteFileStructureRequest struct {
	FileStructureID uuid.UUID
	UserID          uuid.UUID
}

type BulkCreateFileStructuresRequest struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Structures []CreateFileStructureRequest
}

type BulkUpdateFileStructuresRequest struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Structures []UpdateFileStructureRequest
}

// Architecture Diagram request types

type CreateArchitectureDiagramRequest struct {
	ProjectID   uuid.UUID
	UserID      uuid.UUID
	Type        ArchitectureDiagramType
	Title       string
	Description string
	Content     string
	Format      string
	ImageURL    string
}

type UpdateArchitectureDiagramRequest struct {
	DiagramID   uuid.UUID
	UserID      uuid.UUID
	Title       *string
	Description *string
	Content     *string
	ImageURL    *string
}

type DeleteArchitectureDiagramRequest struct {
	DiagramID uuid.UUID
	UserID   uuid.UUID
}

type KanbanColumnWithCards struct {
	Column KanbanColumn
	Cards  []KanbanCard
}

type KanbanBoard struct {
	ProjectID uuid.UUID
	Columns   []KanbanColumnWithCards
}

// Project CRUD

func (s *service) CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error) {
	if req.Status == "" {
		req.Status = ProjectStatusIdea
	}

	project, err := NewProject(req.UserID, req.Name, req.Description, req.Status, req.HealthScore, req.MRR, req.CAC, req.LTV, req.ChurnRate)
	if err != nil {
		return nil, err
	}

	if err := s.assignProjectSlug(ctx, project); err != nil {
		return nil, err
	}

	if err := s.repo.CreateProject(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *service) UpdateProjectStatus(ctx context.Context, req UpdateStatusRequest) (*Project, error) {
	project, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := project.UpdateStatus(req.Status); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateProject(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *service) UpdateProjectMetrics(ctx context.Context, req UpdateMetricsRequest) (*Project, error) {
	project, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := project.UpdateMetrics(req.HealthScore, req.MRR, req.CAC, req.LTV, req.ChurnRate); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateProject(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *service) DeleteProject(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	// Verify project exists and belongs to user
	if _, err := s.repo.GetProject(ctx, projectID, userID); err != nil {
		return err
	}

	return s.repo.DeleteProject(ctx, projectID, userID)
}

func (s *service) ListProjects(ctx context.Context, userID uuid.UUID) ([]Project, error) {
	return s.repo.ListProjects(ctx, userID)
}

func (s *service) GetProjectBySlug(ctx context.Context, userID uuid.UUID, slug string) (*Project, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectSlug)
	}
	return s.repo.GetProjectBySlug(ctx, slug, userID)
}

func (s *service) SearchProjectsBySlug(ctx context.Context, userID uuid.UUID, slug string) ([]Project, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return []Project{}, nil
	}
	return s.repo.SearchProjectsBySlug(ctx, slug, userID)
}

func (s *service) assignProjectSlug(ctx context.Context, project *Project) error {
	if project == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilProject)
	}

	base := strings.TrimSpace(project.Slug)
	if base == "" {
		base = generateProjectSlug(project.Name)
	}

	slug := base
	for attempt := 0; attempt < 50; attempt++ {
		taken, err := s.repo.IsProjectSlugTaken(ctx, project.UserID, slug, project.ID)
		if err != nil {
			return err
		}
		if !taken {
			project.Slug = slug
			return nil
		}
		slug = fmt.Sprintf("%s-%d", base, attempt+2)
	}

	return NewDomainError(ErrCodeRepositoryFailure, "projects: unable to generate unique slug")
}

// Milestones

func (s *service) AddMilestone(ctx context.Context, req AddMilestoneRequest) (*Milestone, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return nil, err
	}

	milestone, err := NewMilestone(req.ProjectID, req.Title, req.Description, req.DueDate)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateMilestone(ctx, milestone); err != nil {
		return nil, err
	}

	return milestone, nil
}

func (s *service) ToggleMilestone(ctx context.Context, req ToggleMilestoneRequest) (*Milestone, error) {
	milestone, err := s.repo.GetMilestone(ctx, req.MilestoneID, req.UserID)
	if err != nil {
		return nil, err
	}

	milestone.MarkCompleted(req.Completed)

	if err := s.repo.UpdateMilestone(ctx, milestone); err != nil {
		return nil, err
	}

	return milestone, nil
}

func (s *service) ListMilestones(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]Milestone, error) {
	return s.repo.ListMilestones(ctx, projectID, userID)
}

func (s *service) BulkUpdateMilestones(ctx context.Context, req BulkUpdateMilestonesRequest) ([]*Milestone, error) {
	if len(req.Updates) == 0 {
		return []*Milestone{}, nil
	}

	updated := make([]*Milestone, 0, len(req.Updates))
	for _, payload := range req.Updates {
		milestone, err := s.repo.GetMilestone(ctx, payload.MilestoneID, req.UserID)
		if err != nil {
			return nil, err
		}

		if payload.Title != nil {
			milestone.Title = strings.TrimSpace(*payload.Title)
		}
		if payload.Description != nil {
			milestone.Description = strings.TrimSpace(*payload.Description)
		}
		if payload.DueDate != nil {
			d := payload.DueDate.UTC()
			milestone.DueDate = d
		}
		if payload.Completed != nil {
			milestone.MarkCompleted(*payload.Completed)
		}
		milestone.UpdatedAt = time.Now().UTC()

		if err := milestone.Validate(); err != nil {
			return nil, err
		}

		updated = append(updated, milestone)
	}

	if err := s.repo.BulkUpdateMilestones(ctx, updated); err != nil {
		return nil, err
	}

	return updated, nil
}

// Kanban board helpers

func (s *service) CreateKanbanColumn(ctx context.Context, req CreateKanbanColumnRequest) (KanbanBoard, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return KanbanBoard{}, err
	}

	existingColumns, err := s.repo.ListKanbanColumns(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return KanbanBoard{}, err
	}

	position := len(existingColumns)
	if req.Position != nil {
		if *req.Position < 0 || *req.Position > len(existingColumns) {
			return KanbanBoard{}, NewDomainError(ErrCodeInvalidPayload, ErrInvalidKanbanPosition)
		}
		position = *req.Position
	}

	limit := 0
	if req.WIPLimit != nil {
		if *req.WIPLimit < 0 {
			return KanbanBoard{}, NewDomainError(ErrCodeInvalidPayload, ErrInvalidWIPLimit)
		}
		limit = *req.WIPLimit
	}

	column, err := NewKanbanColumn(req.ProjectID, req.Name, position, limit)
	if err != nil {
		return KanbanBoard{}, err
	}

	if err := s.repo.CreateKanbanColumn(ctx, column); err != nil {
		return KanbanBoard{}, err
	}

	if err := s.repositionColumn(ctx, req.ProjectID, req.UserID, column.ID, position); err != nil {
		return KanbanBoard{}, err
	}

	return s.buildKanbanBoard(ctx, req.ProjectID, req.UserID)
}

func (s *service) UpdateKanbanColumn(ctx context.Context, req UpdateKanbanColumnRequest) (KanbanBoard, error) {
	column, err := s.repo.GetKanbanColumn(ctx, req.ColumnID, req.UserID)
	if err != nil {
		return KanbanBoard{}, err
	}

	if req.Name != nil {
		if err := column.Rename(*req.Name); err != nil {
			return KanbanBoard{}, err
		}
	}
	if req.WIPLimit != nil {
		if err := column.SetWIPLimit(*req.WIPLimit); err != nil {
			return KanbanBoard{}, err
		}
	}

	if err := s.repo.UpdateKanbanColumn(ctx, column); err != nil {
		return KanbanBoard{}, err
	}

	return s.buildKanbanBoard(ctx, req.ProjectID, req.UserID)
}

func (s *service) ReorderKanbanColumns(ctx context.Context, req ReorderKanbanColumnsRequest) (KanbanBoard, error) {
	columns, err := s.repo.ListKanbanColumns(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return KanbanBoard{}, err
	}

	if len(columns) != len(req.ColumnOrder) {
		return KanbanBoard{}, NewDomainError(ErrCodeInvalidPayload, ErrInvalidColumnOrder)
	}

	lookup := make(map[uuid.UUID]KanbanColumn, len(columns))
	for _, column := range columns {
		lookup[column.ID] = column
	}

	ordered := make([]KanbanColumn, 0, len(columns))
	for _, id := range req.ColumnOrder {
		column, ok := lookup[id]
		if !ok {
			return KanbanBoard{}, NewDomainError(ErrCodeInvalidPayload, ErrInvalidColumnOrder)
		}
		ordered = append(ordered, column)
	}

	if err := s.normalizeColumnPositions(ctx, req.ProjectID, req.UserID, ordered); err != nil {
		return KanbanBoard{}, err
	}

	return s.buildKanbanBoard(ctx, req.ProjectID, req.UserID)
}

func (s *service) DeleteKanbanColumn(ctx context.Context, req DeleteKanbanColumnRequest) (KanbanBoard, error) {
	if err := s.repo.DeleteKanbanColumn(ctx, req.ColumnID, req.UserID); err != nil {
		return KanbanBoard{}, err
	}

	columns, err := s.repo.ListKanbanColumns(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return KanbanBoard{}, err
	}

	if err := s.normalizeColumnPositions(ctx, req.ProjectID, req.UserID, columns); err != nil {
		return KanbanBoard{}, err
	}

	return s.buildKanbanBoard(ctx, req.ProjectID, req.UserID)
}

func (s *service) CreateKanbanCard(ctx context.Context, req CreateKanbanCardRequest) (KanbanBoard, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return KanbanBoard{}, err
	}

	column, err := s.repo.GetKanbanColumn(ctx, req.ColumnID, req.UserID)
	if err != nil {
		return KanbanBoard{}, err
	}
	if column.ProjectID != req.ProjectID {
		return KanbanBoard{}, NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}

	allCards, err := s.repo.ListKanbanCards(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return KanbanBoard{}, err
	}

	columnCards := filterCardsByColumn(allCards, req.ColumnID)
	if column.WIPLimit > 0 && len(columnCards) >= column.WIPLimit {
		return KanbanBoard{}, NewDomainError(ErrCodeConflict, ErrWIPLimitExceeded)
	}

	position := len(columnCards)
	if req.Position != nil {
		if *req.Position < 0 || *req.Position > len(columnCards) {
			return KanbanBoard{}, NewDomainError(ErrCodeInvalidPayload, ErrInvalidKanbanPosition)
		}
		position = *req.Position
	}

	var due *time.Time
	if req.DueDate != nil {
		d := req.DueDate.UTC()
		due = &d
	}

	var milestoneID *uuid.UUID
	if req.MilestoneID != nil {
		id := *req.MilestoneID
		milestoneID = &id
	}

	card, err := NewKanbanCard(req.ProjectID, req.ColumnID, req.Title, req.Description, position, due, milestoneID)
	if err != nil {
		return KanbanBoard{}, err
	}

	if err := s.repo.CreateKanbanCard(ctx, card); err != nil {
		return KanbanBoard{}, err
	}

	pos := position
	cardID := card.ID
	if err := s.repositionCards(ctx, req.ProjectID, req.UserID, req.ColumnID, &cardID, &pos); err != nil {
		return KanbanBoard{}, err
	}

	return s.buildKanbanBoard(ctx, req.ProjectID, req.UserID)
}

func (s *service) UpdateKanbanCard(ctx context.Context, req UpdateKanbanCardRequest) (KanbanBoard, error) {
	card, err := s.repo.GetKanbanCard(ctx, req.CardID, req.UserID)
	if err != nil {
		return KanbanBoard{}, err
	}

	title := card.Title
	if req.Title != nil {
		title = strings.TrimSpace(*req.Title)
	}

	description := card.Description
	if req.Description != nil {
		description = strings.TrimSpace(*req.Description)
	}

	var due *time.Time
	if req.DueDate != nil {
		d := req.DueDate.UTC()
		due = &d
	} else if card.DueDate != nil {
		d := card.DueDate.UTC()
		due = &d
	}

	var milestoneID *uuid.UUID
	if req.MilestoneID != nil {
		id := *req.MilestoneID
		milestoneID = &id
	} else if card.MilestoneID != nil {
		id := *card.MilestoneID
		milestoneID = &id
	}

	if err := card.UpdateDetails(title, description, due, milestoneID); err != nil {
		return KanbanBoard{}, err
	}

	if err := s.repo.UpdateKanbanCard(ctx, card); err != nil {
		return KanbanBoard{}, err
	}

	return s.buildKanbanBoard(ctx, req.ProjectID, req.UserID)
}

func (s *service) MoveKanbanCard(ctx context.Context, req MoveKanbanCardRequest) (KanbanBoard, error) {
	card, err := s.repo.GetKanbanCard(ctx, req.CardID, req.UserID)
	if err != nil {
		return KanbanBoard{}, err
	}

	sourceColumnID := card.ColumnID

	targetColumn, err := s.repo.GetKanbanColumn(ctx, req.TargetColumnID, req.UserID)
	if err != nil {
		return KanbanBoard{}, err
	}

	allCards, err := s.repo.ListKanbanCards(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return KanbanBoard{}, err
	}

	targetCards := filterCardsByColumn(allCards, req.TargetColumnID)
	if targetColumn.WIPLimit > 0 {
		cardCount := len(targetCards)
		if sourceColumnID != req.TargetColumnID {
			cardCount++
		}
		if cardCount > targetColumn.WIPLimit {
			return KanbanBoard{}, NewDomainError(ErrCodeConflict, ErrWIPLimitExceeded)
		}
	}

	if err := card.MoveToColumn(req.TargetColumnID, req.TargetPosition); err != nil {
		return KanbanBoard{}, err
	}

	if err := s.repo.UpdateKanbanCard(ctx, card); err != nil {
		return KanbanBoard{}, err
	}

	if sourceColumnID != req.TargetColumnID {
		if err := s.repositionCards(ctx, req.ProjectID, req.UserID, sourceColumnID, nil, nil); err != nil {
			return KanbanBoard{}, err
		}
	}

	pos := req.TargetPosition
	cardID := card.ID
	if err := s.repositionCards(ctx, req.ProjectID, req.UserID, req.TargetColumnID, &cardID, &pos); err != nil {
		return KanbanBoard{}, err
	}

	return s.buildKanbanBoard(ctx, req.ProjectID, req.UserID)
}

func (s *service) DeleteKanbanCard(ctx context.Context, req DeleteKanbanCardRequest) (KanbanBoard, error) {
	card, err := s.repo.GetKanbanCard(ctx, req.CardID, req.UserID)
	if err != nil {
		return KanbanBoard{}, err
	}

	if err := s.repo.DeleteKanbanCard(ctx, req.CardID, req.UserID); err != nil {
		return KanbanBoard{}, err
	}

	if err := s.repositionCards(ctx, req.ProjectID, req.UserID, card.ColumnID, nil, nil); err != nil {
		return KanbanBoard{}, err
	}

	return s.buildKanbanBoard(ctx, req.ProjectID, req.UserID)
}

func (s *service) GetKanbanBoard(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (KanbanBoard, error) {
	return s.buildKanbanBoard(ctx, projectID, userID)
}

// Dependencies

func (s *service) CreateDependency(ctx context.Context, req CreateDependencyRequest) (*ProjectDependency, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetProject(ctx, req.DependsOnProjectID, req.UserID); err != nil {
		return nil, err
	}

	if exists, err := s.repo.DependencyExists(ctx, req.DependsOnProjectID, req.ProjectID); err == nil && exists {
		return nil, NewDomainError(ErrCodeConflict, ErrDependencyAlreadyExists)
	}

	dependency, err := NewProjectDependency(req.ProjectID, req.DependsOnProjectID, req.Type)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateDependency(ctx, dependency); err != nil {
		return nil, err
	}

	return dependency, nil
}

func (s *service) DeleteDependency(ctx context.Context, req DeleteDependencyRequest) error {
	return s.repo.DeleteDependency(ctx, req.DependencyID, req.UserID)
}

func (s *service) ListDependencies(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectDependency, error) {
	return s.repo.ListDependencies(ctx, projectID, userID)
}

// Duplication

func (s *service) DuplicateProject(ctx context.Context, req DuplicateProjectRequest) (*Project, error) {
	template, err := s.repo.GetProject(ctx, req.TemplateProjectID, req.UserID)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = template.Name + " Copy"
	}
	description := template.Description
	if strings.TrimSpace(req.Description) != "" {
		description = strings.TrimSpace(req.Description)
	}

	status := template.Status
	if req.Status != nil {
		status = *req.Status
	}

	health := template.HealthScore
	if req.HealthScore != nil {
		health = *req.HealthScore
	}

	mrr := template.MRR
	if req.MRR != nil {
		mrr = *req.MRR
	}

	cac := template.CAC
	if req.CAC != nil {
		cac = *req.CAC
	}

	ltv := template.LTV
	if req.LTV != nil {
		ltv = *req.LTV
	}

	churn := template.ChurnRate
	if req.ChurnRate != nil {
		churn = *req.ChurnRate
	}

	project, err := NewProject(template.UserID, name, description, status, health, mrr, cac, ltv, churn)
	if err != nil {
		return nil, err
	}

	if err := s.assignProjectSlug(ctx, project); err != nil {
		return nil, err
	}

	var columns []*KanbanColumn
	var cards []*KanbanCard
	var milestones []*Milestone

	columnMap := make(map[uuid.UUID]uuid.UUID)
	milestoneMap := make(map[uuid.UUID]uuid.UUID)

	if req.CopyBoard {
		existingColumns, err := s.repo.ListKanbanColumns(ctx, template.ID, req.UserID)
		if err != nil {
			return nil, err
		}
		sort.Slice(existingColumns, func(i, j int) bool {
			return existingColumns[i].Position < existingColumns[j].Position
		})
		for _, col := range existingColumns {
			clone, err := NewKanbanColumn(project.ID, col.Name, col.Position, col.WIPLimit)
			if err != nil {
				return nil, err
			}
			clone.Position = col.Position
			columns = append(columns, clone)
			columnMap[col.ID] = clone.ID
		}
	}

	if req.CopyMilestones {
		existingMilestones, err := s.repo.ListMilestones(ctx, template.ID, req.UserID)
		if err != nil {
			return nil, err
		}
		for _, milestone := range existingMilestones {
			due := milestone.DueDate
			clone, err := NewMilestone(project.ID, milestone.Title, milestone.Description, due)
			if err != nil {
				return nil, err
			}
			if milestone.Completed {
				clone.MarkCompleted(true)
			}
			milestones = append(milestones, clone)
			milestoneMap[milestone.ID] = clone.ID
		}
	}

	if req.CopyBoard {
		existingCards, err := s.repo.ListKanbanCards(ctx, template.ID, req.UserID)
		if err != nil {
			return nil, err
		}
		sort.Slice(existingCards, func(i, j int) bool {
			if existingCards[i].ColumnID == existingCards[j].ColumnID {
				return existingCards[i].Position < existingCards[j].Position
			}
			return existingCards[i].ColumnID.String() < existingCards[j].ColumnID.String()
		})
		for _, card := range existingCards {
			newColumnID, ok := columnMap[card.ColumnID]
			if !ok {
				continue
			}

			var due *time.Time
			if card.DueDate != nil {
				d := card.DueDate.UTC()
				due = &d
			}

			var newMilestoneID *uuid.UUID
			if card.MilestoneID != nil {
				if mapped, ok := milestoneMap[*card.MilestoneID]; ok {
					id := mapped
					newMilestoneID = &id
				}
			}

			clone, err := NewKanbanCard(project.ID, newColumnID, card.Title, card.Description, card.Position, due, newMilestoneID)
			if err != nil {
				return nil, err
			}
			clone.Position = card.Position
			cards = append(cards, clone)
		}
	}

	if err := s.repo.CreateProjectWithRelated(ctx, project, columns, cards, milestones); err != nil {
		return nil, err
	}

	if req.CopyDependencies {
		deps, err := s.repo.ListDependencies(ctx, template.ID, req.UserID)
		if err != nil {
			return nil, err
		}
		for _, dep := range deps {
			if dep.DependsOnProjectID == template.ID {
				// Skip self-referential dependencies during duplication
				continue
			}
			if _, err := s.CreateDependency(ctx, CreateDependencyRequest{
				ProjectID:          project.ID,
				UserID:             req.UserID,
				DependsOnProjectID: dep.DependsOnProjectID,
				Type:               dep.Type,
			}); err != nil {
				if s.logger != nil {
					s.logger.Warn("projects: unable to clone dependency", slog.Any("error", err))
				}
			}
		}
	}

	return project, nil
}

// Helpers

func (s *service) normalizeColumnPositions(ctx context.Context, projectID, userID uuid.UUID, columns []KanbanColumn) error {
	for idx := range columns {
		if err := columns[idx].SetPosition(idx); err != nil {
			return err
		}
		if err := s.repo.UpdateKanbanColumn(ctx, &columns[idx]); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) normalizeCardPositions(ctx context.Context, userID uuid.UUID, cards []KanbanCard) error {
	for idx := range cards {
		if err := cards[idx].SetPosition(idx); err != nil {
			return err
		}
		if err := s.repo.UpdateKanbanCard(ctx, &cards[idx]); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) repositionColumn(ctx context.Context, projectID, userID, columnID uuid.UUID, desiredPosition int) error {
	columns, err := s.repo.ListKanbanColumns(ctx, projectID, userID)
	if err != nil {
		return err
	}

	var target *KanbanColumn
	others := make([]KanbanColumn, 0, len(columns))
	for _, column := range columns {
		if column.ID == columnID {
			copy := column
			target = &copy
			continue
		}
		others = append(others, column)
	}

	if target == nil {
		return NewDomainError(ErrCodeNotFound, ErrKanbanColumnNotFound)
	}

	if desiredPosition < 0 {
		desiredPosition = 0
	}
	if desiredPosition > len(others) {
		desiredPosition = len(others)
	}

	ordered := make([]KanbanColumn, 0, len(columns))
	ordered = append(ordered, others[:desiredPosition]...)
	ordered = append(ordered, *target)
	ordered = append(ordered, others[desiredPosition:]...)

	return s.normalizeColumnPositions(ctx, projectID, userID, ordered)
}

func (s *service) repositionCards(ctx context.Context, projectID, userID, columnID uuid.UUID, cardID *uuid.UUID, desiredPosition *int) error {
	cards, err := s.repo.ListKanbanCards(ctx, projectID, userID)
	if err != nil {
		return err
	}

	var target *KanbanCard
	others := make([]KanbanCard, 0)
	for _, card := range cards {
		if card.ColumnID != columnID {
			continue
		}
		if cardID != nil && card.ID == *cardID {
			copy := card
			target = &copy
			continue
		}
		others = append(others, card)
	}

	if cardID != nil && target == nil {
		return NewDomainError(ErrCodeNotFound, ErrKanbanCardNotFound)
	}

	if cardID == nil {
		return s.normalizeCardPositions(ctx, userID, others)
	}

	position := len(others)
	if desiredPosition != nil {
		position = *desiredPosition
		if position < 0 {
			position = 0
		}
		if position > len(others) {
			position = len(others)
		}
	}

	ordered := make([]KanbanCard, 0, len(others)+1)
	ordered = append(ordered, others[:position]...)
	ordered = append(ordered, *target)
	ordered = append(ordered, others[position:]...)

	return s.normalizeCardPositions(ctx, userID, ordered)
}

func filterCardsByColumn(cards []KanbanCard, columnID uuid.UUID) []KanbanCard {
	var filtered []KanbanCard
	for _, card := range cards {
		if card.ColumnID == columnID {
			filtered = append(filtered, card)
		}
	}
	return filtered
}

func (s *service) buildKanbanBoard(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (KanbanBoard, error) {
	columns, err := s.repo.ListKanbanColumns(ctx, projectID, userID)
	if err != nil {
		return KanbanBoard{}, err
	}

	cards, err := s.repo.ListKanbanCards(ctx, projectID, userID)
	if err != nil {
		return KanbanBoard{}, err
	}

	sort.Slice(columns, func(i, j int) bool {
		return columns[i].Position < columns[j].Position
	})

	cardMap := make(map[uuid.UUID][]KanbanCard)
	for _, card := range cards {
		cardMap[card.ColumnID] = append(cardMap[card.ColumnID], card)
	}

	board := KanbanBoard{ProjectID: projectID}
	for _, column := range columns {
		colCards := cardMap[column.ID]
		sort.Slice(colCards, func(i, j int) bool {
			if colCards[i].Position == colCards[j].Position {
				return colCards[i].CreatedAt.Before(colCards[j].CreatedAt)
			}
			return colCards[i].Position < colCards[j].Position
		})
		board.Columns = append(board.Columns, KanbanColumnWithCards{
			Column: column,
			Cards:  colCards,
		})
	}

	return board, nil
}

func insertColumn(columns []KanbanColumn, position int) []KanbanColumn {
	if position >= len(columns)-1 {
		return columns
	}

	inserted := make([]KanbanColumn, 0, len(columns))
	for idx, column := range columns {
		if idx == len(columns)-1 {
			continue
		}
		if idx == position {
			inserted = append(inserted, columns[len(columns)-1])
		}
		inserted = append(inserted, column)
	}
	return append(inserted, columns[len(columns)-1])
}

func insertCard(cards []KanbanCard, position int) []KanbanCard {
	if position >= len(cards)-1 {
		return cards
	}
	for i := len(cards) - 1; i > position; i-- {
		cards[i] = cards[i-1]
	}
	return cards
}

func removeCard(cards []KanbanCard, cardID uuid.UUID) []KanbanCard {
	filtered := cards[:0]
	for _, card := range cards {
		if card.ID != cardID {
			filtered = append(filtered, card)
		}
	}
	return filtered
}

func (b KanbanBoard) findColumnIndex(columnID uuid.UUID) int {
	for idx, column := range b.Columns {
		if column.Column.ID == columnID {
			return idx
		}
	}
	return -1
}

// Documentation operations

func (s *service) CreateDocumentation(ctx context.Context, req CreateDocumentationRequest) (*ProjectDocumentation, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return nil, err
	}

	visibility := req.Visibility
	if visibility == "" {
		visibility = VisibilityCollaborators
	}

	doc, err := NewProjectDocumentation(req.ProjectID, visibility)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateDocumentation(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *service) UpdateDocumentationVisibility(ctx context.Context, req UpdateDocumentationVisibilityRequest) (*ProjectDocumentation, error) {
	doc, err := s.repo.GetDocumentation(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := doc.UpdateVisibility(req.Visibility); err != nil {
		return nil, err
	}

	doc.IncrementVersion()

	if err := s.repo.UpdateDocumentation(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *service) GetDocumentation(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*ProjectDocumentation, error) {
	return s.repo.GetDocumentation(ctx, projectID, userID)
}

func (s *service) GetPublicDocumentation(ctx context.Context, projectSlug string) (*ProjectDocumentation, error) {
	project, err := s.repo.GetProjectBySlugPublic(ctx, projectSlug)
	if err != nil {
		return nil, err
	}

	doc, err := s.repo.GetDocumentationByProjectID(ctx, project.ID)
	if err != nil {
		return nil, err
	}

	if doc.Visibility != VisibilityPublic {
		return nil, NewDomainError(ErrCodeNotFound, ErrDocumentationNotFound)
	}

	return doc, nil
}

func (s *service) DeleteDocumentation(ctx context.Context, req DeleteDocumentationRequest) error {
	doc, err := s.repo.GetDocumentation(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return err
	}

	return s.repo.DeleteDocumentation(ctx, doc.ID, req.UserID)
}

// Documentation Section operations

func (s *service) CreateDocumentationSection(ctx context.Context, req CreateDocumentationSectionRequest) (*DocumentationSection, error) {
	doc, err := s.repo.GetDocumentation(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return nil, err
	}

	sections, err := s.repo.ListDocumentationSections(ctx, doc.ID, req.UserID)
	if err != nil {
		return nil, err
	}

	position := len(sections)
	if req.Position != nil {
		position = *req.Position
		if position < 0 {
			position = 0
		}
		if position > len(sections) {
			position = len(sections)
		}
	}

	section, err := NewDocumentationSection(doc.ID, req.Type, req.Title, req.Content, position)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateDocumentationSection(ctx, section); err != nil {
		return nil, err
	}

	doc.IncrementVersion()
	if err := s.repo.UpdateDocumentation(ctx, doc); err != nil {
		return nil, err
	}

	return section, nil
}

func (s *service) UpdateDocumentationSection(ctx context.Context, req UpdateDocumentationSectionRequest) (*DocumentationSection, error) {
	section, err := s.repo.GetDocumentationSection(ctx, req.SectionID, req.UserID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		if err := section.UpdateTitle(*req.Title); err != nil {
			return nil, err
		}
	}

	if req.Content != nil {
		section.UpdateContent(*req.Content)
	}

	if req.Position != nil {
		section.SetPosition(*req.Position)
	}

	if err := s.repo.UpdateDocumentationSection(ctx, section); err != nil {
		return nil, err
	}

	doc, err := s.repo.GetDocumentationByProjectID(ctx, section.DocumentationID)
	if err == nil {
		doc.IncrementVersion()
		s.repo.UpdateDocumentation(ctx, doc)
	}

	return section, nil
}

func (s *service) DeleteDocumentationSection(ctx context.Context, req DeleteDocumentationSectionRequest) error {
	section, err := s.repo.GetDocumentationSection(ctx, req.SectionID, req.UserID)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteDocumentationSection(ctx, req.SectionID, req.UserID); err != nil {
		return err
	}

	doc, err := s.repo.GetDocumentationByProjectID(ctx, section.DocumentationID)
	if err == nil {
		doc.IncrementVersion()
		s.repo.UpdateDocumentation(ctx, doc)
	}

	return nil
}

func (s *service) ListDocumentationSections(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]DocumentationSection, error) {
	doc, err := s.repo.GetDocumentation(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}

	return s.repo.ListDocumentationSections(ctx, doc.ID, userID)
}

func (s *service) ReorderDocumentationSections(ctx context.Context, req ReorderDocumentationSectionsRequest) ([]DocumentationSection, error) {
	doc, err := s.repo.GetDocumentation(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return nil, err
	}

	sections, err := s.repo.ListDocumentationSections(ctx, doc.ID, req.UserID)
	if err != nil {
		return nil, err
	}

	if len(sections) != len(req.SectionOrder) {
		return nil, NewDomainError(ErrCodeInvalidPayload, "projects: section order must include all sections")
	}

	lookup := make(map[uuid.UUID]DocumentationSection, len(sections))
	for _, section := range sections {
		lookup[section.ID] = section
	}

	ordered := make([]*DocumentationSection, 0, len(sections))
	for idx, id := range req.SectionOrder {
		section, ok := lookup[id]
		if !ok {
			return nil, NewDomainError(ErrCodeInvalidPayload, "projects: invalid section id in order")
		}
		section.SetPosition(idx)
		ordered = append(ordered, &section)
	}

	if err := s.repo.BulkUpdateDocumentationSections(ctx, ordered); err != nil {
		return nil, err
	}

	doc.IncrementVersion()
	if err := s.repo.UpdateDocumentation(ctx, doc); err != nil {
		return nil, err
	}

	result := make([]DocumentationSection, len(ordered))
	for i, s := range ordered {
		result[i] = *s
	}

	return result, nil
}

// Technology operations

func (s *service) CreateTechnology(ctx context.Context, req CreateTechnologyRequest) (*ProjectTechnology, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return nil, err
	}

	tech, err := NewProjectTechnology(req.ProjectID, req.Name, req.Version, req.Category, req.Purpose, req.Link)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateTechnology(ctx, tech); err != nil {
		return nil, err
	}

	return tech, nil
}

func (s *service) UpdateTechnology(ctx context.Context, req UpdateTechnologyRequest) (*ProjectTechnology, error) {
	tech, err := s.repo.GetTechnology(ctx, req.TechID, req.UserID)
	if err != nil {
		return nil, err
	}

	name := tech.Name
	if req.Name != nil {
		name = *req.Name
	}

	version := tech.Version
	if req.Version != nil {
		version = *req.Version
	}

	purpose := tech.Purpose
	if req.Purpose != nil {
		purpose = *req.Purpose
	}

	link := tech.Link
	if req.Link != nil {
		link = *req.Link
	}

	if err := tech.UpdateDetails(name, version, purpose, link); err != nil {
		return nil, err
	}

	if req.Category != nil {
		tech.Category = *req.Category
	}

	if err := s.repo.UpdateTechnology(ctx, tech); err != nil {
		return nil, err
	}

	return tech, nil
}

func (s *service) DeleteTechnology(ctx context.Context, req DeleteTechnologyRequest) error {
	return s.repo.DeleteTechnology(ctx, req.TechID, req.UserID)
}

func (s *service) ListTechnologies(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectTechnology, error) {
	return s.repo.ListTechnologies(ctx, projectID, userID)
}

func (s *service) BulkCreateTechnologies(ctx context.Context, req BulkCreateTechnologiesRequest) ([]ProjectTechnology, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return nil, err
	}

	technologies := make([]*ProjectTechnology, 0, len(req.Technologies))
	for _, techReq := range req.Technologies {
		tech, err := NewProjectTechnology(req.ProjectID, techReq.Name, techReq.Version, techReq.Category, techReq.Purpose, techReq.Link)
		if err != nil {
			return nil, err
		}
		technologies = append(technologies, tech)
	}

	if err := s.repo.BulkCreateTechnologies(ctx, technologies); err != nil {
		return nil, err
	}

	result := make([]ProjectTechnology, len(technologies))
	for i, t := range technologies {
		result[i] = *t
	}

	return result, nil
}

func (s *service) BulkUpdateTechnologies(ctx context.Context, req BulkUpdateTechnologiesRequest) ([]ProjectTechnology, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return nil, err
	}

	technologies := make([]*ProjectTechnology, 0, len(req.Technologies))
	for _, techReq := range req.Technologies {
		tech, err := s.repo.GetTechnology(ctx, techReq.TechID, req.UserID)
		if err != nil {
			return nil, err
		}

		if techReq.Name != nil {
			tech.Name = *techReq.Name
		}
		if techReq.Version != nil {
			tech.Version = *techReq.Version
		}
		if techReq.Category != nil {
			tech.Category = *techReq.Category
		}
		if techReq.Purpose != nil {
			tech.Purpose = *techReq.Purpose
		}
		if techReq.Link != nil {
			tech.Link = *techReq.Link
		}

		if err := tech.Validate(); err != nil {
			return nil, err
		}

		technologies = append(technologies, tech)
	}

	if err := s.repo.BulkUpdateTechnologies(ctx, technologies); err != nil {
		return nil, err
	}

	result := make([]ProjectTechnology, len(technologies))
	for i, t := range technologies {
		result[i] = *t
	}

	return result, nil
}

// File Structure operations

func (s *service) CreateFileStructure(ctx context.Context, req CreateFileStructureRequest) (*ProjectFileStructure, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return nil, err
	}

	position := 0
	if req.Position != nil {
		position = *req.Position
	}

	fs, err := NewProjectFileStructure(req.ProjectID, req.Path, req.Name, req.IsDirectory, req.ParentID, req.Language, req.LineCount, req.Purpose, position)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateFileStructure(ctx, fs); err != nil {
		return nil, err
	}

	return fs, nil
}

func (s *service) UpdateFileStructure(ctx context.Context, req UpdateFileStructureRequest) (*ProjectFileStructure, error) {
	fs, err := s.repo.GetFileStructure(ctx, req.FileStructureID, req.UserID)
	if err != nil {
		return nil, err
	}

	purpose := fs.Purpose
	if req.Purpose != nil {
		purpose = *req.Purpose
	}

	lineCount := fs.LineCount
	if req.LineCount != nil {
		lineCount = *req.LineCount
	}

	language := fs.Language
	if req.Language != nil {
		language = *req.Language
	}

	if err := fs.UpdateDetails(purpose, lineCount, language); err != nil {
		return nil, err
	}

	if req.Position != nil {
		fs.SetPosition(*req.Position)
	}

	if err := s.repo.UpdateFileStructure(ctx, fs); err != nil {
		return nil, err
	}

	return fs, nil
}

func (s *service) DeleteFileStructure(ctx context.Context, req DeleteFileStructureRequest) error {
	return s.repo.DeleteFileStructure(ctx, req.FileStructureID, req.UserID)
}

func (s *service) ListFileStructures(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectFileStructure, error) {
	return s.repo.ListFileStructures(ctx, projectID, userID)
}

func (s *service) BulkCreateFileStructures(ctx context.Context, req BulkCreateFileStructuresRequest) ([]ProjectFileStructure, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return nil, err
	}

	structures := make([]*ProjectFileStructure, 0, len(req.Structures))
	for _, fsReq := range req.Structures {
		position := 0
		if fsReq.Position != nil {
			position = *fsReq.Position
		}

		fs, err := NewProjectFileStructure(req.ProjectID, fsReq.Path, fsReq.Name, fsReq.IsDirectory, fsReq.ParentID, fsReq.Language, fsReq.LineCount, fsReq.Purpose, position)
		if err != nil {
			return nil, err
		}
		structures = append(structures, fs)
	}

	if err := s.repo.BulkCreateFileStructures(ctx, structures); err != nil {
		return nil, err
	}

	result := make([]ProjectFileStructure, len(structures))
	for i, s := range structures {
		result[i] = *s
	}

	return result, nil
}

func (s *service) BulkUpdateFileStructures(ctx context.Context, req BulkUpdateFileStructuresRequest) ([]ProjectFileStructure, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return nil, err
	}

	structures := make([]*ProjectFileStructure, 0, len(req.Structures))
	for _, fsReq := range req.Structures {
		fs, err := s.repo.GetFileStructure(ctx, fsReq.FileStructureID, req.UserID)
		if err != nil {
			return nil, err
		}

		if fsReq.Purpose != nil {
			fs.Purpose = *fsReq.Purpose
		}
		if fsReq.LineCount != nil {
			fs.LineCount = *fsReq.LineCount
		}
		if fsReq.Language != nil {
			fs.Language = *fsReq.Language
		}
		if fsReq.Position != nil {
			fs.SetPosition(*fsReq.Position)
		}

		if err := fs.Validate(); err != nil {
			return nil, err
		}

		structures = append(structures, fs)
	}

	if err := s.repo.BulkUpdateFileStructures(ctx, structures); err != nil {
		return nil, err
	}

	result := make([]ProjectFileStructure, len(structures))
	for i, s := range structures {
		result[i] = *s
	}

	return result, nil
}

// Architecture Diagram operations

func (s *service) CreateArchitectureDiagram(ctx context.Context, req CreateArchitectureDiagramRequest) (*ProjectArchitectureDiagram, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return nil, err
	}

	format := req.Format
	if format == "" {
		format = "mermaid"
	}

	diagram, err := NewProjectArchitectureDiagram(req.ProjectID, req.Type, req.Title, req.Description, req.Content, format, req.ImageURL)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateArchitectureDiagram(ctx, diagram); err != nil {
		return nil, err
	}

	return diagram, nil
}

func (s *service) UpdateArchitectureDiagram(ctx context.Context, req UpdateArchitectureDiagramRequest) (*ProjectArchitectureDiagram, error) {
	diagram, err := s.repo.GetArchitectureDiagram(ctx, req.DiagramID, req.UserID)
	if err != nil {
		return nil, err
	}

	title := diagram.Title
	if req.Title != nil {
		title = *req.Title
	}

	description := diagram.Description
	if req.Description != nil {
		description = *req.Description
	}

	imageURL := diagram.ImageURL
	if req.ImageURL != nil {
		imageURL = *req.ImageURL
	}

	if err := diagram.UpdateDetails(title, description, imageURL); err != nil {
		return nil, err
	}

	if req.Content != nil {
		diagram.UpdateContent(*req.Content)
	}

	if err := s.repo.UpdateArchitectureDiagram(ctx, diagram); err != nil {
		return nil, err
	}

	return diagram, nil
}

func (s *service) DeleteArchitectureDiagram(ctx context.Context, req DeleteArchitectureDiagramRequest) error {
	return s.repo.DeleteArchitectureDiagram(ctx, req.DiagramID, req.UserID)
}

func (s *service) ListArchitectureDiagrams(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectArchitectureDiagram, error) {
	return s.repo.ListArchitectureDiagrams(ctx, projectID, userID)
}

func (s *service) GetArchitectureDiagram(ctx context.Context, diagramID uuid.UUID, userID uuid.UUID) (*ProjectArchitectureDiagram, error) {
	return s.repo.GetArchitectureDiagram(ctx, diagramID, userID)
}
