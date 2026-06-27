# Auth Domain

The Auth domain handles user authentication, authorization, and profile management for the Woragis application. It provides a complete authentication system with JWT tokens, session management, email verification, and user profile functionality.

## Overview

This domain implements a secure authentication system with the following key features:

- **User Registration & Login**: Secure user registration with email validation and password hashing
- **JWT Token Management**: Access and refresh token generation with automatic refresh
- **Session Management**: Multi-device session tracking with logout capabilities
- **Email Verification**: Token-based email verification system
- **Profile Management**: User profile creation and updates
- **Password Management**: Secure password changes with strength validation
- **Admin Functions**: User management and system cleanup for administrators

## Architecture

The auth domain follows a clean architecture pattern with clear separation of concerns:

```
auth/
├── entity.go      # Data models and business logic
├── repository.go  # Data access layer
├── service.go     # Business logic layer
├── handlers.go    # HTTP request handlers
├── routes.go      # Route definitions
├── migration.go   # Database migrations
└── README.md      # This documentation
```

## Data Models

### User

Core user entity with authentication and basic profile information.

```go
type User struct {
    ID         uuid.UUID      `json:"id"`
    Email      string         `json:"email"`
    Password   string         `json:"-"` // Hidden from JSON
    FirstName  string         `json:"first_name"`
    LastName   string         `json:"last_name"`
    Role       string         `json:"role"` // user, moderator, admin
    IsActive   bool           `json:"is_active"`
    IsVerified bool           `json:"is_verified"`
    LastLogin  *time.Time     `json:"last_login"`
    CreatedAt  time.Time      `json:"created_at"`
    UpdatedAt  time.Time      `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `json:"-"`
}
```

### Profile

Extended user profile information.

```go
type Profile struct {
    ID          uuid.UUID      `json:"id"`
    UserID      uuid.UUID      `json:"user_id"`
    Avatar      string         `json:"avatar"`
    Bio         string         `json:"bio"`
    DateOfBirth *time.Time     `json:"date_of_birth"`
    Gender      string         `json:"gender"`
    Phone       string         `json:"phone"`
    Location    string         `json:"location"`
    Website     string         `json:"website"`
    SocialLinks string         `json:"social_links"` // JSON string
    Preferences string         `json:"preferences"`  // JSON string
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-"`
}
```

### Session

User session management for multi-device support.

```go
type Session struct {
    ID           uuid.UUID      `json:"id"`
    UserID       uuid.UUID      `json:"user_id"`
    RefreshToken string         `json:"-"`
    UserAgent    string         `json:"user_agent"`
    IPAddress    string         `json:"ip_address"`
    IsActive     bool           `json:"is_active"`
    ExpiresAt    time.Time      `json:"expires_at"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-"`
}
```

### VerificationToken

Email verification and password reset tokens.

```go
type VerificationToken struct {
    ID        uuid.UUID      `json:"id"`
    UserID    uuid.UUID      `json:"user_id"`
    Token     string         `json:"-"`
    Type      string         `json:"type"` // "email_verification", "password_reset"
    ExpiresAt time.Time      `json:"expires_at"`
    IsUsed    bool           `json:"is_used"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-"`
}
```

## API Endpoints

### Public Endpoints

#### POST /auth/register

Register a new user account.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      /* User object */
    },
    "access_token": "jwt_token_here",
    "refresh_token": "refresh_token_here",
    "expires_at": 1640995200
  }
}
```

#### POST /auth/login

