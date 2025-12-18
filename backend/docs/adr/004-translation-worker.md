# ADR-004: Translation Worker Architecture

## Status
Accepted

## Context
The system needs to translate content (projects, testimonials, certifications, etc.) into multiple languages. We need to decide:
1. Where the translation logic should live (server vs. standalone worker)
2. What language to use (Go vs. Python)
3. How translations are stored (server saves vs. worker saves directly)

**Requirements:**
- Support multiple translation providers (Google Translate, DeepL, LibreTranslate)
- High performance (many translations needed)
- Direct database writes (no round-trip through server)
- Consistency with other workers
- Retry logic for transient failures

**Constraints:**
- Already have Go workers (email, WhatsApp)
- Python is used for AI/ML services
- Translation APIs are HTTP-based (language-agnostic)
- Need to minimize latency and code duplication

## Decision
We will implement a **standalone Go-based translation worker** that:
- Consumes translation jobs from RabbitMQ
- Calls translation APIs (Google/DeepL/LibreTranslate) via HTTP
- Writes translations directly to the database
- Implements retry logic with exponential backoff

**Architecture:**
```
Server → RabbitMQ → Translation Worker → Translation API → Database
```

**Rationale:**
- **Go for Performance**: Consistent with other workers, high performance
- **Standalone Worker**: Independent scaling, fault isolation
- **Direct DB Writes**: Eliminates round-trip, reduces latency
- **HTTP APIs**: Translation APIs are HTTP-based, language doesn't matter

## Consequences

### Positive
- ✅ **Performance**: Go provides high performance for concurrent translations
- ✅ **Consistency**: Matches architecture of other Go workers
- ✅ **Direct Writes**: No round-trip through server, lower latency
- ✅ **Independent Scaling**: Scale translation worker independently
- ✅ **Fault Isolation**: Translation failures don't affect server
- ✅ **Language Flexibility**: Can call any HTTP-based translation API

### Negative
- ⚠️ **Code Duplication**: Some worker patterns duplicated (but acceptable)
- ⚠️ **Database Access**: Worker needs direct database access
- ⚠️ **No Server Validation**: Server doesn't validate translations before saving

### Neutral
- Translation APIs are HTTP-based, so language choice doesn't matter for API calls
- Worker follows same patterns as other workers (health checks, logging, metrics)

## Implementation Details

### Worker Structure
```
translation-worker/
├── cmd/
│   └── translation-worker/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── translator/
│   └── queue/
├── pkg/
│   ├── health/
│   ├── logger/
│   └── metrics/
├── go.mod
├── Dockerfile
└── README.md
```

### Translation Providers

#### Google Translate
- **API**: Google Cloud Translation API
- **Authentication**: API key or service account
- **Rate Limits**: Per project quotas
- **Cost**: Pay-per-character

#### DeepL
- **API**: DeepL API v2
- **Authentication**: API key
- **Rate Limits**: Per subscription tier
- **Cost**: Subscription-based

#### LibreTranslate
- **API**: LibreTranslate API (self-hosted or public)
- **Authentication**: Optional API key
- **Rate Limits**: Self-hosted (unlimited) or public (limited)
- **Cost**: Free (self-hosted) or usage-based (public)

### Message Format
```json
{
  "id": "job-uuid",
  "entityType": "project",
  "entityId": "entity-uuid",
  "language": "pt-BR",
  "fields": ["name", "description"],
  "sourceText": {
    "name": "My Project",
    "description": "Project description"
  }
}
```

### Retry Logic
- **Max Retries**: 3 attempts (configurable)
- **Retry Delay**: Exponential backoff (1s, 2s, 4s)
- **Retry Conditions**: Transient errors (network, timeouts)
- **No Retry**: Permanent errors (invalid input, auth failures)

### Database Writes
- Worker writes directly to database after successful translation
- Updates translation records in database
- No round-trip through server
- Database connection pool per worker instance

## Alternatives Considered

### 1. Translation Logic in Server
- **Pros**: Simpler architecture, no separate service
- **Cons**: Blocks server requests, can't scale independently, language restrictions
- **Rejected**: Doesn't meet performance and scaling requirements

### 2. Python Worker
- **Pros**: Easier AI/ML integration (if needed), consistent with AI services
- **Cons**: Slower than Go, inconsistent with other workers
- **Rejected**: Performance is important, Go is faster for concurrent operations

### 3. Server Receives Translation Back
- **Pros**: Server validates before saving, centralized logic
- **Cons**: Extra round-trip, higher latency, more complex
- **Rejected**: Performance is more important than validation (translations are trusted)

### 4. Hybrid Approach (Go Worker + External APIs)
- **Pros**: Best of both worlds (Go performance + API flexibility)
- **Cons**: None significant
- **Accepted**: This is the chosen approach

## Implementation Example

### Worker Consumption
```go
func (w *Worker) Consume(ctx context.Context) error {
    return w.queue.Consume(ctx, func(job *TranslationJob) error {
        // Translate
        translated, err := w.translator.Translate(ctx, job.SourceText, job.Language)
        if err != nil {
            return err // Will retry if transient
        }
        
        // Write to database
        err = w.database.SaveTranslation(ctx, job.EntityType, job.EntityID, job.Language, translated)
        if err != nil {
            return err // Will retry
        }
        
        return nil // Success, acknowledge message
    })
}
```

### Retry Logic
```go
maxRetries := 3
retryDelay := 1000 * time.Millisecond

for attempt := 0; attempt < maxRetries; attempt++ {
    result, err := translator.Translate(ctx, text, targetLang)
    if err == nil {
        return result, nil
    }
    
    if !isTransientError(err) {
        return "", err // Don't retry permanent errors
    }
    
    if attempt < maxRetries-1 {
        time.Sleep(retryDelay * time.Duration(1<<attempt)) // Exponential backoff
    }
}
```

## Supported Entity Types
- `testimonial`
- `project`
- `certification`
- `skill`
- (More can be added)

## Supported Languages
- `en` - English
- `pt-BR` - Portuguese (Brazil)
- `pt` - Portuguese
- `es` - Spanish
- `fr` - French
- `de` - German
- `ru` - Russian
- `ja` - Japanese
- `ko` - Korean
- `zh-CN` - Chinese (Simplified)
- `el` - Greek
- `la` - Latin

## Future Enhancements

### Planned
- Support for more entity types
- Support for more languages
- Translation caching (avoid re-translating same content)
- Batch translation (translate multiple fields in one API call)
- Translation quality metrics

### Under Consideration
- Multiple translation providers (fallback if one fails)
- Translation memory (reuse previous translations)
- Human review workflow (for critical content)

## Notes
- Worker automatically fetches source text from database if not provided in message
- Worker handles multiple translation providers (configurable via env var)
- Worker implements Dead Letter Queue (DLQ) for failed messages
- Worker exposes health check and metrics endpoints

## Related ADRs
- [ADR-002: Standalone Workers Architecture](./002-standalone-workers.md) - Worker architecture pattern
- [ADR-001: RabbitMQ with Redis Fallback](./001-rabbitmq-redis-fallback.md) - Message queue pattern
- [ADR-003: Structured Logging Implementation](./003-structured-logging.md) - Logging in worker
