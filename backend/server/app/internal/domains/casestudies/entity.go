package casestudies

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

// CaseStudy represents a detailed case study for a project.
type CaseStudy struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	ProjectID   uuid.UUID `gorm:"column:project_id;type:uuid;index;not null" json:"projectId"` // Links to projects table
	ProjectSlug string    `gorm:"column:project_slug;size:160;not null;index" json:"projectSlug"` // For easy lookup
	Title       string    `gorm:"column:title;size:255;not null" json:"title"`
	Problem     string    `gorm:"column:problem;type:text;not null" json:"problem"`
	Context     string    `gorm:"column:context;type:text;not null" json:"context"`
	Solution    string    `gorm:"column:solution;type:text;not null" json:"solution"`
	Approach    JSONArray `gorm:"column:approach;type:jsonb" json:"approach"` // Array of strings
	Architecture *ArchitectureData `gorm:"column:architecture;type:jsonb" json:"architecture,omitempty"`
	Metrics     *MetricsData `gorm:"column:metrics;type:jsonb" json:"metrics,omitempty"`
	LessonsLearned JSONArray `gorm:"column:lessons_learned;type:jsonb" json:"lessonsLearned"` // Array of strings
	Technologies JSONArray `gorm:"column:technologies;type:jsonb" json:"technologies"` // Array of strings
	Featured    bool      `gorm:"column:featured;not null;default:false;index" json:"featured"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// ArchitectureData stores architecture diagram and component information.
type ArchitectureData struct {
	Diagram     string              `json:"diagram,omitempty"`     // Mermaid/PlantUML syntax
	DiagramType string              `json:"diagramType,omitempty"` // "mermaid" or "plantuml"
	Description string              `json:"description,omitempty"`
	Components  []ArchitectureComponent `json:"components,omitempty"`
}

// ArchitectureComponent represents a component in the architecture.
type ArchitectureComponent struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
}

// MetricsData stores before/after metrics and impact.
type MetricsData struct {
	Before []Metric `json:"before,omitempty"`
	After  []Metric `json:"after,omitempty"`
	Impact string   `json:"impact,omitempty"`
}

// Metric represents a single metric.
type Metric struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// JSONArray is a custom type for storing JSON arrays in PostgreSQL.
type JSONArray []string

// Value implements the driver.Valuer interface.
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface.
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), j)
	}
	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface for ArchitectureData.
func (a *ArchitectureData) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface for ArchitectureData.
func (a *ArchitectureData) Scan(value interface{}) error {
	if value == nil {
		*a = ArchitectureData{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), a)
	}
	return json.Unmarshal(bytes, a)
}

// Value implements the driver.Valuer interface for MetricsData.
func (m *MetricsData) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan implements the sql.Scanner interface for MetricsData.
func (m *MetricsData) Scan(value interface{}) error {
	if value == nil {
		*m = MetricsData{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), m)
	}
	return json.Unmarshal(bytes, m)
}

// NewCaseStudy creates a new case study entity.
func NewCaseStudy(userID, projectID uuid.UUID, projectSlug, title, problem, context, solution string) *CaseStudy {
	return &CaseStudy{
		ID:          uuid.New(),
		UserID:      userID,
		ProjectID:   projectID,
		ProjectSlug: strings.TrimSpace(projectSlug),
		Title:       strings.TrimSpace(title),
		Problem:     strings.TrimSpace(problem),
		Context:     strings.TrimSpace(context),
		Solution:    strings.TrimSpace(solution),
		Approach:    JSONArray{},
		LessonsLearned: JSONArray{},
		Technologies: JSONArray{},
		Featured:    false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

// Validate ensures case study invariants hold.
func (c *CaseStudy) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilCaseStudy)
	}
	if c.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if c.ProjectID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}
	if strings.TrimSpace(c.ProjectSlug) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectSlug)
	}
	if strings.TrimSpace(c.Title) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTitle)
	}
	if strings.TrimSpace(c.Problem) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProblem)
	}
	if strings.TrimSpace(c.Context) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyContext)
	}
	if strings.TrimSpace(c.Solution) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySolution)
	}
	return nil
}

