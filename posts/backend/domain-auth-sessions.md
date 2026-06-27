# Auth Domain - Session Management

## Overview
Session management using refresh tokens and device tracking.

## Key Points

### Session Structure
- Session linked to user and device
- Refresh token hash (for token validation)
- Device tracking: DeviceID, UserAgent, IP, Location
- Expiration and revocation support
- Last seen timestamp

### Session Lifecycle
1. User logs in
2. Device ID generated/retrieved
3. Refresh token generated and hashed
4. Session created with device info
5. Session stored in database
6. Session tracked for activity
7. Session revoked on logout or expiration

### Device Tracking
- DeviceID: Unique identifier per device
- UserAgent: Browser/client information
- IP: Client IP address (validated)
- Location: JSON map for geographic data
- LastSeenAt: Last activity timestamp

### Refresh Token Management
- Refresh token hashed before storage
- Unique index on refresh_token_hash
- Token validation via hash lookup
- Token expiration handling

### Session Operations
- Create: New session on login
- Update: Last seen timestamp
- Revoke: Mark session as revoked
- List: User's active sessions
- Delete: Remove session

### Security Features
- IP address validation
- Token hashing (prevents plaintext storage)
- Expiration tracking
- Revocation support
- Device-based tracking

## Potential Improvements
- Add session timeout configuration
- Implement session activity monitoring
- Add suspicious session detection
- Support session limit per user
- Add session geolocation tracking
- Implement session fingerprinting
- Add session migration (device change)
- Support session sharing controls
- Add session analytics
- Implement session notification (new device login)
- Add session remote logout
- Support session encryption at rest

