package projects

import (
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

// Project captures high-level roadmap metadata.
type Project struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID     `gorm:"type:uuid;index;not null"`
	Name        string        `gorm:"size:120;not null"`
	Description string        `gorm:"size:255"`
	Status      ProjectStatus `gorm:"type:varchar(32);not null"`
	HealthScore int           `gorm:"not null"`
	MRR         float64       `gorm:"default:0"`
	CAC         float64       `gorm:"default:0"`
	LTV         float64       `gorm:"default:0"`
	ChurnRate   float64       `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Milestone represents a roadmap milestone.
type Milestone struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProjectID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Title       string    `gorm:"size:120;not null"`
	Description string    `gorm:"size:255"`
	DueDate     time.Time `gorm:"index"`
	Completed   bool      `gorm:"not null;default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
