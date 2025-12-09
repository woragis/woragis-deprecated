package userprofiles

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// UserProfile represents a user's knowledge base/profile information.
type UserProfile struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;uniqueIndex;not null" json:"userId"`
	AboutMe   string    `gorm:"column:about_me;type:text" json:"aboutMe"` // HTML content for about me, languages, hobbies, etc.
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for UserProfile.
func (UserProfile) TableName() string {
	return "user_profiles"
}

// NewUserProfile creates a new user profile entity.
func NewUserProfile(userID uuid.UUID, aboutMe string) (*UserProfile, error) {
	profile := &UserProfile{
		ID:        uuid.New(),
		UserID:    userID,
		AboutMe:   strings.TrimSpace(aboutMe),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return profile, profile.Validate()
}

// Validate ensures user profile invariants hold.
func (p *UserProfile) Validate() error {
	if p == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilProfile)
	}

	if p.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyProfileID)
	}

	if p.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	return nil
}

// UpdateAboutMe updates the about me content.
func (p *UserProfile) UpdateAboutMe(aboutMe string) {
	p.AboutMe = strings.TrimSpace(aboutMe)
	p.UpdatedAt = time.Now().UTC()
}

