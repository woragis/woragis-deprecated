package skills

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	appmetrics "woragis-skills-service/pkg/metrics"
	apptracing "woragis-skills-service/pkg/tracing"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrProfileNotFound   = errors.New("profile not found")
	ErrSessionNotFound   = errors.New("session not found")
	ErrTokenNotFound     = errors.New("verification token not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Repository defines the interface for auth repository operations
type Repository interface {
	createUser(user *User) error
	getUserByID(id uuid.UUID) (*User, error)
	getUserByEmail(email string) (*User, error)
	getUserByUsername(username string) (*User, error)
	updateUser(user *User) error
	deleteUser(id uuid.UUID) error
	updateLastLogin(id uuid.UUID) error
	verifyUserEmail(id uuid.UUID) error
	listUsers(offset, limit int, search string) ([]User, int64, error)
	createProfile(profile *Profile) error
	getProfileByUserID(userID uuid.UUID) (*Profile, error)
	updateProfile(profile *Profile) error
	deleteProfile(userID uuid.UUID) error
	createSession(session *Session) error
	getSessionByRefreshToken(refreshToken string) (*Session, error)
	getSessionsByUserID(userID uuid.UUID) ([]Session, error)
	updateSession(session *Session) error
	deactivateSession(id uuid.UUID) error
	deactivateAllUserSessions(userID uuid.UUID) error
	deleteExpiredSessions() error
	createVerificationToken(token *VerificationToken) error
	getVerificationToken(token string) (*VerificationToken, error)
	markTokenAsUsed(id uuid.UUID) error
	deleteExpiredTokens() error
	deleteUsedTokens() error
}

// repositoryImpl implements the Repository interface
type repositoryImpl struct {
	db *gorm.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

// User operations

// createUser creates a new user
func (r *repositoryImpl) createUser(user *User) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "create", "users", func() error {
		start := time.Now()
		err := r.db.Create(user).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("create", "users", duration)
		
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrUserAlreadyExists
			}
			return fmt.Errorf("failed to create user: %w", err)
		}
		return nil
	})
}

// getUserByID retrieves a user by ID
func (r *repositoryImpl) getUserByID(id uuid.UUID) (*User, error) {
	ctx := context.Background()
	var user User
	var err error
	
	err = apptracing.WithDatabaseSpan(ctx, "select", "users", func() error {
		start := time.Now()
		err := r.db.Preload("Profile").First(&user, "id = ?", id).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("select", "users", duration)
		
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return fmt.Errorf("failed to get user by ID: %w", err)
		}
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// getUserByEmail retrieves a user by email
func (r *repositoryImpl) getUserByEmail(email string) (*User, error) {
	ctx := context.Background()
	var user User
	var err error
	
	err = apptracing.WithDatabaseSpan(ctx, "select", "users", func() error {
		start := time.Now()
		err := r.db.Preload("Profile").First(&user, "email = ?", email).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("select", "users", duration)
		
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return fmt.Errorf("failed to get user by email: %w", err)
		}
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// getUserByUsername retrieves a user by username
func (r *repositoryImpl) getUserByUsername(username string) (*User, error) {
	ctx := context.Background()
	var user User
	var err error
	
	err = apptracing.WithDatabaseSpan(ctx, "select", "users", func() error {
		start := time.Now()
		err := r.db.Preload("Profile").First(&user, "username = ?", username).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("select", "users", duration)
		
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return fmt.Errorf("failed to get user by username: %w", err)
		}
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// updateUser updates a user
func (r *repositoryImpl) updateUser(user *User) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "update", "users", func() error {
		start := time.Now()
		err := r.db.Save(user).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("update", "users", duration)
		
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
		return nil
	})
}

// deleteUser soft deletes a user
func (r *repositoryImpl) deleteUser(id uuid.UUID) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "delete", "users", func() error {
		start := time.Now()
		err := r.db.Delete(&User{}, "id = ?", id).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("delete", "users", duration)
		
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
		return nil
	})
}

// updateLastLogin updates the user's last login time
func (r *repositoryImpl) updateLastLogin(id uuid.UUID) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "update", "users", func() error {
		now := time.Now()
		start := time.Now()
		err := r.db.Model(&User{}).Where("id = ?", id).Update("last_login", now).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("update", "users", duration)
		
		if err != nil {
			return fmt.Errorf("failed to update last login: %w", err)
		}
		return nil
	})
}

