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
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID  `gorm:"column:user_id;type:uuid;not null;index" json:"userId"`
	Type        MFAType    `gorm:"column:type;size:32;not null;index" json:"type"`
	Secret      string     `gorm:"column:secret;size:160;not null" json:"secret"`
	Issuer      string     `gorm:"column:issuer;size:64" json:"issuer"`
	Label       string     `gorm:"column:label;size:128" json:"label"`
	ActivatedAt *time.Time `gorm:"column:activated_at" json:"activatedAt,omitempty"`
	LastUsedAt  *time.Time `gorm:"column:last_used_at" json:"lastUsedAt,omitempty"`
	RevokedAt   *time.Time `gorm:"column:revoked_at" json:"revokedAt,omitempty"`
	BackupCodes []string   `gorm:"column:backup_codes;serializer:json" json:"backupCodes"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updatedAt"`
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
