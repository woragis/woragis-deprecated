# Monitoring & Alerting Guide

**Last Updated:** 2025-12-22  
**Purpose:** Guide for monitoring and alerting beyond basic observability

---

## Overview

This guide covers advanced monitoring and alerting setup for production environments. Basic observability (Loki, Prometheus, Grafana) is already in place.

---

## Current Monitoring Stack

**Already Implemented:**
- ✅ Loki - Log aggregation
- ✅ Prometheus - Metrics collection
- ✅ Grafana - Visualization
- ✅ Jaeger - Distributed tracing
- ✅ Basic dashboards

**Enhancements Needed:**
- Advanced alerting rules
- Notification channels
- Incident response procedures
- Performance monitoring
- Business metrics

---

## Alert Configuration

### Critical Alerts

#### 1. Service Down

**Alert Rule:**
```yaml
# prometheus/alerts.yml
groups:
  - name: service_availability
    interval: 30s
    rules:
      - alert: ServiceDown
        expr: up{job="app"} == 0
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "Service {{ $labels.job }} is down"
          description: "Service {{ $labels.job }} has been down for more than 2 minutes"
```

**Grafana Alert:**
```yaml
# monitoring/grafana/provisioning/alerting/rules.yml
groups:
  - name: service_health
    interval: 30s
    rules:
      - uid: service_down
        title: Service Down
        condition: A
        data:
          - refId: A
            datasourceUid: prometheus
            model:
              expr: up{job="app"} == 0
        for: 2m
        annotations:
          summary: "Service is down"
```

#### 2. High Error Rate

**Alert Rule:**
```yaml
- alert: HighErrorRate
  expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
  for: 5m
  labels:
    severity: critical
  annotations:
    summary: "High error rate detected"
    description: "Error rate is {{ $value }} errors/sec"
```

**LogQL Alert (Grafana):**
```yaml
- uid: high_error_rate
  title: High Error Rate
  condition: A
  data:
    - refId: A
      datasourceUid: loki
      model:
        expr: sum(rate({job="docker", level="error"}[5m])) > 10
  for: 5m
```

#### 3. High Latency

**Alert Rule:**
```yaml
- alert: HighLatency
  expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
  for: 10m
  labels:
    severity: warning
  annotations:
    summary: "High latency detected"
    description: "95th percentile latency is {{ $value }}s"
```

#### 4. Database Connection Issues

**Alert Rule:**
```yaml
- alert: DatabaseConnectionFailure
  expr: pg_up == 0 OR rate(pg_stat_database_numbackends[5m]) == 0
  for: 1m
  labels:
    severity: critical
  annotations:
    summary: "Database connection failure"
```

#### 5. Queue Backup

**Alert Rule:**
```yaml
- alert: QueueBackup
  expr: rabbitmq_queue_messages > 1000
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Queue backup detected"
    description: "Queue {{ $labels.queue }} has {{ $value }} messages"
```

#### 6. Disk Space Low

**Alert Rule:**
```yaml
- alert: DiskSpaceLow
  expr: (node_filesystem_avail_bytes / node_filesystem_size_bytes) < 0.1
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Disk space low"
    description: "Disk {{ $labels.device }} has less than 10% free space"
```

#### 7. Memory Usage High

**Alert Rule:**
```yaml
- alert: HighMemoryUsage
  expr: (container_memory_usage_bytes / container_spec_memory_limit_bytes) > 0.9
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "High memory usage"
    description: "Container {{ $labels.container }} is using {{ $value | humanizePercentage }} of memory"
```

---

## Notification Channels

### Email Notification

**Grafana Configuration:**
```yaml
# monitoring/grafana/provisioning/notifiers/email.yml
notifiers:
  - name: email
    type: email
    uid: email
    org_id: 1
    is_default: true
    settings:
      addresses: admin@woragis.com,ops@woragis.com
      singleEmail: true
```

### Slack Notification

**Grafana Configuration:**
```yaml
notifiers:
  - name: slack
    type: slack
    uid: slack
    org_id: 1
    settings:
      url: https://hooks.slack.com/services/YOUR/WEBHOOK/URL
      username: Grafana
      channel: "#alerts"
      title: "{{ .GroupLabels.alertname }}"
      text: "{{ range .Alerts }}{{ .Annotations.description }}{{ end }}"
```

### PagerDuty Notification

**Grafana Configuration:**
```yaml
notifiers:
  - name: pagerduty
    type: pagerduty
    uid: pagerduty
    org_id: 1
    settings:
      integrationKey: YOUR_PAGERDUTY_KEY
      severity: critical
```

---

## Alert Routing

### Route by Severity

```yaml
# Grafana alert routing
routes:
  - receiver: 'critical-alerts'
    matchers:
      - severity = critical
    continue: false
  
  - receiver: 'warning-alerts'
    matchers:
      - severity = warning
    continue: false
  
  - receiver: 'default'
    matchers:
      - alertname = .*
```

### Route by Service

```yaml
routes:
  - receiver: 'database-team'
    matchers:
      - service = database
    continue: false
  
  - receiver: 'api-team'
    matchers:
      - service = app
    continue: false
```

---

## Incident Response

### Alert Response Procedures

**Critical Alert (Service Down):**
1. **Immediate:** Acknowledge alert
2. **0-5 min:** Check service status
3. **5-10 min:** Review logs and metrics
4. **10-15 min:** Identify root cause
5. **15-30 min:** Implement fix or workaround
6. **Post-incident:** Document and review

