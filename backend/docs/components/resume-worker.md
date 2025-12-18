# Resume Worker Component

## Overview

A Python-based worker service that generates customized resumes (CVs) based on job requirements. It uses AI to generate resume sections, searches through projects by technology tags, and includes relevant certifications.

## Architecture

- **Language**: Python 3.11+
- **Port**: 8080 (health check only)
- **Message Queue**: RabbitMQ
- **Database**: PostgreSQL
- **AI Service**: Calls AI service for content generation
- **PDF Generation**: WeasyPrint for PDF creation

## Responsibilities

1. **Consume Resume Jobs**: Listen to `resumes.queue` for resume generation jobs
2. **Generate Resume Content**: Use AI service to generate resume sections
3. **Match Projects**: Search projects by technology categories matching job requirements
4. **Match Certifications**: Include relevant certifications based on job focus
5. **Generate PDF**: Create professional PDF resumes
6. **Update Database**: Save generation metadata

## Health Check

**Endpoint**: `GET /healthz`

**Checks**:
- RabbitMQ connection - CRITICAL

**Response**:
```json
{
  "status": "healthy|unhealthy",
  "checks": [
    {"name": "rabbitmq", "status": "ok"}
  ]
}
```

## Metrics

**Endpoint**: `GET /metrics`

Exposes Prometheus metrics:
- Job processing rate (success/failed)
- Job processing duration
- Queue depth
- DLQ size
- AI service call metrics

## Configuration

### Environment Variables

#### Required
- `DATABASE_URL` - PostgreSQL connection string
- `RABBITMQ_URL` - RabbitMQ connection URL
- `AI_SERVICE_URL` - AI service URL (default: `http://ai-service:8000`)

#### Optional
- `RESUME_QUEUE_NAME` - Queue name (default: `resumes.queue`)
- `RESUME_EXCHANGE` - Exchange name (default: `woragis.tasks`)
- `RESUME_ROUTING_KEY` - Routing key (default: `resumes.generate`)
- `RESUME_OUTPUT_DIR` - Directory for generated PDFs (default: `/app/output`)
- `RESULTS_LOG_DIR` - Directory for result metadata (default: `/app/results`)
- `ENV` - Environment (development/production)

## Message Format

**Queue**: `resumes.queue`

**Message**:
```json
{
  "id": "job-uuid",
  "user_id": "user-uuid",
  "job_description": "Looking for a backend engineer with Golang, Kubernetes, and microservices experience...",
  "job_title": "Backend Engineer",
  "output_filename": "resume_backend.pdf"
}
```

## Processing Flow

1. **Consume Message**: Worker consumes message from RabbitMQ
2. **Extract Keywords**: Extract technology categories from job description
3. **Match Projects**: Search projects by technology categories
4. **Match Certifications**: Find relevant certifications
5. **Generate Resume Sections**: Use AI service to generate:
   - Professional Summary
   - About Me
   - Technical Skills
   - Projects & Experience
   - Certifications
6. **Generate PDF**: Create PDF using HTML template
7. **Save Results**: Save PDF and metadata
8. **Update Database**: Record generation metadata
9. **Acknowledge**: Acknowledge message on success

## Resume Sections

1. **Professional Summary**: AI-generated 4-5 line summary optimized for the job
2. **About Me**: Brief personal introduction (3-4 sentences)
3. **Technical Skills**: Categorized list of relevant technical skills
4. **Projects & Experience**: Relevant projects with technologies and descriptions
5. **Certifications**: Relevant certifications with issue dates

## Technology Matching

Projects are matched by technology categories:

- **backend**: backend, back-end, server, api, rest, graphql, microservices
- **devops**: devops, ci/cd, docker, kubernetes, k8s, terraform, ansible, jenkins
- **infrastructure**: infrastructure, cloud, aws, azure, gcp
- **database**: database, sql, nosql, postgresql, mysql, mongodb, redis
- **frontend**: frontend, react, vue, angular, javascript, typescript
- **monitoring**: monitoring, observability, prometheus, grafana
- **testing**: testing, qa, test, tdd, bdd

## Certification Matching

Certifications are matched by category:

