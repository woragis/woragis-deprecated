# Server API Documentation

## Overview

The Woragis Server API is a RESTful API built with Go and Fiber framework. It provides endpoints for managing projects, resumes, job applications, authentication, and more.

## Base URL

- **Development**: `http://localhost:8080`
- **Production**: `https://api.woragis.com`

## Authentication

Most endpoints require authentication via JWT tokens.

### Getting a Token

**POST** `/api/auth/login`

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "...",
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "name": "User Name"
  }
}
```

### Using the Token

Include the token in the `Authorization` header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Common Response Formats

### Success Response

```json
{
  "data": { ... },
  "message": "Success message"
}
```

### Error Response

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": { ... }
}
```

### HTTP Status Codes

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## API Endpoints

### Authentication

#### Register

**POST** `/api/auth/register`

Create a new user account.

**Request**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "User Name"
}
```

#### Login

**POST** `/api/auth/login`

Authenticate and get JWT token.

**Request**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Get Current User

**GET** `/api/auth/me`

Get current authenticated user information.

**Headers**: `Authorization: Bearer {token}`

#### Refresh Token

**POST** `/api/auth/refresh`

Refresh access token using refresh token.

**Request**:
```json
{
  "refresh_token": "..."
}
```

#### OAuth Providers

**GET** `/api/auth/oauth/providers`

List available OAuth providers (GitHub, Google, Microsoft).

**POST** `/api/auth/oauth/start`

Start OAuth flow.

**Request**:
```json
{
  "provider": "github"
}
```

**GET** `/api/auth/oauth/callback/:provider`

OAuth callback endpoint.

### Projects

#### List Projects

**GET** `/api/projects`

List all projects for the authenticated user.

**Query Parameters**:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20)
- `status` - Filter by status
- `search` - Search query

#### Get Project

**GET** `/api/projects/:id`

Get project by ID.

#### Create Project

**POST** `/api/projects`

Create a new project.

**Request**:
```json
{
  "name": "Project Name",
  "description": "Project description",
  "status": "active",
  "technologies": ["go", "postgresql", "docker"]
}
```

#### Update Project

**PATCH** `/api/projects/:id`

Update project.

**Request**:
```json
{
  "name": "Updated Name",
  "status": "completed"
}
```

#### Delete Project

**DELETE** `/api/projects/:id`

Delete project.

#### Get Project by Slug

**GET** `/api/projects/slug/:slug`

Get project by slug (public access).

### Resumes

#### List Resumes

**GET** `/api/resumes`

List all resumes for the authenticated user.

#### Get Resume

**GET** `/api/resumes/:id`

Get resume by ID.

#### Create Resume

**POST** `/api/resumes`

Create a new resume.

**Request**:
```json
{
  "name": "Resume Name",
  "job_title": "Software Engineer",
  "job_description": "Job description..."
}
```

#### Generate Resume

**POST** `/api/resumes/:id/generate`

Generate resume PDF (publishes job to queue).

**Request**:
```json
{
  "job_description": "Looking for a backend engineer...",
  "job_title": "Backend Engineer"
}
```

### Job Applications

#### List Job Applications

**GET** `/api/job-applications`

List all job applications for the authenticated user.

#### Get Job Application

**GET** `/api/job-applications/:id`

Get job application by ID.

#### Create Job Application

**POST** `/api/job-applications`

Create a new job application.

**Request**:
```json
{
  "job_url": "https://linkedin.com/jobs/view/123456",
  "website": "linkedin",
  "resume_id": "resume-uuid",
  "cover_letter_template": "template-uuid"
}
```

#### Process Job Application

**POST** `/api/job-applications/:id/process`

Process job application (publishes job to queue).

### Translations

#### Request Translation

**POST** `/api/translations`

Request translation for content.

**Request**:
```json
{
  "entity_type": "project",
  "entity_id": "entity-uuid",
  "language": "pt-BR",
  "fields": ["name", "description"]
}
```

### Skills

#### List Skills

**GET** `/api/skills`

List all skills.

#### Create Skill

**POST** `/api/skills`

Create a new skill.

**Request**:
```json
{
  "name": "Go",
  "category": "programming",
  "level": "expert"
}
```

### Certifications

#### List Certifications

**GET** `/api/certifications`

List all certifications for the authenticated user.

#### Create Certification

**POST** `/api/certifications`

Create a new certification.

**Request**:
```json
{
  "name": "AWS Certified Solutions Architect",
  "issuer": "Amazon Web Services",
  "issue_date": "2024-01-15",
  "expiry_date": "2027-01-15",
  "certificate_url": "https://..."
}
```

### Chats

#### List Chats

**GET** `/api/chats`

List all chats for the authenticated user.

#### Get Chat

**GET** `/api/chats/:id`

Get chat by ID.

#### Create Chat

**POST** `/api/chats`

Create a new chat.

**Request**:
```json
{
  "title": "Chat Title",
  "agent": "economist"
}
```

#### Send Message

**POST** `/api/chats/:id/messages`

Send a message in a chat.

**Request**:
```json
{
  "content": "Message content"
}
```

### WebSocket

**WS** `/ws/chat/:chatId`

WebSocket endpoint for real-time chat.

**Connection**: Include JWT token in query parameter or header.

## Rate Limiting

Rate limiting may be implemented on certain endpoints. Check response headers:

- `X-RateLimit-Limit` - Request limit per window
- `X-RateLimit-Remaining` - Remaining requests
- `X-RateLimit-Reset` - Reset time (Unix timestamp)

## Pagination

List endpoints support pagination:

**Query Parameters**:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20, max: 100)

**Response**:
```json
{
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

## Filtering and Sorting

Many list endpoints support filtering and sorting:

**Query Parameters**:
- `filter[field]` - Filter by field value
- `sort` - Sort field (e.g., `created_at`)
- `order` - Sort order (`asc` or `desc`)

**Example**:
```
GET /api/projects?filter[status]=active&sort=created_at&order=desc
```

## Health Check

**GET** `/healthz`

Check server health and dependencies.

**Response**:
```json
{
  "status": "healthy|degraded|unhealthy",
  "checks": [
    {"name": "database", "status": "ok"},
    {"name": "redis", "status": "ok"},
    {"name": "rabbitmq", "status": "ok"}
  ]
}
```

## Metrics

**GET** `/metrics`

Prometheus metrics endpoint (internal use only).

## Error Codes

Common error codes:

- `AUTH_REQUIRED` - Authentication required
- `AUTH_INVALID` - Invalid credentials
- `AUTH_EXPIRED` - Token expired
- `FORBIDDEN` - Insufficient permissions
- `NOT_FOUND` - Resource not found
- `VALIDATION_ERROR` - Request validation failed
- `RATE_LIMIT_EXCEEDED` - Rate limit exceeded
- `INTERNAL_ERROR` - Internal server error

## Examples

### Complete Flow: Create Project and Request Translation

```bash
# 1. Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Response: {"token":"...","user":{...}}

# 2. Create Project
curl -X POST http://localhost:8080/api/projects \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Project",
    "description": "Project description",
    "status": "active"
  }'

# Response: {"data":{"id":"project-uuid",...}}

# 3. Request Translation
curl -X POST http://localhost:8080/api/translations \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "entity_type": "project",
    "entity_id": "project-uuid",
    "language": "pt-BR",
    "fields": ["name", "description"]
  }'
```

## Related Documentation

- [Component Documentation](../components/server.md) - Server component details
- [Development Guides](../development/) - How to extend the API
- [Architecture Decision Records](../adr/) - Architectural decisions
