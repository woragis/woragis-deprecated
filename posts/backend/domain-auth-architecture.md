# Auth Domain Architecture

## Overview
Authentication and authorization architecture.

## Key Points

### Core Components
- Repository: User data access
- Service: Business logic for auth operations
- JWT Manager: Token generation and validation
- Token Store: Refresh token management
- Email Sender: Email notifications
- OAuth Providers: OAuth integration

### Authentication Flows

#### Email/Password Registration
1. User provides email and password
2. Password hashing (bcrypt)
3. Email confirmation token generation
4. User created with pending status
5. Confirmation email sent
6. User confirms email → account activated

#### Email/Password Login
1. User provides email and password
2. Password verification (bcrypt)
3. JWT access token + refresh token generation
4. Refresh token stored (TokenStore)
5. Session created (optional)
6. Tokens returned to client

#### OAuth Login
1. User initiates OAuth flow
2. OAuth state token generated
3. Redirect to OAuth provider
4. Provider callback with code
5. Exchange code for access token
6. Fetch user info from provider
7. Create or link account
8. Generate JWT tokens

#### Password Reset
1. User requests password reset
2. Reset token generated (time-limited)
3. Reset email sent with token
4. User clicks link with token
5. Token validated
6. New password set

### Security Features
- Password hashing: bcrypt with appropriate cost
- Token expiration: Access tokens (short), refresh tokens (long)
- Email confirmation required
- Password strength validation
- Rate limiting (implicit)
- CSRF protection (OAuth state)
- Secure token storage

## Potential Improvements
- Add MFA (2FA) support
- Implement session management
- Add device tracking
- Implement account lockout after failed attempts
- Add password history (prevent reuse)
- Support password expiration
- Add security audit logging
- Implement IP-based rate limiting
- Add suspicious activity detection
- Support social login providers expansion
- Add passwordless authentication (magic links)
- Implement biometric authentication
- Add account recovery options

