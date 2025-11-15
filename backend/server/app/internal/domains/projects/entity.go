package projects

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ProjectStatus represents the lifecycle stage of a project.
type ProjectStatus string

const (
	ProjectStatusIdea       ProjectStatus = "idea"
	ProjectStatusPlanning   ProjectStatus = "planning"
	ProjectStatusExecuting  ProjectStatus = "executing"
	ProjectStatusMonitoring ProjectStatus = "monitoring"
	ProjectStatusCompleted  ProjectStatus = "completed"
)

// DependencyType describes the relationship between projects.
type DependencyType string

const (
	DependencyTypeBlocks   DependencyType = "blocks"
	DependencyTypeRelates  DependencyType = "relates"
	DependencyTypeSupports DependencyType = "supports"
)

// Project captures high-level roadmap metadata.
type Project struct {
	ID          uuid.UUID     `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID     `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Name        string        `gorm:"column:name;size:120;not null" json:"name"`
	Description string        `gorm:"column:description;size:255" json:"description"`
	Slug        string        `gorm:"column:slug;size:160;uniqueIndex;not null" json:"slug"`
	Status      ProjectStatus `gorm:"column:status;type:varchar(32);not null" json:"status"`
	HealthScore int           `gorm:"column:health_score;not null" json:"healthScore"`
	MRR         float64       `gorm:"column:mrr;default:0" json:"mrr"`
	CAC         float64       `gorm:"column:cac;default:0" json:"cac"`
	LTV         float64       `gorm:"column:ltv;default:0" json:"ltv"`
	ChurnRate   float64       `gorm:"column:churn_rate;default:0" json:"churnRate"`
	CreatedAt   time.Time     `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time     `gorm:"column:updated_at" json:"updatedAt"`
}

// Milestone represents a roadmap milestone.
type Milestone struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ProjectID   uuid.UUID `gorm:"column:project_id;type:uuid;index;not null" json:"projectId"`
	Title       string    `gorm:"column:title;size:120;not null" json:"title"`
	Description string    `gorm:"column:description;size:255" json:"description"`
	DueDate     time.Time `gorm:"column:due_date;index" json:"dueDate"`
	Completed   bool      `gorm:"column:completed;not null;default:false" json:"completed"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// NewProject constructs a new project aggregate.
func NewProject(userID uuid.UUID, name, description string, status ProjectStatus, healthScore int, mrr, cac, ltv, churn float64) (*Project, error) {
	project := &Project{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Status:      status,
		HealthScore: healthScore,
		MRR:         mrr,
		CAC:         cac,
		LTV:         ltv,
		ChurnRate:   churn,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	project.Slug = generateProjectSlug(project.Name, project.ID)

	return project, project.Validate()
}

// Validate ensures project invariants hold.
func (p *Project) Validate() error {
	if p == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilProject)
	}

	if p.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}

	if p.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if p.Name == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyProjectName)
	}

	if strings.TrimSpace(p.Slug) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectSlug)
	}

	switch p.Status {
	case ProjectStatusIdea, ProjectStatusPlanning, ProjectStatusExecuting, ProjectStatusMonitoring, ProjectStatusCompleted:
	default:
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}

	if p.HealthScore < 0 || p.HealthScore > 100 {
		return NewDomainError(ErrCodeInvalidHealthScore, ErrHealthScoreOutOfRange)
	}

	if p.MRR < 0 || p.CAC < 0 || p.LTV < 0 || p.ChurnRate < 0 {
		return NewDomainError(ErrCodeInvalidMetrics, ErrMetricsMustBePositive)
	}

	return nil
}

var slugSanitizer = regexp.MustCompile(`[^a-z0-9]+`)

func generateProjectSlug(name string, id uuid.UUID) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = slugSanitizer.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "project"
	}
	return slug
}

