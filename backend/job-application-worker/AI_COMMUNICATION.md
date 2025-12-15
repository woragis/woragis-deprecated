# AI Service Communication Flow

## Overview

The job application worker communicates with the AI service in **two different ways**:

1. **Cover Letter Generation** - For creating personalized cover letters
2. **Selector Finding** - For self-healing scrapers when HTML selectors break

Both use HTTP REST API calls to the Go backend's AI service.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  Job Application Worker (Node.js)                          │
│                                                             │
│  ┌──────────────────┐    ┌──────────────────────────┐   │
│  │ CoverLetterService│    │  AISelectorFinder         │   │
│  └────────┬─────────┘    └──────────┬─────────────────┘   │
│           │                         │                       │
│           │ HTTP POST               │ HTTP POST             │
│           │ /api/chat/completions   │ /api/chat/completions │
│           │                         │                       │
└───────────┼─────────────────────────┼───────────────────────┘
            │                         │
            │                         │
            ▼                         ▼
┌─────────────────────────────────────────────────────────────┐
│  AI Service (Go Backend) - Port 8000                        │
│  http://ai-service:8000                                     │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  POST /api/chat/completions                         │   │
│  │  - Accepts OpenAI/Anthropic format                  │   │
│  │  - Routes to appropriate AI provider                │   │
│  │  - Returns completion response                      │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
            │
            ▼
┌─────────────────────────────────────────────────────────────┐
│  AI Providers (OpenAI, Anthropic, etc.)                     │
└─────────────────────────────────────────────────────────────┘
```

---

## 1. Cover Letter Generation

### Flow

```
Worker.processApplication()
    ↓
CoverLetterService.generateCoverLetter(profile, jobInfo)
    ↓
Build prompt with user profile + job info
    ↓
HTTP POST to AI Service
    ↓
AI Service → OpenAI/Anthropic
    ↓
Response: Generated cover letter text
    ↓
Return to worker → Use in job application
```

### Code Location
- **File**: `src/coverLetter.js`
- **Method**: `generateCoverLetter()`

### Request Details

**Endpoint**: `POST http://ai-service:8000/api/chat/completions`

**Request Body**:
```json
{
  "provider": "openai",
  "model": "gpt-4o-mini",
  "temperature": 0.7,
  "messages": [
    {
      "role": "user",
      "content": "You are a professional cover letter writer. Write a personalized cover letter for the following job application.\n\nJob Information:\n- Company: Google\n- Position: Senior Software Engineer\n..."
    }
  ],
  "max_tokens": 1500
}
```

**Response**:
```json
{
  "message": {
    "content": "Dear Hiring Manager,\n\nI am writing to express my interest in the Senior Software Engineer position at Google..."
  }
}
```

### Environment Variable
```bash
AI_SERVICE_URL=http://ai-service:8000
```

---

## 2. Selector Finding (Self-Healing)

### Flow

```
Scraper tries to find element
    ↓
Cached selector fails
    ↓
SelfHealingScraper.findElement()
    ↓
AISelectorFinder.findSelectorsFromHTML() or findSelectorsFromScreenshot()
    ↓
HTTP POST to AI Service (with HTML or screenshot)
    ↓
AI Service → OpenAI (gpt-4o-mini or gpt-4o for vision)
    ↓
Response: JSON with selectors
    ↓
Parse and cache selectors in Redis
    ↓
Use new selectors to find element
```

### Code Location
- **File**: `src/aiSelectorFinder.js`
- **Methods**: 
  - `findSelectorsFromHTML()` - For HTML analysis
  - `findSelectorsFromScreenshot()` - For vision analysis

### Request Details

#### HTML Analysis Request

**Endpoint**: `POST http://ai-service:8000/api/chat/completions`

**Request Body**:
```json
{
  "provider": "openai",
  "model": "gpt-4o-mini",
  "temperature": 0.3,
  "messages": [
    {
      "role": "user",
      "content": "You are a web scraping expert. Analyze the HTML below and find the best selectors for: 'Easy Apply button'\n\nWebsite: linkedin\nLooking for: Easy Apply button\n\nHTML:\n<html>...</html>\n\nInstructions:\n1. Find the element that matches the description\n2. Provide multiple selector strategies\n3. Return a JSON object..."
    }
  ],
  "max_tokens": 1000
}
```

