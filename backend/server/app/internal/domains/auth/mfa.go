package auth

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
	"time"

	"github.com/google/uuid"
)

// MFAType identifies the multi-factor mechanism.
type MFAType string

const (
	// MFATypeTOTP represents time-based one-time password.
	MFATypeTOTP MFAType = "totp"
	// MFATypeBackupCode represents static recovery codes.
	MFATypeBackupCode MFAType = "backup_code"
)

// MFAToken holds enrollment information per user.
type MFAToken struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Type        MFAType   `gorm:"size:32;not null;index"`
	Secret      string    `gorm:"size:160;not null"`
	Issuer      string    `gorm:"size:64"`
	Label       string    `gorm:"size:128"`
	ActivatedAt *time.Time
	LastUsedAt  *time.Time
	RevokedAt   *time.Time
	BackupCodes []string `gorm:"serializer:json"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewMFAToken generates a new MFA enrollment.
func NewMFAToken(userID uuid.UUID, tokenType MFAType, issuer, label string) (*MFAToken, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if tokenType == "" {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyMFAType)
	}

	secret, err := generateMFASecret()
	if err != nil {
		return nil, NewDomainError(ErrCodeMFAGenerationFailure, ErrUnableToGenerateMFASecret)
	}

	now := time.Now().UTC()
	token := &MFAToken{
		ID:          uuid.New(),
		UserID:      userID,
		Type:        tokenType,
		Secret:      secret,
		Issuer:      truncate(strings.TrimSpace(issuer), 64),
		Label:       truncate(strings.TrimSpace(label), 128),
		BackupCodes: generateBackupCodes(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return token, nil
}

// Activate marks the MFA token as active.
func (m *MFAToken) Activate() {
	now := time.Now().UTC()
	m.ActivatedAt = &now
	m.UpdatedAt = now
}

// Revoke marks the MFA token as revoked.
func (m *MFAToken) Revoke() {
	now := time.Now().UTC()
	m.RevokedAt = &now
	m.UpdatedAt = now
}

// Touch updates last used timestamp.
func (m *MFAToken) Touch() {
	now := time.Now().UTC()
	m.LastUsedAt = &now
	m.UpdatedAt = now
}

func generateMFASecret() (string, error) {
	secretBytes := make([]byte, 20)
	if _, err := rand.Read(secretBytes); err != nil {
		return "", err
	}

	return strings.TrimRight(base32.StdEncoding.EncodeToString(secretBytes), "="), nil
}

func generateBackupCodes() []string {
	const (
		count = 8
		size  = 10
	)

	codes := make([]string, count)
	raw := make([]byte, size)

	for i := 0; i < count; i++ {
		if _, err := rand.Read(raw); err != nil {
			continue
		}
		codes[i] = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(raw)[:size]
	}

	return codes
}
