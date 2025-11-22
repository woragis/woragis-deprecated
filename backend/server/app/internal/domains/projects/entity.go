package projects

import (
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
	UserID      uuid.UUID     `gorm:"column:user_id;type:uuid;index;not null;uniqueIndex:idx_projects_user_slug" json:"userId"`
	Name        string        `gorm:"column:name;size:120;not null" json:"name"`
	Description string        `gorm:"column:description;size:255" json:"description"`
	Slug        string        `gorm:"column:slug;size:160;not null;uniqueIndex:idx_projects_user_slug" json:"slug"`
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
	project.Slug = generateProjectSlug(project.Name)

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

func generateProjectSlug(name string) string {
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

// Documentation Visibility Types
type DocumentationVisibility string

const (
	VisibilityPublic        DocumentationVisibility = "public"
	VisibilityAuthenticated DocumentationVisibility = "authenticated"
	VisibilityCollaborators DocumentationVisibility = "collaborators"
)

// Documentation Section Types
type DocumentationSectionType string

const (
	SectionTypeOverview      DocumentationSectionType = "overview"
	SectionTypeArchitecture  DocumentationSectionType = "architecture"
	SectionTypeTechStack     DocumentationSectionType = "tech_stack"
	SectionTypeFileStructure DocumentationSectionType = "file_structure"
	SectionTypeAPIDocs       DocumentationSectionType = "api_documentation"
	SectionTypeDeployment    DocumentationSectionType = "deployment"
	SectionTypeContributing  DocumentationSectionType = "contributing"
	SectionTypeCustom        DocumentationSectionType = "custom"
)

// ProjectDocumentation is the main container for project documentation.
type ProjectDocumentation struct {
	ID          uuid.UUID              `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ProjectID   uuid.UUID              `gorm:"column:project_id;type:uuid;index;not null;uniqueIndex:idx_project_documentation" json:"projectId"`
	Visibility  DocumentationVisibility `gorm:"column:visibility;type:varchar(32);not null;default:'collaborators'" json:"visibility"`
	Version     int                    `gorm:"column:version;not null;default:1" json:"version"`
	CreatedAt   time.Time              `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time              `gorm:"column:updated_at" json:"updatedAt"`
}

// NewProjectDocumentation creates a new documentation container.
func NewProjectDocumentation(projectID uuid.UUID, visibility DocumentationVisibility) (*ProjectDocumentation, error) {
	doc := &ProjectDocumentation{
		ID:         uuid.New(),
		ProjectID:  projectID,
		Visibility: visibility,
		Version:    1,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	return doc, doc.Validate()
}

// Validate ensures documentation invariants hold.
func (d *ProjectDocumentation) Validate() error {
	if d == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilDocumentation)
	}
	if d.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDocumentationID)
	}
	if d.ProjectID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}
	switch d.Visibility {
	case VisibilityPublic, VisibilityAuthenticated, VisibilityCollaborators:
	default:
		return NewDomainError(ErrCodeInvalidVisibility, ErrUnsupportedVisibility)
	}
	return nil
}

// UpdateVisibility updates the documentation visibility setting.
func (d *ProjectDocumentation) UpdateVisibility(visibility DocumentationVisibility) error {
	switch visibility {
	case VisibilityPublic, VisibilityAuthenticated, VisibilityCollaborators:
	default:
		return NewDomainError(ErrCodeInvalidVisibility, ErrUnsupportedVisibility)
	}
	d.Visibility = visibility
	d.UpdatedAt = time.Now().UTC()
	return nil
}

// IncrementVersion increments the documentation version.
func (d *ProjectDocumentation) IncrementVersion() {
	d.Version++
	d.UpdatedAt = time.Now().UTC()
}

// DocumentationSection represents an individual section within project documentation.
type DocumentationSection struct {
	ID          uuid.UUID              `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	DocumentationID uuid.UUID          `gorm:"column:documentation_id;type:uuid;index;not null" json:"documentationId"`
	Type        DocumentationSectionType `gorm:"column:type;type:varchar(32);not null" json:"type"`
	Title       string                 `gorm:"column:title;size:160;not null" json:"title"`
	Content     string                 `gorm:"column:content;type:text" json:"content"`
	Position    int                    `gorm:"column:position;not null;default:0;index:idx_section_position" json:"position"`
	CreatedAt   time.Time              `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time              `gorm:"column:updated_at" json:"updatedAt"`
}

// NewDocumentationSection creates a new documentation section.
func NewDocumentationSection(documentationID uuid.UUID, sectionType DocumentationSectionType, title, content string, position int) (*DocumentationSection, error) {
	section := &DocumentationSection{
		ID:            uuid.New(),
		DocumentationID: documentationID,
		Type:          sectionType,
		Title:         strings.TrimSpace(title),
		Content:       content,
		Position:      position,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	return section, section.Validate()
}

// Validate ensures section invariants hold.
func (s *DocumentationSection) Validate() error {
	if s == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilDocumentationSection)
	}
	if s.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySectionID)
	}
	if s.DocumentationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDocumentationID)
	}
	if s.Title == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptySectionTitle)
	}
	switch s.Type {
	case SectionTypeOverview, SectionTypeArchitecture, SectionTypeTechStack, SectionTypeFileStructure,
		SectionTypeAPIDocs, SectionTypeDeployment, SectionTypeContributing, SectionTypeCustom:
	default:
		return NewDomainError(ErrCodeInvalidSectionType, ErrUnsupportedSectionType)
	}
	return nil
}

