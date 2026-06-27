#!/bin/bash

# Health check test script for all services
# This script tests the health endpoints of all services

echo "Testing health checks for all services..."
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to test health endpoint
test_health() {
    local service_name=$1
    local url=$2
    
    echo -n "Testing $service_name... "
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null)
    
    if [ "$response" = "200" ]; then
        echo -e "${GREEN}✓ Healthy (HTTP $response)${NC}"
        return 0
    else
        echo -e "${RED}✗ Unhealthy (HTTP $response)${NC}"
        return 1
    fi
}

# Test all services
echo "=== Health Check Tests ==="
echo ""

# AI Service
test_health "AI Service" "http://localhost:8000/healthz"

# Creative Service
test_health "Creative Service" "http://localhost:8001/healthz"

# Translation Worker
test_health "Translation Worker" "http://localhost:8082/healthz"

# Email Worker
test_health "Email Worker" "http://localhost:8083/healthz"

# WhatsApp Worker
test_health "WhatsApp Worker" "http://localhost:8084/healthz"

echo ""
echo "=== Test Complete ==="

