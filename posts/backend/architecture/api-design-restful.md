# API Design: RESTful Services

## Overview
RESTful API design principles and patterns used in the Woragis backend server.

## Key Points

### API Principles
- RESTful conventions
- Resource-based URLs
- HTTP methods (GET, POST, PUT, PATCH, DELETE)
- JSON request/response
- Status codes

### Endpoint Structure
- `/api/{domain}/{resource}`
- Nested resources: `/api/{domain}/{resource}/{id}/{subresource}`
- Query parameters for filtering
- Pagination support

## Implementation Details

### Endpoint Examples
```
GET    /api/projects              # List projects
POST   /api/projects              # Create project
GET    /api/projects/{id}         # Get project
PUT    /api/projects/{id}         # Update project
PATCH  /api/projects/{id}         # Partial update
DELETE /api/projects/{id}         # Delete project
```

### Authentication
- JWT tokens
- Bearer token in Authorization header
- Refresh token mechanism
- API key support (for public endpoints)

### Error Handling
- Standard error format
- HTTP status codes
- Error codes for client handling
- Error messages (user-friendly)

## Response Format

### Success Response
```json
{
  "data": { ... },
  "meta": {
    "pagination": { ... }
  }
}
```

### Error Response
```json
{
  "error": {
    "code": "ERR_CODE",
    "message": "User-friendly message",
    "details": { ... }
  }
}
```

## Best Practices
- Consistent naming
- Versioning (future)
- Rate limiting
- Request validation
- Response caching

## Lessons Learned
- RESTful conventions help
- Consistent error format important
- Authentication crucial
- Documentation helps

## Future Improvements
- API versioning
- GraphQL (if needed)
- OpenAPI/Swagger docs
- Rate limiting per user
- Response compression
