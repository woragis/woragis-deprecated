#!/bin/bash

# Test script for chat websocket and AI service

echo "=== Testing AI Service Stream ==="
curl -X POST http://localhost:8000/v1/chat/stream \
  -H "Content-Type: application/json" \
  -d '{"agent": "startup", "input": "Hello! What is a startup?", "provider": "openai"}' \
  --no-buffer

echo ""
echo ""
echo "=== Logging in to backend ==="
LOGIN_RESPONSE=$(curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "masteringthecode.woragis@gmail.com", "password": "@Woragis2004"}' \
  -s)

echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4 > /tmp/token.txt

TOKEN=$(cat /tmp/token.txt)

if [ -z "$TOKEN" ]; then
  echo "Login failed or token not found"
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "Token obtained: ${TOKEN:0:20}..."

echo ""
echo "=== Creating a conversation ==="
CONV_RESPONSE=$(curl -X POST http://localhost:8080/api/chats/conversations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Test Chat", "description": "Testing websocket"}' \
  -s)

CONV_ID=$(echo "$CONV_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$CONV_ID" ]; then
  echo "Failed to create conversation"
  echo "Response: $CONV_RESPONSE"
  exit 1
fi

echo "Conversation ID: $CONV_ID"

echo ""
echo "=== Testing WebSocket connection (requires wscat) ==="
echo "To test websocket, run:"
echo "wscat -c 'ws://localhost:8080/api/chats/conversations/$CONV_ID/stream?token=$TOKEN'"
echo ""
echo "Or use the test HTML page in frontend/src/routes/test-chat.html"

