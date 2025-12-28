# API Endpoint Validation Documentation

**Last Updated:** 2025-12-24  
**Status:** ✅ Implemented

---

## Overview

This document describes the validation rules applied to API endpoints across the Woragis backend. All user inputs are validated to ensure data quality, security, and prevent common attacks like SQL injection and XSS.

---

## Validation Utilities

The validation system is built on top of the `pkg/validation` package, which provides:

- **Email validation** - RFC-compliant email format checking
- **UUID validation** - Standard UUID format validation
- **URL validation** - HTTP/HTTPS URL format checking
- **String validation** - Length constraints and sanitization
- **Security validation** - SQL injection and XSS detection
- **File validation** - File extension and size checking

---

## Authentication Endpoints (`/auth`)

### POST `/auth/register`

**Validates:**
- `email` (required): Valid email format, max 254 characters
- `password` (required): 8-128 characters
- `display_name` (optional): 1-100 characters, no SQL injection/XSS
- `locale` (optional): 2-10 characters

**Example:**
```json
{
  "email": "user@example.com",
  "password": "securepass123",
  "display_name": "John Doe",
  "locale": "en-US"
}
```

**Validation Errors:**
- `email: email is required`
- `email: invalid email format`
- `password: password is too short (minimum 8 characters)`
- `display_name: display_name: potentially dangerous input detected`

---

### POST `/auth/login`

**Validates:**
- `email` (required): Valid email format
- `password` (required): Non-empty
- `device_fingerprint` (optional): 1-255 characters
- `device_name` (optional): 1-100 characters, no SQL injection/XSS
- `mfa_code` (optional): Exactly 6 characters

**Example:**
```json
{
  "email": "user@example.com",
  "password": "securepass123",
  "device_name": "Chrome on Windows",
  "mfa_code": "123456"
}
```

---

### POST `/auth/password/reset/request`

**Validates:**
- `email` (required): Valid email format

**Example:**
```json
{
  "email": "user@example.com"
}
```

---

### POST `/auth/password/reset/confirm`

**Validates:**
- `token` (required): Non-empty
- `password` (required): 8-128 characters

**Example:**
```json
{
  "token": "reset-token-here",
  "password": "newsecurepass123"
}
```

---

### POST `/auth/session/refresh`

**Validates:**
- `refresh_token` (required): Non-empty

**Example:**
```json
{
  "refresh_token": "refresh-token-here"
}
```

---

### POST `/auth/mfa/enable`

**Validates:**
- `issuer` (optional): 1-100 characters, no SQL injection/XSS
- `label` (optional): 1-100 characters, no SQL injection/XSS
- `code` (optional): Exactly 6 characters

**Example:**
```json
{
  "issuer": "Woragis",
  "label": "user@example.com",
  "code": "123456"
}
```

---

### POST `/auth/mfa/verify`

**Validates:**
- `code` (required): Exactly 6 characters

**Example:**
```json
{
  "code": "123456"
}
```

---

### POST `/auth/oauth/start`

**Validates:**
- `provider` (required): 1-50 characters
- `mode` (optional): 1-50 characters
- `redirect_origin` (optional): Valid HTTP/HTTPS URL
- `device_fingerprint` (optional): 1-255 characters
- `device_name` (optional): 1-100 characters, no SQL injection/XSS

**Example:**
```json
{
  "provider": "google",
  "mode": "login",
  "redirect_origin": "https://app.woragis.com",
  "device_name": "Chrome on Windows"
}
```

---

## Job Applications Endpoints (`/job-applications`)

### POST `/job-applications`

**Validates:**
- `companyName` (required): 1-200 characters, no SQL injection/XSS
- `location` (optional): 1-200 characters, no SQL injection/XSS
- `jobTitle` (required): 1-200 characters, no SQL injection/XSS
- `jobUrl` (optional): Valid HTTP/HTTPS URL
- `website` (optional): Valid domain format, max 255 characters
- `interestLevel` (optional): Must be one of: `low`, `medium`, `high`, `very_high`
- `tags` (optional): Array of strings, max 20 tags, each 1-50 characters, no SQL injection/XSS
- `notes` (optional): 1-5000 characters, no SQL injection/XSS