**Response**:
```json
{
  "message": {
    "content": "{\n  \"primary\": \"button[data-testid='easy-apply']\",\n  \"alternatives\": [\"button:has-text('Easy Apply')\", \".jobs-apply-button\"],\n  \"xpath\": \"//button[contains(text(), 'Apply')]\",\n  \"text\": \"Easy Apply\",\n  \"explanation\": \"Button with data-testid attribute\"\n}"
  }
}
```

#### Vision Analysis Request (Screenshot)

**Endpoint**: `POST http://ai-service:8000/api/chat/completions`

**Request Body**:
```json
{
  "provider": "openai",
  "model": "gpt-4o",
  "temperature": 0.3,
  "messages": [
    {
      "role": "user",
      "content": [
        {
          "type": "text",
          "text": "You are a web scraping expert. Analyze this screenshot and find the best way to locate: 'Easy Apply button'..."
        },
        {
          "type": "image_url",
          "image_url": {
            "url": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
          }
        }
      ]
    }
  ],
  "max_tokens": 1000
}
```

**Response**: Same JSON structure as HTML analysis

---

## Communication Details

### Protocol
- **Protocol**: HTTP/HTTPS
- **Method**: POST
- **Content-Type**: application/json
- **Timeout**: 
  - Cover letters: 30 seconds
  - HTML analysis: 30 seconds
  - Vision analysis: 60 seconds

### Error Handling

```javascript
try {
  const response = await axios.post(url, data, { timeout: 30000 });
  // Process response
} catch (error) {
  logger.error('AI service call failed', {
    error: error.message,
    response: error.response?.data,
  });
  throw new Error(`AI service failed: ${error.message}`);
}
```

### Response Parsing

The worker handles two possible response formats:

```javascript
// Format 1: Direct message content
const content = response.data.message?.content;

// Format 2: OpenAI-style choices array
const content = response.data.choices?.[0]?.message?.content;
```

---

## Network Configuration

### Docker Compose

The worker and AI service are in the same Docker network:

```yaml
services:
  job-application-worker:
    environment:
      AI_SERVICE_URL: ${AI_SERVICE_URL:-http://ai-service:8000}
    depends_on:
      ai-service:
        condition: service_started

  ai-service:
    ports:
      - "8000:8000"
```

### Service Discovery

- **Container name**: `woragis-ai-service`
- **Internal URL**: `http://ai-service:8000` (Docker network)
- **External URL**: `http://localhost:8000` (if needed)

---

## Data Flow Example

### Complete Flow: Applying to a Job

```
1. Worker dequeues job from Redis
   ↓
2. Fetch user profile from PostgreSQL
   ↓
3. Call AI Service for cover letter
   POST /api/chat/completions
   → AI Service → OpenAI
   ← Generated cover letter
   ↓
4. Launch Playwright browser
   ↓
5. Navigate to job URL
   ↓
6. Try to find "Easy Apply" button
   → Check Redis cache for selectors
   → Cache miss or selector fails
   ↓
7. Call AI Service for selector finding
   POST /api/chat/completions (with HTML)
   → AI Service → OpenAI
   ← JSON with selectors
   ↓
8. Cache selectors in Redis (7-day TTL)
   ↓
9. Use selectors to click button
   ↓
10. Fill form fields (may use AI for each field)
   ↓
11. Submit application
   ↓
12. Save result to PostgreSQL
```

---

## Key Points

1. **Single AI Service**: Both cover letters and selector finding use the same AI service endpoint
2. **Different Models**: 
   - Cover letters: `gpt-4o-mini` (cheaper, faster)
   - HTML analysis: `gpt-4o-mini` (cheaper)
   - Vision analysis: `gpt-4o` (more expensive, for screenshots)
3. **Caching**: Selectors are cached in Redis to reduce AI calls
4. **Error Recovery**: If AI fails, the worker logs the error and marks the job as failed
5. **No Direct AI Provider Access**: The worker never calls OpenAI/Anthropic directly - it always goes through the Go AI service

---

## Environment Variables

```bash
# Required
AI_SERVICE_URL=http://ai-service:8000

# Optional (for debugging)
PLAYWRIGHT_HEADLESS=true
```

---

## Monitoring

The worker logs all AI interactions:

```javascript
logger.info('Generating cover letter', { company, jobTitle });
logger.info('Finding selectors using AI', { website, description });
logger.error('AI service call failed', { error, response });
```

Check logs with:
```bash
docker logs woragis-job-application-worker
```