// UpdateContent updates the section content.
func (s *DocumentationSection) UpdateContent(content string) {
	s.Content = content
	s.UpdatedAt = time.Now().UTC()
}

// UpdateTitle updates the section title.
func (s *DocumentationSection) UpdateTitle(title string) error {
	s.Title = strings.TrimSpace(title)
	s.UpdatedAt = time.Now().UTC()
	return s.Validate()
}

// SetPosition updates the section position.
func (s *DocumentationSection) SetPosition(position int) {
	s.Position = position
	s.UpdatedAt = time.Now().UTC()
}

// Technology Category Types
type TechnologyCategory string

const (
	TechCategoryBackend       TechnologyCategory = "backend"
	TechCategoryDatabase      TechnologyCategory = "database"
	TechCategoryFrontend      TechnologyCategory = "frontend"
	TechCategoryInfrastructure TechnologyCategory = "infrastructure"
	TechCategoryMonitoring    TechnologyCategory = "monitoring"
	TechCategoryDevOps        TechnologyCategory = "devops"
	TechCategoryTesting       TechnologyCategory = "testing"
	TechCategoryOther         TechnologyCategory = "other"
)

// ProjectTechnology represents a technology used in a project.
type ProjectTechnology struct {
	ID          uuid.UUID          `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ProjectID   uuid.UUID          `gorm:"column:project_id;type:uuid;index;not null" json:"projectId"`
	Name        string             `gorm:"column:name;size:120;not null" json:"name"`
	Version     string             `gorm:"column:version;size:80" json:"version"`
	Category    TechnologyCategory `gorm:"column:category;type:varchar(32);not null" json:"category"`
	Purpose     string             `gorm:"column:purpose;size:512" json:"purpose"`
	Link        string             `gorm:"column:link;size:512" json:"link,omitempty"`
	CreatedAt   time.Time          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time          `gorm:"column:updated_at" json:"updatedAt"`
}

// NewProjectTechnology creates a new technology entry.
func NewProjectTechnology(projectID uuid.UUID, name, version string, category TechnologyCategory, purpose, link string) (*ProjectTechnology, error) {
	tech := &ProjectTechnology{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      strings.TrimSpace(name),
		Version:   strings.TrimSpace(version),
		Category:  category,
		Purpose:   strings.TrimSpace(purpose),
		Link:      strings.TrimSpace(link),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return tech, tech.Validate()
}

// Validate ensures technology invariants hold.
func (t *ProjectTechnology) Validate() error {
	if t == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilTechnology)
	}
	if t.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTechnologyID)
	}
	if t.ProjectID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}
	if t.Name == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyTechnologyName)
	}
	switch t.Category {
	case TechCategoryBackend, TechCategoryDatabase, TechCategoryFrontend, TechCategoryInfrastructure,
		TechCategoryMonitoring, TechCategoryDevOps, TechCategoryTesting, TechCategoryOther:
	default:
		return NewDomainError(ErrCodeInvalidTechCategory, ErrUnsupportedTechCategory)
	}
	return nil
}

// UpdateDetails updates technology details.
func (t *ProjectTechnology) UpdateDetails(name, version, purpose, link string) error {
	if name != "" {
		t.Name = strings.TrimSpace(name)
	}
	if version != "" {
		t.Version = strings.TrimSpace(version)
	}
	if purpose != "" {
		t.Purpose = strings.TrimSpace(purpose)
	}
	if link != "" {
		t.Link = strings.TrimSpace(link)
	}
	t.UpdatedAt = time.Now().UTC()
	return t.Validate()
}

// ProjectFileStructure represents a file or directory in the project structure.
type ProjectFileStructure struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ProjectID   uuid.UUID  `gorm:"column:project_id;type:uuid;index;not null" json:"projectId"`
	ParentID    *uuid.UUID `gorm:"column:parent_id;type:uuid;index" json:"parentId,omitempty"`
	Path        string     `gorm:"column:path;size:512;not null" json:"path"`
	Name        string     `gorm:"column:name;size:255;not null" json:"name"`
	IsDirectory bool       `gorm:"column:is_directory;not null;default:false" json:"isDirectory"`
	Language    string     `gorm:"column:language;size:32" json:"language,omitempty"`
	LineCount   int        `gorm:"column:line_count;default:0" json:"lineCount"`
	Purpose     string     `gorm:"column:purpose;size:512" json:"purpose,omitempty"`
	Position    int        `gorm:"column:position;not null;default:0;index:idx_file_structure_position" json:"position"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

// NewProjectFileStructure creates a new file structure entry.
func NewProjectFileStructure(projectID uuid.UUID, path, name string, isDirectory bool, parentID *uuid.UUID, language string, lineCount int, purpose string, position int) (*ProjectFileStructure, error) {
	fs := &ProjectFileStructure{
		ID:          uuid.New(),
		ProjectID:   projectID,
		ParentID:    parentID,
		Path:        strings.TrimSpace(path),
		Name:        strings.TrimSpace(name),
		IsDirectory: isDirectory,
		Language:    strings.TrimSpace(language),
		LineCount:   lineCount,
		Purpose:     strings.TrimSpace(purpose),
		Position:    position,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	return fs, fs.Validate()
}

// Validate ensures file structure invariants hold.
func (f *ProjectFileStructure) Validate() error {
	if f == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilFileStructure)
	}
	if f.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyFileStructureID)
	}
	if f.ProjectID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}
	if f.Path == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyFilePath)
	}
	if f.Name == "" {
		return NewDomainError(ErrCodeInvalidName, "projects: file name cannot be empty")
	}
	if f.LineCount < 0 {
		return NewDomainError(ErrCodeInvalidPayload, "projects: line count cannot be negative")
	}
	return nil
}

