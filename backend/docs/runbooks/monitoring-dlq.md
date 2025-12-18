# Monitoring Dead Letter Queues (DLQ)

## Overview
Dead Letter Queues (DLQ) are used to capture messages that fail processing after all retry attempts. Monitoring DLQ is critical for identifying and resolving issues in the system.

## Prerequisites
- Access to RabbitMQ management UI or CLI
- Understanding of which queues use DLQ
- Access to application logs

## Understanding DLQ

### What is DLQ?
- Messages that fail processing are routed to DLQ
- DLQ exchange: `woragis.dlx`
- DLQ routing key: `{queue-name}.failed`
- Example: `translations.queue` → `translations.queue.failed`

### Why Messages End Up in DLQ
1. **Permanent Errors**: Invalid input, authentication failures
2. **Max Retries Exceeded**: Transient errors that persist after retries
3. **Processing Timeout**: Messages that take too long to process
4. **Worker Crashes**: Worker crashes during processing

## Monitoring DLQ

### 1. Check DLQ Size

#### Via RabbitMQ Management UI
1. Navigate to RabbitMQ Management UI (usually `http://localhost:15672`)
2. Go to **Exchanges** → Find `woragis.dlx`
3. Check **Bindings** to see bound queues
4. Go to **Queues** → Find `{queue-name}.failed` queues
5. Check **Messages** count

#### Via RabbitMQ CLI
```bash
# List all DLQ queues
rabbitmqctl list_queues name messages | grep "\.failed"

# Check specific DLQ size
rabbitmqctl list_queues name messages | grep "translations.queue.failed"
```

#### Via Prometheus Metrics (if available)
```promql
# DLQ size per queue
queue_dlq_size{queue="translations.queue"}
```

### 2. Monitor DLQ Growth

#### Alert Thresholds
- **Warning**: DLQ size > 100 messages
- **Critical**: DLQ size > 1000 messages
- **Alert**: DLQ growing continuously (rate > 10 messages/hour)

#### Check DLQ Growth Rate
```bash
# Check messages in DLQ over time
# Run this command periodically and compare
rabbitmqctl list_queues name messages | grep "\.failed"
```

### 3. Inspect Failed Messages

#### Via RabbitMQ Management UI
1. Go to **Queues** → Select `{queue-name}.failed`
2. Click **Get messages** (preview without ack)
3. Inspect message payload and headers
4. Check `x-death` header for retry information

#### Via RabbitMQ CLI
```bash
# Get messages from DLQ (without consuming)
rabbitmqadmin get queue={queue-name}.failed ackmode=ack_requeue_false count=10
```

## Analyzing DLQ Messages

### Common Patterns

#### 1. Invalid Input Errors
- **Pattern**: Messages with malformed data
- **Solution**: Fix message producer, add validation
- **Action**: Reprocess after fixing producer

#### 2. External API Failures
- **Pattern**: Translation API, AI service failures
- **Solution**: Check external API status, implement circuit breakers
- **Action**: Reprocess after API recovers

#### 3. Database Errors
- **Pattern**: Connection errors, constraint violations
- **Solution**: Check database health, fix data issues
- **Action**: Reprocess after database recovers

#### 4. Timeout Errors
- **Pattern**: Messages taking too long to process
- **Solution**: Optimize processing, increase timeout
- **Action**: Reprocess with increased timeout

## Handling DLQ Messages

### 1. Investigate Root Cause

#### Check Application Logs
```bash
# Search for error logs around the time messages failed
grep "ERROR" /var/log/{service}/app.log | grep "{message-id}"

# Check for specific error patterns
grep "translation.*failed" /var/log/translation-worker/app.log
```

#### Check Message Headers
- `x-death`: Retry count, original queue, failure reason
- `x-original-routing-key`: Original routing key
- `trace_id`: Trace ID for log correlation

### 2. Fix Root Cause

#### Common Fixes
- **Invalid Input**: Fix message producer, add validation
- **External API**: Check API status, implement retries/circuit breakers
- **Database**: Check database health, fix connection issues
- **Timeout**: Optimize processing, increase timeout

