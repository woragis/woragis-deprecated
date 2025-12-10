#!/bin/bash

# Test script for getting chats of a job application
# Usage: ./test_job_application_chats.sh

EMAIL="masteringthecode.woragis@gmail.com"
PASSWORD="@Woragis2004"
JOB_APP_ID="30ef81a9-01fd-41d9-b2dc-4b55cdc2a4dc"
BASE_URL="http://localhost:8080"

echo "=========================================="
echo "Testing Job Application Chats Routes"
echo "=========================================="
echo ""

# Step 1: Login and get token
echo "Step 1: Logging in..."
LOGIN_RESPONSE=$(curl -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}" \
  -s)

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Login failed!"
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "✅ Login successful"
echo "Token: ${TOKEN:0:50}..."
echo ""

# Step 2: Test getting chats for job application (empty result)
echo "Step 2: Getting chats for job application ID: $JOB_APP_ID"
echo "Expected: Empty array (no conversations yet)"
echo "---"
RESPONSE=$(curl -X GET "$BASE_URL/api/chats/conversations/search?job_application_id=$JOB_APP_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -s)

echo "$RESPONSE" | python -m json.tool 2>/dev/null || echo "$RESPONSE"
echo ""

# Step 3: Create a conversation for the job application
echo "Step 3: Creating a conversation for the job application"
echo "---"
CREATE_RESPONSE=$(curl -X POST "$BASE_URL/api/chats/conversations" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"title\": \"Test Chat for Job Application\", \"description\": \"Testing chat route\", \"jobApplicationId\": \"$JOB_APP_ID\"}" \
  -s)

CONV_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$CONV_ID" ]; then
  echo "❌ Failed to create conversation"
  echo "Response: $CREATE_RESPONSE"
else
  echo "✅ Conversation created"
  echo "Conversation ID: $CONV_ID"
  echo "$CREATE_RESPONSE" | python -m json.tool 2>/dev/null || echo "$CREATE_RESPONSE"
fi
echo ""

# Step 4: Test getting chats for job application again (should have the conversation)
echo "Step 4: Getting chats for job application again (should now include the conversation)"
echo "---"
RESPONSE=$(curl -X GET "$BASE_URL/api/chats/conversations/search?job_application_id=$JOB_APP_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -s)

echo "$RESPONSE" | python -m json.tool 2>/dev/null || echo "$RESPONSE"
echo ""

# Step 5: Test with query parameter
echo "Step 5: Testing with query parameter"
echo "---"
RESPONSE=$(curl -X GET "$BASE_URL/api/chats/conversations/search?job_application_id=$JOB_APP_ID&q=Test" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -s)

echo "$RESPONSE" | python -m json.tool 2>/dev/null || echo "$RESPONSE"
echo ""

# Step 6: Test with invalid job_application_id
echo "Step 6: Testing with invalid job_application_id (should return error)"
echo "---"
RESPONSE=$(curl -X GET "$BASE_URL/api/chats/conversations/search?job_application_id=invalid-uuid" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -s)

echo "$RESPONSE" | python -m json.tool 2>/dev/null || echo "$RESPONSE"
echo ""

# Step 7: Test without job_application_id (should return all conversations)
echo "Step 7: Testing without job_application_id filter (should return all conversations)"
echo "---"
RESPONSE=$(curl -X GET "$BASE_URL/api/chats/conversations/search" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -s)

echo "$RESPONSE" | python -m json.tool 2>/dev/null | head -50 || echo "$RESPONSE" | head -50
echo ""

echo "=========================================="
echo "Tests completed!"
echo "=========================================="