// UpdateDetails updates file structure details.
func (f *ProjectFileStructure) UpdateDetails(purpose string, lineCount int, language string) error {
	if purpose != "" {
		f.Purpose = strings.TrimSpace(purpose)
	}
	if lineCount >= 0 {
		f.LineCount = lineCount
	}
	if language != "" {
		f.Language = strings.TrimSpace(language)
	}
	f.UpdatedAt = time.Now().UTC()
	return f.Validate()
}

// SetPosition updates the file position within its parent.
func (f *ProjectFileStructure) SetPosition(position int) {
	f.Position = position
	f.UpdatedAt = time.Now().UTC()
}

// MoveToParent updates the parent directory.
func (f *ProjectFileStructure) MoveToParent(parentID *uuid.UUID) {
	f.ParentID = parentID
	f.UpdatedAt = time.Now().UTC()
}

// Architecture Diagram Types
type ArchitectureDiagramType string

const (
	DiagramTypeDependency   ArchitectureDiagramType = "dependency"
	DiagramTypeComponent   ArchitectureDiagramType = "component"
	DiagramTypeDataFlow    ArchitectureDiagramType = "data_flow"
	DiagramTypeInfrastructure ArchitectureDiagramType = "infrastructure"
	DiagramTypeCustom       ArchitectureDiagramType = "custom"
)

// ProjectArchitectureDiagram represents an architecture diagram for a project.
type ProjectArchitectureDiagram struct {
	ID          uuid.UUID              `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ProjectID   uuid.UUID              `gorm:"column:project_id;type:uuid;index;not null" json:"projectId"`
	Type        ArchitectureDiagramType `gorm:"column:type;type:varchar(32);not null" json:"type"`
	Title       string                 `gorm:"column:title;size:160;not null" json:"title"`
	Description string                 `gorm:"column:description;size:512" json:"description"`
	Content     string                 `gorm:"column:content;type:text" json:"content"`
	Format      string                 `gorm:"column:format;size:32;default:'mermaid'" json:"format"`
	ImageURL    string                 `gorm:"column:image_url;size:512" json:"imageUrl,omitempty"`
	CreatedAt   time.Time              `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time              `gorm:"column:updated_at" json:"updatedAt"`
}

// NewProjectArchitectureDiagram creates a new architecture diagram.
func NewProjectArchitectureDiagram(projectID uuid.UUID, diagramType ArchitectureDiagramType, title, description, content, format, imageURL string) (*ProjectArchitectureDiagram, error) {
	diagram := &ProjectArchitectureDiagram{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Type:        diagramType,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Content:     content,
		Format:      strings.TrimSpace(format),
		ImageURL:    strings.TrimSpace(imageURL),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	if diagram.Format == "" {
		diagram.Format = "mermaid"
	}
	return diagram, diagram.Validate()
}

// Validate ensures diagram invariants hold.
func (d *ProjectArchitectureDiagram) Validate() error {
	if d == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilArchitectureDiagram)
	}
	if d.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDiagramID)
	}
	if d.ProjectID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}
	if d.Title == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyDiagramTitle)
	}
	switch d.Type {
	case DiagramTypeDependency, DiagramTypeComponent, DiagramTypeDataFlow, DiagramTypeInfrastructure, DiagramTypeCustom:
	default:
		return NewDomainError(ErrCodeInvalidDiagramType, ErrUnsupportedDiagramType)
	}
	return nil
}

// UpdateContent updates the diagram content.
func (d *ProjectArchitectureDiagram) UpdateContent(content string) {
	d.Content = content
	d.UpdatedAt = time.Now().UTC()
}

// UpdateDetails updates diagram metadata.
func (d *ProjectArchitectureDiagram) UpdateDetails(title, description, imageURL string) error {
	if title != "" {
		d.Title = strings.TrimSpace(title)
	}
	if description != "" {
		d.Description = strings.TrimSpace(description)
	}
	if imageURL != "" {
		d.ImageURL = strings.TrimSpace(imageURL)
	}
	d.UpdatedAt = time.Now().UTC()
	return d.Validate()
}
