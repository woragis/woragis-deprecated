#!/bin/bash

# Script to create testimonials and posts using curl
# Usage: ./create-content.sh [email] [password]

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

# Extract access_token from JSON response (handle both accessToken and access_token)
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
if [ -z "$ACCESS_TOKEN" ]; then
    ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"accessToken":"[^"]*' | cut -d'"' -f4)
fi

if [ -z "$ACCESS_TOKEN" ]; then
    echo "Failed to login. Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "Login successful! Token: ${ACCESS_TOKEN:0:20}..."

# Create testimonials
echo "Creating testimonials..."

TESTIMONIALS=(
    '{"authorName":"Sarah Chen","authorRole":"Senior Software Engineer","authorCompany":"TechCorp","content":"Working with Jezreel was an absolute pleasure. His deep understanding of distributed systems and ability to architect scalable solutions is impressive. He consistently delivered high-quality code and was always willing to share knowledge with the team. I'\''d work with him again in a heartbeat.","rating":5,"linkedinUrl":"https://linkedin.com/in/sarahchen","displayOrder":1}'
    '{"authorName":"Michael Rodriguez","authorRole":"CTO","authorCompany":"StartupXYZ","content":"Jezreel is one of the most technically skilled developers I'\''ve had the privilege to work with. His expertise in Go, Redis, and microservices architecture helped us build a robust system that handles millions of requests daily. His attention to detail and problem-solving skills are exceptional.","rating":5,"linkedinUrl":"https://linkedin.com/in/michaelrodriguez","displayOrder":2}'
    '{"authorName":"Emily Watson","authorRole":"Product Manager","authorCompany":"InnovateLabs","content":"Jezreel'\''s ability to translate complex technical requirements into elegant solutions is remarkable. He'\''s not just a great coder—he'\''s a great communicator who can explain technical concepts to non-technical stakeholders. His work on our RAG system was groundbreaking.","rating":5,"linkedinUrl":"https://linkedin.com/in/emilywatson","displayOrder":3}'
    '{"authorName":"David Kim","authorRole":"Lead DevOps Engineer","authorCompany":"CloudScale Inc","content":"I'\''ve collaborated with Jezreel on several infrastructure projects, and his knowledge of Kubernetes, Docker, and cloud architecture is top-notch. He'\''s proactive, reliable, and always thinking about scalability and reliability. A true professional.","rating":5,"linkedinUrl":"https://linkedin.com/in/davidkim","displayOrder":4}'
    '{"authorName":"Alexandra Thompson","authorRole":"Engineering Manager","authorCompany":"DataFlow Systems","content":"Jezreel brings a unique combination of technical depth and practical problem-solving. His work on our pub/sub messaging system using Redis was instrumental in improving our system'\''s performance. He'\''s a team player who elevates everyone around him.","rating":5,"linkedinUrl":"https://linkedin.com/in/alexandrathompson","displayOrder":5}'
)

for testimonial in "${TESTIMONIALS[@]}"; do
    RESPONSE=$(curl -s -X POST "$API_BASE/testimonials" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d "$testimonial")
    
    TESTIMONIAL_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    if [ ! -z "$TESTIMONIAL_ID" ]; then
        echo "Created testimonial: $TESTIMONIAL_ID"
        # Approve the testimonial
        curl -s -X POST "$API_BASE/testimonials/$TESTIMONIAL_ID/approve" \
            -H "Authorization: Bearer $ACCESS_TOKEN" > /dev/null
        echo "Approved testimonial: $TESTIMONIAL_ID"
    else
        echo "Failed to create testimonial. Response: $RESPONSE"
    fi
done

echo "Testimonials created!"

# Create blog posts
echo "Creating blog posts..."

