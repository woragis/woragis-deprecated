#!/bin/bash

# Script to seed interests using curl
# Usage: ./seed-interests.sh [email] [password]

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

echo "Login successful! Token: ${ACCESS_TOKEN:0:20}..."

# Create interests
echo "Creating interests..."

# Interest 1: AI & RAG
curl -s -X POST "$API_BASE/interests" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "title": "AI & RAG",
        "description": "Fascinated by Artificial Intelligence and Retrieval-Augmented Generation (RAG) systems. I'\''m exploring how to build intelligent applications that can retrieve and synthesize information effectively. I work with AI servers built in Python, leveraging modern ML frameworks and libraries to create intelligent backend services.",
        "icon": "Brain",
        "color": "pink-purple",
        "bgGradient": "from-pink-900/30 to-purple-900/20",
        "borderColor": "border-pink-700/30",
        "hoverBorderColor": "hover:border-pink-500/50",
        "shadowColor": "hover:shadow-pink-500/20",
        "fullWidth": false,
        "featured": true
    }' > /dev/null

echo "Created interest: AI & RAG"

# Interest 2: Redis & Pub/Sub
curl -s -X POST "$API_BASE/interests" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "title": "Redis & Pub/Sub",
        "description": "Deep interest in Redis Pub/Sub patterns for real-time communication between distributed services. I design systems where multiple servers communicate seamlessly through Redis messaging. I implement inter-service communication architectures using Redis as the backbone, enabling scalable and responsive distributed applications.",
        "icon": "SiRedis",
        "color": "red-orange",
        "bgGradient": "from-red-900/30 to-orange-900/20",
        "borderColor": "border-red-700/30",
        "hoverBorderColor": "hover:border-red-500/50",
        "shadowColor": "hover:shadow-red-500/20",
        "fullWidth": false,
        "featured": true
    }' > /dev/null

echo "Created interest: Redis & Pub/Sub"

# Interest 3: Distributed Architecture
curl -s -X POST "$API_BASE/interests" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "title": "Distributed Architecture",
        "description": "I specialize in building hybrid architectures where the main server is built with Golang for performance and reliability, while AI services are implemented in Python to leverage the rich ML ecosystem. These services communicate through Redis Pub/Sub, enabling a microservices architecture that'\''s both scalable and maintainable. The combination allows each service to use the best tools for its specific domain.",
        "icon": "GitBranch",
        "color": "green-emerald",
        "bgGradient": "from-green-900/30 to-emerald-900/20",
        "borderColor": "border-green-700/30",
        "hoverBorderColor": "hover:border-green-500/50",
        "shadowColor": "hover:shadow-green-500/20",
        "fullWidth": true,
        "featured": true
    }' > /dev/null

echo "Created interest: Distributed Architecture"

echo ""
echo "All interests created successfully!"
echo ""
echo "Summary:"
echo "  - 3 Interests created with full styling information"
echo "  - All interests are featured (publicly visible)"