- **cloud**: cloud, aws, azure, gcp
- **devops**: devops, kubernetes, docker, terraform
- **security**: security, pentesting, penetration testing, cybersecurity
- **programming**: programming, python, golang, java, spring boot
- **database**: database, sql
- **architecture**: architecture, solution architect

## Output

### PDF Files
- **Location**: `RESUME_OUTPUT_DIR` (default: `/app/output`)
- **Naming**: `resume_{job_title}_{timestamp}.pdf`
- **Format**: ATS-friendly PDF (no images, standard fonts)

### Result Metadata
- **Location**: `RESULTS_LOG_DIR` (default: `/app/results`)
- **Format**: JSON files
- **Naming**: `resume_result_{user_id}_{timestamp}.json`
- **Contents**: File path, size, projects/certifications count, keywords, etc.

## Template Customization

Resume templates are located in `templates/`:
- `resume.html` - HTML template
- `style.css` - CSS styling

**Customizable**:
- Fonts and typography
- Colors and styling
- Layout (two-column sections, spacing)
- Section ordering
- Additional sections

## Error Handling

### Transient Errors
- AI service temporarily unavailable
- Database connection errors
- PDF generation errors

**Action**: Retry with exponential backoff (up to 3 attempts)

### Permanent Errors
- Invalid user ID
- Invalid job description
- Template errors

**Action**: Route to Dead Letter Queue (DLQ)

## Dead Letter Queue

**DLQ Exchange**: `woragis.dlx`
**DLQ Routing Key**: `resumes.queue.failed`

Failed messages are automatically routed to DLQ after max retries.

## Logging

**Format**: Structured JSON (production), Text (development)

**Service Name**: `resume-worker`

**Key Log Fields**:
- `job_id` - Job identifier
- `user_id` - User ID
- `job_title` - Job title
- `status` - Success/failed
- `error` - Error message (if failed)

## Deployment

### Local Development

```bash
cd backend/resume-worker
pip install -r requirements.txt
python src/main.py
```

### Docker

```bash
docker build -f Dockerfile.resume-worker -t woragis/resume-worker .
docker run --env-file .env -v ./output:/app/output -v ./results:/app/results woragis/resume-worker
```

### Kubernetes

Deploy as a Deployment:
- Health check probe on `/healthz`
- RabbitMQ connection required
- Database connection required
- AI service connection required
- Volume mounts for output/results directories

## Scaling

### Horizontal Scaling
- Multiple worker replicas consume from same queue
- RabbitMQ distributes messages (round-robin or fair dispatch)
- Each replica has its own database connection pool
- **Note**: Memory-intensive (PDF generation), may need more memory per replica

### Resource Requirements
- **CPU**: 200m-500m (0.2-0.5 core)
- **Memory**: 512Mi-1Gi (PDF generation is memory-intensive)
- **Database Connections**: 5-10 per replica

## AI Service Integration

The worker calls the AI service for:
- Professional Summary generation
- About Me generation
- Project descriptions enhancement

**AI Service Endpoints Used**:
- `POST /v1/chat` - Chat completion for content generation

## Monitoring

### Key Metrics
- Resume generation rate (resumes/second)
- Success rate
- Failure rate
- Queue depth
- DLQ size
- AI service call duration
- PDF generation duration

### Alerts
- DLQ size > 100 messages
- Failure rate > 5%
- Queue depth > 1000 messages
- AI service error rate > 10%

## Troubleshooting

### Common Issues

#### Resumes Not Generating
- Check AI service configuration
- Check database connection
- Check template files exist
- Verify message format

#### High Failure Rate
- Check AI service status
- Verify user ID exists in database
- Check template syntax
- Verify output directory permissions

#### Queue Backlog
- Check AI service availability (may be bottleneck)
- Scale up worker replicas
- Check worker health
- Verify workers are consuming messages

#### PDF Generation Errors
- Check memory limits (PDF generation is memory-intensive)
- Verify template files are valid
- Check WeasyPrint dependencies
- Verify output directory permissions

## Related Documentation

- [Architecture Decision Records](../adr/) - Worker architecture decisions
- [Monitoring DLQ](../runbooks/monitoring-dlq.md) - DLQ monitoring procedures
- [Deploying Services and Workers](../runbooks/deploying-services.md) - Deployment procedures
