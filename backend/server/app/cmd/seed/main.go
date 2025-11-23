package main

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	postsdomain "github.com/woragis/backend/server/app/internal/domains/posts"
	testimonialsdomain "github.com/woragis/backend/server/app/internal/domains/testimonials"
)

func main() {
	// Get database connection string from environment
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get user ID from environment or use default
	userIDStr := os.Getenv("USER_ID")
	if userIDStr == "" {
		userIDStr = "6ad0d828-f605-45fc-a545-3441e17a015c" // Default user ID from projects
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Fatalf("Invalid USER_ID: %v", err)
	}

	ctx := context.Background()

	// Initialize repositories
	testimonialRepo := testimonialsdomain.NewGormRepository(db)
	testimonialService := testimonialsdomain.NewService(testimonialRepo, nil)

	postRepo := postsdomain.NewGormRepository(db)
	postService := postsdomain.NewService(postRepo, nil)

	// Create testimonials
	testimonials := []testimonialsdomain.CreateTestimonialRequest{
		{
			AuthorName:    "Sarah Chen",
			AuthorRole:    "Senior Software Engineer",
			AuthorCompany: "TechCorp",
			Content:       "Working with Jezreel was an absolute pleasure. His deep understanding of distributed systems and ability to architect scalable solutions is impressive. He consistently delivered high-quality code and was always willing to share knowledge with the team. I'd work with him again in a heartbeat.",
			Rating:        intPtr(5),
			LinkedInURL:   "https://linkedin.com/in/sarahchen",
			DisplayOrder:  1,
		},
		{
			AuthorName:    "Michael Rodriguez",
			AuthorRole:    "CTO",
			AuthorCompany: "StartupXYZ",
			Content:       "Jezreel is one of the most technically skilled developers I've had the privilege to work with. His expertise in Go, Redis, and microservices architecture helped us build a robust system that handles millions of requests daily. His attention to detail and problem-solving skills are exceptional.",
			Rating:        intPtr(5),
			LinkedInURL:   "https://linkedin.com/in/michaelrodriguez",
			DisplayOrder:  2,
		},
		{
			AuthorName:    "Emily Watson",
			AuthorRole:    "Product Manager",
			AuthorCompany: "InnovateLabs",
			Content:       "Jezreel's ability to translate complex technical requirements into elegant solutions is remarkable. He's not just a great coder—he's a great communicator who can explain technical concepts to non-technical stakeholders. His work on our RAG system was groundbreaking.",
			Rating:        intPtr(5),
			LinkedInURL:   "https://linkedin.com/in/emilywatson",
			DisplayOrder:  3,
		},
		{
			AuthorName:    "David Kim",
			AuthorRole:    "Lead DevOps Engineer",
			AuthorCompany: "CloudScale Inc",
			Content:       "I've collaborated with Jezreel on several infrastructure projects, and his knowledge of Kubernetes, Docker, and cloud architecture is top-notch. He's proactive, reliable, and always thinking about scalability and reliability. A true professional.",
			Rating:        intPtr(5),
			LinkedInURL:   "https://linkedin.com/in/davidkim",
			DisplayOrder:  4,
		},
		{
			AuthorName:    "Alexandra Thompson",
			AuthorRole:    "Engineering Manager",
			AuthorCompany: "DataFlow Systems",
			Content:       "Jezreel brings a unique combination of technical depth and practical problem-solving. His work on our pub/sub messaging system using Redis was instrumental in improving our system's performance. He's a team player who elevates everyone around him.",
			Rating:        intPtr(5),
			LinkedInURL:   "https://linkedin.com/in/alexandrathompson",
			DisplayOrder:  5,
		},
	}

	log.Println("Creating testimonials...")
	for i, req := range testimonials {
		testimonial, err := testimonialService.CreateTestimonial(ctx, userID, req)
		if err != nil {
			log.Printf("Error creating testimonial %d: %v", i+1, err)
			continue
		}
		log.Printf("Created testimonial: %s by %s", testimonial.ID, testimonial.AuthorName)

		// Approve the testimonial
		approvedStatus := testimonialsdomain.TestimonialStatusApproved
		_, err = testimonialService.UpdateTestimonial(ctx, testimonial.ID, userID, testimonialsdomain.UpdateTestimonialRequest{
			Status: &approvedStatus,
		})
		if err != nil {
			log.Printf("Error approving testimonial: %v", err)
		}
	}

	// Create blog posts
	posts := []postsdomain.CreatePostRequest{
		{
			Title:   "Building Scalable Microservices with Go and Redis Pub/Sub",
			Content: `# Building Scalable Microservices with Go and Redis Pub/Sub

In modern distributed systems, efficient communication between services is crucial. Redis Pub/Sub provides a lightweight, fast messaging solution that's perfect for microservices architectures.

## Why Redis Pub/Sub?

Redis Pub/Sub offers several advantages:
- **Low Latency**: Sub-millisecond message delivery
- **Decoupling**: Services don't need to know about each other
- **Scalability**: Handle millions of messages per second
- **Simplicity**: Easy to implement and maintain

## Implementation Pattern

Here's a common pattern I use in production:

` + "```" + `go
type MessageBroker struct {
    client *redis.Client
    pubsub *redis.PubSub
}

func (mb *MessageBroker) Publish(channel string, message []byte) error {
    return mb.client.Publish(context.Background(), channel, message).Err()
}

func (mb *MessageBroker) Subscribe(channel string, handler func([]byte)) error {
    pubsub := mb.client.Subscribe(context.Background(), channel)
    mb.pubsub = pubsub
    
    ch := pubsub.Channel()
    for msg := range ch {
        handler([]byte(msg.Payload))
    }
    return nil
}
` + "```" + `

## Best Practices

1. **Error Handling**: Always implement retry logic for failed messages
2. **Message Format**: Use structured formats like JSON or Protocol Buffers
3. **Monitoring**: Track message rates and latency
4. **Backpressure**: Implement rate limiting to prevent overload

## Real-World Example

In a recent project, I used Redis Pub/Sub to handle real-time notifications across 50+ microservices. The system processes over 100,000 messages per second with sub-10ms latency.

The key was implementing proper connection pooling, message batching, and circuit breakers to handle failures gracefully.`,
			Excerpt:        "Learn how to build scalable microservices using Go and Redis Pub/Sub for efficient inter-service communication.",
			Status:         postsdomain.PostStatusPublished,
			Featured:       true,
			MetaTitle:      "Building Scalable Microservices with Go and Redis Pub/Sub",
			MetaDescription: "A comprehensive guide to implementing microservices communication using Go and Redis Pub/Sub patterns.",
			TagNames:       []string{"go", "redis", "microservices", "distributed-systems"},
		},
		{
			Title:   "Implementing RAG Systems: From Concept to Production",
			Content: `# Implementing RAG Systems: From Concept to Production

Retrieval-Augmented Generation (RAG) has revolutionized how we build AI applications. Let me share insights from building production RAG systems.

## What is RAG?

RAG combines information retrieval with language models to provide accurate, context-aware responses. Instead of relying solely on the model's training data, RAG retrieves relevant documents and uses them as context.

## Architecture Components

A typical RAG system includes:

1. **Document Store**: Vector database (e.g., Pinecone, Weaviate, or pgvector)
2. **Embedding Model**: Converts text to vectors (e.g., OpenAI, Cohere, or open-source models)
3. **Retrieval System**: Finds relevant documents based on queries
4. **LLM**: Generates responses using retrieved context

## Implementation Challenges

### Chunking Strategy
The way you chunk documents significantly impacts retrieval quality:
- **Fixed-size chunks**: Simple but may split important context
- **Semantic chunks**: Better context preservation but more complex
- **Hybrid approach**: Combine both for optimal results

### Embedding Quality
Not all embeddings are created equal:
- Use domain-specific models when available
- Fine-tune embeddings for your use case
- Consider multi-vector approaches for complex documents

### Retrieval Optimization
- Implement re-ranking to improve relevance
- Use hybrid search (keyword + semantic)
- Cache frequent queries
- Implement query expansion

## Production Considerations

1. **Latency**: RAG adds retrieval overhead—optimize your vector search
2. **Cost**: Balance embedding costs with retrieval quality
3. **Monitoring**: Track retrieval quality, latency, and user satisfaction
4. **Fallbacks**: Handle cases where no relevant documents are found

## Lessons Learned

In my experience building RAG systems, the biggest wins come from:
- Careful chunking and metadata design
- Iterative improvement of retrieval strategies
- Comprehensive testing with real user queries
- Monitoring and continuous optimization

The key is to start simple and iterate based on real-world performance.`,
			Excerpt:        "A deep dive into building production-ready RAG systems, covering architecture, challenges, and best practices.",
			Status:         postsdomain.PostStatusPublished,
			Featured:       true,
			MetaTitle:      "Implementing RAG Systems: From Concept to Production",
			MetaDescription: "Learn how to build production-ready RAG systems with practical insights and real-world examples.",
			TagNames:       []string{"ai", "rag", "llm", "machine-learning", "nlp"},
		},
		{
			Title:   "Kubernetes Patterns for High-Availability Applications",
			Content: `# Kubernetes Patterns for High-Availability Applications

Running applications in Kubernetes requires understanding various patterns to ensure high availability and reliability.

## Pod Disruption Budgets

Pod Disruption Budgets (PDBs) ensure a minimum number of pods remain available during voluntary disruptions:

` + "```" + `yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: app-pdb
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: my-app
` + "```" + `

## Health Checks

Proper health checks are critical:

` + "```" + `yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
` + "```" + `

## Resource Management

Always set resource requests and limits:

` + "```" + `yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
` + "```" + `

## Deployment Strategies

### Rolling Updates
Default strategy, updates pods gradually.

### Blue-Green
Maintain two identical environments, switch traffic instantly.

### Canary
Gradually roll out changes to a subset of users.

## Monitoring and Observability

- Use Prometheus for metrics
- Implement distributed tracing
- Set up alerting for critical metrics
- Monitor resource utilization

## Best Practices

1. Use HorizontalPodAutoscaler for automatic scaling
2. Implement proper logging and monitoring
3. Use ConfigMaps and Secrets appropriately
4. Design for failure—assume pods will crash
5. Implement graceful shutdowns

## Real-World Example

In a production system handling 10M+ requests daily, I implemented:
- Multi-zone deployments for redundancy
- PodDisruptionBudgets to ensure availability
- Comprehensive health checks
- Automatic scaling based on CPU and custom metrics

The result: 99.99% uptime with zero-downtime deployments.`,
			Excerpt:        "Essential Kubernetes patterns and best practices for building highly available applications in production.",
			Status:         postsdomain.PostStatusPublished,
			Featured:       false,
			MetaTitle:      "Kubernetes Patterns for High-Availability Applications",
			MetaDescription: "Learn essential Kubernetes patterns for building reliable, high-availability applications.",
			TagNames:       []string{"kubernetes", "devops", "containers", "infrastructure"},
		},
		{
			Title:   "Optimizing Database Queries: Lessons from Production",
			Content: `# Optimizing Database Queries: Lessons from Production

Database performance is often the bottleneck in web applications. Here are practical optimization strategies I've learned from production systems.

## Indexing Strategy

Indexes are your best friend, but use them wisely:

- **Composite indexes**: Order matters—put most selective columns first
- **Partial indexes**: Index only relevant rows (e.g., WHERE status = 'active')
- **Covering indexes**: Include columns in the index to avoid table lookups

## Query Patterns to Avoid

### N+1 Queries
` + "```" + `go
// Bad: N+1 queries
for _, user := range users {
    orders := db.Where("user_id = ?", user.ID).Find(&orders)
}

// Good: Eager loading
db.Preload("Orders").Find(&users)
` + "```" + `

### SELECT *
Always select only needed columns to reduce data transfer.

### Missing WHERE Clauses
Always filter at the database level, not in application code.

## Connection Pooling

Proper connection pooling is crucial:
- Set appropriate min/max connections
- Monitor connection usage
- Use read replicas for read-heavy workloads
- Implement connection timeouts

## Query Analysis

Use EXPLAIN ANALYZE to understand query execution:
- Look for sequential scans (bad)
- Check index usage
- Identify expensive operations
- Monitor query execution time

## Caching Strategies

1. **Application-level caching**: Cache frequently accessed data
2. **Query result caching**: Cache expensive query results
3. **Cache invalidation**: Implement proper invalidation strategies
4. **Cache warming**: Pre-populate cache for known queries

## Real-World Optimization

In a recent project, I reduced query time from 2.5s to 50ms by:
- Adding composite indexes on frequently queried columns
- Implementing query result caching with Redis
- Using database views for complex queries
- Optimizing JOIN operations

The key is to measure first, optimize based on data, and continuously monitor.`,
			Excerpt:        "Practical database optimization strategies learned from production systems, covering indexing, query patterns, and caching.",
			Status:         postsdomain.PostStatusPublished,
			Featured:       false,
			MetaTitle:      "Optimizing Database Queries: Lessons from Production",
			MetaDescription: "Learn practical database optimization strategies from real production experiences.",
			TagNames:       []string{"database", "performance", "postgresql", "optimization"},
		},
		{
			Title:   "Building Event-Driven Architectures with Message Queues",
			Content: `# Building Event-Driven Architectures with Message Queues

Event-driven architectures provide excellent scalability and decoupling. Let's explore how to build them effectively.

## Why Event-Driven?

Event-driven architectures offer:
- **Loose Coupling**: Services communicate through events
- **Scalability**: Easy to scale individual components
- **Resilience**: Failures in one service don't cascade
- **Flexibility**: Easy to add new consumers

## Message Queue Options

### Redis Streams
- Fast and lightweight
- Good for high-throughput scenarios
- Built-in consumer groups
- Persistence options

### RabbitMQ
- Mature and feature-rich
- Excellent management UI
- Supports complex routing
- Good for enterprise use cases

### Apache Kafka
- Designed for high throughput
- Excellent for event sourcing
- Strong ordering guarantees
- More complex setup

## Implementation Patterns

### Event Sourcing
Store all changes as a sequence of events:
` + "```" + `go
type Event struct {
    ID        uuid.UUID
    Type      string
    Payload   []byte
    Timestamp time.Time
}
` + "```" + `

### CQRS (Command Query Responsibility Segregation)
Separate read and write models for better scalability.

### Saga Pattern
Manage distributed transactions across services using events.

## Best Practices

1. **Idempotency**: Design handlers to be idempotent
2. **Ordering**: Understand ordering guarantees of your queue
3. **Dead Letter Queues**: Handle failed messages gracefully
4. **Monitoring**: Track message rates, latency, and errors
5. **Schema Evolution**: Plan for schema changes

## Real-World Example

I built an event-driven system processing 1M+ events daily:
- Redis Streams for real-time events
- RabbitMQ for reliable delivery
- Event sourcing for audit trails
- CQRS for read optimization

The system handles peak loads gracefully and provides excellent observability.`,
			Excerpt:        "Learn how to build scalable event-driven architectures using message queues, covering patterns and best practices.",
			Status:         postsdomain.PostStatusPublished,
			Featured:       false,
			MetaTitle:      "Building Event-Driven Architectures with Message Queues",
			MetaDescription: "A comprehensive guide to building event-driven architectures with practical examples.",
			TagNames:       []string{"architecture", "event-driven", "message-queues", "microservices"},
		},
	}

	log.Println("Creating blog posts...")
	for i, req := range posts {
		post, err := postService.CreatePost(ctx, userID, req)
		if err != nil {
			log.Printf("Error creating post %d: %v", i+1, err)
			continue
		}
		log.Printf("Created post: %s - %s", post.ID, post.Title)
	}

	log.Println("Seed completed successfully!")
}

func intPtr(i int) *int {
	return &i
}

