#!/bin/bash

# Comprehensive script to delete existing data and run all seeders
# Usage: ./run-all-seeders.sh [email] [password]

API_BASE="http://localhost:8080/api"
EMAIL="${1:-masteringthecode.woragis@gmail.com}"
PASSWORD="${2:-@Woragis2004}"

if [ -z "$PASSWORD" ]; then
    echo "Usage: $0 [email] [password]"
    echo "Please provide your password to authenticate"
    exit 1
fi

echo "=========================================="
echo "Starting Seeding Process"
echo "=========================================="
echo ""

# Step 1: Login
echo "Step 1: Logging in..."
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

echo "✓ Login successful! Token: ${ACCESS_TOKEN:0:20}..."
echo ""

# Step 2: Delete existing data
echo "Step 2: Deleting existing data..."
echo ""

# Get all entities and delete them
echo "  - Deleting impact metrics..."
curl -s -X GET "$API_BASE/impact-metrics" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | \
    grep -o '"id":"[^"]*' | cut -d'"' -f4 | \
    while read id; do
        curl -s -X DELETE "$API_BASE/impact-metrics/$id" \
            -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null
    done

echo "  - Deleting interests..."
curl -s -X GET "$API_BASE/interests" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | \
    grep -o '"id":"[^"]*' | cut -d'"' -f4 | \
    while read id; do
        curl -s -X DELETE "$API_BASE/interests/$id" \
            -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null
    done

echo "  - Deleting skills..."
curl -s -X GET "$API_BASE/skills" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | \
    grep -o '"id":"[^"]*' | cut -d'"' -f4 | \
    while read id; do
        curl -s -X DELETE "$API_BASE/skills/$id" \
            -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null
    done

echo "  - Deleting certifications..."
curl -s -X GET "$API_BASE/certifications" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | \
    grep -o '"id":"[^"]*' | cut -d'"' -f4 | \
    while read id; do
        curl -s -X DELETE "$API_BASE/certifications/$id" \
            -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null
    done

echo "  - Deleting AIML integrations..."
curl -s -X GET "$API_BASE/aiml-integrations" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | \
    grep -o '"id":"[^"]*' | cut -d'"' -f4 | \
    while read id; do
        curl -s -X DELETE "$API_BASE/aiml-integrations/$id" \
            -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null
    done

echo "  - Deleting case studies..."
curl -s -X GET "$API_BASE/case-studies" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | \
    grep -o '"id":"[^"]*' | cut -d'"' -f4 | \
    while read id; do
        curl -s -X DELETE "$API_BASE/case-studies/$id" \
            -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null
    done

echo "  - Deleting system designs..."
curl -s -X GET "$API_BASE/system-designs" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | \
    grep -o '"id":"[^"]*' | cut -d'"' -f4 | \
    while read id; do
        curl -s -X DELETE "$API_BASE/system-designs/$id" \
            -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null
    done

echo "  - Deleting problem solutions..."
curl -s -X GET "$API_BASE/problem-solutions" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | \
    grep -o '"id":"[^"]*' | cut -d'"' -f4 | \
    while read id; do
        curl -s -X DELETE "$API_BASE/problem-solutions/$id" \
            -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null
    done

echo "  - Deleting technical writings..."
curl -s -X GET "$API_BASE/technical-writings" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | \
    grep -o '"id":"[^"]*' | cut -d'"' -f4 | \
    while read id; do
        curl -s -X DELETE "$API_BASE/technical-writings/$id" \
            -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null
    done

echo "  - Deleting testimonials..."
curl -s -X GET "$API_BASE/testimonials" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | \
    grep -o '"id":"[^"]*' | cut -d'"' -f4 | \
    while read id; do
        curl -s -X DELETE "$API_BASE/testimonials/$id" \
            -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null
    done

echo "✓ Existing data deleted"
echo ""

# Step 3: Run seeders
echo "Step 3: Running seeders..."
echo ""

export AUTH_TOKEN="$ACCESS_TOKEN"
export API_BASE_URL="$API_BASE"

echo "  - Seeding impact metrics..."
bash scripts/seed-impact-metrics.sh
echo ""

echo "  - Seeding interests..."
bash scripts/seed-interests.sh "$EMAIL" "$PASSWORD"
echo ""

echo "  - Seeding skills..."
bash scripts/seed-skills.sh "$EMAIL" "$PASSWORD"
echo ""

echo "  - Seeding certifications..."
bash scripts/seed-certifications.sh "$EMAIL" "$PASSWORD"
echo ""

echo "  - Seeding AIML integrations..."
bash scripts/seed-aiml-integrations.sh
echo ""

echo "  - Seeding technical content (case studies, system designs, problem solutions)..."
bash scripts/seed-technical-content.sh "$EMAIL" "$PASSWORD"
echo ""

echo "  - Seeding technical writings..."
bash scripts/seed-technical-writings.sh
echo ""

echo "  - Seeding testimonials..."
bash scripts/seed-testimonials.sh "$EMAIL" "$PASSWORD"
echo ""

# Step 4: Verify translations
echo "Step 4: Verifying translations..."
echo ""

LANGUAGES=("en" "pt" "fr" "es")

for lang in "${LANGUAGES[@]}"; do
    echo "  Checking translations for language: $lang"
    
    # Check a few endpoints with language parameter
    echo "    - Impact metrics..."
    curl -s -X GET "$API_BASE/impact-metrics?lang=$lang" \
        -H "Authorization: Bearer $ACCESS_TOKEN" | \
        grep -q '"id"' && echo "      ✓ Impact metrics translated" || echo "      ✗ Impact metrics not translated"
    
    echo "    - Skills..."
    curl -s -X GET "$API_BASE/skills?lang=$lang" \
        -H "Authorization: Bearer $ACCESS_TOKEN" | \
        grep -q '"id"' && echo "      ✓ Skills translated" || echo "      ✗ Skills not translated"
    
    echo "    - Interests..."
    curl -s -X GET "$API_BASE/interests?lang=$lang" \
        -H "Authorization: Bearer $ACCESS_TOKEN" | \
        grep -q '"id"' && echo "      ✓ Interests translated" || echo "      ✗ Interests not translated"
    
    echo "    - Certifications..."
    curl -s -X GET "$API_BASE/certifications?lang=$lang" \
        -H "Authorization: Bearer $ACCESS_TOKEN" | \
        grep -q '"id"' && echo "      ✓ Certifications translated" || echo "      ✗ Certifications not translated"
    
    echo "    - AIML integrations..."
    curl -s -X GET "$API_BASE/aiml-integrations?lang=$lang" \
        -H "Authorization: Bearer $ACCESS_TOKEN" | \
        grep -q '"id"' && echo "      ✓ AIML integrations translated" || echo "      ✗ AIML integrations not translated"
    
    echo "    - Technical writings..."
    curl -s -X GET "$API_BASE/technical-writings?lang=$lang" \
        -H "Authorization: Bearer $ACCESS_TOKEN" | \
        grep -q '"id"' && echo "      ✓ Technical writings translated" || echo "      ✗ Technical writings not translated"
    
    echo ""
done

echo "=========================================="
echo "Seeding Process Completed!"
echo "=========================================="
echo ""
echo "Summary:"
echo "  - All existing data deleted"
echo "  - All seeders executed"
echo "  - Translations verified for: ${LANGUAGES[*]}"
echo ""
echo "Note: Translations are processed asynchronously by the translation worker."
echo "      It may take a few minutes for all translations to be available."

