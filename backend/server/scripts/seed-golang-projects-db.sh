#!/bin/bash

# Direct database script to seed Golang projects
# This bypasses the API and inserts directly into the database

USER_ID="${1:-6ad0d828-f605-45fc-a545-3441e17a015c}"

echo "Seeding Golang projects directly to database..."
echo "User ID: $USER_ID"
echo ""

# Get Golang skill ID
GOLANG_SKILL_ID=$(docker exec woragis-database psql -U postgres -d woragis -t -c "SELECT id FROM skills WHERE name = 'Golang' LIMIT 1;" | tr -d ' ')

if [ -z "$GOLANG_SKILL_ID" ]; then
    echo "Error: Golang skill not found in database"
    exit 1
fi

echo "Found Golang skill ID: $GOLANG_SKILL_ID"
echo ""

# Function to create project and link technologies/skills
create_golang_project() {
    local name="$1"
    local description="$2"
    local slug=$(echo "$name" | tr '[:upper:]' '[:lower:]' | tr ' ' '-' | sed 's/[^a-z0-9-]//g')
    
    echo "Creating project: $name"
    
    # Insert project
    PROJECT_ID=$(docker exec woragis-database psql -U postgres -d woragis -t -c "
        INSERT INTO projects (id, user_id, name, description, slug, status, health_score, mrr, cac, ltv, churn_rate, created_at, updated_at)
        VALUES (
            gen_random_uuid(),
            '$USER_ID',
            '$name',
            '$description',
            '$slug',
            'completed',
            95,
            0,
            0,
            0,
            0,
            NOW(),
            NOW()
        )
        RETURNING id;
    " | tr -d ' ')
    
    if [ -z "$PROJECT_ID" ]; then
        echo "Failed to create project: $name"
        return 1
    fi
    
    echo "✓ Project created: $PROJECT_ID"
    
    # Add Golang technology (if table exists)
    docker exec woragis-database psql -U postgres -d woragis -q -c "
        INSERT INTO project_technologies (id, project_id, name, version, category, purpose, link, created_at, updated_at)
        VALUES (
            gen_random_uuid(),
            '$PROJECT_ID',
            'Golang',
            '1.21+',
            'backend',
            'Primary programming language',
            'https://go.dev/',
            NOW(),
            NOW()
        );
    " > /dev/null 2>&1 && echo "✓ Added Golang technology" || echo "  (Skipped technology - table may not exist)"
    
    # Link Golang skill
    docker exec woragis-database psql -U postgres -d woragis -q -c "
        INSERT INTO project_skills (project_id, skill_id, created_at)
        VALUES ('$PROJECT_ID', '$GOLANG_SKILL_ID', NOW())
        ON CONFLICT DO NOTHING;
    " > /dev/null 2>&1
    
    echo "✓ Linked Golang skill"
    echo ""
}

# Create projects
create_golang_project \
    "High-Performance API Gateway" \
    "Production-ready API gateway in Go handling 50K+ req/s with sub-millisecond latency. Features rate limiting, auth middleware, routing, and circuit breakers using goroutines and channels."

create_golang_project \
    "Distributed Task Queue System" \
    "Distributed task queue using Go channels, Redis, and worker pools. Handles millions of background jobs with guaranteed delivery, retries, and priority queuing. Horizontally scalable."

create_golang_project \
    "Microservices Orchestration Platform" \
    "Microservices platform in Go using gRPC. Features service discovery, load balancing, health checks, and distributed tracing. Supports 20+ microservices with zero-downtime deployments."

create_golang_project \
    "Real-time WebSocket Server" \
    "High-performance WebSocket server in Go supporting 100K+ concurrent connections. Features connection pooling, message broadcasting, and graceful shutdown using goroutines."

create_golang_project \
    "Blockchain Transaction Processor" \
    "Blockchain transaction engine in Go with cryptographic validation and consensus algorithms. Processes 10K+ tps with integrity guarantees. Implements Merkle trees and proof-of-work."

echo "=========================================="
echo "All Golang projects created successfully!"
echo "=========================================="

