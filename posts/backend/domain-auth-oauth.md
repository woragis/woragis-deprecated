# Auth Domain - OAuth Integration

## Overview
How OAuth authentication is integrated with multiple providers.

## Key Points

### Supported Providers
- Google OAuth
- GitHub OAuth
- Microsoft OAuth
- Extensible provider system

### OAuth Account Structure
- Links external provider identity to user account
- Stores: provider, provider_user_id, access_token, refresh_token, expires_at, scopes
- Unique index on provider + provider_user_id
- User can have multiple OAuth accounts

### OAuth Flow
1. User initiates OAuth login
2. OAuth state token generated (CSRF protection)
3. Redirect to provider authorization URL
4. Provider callback with authorization code
5. Exchange code for access/refresh tokens
6. Fetch user info from provider
7. Create or link OAuth account
8. Generate JWT tokens

### Token Management
- UpdateTokens method refreshes OAuth tokens
- Stores access_token, refresh_token, expires_at, scopes
- Automatic token updates when expired

### Account Linking
- Multiple OAuth accounts per user
- Provider-specific user IDs
- Account unlinking support

### Security Features
- OAuth state token for CSRF protection
- State TTL: 10 minutes (default)
- Secure token storage
- Provider validation

## Potential Improvements
- Add more OAuth providers (Facebook, Twitter, LinkedIn)
- Implement OAuth token refresh automation
- Add OAuth account merging
- Support OAuth scope management
- Add OAuth consent screen customization
- Implement OAuth account verification
- Support OAuth account migration
- Add OAuth token encryption at rest
- Support OAuth token revocation
- Add OAuth account analytics
- Implement OAuth account backup
- Support OAuth account recovery

