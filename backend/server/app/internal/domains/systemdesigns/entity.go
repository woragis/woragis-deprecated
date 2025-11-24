package systemdesigns

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TableName specifies the table name for SystemDesign.
func (SystemDesign) TableName() string {
	return "system_designs"
}

// SystemDesign represents a system design document.
type SystemDesign struct {
	ID          uuid.UUID        `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID        `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Title       string           `gorm:"column:title;size:255;not null" json:"title"`
	Description string           `gorm:"column:description;type:text;not null" json:"description"`
	Components  *ComponentsData  `gorm:"column:components;type:jsonb" json:"components,omitempty"`
	DataFlow    string           `gorm:"column:data_flow;type:text" json:"dataFlow,omitempty"`
	Scalability string           `gorm:"column:scalability;type:text" json:"scalability,omitempty"`
	Reliability string           `gorm:"column:reliability;type:text" json:"reliability,omitempty"`
	Diagram     string           `gorm:"column:diagram;size:512" json:"diagram,omitempty"`
	Featured    bool             `gorm:"column:featured;not null;default:false;index" json:"featured"`
	CreatedAt   time.Time        `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time        `gorm:"column:updated_at" json:"updatedAt"`
}

// ComponentsData stores the array of components.
type ComponentsData struct {
	Components []Component `json:"components"`
}

// Component represents a component in the system design.
type Component struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Technology  string `json:"technology"`
}

// Value implements the driver.Valuer interface for ComponentsData.
func (c *ComponentsData) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

// Scan implements the sql.Scanner interface for ComponentsData.
func (c *ComponentsData) Scan(value interface{}) error {
	if value == nil {
		*c = ComponentsData{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), c)
	}
	return json.Unmarshal(bytes, c)
}

// NewSystemDesign creates a new system design entity.
func NewSystemDesign(userID uuid.UUID, title, description string) (*SystemDesign, error) {
	sd := &SystemDesign{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Featured:    false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	return sd, sd.Validate()
}

// Validate ensures system design invariants hold.
func (s *SystemDesign) Validate() error {
	if s == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilSystemDesign)
	}
	if s.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySystemDesignID)
	}
	if s.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if strings.TrimSpace(s.Title) == "" {
		return NewDomainError(ErrCodeInvalidTitle, ErrEmptyTitle)
	}
	if strings.TrimSpace(s.Description) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDescription)
	}
	return nil
}

// UpdateDetails updates system design details.
func (s *SystemDesign) UpdateDetails(title, description, dataFlow, scalability, reliability, diagram string) error {
	if title != "" {
		s.Title = strings.TrimSpace(title)
	}
	if description != "" {
		s.Description = strings.TrimSpace(description)
	}
	if dataFlow != "" {
		s.DataFlow = strings.TrimSpace(dataFlow)
	}
	if scalability != "" {
		s.Scalability = strings.TrimSpace(scalability)
	}
	if reliability != "" {
		s.Reliability = strings.TrimSpace(reliability)
	}
	if diagram != "" {
		s.Diagram = strings.TrimSpace(diagram)
	}
	s.UpdatedAt = time.Now().UTC()
	return s.Validate()
}

// SetComponents updates the components data.
func (s *SystemDesign) SetComponents(components *ComponentsData) {
	s.Components = components
	s.UpdatedAt = time.Now().UTC()
}

// SetFeatured updates the featured flag.
func (s *SystemDesign) SetFeatured(featured bool) {
	s.Featured = featured
	s.UpdatedAt = time.Now().UTC()
}

