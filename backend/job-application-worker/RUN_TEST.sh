#!/bin/bash
# Quick test script - run inside Docker container

echo "🧪 Running Job Application Worker Tests"
echo "========================================"
echo ""

# Set environment variables
export DATABASE_URL="postgres://postgres:postgres@database:5432/woragis?sslmode=disable"
export REDIS_URL="redis://redis:6379/0"
export AI_SERVICE_URL="http://ai-service:8000"
export PLAYWRIGHT_HEADLESS="false"
export PLAYWRIGHT_SLOW_MO="500"

# Run test
cd /app
node test-all.js

