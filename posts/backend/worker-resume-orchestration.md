# Resume Worker Orchestration Pattern

## Overview
How the resume worker orchestrates the resume generation workflow including AI service integration, database operations, and file generation.

## Key Points

### Worker Components
- Queue: Redis queue for job management
- Database: PostgreSQL connection for data retrieval
- AI Service: External service for resume content generation
- Translation Helper: Database connection for translation lookups
- Resume Generator: Main orchestration logic for PDF generation

### Initialization Flow
1. Connect to Redis queue
2. Connect to PostgreSQL database
3. Initialize AI service client (HTTP client)
4. Initialize translation helper
5. Initialize resume generator with all dependencies

### Processing Loop
- Continuous polling loop with configurable timeout (default 5 seconds)
- Graceful shutdown handling (KeyboardInterrupt, stop signal)
- Error resilience (continue processing after individual job failures)
- Status updates at each stage

### Job Processing Flow
1. Dequeue job from Redis
2. Update job status to "processing"
3. Extract job parameters (user_id, job_description, job_title, language)
4. Call ResumeGenerator.generate_resume()
5. Update job status to "completed" with result data
6. Handle errors with retry logic and status updates

### Resume Generation Orchestration
- Fetches user profile data from database
- Fetches projects, certifications, experiences
- Calls AI service for resume content generation
- Handles translations for multilingual resumes
- Generates PDF using templates
- Saves output files to configured directory
- Returns result metadata (file path, size, tags, duration)

### Error Handling & Retry Logic
- Error classification: transient vs permanent
- Exponential backoff calculation (2^retry_count, max 60s)
- Retry count tracking per job
- Max retries: 3 (configurable)
- Error type and message preservation
- Failed jobs tracked in job status

### Integration Points
- **Database**: User profiles, projects, certifications, experiences
- **AI Service**: HTTP REST API for content generation
- **Redis**: Job queue and status tracking
- **File System**: PDF output storage
- **Translation System**: Language-specific content retrieval

## Implementation Details

### Dependencies
- Queue: Redis connection
- Database: PostgreSQL with SQLAlchemy
- AI Service: HTTP client with retry logic
- Translation Helper: Database queries for translations
- Resume Generator: PDF generation with templates

### Resource Management
- Proper connection cleanup on shutdown
- Database connection pooling
- File handle management
- Graceful resource release

## Potential Improvements
- Add worker pool for parallel processing
- Implement job prioritization (urgent resumes first)
- Add rate limiting for AI service calls
- Implement circuit breaker for AI service failures
- Add health check endpoint
- Implement job caching (same user + job description)
- Add partial generation support (incremental updates)
- Implement job preemption/cancellation
- Add metrics collection (throughput, latency, error rate)
- Support job scheduling (delayed generation)
- Implement job result streaming (partial results)
- Add job dependency management (wait for translations)
- Implement job batching for efficiency
- Add resume preview generation
- Support template selection per job
- Implement A/B testing for different templates