**Example:**
```json
{
  "companyName": "Tech Corp",
  "location": "San Francisco, CA",
  "jobTitle": "Senior Software Engineer",
  "jobUrl": "https://techcorp.com/jobs/123",
  "website": "techcorp.com",
  "interestLevel": "high",
  "tags": ["remote", "python", "go"],
  "notes": "Great opportunity with interesting tech stack"
}
```

**Validation Errors:**
- `companyName: companyName is required`
- `jobTitle: jobTitle is too long (maximum 200 characters)`
- `jobUrl: jobUrl: invalid URL format`
- `interestLevel: interestLevel: must be one of: low, medium, high, very_high`
- `tags: tags: too many tags (maximum 20)`
- `notes: notes: notes is too long (maximum 5000 characters)`

---

### PATCH `/job-applications/:id/status`

**Validates:**
- `status` (required): Valid application status

**Example:**
```json
{
  "status": "applied"
}
```

---

## Path Parameter Validation

All endpoints that accept UUIDs in path parameters validate the format:

- `/auth/users/:id` - Validates UUID format
- `/job-applications/:id` - Validates UUID format
- `/resumes/:id` - Validates UUID format

**Validation Error:**
```json
{
  "error": "Invalid path parameter",
  "field": "id",
  "message": "invalid UUID format"
}
```

---

## Security Validations

### SQL Injection Detection

All string inputs are checked for common SQL injection patterns:
- `' OR '1'='1`
- `'; DROP TABLE`
- `'; DELETE FROM`
- `UNION SELECT`
- And other common patterns

**Error Response:**
```json
{
  "error": "Invalid query parameter",
  "field": "companyName",
  "message": "potentially dangerous input detected"
}
```

### XSS Detection

All string inputs are checked for XSS patterns:
- `<script>` tags
- `javascript:` protocol
- Event handlers (`onerror=`, `onload=`, `onclick=`)
- `<iframe>`, `<img>`, `<svg>` tags

**Error Response:**
```json
{
  "error": "Invalid query parameter",
  "field": "display_name",
  "message": "potentially dangerous input detected"
}
```

---

## Validation Error Response Format

All validation errors return a consistent format:

```json
{
  "error": "Invalid payload",
  "message": "field_name: validation error message"
}
```

**HTTP Status Code:** `400 Bad Request`

---

## Implementation Details

### Validation Functions Location

- **Auth domain:** `server/app/internal/domains/auth/validation.go`
- **Job Applications domain:** `server/app/internal/domains/jobapplications/validation.go`
- **Core utilities:** `server/app/pkg/validation/`

### Adding Validation to New Endpoints

1. Create validation function in domain's `validation.go` file:
```go
func ValidateCreateXPayload(payload *createXPayload) error {
    if err := validation.ValidateEmail(payload.Email); err != nil {
        return fmt.Errorf("email: %w", err)
    }
    // ... more validations
    return nil
}
```

2. Call validation in handler:
```go
func (h *Handler) CreateX(c *fiber.Ctx) error {
    var payload createXPayload
    if err := c.BodyParser(&payload); err != nil {
        return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
    }
    
    // Validate payload
    if err := ValidateCreateXPayload(&payload); err != nil {
        return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
            "message": err.Error(),
        })
    }
    
    // ... rest of handler logic
}
```

---

## Testing

Validation functions have comprehensive unit tests in:
- `server/app/pkg/validation/validator_test.go`
- Domain-specific validation tests (as needed)

**Test Coverage:** 33.1%+ for validation utilities

---

## Best Practices

1. **Always validate user input** - Never trust client input
2. **Validate early** - Validate before processing
3. **Return clear error messages** - Help users fix their input
4. **Sanitize before storing** - Use validation utilities' sanitization
5. **Log validation failures** - For security monitoring

---

## Future Enhancements

- [ ] Add validation to Resume endpoints
- [ ] Add validation to Social Media Posts endpoints
- [ ] Add validation to Creative Assets endpoints
- [ ] Add rate limiting per endpoint
- [ ] Add request size validation per endpoint
- [ ] Add custom validation rules for business logic

---

## Related Documentation

- [Input Validation Guide](../../deployment/input-validation.md)
- [Security Middleware](../../deployment/security-middleware.md)
- [API Documentation](../api/)

