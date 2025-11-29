# Resume Generation Worker

A Python-based worker service that generates customized resumes (CVs) based on job requirements. It uses AI to generate resume sections, searches through projects by technology tags, and includes relevant certifications.

## Features

- **AI-Powered Content Generation**: Uses the AI service to generate professional resume sections (Profile, About, Experience, Skills)
- **Smart Project Matching**: Searches projects by technology categories/tags that match job requirements
- **Certification Filtering**: Automatically includes relevant certifications based on job focus areas
- **PDF Generation**: Creates professional PDF resumes using WeasyPrint with customizable HTML/CSS templates
- **ATS-Friendly**: Generates resumes optimized for Applicant Tracking Systems
- **Result Tracking**: Saves generation metadata and results to JSON files

## Architecture

The codebase is organized into separate modules:

- **`database.py`**: Database connection and query operations
- **`ai_service.py`**: AI service client for generating resume content
- **`keyword_extractor.py`**: Extracts keywords and matches categories from job descriptions
- **`resume_generator.py`**: Main orchestration logic for resume generation
- **`main.py`**: Entry point and CLI interface

## File Storage

### Generated PDFs
- **Location**: `./resume-worker/output/` (on your local machine)
- **Naming**: `resume_<job_title>_<timestamp>.pdf`
- **Access**: Files are saved directly to your PC in the `resume-worker/output/` directory
- **Example**: `backend/server/resume-worker/output/resume_Backend_Engineer_20241129_120000.pdf`

### Result Metadata
- **Location**: `./resume-worker/results/` (on your local machine)
- **Format**: JSON files with generation metadata
- **Naming**: `resume_result_<user_id>_<timestamp>.json`
- **Contents**: Includes file path, size, projects/certifications count, keywords, etc.
- **Example**: `backend/server/resume-worker/results/resume_result_550e8400_20241129_120000.json`

### Local Directories
Files are saved directly to your local filesystem using bind mounts:
- `./resume-worker/output/`: Contains all generated PDF files (accessible on your PC)
- `./resume-worker/results/`: Contains JSON metadata for each generation (accessible on your PC)

**Note**: These directories are created automatically when the container starts. They are relative to the `docker-compose.yml` file location.

## Usage

### Command Line

```bash
python src/main.py <user_id> <job_description> [job_title] [output_filename]
```

Example:
```bash
python src/main.py "550e8400-e29b-41d4-a716-446655440000" "Looking for a backend engineer with Golang, Kubernetes, and microservices experience..." "Backend Engineer" "resume_backend.pdf"
```

### Docker

The worker is configured in `docker-compose.yml` and runs automatically when the stack is started.

To access generated files:
```bash
# List generated PDFs (from your PC)
ls -lh resume-worker/output/

# Or from inside container
docker exec woragis-resume-worker ls -lh /app/output/

# View result metadata (from your PC)
cat resume-worker/results/resume_result_*.json

# Or from inside container
docker exec woragis-resume-worker cat /app/results/resume_result_*.json
```

**Note**: Files are saved directly to `resume-worker/output/` and `resume-worker/results/` on your local machine, so you can access them directly without using Docker commands.

## Configuration

Environment variables:

- `DATABASE_URL`: PostgreSQL connection string
- `AI_SERVICE_URL`: URL of the AI service (default: `http://ai-service:8000`)
- `RESUME_OUTPUT_DIR`: Directory for generated PDFs (default: `/app/output`)
- `RESULTS_LOG_DIR`: Directory for result metadata JSON files (default: `/app/results`)

## Technology Matching

The worker matches projects to job requirements using technology categories:

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

## Resume Sections

1. **Professional Summary**: AI-generated 4-5 line summary optimized for the job
2. **About Me**: Brief personal introduction (3-4 sentences)
3. **Technical Skills**: Categorized list of relevant technical skills
4. **Projects & Experience**: Relevant projects with technologies and descriptions
5. **Certifications**: Relevant certifications with issue dates

## Template Customization

The resume template is located in `templates/resume.html` and `templates/style.css`. You can customize:

- Fonts and typography
- Colors and styling
- Layout (two-column sections, spacing)
- Section ordering
- Additional sections

## Output

### PDF Files
Generated PDFs are saved to the `RESUME_OUTPUT_DIR` (default: `/app/output`). The filename format is:

```
resume_<job_title>_<timestamp>.pdf
```

### Result Metadata
Each generation creates a JSON file with metadata:

```json
{
  "output_path": "/app/output/resume_Backend_Engineer_20241129_120000.pdf",
  "filename": "resume_Backend_Engineer_20241129_120000.pdf",
  "file_size": 123456,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "job_title": "Backend Engineer",
  "projects_count": 5,
  "certifications_count": 3,
  "generated_at": "2024-11-29T12:00:00",
  "keywords": {
    "tech_categories": ["backend", "devops"],
    "cert_categories": ["cloud", "devops"]
  }
}
```

## Development

### Local Setup

1. Install dependencies:
```bash
pip install -r requirements.txt
```

2. Set environment variables:
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/woragis?sslmode=disable"
export AI_SERVICE_URL="http://localhost:8000"
export RESUME_OUTPUT_DIR="./output"
export RESULTS_LOG_DIR="./results"
```

3. Run the worker:
```bash
python src/main.py <user_id> <job_description> [job_title]
```

### Docker Build

```bash
docker build -f Dockerfile.resume-worker -t woragis-resume-worker .
```

## Integration

The worker can be integrated with:

- **Job Application System**: Automatically generate resumes for job applications
- **Queue System**: Process resume generation requests from a queue (Redis, RabbitMQ, etc.)
- **API Endpoint**: Expose as an API endpoint for on-demand resume generation

## Notes

- The worker requires access to the PostgreSQL database and AI service
- Generated PDFs are ATS-friendly (no images, standard fonts, proper structure)
- The AI service should be configured with Anthropic (Claude) for best results
- Projects are filtered by technology categories, not exact keyword matching
- Only active certifications are included in the resume
- All generated files and metadata are persisted in Docker volumes
