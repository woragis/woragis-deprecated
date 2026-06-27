# Resumes Domain - Resume Generation Flow

## Overview
How resume generation is orchestrated from request to PDF delivery.

## Key Points

### Generation Request Flow
1. User requests resume generation via API
2. Handler validates request (user_id, job_application_id, language)
3. Job created and enqueued to Redis queue
4. Response returned to user with job_id for status tracking
5. Worker picks up job asynchronously

### Worker Processing Flow
1. Worker dequeues job from Redis
2. Updates job status to "processing"
3. Fetches user profile data (experiences, projects, certifications)
4. Calls AI service for resume content generation
5. Handles translations for multilingual resumes
6. Generates PDF using templates
7. Saves PDF to file system
8. Updates job status to "completed" with result metadata

### Components Integration
- **Queue**: Redis queue for job management
- **Database**: User profiles, projects, certifications, experiences
- **AI Service**: Content generation via HTTP API
- **Translation Helper**: Language-specific content retrieval
- **Resume Generator**: PDF generation with templates
- **File System**: PDF storage

### Job Status Tracking
- Status transitions: pending → processing → completed/failed/retrying
- Result metadata: output_path, file_name, file_size, tags, duration_ms
- Error tracking: error message, error type, retry count

## Potential Improvements
- Add resume generation caching (same user + job description)
- Implement incremental resume updates
- Support multiple resume templates
- Add resume preview before full generation
- Implement resume versioning
- Add resume customization options
- Support A/B testing for different templates
- Add resume quality scoring
- Implement resume generation scheduling
- Support resume generation from templates
- Add resume generation analytics

