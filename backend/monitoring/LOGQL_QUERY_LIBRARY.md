# LogQL Query Library

**Version:** 1.0  
**Last Updated:** 2025-12-22

This document contains common LogQL queries for searching and analyzing logs in Grafana Loki.

## Basic Queries

### View All Logs
```
{job="docker"}
```

### Filter by Service
```
{job="docker", service="app"}
{job="docker", service="ai-service"}
{job="docker", service="resume-worker"}
{job="docker", service="translation-worker"}
{job="docker", service="email-worker"}
{job="docker", service="whatsapp-worker"}
{job="docker", service="job-application-worker"}
{job="docker", service="docs-service"}
{job="docker", service="creative-service"}
```

### Filter by Log Level
```
{job="docker", level="error"}
{job="docker", level="warn"}
{job="docker", level="info"}
{job="docker", level="debug"}
```

### Combine Filters
```
{job="docker", service="app", level="error"}
{job="docker", service="resume-worker", level="warn"}
```

## Aggregation Queries

### Log Volume by Service
```
sum(count_over_time({job="docker"}[1m])) by (service)
```

### Error Rate by Service
```
sum(rate({job="docker", level="error"}[5m])) by (service)
```

### Warning Rate by Service
```
sum(rate({job="docker", level="warn"}[5m])) by (service)
```

### Total Log Volume Over Time
```
sum(rate({job="docker"}[1m]))
```

### Error Count by Service (Last Hour)
```
sum(count_over_time({job="docker", level="error"}[1h])) by (service)
```

## Pattern Matching

### Search for Specific Text
```
{job="docker"} |= "database"
{job="docker"} |= "connection"
{job="docker"} |= "timeout"
```

### Search with Regex
```
{job="docker"} |~ "error|exception|failure"
{job="docker"} |~ "user_id=\\d+"
```

### Exclude Patterns
```
{job="docker"} != "debug"
{job="docker"} !~ "healthcheck"
```

## Time-Based Queries

### Logs in Last 5 Minutes
```
{job="docker"} [5m]
```

### Logs in Last Hour
```
{job="docker"} [1h]
```

### Logs in Last 24 Hours
```
{job="docker"} [24h]
```

## Trace ID Queries

### Find All Logs for a Trace
```
{job="docker", trace_id="550e8400-e29b-41d4-a716-446655440000"}
```

### Find Logs with Any Trace ID
```
{job="docker"} | json | trace_id != ""
```

## Request Tracing

### Find All Logs for a Request
```
{job="docker", request_id="550e8400-e29b-41d4-a716-446655440000"}
```

### HTTP Requests by Status Code
```
{job="docker", service="app"} | json | status_code="500"
{job="docker", service="app"} | json | status_code="404"
```

### Slow Requests (> 1 second)
```
{job="docker", service="app"} | json | duration_ms > 1000
```

## Error Analysis

### Top Error Messages
```
topk(10, sum(count_over_time({job="docker", level="error"}[1h])) by (message))
```

### Errors by Service (Last Hour)
```
sum(count_over_time({job="docker", level="error"}[1h])) by (service)
```

### Errors with Stack Traces
```
{job="docker", level="error"} | json | stack_trace != ""
```

### Database Errors
```
{job="docker", level="error"} |~ "database|postgres|connection"
```

### Authentication Errors
```
{job="docker", level="error"} |~ "auth|authentication|unauthorized|forbidden"
```

## Service-Specific Queries

### Resume Worker - Job Processing
```
{job="docker", service="resume-worker"} | json | message=~".*resume.*"
```

### Translation Worker - Translation Jobs
```
{job="docker", service="translation-worker"} | json | message=~".*translation.*"
```

### Email Worker - Email Sending
```
{job="docker", service="email-worker"} | json | message=~".*email.*"
```

### Job Application Worker - Applications
```
{job="docker", service="job-application-worker"} | json | message=~".*application.*"
```

## User Activity

### Logs for Specific User
```
{job="docker"} | json | user_id="123"
```

### User Activity by Service
```
{job="docker"} | json | user_id != "" | sum(count_over_time([1h])) by (service, user_id)
```

## Performance Queries

### Average Request Duration by Service
```
{job="docker", service="app"} | json | avg(duration_ms) by (service)
```

### Requests Slower Than 500ms
```
{job="docker", service="app"} | json | duration_ms > 500
```

### Request Rate by Endpoint
```
sum(rate({job="docker", service="app"} | json | path != "" [5m])) by (path)
```

## Health Checks

### Service Startup Events
```
{job="docker"} |~ "started|starting|ready"
```

### Service Shutdown Events
```
{job="docker"} |~ "shutdown|stopping|stopped"
```

### Health Check Logs
```
{job="docker"} |~ "health|healthcheck|ping"
```

## Advanced Queries

### Error Rate Percentage
```
sum(rate({job="docker", level="error"}[5m])) / sum(rate({job="docker"}[5m])) * 100
```

### Errors per Minute by Service
```
sum(rate({job="docker", level="error"}[1m])) by (service) * 60
```

### Unique Trace IDs (Last Hour)
```
count(count_over_time({job="docker"} | json | trace_id != "" [1h]) > 0) by (trace_id)
```

### Log Volume Trend
```
sum(rate({job="docker"}[5m])) by (service)
```

## Query Tips

1. **Use time ranges**: Always specify appropriate time ranges for aggregations
2. **Filter early**: Apply label filters before pattern matching for better performance
3. **Use rate() for metrics**: Use `rate()` for calculating rates over time
4. **Use count_over_time() for counts**: Use `count_over_time()` for counting log lines
5. **Combine filters**: Use multiple label filters to narrow down results
6. **Use JSON parsing**: Use `| json` to parse structured JSON logs
7. **Use regex sparingly**: Regex queries are slower, use label filters when possible

## Performance Optimization

- Use label filters (`service`, `level`) before pattern matching
- Use shorter time ranges when possible
- Use `topk()` to limit results
- Avoid very broad regex patterns
- Use `rate()` instead of `count_over_time()` for time series

## Examples for Common Tasks

### Find all errors from app service in the last hour
```
{job="docker", service="app", level="error"} [1h]
```

### Calculate error rate per service
```
sum(rate({job="docker", level="error"}[5m])) by (service)
```

### Find slow API requests
```
{job="docker", service="app"} | json | duration_ms > 1000
```

### Track a request through all services
```
{job="docker", trace_id="550e8400-e29b-41d4-a716-446655440000"}
```

### Find database connection errors
```
{job="docker", level="error"} |~ "database.*connection"
```

## References

- [LogQL Documentation](https://grafana.com/docs/loki/latest/logql/)
- [LogQL Query Examples](https://grafana.com/docs/loki/latest/logql/log_queries/)
