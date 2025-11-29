#!/bin/bash

# Script to seed Golang projects
# Usage: ./seed-golang-projects.sh [email] [password]

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

# Helper function to extract ID from response
extract_id() {
    local response=$1
    echo "$response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4
}

# Get or create Golang skill
echo "Getting Golang skill..."
GOLANG_SKILL_RESPONSE=$(curl -s -X GET "$API_BASE/skills" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H "Content-Type: application/json")

GOLANG_SKILL_ID=$(echo "$GOLANG_SKILL_RESPONSE" | grep -o '"id":"[^"]*"[^}]*"name":"Golang"' | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$GOLANG_SKILL_ID" ]; then
    echo "Golang skill not found. Please create it first using seed-skills.sh"
    exit 1
fi

echo "Found Golang skill ID: $GOLANG_SKILL_ID"
echo ""

# ==========================================
# PROJECT 1: High-Performance API Gateway
# ==========================================
echo "Creating Project 1: High-Performance API Gateway..."

PROJECT_1_PAYLOAD='{
    "name": "High-Performance API Gateway",
    "description": "Built a production-ready API gateway in Go handling 50K+ requests per second with sub-millisecond latency. Implemented rate limiting, authentication middleware, request routing, and circuit breakers. Uses goroutines for concurrent request processing and channels for inter-service communication.",
    "status": "completed",
    "healthScore": 98,
    "mrr": 0,
    "cac": 0,
    "ltv": 0,
    "churnRate": 0
}'

PROJECT_1_RESPONSE=$(curl -s -X POST "$API_BASE/projects" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d "$PROJECT_1_PAYLOAD")

PROJECT_1_ID=$(extract_id "$PROJECT_1_RESPONSE")

if [ -z "$PROJECT_1_ID" ]; then
    echo "Failed to create project 1"
else
    echo "✓ Project 1 created: $PROJECT_1_ID"
    
    # Add Golang technology to project
    TECH_PAYLOAD='{
        "name": "Golang",
        "version": "1.21+",
        "category": "backend",
        "purpose": "Primary programming language for high-performance backend services",
        "link": "https://go.dev/"
    }'
    
    curl -s -X POST "$API_BASE/projects/$PROJECT_1_ID/technologies" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d "$TECH_PAYLOAD" > /dev/null
    echo "✓ Added Golang technology to project"
    
    # Link Golang skill to project (if endpoint exists)
    curl -s -X POST "$API_BASE/projects/$PROJECT_1_ID/skills/$GOLANG_SKILL_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null 2>&1
fi

echo ""

# ==========================================
# PROJECT 2: Distributed Task Queue System
# ==========================================
echo "Creating Project 2: Distributed Task Queue System..."

PROJECT_2_PAYLOAD='{
    "name": "Distributed Task Queue System",
    "description": "Developed a distributed task queue system using Go channels, Redis for job persistence, and worker pools. Handles millions of background jobs with guaranteed delivery, retry mechanisms, and priority queuing. Features horizontal scaling and fault tolerance.",
    "status": "completed",
    "healthScore": 95,
    "mrr": 0,
    "cac": 0,
    "ltv": 0,
    "churnRate": 0
}'

PROJECT_2_RESPONSE=$(curl -s -X POST "$API_BASE/projects" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d "$PROJECT_2_PAYLOAD")

PROJECT_2_ID=$(extract_id "$PROJECT_2_RESPONSE")

if [ -z "$PROJECT_2_ID" ]; then
    echo "Failed to create project 2"
else
    echo "✓ Project 2 created: $PROJECT_2_ID"
    
    # Add Golang technology to project
    TECH_PAYLOAD='{
        "name": "Golang",
        "version": "1.21+",
        "category": "backend",
        "purpose": "Primary programming language for distributed systems",
        "link": "https://go.dev/"
    }'
    
    curl -s -X POST "$API_BASE/projects/$PROJECT_2_ID/technologies" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d "$TECH_PAYLOAD" > /dev/null
    echo "✓ Added Golang technology to project"
    
    # Link Golang skill to project (if endpoint exists)
    curl -s -X POST "$API_BASE/projects/$PROJECT_2_ID/skills/$GOLANG_SKILL_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null 2>&1
fi

echo ""

# ==========================================
# PROJECT 3: Microservices Orchestration Platform
# ==========================================
echo "Creating Project 3: Microservices Orchestration Platform..."

PROJECT_3_PAYLOAD='{
    "name": "Microservices Orchestration Platform",
    "description": "Architected and implemented a microservices orchestration platform in Go using gRPC for inter-service communication. Features service discovery, load balancing, health checks, and distributed tracing. Supports 20+ microservices with zero-downtime deployments.",
    "status": "completed",
    "healthScore": 97,
    "mrr": 0,
    "cac": 0,
    "ltv": 0,
    "churnRate": 0
}'

