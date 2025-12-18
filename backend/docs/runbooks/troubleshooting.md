# Troubleshooting Common Issues

## Overview
This runbook covers common issues and their solutions for the backend services and workers.

## General Troubleshooting Steps

### 1. Check Health Status
```bash
# Server
curl https://api.woragis.com/healthz

# Workers
curl http://{worker-name}:8080/healthz

# Services
curl http://{service-name}:8000/healthz
```

### 2. Check Logs
```bash
# Railway
railway logs --service {service-name}

# Kubernetes
kubectl logs -f deployment/{service-name}

# Docker Compose
docker-compose logs -f {service-name}
```

### 3. Check Metrics
```bash
# Prometheus metrics
curl http://{service-name}:8080/metrics

# Check specific metrics
curl http://{service-name}:8080/metrics | grep "error_rate"
```

### 4. Check Dependencies
- Database: Connection, query performance
- Redis: Connection, memory usage
- RabbitMQ: Connection, queue depths
- External APIs: Status, rate limits

## Common Issues

### Issue 1: Service Unhealthy

#### Symptoms
- Health check returns 503
- Service not responding to requests
- Errors in logs

#### Diagnosis
```bash
# Check health check response
curl -v http://{service-name}:8080/healthz

# Check logs
railway logs --service {service-name} | grep -i error

# Check dependencies
# Database
psql $DATABASE_URL -c "SELECT 1"

# Redis
redis-cli -u $REDIS_URL ping

# RabbitMQ
rabbitmqctl status
```

#### Solutions

##### Database Connection Failed
```bash
# Check database is running
docker ps | grep postgres

# Check connection string
echo $DATABASE_URL

# Test connection
psql $DATABASE_URL -c "SELECT 1"

# Fix: Update DATABASE_URL or restart database
```

##### Redis Connection Failed
```bash
# Check Redis is running
docker ps | grep redis

# Check connection string
echo $REDIS_URL

# Test connection
redis-cli -u $REDIS_URL ping

# Fix: Update REDIS_URL or restart Redis
```

##### RabbitMQ Connection Failed
```bash
# Check RabbitMQ is running
docker ps | grep rabbitmq

# Check connection string
echo $RABBITMQ_URL

# Test connection
rabbitmqctl status

# Fix: Update RABBITMQ_URL or restart RabbitMQ
```

### Issue 2: High Error Rate

#### Symptoms
- High error rate in metrics
- Many errors in logs
- Users reporting issues

#### Diagnosis
```bash
# Check error rate
curl http://{service-name}:8080/metrics | grep "error_rate"

# Check error logs
railway logs --service {service-name} | grep -i error | tail -100

# Check specific error patterns
railway logs --service {service-name} | grep "database.*error"
railway logs --service {service-name} | grep "timeout"
railway logs --service {service-name} | grep "connection.*refused"
```

#### Solutions

##### Database Errors
```bash
# Check database health
psql $DATABASE_URL -c "SELECT count(*) FROM pg_stat_activity"

# Check connection pool
# (if using connection pooling)

# Fix: Increase connection pool, optimize queries, restart database
```

##### External API Errors
```bash
# Check external API status
curl https://api.external-service.com/health

# Check rate limits
# (check API documentation)

# Fix: Implement retries, circuit breakers, rate limiting
```

##### Timeout Errors
```bash
# Check timeout configuration
grep -r "timeout" {service-name}/src/

# Check external API response times
curl -w "@curl-format.txt" https://api.external-service.com/endpoint

# Fix: Increase timeout, optimize queries, add caching
```

### Issue 3: Workers Not Processing Jobs

#### Symptoms
- Queue depth increasing
- No job processing logs
- Workers healthy but idle

#### Diagnosis
```bash
# Check queue depth
rabbitmqctl list_queues name messages | grep "{queue-name}"

# Check worker logs
railway logs --service {worker-name} | grep -i "consume\|process"

# Check RabbitMQ connection
curl http://{worker-name}:8080/healthz

# Check if workers are consuming
rabbitmqctl list_consumers | grep "{queue-name}"
```

#### Solutions

##### Workers Not Connected to RabbitMQ
```bash
# Check RabbitMQ connection
rabbitmqctl list_connections

# Check worker configuration
echo $RABBITMQ_URL

# Fix: Update RABBITMQ_URL, restart worker
```

##### Workers Consuming But Not Processing
```bash
# Check worker logs for errors
railway logs --service {worker-name} | grep -i error

# Check message format
rabbitmqadmin get queue={queue-name} ackmode=ack_requeue_false count=1

# Fix: Fix processing logic, update message format
```

##### Queue Configuration Issues
```bash
# Check queue exists
rabbitmqctl list_queues name | grep "{queue-name}"

# Check queue bindings
rabbitmqctl list_bindings | grep "{queue-name}"

# Fix: Create queue, fix bindings
```

### Issue 4: High Latency

#### Symptoms
- Slow response times
- High p95/p99 latency
- Users reporting slowness

#### Diagnosis
```bash
# Check latency metrics
curl http://{service-name}:8080/metrics | grep "latency\|duration"

# Check slow queries (database)
psql $DATABASE_URL -c "SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10"

# Check external API latency
curl -w "@curl-format.txt" https://api.external-service.com/endpoint
```

