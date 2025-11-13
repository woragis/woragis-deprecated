package auth

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User represents an authenticated account inside Woragis.
type User struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email             string    `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash      string    `gorm:"size:255;not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	EmailConfirmedAt  *time.Time
	LastLoginAt       *time.Time
	PasswordUpdatedAt *time.Time
	Role              string `gorm:"size:50;default:user"`
	MFAEnabled        bool
	MFAMethod         string `gorm:"size:32"`
	PreferredLocale   string `gorm:"size:10;default:en"`
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
