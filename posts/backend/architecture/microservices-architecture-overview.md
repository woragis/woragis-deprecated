# Microservices Architecture: Server, Services, Workers Overview

## Overview
High-level overview of the Woragis backend microservices architecture, including the main server, services (AI, Creative), and workers (Email, WhatsApp, Translation, Resume, Job Application).

## Key Points

### Architecture Components
- **Server**: Main API server (Go, Fiber framework)
- **Services**: AI Service (Python, FastAPI), Creative Service (Python, FastAPI)
- **Workers**: Email, WhatsApp, Translation, Resume, Job Application workers
- **Message Queue**: RabbitMQ (primary), Redis (fallback)
- **Database**: PostgreSQL
- **Cache**: Redis

### Communication Patterns
- Server ↔ Services: HTTP/REST
- Server ↔ Workers: Message queues (RabbitMQ/Redis)
- Workers ↔ Services: HTTP/REST (for AI/Creative services)
- All ↔ Database: Direct PostgreSQL connections

### Data Flow
- API requests → Server → Domain logic → Queue/Response
- Queue → Worker → Process → Database/External APIs
- Services → External APIs (OpenAI, Anthropic, etc.)

## Implementation Details

### Server Architecture
- Domain-driven design (7 domains)
- RESTful API endpoints
- JWT authentication
- Structured logging
- Health checks

### Services Architecture
- FastAPI-based microservices
- Provider abstraction (OpenAI, Anthropic, etc.)
- Agent system (AI Service)
- Image/Video generation (Creative Service)

### Workers Architecture
- Standalone processes
- RabbitMQ consumers
- Direct database writes
- Health check endpoints

## Visual Diagram
[Architecture diagram showing: Server (center), Services (AI, Creative), Workers (5 workers), RabbitMQ, Redis, PostgreSQL]

## Benefits
- Independent scaling
- Technology diversity (Go, Python, Node.js)
- Fault isolation
- Independent deployment

## Challenges
- Increased complexity
- Distributed system challenges
- Service coordination
- Data consistency

## Future Improvements
- Service mesh (if needed)
- API gateway
- Centralized configuration
- Distributed tracing
