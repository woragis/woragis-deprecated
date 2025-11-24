#!/bin/bash

# Script to seed technical content (case studies, system designs, problem solutions)
# Usage: ./seed-technical-content.sh [email] [password]

API_BASE="http://localhost:8080/api"
EMAIL="${1:-masteringthecode.woragis@gmail.com}"
PASSWORD="${2}"

if [ -z "$PASSWORD" ]; then
    echo "Usage: $0 [email] [password]"
    echo "Please provide your password to authenticate"
    exit 1
fi

echo "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

# Extract access_token from JSON response
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
if [ -z "$ACCESS_TOKEN" ]; then
    ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"accessToken":"[^"]*' | cut -d'"' -f4)
fi

if [ -z "$ACCESS_TOKEN" ]; then
    echo "Failed to login. Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "Login successful! Token: ${ACCESS_TOKEN:0:20}..."

# Get or create a project first (needed for case studies)
echo "Fetching projects..."
PROJECTS_RESPONSE=$(curl -s -X GET "$API_BASE/projects" \
    -H "Authorization: Bearer $ACCESS_TOKEN")

# Try to extract first project ID, or create one
PROJECT_ID=$(echo $PROJECTS_RESPONSE | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)

if [ -z "$PROJECT_ID" ]; then
    echo "No projects found. Creating a project..."
    PROJECT_RESPONSE=$(curl -s -X POST "$API_BASE/projects" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d '{
            "name": "Woragis Backend",
            "description": "Scalable backend system with domain-driven design",
            "status": "executing",
            "health_score": 85,
            "mrr": 0,
            "cac": 0,
            "ltv": 0,
            "churn_rate": 0
        }')
    
    PROJECT_ID=$(echo $PROJECT_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    if [ -z "$PROJECT_ID" ]; then
        echo "Failed to create project. Response: $PROJECT_RESPONSE"
        exit 1
    fi
    echo "Created project: $PROJECT_ID"
else
    echo "Using existing project: $PROJECT_ID"
fi

# ============================================================================
# PROJECT CASE STUDIES
# ============================================================================
echo ""
echo "Creating project case studies..."

# Case Study 1: Woragis Backend Architecture
curl -s -X POST "$API_BASE/projects/$PROJECT_ID/case-studies" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "title": "Woragis Backend Architecture",
        "description": "Built a scalable backend system using Go, Fiber, and PostgreSQL with microservices architecture supporting multiple domains including projects, finances, skills, and AI integrations.",
        "challenge": "Need to build a robust, scalable backend that can handle multiple business domains while maintaining clean architecture, API key authentication, and real-time capabilities.",
        "solution": "Implemented a domain-driven design with separate modules for each business domain, JWT and API key authentication, Redis for caching and pub/sub, and a clean separation of concerns with repositories, services, and handlers.",
        "technologies": ["Go", "Fiber", "PostgreSQL", "Redis", "Docker", "Kubernetes"],
        "architecture": "Microservices-ready architecture with domain-driven design. Each domain (projects, skills, finances, etc.) has its own repository, service, and handler layers. API key authentication for public read access, JWT for authenticated operations.",
        "metrics": {
            "metrics": [
                {"label": "Response Time", "value": "< 50ms", "improvement": "P95 latency"},
                {"label": "Uptime", "value": "99.9%", "improvement": "Target SLA"},
                {"label": "API Endpoints", "value": "50+", "improvement": "RESTful APIs"}
            ]
        },
        "tradeoffs": {
            "tradeoffs": [
                {
                    "decision": "API Key vs JWT for public endpoints",
                    "pros": [
                        "Simpler for public consumption",
                        "No token expiration management",
                        "Better for static site integration"
                    ],
                    "cons": [
                        "Less granular permissions",
                        "Requires separate key management"
                    ]
                },
                {
                    "decision": "Monolithic vs Microservices",
                    "pros": [
                        "Easier development and deployment",
                        "Simpler debugging",
                        "Lower operational complexity"
                    ],
                    "cons": [
                        "Potential scaling bottlenecks",
                        "Shared database can become a bottleneck"
                    ]
                }
            ]
        },
        "lessonsLearned": [
            "Domain-driven design makes code more maintainable and testable",
            "API key authentication is essential for public-facing APIs",
            "Redis pub/sub enables real-time features without WebSocket complexity",
            "Clean architecture pays off in long-term maintainability"
        ]
    }' > /dev/null

echo "Created case study: Woragis Backend Architecture"

# ============================================================================
# SYSTEM DESIGNS
# ============================================================================
echo ""
echo "Creating system designs..."