// verifyUserEmail marks user as verified
func (r *repositoryImpl) verifyUserEmail(id uuid.UUID) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "update", "users", func() error {
		start := time.Now()
		err := r.db.Model(&User{}).Where("id = ?", id).Update("is_verified", true).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("update", "users", duration)
		
		if err != nil {
			return fmt.Errorf("failed to verify user email: %w", err)
		}
		return nil
	})
}

// listUsers retrieves users with pagination
func (r *repositoryImpl) listUsers(offset, limit int, search string) ([]User, int64, error) {
	ctx := context.Background()
	var users []User
	var total int64

	query := r.db.Model(&User{})

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?", 
			searchPattern, searchPattern, searchPattern)
	}

	var countErr, findErr error
	
	// Get total count with tracing and metrics
	countErr = apptracing.WithDatabaseSpan(ctx, "count", "users", func() error {
		start := time.Now()
		err := query.Count(&total).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("count", "users", duration)
		
		if err != nil {
			return fmt.Errorf("failed to count users: %w", err)
		}
		return nil
	})
	
	if countErr != nil {
		return nil, 0, countErr
	}

	// Get users with pagination with tracing and metrics
	findErr = apptracing.WithDatabaseSpan(ctx, "select", "users", func() error {
		start := time.Now()
		err := query.Preload("Profile").Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("select", "users", duration)
		
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}
		return nil
	})
	
	if findErr != nil {
		return nil, 0, findErr
	}

	return users, total, nil
}

// Profile operations

// createProfile creates a new profile
func (r *repositoryImpl) createProfile(profile *Profile) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "create", "profiles", func() error {
		start := time.Now()
		err := r.db.Create(profile).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("create", "profiles", duration)
		
		if err != nil {
			return fmt.Errorf("failed to create profile: %w", err)
		}
		return nil
	})
}

// getProfileByUserID retrieves a profile by user ID
func (r *repositoryImpl) getProfileByUserID(userID uuid.UUID) (*Profile, error) {
	ctx := context.Background()
	var profile Profile
	var err error
	
	err = apptracing.WithDatabaseSpan(ctx, "select", "profiles", func() error {
		start := time.Now()
		err := r.db.Preload("User").First(&profile, "user_id = ?", userID).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("select", "profiles", duration)
		
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrProfileNotFound
			}
			return fmt.Errorf("failed to get profile by user ID: %w", err)
		}
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// updateProfile updates a profile
func (r *repositoryImpl) updateProfile(profile *Profile) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "update", "profiles", func() error {
		start := time.Now()
		err := r.db.Save(profile).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("update", "profiles", duration)
		
		if err != nil {
			return fmt.Errorf("failed to update profile: %w", err)
		}
		return nil
	})
}

// deleteProfile soft deletes a profile
func (r *repositoryImpl) deleteProfile(userID uuid.UUID) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "delete", "profiles", func() error {
		start := time.Now()
		err := r.db.Delete(&Profile{}, "user_id = ?", userID).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("delete", "profiles", duration)
		
		if err != nil {
			return fmt.Errorf("failed to delete profile: %w", err)
		}
		return nil
	})
}

// Session operations

// createSession creates a new session
func (r *repositoryImpl) createSession(session *Session) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "create", "sessions", func() error {
		start := time.Now()
		err := r.db.Create(session).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("create", "sessions", duration)
		
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
		return nil
	})
}