# Post 1: Redis Pub/Sub
curl -s -X POST "$API_BASE/posts" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "title": "Building Scalable Microservices with Go and Redis Pub/Sub",
        "content": "# Building Scalable Microservices with Go and Redis Pub/Sub\n\nIn modern distributed systems, efficient communication between services is crucial. Redis Pub/Sub provides a lightweight, fast messaging solution that'\''s perfect for microservices architectures.\n\n## Why Redis Pub/Sub?\n\nRedis Pub/Sub offers several advantages:\n- **Low Latency**: Sub-millisecond message delivery\n- **Decoupling**: Services don'\''t need to know about each other\n- **Scalability**: Handle millions of messages per second\n- **Simplicity**: Easy to implement and maintain\n\n## Implementation Pattern\n\nHere'\''s a common pattern I use in production using Go and Redis.\n\n## Best Practices\n\n1. **Error Handling**: Always implement retry logic for failed messages\n2. **Message Format**: Use structured formats like JSON or Protocol Buffers\n3. **Monitoring**: Track message rates and latency\n4. **Backpressure**: Implement rate limiting to prevent overload\n\n## Real-World Example\n\nIn a recent project, I used Redis Pub/Sub to handle real-time notifications across 50+ microservices. The system processes over 100,000 messages per second with sub-10ms latency.\n\nThe key was implementing proper connection pooling, message batching, and circuit breakers to handle failures gracefully.",
        "excerpt": "Learn how to build scalable microservices using Go and Redis Pub/Sub for efficient inter-service communication.",
        "status": "published",
        "featured": true,
        "metaTitle": "Building Scalable Microservices with Go and Redis Pub/Sub",
        "metaDescription": "A comprehensive guide to implementing microservices communication using Go and Redis Pub/Sub patterns.",
        "tagNames": ["go", "redis", "microservices", "distributed-systems"]
    }' > /dev/null

echo "Created post: Building Scalable Microservices with Go and Redis Pub/Sub"

# Post 2: RAG Systems
curl -s -X POST "$API_BASE/posts" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "title": "Implementing RAG Systems: From Concept to Production",
        "content": "# Implementing RAG Systems: From Concept to Production\n\nRetrieval-Augmented Generation (RAG) has revolutionized how we build AI applications. Let me share insights from building production RAG systems.\n\n## What is RAG?\n\nRAG combines information retrieval with language models to provide accurate, context-aware responses. Instead of relying solely on the model'\''s training data, RAG retrieves relevant documents and uses them as context.\n\n## Architecture Components\n\nA typical RAG system includes:\n\n1. **Document Store**: Vector database (e.g., Pinecone, Weaviate, or pgvector)\n2. **Embedding Model**: Converts text to vectors\n3. **Retrieval System**: Finds relevant documents based on queries\n4. **LLM**: Generates responses using retrieved context\n\n## Implementation Challenges\n\n### Chunking Strategy\nThe way you chunk documents significantly impacts retrieval quality.\n\n### Embedding Quality\nNot all embeddings are created equal. Use domain-specific models when available.\n\n### Retrieval Optimization\n- Implement re-ranking to improve relevance\n- Use hybrid search (keyword + semantic)\n- Cache frequent queries\n\n## Production Considerations\n\n1. **Latency**: RAG adds retrieval overhead—optimize your vector search\n2. **Cost**: Balance embedding costs with retrieval quality\n3. **Monitoring**: Track retrieval quality, latency, and user satisfaction\n4. **Fallbacks**: Handle cases where no relevant documents are found",
        "excerpt": "A deep dive into building production-ready RAG systems, covering architecture, challenges, and best practices.",
        "status": "published",
        "featured": true,
        "metaTitle": "Implementing RAG Systems: From Concept to Production",
        "metaDescription": "Learn how to build production-ready RAG systems with practical insights and real-world examples.",
        "tagNames": ["ai", "rag", "llm", "machine-learning", "nlp"]
    }' > /dev/null

echo "Created post: Implementing RAG Systems: From Concept to Production"

# Post 3: Kubernetes
curl -s -X POST "$API_BASE/posts" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "title": "Kubernetes Patterns for High-Availability Applications",
        "content": "# Kubernetes Patterns for High-Availability Applications\n\nRunning applications in Kubernetes requires understanding various patterns to ensure high availability and reliability.\n\n## Pod Disruption Budgets\n\nPod Disruption Budgets (PDBs) ensure a minimum number of pods remain available during voluntary disruptions.\n\n## Health Checks\n\nProper health checks are critical for maintaining application reliability.\n\n## Resource Management\n\nAlways set resource requests and limits to ensure proper scheduling and prevent resource starvation.\n\n## Deployment Strategies\n\n### Rolling Updates\nDefault strategy, updates pods gradually.\n\n### Blue-Green\nMaintain two identical environments, switch traffic instantly.\n\n### Canary\nGradually roll out changes to a subset of users.\n\n## Monitoring and Observability\n\n- Use Prometheus for metrics\n- Implement distributed tracing\n- Set up alerting for critical metrics\n- Monitor resource utilization\n\n## Best Practices\n\n1. Use HorizontalPodAutoscaler for automatic scaling\n2. Implement proper logging and monitoring\n3. Use ConfigMaps and Secrets appropriately\n4. Design for failure—assume pods will crash\n5. Implement graceful shutdowns",
        "excerpt": "Essential Kubernetes patterns and best practices for building highly available applications in production.",
        "status": "published",
        "featured": false,
        "metaTitle": "Kubernetes Patterns for High-Availability Applications",
        "metaDescription": "Learn essential Kubernetes patterns for building reliable, high-availability applications.",
        "tagNames": ["kubernetes", "devops", "containers", "infrastructure"]
    }' > /dev/null

