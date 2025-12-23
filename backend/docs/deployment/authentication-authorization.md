# Authentication & Authorization Guide

**Last Updated:** 2025-12-22  
**Purpose:** Guide for authentication and authorization in production

---

## Overview

This guide covers authentication and authorization best practices for the Woragis backend services.

---

## Current State

**Authentication:** Basic session-based and token-based authentication  
**Authorization:** Role-based access control (if implemented)

---

## Authentication Methods

### 1. API Key Authentication

**Use Case:** Service-to-service communication, public API access

**Implementation:**
```go
// Go example
func APIKeyMiddleware(c *fiber.Ctx) error {
    apiKey := c.Get("X-API-Key")
    if apiKey == "" {
        return c.Status(401).JSON(fiber.Map{"error": "API key required"})
    }
    
    // Validate API key
    if !isValidAPIKey(apiKey) {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid API key"})
    }
    
    return c.Next()
}
```

**Best Practices:**
- Store API keys securely (hashed in database)
- Rotate keys regularly
- Use different keys per environment
- Implement rate limiting per key
- Log API key usage

---

### 2. JWT (JSON Web Tokens)

**Use Case:** User authentication, stateless sessions

**Implementation:**
```go
// Generate JWT
func GenerateToken(userID string, expiresIn time.Duration) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(expiresIn).Unix(),
        "iat":     time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Validate JWT
func ValidateToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
}
```

**Best Practices:**
- Use short expiration times (15-60 minutes)
- Implement refresh tokens for long-lived sessions
- Store JWT secret securely
- Use HTTPS only
- Include user ID and roles in claims
- Validate token on every request

---

### 3. Session-Based Authentication

**Use Case:** Web applications, stateful sessions

**Implementation:**
```go
// Create session
func CreateSession(userID string) (string, error) {
    sessionID := uuid.New().String()
    
    // Store in Redis with expiration
    err := redis.Set(ctx, fmt.Sprintf("session:%s", sessionID), userID, 24*time.Hour).Err()
    if err != nil {
        return "", err
    }
    
    return sessionID, nil
}

// Validate session
func ValidateSession(sessionID string) (string, error) {
    userID, err := redis.Get(ctx, fmt.Sprintf("session:%s", sessionID)).Result()
    if err == redis.Nil {
        return "", errors.New("session not found")
    }
    return userID, err
}
```

**Best Practices:**
- Use secure, HTTP-only cookies
- Set appropriate expiration times
- Implement session rotation
- Store sessions in Redis (not in-memory)
- Invalidate sessions on logout
- Implement session timeout

---

## Authorization (RBAC)

### Role-Based Access Control

**Roles:**
- `admin` - Full system access
- `user` - Standard user access
- `service` - Service-to-service access
- `readonly` - Read-only access

**Implementation:**
```go
// Role-based middleware
func RequireRole(roles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userRole := c.Locals("user_role")
        if userRole == nil {
            return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
        }
        
        for _, role := range roles {
            if userRole == role {
                return c.Next()
            }
        }
        
        return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
    }
}

// Usage
app.Get("/admin/users", RequireRole("admin"), GetUsers)
```

---

## Security Best Practices

### 1. Password Security

```go
// Hash passwords with bcrypt
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hash), err
}

func VerifyPassword(hashedPassword, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
```

**Requirements:**
- Minimum 12 characters
- Require uppercase, lowercase, numbers, symbols
- Use bcrypt with cost factor 10+
- Never store plaintext passwords
- Implement password reset with secure tokens

---

### 2. Rate Limiting

**Implementation:**
```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

// Global rate limiter
app.Use(limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP()
    },
}))

// Per-user rate limiter
app.Use(limiter.New(limiter.Config{
    Max:        20,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        userID := c.Locals("user_id")
        if userID == nil {
            return c.IP()
        }
        return fmt.Sprintf("user:%v", userID)
    },
}))
```

**Rate Limits:**
- Login attempts: 5 per 15 minutes
- API requests: 100 per minute (per IP)
- User requests: 1000 per hour (per user)
- Password reset: 3 per hour

---

### 3. Token Expiration

**Short-lived Access Tokens:**
- Expiration: 15-60 minutes
- Stored in memory only
- Cannot be revoked (stateless)

**Long-lived Refresh Tokens:**
- Expiration: 7-30 days
- Stored in database/Redis
- Can be revoked
- Used to obtain new access tokens

**Implementation:**
```go
// Refresh token flow
func RefreshToken(c *fiber.Ctx) error {
    refreshToken := c.Get("X-Refresh-Token")
    
    // Validate refresh token
    token, err := ValidateRefreshToken(refreshToken)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid refresh token"})
    }
    
    // Check if token is revoked
    if isTokenRevoked(refreshToken) {
        return c.Status(401).JSON(fiber.Map{"error": "Token revoked"})
    }
    
    // Generate new access token
    userID := token.Claims.(jwt.MapClaims)["user_id"].(string)
    accessToken, err := GenerateToken(userID, 15*time.Minute)
    if err != nil {
        return err
    }
    
    return c.JSON(fiber.Map{"access_token": accessToken})
}
```

---

### 4. Secure Headers

```go
import "github.com/gofiber/fiber/v2/middleware/helmet"

app.Use(helmet.New())
```

**Headers to Set:**
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security: max-age=31536000`
- `Content-Security-Policy: default-src 'self'`

---

## Implementation Checklist

### Phase 1: Basic Security (Week 1)
- [ ] Implement password hashing (bcrypt)
- [ ] Add rate limiting
- [ ] Set secure headers
- [ ] Implement token expiration
- [ ] Add input validation

### Phase 2: Advanced Features (Week 2)
- [ ] Implement refresh tokens
- [ ] Add role-based access control
- [ ] Implement session management
- [ ] Add API key management
- [ ] Set up audit logging

### Phase 3: Hardening (Week 3)
- [ ] Implement token rotation
- [ ] Add account lockout after failed attempts
- [ ] Implement password complexity requirements
- [ ] Add two-factor authentication (optional)
- [ ] Set up security monitoring

---

## Related Documentation

- [Secrets Management](./secrets-management.md)
- [SSL/TLS Configuration](./ssl-tls-configuration.md)
- [Input Validation](./input-validation.md) (when created)
- [Security Guide](../operations/security.md) (when created)

---

**Next Steps:**
1. Review current authentication implementation
2. Implement password hashing if not present
3. Add rate limiting
4. Set up token expiration
5. Implement RBAC if needed

---

**Last Updated:** 2025-12-22
