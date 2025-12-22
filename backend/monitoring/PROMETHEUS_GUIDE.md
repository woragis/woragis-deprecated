# Prometheus Metrics Guide

**Version:** 1.0  
**Last Updated:** 2025-12-22

## Overview

This guide covers the Prometheus metrics collection system for Woragis. Prometheus collects metrics from all services and provides them to Grafana for visualization and alerting.

## Architecture

- **Prometheus**: Metrics collection and storage (port 9090)
- **Grafana**: Visualization and dashboards (port 3000)
- **Services**: Expose `/metrics` endpoints for scraping

## Quick Start

1. **Start Prometheus:**
   ```bash
   docker-compose up -d prometheus
   ```

2. **Access Prometheus UI:**
   - URL: http://localhost:9090
   - View targets: http://localhost:9090/targets
   - Query metrics: http://localhost:9090/graph

3. **Access Grafana:**
   - URL: http://localhost:3000
   - Prometheus data source is auto-configured
   - View dashboards: Dashboards → System Overview

## Services Monitored

All services expose metrics at `/metrics`:

- `app` (main server) - port 8080
- `ai-service` - port 8000
- `creative-service` - port 8000
- `docs-service` - port 8000
- `resume-worker` - port 8000
- `translation-worker` - port 8080
- `email-worker` - port 8080
- `whatsapp-worker` - port 8080
- `job-application-worker` - port 3000

## Available Metrics

### HTTP Metrics

**Request Rate:**
```
http_requests_total
```

**Request Duration:**
```
http_request_duration_seconds_bucket
http_request_duration_seconds_sum
http_request_duration_seconds_count
```

**Status Codes:**
```
http_requests_total{status="200"}
http_requests_total{status="404"}
http_requests_total{status="500"}
```

### Queue Metrics

**Queue Depth:**
```
rabbitmq_queue_messages
rabbitmq_queue_messages_ready
rabbitmq_queue_messages_unacked
```

**DLQ Size:**
```
rabbitmq_queue_messages{queue=~".*dlq.*"}
```

### Worker Metrics

**Jobs Processed:**
```
worker_jobs_processed_total
```

**Job Duration:**
```
worker_job_duration_seconds_bucket
worker_job_duration_seconds_sum
worker_job_duration_seconds_count
```

**Jobs Failed:**
```
worker_jobs_failed_total
```

### Service Health

**Service Up/Down:**
```
up{job="app"}
up{job="ai-service"}
```

## PromQL Queries

### Request Rate by Service
```
sum(rate(http_requests_total[5m])) by (service)
```

### Error Rate
```
sum(rate(http_requests_total{status=~"5.."}[5m])) by (service)
```

### P95 Latency
```
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, service))
```

### Queue Depth
```
sum(rabbitmq_queue_messages) by (queue)
```

### Jobs Processed Rate
```
sum(rate(worker_jobs_processed_total[5m])) by (service)
```

### Service Availability
```
up{job!="prometheus"}
```

## Dashboards

### System Overview
- Total requests/sec
- Error rate
- P95 latency
- Active services
- Request rate by service
- Error rate by service
- Request latency (P50, P95, P99)
- Service health

### Queue Monitoring
- Total queue depth
- DLQ size
- Queue throughput
- Queue depth by queue
- Jobs processed by worker
- Job processing duration

## Alerting

Alert rules are configured in `monitoring/prometheus/alerts.yml`:

1. **ServiceDown** - Service is down for > 1 minute
2. **HighErrorRate** - Error rate > 10 errors/sec for 5 minutes
3. **HighLatency** - P95 latency > 1s for 5 minutes
4. **HighQueueDepth** - Queue depth > 1000 messages for 5 minutes
5. **DeadLetterQueueMessages** - Messages in DLQ
6. **LowJobProcessingRate** - Job processing rate < 0.1 jobs/sec for 10 minutes

## Configuration

### Prometheus Config

Location: `monitoring/prometheus-config.yml`

Key settings:
- **Scrape interval**: 15s
- **Evaluation interval**: 15s
- **Retention**: 30 days

### Adding New Services

1. Add service to `prometheus-config.yml`:
   ```yaml
   - job_name: 'new-service'
     metrics_path: '/metrics'
     static_configs:
       - targets: ['new-service:8080']
         labels:
           service: 'new-service'
   ```

2. Restart Prometheus:
   ```bash
   docker-compose restart prometheus
   ```

3. Verify in Prometheus UI: http://localhost:9090/targets

## Troubleshooting

### Service Not Scraping

1. **Check service is running:**
   ```bash
   docker-compose ps <service-name>
   ```

2. **Check metrics endpoint:**
   ```bash
   curl http://localhost:<port>/metrics
   ```

3. **Check Prometheus targets:**
   - Go to http://localhost:9090/targets
   - Look for service status

4. **Check Prometheus logs:**
   ```bash
   docker-compose logs prometheus
   ```

### No Metrics Appearing

1. **Verify service exposes metrics:**
   - Check service has `/metrics` endpoint
   - Verify Prometheus client library is installed

2. **Check Prometheus config:**
   - Verify service is in `prometheus-config.yml`
   - Check target URL is correct

3. **Check network:**
   - Verify services are on same Docker network
   - Test connectivity: `docker exec woragis-prometheus ping <service-name>`

### High Memory Usage

1. **Reduce retention:**
   Edit `docker-compose.yml`:
   ```yaml
   - "--storage.tsdb.retention.time=7d"  # Reduce from 30d
   ```

2. **Reduce scrape interval:**
   Edit `prometheus-config.yml`:
   ```yaml
   scrape_interval: 30s  # Increase from 15s
   ```

## Best Practices

1. **Use appropriate metric names:**
   - Follow Prometheus naming conventions
   - Use descriptive labels

2. **Set reasonable retention:**
   - 30 days for development
   - Adjust based on storage capacity

3. **Monitor Prometheus itself:**
   - Check Prometheus metrics
   - Monitor scrape failures

4. **Use recording rules:**
   - Pre-compute expensive queries
   - Reduce query load

5. **Set up alerts:**
   - Alert on critical issues
   - Use appropriate thresholds

## Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [PromQL Guide](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Grafana Prometheus Data Source](https://grafana.com/docs/grafana/latest/datasources/prometheus/)
