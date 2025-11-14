package auth

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditAction enumerates domain events persisted for compliance.
type AuditAction string

const (
	AuditActionUserRegistered         AuditAction = "user_registered"
	AuditActionEmailConfirmed         AuditAction = "email_confirmed"
	AuditActionLoginSucceeded         AuditAction = "login_succeeded"
	AuditActionLoginFailed            AuditAction = "login_failed"
	AuditActionPasswordResetRequested AuditAction = "password_reset_requested"
	AuditActionPasswordResetCompleted AuditAction = "password_reset_completed"
	AuditActionSessionRevoked         AuditAction = "session_revoked"
	AuditActionMFAEnabled             AuditAction = "mfa_enabled"
	AuditActionMFADisabled            AuditAction = "mfa_disabled"
	AuditActionOAuthLinked            AuditAction = "oauth_linked"
	AuditActionOAuthUnlinked          AuditAction = "oauth_unlinked"
	AuditActionBulkUserAction         AuditAction = "bulk_user_action"
)

// AuditLog captures security-sensitive operations.
type AuditLog struct {
	ID        uuid.UUID   `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID    *uuid.UUID  `gorm:"column:user_id;type:uuid;index" json:"userId,omitempty"`
	Action    AuditAction `gorm:"column:action;size:64;index" json:"action"`
	Metadata  []byte      `gorm:"column:metadata;type:jsonb" json:"metadata"`
	IP        string      `gorm:"column:ip;size:64" json:"ip"`
	UserAgent string      `gorm:"column:user_agent;size:512" json:"userAgent"`
	CreatedAt time.Time   `gorm:"column:created_at" json:"createdAt"`
}

// NewAuditLog constructs a new audit entry.
func NewAuditLog(userID *uuid.UUID, action AuditAction, metadata map[string]any, ip, userAgent string) (*AuditLog, error) {
	if action == "" {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyAuditAction)
	}

	metaBytes, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	var uidPtr *uuid.UUID
	if userID != nil && *userID != uuid.Nil {
		uid := *userID
		uidPtr = &uid
	}

	return &AuditLog{
		ID:        uuid.New(),
		UserID:    uidPtr,
		Action:    action,
		Metadata:  metaBytes,
		IP:        truncate(ip, 64),
		UserAgent: truncate(userAgent, 512),
		CreatedAt: now,
	}, nil
}
