# Worker Validation Implementation

This document describes the validation implementation for all standalone workers in the backend.

## Overview

All standalone workers now include comprehensive input validation to ensure:
- Data integrity
- Security (SQL injection, XSS prevention)
- Format validation (UUIDs, URLs, emails, phone numbers)
- Length constraints
- Type validation

## Implemented Workers

### 1. Resume Worker (Python)

**Location**: `backend/resume-worker/src/validation.py`

**Validations**:
- UUID validation for user IDs and job IDs
- String length validation (job descriptions, titles, filenames)
- Language code validation (ISO 639-1)
- SQL injection prevention
- XSS prevention
- Path traversal prevention for filenames
- URL validation

**Integration Points**:
- `main.py`: Validates queue messages and CLI arguments
- `process_resume_job()`: Validates job messages before processing
- `run_cli_mode()`: Validates command-line arguments

**Example**:
```python
from validation import validate_resume_job_message, validate_cli_arguments

# In queue consumer
validate_resume_job_message(message)

# In CLI mode
validate_cli_arguments(sys.argv)
```

### 2. Translation Worker (Go)

**Location**: `backend/translation-worker/internal/queue/validation.go`

**Validations**:
- UUID validation for job IDs and entity IDs
- Entity type validation (whitelist of allowed types)
- Language code validation (ISO 639-1 or locale codes like "pt-BR")
- Field array validation (max 50 fields)
- Source text validation (max 50 entries, length constraints)
- SQL injection prevention
- XSS prevention

**Integration Points**:
- `cmd/translation-worker/main.go`: Validates translation jobs before processing
- `processTranslationJob()`: Job is validated before any processing

**Example**:
```go
// In queue consumer
if err := queue.ValidateTranslationJob(job); err != nil {
    logger.Error("Invalid translation job", slog.Any("error", err))
    return fmt.Errorf("validation failed: %w", err)
}
```

### 3. WhatsApp Worker (Go)

**Location**: `backend/whatsapp-worker/internal/queue/validation.go`

**Validations**:
- UUID validation for user IDs
- Phone number validation (international format with + prefix, 8-15 digits)
- Message text validation (1-4096 characters)
- SQL injection prevention
- XSS prevention

**Integration Points**:
- `cmd/whatsapp-worker/main.go`: Validates WhatsApp envelopes before sending

**Example**:
```go
// In queue consumer
if err := queue.ValidateWhatsAppEnvelope(envelope); err != nil {
    logger.Error("Invalid WhatsApp message", slog.Any("error", err))
    return fmt.Errorf("validation failed: %w", err)
}
```

### 4. Job Application Worker (Node.js)

**Location**: `backend/job-application-worker/src/validation.js`

**Validations**:
- UUID validation for job IDs and user IDs
- String length validation
- URL validation for job URLs
- Website format validation (no spaces)
- SQL injection prevention
- XSS prevention
- Scraped job info validation
- Cover letter content validation

**Integration Points**:
- `worker.js`: Validates job application jobs and cover letters
- `coverLetter.js`: Validates job info inputs before generating cover letters

**Example**:
```javascript
// In worker
import { validateJobApplicationJob, validateCoverLetter } from './validation.js';

// Validate job
validateJobApplicationJob(job);

// Validate cover letter
validateCoverLetter(coverLetter);
```

## Common Validation Patterns

### UUID Validation
All workers validate UUIDs using regex pattern:
```
^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$
```

### SQL Injection Prevention
All workers check for common SQL injection patterns:
- SQL keywords (SELECT, INSERT, UPDATE, DELETE, etc.)
- Comment patterns (--, #, /* */)
- Union-based attacks
- Special characters that could break SQL queries

### XSS Prevention
All workers check for common XSS patterns:
- `<script>` tags
- `javascript:` protocol
- Event handlers (`onclick`, `onerror`, etc.)
- `<iframe>`, `<object>`, `<embed>` tags

### String Length Validation
All workers enforce minimum and maximum length constraints:
- Job descriptions: 10-50000 characters
- Job titles: 1-200 characters
- Company names: 1-200 characters
- URLs: Valid HTTP/HTTPS format
- Phone numbers: 8-15 digits after country code

## Error Handling

All validation errors are:
1. Logged with appropriate context
2. Returned as structured errors
3. Tracked in metrics (validation_error)
4. Prevent processing of invalid data

## Security Considerations

1. **Input Sanitization**: All inputs are sanitized before processing
2. **Path Traversal Prevention**: Filenames are validated to prevent directory traversal
3. **Type Safety**: All types are validated before use
4. **Length Limits**: All inputs have maximum length constraints to prevent DoS attacks

## Testing

Each worker's validation can be tested by:
1. Sending invalid messages to the queue
2. Verifying that validation errors are logged
3. Confirming that invalid messages are rejected
4. Checking that metrics record validation errors

## Future Enhancements

- Add rate limiting validation
- Add content validation for generated content
- Add schema validation for complex nested structures
- Add validation for file uploads (if applicable)

