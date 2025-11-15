package auth

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User represents an authenticated account inside Woragis.
type User struct {
	ID                uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Email             string     `gorm:"column:email;uniqueIndex;size:255;not null" json:"email"`
	PasswordHash      string     `gorm:"column:password_hash;size:255;not null" json:"passwordHash"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	EmailConfirmedAt  *time.Time `gorm:"column:email_confirmed_at" json:"emailConfirmedAt,omitempty"`
	LastLoginAt       *time.Time `gorm:"column:last_login_at" json:"lastLoginAt,omitempty"`
	PasswordUpdatedAt *time.Time `gorm:"column:password_updated_at" json:"passwordUpdatedAt,omitempty"`
	Role              string     `gorm:"column:role;size:50;default:user" json:"role"`
	MFAEnabled        bool       `gorm:"column:mfa_enabled" json:"mfaEnabled"`
	MFAMethod         string     `gorm:"column:mfa_method;size:32" json:"mfaMethod"`
	PreferredLocale   string     `gorm:"column:preferred_locale;size:10;default:en" json:"preferredLocale"`
	PhoneNumber       string     `gorm:"column:phone_number;size:20;index" json:"phoneNumber,omitempty"`
}

// NewUser constructs a User aggregate with the provided e-mail and password hash.
func NewUser(email, passwordHash string) (*User, error) {
	now := time.Now().UTC()
	user := &User{
		ID:                uuid.New(),
		Email:             strings.TrimSpace(email),
		PasswordHash:      strings.TrimSpace(passwordHash),
		CreatedAt:         now,
		UpdatedAt:         now,
		PasswordUpdatedAt: &now,
		Role:              "user",
	}

	return user, user.Validate()
}

// Validate enforces domain invariants for the user entity.
func (u *User) Validate() error {
	if u == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilUser)
	}

	if u.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if u.Email == "" {
		return NewDomainError(ErrCodeInvalidEmail, ErrEmptyEmail)
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return NewDomainError(ErrCodeInvalidEmail, ErrInvalidEmailFormat)
	}

	if len(u.PasswordHash) == 0 {
		return NewDomainError(ErrCodeInvalidPassword, ErrEmptyPasswordHash)
	}

	if len(u.PasswordHash) > 255 {
		return NewDomainError(ErrCodeInvalidPassword, ErrPasswordTooLong)
	}

	return nil
}

// UpdatePassword updates the password hash while keeping validation in place.
func (u *User) UpdatePassword(newHash string) error {
	if newHash == "" {
		return NewDomainError(ErrCodeInvalidPassword, ErrEmptyPasswordHash)
	}

	u.PasswordHash = strings.TrimSpace(newHash)
	now := time.Now().UTC()
	u.UpdatedAt = now
	u.PasswordUpdatedAt = &now

	return nil
}

// ConfirmEmail marks the account as verified.
func (u *User) ConfirmEmail() {
	now := time.Now().UTC()
	u.EmailConfirmedAt = &now
	u.UpdatedAt = now
}

// MarkLogin updates bookkeeping fields for successful authentication.
func (u *User) MarkLogin() {
	now := time.Now().UTC()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}
