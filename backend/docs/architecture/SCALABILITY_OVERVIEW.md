# Scalability Documentation - Backend Architecture

## General Architecture

Scalability is the ability of a system to handle increased load by adding resources (horizontal scaling) or increasing resource capacity (vertical scaling). This document outlines our scalability strategy, benchmarks, and scaling approaches.

### Current State

- ✅ **Horizontal Scaling Ready**: Stateless services and workers
- ✅ **Container Orchestration**: Kubernetes support (partial)
- ⚠️ **Auto-scaling**: Not implemented
- ⚠️ **Load Testing**: Not performed
- ⚠️ **Performance Benchmarks**: Not documented
- ⚠️ **Scaling Documentation**: Not documented

---

## Scalability Strategy

### Horizontal Scaling (Scale Out)

**Principle**: Add more instances to handle increased load

**Components Suitable for Horizontal Scaling:**

1. **Server (Go):**
   - ✅ Stateless (no session state)
   - ✅ Can run multiple replicas
   - ✅ Load balancer distributes requests
   - **Scaling Factor**: 2-5 replicas (based on load)

2. **Workers (Go, Python, Node.js):**
   - ✅ Stateless (process messages independently)
   - ✅ Multiple workers consume from same queue
   - ✅ RabbitMQ distributes messages
   - **Scaling Factor**: 1-10 replicas per worker type (based on queue depth)

3. **Services (AI, Creative):**
   - ✅ Stateless (no shared state)
   - ✅ Can run multiple replicas
   - ✅ Load balancer distributes requests
   - **Scaling Factor**: 1-3 replicas (based on load)

**Scaling Triggers:**
- Queue depth > threshold
- CPU usage > 70%
- Memory usage > 70%
- Request latency > threshold
- Error rate > threshold

**Scaling Limits:**
- **Server**: Limited by database connection pool
- **Workers**: Limited by RabbitMQ connection limits
- **Services**: Limited by external API rate limits

---

### Vertical Scaling (Scale Up)

**Principle**: Increase resources (CPU, memory) of existing instances

**When to Use:**
- Single-threaded bottlenecks
- Memory-intensive operations
- Database connection limits
- Before horizontal scaling (cheaper)

**Components Suitable for Vertical Scaling:**

1. **Database (PostgreSQL):**
   - Increase CPU for complex queries
   - Increase memory for query cache
   - Increase storage for data growth

2. **Workers (Memory-intensive):**
   - Resume worker (PDF generation)
   - Job application worker (browser instances)

**Scaling Limits:**
- Maximum instance size (cloud provider limits)
- Cost efficiency (vertical scaling becomes expensive)
- Single point of failure (not addressed by vertical scaling)

---

## Component Scaling Strategies

### Server

**Current Configuration:**
- Single instance (or 2-3 replicas)
- Database connection pool: ~20-50 connections
- Redis connection pool: ~10-20 connections

**Scaling Strategy:**

**Horizontal:**
- Add more server replicas (2-5)
- Use load balancer (Kubernetes Service or external LB)
- Database connection pool per replica
- Redis connection pool per replica

**Bottlenecks:**
- Database connection pool (limit: ~100-200 total)
- Database query performance
- External API rate limits

**Scaling Metrics:**
- Request rate per replica: ~100-500 req/s
- CPU usage: Scale when > 70%
- Memory usage: Scale when > 70%
- Database connections: Monitor pool usage

---

### Workers

**Current Configuration:**
- Single instance per worker type
- RabbitMQ connection per worker
- Database connection pool per worker

**Scaling Strategy:**

**Horizontal:**
- Add more worker replicas (1-10 per type)
- Multiple workers consume from same queue
- RabbitMQ distributes messages (round-robin or fair dispatch)
- Each worker has its own database connection pool

**Queue-Based Scaling:**
- Scale based on queue depth
- Target: Queue depth < 100 messages
- Scale up when queue depth > 500
- Scale down when queue depth < 10

**Worker-Specific Considerations:**

1. **Email Worker:**
   - SMTP connection limits
   - Email sending rate limits
   - Scaling: 1-5 replicas

2. **WhatsApp Worker:**
   - WhatsApp API rate limits
   - Session management
   - Scaling: 1-3 replicas

3. **Translation Worker:**
   - Translation API rate limits
   - Cost per translation
   - Scaling: 1-5 replicas

4. **Resume Worker:**
   - Memory-intensive (PDF generation)
   - AI service rate limits
   - Scaling: 1-3 replicas

5. **Job Application Worker:**
   - Browser instances (memory-intensive)
   - Website rate limits
   - Scaling: 1-5 replicas

**Scaling Metrics:**
- Queue depth: Scale when > 500
- Processing rate: Monitor jobs/second
- Error rate: Scale if errors due to overload

---

