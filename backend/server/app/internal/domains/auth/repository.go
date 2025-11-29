package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// Repository abstraction for persistence.
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdatePasswordHash(ctx context.Context, id uuid.UUID, hash string) error
	Update(ctx context.Context, user *User) error

	CreateSession(ctx context.Context, session *Session) error
	UpdateSession(ctx context.Context, session *Session) error
	FindSessionByID(ctx context.Context, id uuid.UUID) (*Session, error)
	FindSessionByRefreshHash(ctx context.Context, hash string) (*Session, error)
	ListActiveSessions(ctx context.Context, userID uuid.UUID) ([]Session, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RevokeUserSessions(ctx context.Context, userID uuid.UUID, exclude uuid.UUID) error

	UpsertDevice(ctx context.Context, device *Device) (*Device, error)
	ListDevices(ctx context.Context, userID uuid.UUID) ([]Device, error)

	CreateMFAToken(ctx context.Context, token *MFAToken) error
	FindActiveMFAToken(ctx context.Context, userID uuid.UUID, tokenType MFAType) (*MFAToken, error)
	UpdateMFAToken(ctx context.Context, token *MFAToken) error
	DeleteMFAToken(ctx context.Context, id uuid.UUID) error

	InsertAuditLog(ctx context.Context, entry *AuditLog) error
	ListAuditLogs(ctx context.Context, userID uuid.UUID, limit int) ([]AuditLog, error)

	FindOAuthAccount(ctx context.Context, provider OAuthProvider, providerUserID string) (*OAuthAccount, error)
	UpsertOAuthAccount(ctx context.Context, account *OAuthAccount) error
	ListOAuthAccounts(ctx context.Context, userID uuid.UUID) ([]OAuthAccount, error)
	DeleteOAuthAccount(ctx context.Context, userID uuid.UUID, provider OAuthProvider) error

	CreateEmailToken(ctx context.Context, token *EmailToken) error
	FindEmailToken(ctx context.Context, tokenType EmailTokenType, tokenHash string) (*EmailToken, error)
	UpdateEmailToken(ctx context.Context, token *EmailToken) error
	DeleteEmailTokensByUser(ctx context.Context, userID uuid.UUID, tokenType EmailTokenType) error

	BulkUpdateUserStatus(ctx context.Context, userIDs []uuid.UUID, updates map[string]any) error

	// Admin operations
	ListUsers(ctx context.Context, limit, offset int, search string) ([]User, int64, error)
}

// GormRepository implements Repository using PostgreSQL via GORM.
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new repository instance.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// Create persists a new user.
func (r *GormRepository) Create(ctx context.Context, user *User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return mapPersistenceError(err)
	}

	return nil
}

// Update persists user changes.
func (r *GormRepository) Update(ctx context.Context, user *User) error {
	if user == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilUser)
	}

	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return mapPersistenceError(err)
	}

	return nil
}

// FindByEmail retrieves a user by email, returning a domain error when not found.
func (r *GormRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeUserNotFound, ErrUserNotFound)
		}
		return nil, mapPersistenceError(err)
	}

	return &user, nil
}

// FindByID retrieves a user by ID.
func (r *GormRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeUserNotFound, ErrUserNotFound)
		}
		return nil, mapPersistenceError(err)
	}

	return &user, nil
}

// UpdatePasswordHash updates the stored password hash.
func (r *GormRepository) UpdatePasswordHash(ctx context.Context, id uuid.UUID, hash string) error {
	if hash == "" {
		return NewDomainError(ErrCodeInvalidPassword, ErrEmptyPasswordHash)
	}

	result := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"password_hash": hash,
			"updated_at":    time.Now().UTC(),
		})

	if result.Error != nil {
		return mapPersistenceError(result.Error)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeUserNotFound, ErrUserNotFound)
	}

	return nil
}

// CreateSession saves a new session entry.
func (r *GormRepository) CreateSession(ctx context.Context, session *Session) error {
	if err := session.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return mapPersistenceError(err)
	}
	return nil
}

// UpdateSession persists session changes.
func (r *GormRepository) UpdateSession(ctx context.Context, session *Session) error {
	if err := session.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(session).Error; err != nil {
		return mapPersistenceError(err)
	}

	return nil
}

// FindSessionByID fetches a session by ID.
func (r *GormRepository) FindSessionByID(ctx context.Context, id uuid.UUID) (*Session, error) {
	var session Session
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeSessionRevoked, ErrSessionNotFound)
		}
		return nil, mapPersistenceError(err)
	}

	return &session, nil
}