Authenticate user and return tokens.

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      /* User object */
    },
    "access_token": "jwt_token_here",
    "refresh_token": "refresh_token_here",
    "expires_at": 1640995200
  }
}
```

#### POST /auth/refresh

Refresh access token using refresh token.

**Request Body:**

```json
{
  "refresh_token": "refresh_token_here"
}
```

#### POST /auth/logout

Logout user and invalidate session.

**Request Body:**

```json
{
  "refresh_token": "refresh_token_here"
}
```

#### GET /auth/verify-email?token=verification_token

Verify user email using verification token.

### Protected Endpoints (Require Authentication)

#### GET /auth/profile

Get current user's profile information.

**Headers:**

```
Authorization: Bearer <access_token>
```

#### PUT /auth/profile

Update current user's profile.

**Request Body:**

```json
{
  "avatar": "https://example.com/avatar.jpg",
  "bio": "User bio here",
  "date_of_birth": "1990-01-01T00:00:00Z",
  "gender": "male",
  "phone": "+1234567890",
  "location": "New York, NY",
  "website": "https://example.com",
  "social_links": "{\"twitter\": \"@username\"}",
  "preferences": "{\"theme\": \"dark\"}"
}
```

#### POST /auth/change-password

Change user password.

**Request Body:**

```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword123"
}
```

#### POST /auth/logout-all

Logout from all devices.

### Admin Endpoints (Require Admin Role)

#### GET /auth/admin/users/:id

Get user by ID (admin only).

#### GET /auth/admin/users

List all users with pagination (admin only).

**Query Parameters:**

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `search`: Search query for name or email

#### POST /auth/admin/cleanup

Cleanup expired sessions and tokens (admin only).

## Security Features

### Password Security

- **Bcrypt Hashing**: Passwords are hashed using bcrypt with configurable cost
- **Strength Validation**: Password strength requirements enforced
- **Secure Storage**: Passwords are never returned in API responses

### Token Security

- **JWT Tokens**: Secure access tokens with configurable expiration
- **Refresh Tokens**: Long-lived refresh tokens for seamless re-authentication
- **Session Tracking**: Multi-device session management with device information
- **Token Rotation**: Refresh tokens can be rotated for enhanced security

### Session Management

- **Multi-Device Support**: Users can be logged in on multiple devices
- **Session Invalidation**: Individual or all-device logout capabilities
- **Automatic Cleanup**: Expired sessions are automatically cleaned up
- **Device Tracking**: User agent and IP address tracking for security

### Email Verification

- **Token-Based**: Secure email verification using random tokens
- **Time-Limited**: Verification tokens expire after 24 hours
- **One-Time Use**: Tokens are invalidated after successful verification

## Error Handling

The auth domain provides comprehensive error handling with specific error types:

```go
var (
    ErrInvalidPassword     = errors.New("invalid password")
    ErrUserNotVerified     = errors.New("user email not verified")
    ErrUserInactive        = errors.New("user account is inactive")
    ErrSessionExpired      = errors.New("session has expired")
    ErrTokenExpired        = errors.New("verification token has expired")
    ErrTokenAlreadyUsed    = errors.New("verification token already used")
    ErrPasswordTooWeak     = errors.New("password does not meet strength requirements")
    ErrUserNotFound        = errors.New("user not found")
    ErrProfileNotFound     = errors.New("profile not found")
    ErrSessionNotFound     = errors.New("session not found")
    ErrTokenNotFound       = errors.New("verification token not found")
    ErrUserAlreadyExists   = errors.New("user already exists")
    ErrInvalidCredentials  = errors.New("invalid credentials")
)
```

## Database Schema

### Tables Created

- `users`: Core user information
- `profiles`: Extended user profile data
- `sessions`: User authentication sessions
- `verification_tokens`: Email verification and password reset tokens

### Indexes

The migration creates optimized indexes for:

- User email lookups
- Session token lookups
- User role filtering
- Expiration date queries
- User ID foreign key relationships

## Usage Example

### Initialization

```go
// Initialize auth handler
authHandler := auth.NewHandler(db, jwtManager, bcryptCost)

// Setup routes
authHandler.SetupRoutes(api)

// Run migrations
err := auth.MigrateAuthTables(db)
```

### Service Usage

```go
// Register a new user
registerReq := &auth.RegisterRequest{
    Email:     "user@example.com",
    Password:  "securepassword123",
    FirstName: "John",
    LastName:  "Doe",
}

response, err := authService.Register(registerReq)
```

## Dependencies

- **GORM**: Database ORM for PostgreSQL
- **Fiber**: HTTP web framework
- **JWT**: JSON Web Token implementation
- **UUID**: UUID generation and handling
- **Bcrypt**: Password hashing
- **Crypto**: Random string generation for tokens

## Configuration

The auth domain requires the following configuration:

- **Database Connection**: GORM database instance
- **JWT Manager**: Configured JWT manager with secret and expiration
- **Bcrypt Cost**: Password hashing cost (recommended: 12-14)

## Best Practices

1. **Password Security**: Always use strong passwords and proper hashing
2. **Token Management**: Implement proper token rotation and expiration
3. **Session Cleanup**: Regularly clean up expired sessions and tokens
4. **Input Validation**: Validate all user inputs before processing
5. **Error Handling**: Provide meaningful error messages without exposing sensitive information
6. **Rate Limiting**: Implement rate limiting for authentication endpoints
7. **Logging**: Log authentication events for security monitoring

## Future Enhancements

- **Two-Factor Authentication**: SMS or TOTP-based 2FA
- **OAuth Integration**: Social login providers
- **Password Reset**: Email-based password reset functionality
- **Account Lockout**: Brute force protection
- **Audit Logging**: Comprehensive authentication audit trail
- **Role-Based Permissions**: Granular permission system
