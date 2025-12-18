#!/bin/bash

# Script to run performance tests for all workers
# Usage: ./run-worker-performance-tests.sh [worker-name]
# If worker-name is provided, only that worker's tests will run
# Otherwise, all workers' performance tests will run

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
SERVER_DIR="$BACKEND_DIR/server"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting Worker Performance Tests${NC}"
echo "=========================================="

# Check if Docker Compose is available
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: docker-compose not found${NC}"
    exit 1
fi

# Start test dependencies
echo -e "\n${YELLOW}Starting test dependencies...${NC}"
cd "$SERVER_DIR"
docker-compose -f docker-compose.test.yml up -d

# Wait for services to be ready
echo -e "${YELLOW}Waiting for services to be ready...${NC}"
sleep 5

# Function to run performance tests for a worker
run_worker_tests() {
    local worker_name=$1
    local worker_dir="$BACKEND_DIR/$worker_name"
    
    if [ ! -d "$worker_dir" ]; then
        echo -e "${RED}Error: Worker directory not found: $worker_dir${NC}"
        return 1
    fi
    
    echo -e "\n${GREEN}Running performance tests for $worker_name...${NC}"
    echo "----------------------------------------"
    
    cd "$worker_dir"
    
    if [ -d "internal/integration" ]; then
        # Run performance tests (skip short tests)
        go test -tags=integration -run "Test.*Performance|Test.*Load|Test.*Latency|Test.*Rate|Benchmark" ./internal/integration/... -v
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ $worker_name performance tests passed${NC}"
        else
            echo -e "${RED}✗ $worker_name performance tests failed${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}No integration tests found for $worker_name${NC}"
    fi
}

# Run tests based on argument
if [ $# -eq 0 ]; then
    # Run all worker performance tests
    run_worker_tests "email-worker"
    run_worker_tests "translation-worker"
    run_worker_tests "whatsapp-worker"
else
    # Run specific worker tests
    run_worker_tests "$1"
fi

# Cleanup
echo -e "\n${YELLOW}Cleaning up...${NC}"
cd "$SERVER_DIR"
docker-compose -f docker-compose.test.yml down

echo -e "\n${GREEN}Performance tests completed!${NC}"
