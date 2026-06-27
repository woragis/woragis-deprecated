# Service Boundaries: When to Split Services

## Overview
Decision criteria and patterns for when to split functionality into separate services vs keeping it in the main server.

## Key Points

### Service Split Criteria
- Independent scaling needs
- Different technology requirements
- Fault isolation
- Team ownership
- Deployment frequency

### Current Services
- **AI Service**: AI/LLM operations (Python, FastAPI)
- **Creative Service**: Image/video generation (Python, FastAPI)
- **Workers**: Async job processing (Go, Python, Node.js)

### What Stayed in Server
- Domain logic (7 domains)
- API endpoints
- Authentication
- Database operations

## Implementation Details

### Why AI Service Separate?
- Heavy computation
- Different scaling needs
- Python ecosystem (LangChain)
- Provider abstraction

### Why Creative Service Separate?
- Image/video processing
- Different scaling needs
- Python ecosystem
- Provider abstraction

### Why Workers Separate?
- Async processing
- Independent scaling
- Fault isolation
- Technology diversity

## Decision Framework
1. Does it need different scaling?
2. Does it need different technology?
3. Does it need fault isolation?
4. Does it have different deployment needs?

## Benefits
- Independent scaling
- Technology diversity
- Fault isolation
- Team autonomy

## Trade-offs
- More complexity
- Network overhead
- Deployment complexity
- Service coordination

## Lessons Learned
- Split when scaling needs differ
- Keep together when tightly coupled
- Technology needs justify split
- Fault isolation important

## Future Improvements
- Service mesh (if needed)
- API gateway
- Service discovery
- Centralized configuration
