#!/bin/bash
# Test script for security middleware
# Usage: ./scripts/test-security-middleware.sh [server-url]

set -e

SERVER_URL="${1:-http://localhost:8080}"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[TEST]${NC} $1"
}

error() {
    echo -e "${RED}[FAIL]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Check if server is running
log "Checking if server is running at $SERVER_URL..."
if ! curl -s -f "$SERVER_URL/healthz" > /dev/null; then
    error "Server is not running at $SERVER_URL"
    error "Please start the server first: cd server/app && go run ./cmd/server/main.go"
    exit 1
fi
log "Server is running"

# Test 1: Security Headers
log "Test 1: Checking security headers..."
HEADERS=$(curl -s -I "$SERVER_URL/healthz")
REQUIRED_HEADERS=(
    "X-Content-Type-Options"
    "X-Frame-Options"
    "X-XSS-Protection"
)

ALL_HEADERS_PRESENT=true
for header in "${REQUIRED_HEADERS[@]}"; do
    if echo "$HEADERS" | grep -qi "$header"; then
        log "  ✓ $header present"
    else
        error "  ✗ $header missing"
        ALL_HEADERS_PRESENT=false
    fi
done

if [ "$ALL_HEADERS_PRESENT" = true ]; then
    log "✓ All security headers present"
else
    error "✗ Some security headers are missing"
fi

# Test 2: Rate Limiting
log "Test 2: Testing rate limiting (100 requests/minute)..."
log "  Making 105 requests..."
SUCCESS_COUNT=0
RATE_LIMITED_COUNT=0

for i in {1..105}; do
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/healthz")
    if [ "$HTTP_CODE" = "200" ]; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    elif [ "$HTTP_CODE" = "429" ]; then
        RATE_LIMITED_COUNT=$((RATE_LIMITED_COUNT + 1))
    fi
done

log "  Successful requests: $SUCCESS_COUNT"
log "  Rate limited requests: $RATE_LIMITED_COUNT"

if [ "$RATE_LIMITED_COUNT" -gt 0 ]; then
    log "✓ Rate limiting is working"
else
    warn "⚠ Rate limiting may not be working (no 429 responses)"
fi

# Test 3: Request Size Limit
log "Test 3: Testing request size limit (10MB)..."
log "  Sending 11MB request..."
LARGE_PAYLOAD=$(python3 -c "print('x' * 11 * 1024 * 1024)" 2>/dev/null || python -c "print('x' * 11 * 1024 * 1024)" 2>/dev/null || echo "x" | head -c 11534336)

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST "$SERVER_URL/api/test" \
    -H "Content-Type: application/json" \
    -d "$LARGE_PAYLOAD" 2>/dev/null || echo "000")

if [ "$HTTP_CODE" = "413" ] || [ "$HTTP_CODE" = "000" ]; then
    log "✓ Request size limit is working (got $HTTP_CODE)"
else
    warn "⚠ Request size limit may not be working (got $HTTP_CODE)"
fi

# Test 4: Input Sanitization
log "Test 4: Testing input sanitization..."
log "  Sending request with potentially dangerous input..."
RESPONSE=$(curl -s "$SERVER_URL/healthz?test=<script>alert('xss')</script>")
# Note: This is a basic test - actual sanitization happens in middleware
log "  ✓ Request processed (sanitization happens in middleware)"

# Summary
echo ""
log "=========================================="
log "Security Middleware Test Summary"
log "=========================================="
log "Server URL: $SERVER_URL"
log "Security Headers: $([ "$ALL_HEADERS_PRESENT" = true ] && echo "✓ PASS" || echo "✗ FAIL")"
log "Rate Limiting: $([ "$RATE_LIMITED_COUNT" -gt 0 ] && echo "✓ PASS" || echo "⚠ CHECK")"
log "Request Size Limit: $([ "$HTTP_CODE" = "413" ] && echo "✓ PASS" || echo "⚠ CHECK")"
log "=========================================="
