#!/bin/bash

# Script to seed skills using curl
# Usage: ./seed-skills.sh [email] [password]

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

# Create skills
echo "Creating skills..."

# Skill 1: Golang
curl -s -X POST "$API_BASE/skills" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "Golang",
        "description": "My primary language for building high-performance backend services. I develop main server applications with Go, leveraging its concurrency model and efficiency for scalable distributed systems.",
        "icon": "SiGo",
        "color": "cyan",
        "bgGradient": "from-cyan-900/30 to-cyan-800/20",
        "borderColor": "border-cyan-700/30",
        "hoverBorderColor": "hover:border-cyan-500/50",
        "shadowColor": "hover:shadow-cyan-500/20",
        "category": "language"
    }' > /dev/null

echo "Created skill: Golang"

# Skill 2: Python
curl -s -X POST "$API_BASE/skills" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "Python",
        "description": "Building AI servers and intelligent backend services with Python. I leverage modern ML frameworks and libraries like LangChain, OpenAI, and other AI tools to create intelligent applications that integrate seamlessly with my Golang services.",
        "icon": "SiPython",
        "color": "yellow",
        "bgGradient": "from-yellow-900/30 to-yellow-800/20",
        "borderColor": "border-yellow-700/30",
        "hoverBorderColor": "hover:border-yellow-500/50",
        "shadowColor": "hover:shadow-yellow-500/20",
        "category": "language"
    }' > /dev/null

echo "Created skill: Python"

# Skill 3: Docker
curl -s -X POST "$API_BASE/skills" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "Docker",
        "description": "Expertise in containerization for consistent deployments. I containerize applications for development, testing, and production environments, ensuring reproducibility and portability across different platforms.",
        "icon": "SiDocker",
        "color": "blue",
        "bgGradient": "from-blue-900/30 to-blue-800/20",
        "borderColor": "border-blue-700/30",
        "hoverBorderColor": "hover:border-blue-500/50",
        "shadowColor": "hover:shadow-blue-500/20",
        "category": "tool"
    }' > /dev/null

echo "Created skill: Docker"

# Skill 4: Kubernetes
curl -s -X POST "$API_BASE/skills" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "Kubernetes",
        "description": "Orchestrating containerized applications at scale. I design and manage K8s clusters for production-grade infrastructure and deployments, implementing auto-scaling, service discovery, and high availability patterns.",
        "icon": "SiKubernetes",
        "color": "indigo",
        "bgGradient": "from-indigo-900/30 to-indigo-800/20",
        "borderColor": "border-indigo-700/30",
        "hoverBorderColor": "hover:border-indigo-500/50",
        "shadowColor": "hover:shadow-indigo-500/20",
        "category": "infrastructure"
    }' > /dev/null

echo "Created skill: Kubernetes"

# Skill 5: DevOps
curl -s -X POST "$API_BASE/skills" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "DevOps",
        "description": "Bridging development and operations. I focus on automation, CI/CD pipelines, infrastructure as code, and cloud-native practices to streamline deployments and ensure reliable, scalable systems.",
        "icon": "Settings",
        "color": "purple",
        "bgGradient": "from-purple-900/30 to-purple-800/20",
        "borderColor": "border-purple-700/30",
        "hoverBorderColor": "hover:border-purple-500/50",
        "shadowColor": "hover:shadow-purple-500/20",
        "category": "devops"
    }' > /dev/null

echo "Created skill: DevOps"

# Additional skills based on the project

# Skill 6: PostgreSQL
curl -s -X POST "$API_BASE/skills" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "PostgreSQL",
        "description": "Expertise in PostgreSQL for robust data persistence. I design normalized schemas, optimize queries, implement JSONB for flexible data structures, and ensure data integrity with proper indexing and constraints.",
        "icon": "SiPostgresql",
        "color": "blue",
        "bgGradient": "from-blue-900/30 to-blue-800/20",
        "borderColor": "border-blue-700/30",
        "hoverBorderColor": "hover:border-blue-500/50",
        "shadowColor": "hover:shadow-blue-500/20",
        "category": "database"
    }' > /dev/null

echo "Created skill: PostgreSQL"

# Skill 7: Redis
curl -s -X POST "$API_BASE/skills" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "Redis",
        "description": "Using Redis for caching, session storage, and pub/sub messaging. I implement connection pooling, design efficient cache strategies, and leverage Redis streams for real-time event processing.",
        "icon": "SiRedis",
        "color": "red",
        "bgGradient": "from-red-900/30 to-red-800/20",
        "borderColor": "border-red-700/30",
        "hoverBorderColor": "hover:border-red-500/50",
        "shadowColor": "hover:shadow-red-500/20",
        "category": "database"
    }' > /dev/null

echo "Created skill: Redis"

# Skill 8: Fiber (Go Framework)
curl -s -X POST "$API_BASE/skills" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "Fiber",
        "description": "Building high-performance REST APIs with Fiber framework. I leverage middleware chains, route grouping, and Fiber'\''s Express-like API to create scalable backend services with excellent performance.",
        "icon": "SiGo",
        "color": "cyan",
        "bgGradient": "from-cyan-900/30 to-cyan-800/20",
        "borderColor": "border-cyan-700/30",
        "hoverBorderColor": "hover:border-cyan-500/50",
        "shadowColor": "hover:shadow-cyan-500/20",
        "category": "framework"
    }' > /dev/null

echo "Created skill: Fiber"

# Skill 9: LangChain
curl -s -X POST "$API_BASE/skills" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "LangChain",
        "description": "Building AI applications with LangChain for orchestration and RAG systems. I integrate LLMs, implement retrieval-augmented generation, and create intelligent workflows that combine language models with external data sources.",
        "icon": "SiOpenai",
        "color": "green",
        "bgGradient": "from-green-900/30 to-green-800/20",
        "borderColor": "border-green-700/30",
        "hoverBorderColor": "hover:border-green-500/50",
        "shadowColor": "hover:shadow-green-500/20",
        "category": "library"
    }' > /dev/null

echo "Created skill: LangChain"

# Skill 10: GORM
curl -s -X POST "$API_BASE/skills" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "GORM",
        "description": "Using GORM for database operations in Go. I leverage migrations, relationships, custom types for JSONB, and query optimization to build maintainable data access layers.",
        "icon": "SiGo",
        "color": "cyan",
        "bgGradient": "from-cyan-900/30 to-cyan-800/20",
        "borderColor": "border-cyan-700/30",
        "hoverBorderColor": "hover:border-cyan-500/50",
        "shadowColor": "hover:shadow-cyan-500/20",
        "category": "library"
    }' > /dev/null

echo "Created skill: GORM"

echo ""
echo "All skills created successfully!"
echo ""
echo "Summary:"
echo "  - 10 Skills created with full styling information"

