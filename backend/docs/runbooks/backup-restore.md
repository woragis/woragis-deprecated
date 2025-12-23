# Backup & Disaster Recovery Guide

**Last Updated:** 2025-12-22  
**Purpose:** Procedures for backing up and restoring system data

---

## Overview

This guide covers backup and disaster recovery procedures for the Woragis backend services.

---

## What to Backup

### Critical Data

1. **PostgreSQL Database**
   - User data
   - Application data
   - Configuration data

2. **Redis Data** (if persistence enabled)
   - Cache data
   - Session data
   - Temporary data

3. **File Storage**
   - Generated resumes (resume-worker output)
   - Uploaded files
   - Configuration files

4. **Secrets & Configuration**
   - Environment variables
   - SSL certificates
   - API keys (encrypted)

---

## Backup Strategy

### Database Backups

**Frequency:**
- **Full Backup:** Daily at 2 AM
- **Incremental Backup:** Every 6 hours
- **Retention:** 30 days for daily, 7 days for incremental

**PostgreSQL Backup:**
```bash
#!/bin/bash
# backup-database.sh

BACKUP_DIR="/backups/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="woragis"
DB_USER="postgres"

# Create backup directory
mkdir -p $BACKUP_DIR

# Full backup
pg_dump -U $DB_USER -d $DB_NAME -F c -f $BACKUP_DIR/woragis_$DATE.dump

# Compress
gzip $BACKUP_DIR/woragis_$DATE.dump

# Remove backups older than 30 days
find $BACKUP_DIR -name "*.dump.gz" -mtime +30 -delete

echo "Backup completed: woragis_$DATE.dump.gz"
```

**Docker Backup:**
```bash
# Backup from Docker container
docker exec woragis-database pg_dump -U postgres woragis | gzip > backup_$(date +%Y%m%d).sql.gz
```

### Redis Backups

**If Persistence Enabled:**
```bash
# Copy RDB file
docker cp woragis-redis:/data/dump.rdb ./backups/redis/dump_$(date +%Y%m%d).rdb
```

**Note:** Redis is primarily a cache. If data loss is acceptable, backups may not be necessary.

### File Backups

```bash
#!/bin/bash
# backup-files.sh

BACKUP_DIR="/backups/files"
DATE=$(date +%Y%m%d_%H%M%S)

# Backup resume-worker output
tar -czf $BACKUP_DIR/resumes_$DATE.tar.gz resume-worker/output/

# Backup uploaded files (if any)
tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz uploads/

# Remove backups older than 30 days
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete
```

---

## Automated Backup Script

```bash
#!/bin/bash
# automated-backup.sh

set -e

BACKUP_ROOT="/backups"
DATE=$(date +%Y%m%d_%H%M%S)
LOG_FILE="/var/log/backup.log"

log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" | tee -a $LOG_FILE
}

log "Starting backup process"

# Create backup directory
mkdir -p $BACKUP_ROOT/$DATE

# Backup database
log "Backing up database..."
docker exec woragis-database pg_dump -U postgres woragis | gzip > $BACKUP_ROOT/$DATE/database.sql.gz
log "Database backup completed"

# Backup Redis (if needed)
log "Backing up Redis..."
docker cp woragis-redis:/data/dump.rdb $BACKUP_ROOT/$DATE/redis.rdb 2>/dev/null || log "Redis backup skipped (no persistence)"
log "Redis backup completed"

# Backup files
log "Backing up files..."
tar -czf $BACKUP_ROOT/$DATE/files.tar.gz resume-worker/output/ 2>/dev/null || log "File backup skipped (no files)"
log "File backup completed"

# Backup configuration
log "Backing up configuration..."
tar -czf $BACKUP_ROOT/$DATE/config.tar.gz docker-compose.yml .env.example monitoring/ 2>/dev/null
log "Configuration backup completed"

# Compress entire backup
log "Compressing backup..."
tar -czf $BACKUP_ROOT/woragis_backup_$DATE.tar.gz -C $BACKUP_ROOT $DATE
rm -rf $BACKUP_ROOT/$DATE
log "Backup compression completed"

# Upload to remote storage (optional)
# aws s3 cp $BACKUP_ROOT/woragis_backup_$DATE.tar.gz s3://woragis-backups/

log "Backup process completed: woragis_backup_$DATE.tar.gz"

# Remove backups older than 30 days
find $BACKUP_ROOT -name "woragis_backup_*.tar.gz" -mtime +30 -delete
log "Old backups cleaned up"
```

**Schedule with Cron:**
```bash
# Add to crontab
0 2 * * * /path/to/automated-backup.sh
```

---

## Restore Procedures

### Database Restore

**From Dump File:**
```bash
# Stop application
docker-compose stop app

# Restore database
gunzip < backup_20241222.sql.gz | docker exec -i woragis-database psql -U postgres -d woragis

# Verify restore
docker exec woragis-database psql -U postgres -d woragis -c "SELECT COUNT(*) FROM users;"

# Start application
docker-compose start app
```

