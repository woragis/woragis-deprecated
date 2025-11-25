#!/bin/bash

# Script to request translations for entities that are missing translations
# Usage: ./request-translations.sh [email] [password]

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

echo "Login successful!"
echo ""

# Request translations for technical writings
echo "Requesting translations for technical writings..."
TECHNICAL_WRITINGS=$(curl -s -X GET "$API_BASE/technical-writings" \
    -H "Authorization: Bearer $ACCESS_TOKEN")

echo "$TECHNICAL_WRITINGS" | grep -o '"id":"[^"]*' | cut -d'"' -f4 | while read id; do
    if [ ! -z "$id" ]; then
        echo "  - Requesting translations for technical writing: $id"
        RESPONSE=$(curl -s -X POST "$API_BASE/translations/translate-entity" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $ACCESS_TOKEN" \
            -d "{\"entityType\":\"technical_writing\",\"entityId\":\"$id\"}")
        echo "$RESPONSE" | grep -q "queuedCount" && echo "    ✓ Queued" || echo "    ✗ Failed: $RESPONSE"
    fi
done

echo ""

# Request translations for impact metrics
echo "Requesting translations for impact metrics..."
IMPACT_METRICS=$(curl -s -X GET "$API_BASE/impact-metrics" \
    -H "Authorization: Bearer $ACCESS_TOKEN")

echo "$IMPACT_METRICS" | grep -o '"id":"[^"]*' | cut -d'"' -f4 | while read id; do
    if [ ! -z "$id" ]; then
        echo "  - Requesting translations for impact metric: $id"
        RESPONSE=$(curl -s -X POST "$API_BASE/translations/translate-entity" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $ACCESS_TOKEN" \
            -d "{\"entityType\":\"impact_metric\",\"entityId\":\"$id\"}")
        echo "$RESPONSE" | grep -q "queuedCount" && echo "    ✓ Queued" || echo "    ✗ Failed: $RESPONSE"
    fi
done

echo ""
echo "=========================================="
echo "Translation requests completed!"
echo "=========================================="
echo ""
echo "The translation worker will process these jobs asynchronously."
echo "Check translation status with: curl -X GET \"$API_BASE/translations?entityType=technical_writing\" -H \"Authorization: Bearer $ACCESS_TOKEN\""

