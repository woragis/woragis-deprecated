#!/bin/bash

# Script to run Python service/worker tests in Docker
# Usage: ./scripts/run-python-tests-docker.sh [service-name]
# If service-name is provided, only that service's tests will run
# Otherwise, all Python services' tests will run

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}Running Python Tests in Docker${NC}"
echo "=========================================="

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running${NC}"
    exit 1
fi

# Function to run tests for a Python service
run_service_tests() {
    local service_name=$1
    local service_dir="$BACKEND_DIR/$service_name"
    
    if [ ! -d "$service_dir" ]; then
        echo -e "${RED}Error: Service directory not found: $service_dir${NC}"
        return 1
    fi
    
    if [ ! -f "$service_dir/Dockerfile.test" ]; then
        echo -e "${YELLOW}No Dockerfile.test found for $service_name, skipping${NC}"
        return 1
    fi
    
    echo -e "\n${BLUE}Testing $service_name...${NC}"
    echo "----------------------------------------"
    
    cd "$service_dir"
    
    # Build test image
    echo -e "${YELLOW}Building test image...${NC}"
    docker build -f Dockerfile.test -t "${service_name}-test:latest" . || {
        echo -e "${RED}Failed to build test image for $service_name${NC}"
        return 1
    }
    
    # Run tests
    echo -e "${YELLOW}Running tests...${NC}"
    if docker run --rm "${service_name}-test:latest"; then
        echo -e "${GREEN}✓ $service_name tests passed${NC}"
        return 0
    else
        echo -e "${RED}✗ $service_name tests failed${NC}"
        return 1
    fi
}

# Run tests based on argument
if [ $# -eq 0 ]; then
    # Run all Python service tests
    run_service_tests "ai-service"
    run_service_tests "creative-service"
    run_service_tests "docs-service"
    run_service_tests "resume-worker"
else
    # Run specific service tests
    run_service_tests "$1"
fi

echo -e "\n${GREEN}Python Docker tests completed!${NC}"