// UpdateStatus updates the stage and timestamp.
func (p *Project) UpdateStatus(status ProjectStatus) error {
	switch status {
	case ProjectStatusIdea, ProjectStatusPlanning, ProjectStatusExecuting, ProjectStatusMonitoring, ProjectStatusCompleted:
	default:
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}

	p.Status = status
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateMetrics adjusts KPI metrics for the project.
func (p *Project) UpdateMetrics(healthScore int, mrr, cac, ltv, churn float64) error {
	if healthScore < 0 || healthScore > 100 {
		return NewDomainError(ErrCodeInvalidHealthScore, ErrHealthScoreOutOfRange)
	}

	if mrr < 0 || cac < 0 || ltv < 0 || churn < 0 {
		return NewDomainError(ErrCodeInvalidMetrics, ErrMetricsMustBePositive)
	}

	p.HealthScore = healthScore
	p.MRR = mrr
	p.CAC = cac
	p.LTV = ltv
	p.ChurnRate = churn
	p.UpdatedAt = time.Now().UTC()

	return nil
}

// NewMilestone constructs a new milestone entry.
func NewMilestone(projectID uuid.UUID, title, description string, dueDate time.Time) (*Milestone, error) {
	m := &Milestone{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		DueDate:     dueDate.UTC(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	return m, m.Validate()
}

// Validate ensures milestone data integrity.
func (m *Milestone) Validate() error {
	if m == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilMilestone)
	}

	if m.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyMilestoneID)
	}

	if m.ProjectID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}

	if m.Title == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyMilestoneTitle)
	}

	return nil
}

// MarkCompleted toggles milestone completion.
func (m *Milestone) MarkCompleted(completed bool) {
	m.Completed = completed
	m.UpdatedAt = time.Now().UTC()
}

// KanbanColumn represents a column on the kanban board for a project.
type KanbanColumn struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ProjectID uuid.UUID `gorm:"column:project_id;type:uuid;index;not null" json:"projectId"`
	Name      string    `gorm:"column:name;size:80;not null" json:"name"`
	WIPLimit  int       `gorm:"column:wip_limit;not null;default:0" json:"wipLimit"`
	Position  int       `gorm:"column:position;not null;index:idx_kanban_column_position" json:"position"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// NewKanbanColumn constructs a new kanban column with sensible defaults.
func NewKanbanColumn(projectID uuid.UUID, name string, position int, wipLimit int) (*KanbanColumn, error) {
	column := &KanbanColumn{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      strings.TrimSpace(name),
		WIPLimit:  wipLimit,
		Position:  position,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return column, column.Validate()
}

// Validate enforces invariants for a kanban column.
func (c *KanbanColumn) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilKanbanColumn)
	}

	if c.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyKanbanColumnID)
	}

	if c.ProjectID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}

	if c.Name == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyKanbanColumnName)
	}

	if c.Position < 0 {
		return NewDomainError(ErrCodeInvalidPayload, ErrInvalidKanbanPosition)
	}

	if c.WIPLimit < 0 {
		return NewDomainError(ErrCodeInvalidPayload, ErrInvalidWIPLimit)
	}

	return nil
}

// Rename updates the column name.
func (c *KanbanColumn) Rename(name string) error {
	c.Name = strings.TrimSpace(name)
	c.UpdatedAt = time.Now().UTC()
	return c.Validate()
}

// SetWIPLimit updates the WIP limit value.
func (c *KanbanColumn) SetWIPLimit(limit int) error {
	c.WIPLimit = limit
	c.UpdatedAt = time.Now().UTC()
	return c.Validate()
}

// SetPosition updates the board order position.
func (c *KanbanColumn) SetPosition(position int) error {
	c.Position = position
	c.UpdatedAt = time.Now().UTC()
	return c.Validate()
}

// KanbanCard represents a task card on the kanban board.
type KanbanCard struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ProjectID   uuid.UUID  `gorm:"column:project_id;type:uuid;index;not null" json:"projectId"`
	ColumnID    uuid.UUID  `gorm:"column:column_id;type:uuid;index;not null" json:"columnId"`
	MilestoneID *uuid.UUID `gorm:"column:milestone_id;type:uuid" json:"milestoneId,omitempty"`
	Title       string     `gorm:"column:title;size:160;not null" json:"title"`
	Description string     `gorm:"column:description;size:512" json:"description"`
	DueDate     *time.Time `gorm:"column:due_date" json:"dueDate,omitempty"`
	Position    int        `gorm:"column:position;not null;index:idx_kanban_card_position" json:"position"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

