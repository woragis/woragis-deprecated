package auth

import (
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Session represents a refresh-token backed session scoped to a device.
type Session struct {
	ID               uuid.UUID         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID         `gorm:"column:user_id;type:uuid;not null;index" json:"userId"`
	DeviceID         uuid.UUID         `gorm:"column:device_id;type:uuid;not null;index" json:"deviceId"`
	RefreshTokenHash string            `gorm:"column:refresh_token_hash;size:255;not null;uniqueIndex" json:"refreshTokenHash"`
	UserAgent        string            `gorm:"column:user_agent;size:512" json:"userAgent"`
	IP               string            `gorm:"column:ip;size:64" json:"ip"`
	Location         datatypes.JSONMap `gorm:"column:location;type:jsonb" json:"location"`
	CreatedAt        time.Time         `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time         `gorm:"column:updated_at" json:"updatedAt"`
	ExpiresAt        time.Time         `gorm:"column:expires_at" json:"expiresAt"`
	RevokedAt        *time.Time        `gorm:"column:revoked_at" json:"revokedAt,omitempty"`
	LastSeenAt       time.Time         `gorm:"column:last_seen_at" json:"lastSeenAt"`
}

// NewSession constructs a new session and validates its invariants.
func NewSession(userID, deviceID uuid.UUID, refreshTokenHash, userAgent, ip string, ttl time.Duration) (*Session, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if deviceID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyDeviceID)
	}
	refreshTokenHash = strings.TrimSpace(refreshTokenHash)
	if refreshTokenHash == "" {
		return nil, NewDomainError(ErrCodeInvalidToken, ErrEmptyRefreshTokenHash)
	}

	if ip != "" && net.ParseIP(ip) == nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrInvalidIPAddress)
	}

	now := time.Now().UTC()
	session := &Session{
		ID:               uuid.New(),
		UserID:           userID,
		DeviceID:         deviceID,
		RefreshTokenHash: refreshTokenHash,
		UserAgent:        truncate(userAgent, 512),
		IP:               truncate(ip, 64),
		CreatedAt:        now,
		UpdatedAt:        now,
		ExpiresAt:        now.Add(ttl),
		LastSeenAt:       now,
		Location:         datatypes.JSONMap{},
	}

	return session, session.Validate()
}

// Validate enforces session invariants.
func (s *Session) Validate() error {
	if s == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilSession)
	}
	if s.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySessionID)
	}
	if s.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if s.DeviceID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDeviceID)
	}
	if s.RefreshTokenHash == "" {
		return NewDomainError(ErrCodeInvalidToken, ErrEmptyRefreshTokenHash)
	}
	if s.ExpiresAt.Before(s.CreatedAt) {
		return NewDomainError(ErrCodeInvalidToken, ErrInvalidExpiry)
	}

	return nil
}

// Revoke marks the session as revoked.
func (s *Session) Revoke() {
	now := time.Now().UTC()
	s.RevokedAt = &now
	s.UpdatedAt = now
}

// Touch updates last seen information.
func (s *Session) Touch(ip, userAgent string) {
	now := time.Now().UTC()
	s.IP = truncate(ip, 64)
	s.UserAgent = truncate(userAgent, 512)
	s.LastSeenAt = now
	s.UpdatedAt = now
}

// IsActive reports if session is active and not expired.
func (s *Session) IsActive(reference time.Time) bool {
	if s.RevokedAt != nil {
		return false
	}
	return reference.UTC().Before(s.ExpiresAt)
}

// Device describes a physical or virtual client accessing the system.
type Device struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;not null;index" json:"userId"`
	Name        string    `gorm:"column:name;size:120" json:"name"`
	Fingerprint string    `gorm:"column:fingerprint;size:255;uniqueIndex" json:"fingerprint"`
	LastSeenAt  time.Time `gorm:"column:last_seen_at" json:"lastSeenAt"`
	IsTrusted   bool      `gorm:"column:is_trusted" json:"isTrusted"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// NewDevice constructs a device aggregate.
func NewDevice(userID uuid.UUID, name, fingerprint string) (*Device, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Unknown device"
	}
	fingerprint = strings.TrimSpace(fingerprint)
	if fingerprint == "" {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyDeviceFingerprint)
	}

	now := time.Now().UTC()
	device := &Device{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        truncate(name, 120),
		Fingerprint: fingerprint,
		LastSeenAt:  now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return device, nil
}

// Touch updates the device metadata.
func (d *Device) Touch(name string) {
	d.Name = truncate(strings.TrimSpace(name), 120)
	now := time.Now().UTC()
	d.LastSeenAt = now
	d.UpdatedAt = now
}

// PromoteTrusted marks the device as trusted.
func (d *Device) PromoteTrusted() {
	d.IsTrusted = true
	d.UpdatedAt = time.Now().UTC()
}

func truncate(value string, limit int) string {
	if len(value) <= limit {
		return value
	}
	return value[:limit]
}
