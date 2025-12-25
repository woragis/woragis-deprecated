#!/bin/bash

# Master script to run all tests in the recommended order
# Usage: ./scripts/run-all-tests.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Track results
PASSED=0
FAILED=0
SKIPPED=0

# Function to print section header
print_section() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

# Function to print result
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
        ((PASSED++))
    else
        echo -e "${RED}✗ $2${NC}"
        ((FAILED++))
    fi
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo -e "${RED}Error: Docker is not running${NC}"
        exit 1
    fi
}

# Start test dependencies
start_test_deps() {
    print_section "Starting Test Dependencies"
    cd "$BACKEND_DIR"
    
    if [ -f "docker-compose.test.yml" ]; then
        echo "Starting docker-compose.test.yml..."
        docker-compose -f docker-compose.test.yml up -d
        echo "Waiting for services to be healthy..."
        sleep 10
        
        # Check service health
        docker-compose -f docker-compose.test.yml ps
        print_result $? "Test dependencies started"
    else
        echo -e "${YELLOW}docker-compose.test.yml not found, using server's docker-compose.test.yml${NC}"
        cd "$BACKEND_DIR/server"
        if [ -f "docker-compose.test.yml" ]; then
            docker-compose -f docker-compose.test.yml up -d
            sleep 10
            docker-compose -f docker-compose.test.yml ps
            print_result $? "Test dependencies started"
        else
            echo -e "${RED}No docker-compose.test.yml found${NC}"
            exit 1
        fi
    fi
}

# Stop test dependencies
stop_test_deps() {
    print_section "Stopping Test Dependencies"
    cd "$BACKEND_DIR"
    
    if [ -f "docker-compose.test.yml" ]; then
        docker-compose -f docker-compose.test.yml down
    else
        cd "$BACKEND_DIR/server"
        if [ -f "docker-compose.test.yml" ]; then
            docker-compose -f docker-compose.test.yml down
        fi
    fi
    echo -e "${GREEN}Test dependencies stopped${NC}"
}

# Run integration tests
run_integration_tests() {
    print_section "1. Running Integration Tests"
    
    # Server integration tests
    echo -e "\n${YELLOW}Server Integration Tests...${NC}"
    cd "$BACKEND_DIR/server"
    if [ -f "scripts/run-integration-tests.sh" ]; then
        bash scripts/run-integration-tests.sh || print_result 1 "Server integration tests"
    else
        go test ./app/internal/integration/... -tags=integration -v || print_result 1 "Server integration tests"
    fi
    print_result $? "Server integration tests"
    
    # Python services integration tests
    for service in ai-service creative-service docs-service; do
        echo -e "\n${YELLOW}${service} Integration Tests...${NC}"
        cd "$BACKEND_DIR/$service"
        if [ -d "tests/integration" ]; then
            make test-integration 2>/dev/null || pytest tests/integration/ -v -m integration || print_result 1 "$service integration tests"
            print_result $? "$service integration tests"
        else
            echo -e "${YELLOW}No integration tests found for $service${NC}"
            ((SKIPPED++))
        fi
    done
    
    # Go workers integration tests
    for worker in email-worker translation-worker whatsapp-worker; do
        echo -e "\n${YELLOW}${worker} Integration Tests...${NC}"
        cd "$BACKEND_DIR/$worker"
        if [ -d "internal/integration" ]; then
            go test ./internal/integration/... -tags=integration -v || print_result 1 "$worker integration tests"
            print_result $? "$worker integration tests"
        else
            echo -e "${YELLOW}No integration tests found for $worker${NC}"
            ((SKIPPED++))
        fi
    done
    
    # Python worker integration tests
    echo -e "\n${YELLOW}resume-worker Integration Tests...${NC}"
    cd "$BACKEND_DIR/resume-worker"
    if [ -d "tests/integration" ]; then
        make test-integration 2>/dev/null || pytest tests/integration/ -v || print_result 1 "resume-worker integration tests"
        print_result $? "resume-worker integration tests"
    else
        echo -e "${YELLOW}No integration tests found for resume-worker${NC}"
        ((SKIPPED++))
    fi
    
    # Node.js worker integration tests
    echo -e "\n${YELLOW}job-application-worker Integration Tests...${NC}"
    cd "$BACKEND_DIR/job-application-worker"
    if [ -d "tests/integration" ]; then
        npm run test:integration 2>/dev/null || npm test -- tests/integration/ || print_result 1 "job-application-worker integration tests"
        print_result $? "job-application-worker integration tests"
    else
        echo -e "${YELLOW}No integration tests found for job-application-worker${NC}"
        ((SKIPPED++))
    fi
}