# System Design 1: RESTful API Architecture
curl -s -X POST "$API_BASE/system-designs" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "title": "RESTful API Architecture",
        "description": "Clean, scalable API architecture with middleware-based authentication, CORS support, and domain-driven design.",
        "components": {
            "components": [
                {
                    "name": "API Gateway",
                    "description": "Fiber router with CORS, logging, and recovery middleware",
                    "technology": "Go Fiber"
                },
                {
                    "name": "Authentication Layer",
                    "description": "JWT for authenticated users, API keys for public read access",
                    "technology": "JWT, API Keys"
                },
                {
                    "name": "Domain Services",
                    "description": "Business logic separated by domain (projects, skills, finances)",
                    "technology": "Go"
                },
                {
                    "name": "Data Layer",
                    "description": "Repository pattern with GORM for database access",
                    "technology": "GORM, PostgreSQL"
                },
                {
                    "name": "Cache Layer",
                    "description": "Redis for caching and pub/sub messaging",
                    "technology": "Redis"
                }
            ]
        },
        "dataFlow": "Request → CORS Middleware → Auth Middleware → Domain Handler → Service → Repository → Database. Responses cached in Redis when appropriate.",
        "scalability": "Horizontal scaling ready with stateless services. Database connection pooling and Redis clustering support. Kubernetes-ready deployment.",
        "reliability": "Error recovery middleware, structured logging, health checks, and graceful shutdown. Database migrations and rollback support.",
        "featured": true
    }' > /dev/null

echo "Created system design: RESTful API Architecture"

# System Design 2: Dual Authentication System
curl -s -X POST "$API_BASE/system-designs" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "title": "Dual Authentication System",
        "description": "Flexible authentication supporting both JWT tokens for authenticated users and API keys for public read-only access.",
        "components": {
            "components": [
                {
                    "name": "JWT Authentication",
                    "description": "Token-based auth for authenticated operations (POST, PATCH, DELETE)",
                    "technology": "JWT, Go"
                },
                {
                    "name": "API Key Authentication",
                    "description": "Key-based auth for public read operations (GET)",
                    "technology": "SHA256 Hashing"
                },
                {
                    "name": "Middleware Chain",
                    "description": "RequireAPIKeyOrAuth middleware checks API key first, falls back to JWT",
                    "technology": "Go Fiber Middleware"
                },
                {
                    "name": "Context Storage",
                    "description": "User ID and API key stored in Fiber context for downstream handlers",
                    "technology": "Fiber Context"
                }
            ]
        },
        "dataFlow": "Request → Extract API Key/JWT → Validate → Store in Context → Handler → Service. API keys validated via SHA256 hash comparison.",
        "scalability": "Stateless authentication allows horizontal scaling. API keys can be rate-limited per key. JWT tokens validated without database lookup.",
        "reliability": "Graceful fallback from API key to JWT. Comprehensive error handling and logging. Secure key storage with hashing.",
        "featured": true
    }' > /dev/null

echo "Created system design: Dual Authentication System"

# System Design 3: Translation Pipeline Architecture
curl -s -X POST "$API_BASE/system-designs" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "title": "Translation Pipeline Architecture",
        "description": "Asynchronous translation system using Redis queues and AI services to provide multi-language support for dynamic content.",
        "components": {
            "components": [
                {
                    "name": "Translation Queue",
                    "description": "Redis-based queue for managing translation jobs",
                    "technology": "Redis Streams"
                },
                {
                    "name": "AI Translation Service",
                    "description": "LangChain integration for translating content using LLM",
                    "technology": "LangChain, OpenAI"
                },
                {
                    "name": "Translation Worker",
                    "description": "Background worker processing translation jobs from queue",
                    "technology": "Go, Redis"
                },
                {
                    "name": "Translation Enricher",
                    "description": "Middleware that enriches API responses with translated content based on Accept-Language header",
                    "technology": "Go, PostgreSQL"
                },
                {
                    "name": "Translation Storage",
                    "description": "PostgreSQL table storing translations with entity type, language, and field mappings",
                    "technology": "PostgreSQL, JSONB"
                }
            ]
        },
        "dataFlow": "Entity Created → Queue Translation Job → Worker Processes → AI Translation → Store in DB → Enricher Applies on GET Requests",
        "scalability": "Horizontal scaling with multiple workers. Redis queue handles high throughput. Database indexes on entity_type, entity_id, language for fast lookups.",
        "reliability": "Retry logic for failed translations. Fallback to source language if translation unavailable. Comprehensive error logging and monitoring.",
        "featured": false
    }' > /dev/null

echo "Created system design: Translation Pipeline Architecture"

# ============================================================================
# PROBLEM SOLUTIONS
# ============================================================================
echo ""
echo "Creating problem solutions..."

# Problem Solution 1: CORS Preflight Issues
curl -s -X POST "$API_BASE/problem-solutions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "problem": "CORS preflight requests failing for custom headers",
        "context": "Frontend landing page could not make API requests with X-API-Key header due to CORS preflight failures. Browser normalizes custom headers to lowercase for preflight checks.",
        "solution": "Updated CORS configuration to accept both X-API-Key and x-api-key in allowed headers. Browsers send lowercase headers in preflight, but the actual request can use either case.",
        "technologies": ["CORS", "Go Fiber", "HTTP Headers"],
        "impact": "Enabled seamless API integration from static landing page without CORS errors. Public API access now works reliably.",
        "metrics": {
            "before": "CORS errors blocking all requests",
            "after": "100% successful preflight requests",
            "improvement": "Zero CORS-related failures"
        },
        "featured": true
    }' > /dev/null

