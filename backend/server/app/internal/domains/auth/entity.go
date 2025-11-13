package auth

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User represents an authenticated account inside Woragis.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser constructs a User aggregate with the provided e-mail and password hash.
func NewUser(email, passwordHash string) (*User, error) {
	user := &User{
		ID:           uuid.New(),
		Email:        strings.TrimSpace(email),
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
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

	return nil
}

// UpdatePassword updates the password hash while keeping validation in place.
func (u *User) UpdatePassword(newHash string) error {
	if newHash == "" {
		return NewDomainError(ErrCodeInvalidPassword, ErrEmptyPasswordHash)
	}

	u.PasswordHash = newHash
	u.UpdatedAt = time.Now().UTC()

	return nil
}
