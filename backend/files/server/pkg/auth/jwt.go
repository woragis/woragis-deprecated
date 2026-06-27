package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrTokenInvalid = errors.New("invalid or expired token")
	ErrTokenExpired = errors.New("token has expired")
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Name   string    `json:"name"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey        string
	issuer           string
	accessExpiry     time.Duration
	refreshExpiry    time.Duration
	redisClient      *redis.Client // Optional: for token blacklist
}

func NewJWTManager(secret string, issuer string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secret,
		issuer:        issuer,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// SetRedisClient sets the Redis client for token blacklist functionality
func (j *JWTManager) SetRedisClient(client *redis.Client) {
	j.redisClient = client
}

// Generate creates both access and refresh tokens
func (j *JWTManager) Generate(userID uuid.UUID, email string, role string, name string) (string, string, error) {
	now := time.Now()
	
	// Access token claims
	accessClaims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   userID.String(),
		},
	}

	// Refresh token claims (longer expiry, no role/name for security)
	refreshClaims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   userID.String(),
		},
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// Validate validates a token and returns the claims
func (j *JWTManager) Validate(tokenString string) (*Claims, error) {
	// Check if token is blacklisted (if Redis is available)
	if j.redisClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		
		blacklisted, err := j.redisClient.Exists(ctx, fmt.Sprintf("token:blacklist:%s", tokenString)).Result()
		if err == nil && blacklisted > 0 {
			return nil, ErrTokenInvalid
		}
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})
	
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}
	
	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// Refresh generates a new access token from a refresh token
func (j *JWTManager) Refresh(refreshTokenString string) (string, error) {
	claims, err := j.Validate(refreshTokenString)
	if err != nil {
		return "", err
	}

	// Generate new access token with same user info
	now := time.Now()
	accessClaims := Claims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
		Name:   claims.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   claims.UserID.String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

// RevokeToken adds a token to the blacklist until its expiry
func (j *JWTManager) RevokeToken(tokenString string, expiry time.Duration) error {
	if j.redisClient == nil {
		return errors.New("Redis client not configured for token revocation")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := fmt.Sprintf("token:blacklist:%s", tokenString)
	return j.redisClient.Set(ctx, key, "1", expiry).Err()
}

// RevokeUserTokens revokes all tokens for a user (useful for logout-all)
func (j *JWTManager) RevokeUserTokens(userID uuid.UUID, expiry time.Duration) error {
	if j.redisClient == nil {
		return errors.New("Redis client not configured for token revocation")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Store a pattern key that can be checked during validation
	key := fmt.Sprintf("user:blacklist:%s", userID.String())
	return j.redisClient.Set(ctx, key, "1", expiry).Err()
}

// IsUserTokenRevoked checks if a user's tokens are revoked
func (j *JWTManager) IsUserTokenRevoked(userID uuid.UUID) (bool, error) {
	if j.redisClient == nil {
		return false, nil // If no Redis, assume not revoked
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := fmt.Sprintf("user:blacklist:%s", userID.String())
	exists, err := j.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// ExtractTokenFromHeader extracts token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	return authHeader[len(bearerPrefix):], nil
}