// getSessionByRefreshToken retrieves a session by refresh token
func (r *repositoryImpl) getSessionByRefreshToken(refreshToken string) (*Session, error) {
	ctx := context.Background()
	var session Session
	var err error
	
	err = apptracing.WithDatabaseSpan(ctx, "select", "sessions", func() error {
		start := time.Now()
		err := r.db.Preload("User").First(&session, "refresh_token = ? AND is_active = ?", refreshToken, true).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("select", "sessions", duration)
		
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrSessionNotFound
			}
			return fmt.Errorf("failed to get session by refresh token: %w", err)
		}
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// getSessionsByUserID retrieves all active sessions for a user
func (r *repositoryImpl) getSessionsByUserID(userID uuid.UUID) ([]Session, error) {
	ctx := context.Background()
	var sessions []Session
	var err error
	
	err = apptracing.WithDatabaseSpan(ctx, "select", "sessions", func() error {
		start := time.Now()
		err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&sessions).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("select", "sessions", duration)
		
		if err != nil {
			return fmt.Errorf("failed to get sessions by user ID: %w", err)
		}
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// updateSession updates a session
func (r *repositoryImpl) updateSession(session *Session) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "update", "sessions", func() error {
		start := time.Now()
		err := r.db.Save(session).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("update", "sessions", duration)
		
		if err != nil {
			return fmt.Errorf("failed to update session: %w", err)
		}
		return nil
	})
}

// deactivateSession deactivates a session
func (r *repositoryImpl) deactivateSession(id uuid.UUID) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "update", "sessions", func() error {
		start := time.Now()
		err := r.db.Model(&Session{}).Where("id = ?", id).Update("is_active", false).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("update", "sessions", duration)
		
		if err != nil {
			return fmt.Errorf("failed to deactivate session: %w", err)
		}
		return nil
	})
}

// deactivateAllUserSessions deactivates all sessions for a user
func (r *repositoryImpl) deactivateAllUserSessions(userID uuid.UUID) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "update", "sessions", func() error {
		start := time.Now()
		err := r.db.Model(&Session{}).Where("user_id = ?", userID).Update("is_active", false).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("update", "sessions", duration)
		
		if err != nil {
			return fmt.Errorf("failed to deactivate all user sessions: %w", err)
		}
		return nil
	})
}

// deleteExpiredSessions deletes expired sessions
func (r *repositoryImpl) deleteExpiredSessions() error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "delete", "sessions", func() error {
		start := time.Now()
		err := r.db.Where("expires_at < ?", time.Now()).Delete(&Session{}).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("delete", "sessions", duration)
		
		if err != nil {
			return fmt.Errorf("failed to delete expired sessions: %w", err)
		}
		return nil
	})
}

// VerificationToken operations

// createVerificationToken creates a new verification token
func (r *repositoryImpl) createVerificationToken(token *VerificationToken) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "create", "verification_tokens", func() error {
		start := time.Now()
		err := r.db.Create(token).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("create", "verification_tokens", duration)
		
		if err != nil {
			return fmt.Errorf("failed to create verification token: %w", err)
		}
		return nil
	})
}

// getVerificationToken retrieves a verification token
func (r *repositoryImpl) getVerificationToken(token string) (*VerificationToken, error) {
	ctx := context.Background()
	var verificationToken VerificationToken
	var err error
	
	err = apptracing.WithDatabaseSpan(ctx, "select", "verification_tokens", func() error {
		start := time.Now()
		err := r.db.Preload("User").First(&verificationToken, "token = ?", token).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("select", "verification_tokens", duration)
		
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTokenNotFound
			}
			return fmt.Errorf("failed to get verification token: %w", err)
		}
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	return &verificationToken, nil
}

// markTokenAsUsed marks a verification token as used
func (r *repositoryImpl) markTokenAsUsed(id uuid.UUID) error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "update", "verification_tokens", func() error {
		start := time.Now()
		err := r.db.Model(&VerificationToken{}).Where("id = ?", id).Update("is_used", true).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("update", "verification_tokens", duration)
		
		if err != nil {
			return fmt.Errorf("failed to mark token as used: %w", err)
		}
		return nil
	})
}

// deleteExpiredTokens deletes expired verification tokens
func (r *repositoryImpl) deleteExpiredTokens() error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "delete", "verification_tokens", func() error {
		start := time.Now()
		err := r.db.Where("expires_at < ?", time.Now()).Delete(&VerificationToken{}).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("delete", "verification_tokens", duration)
		
		if err != nil {
			return fmt.Errorf("failed to delete expired tokens: %w", err)
		}
		return nil
	})
}

// deleteUsedTokens deletes used verification tokens
func (r *repositoryImpl) deleteUsedTokens() error {
	ctx := context.Background()
	return apptracing.WithDatabaseSpan(ctx, "delete", "verification_tokens", func() error {
		start := time.Now()
		err := r.db.Where("is_used = ?", true).Delete(&VerificationToken{}).Error
		duration := time.Since(start).Seconds()
		
		appmetrics.RecordDatabaseQuery("delete", "verification_tokens", duration)
		
		if err != nil {
			return fmt.Errorf("failed to delete used tokens: %w", err)
		}
		return nil
	})
}
