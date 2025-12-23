# Input Validation & Security Hardening Guide

**Last Updated:** 2025-12-22  
**Purpose:** Guide for input validation and security hardening

---

## Overview

This guide covers input validation, sanitization, and security measures to protect against common vulnerabilities.

---

## Common Vulnerabilities

### 1. SQL Injection

**Prevention:** Always use parameterized queries

**Go Example:**
```go
// ❌ BAD - Vulnerable to SQL injection
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
db.Query(query)

// ✅ GOOD - Parameterized query
query := "SELECT * FROM users WHERE email = $1"
db.Query(query, email)
```

**Python Example:**
```python
# ❌ BAD
cursor.execute(f"SELECT * FROM users WHERE email = '{email}'")

# ✅ GOOD
cursor.execute("SELECT * FROM users WHERE email = %s", (email,))
```

---

### 2. XSS (Cross-Site Scripting)

**Prevention:** Escape user input, use Content Security Policy

**Go Example:**
```go
import "html"

// Escape HTML
safeOutput := html.EscapeString(userInput)

// Or use template engine with auto-escaping
tmpl.Execute(w, data)  // Auto-escapes by default
```

**Python Example:**
```python
from markupsafe import escape

# Escape HTML
safe_output = escape(user_input)

# In templates (Jinja2)
{{ user_input | escape }}
```

---

### 3. CSRF (Cross-Site Request Forgery)

**Prevention:** Use CSRF tokens

**Go Example:**
```go
import "github.com/gofiber/fiber/v2/middleware/csrf"

app.Use(csrf.New(csrf.Config{
    KeyLookup: "header:X-Csrf-Token",
    CookieName: "csrf_",
    CookieSameSite: "Strict",
}))
```

---

### 4. Path Traversal

**Prevention:** Validate and sanitize file paths

```go
import "path/filepath"

func SafePath(userPath string, baseDir string) (string, error) {
    // Resolve to absolute path
    absPath, err := filepath.Abs(filepath.Join(baseDir, userPath))
    if err != nil {
        return "", err
    }
    
    // Ensure path is within base directory
    baseAbs, _ := filepath.Abs(baseDir)
    if !strings.HasPrefix(absPath, baseAbs) {
        return "", errors.New("path traversal detected")
    }
    
    return absPath, nil
}
```

---

## Input Validation

### 1. Email Validation

```go
import "regexp"

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) bool {
    if len(email) > 254 {
        return false
    }
    return emailRegex.MatchString(email)
}
```

**Python:**
```python
import re
from email.utils import parseaddr

def validate_email(email: str) -> bool:
    if len(email) > 254:
        return False
    pattern = r'^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$'
    return bool(re.match(pattern, email))
```

---

### 2. URL Validation

```go
import "net/url"

func ValidateURL(urlString string) bool {
    u, err := url.Parse(urlString)
    if err != nil {
        return false
    }
    return u.Scheme == "http" || u.Scheme == "https"
}
```

---

### 3. String Length Validation

```go
func ValidateString(input string, min, max int) error {
    length := len([]rune(input))  // Count runes, not bytes
    if length < min {
        return fmt.Errorf("input too short (minimum %d characters)", min)
    }
    if length > max {
        return fmt.Errorf("input too long (maximum %d characters)", max)
    }
    return nil
}
```

---

### 4. Numeric Validation

```go
func ValidateInt(value int, min, max int) error {
    if value < min {
        return fmt.Errorf("value too small (minimum %d)", min)
    }
    if value > max {
        return fmt.Errorf("value too large (maximum %d)", max)
    }
    return nil
}
```

---

## Request Size Limits

### Go (Fiber)

```go
app.Use(limiter.New(limiter.Config{
    MaxRequestBodySize: 10 * 1024 * 1024,  // 10MB
}))

// Or per route
app.Post("/upload", func(c *fiber.Ctx) error {
    // Max 5MB for this route
    c.Request().Header.Set("Content-Length", "5242880")
    // ...
})
```

### Python (FastAPI)

```python
from fastapi import Request
from fastapi.exceptions import RequestValidationError

@app.middleware("http")
async def limit_upload_size(request: Request, call_next):
    if request.method == "POST":
        content_length = request.headers.get("content-length")
        if content_length and int(content_length) > 10 * 1024 * 1024:  # 10MB
            raise RequestValidationError("Request too large")
    return await call_next(request)
```

---

## File Upload Validation

