#!/usr/bin/env python3
"""
Test script for WebSocket chat stream connection
"""
import asyncio
import json
import sys
from websockets.client import connect

async def test_websocket(conversation_id: str, token: str):
    """Test WebSocket connection to chat stream endpoint"""
    url = f"ws://localhost:8080/api/chats/conversations/{conversation_id}/stream?token={token}"
    
    print(f"Connecting to: {url}")
    print("Waiting for messages... (Press Ctrl+C to exit)\n")
    
    try:
        async with connect(url) as websocket:
            print("✅ WebSocket connected successfully!\n")
            print("Listening for messages...\n")
            
            # Keep connection alive and listen for messages
            while True:
                try:
                    message = await asyncio.wait_for(websocket.recv(), timeout=30.0)
                    data = json.loads(message)
                    print(f"📨 Received: {json.dumps(data, indent=2)}")
                except asyncio.TimeoutError:
                    print("⏱️  No message received in 30 seconds, still listening...")
                    # Send ping to keep connection alive
                    await websocket.ping()
                except json.JSONDecodeError:
                    print(f"📨 Received (non-JSON): {message}")
                    
    except Exception as e:
        print(f"❌ Error: {e}")
        return 1
    
    return 0

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python test_websocket.py <conversation_id> <token>")
        print("\nExample:")
        print("  python test_websocket.py fc9688d6-43b7-4e2e-b9b7-7a073733b375 'eyJhbGciOiJIUzI1NiIs...'")
        sys.exit(1)
    
    conversation_id = sys.argv[1]
    token = sys.argv[2]
    
    exit_code = asyncio.run(test_websocket(conversation_id, token))
    sys.exit(exit_code)

