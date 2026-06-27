#!/usr/bin/env python3
"""
WebSocket Test for Posts-AI Service
Tests the /ws/chats/:id endpoint
"""

import asyncio
import websockets
import json
import uuid
from datetime import datetime


async def test_websocket():
    chat_id = str(uuid.uuid4())
    user_id = str(uuid.uuid4())

    uri = f"ws://localhost:3014/ws/chats/{chat_id}?user_id={user_id}"

    print("\n=== WebSocket Test for Posts-AI Service ===")
    print(f"Service URL: ws://localhost:3014")
    print(f"Chat ID: {chat_id}")
    print(f"User ID: {user_id}")
    print(f"\nConnecting to WebSocket...\n")

    try:
        async with websockets.connect(uri) as websocket:
            print("✓ WebSocket connected successfully!")
            print(f"  Server: {websocket.remote_address}")
            print(
                f"  Connection time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")

            # Send test message
            test_message = {
                "prompt": "How to write effective blog posts?",
                "agent": "economist"
            }

            print(f"\nSending test message:")
            print(f"  {json.dumps(test_message, indent=2)}")

            await websocket.send(json.dumps(test_message))

            # Wait for response with timeout
            print(f"\nWaiting for response...")
            try:
                response = await asyncio.wait_for(websocket.recv(), timeout=5.0)
                print(f"\n✓ Received message from backend:")
                try:
                    parsed = json.loads(response)
                    print(f"  {json.dumps(parsed, indent=2)}")
                except json.JSONDecodeError:
                    print(f"  {response}")

            except asyncio.TimeoutError:
                print(f"  (No response received within 5 seconds)")
                print(f"  This is expected if AI service is not configured")

            print(f"\n✓ WebSocket test completed")

    except ConnectionRefusedError:
        print(f"✗ Connection refused - service not running on localhost:3014")
    except Exception as e:
        print(f"✗ Error: {type(e).__name__}: {e}")


if __name__ == "__main__":
    asyncio.run(test_websocket())