echo "Created problem solution: CORS Preflight Issues"

# Problem Solution 2: API Key Middleware
curl -s -X POST "$API_BASE/problem-solutions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "problem": "API key authentication not working for GET requests",
        "context": "Needed to allow public read access to projects and skills endpoints using API keys, while maintaining JWT authentication for write operations.",
        "solution": "Created RequireAPIKeyOrAuth middleware that checks for API key first on GET requests, then falls back to JWT. Ensured middleware order so API key routes are registered before protected routes.",
        "technologies": ["Go", "Fiber Middleware", "Authentication"],
        "impact": "Landing page can now fetch projects and skills data without requiring user authentication. Clean separation between public read and authenticated write operations.",
        "metrics": {
            "before": "All endpoints required JWT authentication",
            "after": "GET endpoints support API keys, write operations require JWT",
            "improvement": "Flexible authentication model"
        },
        "featured": true
    }' > /dev/null

echo "Created problem solution: API Key Middleware"

# Problem Solution 3: Domain Architecture
curl -s -X POST "$API_BASE/problem-solutions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "problem": "Code organization and maintainability as system grows",
        "context": "Multiple business domains (projects, skills, finances, etc.) needed clean separation to maintain code quality and enable team collaboration.",
        "solution": "Implemented domain-driven design with each domain having its own package structure: repository (data access), service (business logic), handler (HTTP layer), and routes. Clear boundaries and interfaces.",
        "technologies": ["Go", "Domain-Driven Design", "Clean Architecture"],
        "impact": "Code is more maintainable, testable, and scalable. New domains can be added without affecting existing ones. Clear ownership and boundaries.",
        "metrics": {
            "before": "Monolithic structure with mixed concerns",
            "after": "8+ independent domains with clear boundaries",
            "improvement": "Modular, maintainable architecture"
        },
        "featured": true
    }' > /dev/null

echo "Created problem solution: Domain Architecture"

# Problem Solution 4: JSONB Data Storage
curl -s -X POST "$API_BASE/problem-solutions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "problem": "Storing complex nested data structures in relational database",
        "context": "Need to store arrays and nested objects (metrics, tradeoffs, components) for case studies and system designs. Traditional normalized approach would require many join tables and complex queries.",
        "solution": "Used PostgreSQL JSONB columns with custom GORM Value() and Scan() methods. This allows storing complex nested structures while maintaining queryability and type safety in Go code.",
        "technologies": ["PostgreSQL", "JSONB", "GORM", "Go"],
        "impact": "Simplified data model, faster queries for complex structures, and easier to evolve schema. Can still query JSONB fields when needed using PostgreSQL operators.",
        "metrics": {
            "before": "Multiple join tables, complex queries, harder to maintain",
            "after": "Single JSONB column, simple queries, easy to extend",
            "improvement": "Reduced complexity, improved performance"
        },
        "featured": false
    }' > /dev/null

echo "Created problem solution: JSONB Data Storage"

# Problem Solution 5: Translation Performance
curl -s -X POST "$API_BASE/problem-solutions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "problem": "Translation lookups causing slow API responses",
        "context": "Initially fetching translations synchronously on every request caused 200-300ms latency. With multiple languages and entities, this became a bottleneck.",
        "solution": "Implemented translation enricher that batches lookups and uses indexed queries. Added caching layer for frequently accessed translations. Made translation fetching non-blocking where possible.",
        "technologies": ["PostgreSQL", "Indexing", "Caching", "Go"],
        "impact": "Reduced translation lookup overhead from 200ms to <10ms. API responses remain fast even with multiple language support.",
        "metrics": {
            "before": "200-300ms translation lookup overhead",
            "after": "<10ms with indexing and caching",
            "improvement": "20-30x performance improvement"
        },
        "featured": false
    }' > /dev/null

echo "Created problem solution: Translation Performance"

# Problem Solution 6: Redis Connection Pooling
curl -s -X POST "$API_BASE/problem-solutions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "problem": "Redis connection exhaustion under high load",
        "context": "Creating new Redis connections for each request caused connection pool exhaustion and timeouts during traffic spikes.",
        "solution": "Implemented connection pooling with proper connection limits, idle timeout, and connection reuse. Used Redis client with built-in pooling and configured appropriate pool sizes based on expected load.",
        "technologies": ["Redis", "Connection Pooling", "Go"],
        "impact": "Eliminated connection timeouts, improved throughput, and reduced connection overhead. System now handles 10x more concurrent requests.",
        "metrics": {
            "before": "Connection timeouts during traffic spikes",
            "after": "Stable connections, no timeouts",
            "improvement": "10x improvement in concurrent request handling"
        },
        "featured": false
    }' > /dev/null

echo "Created problem solution: Redis Connection Pooling"

echo ""
echo "All technical content seeded successfully!"
echo ""
echo "Summary:"
echo "  - 1 Project Case Study created"
echo "  - 3 System Designs created"
echo "  - 6 Problem Solutions created"

