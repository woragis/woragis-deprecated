package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// TokenStore defines operations for password reset tokens.
type TokenStore interface {
	CreateToken(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error)
	ValidateToken(ctx context.Context, token string) (uuid.UUID, error)
	InvalidateToken(ctx context.Context, token string) error
}

const (
	tokenPrefix = "auth:reset:token:"
	userPrefix  = "auth:reset:user:"
	tokenBytes  = 32
)

// RedisTokenStore implements TokenStore backed by Redis.
type RedisTokenStore struct {
	client *redis.Client
}

// NewRedisTokenStore constructs a Redis-backed token store.
func NewRedisTokenStore(client *redis.Client) *RedisTokenStore {
	return &RedisTokenStore{client: client}
}

// CreateToken issues a new token for a user, invalidating previous ones.
func (r *RedisTokenStore) CreateToken(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error) {
	if r.client == nil {
		return generateToken()
	}

	existing, err := r.client.Get(ctx, userPrefix+userID.String()).Result()
	if err == nil && existing != "" {
		_ = r.client.Del(ctx, tokenPrefix+existing).Err()
	}

	token, err := generateToken()
	if err != nil {
		return "", err
	}

	tokenKey := tokenPrefix + token
	userKey := userPrefix + userID.String()

	if err := r.client.Set(ctx, tokenKey, userID.String(), ttl).Err(); err != nil {
		return "", err
	}

	if err := r.client.Set(ctx, userKey, token, ttl).Err(); err != nil {
		_ = r.client.Del(ctx, tokenKey).Err()
		return "", err
	}

	return token, nil
}

// ValidateToken returns the associated user ID when token is valid.
func (r *RedisTokenStore) ValidateToken(ctx context.Context, token string) (uuid.UUID, error) {
	if token == "" {
		return uuid.Nil, fmt.Errorf("token required")
	}

	if r.client == nil {
		return uuid.Nil, fmt.Errorf("token store not configured")
	}

	val, err := r.client.Get(ctx, tokenPrefix+token).Result()
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

// InvalidateToken removes token entries from Redis.
func (r *RedisTokenStore) InvalidateToken(ctx context.Context, token string) error {
	if token == "" || r.client == nil {
		return nil
	}

	val, err := r.client.Get(ctx, tokenPrefix+token).Result()
	if err == nil && val != "" {
		userKey := userPrefix + val
		_ = r.client.Del(ctx, userKey).Err()
	}

	return r.client.Del(ctx, tokenPrefix+token).Err()
}

func generateToken() (string, error) {
	buf := make([]byte, tokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}
