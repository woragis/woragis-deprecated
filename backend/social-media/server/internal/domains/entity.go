package socialmedia

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"` // Hidden from JSON
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Role      string         `json:"role" gorm:"default:'user'"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	IsVerified bool          `json:"is_verified" gorm:"default:false"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Profile  *Profile   `json:"profile,omitempty" gorm:"foreignKey:UserID"`
	Sessions []Session  `json:"-" gorm:"foreignKey:UserID"`
}

// Profile represents user profile information
type Profile struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	Avatar      string         `json:"avatar"`
	Bio         string         `json:"bio" gorm:"type:text"`
	DateOfBirth *time.Time     `json:"date_of_birth"`
	Gender      string         `json:"gender" gorm:"type:varchar(20)"`
	Phone       string         `json:"phone"`
	Location    string         `json:"location"`
	Website     string         `json:"website"`
	SocialLinks string         `json:"social_links" gorm:"type:jsonb;serializer:json"` // JSON string for social media links
	Preferences string         `json:"preferences" gorm:"type:jsonb;serializer:json"` // JSON string for user preferences
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Session represents user authentication sessions
type Session struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	RefreshToken string         `json:"-" gorm:"not null;uniqueIndex"`
	UserAgent    string         `json:"user_agent"`
	IPAddress    string         `json:"ip_address"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	ExpiresAt    time.Time      `json:"expires_at" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// VerificationToken represents email verification tokens
type VerificationToken struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Token     string         `json:"-" gorm:"not null;uniqueIndex"`
	Type      string         `json:"type" gorm:"not null"` // "email_verification", "password_reset"
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	IsUsed    bool           `json:"is_used" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// TableName returns the table name for Profile
func (Profile) TableName() string {
	return "profiles"
}

// TableName returns the table name for Session
func (Session) TableName() string {
	return "sessions"
}

// TableName returns the table name for VerificationToken
func (VerificationToken) TableName() string {
	return "verification_tokens"
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsModerator checks if the user has moderator role
func (u *User) IsModerator() bool {
	return u.Role == "moderator" || u.Role == "admin"
}

// IsSessionValid checks if the session is still valid
func (s *Session) IsSessionValid() bool {
	return s.IsActive && s.ExpiresAt.After(time.Now())
}

// IsTokenValid checks if the verification token is still valid
func (vt *VerificationToken) IsTokenValid() bool {
	return !vt.IsUsed && vt.ExpiresAt.After(time.Now())
}