PROJECT_3_RESPONSE=$(curl -s -X POST "$API_BASE/projects" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d "$PROJECT_3_PAYLOAD")

PROJECT_3_ID=$(extract_id "$PROJECT_3_RESPONSE")

if [ -z "$PROJECT_3_ID" ]; then
    echo "Failed to create project 3"
else
    echo "✓ Project 3 created: $PROJECT_3_ID"
    
    # Add Golang technology to project
    TECH_PAYLOAD='{
        "name": "Golang",
        "version": "1.21+",
        "category": "backend",
        "purpose": "Core language for microservices architecture",
        "link": "https://go.dev/"
    }'
    
    curl -s -X POST "$API_BASE/projects/$PROJECT_3_ID/technologies" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d "$TECH_PAYLOAD" > /dev/null
    echo "✓ Added Golang technology to project"
    
    # Link Golang skill to project (if endpoint exists)
    curl -s -X POST "$API_BASE/projects/$PROJECT_3_ID/skills/$GOLANG_SKILL_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null 2>&1
fi

echo ""

# ==========================================
# PROJECT 4: Real-time WebSocket Server
# ==========================================
echo "Creating Project 4: Real-time WebSocket Server..."

PROJECT_4_PAYLOAD='{
    "name": "Real-time WebSocket Server",
    "description": "Built a high-performance WebSocket server in Go supporting 100K+ concurrent connections. Implements connection pooling, message broadcasting, room management, and graceful shutdown. Uses goroutines for handling each connection efficiently with minimal memory footprint.",
    "status": "completed",
    "healthScore": 96,
    "mrr": 0,
    "cac": 0,
    "ltv": 0,
    "churnRate": 0
}'

PROJECT_4_RESPONSE=$(curl -s -X POST "$API_BASE/projects" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d "$PROJECT_4_PAYLOAD")

PROJECT_4_ID=$(extract_id "$PROJECT_4_RESPONSE")

if [ -z "$PROJECT_4_ID" ]; then
    echo "Failed to create project 4"
else
    echo "✓ Project 4 created: $PROJECT_4_ID"
    
    # Add Golang technology to project
    TECH_PAYLOAD='{
        "name": "Golang",
        "version": "1.21+",
        "category": "backend",
        "purpose": "High-performance WebSocket server implementation",
        "link": "https://go.dev/"
    }'
    
    curl -s -X POST "$API_BASE/projects/$PROJECT_4_ID/technologies" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d "$TECH_PAYLOAD" > /dev/null
    echo "✓ Added Golang technology to project"
    
    # Link Golang skill to project (if endpoint exists)
    curl -s -X POST "$API_BASE/projects/$PROJECT_4_ID/skills/$GOLANG_SKILL_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null 2>&1
fi

echo ""

# ==========================================
# PROJECT 5: Blockchain Transaction Processor
# ==========================================
echo "Creating Project 5: Blockchain Transaction Processor..."

PROJECT_5_PAYLOAD='{
    "name": "Blockchain Transaction Processor",
    "description": "Developed a blockchain transaction processing engine in Go with cryptographic validation, consensus algorithms, and distributed ledger management. Processes 10K+ transactions per second with cryptographic integrity guarantees. Implements Merkle trees and proof-of-work mechanisms.",
    "status": "completed",
    "healthScore": 94,
    "mrr": 0,
    "cac": 0,
    "ltv": 0,
    "churnRate": 0
}'

PROJECT_5_RESPONSE=$(curl -s -X POST "$API_BASE/projects" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d "$PROJECT_5_PAYLOAD")

PROJECT_5_ID=$(extract_id "$PROJECT_5_RESPONSE")

if [ -z "$PROJECT_5_ID" ]; then
    echo "Failed to create project 5"
else
    echo "✓ Project 5 created: $PROJECT_5_ID"
    
    # Add Golang technology to project
    TECH_PAYLOAD='{
        "name": "Golang",
        "version": "1.21+",
        "category": "backend",
        "purpose": "Blockchain and cryptographic operations",
        "link": "https://go.dev/"
    }'
    
    curl -s -X POST "$API_BASE/projects/$PROJECT_5_ID/technologies" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d "$TECH_PAYLOAD" > /dev/null
    echo "✓ Added Golang technology to project"
    
    # Link Golang skill to project (if endpoint exists)
    curl -s -X POST "$API_BASE/projects/$PROJECT_5_ID/skills/$GOLANG_SKILL_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null 2>&1
fi

echo ""

echo "=========================================="
echo "All Golang projects created successfully!"
echo "=========================================="
echo ""
echo "Projects created:"
echo "  1. High-Performance API Gateway"
echo "  2. Distributed Task Queue System"
echo "  3. Microservices Orchestration Platform"
echo "  4. Real-time WebSocket Server"
echo "  5. Blockchain Transaction Processor"
echo ""