# Run performance tests
run_performance_tests() {
    print_section "2. Running Performance Tests"
    
    # Server performance tests
    echo -e "\n${YELLOW}Server Performance Tests...${NC}"
    cd "$BACKEND_DIR/server"
    go test ./app/internal/integration/... -tags="integration,performance_test" -v -run "Test.*Performance|Test.*Load|Test.*Latency|Benchmark" || print_result 1 "Server performance tests"
    print_result $? "Server performance tests"
    
    # Python services performance tests
    for service in ai-service creative-service docs-service; do
        echo -e "\n${YELLOW}${service} Performance Tests...${NC}"
        cd "$BACKEND_DIR/$service"
        if [ -f "app/tests/performance_test.py" ] || [ -f "tests/performance_test.py" ]; then
            pytest app/tests/performance_test.py tests/performance_test.py -v 2>/dev/null || print_result 1 "$service performance tests"
            print_result $? "$service performance tests"
        else
            echo -e "${YELLOW}No performance tests found for $service${NC}"
            ((SKIPPED++))
        fi
    done
    
    # Go workers performance tests
    for worker in email-worker translation-worker whatsapp-worker; do
        echo -e "\n${YELLOW}${worker} Performance Tests...${NC}"
        cd "$BACKEND_DIR/$worker"
        if [ -d "internal/integration" ]; then
            go test ./internal/integration/... -tags=integration -v -run "Test.*Performance|Test.*Load|Test.*Latency|Test.*Rate|Benchmark" || print_result 1 "$worker performance tests"
            print_result $? "$worker performance tests"
        else
            echo -e "${YELLOW}No performance tests found for $worker${NC}"
            ((SKIPPED++))
        fi
    done
    
    # Python worker performance tests
    echo -e "\n${YELLOW}resume-worker Performance Tests...${NC}"
    cd "$BACKEND_DIR/resume-worker"
    if [ -f "tests/performance_test.py" ]; then
        pytest tests/performance_test.py -v || print_result 1 "resume-worker performance tests"
        print_result $? "resume-worker performance tests"
    else
        echo -e "${YELLOW}No performance tests found for resume-worker${NC}"
        ((SKIPPED++))
    fi
    
    # Node.js worker performance tests
    echo -e "\n${YELLOW}job-application-worker Performance Tests...${NC}"
    cd "$BACKEND_DIR/job-application-worker"
    if [ -f "tests/performance.test.js" ]; then
        npm run test:performance 2>/dev/null || npm test -- tests/performance.test.js || print_result 1 "job-application-worker performance tests"
        print_result $? "job-application-worker performance tests"
    else
        echo -e "${YELLOW}No performance tests found for job-application-worker${NC}"
        ((SKIPPED++))
    fi
}

# Test builds
test_builds() {
    print_section "3. Testing Builds (Compilation)"
    
    # Server build
    echo -e "\n${YELLOW}Building server...${NC}"
    cd "$BACKEND_DIR/server"
    go build ./app/cmd/server/... || print_result 1 "Server build"
    print_result $? "Server build"
    
    # Go workers build
    for worker in email-worker translation-worker whatsapp-worker; do
        echo -e "\n${YELLOW}Building $worker...${NC}"
        cd "$BACKEND_DIR/$worker"
        go build ./cmd/... 2>/dev/null || go build ./... || print_result 1 "$worker build"
        print_result $? "$worker build"
    done
    
    # Python services - check syntax
    for service in ai-service creative-service docs-service resume-worker; do
        echo -e "\n${YELLOW}Checking $service syntax...${NC}"
        cd "$BACKEND_DIR/$service"
        python -m py_compile app/main.py 2>/dev/null || python -m py_compile src/main.py 2>/dev/null || echo -e "${YELLOW}No main.py found, skipping${NC}"
        print_result $? "$service syntax check"
    done
    
    # Node.js worker - check syntax
    echo -e "\n${YELLOW}Checking job-application-worker syntax...${NC}"
    cd "$BACKEND_DIR/job-application-worker"
    node --check src/worker.js 2>/dev/null || print_result 1 "job-application-worker syntax check"
    print_result $? "job-application-worker syntax check"
}

# Test security middleware
test_security_middleware() {
    print_section "4. Testing Security Middleware"
    
    cd "$BACKEND_DIR"
    if [ -f "scripts/test-security-middleware.sh" ]; then
        bash scripts/test-security-middleware.sh || print_result 1 "Security middleware tests"
        print_result $? "Security middleware tests"
    else
        echo -e "${YELLOW}test-security-middleware.sh not found, skipping${NC}"
        ((SKIPPED++))
    fi
}

# Test backup scripts
test_backup_scripts() {
    print_section "5. Testing Backup Scripts"
    
    cd "$BACKEND_DIR"
    if [ -f "scripts/test-backups.sh" ]; then
        bash scripts/test-backups.sh || print_result 1 "Backup scripts tests"
        print_result $? "Backup scripts tests"
    else
        echo -e "${YELLOW}test-backups.sh not found, skipping${NC}"
        ((SKIPPED++))
    fi
}

# Print summary
print_summary() {
    print_section "Test Summary"
    echo -e "${GREEN}Passed: $PASSED${NC}"
    echo -e "${RED}Failed: $FAILED${NC}"
    echo -e "${YELLOW}Skipped: $SKIPPED${NC}"
    
    if [ $FAILED -eq 0 ]; then
        echo -e "\n${GREEN}All tests passed! ✓${NC}"
        return 0
    else
        echo -e "\n${RED}Some tests failed. Please review the output above.${NC}"
        return 1
    fi
}

# Main execution
main() {
    echo -e "${GREEN}Starting Comprehensive Test Suite${NC}"
    echo "=========================================="
    
    check_docker
    start_test_deps
    
    # Run tests
    run_integration_tests
    run_performance_tests
    test_builds
    test_security_middleware
    test_backup_scripts
    
    # Cleanup
    stop_test_deps
    
    # Print summary
    print_summary
    exit $?
}

# Run main
main

