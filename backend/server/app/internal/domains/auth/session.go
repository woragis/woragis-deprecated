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
	ID               uuid.UUID         `gorm:"type:uuid;primaryKey"`
	UserID           uuid.UUID         `gorm:"type:uuid;not null;index"`
	DeviceID         uuid.UUID         `gorm:"type:uuid;not null;index"`
	RefreshTokenHash string            `gorm:"size:255;not null;uniqueIndex"`
	UserAgent        string            `gorm:"size:512"`
	IP               string            `gorm:"size:64"`
	Location         datatypes.JSONMap `gorm:"type:jsonb"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ExpiresAt        time.Time
	RevokedAt        *time.Time
	LastSeenAt       time.Time
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
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"size:120"`
	Fingerprint string    `gorm:"size:255;uniqueIndex"`
	LastSeenAt  time.Time
	IsTrusted   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