// NewKanbanCard creates a new kanban card instance.
func NewKanbanCard(projectID, columnID uuid.UUID, title, description string, position int, dueDate *time.Time, milestoneID *uuid.UUID) (*KanbanCard, error) {
	card := &KanbanCard{
		ID:          uuid.New(),
		ProjectID:   projectID,
		ColumnID:    columnID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Position:    position,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if dueDate != nil {
		d := dueDate.UTC()
		card.DueDate = &d
	}

	if milestoneID != nil {
		id := *milestoneID
		card.MilestoneID = &id
	}

	return card, card.Validate()
}

// Validate enforces invariants for a kanban card.
func (c *KanbanCard) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilKanbanCard)
	}

	if c.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyKanbanCardID)
	}

	if c.ProjectID == uuid.Nil || c.ColumnID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}

	if c.Title == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyKanbanCardTitle)
	}

	if c.Position < 0 {
		return NewDomainError(ErrCodeInvalidPayload, ErrInvalidKanbanPosition)
	}

	return nil
}

// SetPosition adjusts the card order inside a column.
func (c *KanbanCard) SetPosition(position int) error {
	c.Position = position
	c.UpdatedAt = time.Now().UTC()
	return c.Validate()
}

// MoveToColumn reassigns card to another column.
func (c *KanbanCard) MoveToColumn(columnID uuid.UUID, position int) error {
	if columnID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyKanbanColumnID)
	}
	c.ColumnID = columnID
	return c.SetPosition(position)
}

// UpdateDetails updates mutable card fields.
func (c *KanbanCard) UpdateDetails(title, description string, dueDate *time.Time, milestoneID *uuid.UUID) error {
	if title != "" {
		c.Title = strings.TrimSpace(title)
	}
	if description != "" {
		c.Description = strings.TrimSpace(description)
	}
	if dueDate != nil {
		d := dueDate.UTC()
		c.DueDate = &d
	}
	if milestoneID != nil {
		id := *milestoneID
		c.MilestoneID = &id
	}
	c.UpdatedAt = time.Now().UTC()
	return c.Validate()
}

// ProjectDependency models dependency relationships between projects.
type ProjectDependency struct {
	ID                 uuid.UUID      `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ProjectID          uuid.UUID      `gorm:"column:project_id;type:uuid;index:idx_project_dependency,unique;not null" json:"projectId"`
	DependsOnProjectID uuid.UUID      `gorm:"column:depends_on_project_id;type:uuid;index:idx_project_dependency,unique;not null" json:"dependsOnProjectId"`
	Type               DependencyType `gorm:"column:type;size:32;not null" json:"type"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updatedAt"`
}

// NewProjectDependency constructs a dependency edge.
func NewProjectDependency(projectID, dependsOn uuid.UUID, depType DependencyType) (*ProjectDependency, error) {
	dependency := &ProjectDependency{
		ID:                 uuid.New(),
		ProjectID:          projectID,
		DependsOnProjectID: dependsOn,
		Type:               depType,
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}

	return dependency, dependency.Validate()
}

// Validate enforces invariants for dependency edges.
func (d *ProjectDependency) Validate() error {
	if d == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilDependency)
	}

	if d.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDependencyID)
	}

	if d.ProjectID == uuid.Nil || d.DependsOnProjectID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}

	if d.ProjectID == d.DependsOnProjectID {
		return NewDomainError(ErrCodeInvalidPayload, ErrSelfDependencyNotAllowed)
	}

	switch d.Type {
	case DependencyTypeBlocks, DependencyTypeRelates, DependencyTypeSupports:
	default:
		return NewDomainError(ErrCodeInvalidPayload, ErrUnsupportedDependencyType)
	}

	return nil
}

// UpdateType updates the dependency classification.
func (d *ProjectDependency) UpdateType(depType DependencyType) error {
	d.Type = depType
	d.UpdatedAt = time.Now().UTC()
	return d.Validate()
}
