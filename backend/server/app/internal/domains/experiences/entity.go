package experiences

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// ExperienceType represents the type of work experience.
type ExperienceType string

const (
	ExperienceTypeFullTime  ExperienceType = "full-time"
	ExperienceTypeFreelance ExperienceType = "freelance"
	ExperienceTypeContract  ExperienceType = "contract"
	ExperienceTypeInternship ExperienceType = "internship"
)

// Experience represents a professional work experience.
type Experience struct {
	ID           uuid.UUID      `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID      `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Company      string         `gorm:"column:company;size:200;not null" json:"company"`
	Position     string         `gorm:"column:position;size:200;not null" json:"position"`
	PeriodStart  *time.Time     `gorm:"column:period_start" json:"periodStart,omitempty"`
	PeriodEnd    *time.Time     `gorm:"column:period_end" json:"periodEnd,omitempty"`
	PeriodText   string         `gorm:"column:period_text;size:100" json:"periodText,omitempty"` // e.g., "2022 - Presente"
	Location     string         `gorm:"column:location;size:200" json:"location,omitempty"`
	Description  string         `gorm:"column:description;type:text" json:"description,omitempty"`
	Type         ExperienceType `gorm:"column:type;type:varchar(32);not null;default:'full-time';index" json:"type"`
	CompanyURL   string         `gorm:"column:company_url;size:512" json:"companyUrl,omitempty"`
	LinkedInURL  string         `gorm:"column:linkedin_url;size:512" json:"linkedinUrl,omitempty"`
	DisplayOrder int            `gorm:"column:display_order;not null;default:0;index" json:"displayOrder"`
	IsCurrent    bool           `gorm:"column:is_current;default:false;index" json:"isCurrent"` // Currently working here
	CreatedAt    time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updatedAt"`
}

// ExperienceTechnology represents a technology used in an experience.
type ExperienceTechnology struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ExperienceID uuid.UUID `gorm:"column:experience_id;type:uuid;not null;index:idx_experience_tech" json:"experienceId"`
	Technology   string    `gorm:"column:technology;size:100;not null;index:idx_experience_tech" json:"technology"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for ExperienceTechnology.
func (ExperienceTechnology) TableName() string {
	return "experience_technologies"
}

// ExperienceProject represents a project related to an experience.
type ExperienceProject struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ExperienceID uuid.UUID `gorm:"column:experience_id;type:uuid;not null;index" json:"experienceId"`
	Name         string    `gorm:"column:name;size:200;not null" json:"name"`
	URL          string    `gorm:"column:url;size:512" json:"url,omitempty"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for ExperienceProject.
func (ExperienceProject) TableName() string {
	return "experience_projects"
}

// ExperienceAchievement represents an achievement or metric for an experience.
type ExperienceAchievement struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ExperienceID uuid.UUID `gorm:"column:experience_id;type:uuid;not null;index" json:"experienceId"`
	Metric       string    `gorm:"column:metric;size:100;not null" json:"metric"` // e.g., "40%", "50K+"
	Description  string    `gorm:"column:description;size:200;not null" json:"description"`
	Icon         string    `gorm:"column:icon;size:50" json:"icon,omitempty"` // Icon identifier (e.g., "TrendingUp", "Trophy")
	DisplayOrder int       `gorm:"column:display_order;not null;default:0" json:"displayOrder"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for ExperienceAchievement.
func (ExperienceAchievement) TableName() string {
	return "experience_achievements"
}

// NewExperience creates a new experience entity.
func NewExperience(userID uuid.UUID, company, position string) *Experience {
	return &Experience{
		ID:          uuid.New(),
		UserID:      userID,
		Company:     strings.TrimSpace(company),
		Position:    strings.TrimSpace(position),
		Type:        ExperienceTypeFullTime,
		IsCurrent:   false,
		DisplayOrder: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// SetType updates the experience type.
func (e *Experience) SetType(expType ExperienceType) error {
	if !isValidExperienceType(expType) {
		return NewDomainError(ErrCodeInvalidType, ErrUnsupportedExperienceType)
	}
	e.Type = expType
	e.UpdatedAt = time.Now()
	return nil
}

// SetDescription updates the description field.
func (e *Experience) SetDescription(description string) {
	e.Description = strings.TrimSpace(description)
	e.UpdatedAt = time.Now()
}

// SetPeriodText updates the period text field.
func (e *Experience) SetPeriodText(periodText string) {
	e.PeriodText = strings.TrimSpace(periodText)
	e.UpdatedAt = time.Now()
}

// SetIsCurrent updates the isCurrent field.
func (e *Experience) SetIsCurrent(isCurrent bool) {
	e.IsCurrent = isCurrent
	e.UpdatedAt = time.Now()
}

// Validate validates the experience entity.
func (e *Experience) Validate() error {
	if e.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if strings.TrimSpace(e.Company) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCompany)
	}
	if strings.TrimSpace(e.Position) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyPosition)
	}
	if !isValidExperienceType(e.Type) {
		return NewDomainError(ErrCodeInvalidType, ErrUnsupportedExperienceType)
	}
	return nil
}

// Validation helpers

func isValidExperienceType(et ExperienceType) bool {
	switch et {
	case ExperienceTypeFullTime, ExperienceTypeFreelance, ExperienceTypeContract, ExperienceTypeInternship:
		return true
	}
	return false
}