**Warning Alert (High Latency):**
1. **Immediate:** Acknowledge alert
2. **0-15 min:** Investigate cause
3. **15-30 min:** Determine if action needed
4. **30-60 min:** Implement fix if needed
5. **Monitor:** Verify resolution

### On-Call Rotation

**Schedule:**
- Primary: Week 1
- Secondary: Week 2
- Rotation: Weekly

**Responsibilities:**
- Respond to critical alerts within 15 minutes
- Escalate if unable to resolve
- Document incidents
- Hand off to next on-call

---

## Business Metrics

### Key Performance Indicators (KPIs)

**Track:**
- User registrations per day
- Active users
- API request volume
- Job processing rate
- Error rate by endpoint
- Response time by endpoint

**Grafana Dashboard:**
```yaml
# Create custom dashboard for business metrics
# Track:
# - User signups
# - API usage
# - Feature adoption
# - Revenue metrics (if applicable)
```

---

## Performance Monitoring

### Application Performance

**Metrics to Track:**
- Request rate (requests/second)
- Response time (p50, p95, p99)
- Error rate (errors/second)
- Throughput (bytes/second)

**Grafana Queries:**
```promql
# Request rate
rate(http_requests_total[5m])

# Response time (p95)
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Error rate
rate(http_requests_total{status=~"5.."}[5m])
```

### Database Performance

**Metrics to Track:**
- Query execution time
- Connection pool usage
- Slow queries
- Database size

**PostgreSQL Queries:**
```sql
-- Slow queries
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- Connection usage
SELECT count(*) FROM pg_stat_activity;
```

### Queue Performance

**Metrics to Track:**
- Queue depth
- Processing rate
- Message age
- Dead letter queue size

**RabbitMQ Metrics:**
```promql
# Queue depth
rabbitmq_queue_messages

# Processing rate
rate(rabbitmq_queue_messages_ready[5m])
```

---

## Log-Based Alerts

### Error Pattern Detection

**LogQL Alert:**
```yaml
- uid: error_pattern
  title: Error Pattern Detected
  condition: A
  data:
    - refId: A
      datasourceUid: loki
      model:
        expr: |
          sum(count_over_time(
            {job="docker", level="error"} 
            |~ "database|connection|timeout" [5m]
          )) > 5
  for: 3m
```

### Security Event Alerts

**LogQL Alert:**
```yaml
- uid: security_event
  title: Security Event Detected
  condition: A
  data:
    - refId: A
      datasourceUid: loki
      model:
        expr: |
          sum(count_over_time(
            {job="docker"} 
            |~ "unauthorized|forbidden|authentication failed" [5m]
          )) > 10
  for: 1m
```

---

## Alert Tuning

### Reduce False Positives

**Strategies:**
1. **Increase thresholds** - Don't alert on minor issues
2. **Add conditions** - Alert only when multiple conditions met
3. **Increase duration** - Alert only if issue persists
4. **Use different severity** - Not all issues are critical

### Example: Tuned Alert

```yaml
# Before: Too sensitive
- alert: HighErrorRate
  expr: rate(http_requests_total{status=~"5.."}[1m]) > 0.01
  for: 1m

# After: Tuned
- alert: HighErrorRate
  expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
  for: 5m
  labels:
    severity: critical
  annotations:
    summary: "High error rate: {{ $value }} errors/sec for 5+ minutes"
```

---

## Monitoring Best Practices

1. **Monitor What Matters**
   - User-facing metrics
   - Business metrics
   - System health metrics

2. **Set Appropriate Thresholds**
   - Based on historical data
   - Account for normal variations
   - Different thresholds per environment

3. **Use Multiple Alert Levels**
   - Critical: Immediate action needed
   - Warning: Monitor, investigate
   - Info: Logged, no action

4. **Document Alert Procedures**
   - What each alert means
   - How to investigate
   - How to resolve

5. **Review and Tune Regularly**
   - Weekly alert review
   - Adjust thresholds
   - Remove unused alerts

---

## Implementation Checklist

### Phase 1: Basic Alerts (Week 1)
- [ ] Service down alerts
- [ ] High error rate alerts
- [ ] Database connection alerts
- [ ] Configure email notifications
- [ ] Test alert delivery

### Phase 2: Advanced Alerts (Week 2)
- [ ] Performance alerts (latency, throughput)
- [ ] Resource alerts (CPU, memory, disk)
- [ ] Queue alerts
- [ ] Configure Slack/PagerDuty
- [ ] Set up alert routing

### Phase 3: Business Metrics (Week 3)
- [ ] Define business KPIs
- [ ] Create business metrics dashboard
- [ ] Set up business alerts
- [ ] Document procedures

---

## Related Documentation

- [Logging Aggregation Plan](../PLANNING/01-logging-aggregation-plan.md)
- [Metrics Implementation](../monitoring/METRICS_IMPLEMENTATION_SUMMARY.md)
- [Troubleshooting Guide](./troubleshooting.md)
- [Backup & Disaster Recovery](./backup-restore.md)

---

**Next Steps:**
1. Configure notification channels
2. Set up alert rules
3. Test alert delivery
4. Document incident response procedures
5. Train team on alerting

---

**Last Updated:** 2025-12-22