echo "Created post: Kubernetes Patterns for High-Availability Applications"

# Post 4: Database Optimization
curl -s -X POST "$API_BASE/posts" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "title": "Optimizing Database Queries: Lessons from Production",
        "content": "# Optimizing Database Queries: Lessons from Production\n\nDatabase performance is often the bottleneck in web applications. Here are practical optimization strategies I'\''ve learned from production systems.\n\n## Indexing Strategy\n\nIndexes are your best friend, but use them wisely:\n\n- **Composite indexes**: Order matters—put most selective columns first\n- **Partial indexes**: Index only relevant rows\n- **Covering indexes**: Include columns in the index to avoid table lookups\n\n## Query Patterns to Avoid\n\n### N+1 Queries\nAlways use eager loading instead of querying in loops.\n\n### SELECT *\nAlways select only needed columns to reduce data transfer.\n\n### Missing WHERE Clauses\nAlways filter at the database level, not in application code.\n\n## Connection Pooling\n\nProper connection pooling is crucial for performance.\n\n## Query Analysis\n\nUse EXPLAIN ANALYZE to understand query execution and identify bottlenecks.\n\n## Caching Strategies\n\n1. Application-level caching for frequently accessed data\n2. Query result caching for expensive queries\n3. Proper cache invalidation strategies\n4. Cache warming for known queries",
        "excerpt": "Practical database optimization strategies learned from production systems, covering indexing, query patterns, and caching.",
        "status": "published",
        "featured": false,
        "metaTitle": "Optimizing Database Queries: Lessons from Production",
        "metaDescription": "Learn practical database optimization strategies from real production experiences.",
        "tagNames": ["database", "performance", "postgresql", "optimization"]
    }' > /dev/null

echo "Created post: Optimizing Database Queries: Lessons from Production"

# Post 5: Event-Driven Architecture
curl -s -X POST "$API_BASE/posts" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "title": "Building Event-Driven Architectures with Message Queues",
        "content": "# Building Event-Driven Architectures with Message Queues\n\nEvent-driven architectures provide excellent scalability and decoupling. Let'\''s explore how to build them effectively.\n\n## Why Event-Driven?\n\nEvent-driven architectures offer:\n- **Loose Coupling**: Services communicate through events\n- **Scalability**: Easy to scale individual components\n- **Resilience**: Failures in one service don'\''t cascade\n- **Flexibility**: Easy to add new consumers\n\n## Message Queue Options\n\n### Redis Streams\nFast and lightweight, good for high-throughput scenarios.\n\n### RabbitMQ\nMature and feature-rich, excellent management UI.\n\n### Apache Kafka\nDesigned for high throughput, excellent for event sourcing.\n\n## Implementation Patterns\n\n### Event Sourcing\nStore all changes as a sequence of events.\n\n### CQRS (Command Query Responsibility Segregation)\nSeparate read and write models for better scalability.\n\n### Saga Pattern\nManage distributed transactions across services using events.\n\n## Best Practices\n\n1. **Idempotency**: Design handlers to be idempotent\n2. **Ordering**: Understand ordering guarantees of your queue\n3. **Dead Letter Queues**: Handle failed messages gracefully\n4. **Monitoring**: Track message rates, latency, and errors\n5. **Schema Evolution**: Plan for schema changes",
        "excerpt": "Learn how to build scalable event-driven architectures using message queues, covering patterns and best practices.",
        "status": "published",
        "featured": false,
        "metaTitle": "Building Event-Driven Architectures with Message Queues",
        "metaDescription": "A comprehensive guide to building event-driven architectures with practical examples.",
        "tagNames": ["architecture", "event-driven", "message-queues", "microservices"]
    }' > /dev/null

echo "Created post: Building Event-Driven Architectures with Message Queues"

echo "All content created successfully!"

