package skills

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SkillCategory represents the category of a skill.
type SkillCategory string

const (
	SkillCategoryBackend       SkillCategory = "backend"
	SkillCategoryFrontend      SkillCategory = "frontend"
	SkillCategoryDatabase     SkillCategory = "database"
	SkillCategoryInfrastructure SkillCategory = "infrastructure"
	SkillCategoryDevOps        SkillCategory = "devops"
	SkillCategoryLanguage      SkillCategory = "language"
	SkillCategoryFramework     SkillCategory = "framework"
	SkillCategoryTool          SkillCategory = "tool"
	SkillCategoryService       SkillCategory = "service"
	SkillCategoryLibrary       SkillCategory = "library"
	SkillCategoryOther         SkillCategory = "other"
)

// Skill represents a global skill that can be attached to projects.
type Skill struct {
	ID          uuid.UUID    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name        string       `gorm:"column:name;size:120;not null;uniqueIndex:idx_skill_name" json:"name"`
	Slug        string       `gorm:"column:slug;size:160;not null;uniqueIndex:idx_skill_slug" json:"slug"`
	Category    SkillCategory `gorm:"column:category;type:varchar(32);not null;index" json:"category"`
	Description string       `gorm:"column:description;size:512" json:"description,omitempty"`
	Icon        string       `gorm:"column:icon;size:255" json:"icon,omitempty"` // Icon name/identifier (e.g., "redis", "postgresql")
	CreatedAt   time.Time    `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time    `gorm:"column:updated_at" json:"updatedAt"`
}

// ProjectSkill represents the many-to-many relationship between projects and skills.
type ProjectSkill struct {
	ProjectID uuid.UUID `gorm:"column:project_id;type:uuid;primaryKey;index" json:"projectId"`
	SkillID   uuid.UUID `gorm:"column:skill_id;type:uuid;primaryKey;index" json:"skillId"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

// TableName specifies the table name for ProjectSkill.
func (ProjectSkill) TableName() string {
	return "project_skills"
}

// NewSkill creates a new skill entity.
func NewSkill(name, description, icon string, category SkillCategory) (*Skill, error) {
	skill := &Skill{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Icon:        strings.TrimSpace(icon),
		Category:    category,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	skill.Slug = generateSkillSlug(skill.Name)

	return skill, skill.Validate()
}

// Validate ensures skill invariants hold.
func (s *Skill) Validate() error {
	if s == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilSkill)
	}

	if s.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySkillID)
	}

	if s.Name == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptySkillName)
	}

	if strings.TrimSpace(s.Slug) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySkillSlug)
	}

	switch s.Category {
	case SkillCategoryBackend, SkillCategoryFrontend, SkillCategoryDatabase,
		SkillCategoryInfrastructure, SkillCategoryDevOps, SkillCategoryLanguage,
		SkillCategoryFramework, SkillCategoryTool, SkillCategoryService,
		SkillCategoryLibrary, SkillCategoryOther:
	default:
		return NewDomainError(ErrCodeInvalidCategory, ErrUnsupportedCategory)
	}

	return nil
}

// UpdateDetails updates skill details.
func (s *Skill) UpdateDetails(name, description, icon string, category SkillCategory) error {
	if name != "" {
		s.Name = strings.TrimSpace(name)
		s.Slug = generateSkillSlug(s.Name)
	}
	if description != "" {
		s.Description = strings.TrimSpace(description)
	}
	if icon != "" {
		s.Icon = strings.TrimSpace(icon)
	}
	if category != "" {
		s.Category = category
	}
	s.UpdatedAt = time.Now().UTC()
	return s.Validate()
}

// generateSkillSlug creates a URL-friendly slug from a skill name.
func generateSkillSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = regexp.MustCompile(`^-+|-+$`).ReplaceAllString(slug, "")
	return slug
}

