# Auth Domain - Audit Logging

## Overview
Security audit logging for compliance and security monitoring.

## Key Points

### Audit Actions
- User registration, email confirmation
- Login success/failure
- Password reset requested/completed
- Session revocation
- MFA enabled/disabled
- OAuth account linked/unlinked
- Profile updates
- Bulk user actions

### AuditLog Structure
- UserID: Optional (nil for anonymous actions)
- Action: Audit action type
- Metadata: JSON map of additional data
- IP: Client IP address
- UserAgent: Client user agent
- CreatedAt: Timestamp of action

### Audit Events
- Security-sensitive operations logged
- User actions tracked
- System actions tracked
- Anonymous actions supported

### Metadata Storage
- JSONB storage for flexible metadata
- Action-specific data stored
- Structured for easy querying

### Use Cases
- Compliance requirements (GDPR, SOC2)
- Security incident investigation
- User activity tracking
- Fraud detection
- Access control auditing

## Potential Improvements
- Add audit log retention policies
- Implement audit log archival
- Add audit log search and filtering
- Support audit log export
- Add audit log alerting (suspicious activities)
- Implement audit log encryption
- Add audit log integrity verification
- Support audit log streaming
- Add audit log analytics
- Implement audit log compression
- Add audit log deduplication
- Support audit log correlation

