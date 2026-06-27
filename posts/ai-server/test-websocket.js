#!/usr/bin/env node

const WebSocket = require('ws')
const { v4: uuidv4 } = require('uuid')

const chatId = uuidv4()
const userId = uuidv4()

console.log('\n=== WebSocket Test for Posts-AI Service ===')
console.log(`Service URL: ws://localhost:3014`)
console.log(`Chat ID: ${chatId}`)
console.log(`User ID: ${userId}`)
console.log(`\nConnecting to WebSocket...`)

const ws = new WebSocket(
  `ws://localhost:3014/ws/chats/${chatId}?user_id=${userId}`,
)

ws.on('open', function open() {
  console.log('✓ WebSocket connected successfully!')
  console.log('\nSending test message...')

  const testMessage = {
    prompt: 'How to write effective blog posts?',
    agent: 'economist',
  }

  console.log(`Message: ${JSON.stringify(testMessage, null, 2)}`)
  ws.send(JSON.stringify(testMessage))
})

ws.on('message', function message(data) {
  console.log(`\nReceived message from backend:`)
  try {
    const parsed = JSON.parse(data)
    console.log(JSON.stringify(parsed, null, 2))
  } catch (e) {
    console.log(data)
  }
})

ws.on('error', function error(err) {
  console.error('✗ WebSocket Error:', err.message)
})

ws.on('close', function close() {
  console.log('\n✓ WebSocket closed')
  process.exit(0)
})

// Auto-close after 10 seconds
setTimeout(() => {
  console.log('\n--- Test timeout (10s) ---')
  ws.close()
}, 10000)
