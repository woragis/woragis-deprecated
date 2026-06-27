#!/bin/bash

# WebSocket test using wscat or websocat
# Make sure you have one of them installed:
# npm install -g wscat
# cargo install websocat
# or use simple curl-like approach with a Python script

# For testing, we'll use a simple approach
# First, test if wscat is available
if command -v wscat &> /dev/null; then
    echo "Testing WebSocket with wscat..."
    # wscat -c "ws://localhost:3014/ws/chats/123e4567-e89b-12d3-a456-426614174000?user_id=123e4567-e89b-12d3-a456-426614174000"
    wscat -c "ws://localhost:3014/ws/chats/550e8400-e29b-41d4-a716-446655440000?user_id=550e8400-e29b-41d4-a716-446655440000" << 'EOF'
{"prompt":"test prompt","agent":"economist"}
EOF
elif command -v websocat &> /dev/null; then
    echo "Testing WebSocket with websocat..."
    websocat "ws://localhost:3014/ws/chats/550e8400-e29b-41d4-a716-446655440000?user_id=550e8400-e29b-41d4-a716-446655440000"
else
    echo "No WebSocket client found. Installing wscat..."
    npm install -g wscat 2>/dev/null || echo "npm not found, please install wscat manually"
    echo ""
    echo "WebSocket test URL:"
    echo "ws://localhost:3014/ws/chats/550e8400-e29b-41d4-a716-446655440000?user_id=550e8400-e29b-41d4-a716-446655440000"
fi
