package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Context keys for storing authentication metadata.
const (
	ContextUserIDKey    = "auth.user_id"
	ContextUserEmailKey = "auth.user_email"
)

// TokenClaims describes the JWT payload used in Woragis.
type TokenClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// JWTManager handles signing and verifying JWT access tokens.
type JWTManager struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

// NewJWTManager constructs a JWT manager for the provided secret and TTL.
func NewJWTManager(secret string, ttl time.Duration, issuer string) (*JWTManager, error) {
	if secret == "" {
		return nil, fmt.Errorf("jwt manager: secret cannot be empty")
	}

	if ttl <= 0 {
		return nil, fmt.Errorf("jwt manager: ttl must be positive")
	}

	if issuer == "" {
		issuer = "woragis"
	}

	return &JWTManager{
		secret: []byte(secret),
		ttl:    ttl,
		issuer: issuer,
	}, nil
}

// Generate issues a new signed JWT for the supplied user.
func (m *JWTManager) Generate(userID uuid.UUID, email string) (string, error) {
	if userID == uuid.Nil {
		return "", fmt.Errorf("jwt manager: user id cannot be empty")
	}

	now := time.Now().UTC()
	claims := TokenClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Verify parses and validates the provided token string, returning its claims.
func (m *JWTManager) Verify(tokenString string) (*TokenClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("jwt manager: token cannot be empty")
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("jwt manager: unexpected signing method %s", t.Method.Alg())
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*TokenClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, errors.New("jwt manager: invalid token claims")
}
