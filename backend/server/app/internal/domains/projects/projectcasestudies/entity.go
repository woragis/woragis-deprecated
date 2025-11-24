package projectcasestudies

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TableName specifies the table name for ProjectCaseStudy.
func (ProjectCaseStudy) TableName() string {
	return "project_case_studies"
}

// ProjectCaseStudy represents a detailed case study for a project.
type ProjectCaseStudy struct {
	ID            uuid.UUID        `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ProjectID     uuid.UUID        `gorm:"column:project_id;type:uuid;index;not null" json:"projectId"`
	Title         string           `gorm:"column:title;size:255;not null" json:"title"`
	Description   string           `gorm:"column:description;type:text;not null" json:"description"`
	Challenge     string           `gorm:"column:challenge;type:text;not null" json:"challenge"`
	Solution      string           `gorm:"column:solution;type:text;not null" json:"solution"`
	Technologies  JSONArray        `gorm:"column:technologies;type:jsonb" json:"technologies"` // Array of strings
	Architecture  string           `gorm:"column:architecture;type:text" json:"architecture"`
	Metrics       *MetricsData     `gorm:"column:metrics;type:jsonb" json:"metrics,omitempty"`
	Tradeoffs     *TradeoffsData   `gorm:"column:tradeoffs;type:jsonb" json:"tradeoffs,omitempty"`
	LessonsLearned JSONArray       `gorm:"column:lessons_learned;type:jsonb" json:"lessonsLearned"` // Array of strings
	CreatedAt     time.Time        `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time        `gorm:"column:updated_at" json:"updatedAt"`
}

// MetricsData stores metrics with label, value, and improvement.
type MetricsData struct {
	Metrics []Metric `json:"metrics"`
}

// Metric represents a single metric.
type Metric struct {
	Label      string `json:"label"`
	Value      string `json:"value"`
	Improvement string `json:"improvement,omitempty"`
}

// TradeoffsData stores tradeoff decisions with pros and cons.
type TradeoffsData struct {
	Tradeoffs []Tradeoff `json:"tradeoffs"`
}

// Tradeoff represents a decision with pros and cons.
type Tradeoff struct {
	Decision string   `json:"decision"`
	Pros     []string `json:"pros"`
	Cons     []string `json:"cons"`
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

// Value implements the driver.Valuer interface for TradeoffsData.
func (t *TradeoffsData) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return json.Marshal(t)
}

// Scan implements the sql.Scanner interface for TradeoffsData.
func (t *TradeoffsData) Scan(value interface{}) error {
	if value == nil {
		*t = TradeoffsData{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), t)
	}
	return json.Unmarshal(bytes, t)
}

// NewProjectCaseStudy creates a new project case study entity.
func NewProjectCaseStudy(projectID uuid.UUID, title, description, challenge, solution, architecture string) (*ProjectCaseStudy, error) {
	cs := &ProjectCaseStudy{
		ID:            uuid.New(),
		ProjectID:     projectID,
		Title:         strings.TrimSpace(title),
		Description:   strings.TrimSpace(description),
		Challenge:     strings.TrimSpace(challenge),
		Solution:      strings.TrimSpace(solution),
		Architecture:  strings.TrimSpace(architecture),
		Technologies:  JSONArray{},
		LessonsLearned: JSONArray{},
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	return cs, cs.Validate()
}

// Validate ensures case study invariants hold.
func (c *ProjectCaseStudy) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilCaseStudy)
	}
	if c.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCaseStudyID)
	}
	if c.ProjectID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}
	if strings.TrimSpace(c.Title) == "" {
		return NewDomainError(ErrCodeInvalidTitle, ErrEmptyTitle)
	}
	if strings.TrimSpace(c.Description) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDescription)
	}
	if strings.TrimSpace(c.Challenge) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyChallenge)
	}
	if strings.TrimSpace(c.Solution) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySolution)
	}
	return nil
}

// UpdateDetails updates case study details.
func (c *ProjectCaseStudy) UpdateDetails(title, description, challenge, solution, architecture string) error {
	if title != "" {
		c.Title = strings.TrimSpace(title)
	}
	if description != "" {
		c.Description = strings.TrimSpace(description)
	}
	if challenge != "" {
		c.Challenge = strings.TrimSpace(challenge)
	}
	if solution != "" {
		c.Solution = strings.TrimSpace(solution)
	}
	if architecture != "" {
		c.Architecture = strings.TrimSpace(architecture)
	}
	c.UpdatedAt = time.Now().UTC()
	return c.Validate()
}

// SetTechnologies updates the technologies array.
func (c *ProjectCaseStudy) SetTechnologies(technologies []string) {
	c.Technologies = JSONArray(technologies)
	c.UpdatedAt = time.Now().UTC()
}

// SetMetrics updates the metrics data.
func (c *ProjectCaseStudy) SetMetrics(metrics *MetricsData) {
	c.Metrics = metrics
	c.UpdatedAt = time.Now().UTC()
}

// SetTradeoffs updates the tradeoffs data.
func (c *ProjectCaseStudy) SetTradeoffs(tradeoffs *TradeoffsData) {
	c.Tradeoffs = tradeoffs
	c.UpdatedAt = time.Now().UTC()
}

// SetLessonsLearned updates the lessons learned array.
func (c *ProjectCaseStudy) SetLessonsLearned(lessons []string) {
	c.LessonsLearned = JSONArray(lessons)
	c.UpdatedAt = time.Now().UTC()
}
