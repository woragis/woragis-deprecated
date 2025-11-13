package auth

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// OAuthProvider enumerates supported providers.
type OAuthProvider string

const (
	OAuthProviderGoogle    OAuthProvider = "google"
	OAuthProviderGithub    OAuthProvider = "github"
	OAuthProviderMicrosoft OAuthProvider = "microsoft"
)

// OAuthAccount keeps external identity binding.
type OAuthAccount struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID     `gorm:"type:uuid;not null;index"`
	Provider       OAuthProvider `gorm:"size:32;not null;index:idx_oauth_provider_user,priority:1"`
	ProviderUserID string        `gorm:"size:191;not null;index:idx_oauth_provider_user,priority:2"`
	AccessToken    string        `gorm:"type:text"`
	RefreshToken   string        `gorm:"type:text"`
	ExpiresAt      *time.Time
	Scopes         string `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewOAuthAccount constructs a new OAuth account aggregate.
func NewOAuthAccount(userID uuid.UUID, provider OAuthProvider, providerUserID string) (*OAuthAccount, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if provider == "" {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyOAuthProvider)
	}
	providerUserID = strings.TrimSpace(providerUserID)
	if providerUserID == "" {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyOAuthSubject)
	}

	now := time.Now().UTC()
	return &OAuthAccount{
		ID:             uuid.New(),
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: providerUserID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// UpdateTokens refreshes OAuth tokens.
func (o *OAuthAccount) UpdateTokens(accessToken, refreshToken string, expiresAt *time.Time, scopes []string) {
	o.AccessToken = strings.TrimSpace(accessToken)
	o.RefreshToken = strings.TrimSpace(refreshToken)
	o.ExpiresAt = expiresAt
	o.Scopes = strings.Join(scopes, " ")
	o.UpdatedAt = time.Now().UTC()
}