// FindSessionByRefreshHash fetches session by refresh token hash.
func (r *GormRepository) FindSessionByRefreshHash(ctx context.Context, hash string) (*Session, error) {
	var session Session
	err := r.db.WithContext(ctx).Where("refresh_token_hash = ?", hash).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeSessionRevoked, ErrSessionNotFound)
		}
		return nil, mapPersistenceError(err)
	}

	return &session, nil
}

// ListActiveSessions returns all active sessions for a user.
func (r *GormRepository) ListActiveSessions(ctx context.Context, userID uuid.UUID) ([]Session, error) {
	var sessions []Session
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Find(&sessions).Error
	if err != nil {
		return nil, mapPersistenceError(err)
	}
	return sessions, nil
}

// RevokeSession marks session as revoked.
func (r *GormRepository) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]any{
			"revoked_at": now,
			"updated_at": now,
		})

	if result.Error != nil {
		return mapPersistenceError(result.Error)
	}
	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeSessionRevoked, ErrSessionNotFound)
	}
	return nil
}

// RevokeUserSessions revokes other sessions for safety.
func (r *GormRepository) RevokeUserSessions(ctx context.Context, userID uuid.UUID, exclude uuid.UUID) error {
	now := time.Now().UTC()
	query := r.db.WithContext(ctx).Model(&Session{}).
		Where("user_id = ? AND revoked_at IS NULL", userID)

	if exclude != uuid.Nil {
		query = query.Where("id <> ?", exclude)
	}

	if err := query.Updates(map[string]any{
		"revoked_at": now,
		"updated_at": now,
	}).Error; err != nil {
		return mapPersistenceError(err)
	}
	return nil
}

// UpsertDevice registers or updates device metadata.
func (r *GormRepository) UpsertDevice(ctx context.Context, device *Device) (*Device, error) {
	if device == nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrNilDevice)
	}

	var existing Device
	err := r.db.WithContext(ctx).
		Where("fingerprint = ?", device.Fingerprint).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := r.db.WithContext(ctx).Create(device).Error; err != nil {
				return nil, mapPersistenceError(err)
			}
			return device, nil
		}
		return nil, mapPersistenceError(err)
	}

	existing.Touch(device.Name)
	existing.UserID = device.UserID

	if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return nil, mapPersistenceError(err)
	}

	return &existing, nil
}

// ListDevices returns devices for a given user.
func (r *GormRepository) ListDevices(ctx context.Context, userID uuid.UUID) ([]Device, error) {
	var devices []Device
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&devices).Error; err != nil {
		return nil, mapPersistenceError(err)
	}
	return devices, nil
}

// CreateMFAToken persists a new MFA token.
func (r *GormRepository) CreateMFAToken(ctx context.Context, token *MFAToken) error {
	if token == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilMFAToken)
	}

	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return mapPersistenceError(err)
	}
	return nil
}

// FindActiveMFAToken retrieves the active MFA token.
func (r *GormRepository) FindActiveMFAToken(ctx context.Context, userID uuid.UUID, tokenType MFAType) (*MFAToken, error) {
	var token MFAToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND revoked_at IS NULL", userID, tokenType).
		First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeUserNotFound, ErrMFATokenNotFound)
		}
		return nil, mapPersistenceError(err)
	}
	return &token, nil
}

// UpdateMFAToken saves updates.
func (r *GormRepository) UpdateMFAToken(ctx context.Context, token *MFAToken) error {
	if token == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilMFAToken)
	}
	if err := r.db.WithContext(ctx).Save(token).Error; err != nil {
		return mapPersistenceError(err)
	}
	return nil
}

// DeleteMFAToken removes mfa token.
func (r *GormRepository) DeleteMFAToken(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&MFAToken{}, "id = ?", id).Error; err != nil {
		return mapPersistenceError(err)
	}
	return nil
}

// InsertAuditLog inserts a new audit entry.
func (r *GormRepository) InsertAuditLog(ctx context.Context, entry *AuditLog) error {
	if entry == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilAuditLog)
	}
	if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
		return mapPersistenceError(err)
	}
	return nil
}

// ListAuditLogs fetches audit logs for a user.
func (r *GormRepository) ListAuditLogs(ctx context.Context, userID uuid.UUID, limit int) ([]AuditLog, error) {
	if limit <= 0 {
		limit = 50
	}
	var entries []AuditLog
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&entries).Error
	if err != nil {
		return nil, mapPersistenceError(err)
	}
	return entries, nil
}

