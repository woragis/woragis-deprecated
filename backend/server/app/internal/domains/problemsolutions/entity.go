package problemsolutions

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TableName specifies the table name for ProblemSolution.
func (ProblemSolution) TableName() string {
	return "problem_solutions"
}

// ProblemSolution represents a problem-solution document.
type ProblemSolution struct {
	ID          uuid.UUID     `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID     `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Problem     string        `gorm:"column:problem;type:text;not null" json:"problem"`
	Context     string        `gorm:"column:context;type:text;not null" json:"context"`
	Solution    string        `gorm:"column:solution;type:text;not null" json:"solution"`
	Technologies JSONArray    `gorm:"column:technologies;type:jsonb" json:"technologies"` // Array of strings
	Impact      string        `gorm:"column:impact;type:text" json:"impact"`
	Metrics     *MetricsData  `gorm:"column:metrics;type:jsonb" json:"metrics,omitempty"`
	Featured    bool          `gorm:"column:featured;not null;default:false;index" json:"featured"`
	CreatedAt   time.Time     `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
}

// MetricsData stores before/after metrics and improvement.
type MetricsData struct {
	Before     string `json:"before"`
	After      string `json:"after"`
	Improvement string `json:"improvement"`
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

// NewProblemSolution creates a new problem solution entity.
func NewProblemSolution(userID uuid.UUID, problem, context, solution string) (*ProblemSolution, error) {
	ps := &ProblemSolution{
		ID:        uuid.New(),
		UserID:    userID,
		Problem:   strings.TrimSpace(problem),
		Context:   strings.TrimSpace(context),
		Solution:  strings.TrimSpace(solution),
		Technologies: JSONArray{},
		Featured:   false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return ps, ps.Validate()
}

// Validate ensures problem solution invariants hold.
func (p *ProblemSolution) Validate() error {
	if p == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilProblemSolution)
	}
	if p.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProblemSolutionID)
	}
	if p.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if strings.TrimSpace(p.Problem) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProblem)
	}
	if strings.TrimSpace(p.Context) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyContext)
	}
	if strings.TrimSpace(p.Solution) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySolution)
	}
	return nil
}

// UpdateDetails updates problem solution details.
func (p *ProblemSolution) UpdateDetails(problem, context, solution, impact string) error {
	if problem != "" {
		p.Problem = strings.TrimSpace(problem)
	}
	if context != "" {
		p.Context = strings.TrimSpace(context)
	}
	if solution != "" {
		p.Solution = strings.TrimSpace(solution)
	}
	if impact != "" {
		p.Impact = strings.TrimSpace(impact)
	}
	p.UpdatedAt = time.Now().UTC()
	return p.Validate()
}

// SetTechnologies updates the technologies array.
func (p *ProblemSolution) SetTechnologies(technologies []string) {
	p.Technologies = JSONArray(technologies)
	p.UpdatedAt = time.Now().UTC()
}

// SetMetrics updates the metrics data.
func (p *ProblemSolution) SetMetrics(metrics *MetricsData) {
	p.Metrics = metrics
	p.UpdatedAt = time.Now().UTC()
}

// SetFeatured updates the featured flag.
func (p *ProblemSolution) SetFeatured(featured bool) {
	p.Featured = featured
	p.UpdatedAt = time.Now().UTC()
}

