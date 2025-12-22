# Log Retention and Archival Policy

**Version:** 1.0  
**Last Updated:** 2025-12-22  
**Status:** Active

## Overview

This document defines the log retention and archival policies for the Woragis logging aggregation system.

## Retention Periods

### Current Configuration

- **Hot Storage (Loki)**: 30 days (720 hours)
- **Warm Storage**: Not configured (future enhancement)
- **Cold Storage**: Not configured (future enhancement)

### Retention by Log Level

| Log Level | Retention Period | Rationale |
|-----------|------------------|-----------|
| ERROR | 90 days | Critical for debugging and compliance |
| WARN | 30 days | Useful for trend analysis |
| INFO | 30 days | Standard operational logs |
| DEBUG | 7 days | Development/debugging only |

### Retention by Service

All services follow the same retention policy. Future enhancements may include:
- Extended retention for critical services (app, database)
- Shorter retention for high-volume services (workers)

## Configuration

### Loki Configuration

Retention is configured in `monitoring/loki-config.yml`:

```yaml
limits_config:
  retention_period: 720h  # 30 days
```

### Updating Retention

To change retention period:

1. Edit `monitoring/loki-config.yml`
2. Update `retention_period` value
3. Restart Loki service:
   ```bash
   docker-compose restart loki
   ```

## Archival Strategy

### Current State

- **Archival**: Not implemented
- **Backup**: Not implemented
- **Restoration**: Not implemented

### Future Enhancements

#### Phase 1: Basic Archival (Planned)
- Export logs older than retention period to compressed files
- Store in object storage (S3, GCS, Azure Blob)
- Monthly archival jobs

#### Phase 2: Tiered Storage (Planned)
- Hot: 7 days in Loki
- Warm: 30 days in object storage (frequently accessed)
- Cold: 90+ days in object storage (rarely accessed)

#### Phase 3: Automated Archival (Planned)
- Automated daily archival jobs
- Compression (gzip)
- Encryption at rest
- Retention policies per service/level

## Backup Strategy

### Current State

- **Backup**: Not implemented
- **Disaster Recovery**: Not implemented

### Recommended Backup Strategy

1. **Daily Backups**
   - Export critical logs (ERROR level) daily
   - Store in separate backup location
   - Retain for 1 year

2. **Weekly Full Backups**
   - Export all logs weekly
   - Compressed and encrypted
   - Retain for 3 months

3. **Monthly Archives**
   - Export all logs monthly
   - Long-term storage
   - Retain for compliance period (varies by requirement)

## Disk Space Management

### Current Limits

- **Loki Data Volume**: No explicit limit (uses Docker volume)
- **Log Size Limits**: 256KB per log line (configured in Loki)

### Monitoring

Monitor disk usage:
```bash
docker system df
docker volume inspect woragis_loki-data
```

### Cleanup

To manually clean old logs:
```bash
# Restart Loki to apply retention
docker-compose restart loki

# Or manually delete old data (use with caution)
docker exec woragis-loki loki -config.file=/etc/loki/local-config.yaml -target=compactor
```

## Compliance Requirements

### Audit Logs

- **Retention**: 1 year minimum
- **Access Logs**: 90 days minimum
- **Security Events**: 1 year minimum

### Data Privacy

- **PII Removal**: Logs should not contain PII
- **GDPR Compliance**: Logs containing user data must be removable
- **Data Minimization**: Only log necessary information

## Cost Considerations

### Storage Costs

- **30 days retention**: ~10-50GB (estimated)
- **90 days retention**: ~30-150GB (estimated)
- **1 year retention**: ~120-600GB (estimated)

### Optimization Strategies

1. **Reduce log volume**
   - Use appropriate log levels
   - Avoid verbose debug logs in production
   - Filter out health check logs

2. **Compression**
   - Enable compression in Loki
   - Compress archived logs

3. **Selective Retention**
   - Longer retention for errors only
   - Shorter retention for info/debug logs

## Implementation Checklist

- [x] Configure retention period in Loki
- [ ] Set up automated archival jobs
- [ ] Configure object storage for archives
- [ ] Implement backup strategy
- [ ] Set up disk space monitoring
- [ ] Create restoration procedures
- [ ] Document compliance requirements
- [ ] Test archival and restoration

## Future Enhancements

1. **Automated Archival**
   - Daily export of logs older than retention
   - Compression and encryption
   - Upload to object storage

2. **Tiered Storage**
   - Hot/Warm/Cold storage tiers
   - Automatic tier migration
   - Cost optimization

3. **Backup Automation**
   - Scheduled backups
   - Backup verification
   - Automated restoration testing

4. **Compliance Features**
   - PII detection and redaction
   - Automated data deletion
   - Audit trail for log access

## References

- [Loki Retention Documentation](https://grafana.com/docs/loki/latest/operations/storage/retention/)
- [Loki Compactor](https://grafana.com/docs/loki/latest/operations/storage/compactor/)