// FindOAuthAccount retrieves an OAuth account by provider and subject.
func (r *GormRepository) FindOAuthAccount(ctx context.Context, provider OAuthProvider, providerUserID string) (*OAuthAccount, error) {
	var account OAuthAccount
	err := r.db.WithContext(ctx).
		Where("provider = ? AND provider_user_id = ?", provider, providerUserID).
		First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, mapPersistenceError(err)
	}
	return &account, nil
}

// UpsertOAuthAccount stores external identity binding.
func (r *GormRepository) UpsertOAuthAccount(ctx context.Context, account *OAuthAccount) error {
	if account == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilOAuthAccount)
	}

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ?", account.UserID, account.Provider).
		Assign(map[string]any{
			"provider_user_id": account.ProviderUserID,
			"access_token":     account.AccessToken,
			"refresh_token":    account.RefreshToken,
			"expires_at":       account.ExpiresAt,
			"scopes":           account.Scopes,
			"updated_at":       time.Now().UTC(),
		}).
		FirstOrCreate(account).Error

	if err != nil {
		return mapPersistenceError(err)
	}

	return nil
}

// ListOAuthAccounts retrieves linked OAuth accounts.
func (r *GormRepository) ListOAuthAccounts(ctx context.Context, userID uuid.UUID) ([]OAuthAccount, error) {
	var accounts []OAuthAccount
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, mapPersistenceError(err)
	}
	return accounts, nil
}

// DeleteOAuthAccount removes an OAuth link.
func (r *GormRepository) DeleteOAuthAccount(ctx context.Context, userID uuid.UUID, provider OAuthProvider) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ?", userID, provider).
		Delete(&OAuthAccount{}).Error; err != nil {
		return mapPersistenceError(err)
	}
	return nil
}

// CreateEmailToken persists a new email token entry.
func (r *GormRepository) CreateEmailToken(ctx context.Context, token *EmailToken) error {
	if token == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilEmailToken)
	}

	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return mapPersistenceError(err)
	}
	return nil
}

// FindEmailToken fetches an email token by hash.
func (r *GormRepository) FindEmailToken(ctx context.Context, tokenType EmailTokenType, tokenHash string) (*EmailToken, error) {
	var token EmailToken
	err := r.db.WithContext(ctx).
		Where("type = ? AND token_hash = ?", tokenType, tokenHash).
		First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeInvalidToken, ErrInvalidResetToken)
		}
		return nil, mapPersistenceError(err)
	}
	return &token, nil
}

// UpdateEmailToken updates token counters or consumption.
func (r *GormRepository) UpdateEmailToken(ctx context.Context, token *EmailToken) error {
	if token == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilEmailToken)
	}

	if err := r.db.WithContext(ctx).Save(token).Error; err != nil {
		return mapPersistenceError(err)
	}
	return nil
}

// DeleteEmailTokensByUser deletes tokens of a given type for user.
func (r *GormRepository) DeleteEmailTokensByUser(ctx context.Context, userID uuid.UUID, tokenType EmailTokenType) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, tokenType).
		Delete(&EmailToken{}).Error; err != nil {
		return mapPersistenceError(err)
	}
	return nil
}

// BulkUpdateUserStatus performs mass updates for administrative actions.
func (r *GormRepository) BulkUpdateUserStatus(ctx context.Context, userIDs []uuid.UUID, updates map[string]any) error {
	if len(userIDs) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id IN ?", userIDs).
		Updates(updates)

	if result.Error != nil {
		return mapPersistenceError(result.Error)
	}
	return nil
}

// ListUsers retrieves users with pagination and optional search.
func (r *GormRepository) ListUsers(ctx context.Context, limit, offset int, search string) ([]User, int64, error) {
	var users []User
	var total int64

	query := r.db.WithContext(ctx).Model(&User{})

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(email) LIKE ? OR LOWER(role) LIKE ?", searchPattern, searchPattern)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, mapPersistenceError(err)
	}

	// Apply pagination and fetch users
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, mapPersistenceError(err)
	}

	return users, total, nil
}

func mapPersistenceError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return NewDomainError(ErrCodeEmailAlreadyExists, ErrUserAlreadyExists)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return NewDomainError(ErrCodeEmailAlreadyExists, ErrUserAlreadyExists)
	}

	// TODO: inspect pq error codes for duplicate violations and map precisely.
	return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
}