#### Solutions

##### Database Slow Queries
```bash
# Identify slow queries
psql $DATABASE_URL -c "SELECT * FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10"

# Check indexes
psql $DATABASE_URL -c "\d {table-name}"

# Fix: Add indexes, optimize queries, add caching
```

##### External API Slow Responses
```bash
# Check external API status
curl https://api.external-service.com/health

# Check rate limiting
# (may be throttling)

# Fix: Add caching, implement async processing, use faster APIs
```

##### High CPU/Memory Usage
```bash
# Check resource usage
# Railway: Check dashboard
# Kubernetes: kubectl top pod {pod-name}

# Fix: Scale up, optimize code, add caching
```

### Issue 5: Memory Leaks

#### Symptoms
- Memory usage increasing over time
- Service crashes with OOM errors
- Performance degradation

#### Diagnosis
```bash
# Check memory usage over time
# Railway: Check metrics dashboard
# Kubernetes: kubectl top pod {pod-name} --watch

# Check for memory leaks in code
# (use profiling tools)
```

#### Solutions

##### Go Services
```bash
# Enable memory profiling
# Add to code:
import _ "net/http/pprof"

# Check memory profile
go tool pprof http://localhost:8080/debug/pprof/heap
```

##### Python Services
```bash
# Use memory profiler
pip install memory-profiler

# Profile code
python -m memory_profiler script.py
```

##### General Fixes
- Fix memory leaks in code
- Increase memory limits (temporary)
- Restart service periodically (temporary)
- Optimize data structures

### Issue 6: Database Connection Pool Exhausted

#### Symptoms
- Database connection errors
- "too many connections" errors
- Service degradation

#### Diagnosis
```bash
# Check active connections
psql $DATABASE_URL -c "SELECT count(*) FROM pg_stat_activity"

# Check max connections
psql $DATABASE_URL -c "SHOW max_connections"

# Check connection pool configuration
grep -r "max.*connection" {service-name}/
```

#### Solutions

##### Increase Connection Pool
```bash
# Update connection pool size in code
# Example: max_connections=50 -> max_connections=100
```

##### Reduce Connection Usage
```bash
# Check for connection leaks
# (connections not being closed)

# Fix: Ensure connections are closed, use connection pooling
```

##### Increase Database Max Connections
```bash
# Update PostgreSQL configuration
# max_connections = 200

# Restart database
```

### Issue 7: RabbitMQ Queue Backlog

#### Symptoms
- Queue depth increasing
- Jobs not being processed fast enough
- High latency for async operations

#### Diagnosis
```bash
# Check queue depths
rabbitmqctl list_queues name messages | sort -k2 -n

# Check processing rate
# (check metrics or logs)

# Check worker count
# Railway: Check dashboard
# Kubernetes: kubectl get deployment {worker-name}
```

#### Solutions

##### Scale Workers
```bash
# Railway: Increase instance count
# Kubernetes: kubectl scale deployment {worker-name} --replicas=5
```

##### Optimize Processing
```bash
# Check processing time per job
# (check logs or metrics)

# Fix: Optimize processing logic, parallelize work
```

##### Increase Prefetch
```bash
# Increase prefetch count (if applicable)
# (allows workers to process multiple messages concurrently)
```

## Service-Specific Issues

### Server Issues

#### High Request Rate
- **Solution**: Scale server, add caching, optimize queries

#### Authentication Issues
- **Solution**: Check JWT configuration, check token expiration

#### Database Query Performance
- **Solution**: Add indexes, optimize queries, add caching

### Worker Issues

#### Email Worker: SMTP Errors
- **Solution**: Check SMTP configuration, check rate limits

#### Translation Worker: API Errors
- **Solution**: Check API keys, check rate limits, implement retries

#### Resume Worker: PDF Generation Errors
- **Solution**: Check memory limits, check template files

#### Job Application Worker: Scraping Errors
- **Solution**: Check website changes, update selectors, check rate limits

### Service Issues

#### AI Service: LLM API Errors
- **Solution**: Check API keys, check rate limits, implement retries

#### Creative Service: Generation Errors
- **Solution**: Check API keys, check rate limits, check input validation

## Escalation

### When to Escalate
- Issue persists after following troubleshooting steps
- Issue affects production users
- Issue requires code changes
- Issue requires infrastructure changes

### Escalation Steps
1. Document issue (symptoms, diagnosis, attempted solutions)
2. Check if issue is known (check issue tracker)
3. Escalate to team lead or on-call engineer
4. Create incident report (if production issue)

## Prevention

### Monitoring
- Set up alerts for common issues
- Monitor metrics regularly
- Review logs periodically

### Testing
- Test changes in staging first
- Run load tests before deployment
- Test failure scenarios

### Documentation
- Document common issues and solutions
- Update runbooks when procedures change
- Share knowledge with team

## Related Documentation
- [Monitoring DLQ](./monitoring-dlq.md) - DLQ troubleshooting
- [Deploying Services and Workers](./deploying-services.md) - Deployment issues
- [Architecture Decision Records](../adr/) - Architecture context