### Services (AI, Creative)

**Current Configuration:**
- Single instance per service
- External API rate limits
- No internal rate limiting

**Scaling Strategy:**

**Horizontal:**
- Add more service replicas (1-3)
- Use load balancer
- Each replica has its own external API connections

**Bottlenecks:**
- External API rate limits (OpenAI, Anthropic, etc.)
- Cost per API call
- Response time (external APIs)

**Scaling Metrics:**
- Request rate per replica: ~10-50 req/s
- CPU usage: Scale when > 70%
- Memory usage: Scale when > 70%
- External API rate limits: Monitor usage

---

### Database (PostgreSQL)

**Scaling Strategy:**

**Vertical Scaling:**
- Increase CPU (for complex queries)
- Increase memory (for query cache)
- Increase storage (for data growth)

**Horizontal Scaling (Read Replicas):**
- Create read replicas for read-heavy workloads
- Route read queries to replicas
- Write queries to primary
- Replication lag monitoring

**Connection Pooling:**
- Limit connections per application instance
- Use connection pooler (PgBouncer)
- Monitor connection usage

**Scaling Metrics:**
- Query performance: p95 latency < 100ms
- Connection pool usage: < 80%
- Replication lag: < 1 second
- CPU usage: < 70%
- Memory usage: < 70%

---

### Message Queue (RabbitMQ)

**Scaling Strategy:**

**Vertical Scaling:**
- Increase memory (for message storage)
- Increase CPU (for message processing)
- Increase disk (for message persistence)

**Horizontal Scaling (Clustering):**
- RabbitMQ cluster (3 nodes)
- Queue mirroring across nodes
- High availability

**Queue Partitioning:**
- Partition queues by user ID or job type
- Distribute load across partitions
- Independent scaling per partition

**Scaling Metrics:**
- Queue depth: Alert when > 1000
- Message rate: Monitor throughput
- Connection count: Monitor limits
- Memory usage: < 70%

---

## Performance Benchmarks

### Target Metrics

**Server:**
- Request rate: 500 req/s per replica
- Latency (p95): < 200ms
- Latency (p99): < 500ms
- Error rate: < 0.1%

**Workers:**
- Processing rate: 10-100 jobs/s per replica (depending on worker)
- Processing latency (p95): < 5s
- Error rate: < 1%

**Services:**
- Request rate: 50 req/s per replica
- Latency (p95): < 2s (depends on external APIs)
- Error rate: < 1%

**Database:**
- Query latency (p95): < 100ms
- Connection pool usage: < 80%
- CPU usage: < 70%

### Load Testing Plan

**Tools:**
- `k6` for HTTP load testing
- `artillery` for API load testing
- `wrk` for simple HTTP benchmarking
- Custom scripts for queue load testing

**Test Scenarios:**

1. **Normal Load:**
   - Baseline traffic
   - Measure baseline metrics
   - Identify bottlenecks

2. **Peak Load (2x Normal):**
   - Double the traffic
   - Measure performance degradation
   - Identify scaling needs

3. **Stress Test (5x Normal):**
   - 5x the traffic
   - Measure failure points
   - Identify maximum capacity

4. **Spike Test:**
   - Sudden traffic increase
   - Measure recovery time
   - Test auto-scaling

**Metrics to Measure:**
- Request rate (req/s)
- Latency (p50, p95, p99)
- Error rate
- CPU usage
- Memory usage
- Queue depth
- Database connections
- External API calls

---

## Auto-Scaling Strategy

### Kubernetes Horizontal Pod Autoscaler (HPA)

**Server HPA:**
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: server-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: server
  minReplicas: 2
  maxReplicas: 5
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 70
```

**Worker HPA (Queue-Based):**
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: email-worker-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: email-worker
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: External
    external:
      metric:
        name: queue_depth
      target:
        type: AverageValue
        averageValue: "100"
```

**Scaling Policies:**
- Scale up: When metric > target for 2 minutes
- Scale down: When metric < target for 5 minutes
- Cooldown: 1 minute between scaling events

---

## Sharding Strategy

### When Sharding is Needed

**Database Sharding:**
- Database size > 100GB
- Query performance degradation
- Connection pool exhaustion
- Single database becomes bottleneck

**Queue Sharding:**
- Queue depth consistently > 10,000
- Single queue becomes bottleneck
- Need for independent scaling

### Sharding Approaches

**Database Sharding:**

1. **User ID Sharding:**
   - Shard by user ID hash
   - Each shard handles subset of users
   - Cross-shard queries require aggregation

2. **Geographic Sharding:**
   - Shard by region/country
   - Data locality benefits
   - Cross-shard queries for global features

3. **Functional Sharding:**
   - Separate databases for different domains
   - Projects DB, Users DB, etc.
   - Simpler than user-based sharding

