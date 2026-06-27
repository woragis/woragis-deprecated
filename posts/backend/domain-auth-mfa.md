# Auth Domain - MFA Implementation

## Overview
Multi-factor authentication implementation using TOTP and backup codes.

## Key Points

### MFA Types
- TOTP (Time-based One-Time Password): Authenticator apps
- Backup Codes: Static recovery codes

### MFAToken Structure
- User-scoped MFA enrollment
- Secret: Base32 encoded secret for TOTP
- Issuer and Label: For QR code generation
- Activation tracking (ActivatedAt)
- Usage tracking (LastUsedAt)
- Revocation support (RevokedAt)
- Backup codes: JSON array of recovery codes

### MFA Enrollment Flow
1. User requests MFA enrollment
2. Generate secret (32 bytes, base32 encoded)
3. Generate backup codes (8 codes, 8 characters each)
4. Create MFAToken record (inactive)
5. Return secret for QR code generation
6. User verifies with TOTP code
7. Activate MFA token
8. Store backup codes securely

### MFA Verification
- Verify TOTP code during login
- Verify backup code as fallback
- Update LastUsedAt on successful verification
- Prevent reuse of backup codes

### Backup Codes
- 8 codes generated per enrollment
- 8 characters each (base32)
- Single-use codes
- Should be shown once (during enrollment)

### Security Features
- Cryptographically secure random secret generation
- Activation required before use
- Revocation support
- Last used timestamp tracking

## Potential Improvements
- Add SMS-based MFA
- Support hardware security keys (WebAuthn/FIDO2)
- Add biometric MFA
- Implement MFA policy enforcement
- Add MFA backup code regeneration
- Support multiple MFA methods per user
- Add MFA recovery options
- Implement MFA rate limiting
- Add MFA audit logging
- Support MFA device management
- Add MFA trust devices
- Implement MFA step-up authentication

