package interests

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Interest represents a personal interest or area of focus.
type Interest struct {
	ID              uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Title           string    `gorm:"column:title;size:120;not null;uniqueIndex:idx_interest_title" json:"title"`
	Slug            string    `gorm:"column:slug;size:160;not null;uniqueIndex:idx_interest_slug" json:"slug"`
	Description     string    `gorm:"column:description;type:text;not null" json:"description"`
	Icon            string    `gorm:"column:icon;size:255" json:"icon,omitempty"` // Icon name/identifier (e.g., "Brain", "SiRedis", "GitBranch")
	Color           string    `gorm:"column:color;size:50" json:"color,omitempty"` // Color name (e.g., "pink-purple", "red-orange")
	BgGradient      string    `gorm:"column:bg_gradient;size:255" json:"bgGradient,omitempty"` // Tailwind gradient classes
	BorderColor     string    `gorm:"column:border_color;size:255" json:"borderColor,omitempty"` // Tailwind border classes
	HoverBorderColor string   `gorm:"column:hover_border_color;size:255" json:"hoverBorderColor,omitempty"` // Tailwind hover border classes
	ShadowColor     string    `gorm:"column:shadow_color;size:255" json:"shadowColor,omitempty"` // Tailwind shadow classes
	FullWidth       bool      `gorm:"column:full_width;not null;default:false" json:"fullWidth"`
	Featured        bool      `gorm:"column:featured;not null;default:false;index" json:"featured"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for Interest.
func (Interest) TableName() string {
	return "interests"
}

// NewInterest creates a new interest entity.
func NewInterest(title, description, icon, color, bgGradient, borderColor, hoverBorderColor, shadowColor string, fullWidth, featured bool) (*Interest, error) {
	interest := &Interest{
		ID:              uuid.New(),
		Title:           strings.TrimSpace(title),
		Description:     strings.TrimSpace(description),
		Icon:            strings.TrimSpace(icon),
		Color:           strings.TrimSpace(color),
		BgGradient:      strings.TrimSpace(bgGradient),
		BorderColor:     strings.TrimSpace(borderColor),
		HoverBorderColor: strings.TrimSpace(hoverBorderColor),
		ShadowColor:     strings.TrimSpace(shadowColor),
		FullWidth:       fullWidth,
		Featured:        featured,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	interest.Slug = generateInterestSlug(interest.Title)

	return interest, interest.Validate()
}

// Validate ensures interest invariants hold.
func (i *Interest) Validate() error {
	if i == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilInterest)
	}

	if i.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyInterestID)
	}

	if i.Title == "" {
		return NewDomainError(ErrCodeInvalidTitle, ErrEmptyInterestTitle)
	}

	if strings.TrimSpace(i.Slug) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyInterestSlug)
	}

	if strings.TrimSpace(i.Description) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyInterestDescription)
	}

	return nil
}

// UpdateDetails updates interest details.
func (i *Interest) UpdateDetails(title, description, icon, color, bgGradient, borderColor, hoverBorderColor, shadowColor string, fullWidth, featured *bool) error {
	if title != "" {
		i.Title = strings.TrimSpace(title)
		i.Slug = generateInterestSlug(i.Title)
	}
	if description != "" {
		i.Description = strings.TrimSpace(description)
	}
	if icon != "" {
		i.Icon = strings.TrimSpace(icon)
	}
	if color != "" {
		i.Color = strings.TrimSpace(color)
	}
	if bgGradient != "" {
		i.BgGradient = strings.TrimSpace(bgGradient)
	}
	if borderColor != "" {
		i.BorderColor = strings.TrimSpace(borderColor)
	}
	if hoverBorderColor != "" {
		i.HoverBorderColor = strings.TrimSpace(hoverBorderColor)
	}
	if shadowColor != "" {
		i.ShadowColor = strings.TrimSpace(shadowColor)
	}
	if fullWidth != nil {
		i.FullWidth = *fullWidth
	}
	if featured != nil {
		i.Featured = *featured
	}
	i.UpdatedAt = time.Now().UTC()
	return i.Validate()
}

// generateInterestSlug creates a URL-friendly slug from an interest title.
func generateInterestSlug(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = regexp.MustCompile(`^-+|-+$`).ReplaceAllString(slug, "")
	return slug
}

