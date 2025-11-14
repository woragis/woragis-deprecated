package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// EmailTokenType represents the semantic of the token.
type EmailTokenType string

const (
	EmailTokenTypeConfirmation  EmailTokenType = "confirmation"
	EmailTokenTypePasswordReset EmailTokenType = "password_reset"
	EmailTokenTypeMagicLink     EmailTokenType = "magic_link"
)

// EmailToken persists single-use tokens for email flows.
type EmailToken struct {
	ID         uuid.UUID      `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID      `gorm:"column:user_id;type:uuid;not null;index" json:"userId"`
	Type       EmailTokenType `gorm:"column:type;size:32;index" json:"type"`
	TokenHash  string         `gorm:"column:token_hash;size:128;uniqueIndex" json:"tokenHash"`
	ExpiresAt  time.Time      `gorm:"column:expires_at;index" json:"expiresAt"`
	ConsumedAt *time.Time     `gorm:"column:consumed_at" json:"consumedAt,omitempty"`
	SentCount  int            `gorm:"column:sent_count" json:"sentCount"`
	LastSentAt *time.Time     `gorm:"column:last_sent_at" json:"lastSentAt,omitempty"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updatedAt"`
}

// NewEmailToken creates a new email token entry.
func NewEmailToken(userID uuid.UUID, tokenType EmailTokenType, rawToken string, ttl time.Duration) (*EmailToken, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if tokenType == "" {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyEmailTokenType)
	}
	if rawToken == "" {
		return nil, NewDomainError(ErrCodeInvalidToken, ErrEmptyEmailToken)
	}

	now := time.Now().UTC()
	token := &EmailToken{
		ID:         uuid.New(),
		UserID:     userID,
		Type:       tokenType,
		TokenHash:  hashToken(rawToken),
		ExpiresAt:  now.Add(ttl),
		CreatedAt:  now,
		UpdatedAt:  now,
		SentCount:  1,
		LastSentAt: &now,
	}

	return token, nil
}

// Consume marks the token as used.
func (t *EmailToken) Consume() {
	now := time.Now().UTC()
	t.ConsumedAt = &now
	t.UpdatedAt = now
}

// Touch dispatch bookkeeping for resend limits.
func (t *EmailToken) Touch() {
	now := time.Now().UTC()
	t.SentCount++
	t.LastSentAt = &now
	t.UpdatedAt = now
}

// IsExpired returns true when token expired or consumed.
func (t *EmailToken) IsExpired(reference time.Time) bool {
	if t.ConsumedAt != nil {
		return true
	}
	return reference.UTC().After(t.ExpiresAt)
}

// Matches compares raw token with stored hash.
func (t *EmailToken) Matches(rawToken string) bool {
	return t.TokenHash == hashToken(rawToken)
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
