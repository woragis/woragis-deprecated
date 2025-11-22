package apikeys

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"time"

	"github.com/google/uuid"
)

// APIKey represents an API key for public read-only access.
type APIKey struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"column:name;size:120;not null" json:"name"`
	KeyHash     string    `gorm:"column:key_hash;size:255;not null;uniqueIndex:idx_api_key_hash" json:"-"`
	Prefix      string    `gorm:"column:prefix;size:20;not null;index" json:"prefix"` // First 8 chars for display
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	LastUsedAt  *time.Time `gorm:"column:last_used_at" json:"lastUsedAt,omitempty"`
	ExpiresAt   *time.Time `gorm:"column:expires_at;index" json:"expiresAt,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// APIKeyWithToken represents an API key with its plaintext token (only shown once on creation).
type APIKeyWithToken struct {
	APIKey
	Token string `json:"token"` // Plaintext token, only included on creation
}

// NewAPIKey creates a new API key entity.
func NewAPIKey(userID uuid.UUID, name string, expiresAt *time.Time) (*APIKey, string, error) {
	// Generate a secure random key
	token, err := generateSecureToken()
	if err != nil {
		return nil, "", err
	}

	// Hash the token for storage
	keyHash := hashToken(token)
	prefix := token[:8] // First 8 characters for display

	apiKey := &APIKey{
		ID:        uuid.New(),
		Name:      strings.TrimSpace(name),
		KeyHash:   keyHash,
		Prefix:    prefix,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return apiKey, token, apiKey.Validate()
}

// Validate ensures API key invariants hold.
func (k *APIKey) Validate() error {
	if k == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilAPIKey)
	}

	if k.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyAPIKeyID)
	}

	if k.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if k.Name == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyAPIKeyName)
	}

	if k.KeyHash == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyKeyHash)
	}

	if k.Prefix == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyPrefix)
	}

	return nil
}

// IsExpired checks if the API key has expired.
func (k *APIKey) IsExpired() bool {
	if k.ExpiresAt == nil {
		return false
	}
	return time.Now().UTC().After(*k.ExpiresAt)
}

// UpdateLastUsed updates the last used timestamp.
func (k *APIKey) UpdateLastUsed() {
	now := time.Now().UTC()
	k.LastUsedAt = &now
	k.UpdatedAt = now
}

// generateSecureToken generates a cryptographically secure random token.
func generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Base64 encode to get a URL-safe string
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
}

// hashToken creates a SHA256 hash of the token for storage.
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// VerifyToken checks if a provided token matches the stored hash.
func VerifyToken(providedToken string, storedHash string) bool {
	computedHash := hashToken(providedToken)
	return computedHash == storedHash
}