**Point-in-Time Recovery:**
```bash
# Requires WAL archiving enabled
# Restore base backup
pg_basebackup -D /var/lib/postgresql/data -Ft -z -P

# Replay WAL logs to specific time
pg_recovery -D /var/lib/postgresql/data --target-time="2024-12-22 10:00:00"
```

### Redis Restore

```bash
# Stop Redis
docker-compose stop redis

# Copy RDB file
docker cp ./backups/redis/dump_20241222.rdb woragis-redis:/data/dump.rdb

# Start Redis
docker-compose start redis
```

### File Restore

```bash
# Extract backup
tar -xzf backups/files/resumes_20241222.tar.gz -C resume-worker/output/

# Verify files
ls -lh resume-worker/output/
```

---

## Disaster Recovery Plan

### RTO (Recovery Time Objective): 4 hours
### RPO (Recovery Point Objective): 1 hour

### Recovery Scenarios

#### Scenario 1: Database Corruption

**Steps:**
1. Stop all services
2. Identify last good backup
3. Restore database from backup
4. Verify data integrity
5. Restart services
6. Monitor for issues

**Time Estimate:** 1-2 hours

#### Scenario 2: Complete System Failure

**Steps:**
1. Provision new infrastructure
2. Restore database from backup
3. Restore configuration files
4. Restore file storage
5. Deploy application code
6. Verify all services
7. Test critical functionality

**Time Estimate:** 2-4 hours

#### Scenario 3: Data Center Outage

**Steps:**
1. Failover to secondary region (if available)
2. Restore from latest backup
3. Update DNS/routing
4. Verify services
5. Monitor for issues

**Time Estimate:** 2-4 hours

---

## Backup Testing

### Regular Testing Schedule

- **Weekly:** Test database restore on staging
- **Monthly:** Full disaster recovery drill
- **Quarterly:** Test backup restoration procedures

### Test Procedure

```bash
# 1. Create test environment
docker-compose -f docker-compose.test.yml up -d

# 2. Restore backup
./restore-backup.sh backup_20241222.tar.gz

# 3. Verify data
docker exec test-database psql -U postgres -d woragis -c "SELECT COUNT(*) FROM users;"

# 4. Test application
curl http://localhost:8080/healthz

# 5. Document results
echo "Backup test completed successfully" >> backup-test.log
```

---

## Backup Storage

### Local Storage
- **Location:** `/backups/`
- **Retention:** 30 days
- **Encryption:** Optional (recommended for sensitive data)

### Remote Storage Options

**AWS S3:**
```bash
aws s3 cp backup.tar.gz s3://woragis-backups/
aws s3 ls s3://woragis-backups/
```

**Azure Blob Storage:**
```bash
az storage blob upload --container-name backups --name backup.tar.gz --file backup.tar.gz
```

**Google Cloud Storage:**
```bash
gsutil cp backup.tar.gz gs://woragis-backups/
```

### Backup Encryption

```bash
# Encrypt backup before upload
gpg --symmetric --cipher-algo AES256 backup.tar.gz
# Creates backup.tar.gz.gpg

# Decrypt
gpg --decrypt backup.tar.gz.gpg > backup.tar.gz
```

---

## Monitoring & Alerts

### Backup Monitoring

**Check Backup Success:**
```bash
# Check last backup time
ls -lh /backups/ | tail -1

# Check backup size
du -sh /backups/

# Verify backup integrity
gzip -t /backups/woragis_backup_20241222.tar.gz
```

**Alert on Backup Failure:**
- Monitor backup script exit codes
- Alert if backup older than 25 hours
- Alert if backup size is suspiciously small
- Alert on backup script errors

---

## Backup Retention Policy

| Backup Type | Retention | Location |
|------------|-----------|----------|
| Daily Full | 30 days | Local + Remote |
| Incremental | 7 days | Local |
| Weekly Full | 12 weeks | Remote |
| Monthly Full | 12 months | Remote |
| Yearly Full | 7 years | Remote (compliance) |

---

## Implementation Checklist

### Phase 1: Setup (Week 1)
- [ ] Create backup scripts
- [ ] Set up backup storage
- [ ] Configure automated backups
- [ ] Test backup process
- [ ] Document procedures

### Phase 2: Testing (Week 2)
- [ ] Test database restore
- [ ] Test file restore
- [ ] Test disaster recovery
- [ ] Document test results
- [ ] Refine procedures

### Phase 3: Monitoring (Week 2)
- [ ] Set up backup monitoring
- [ ] Configure alerts
- [ ] Create dashboards
- [ ] Schedule regular tests

---

## Related Documentation

- [Deployment Procedures](../deployment/deployment-procedures.md) (when created)
- [Disaster Recovery Overview](../architecture/DISASTER_RECOVERY_OVERVIEW.md)
- [Monitoring Guide](./monitoring.md) (when created)

---

**Next Steps:**
1. Create backup scripts
2. Set up automated backups
3. Test restore procedures
4. Document disaster recovery plan
5. Schedule regular backup tests

---

**Last Updated:** 2025-12-22