```go
func ValidateFileUpload(file *multipart.FileHeader) error {
    // Check file size (5MB max)
    if file.Size > 5*1024*1024 {
        return errors.New("file too large (maximum 5MB)")
    }
    
    // Check file extension
    ext := filepath.Ext(file.Filename)
    allowedExts := []string{".jpg", ".jpeg", ".png", ".pdf"}
    valid := false
    for _, allowedExt := range allowedExts {
        if ext == allowedExt {
            valid = true
            break
        }
    }
    if !valid {
        return errors.New("invalid file type")
    }
    
    // Check MIME type
    fileType := file.Header.Get("Content-Type")
    allowedTypes := []string{"image/jpeg", "image/png", "application/pdf"}
    valid = false
    for _, allowedType := range allowedTypes {
        if fileType == allowedType {
            valid = true
            break
        }
    }
    if !valid {
        return errors.New("invalid file type")
    }
    
    return nil
}
```

---

## Input Sanitization

### HTML Sanitization

```go
import "github.com/microcosm-cc/bluemonday"

p := bluemonday.UGCPolicy()  // User Generated Content policy
sanitized := p.Sanitize(userInput)
```

### SQL Sanitization

**Don't sanitize SQL - use parameterized queries instead!**

```go
// ✅ Always use parameterized queries
db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
```

---

## Validation Middleware

### Go Example

```go
type ValidationMiddleware struct {
    validator *validator.Validate
}

func NewValidationMiddleware() *ValidationMiddleware {
    return &ValidationMiddleware{
        validator: validator.New(),
    }
}

func (v *ValidationMiddleware) Validate(schema interface{}) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var body interface{}
        if err := c.BodyParser(&body); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
        }
        
        if err := v.validator.Struct(body); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": err.Error()})
        }
        
        c.Locals("validated", body)
        return c.Next()
    }
}
```

### Python Example (Pydantic)

```python
from pydantic import BaseModel, EmailStr, validator
from typing import Optional

class UserCreate(BaseModel):
    email: EmailStr
    name: str
    age: Optional[int] = None
    
    @validator('name')
    def validate_name(cls, v):
        if len(v) < 2 or len(v) > 100:
            raise ValueError('Name must be between 2 and 100 characters')
        return v.strip()
    
    @validator('age')
    def validate_age(cls, v):
        if v is not None and (v < 0 or v > 150):
            raise ValueError('Age must be between 0 and 150')
        return v

@app.post("/users")
async def create_user(user: UserCreate):
    # user is already validated
    return {"message": "User created", "user": user}
```

---

## Security Headers

```go
import "github.com/gofiber/fiber/v2/middleware/helmet"

app.Use(helmet.New(helmet.Config{
    XSSProtection:             "1; mode=block",
    ContentTypeNosniff:        "nosniff",
    XFrameOptions:             "DENY",
    ReferrerPolicy:            "no-referrer",
    ContentSecurityPolicy:     "default-src 'self'",
    PermissionsPolicy:         "geolocation=(), microphone=(), camera=()",
}))
```

---

## Implementation Checklist

### Phase 1: Basic Validation (Week 1)
- [ ] Add input validation to all API endpoints
- [ ] Implement request size limits
- [ ] Add file upload validation
- [ ] Set security headers
- [ ] Validate email, URL, and other common inputs

### Phase 2: Security Hardening (Week 2)
- [ ] Implement CSRF protection
- [ ] Add XSS prevention
- [ ] Sanitize HTML input
- [ ] Validate file uploads
- [ ] Add rate limiting

### Phase 3: Advanced Protection (Week 3)
- [ ] Implement input sanitization
- [ ] Add path traversal protection
- [ ] Set up security monitoring
- [ ] Add security testing
- [ ] Document security measures

---

## Testing

### Security Testing Tools

- **OWASP ZAP** - Security scanner
- **Burp Suite** - Web vulnerability scanner
- **SQLMap** - SQL injection testing
- **XSStrike** - XSS vulnerability scanner

### Manual Testing

```bash
# SQL Injection test
curl -X POST http://api/endpoint -d "email=' OR '1'='1"

# XSS test
curl -X POST http://api/endpoint -d "name=<script>alert('XSS')</script>"

# Path traversal test
curl http://api/files/../../../etc/passwd
```

---

## Related Documentation

- [Authentication & Authorization](./authentication-authorization.md)
- [Secrets Management](./secrets-management.md)
- [Security Guide](../operations/security.md) (when created)

---

**Next Steps:**
1. Audit all API endpoints for input validation
2. Implement validation middleware
3. Add security headers
4. Set up security testing
5. Document security measures

---

**Last Updated:** 2025-12-22
