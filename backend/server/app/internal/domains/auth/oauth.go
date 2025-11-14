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
	ID             uuid.UUID     `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID     `gorm:"column:user_id;type:uuid;not null;index" json:"userId"`
	Provider       OAuthProvider `gorm:"column:provider;size:32;not null;index:idx_oauth_provider_user,priority:1" json:"provider"`
	ProviderUserID string        `gorm:"column:provider_user_id;size:191;not null;index:idx_oauth_provider_user,priority:2" json:"providerUserId"`
	AccessToken    string        `gorm:"column:access_token;type:text" json:"accessToken"`
	RefreshToken   string        `gorm:"column:refresh_token;type:text" json:"refreshToken"`
	ExpiresAt      *time.Time    `gorm:"column:expires_at" json:"expiresAt,omitempty"`
	Scopes         string        `gorm:"column:scopes;type:text" json:"scopes"`
	CreatedAt      time.Time     `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time     `gorm:"column:updated_at" json:"updatedAt"`
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
