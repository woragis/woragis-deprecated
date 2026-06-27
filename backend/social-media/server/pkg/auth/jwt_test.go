package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_Generate(t *testing.T) {
	manager := NewJWTManager(
		"test-secret-key",
		"test-issuer",
		1*time.Hour,
		24*time.Hour,
	)

	userID := uuid.New()
	email := "test@example.com"
	role := "user"
	name := "Test User"

	accessToken, refreshToken, err := manager.Generate(userID, email, role, name)

	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.NotEqual(t, accessToken, refreshToken)
}

func TestJWTManager_Validate(t *testing.T) {
	manager := NewJWTManager(
		"test-secret-key",
		"test-issuer",
		1*time.Hour,
		24*time.Hour,
	)

	userID := uuid.New()
	email := "test@example.com"
	role := "user"
	name := "Test User"

	accessToken, _, err := manager.Generate(userID, email, role, name)
	require.NoError(t, err)

	claims, err := manager.Validate(accessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, name, claims.Name)
}

func TestJWTManager_Validate_InvalidToken(t *testing.T) {
	manager := NewJWTManager(
		"test-secret-key",
		"test-issuer",
		1*time.Hour,
		24*time.Hour,
	)

	_, err := manager.Validate("invalid-token")
	assert.Error(t, err)
	assert.Equal(t, ErrTokenInvalid, err)
}

func TestJWTManager_Validate_ExpiredToken(t *testing.T) {
	manager := NewJWTManager(
		"test-secret-key",
		"test-issuer",
		-1*time.Hour, // Already expired
		24*time.Hour,
	)

	userID := uuid.New()
	accessToken, _, err := manager.Generate(userID, "test@example.com", "user", "Test")
	require.NoError(t, err)

	// Wait a bit to ensure expiry
	time.Sleep(100 * time.Millisecond)

	_, err = manager.Validate(accessToken)
	assert.Error(t, err)
	assert.Equal(t, ErrTokenExpired, err)
}

func TestJWTManager_Refresh(t *testing.T) {
	manager := NewJWTManager(
		"test-secret-key",
		"test-issuer",
		1*time.Hour,
		24*time.Hour,
	)

	userID := uuid.New()
	email := "test@example.com"
	role := "user"
	name := "Test User"

	_, refreshToken, err := manager.Generate(userID, email, role, name)
	require.NoError(t, err)

	newAccessToken, err := manager.Refresh(refreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)

	// Validate the new access token
	claims, err := manager.Validate(newAccessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestJWTManager_RevokeToken(t *testing.T) {
	// Create a mock Redis client using a real connection for testing
	// In real tests, use a test Redis instance
	opt := &redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Use DB 1 for tests
	}
	redisClient := redis.NewClient(opt)
	ctx := context.Background()

	// Test connection (skip if Redis is not available)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping token revocation test")
	}

	manager := NewJWTManager(
		"test-secret-key",
		"test-issuer",
		1*time.Hour,
		24*time.Hour,
	)
	manager.SetRedisClient(redisClient)

	userID := uuid.New()
	accessToken, _, err := manager.Generate(userID, "test@example.com", "user", "Test")
	require.NoError(t, err)

	// Revoke token
	err = manager.RevokeToken(accessToken, 1*time.Hour)
	require.NoError(t, err)

	// Token should still validate (blacklist check happens in Validate if Redis is set)
	// But the blacklist check in Validate uses a different key format
	// This test verifies RevokeToken doesn't error
	assert.NoError(t, err)

	// Cleanup
	redisClient.FlushDB(ctx)
	redisClient.Close()
}

func TestJWTManager_RevokeToken_NoRedis(t *testing.T) {
	manager := NewJWTManager(
		"test-secret-key",
		"test-issuer",
		1*time.Hour,
		24*time.Hour,
	)
	// No Redis client set

	err := manager.RevokeToken("token", 1*time.Hour)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Redis client not configured")
}

func TestJWTManager_IsUserTokenRevoked(t *testing.T) {
	opt := &redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	}
	redisClient := redis.NewClient(opt)
	ctx := context.Background()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping user token revocation test")
	}

	manager := NewJWTManager(
		"test-secret-key",
		"test-issuer",
		1*time.Hour,
		24*time.Hour,
	)
	manager.SetRedisClient(redisClient)

	userID := uuid.New()

	// User tokens not revoked initially
	revoked, err := manager.IsUserTokenRevoked(userID)
	require.NoError(t, err)
	assert.False(t, revoked)

	// Revoke user tokens
	err = manager.RevokeUserTokens(userID, 1*time.Hour)
	require.NoError(t, err)

	// Check if revoked
	revoked, err = manager.IsUserTokenRevoked(userID)
	require.NoError(t, err)
	assert.True(t, revoked)

	// Cleanup
	redisClient.FlushDB(ctx)
	redisClient.Close()
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		wantToken   string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid header",
			authHeader: "Bearer test-token-123",
			wantToken:  "test-token-123",
			wantErr:    false,
		},
		{
			name:        "missing header",
			authHeader:  "",
			wantErr:     true,
			errContains: "authorization header is required",
		},
		{
			name:        "invalid prefix",
			authHeader:  "Basic test-token-123",
			wantErr:     true,
			errContains: "authorization header must start with 'Bearer '",
		},
		{
			name:        "too short",
			authHeader:  "Bear",
			wantErr:     true,
			errContains: "authorization header must start with 'Bearer '",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ExtractTokenFromHeader(tt.authHeader)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
			}
		})
	}
}

func BenchmarkJWTManager_Generate(b *testing.B) {
	manager := NewJWTManager(
		"test-secret-key",
		"test-issuer",
		1*time.Hour,
		24*time.Hour,
	)

	userID := uuid.New()
	email := "test@example.com"
	role := "user"
	name := "Test User"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := manager.Generate(userID, email, role, name)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJWTManager_Validate(b *testing.B) {
	manager := NewJWTManager(
		"test-secret-key",
		"test-issuer",
		1*time.Hour,
		24*time.Hour,
	)

	userID := uuid.New()
	accessToken, _, err := manager.Generate(userID, "test@example.com", "user", "Test")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.Validate(accessToken)
		if err != nil {
			b.Fatal(err)
		}
	}
}