### 3. Reprocess Messages

#### Option 1: Manual Reprocessing (Recommended)
1. Fix root cause
2. Move messages from DLQ back to original queue
3. Monitor processing

```bash
# Move messages from DLQ to original queue
# Using RabbitMQ Management UI:
# 1. Go to DLQ queue
# 2. Get messages (with ack)
# 3. Publish to original queue with original routing key

# Using rabbitmqadmin:
rabbitmqadmin publish exchange=woragis.tasks routing_key=translations.process payload='{"id":"..."}'
```

#### Option 2: Automatic Reprocessing Script
```bash
#!/bin/bash
# reprocess-dlq.sh

QUEUE_NAME="translations.queue"
DLQ_NAME="${QUEUE_NAME}.failed"
ORIGINAL_ROUTING_KEY="translations.process"

# Get messages from DLQ and republish to original queue
rabbitmqadmin get queue=$DLQ_NAME ackmode=ack_requeue_false count=100 | \
  jq -r '.payload' | \
  while read payload; do
    rabbitmqadmin publish exchange=woragis.tasks routing_key=$ORIGINAL_ROUTING_KEY payload="$payload"
  done
```

#### Option 3: Delete Messages (If Not Needed)
```bash
# Delete messages from DLQ (use with caution!)
rabbitmqadmin delete queue={queue-name}.failed
```

## Prevention

### 1. Implement Proper Error Handling
- Distinguish between transient and permanent errors
- Retry transient errors with exponential backoff
- Don't retry permanent errors

### 2. Set Appropriate Retry Limits
- **Transient Errors**: 3-5 retries
- **Permanent Errors**: No retries (go directly to DLQ)

### 3. Monitor Error Rates
- Alert on high error rates
- Alert on DLQ growth
- Investigate errors before they accumulate

### 4. Implement Circuit Breakers
- Prevent cascading failures
- Fail fast when downstream services are down
- Reduce load on failing services

## Alerting

### Recommended Alerts

#### Critical Alerts
- DLQ size > 1000 messages
- DLQ growth rate > 100 messages/hour
- DLQ size > 500 for > 1 hour

#### Warning Alerts
- DLQ size > 100 messages
- DLQ growth rate > 10 messages/hour
- Any messages in DLQ (for monitoring)

### Alert Actions
1. **Immediate**: Investigate root cause
2. **Short-term**: Fix root cause, reprocess messages
3. **Long-term**: Implement prevention measures

## Troubleshooting

### DLQ Growing Rapidly

#### Symptoms
- DLQ size increasing quickly
- High error rate in logs
- Service degradation

#### Steps
1. Check application logs for error patterns
2. Check external API status
3. Check database health
4. Check worker health (are workers running?)
5. Check message format (are messages valid?)

### DLQ Not Clearing

#### Symptoms
- Messages in DLQ not being reprocessed
- Root cause fixed but messages still in DLQ

#### Steps
1. Verify root cause is actually fixed
2. Manually reprocess messages (see above)
3. Check if workers are consuming from original queue
4. Verify message format is correct

### Messages Disappearing from DLQ

#### Symptoms
- DLQ size decreasing unexpectedly
- Messages not in DLQ or original queue

#### Steps
1. Check if messages were manually deleted
2. Check if messages expired (if TTL set)
3. Check RabbitMQ logs for errors
4. Verify DLQ configuration

## Best Practices

### 1. Regular Monitoring
- Check DLQ size daily
- Monitor DLQ growth trends
- Investigate any messages in DLQ

### 2. Document Common Issues
- Keep a log of common DLQ issues
- Document solutions for future reference
- Share knowledge with team

### 3. Automate Where Possible
- Automate DLQ monitoring
- Automate alerting
- Automate reprocessing (if safe)

### 4. Review and Improve
- Review DLQ patterns regularly
- Identify root causes
- Implement prevention measures

## Related Documentation
- [Architecture Decision Records](../adr/) - DLQ implementation details
- [Troubleshooting Guide](./troubleshooting.md) - General troubleshooting
- [Deploying Services and Workers](./deploying-services.md) - Deployment procedures