**Queue Sharding:**

1. **Queue Partitioning:**
   - Partition by user ID or job type
   - Multiple queues, multiple workers
   - Independent scaling per partition

2. **Priority Queues:**
   - Separate queues by priority
   - High-priority queue processed faster
   - Low-priority queue can backlog

---

## Capacity Planning

### Resource Requirements

**Server (per replica):**
- CPU: 500m-1000m (0.5-1 core)
- Memory: 512Mi-1Gi
- Database connections: 20-50
- Redis connections: 10-20

**Workers (per replica):**
- CPU: 200m-500m (0.2-0.5 core)
- Memory: 256Mi-512Mi (1Gi for memory-intensive)
- Database connections: 5-10
- RabbitMQ connections: 1

**Services (per replica):**
- CPU: 300m-1000m (0.3-1 core)
- Memory: 512Mi-1Gi
- External API connections: 10-50

**Database:**
- CPU: 2-4 cores
- Memory: 4-8Gi
- Storage: 100GB-1TB (grows with data)
- Connections: 100-200

### Growth Projections

**Year 1:**
- Users: 1,000-10,000
- Requests/day: 100K-1M
- Database size: 10-50GB
- Scaling: 2-3 server replicas, 1-2 worker replicas

**Year 2:**
- Users: 10,000-100,000
- Requests/day: 1M-10M
- Database size: 50-200GB
- Scaling: 3-5 server replicas, 2-5 worker replicas

**Year 3:**
- Users: 100,000-1M
- Requests/day: 10M-100M
- Database size: 200GB-1TB
- Scaling: 5-10 server replicas, 5-10 worker replicas, read replicas

---

## Monitoring and Alerts

### Scaling Metrics

**Metrics to Monitor:**
- Request rate (req/s)
- Queue depth
- CPU usage (%)
- Memory usage (%)
- Database connections
- Latency (p95, p99)
- Error rate (%)

### Scaling Alerts

**Scale Up Alerts:**
- CPU usage > 70% for 5 minutes
- Memory usage > 70% for 5 minutes
- Queue depth > 500 for 5 minutes
- Latency p95 > threshold for 5 minutes

**Scale Down Alerts:**
- CPU usage < 30% for 15 minutes
- Memory usage < 30% for 15 minutes
- Queue depth < 10 for 15 minutes

**Capacity Alerts:**
- Approaching max replicas
- Database connection pool > 80%
- External API rate limit > 80%

---

## Cost Optimization

### Scaling Costs

**Horizontal Scaling:**
- Linear cost increase
- Pay for what you use
- Can scale down during low traffic

**Vertical Scaling:**
- Exponential cost increase
- Pay for peak capacity
- Can't scale down easily

### Cost Optimization Strategies

1. **Right-Sizing:**
   - Don't over-provision
   - Monitor actual usage
   - Adjust resource requests/limits

2. **Auto-Scaling:**
   - Scale down during low traffic
   - Scale up proactively
   - Use predictive scaling

3. **Reserved Instances:**
   - Commit to 1-3 year terms
   - 30-50% cost savings
   - For predictable workloads

4. **Spot Instances:**
   - 50-90% cost savings
   - For fault-tolerant workloads
   - Workers can use spot instances

---

## Implementation Roadmap

### Phase 1: Documentation and Baseline (Week 1-2)

- Document current scaling strategy
- Perform baseline load testing
- Establish performance benchmarks
- Document resource requirements

### Phase 2: Auto-Scaling (Week 3-4)

- Implement Kubernetes HPA for server
- Implement Kubernetes HPA for workers
- Configure scaling policies
- Test auto-scaling behavior

### Phase 3: Load Testing (Week 5-6)

- Perform normal load testing
- Perform peak load testing
- Perform stress testing
- Document results and bottlenecks

### Phase 4: Optimization (Week 7-8)

- Optimize based on load test results
- Right-size resources
- Implement read replicas (if needed)
- Cost optimization

---

## Summary

**Current State:**
- ✅ Horizontal scaling ready (stateless components)
- ✅ Container orchestration (Kubernetes partial)
- ❌ Auto-scaling (not implemented)
- ❌ Load testing (not performed)
- ❌ Performance benchmarks (not documented)

**Priority:**
1. **Documentation** - Document scaling strategy and benchmarks
2. **Auto-Scaling** - Implement Kubernetes HPA
3. **Load Testing** - Perform comprehensive load testing
4. **Optimization** - Optimize based on results

**Scaling Strategy:**
- **Horizontal**: Preferred (stateless components)
- **Vertical**: For database and memory-intensive workers
- **Auto-Scaling**: Based on CPU, memory, queue depth

**Key Metrics:**
- Request rate, latency, error rate
- Queue depth, CPU, memory usage
- Database connections, external API usage
