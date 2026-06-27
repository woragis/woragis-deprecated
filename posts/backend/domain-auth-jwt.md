# Auth Domain - JWT Management

## Overview
How JWT tokens are generated, verified, and managed for authentication.

## Key Points

### JWT Manager Architecture
- JWTManager handles signing and verification
- HS256 signing method (HMAC-SHA256)
- Configurable secret, TTL, and issuer
- TokenClaims with email and standard claims

### Token Generation
- Generate method creates signed JWTs
- Includes: issuer, subject (user_id), email, issued_at, not_before, expires_at
- TTL configurable (default from service)
- Validation: user_id cannot be empty

### Token Verification
- Verify method parses and validates tokens
- Validates signing method (must be HS256)
- Checks token signature with secret
- Returns TokenClaims if valid

### Token Claims Structure
- Email: User email address
- RegisteredClaims: Standard JWT claims (issuer, subject, issued_at, expires_at, not_before)
- Subject: User ID (UUID string)

### Context Integration
- ContextUserIDKey: Stores user_id in context
- ContextUserEmailKey: Stores user_email in context
- Used by middleware for authentication

### Security Features
- Secret validation (cannot be empty)
- TTL validation (must be positive)
- Token expiration checking
- Signature verification

## Potential Improvements
- Support multiple signing algorithms (RS256, ES256)
- Add token refresh mechanism
- Implement token revocation list
- Add token versioning
- Support token rotation
- Add token claims customization
- Implement token encryption (JWE)
- Add token audience validation
- Support multiple issuers
- Add token usage tracking
- Implement token compression
- Support token bundling (access + refresh)

