# Alerting: When to Alert, When Not To

## Overview
Alerting strategies for when to send alerts vs when to just log issues.

## Key Points

### Alert Criteria
- **Critical**: Service down, database down
- **Warning**: High error rate, degraded performance
- **Info**: Deployment, configuration changes

### Don't Alert
- Expected errors (user errors)
- Transient issues (auto-recover)
- Low-severity issues
- Noisy alerts

## Alert Types

### Critical Alerts
- Service unavailable
- Database connection failure
- High error rate (>5%)
- Queue backup (>1000 messages)

### Warning Alerts
- Degraded performance
- Increased latency
- Moderate error rate (1-5%)
- Queue depth increasing

## Implementation

### Alert Rules
- Prometheus alerting rules
- Threshold-based
- Duration-based
- Rate-based

### Alert Channels
- Email
- Slack
- PagerDuty (critical)
- Dashboard

## Benefits
- Proactive issue detection
- Faster response
- Better reliability
- Team awareness

## Challenges
- Alert fatigue
- Threshold tuning
- False positives
- Alert routing

## Lessons Learned
- Alert on actionable issues
- Avoid noisy alerts
- Threshold tuning important
- Alert fatigue real

## Future Improvements
- Alert grouping
- Alert routing
- Alert analytics
- SLO-based alerting
