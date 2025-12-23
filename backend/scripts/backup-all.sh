#!/bin/bash
# Complete backup script for Woragis backend
# Backs up: database, Redis, files, configuration
# Usage: ./scripts/backup-all.sh

set -e

BACKUP_ROOT="${BACKUP_ROOT:-./backups}"
DATE=$(date +%Y%m%d_%H%M%S)
LOG_FILE="${LOG_FILE:-./backups/backup.log}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE" >&2
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1" | tee -a "$LOG_FILE"
}

# Create backup directory
BACKUP_DIR="$BACKUP_ROOT/$DATE"
mkdir -p "$BACKUP_DIR"
mkdir -p "$(dirname "$LOG_FILE")"

log "=========================================="
log "Starting Woragis backup process"
log "Backup directory: $BACKUP_DIR"
log "=========================================="

# Backup database
log "Backing up PostgreSQL database..."
if docker exec woragis-database pg_dump -U postgres woragis 2>/dev/null | gzip > "$BACKUP_DIR/database.sql.gz"; then
    DB_SIZE=$(du -h "$BACKUP_DIR/database.sql.gz" | cut -f1)
    log "✓ Database backup completed ($DB_SIZE)"
else
    error "✗ Database backup failed"
    exit 1
fi

# Backup Redis (if persistence enabled)
log "Backing up Redis..."
if docker cp woragis-redis:/data/dump.rdb "$BACKUP_DIR/redis.rdb" 2>/dev/null; then
    log "✓ Redis backup completed"
else
    warn "Redis backup skipped (no persistence or container not running)"
fi

# Backup files
log "Backing up files..."
if [ -d "resume-worker/output" ] && [ "$(ls -A resume-worker/output 2>/dev/null)" ]; then
    tar -czf "$BACKUP_DIR/files.tar.gz" resume-worker/output/ resume-worker/results/ 2>/dev/null || warn "Some files may not exist"
    FILES_SIZE=$(du -h "$BACKUP_DIR/files.tar.gz" 2>/dev/null | cut -f1 || echo "0")
    log "✓ Files backup completed ($FILES_SIZE)"
else
    warn "No files to backup"
fi

# Backup configuration
log "Backing up configuration..."
tar -czf "$BACKUP_DIR/config.tar.gz" \
    docker-compose.yml \
    .env.example \
    monitoring/ \
    docs/PLANNING/ \
    2>/dev/null || warn "Some config files may not exist"
log "✓ Configuration backup completed"

# Create backup manifest
cat > "$BACKUP_DIR/manifest.json" <<EOF
{
  "backup_date": "$(date -Iseconds)",
  "backup_version": "1.0",
  "components": {
    "database": "$(test -f "$BACKUP_DIR/database.sql.gz" && echo "present" || echo "missing")",
    "redis": "$(test -f "$BACKUP_DIR/redis.rdb" && echo "present" || echo "missing")",
    "files": "$(test -f "$BACKUP_DIR/files.tar.gz" && echo "present" || echo "missing")",
    "config": "$(test -f "$BACKUP_DIR/config.tar.gz" && echo "present" || echo "missing")"
  },
  "backup_size": "$(du -sh "$BACKUP_DIR" | cut -f1)"
}
EOF

# Compress entire backup
log "Compressing backup archive..."
cd "$BACKUP_ROOT"
tar -czf "woragis_backup_$DATE.tar.gz" "$DATE"
BACKUP_SIZE=$(du -h "woragis_backup_$DATE.tar.gz" | cut -f1)
log "✓ Backup archive created: woragis_backup_$DATE.tar.gz ($BACKUP_SIZE)"

# Remove uncompressed directory
rm -rf "$DATE"

# Clean up old backups
log "Cleaning up backups older than $RETENTION_DAYS days..."
find "$BACKUP_ROOT" -name "woragis_backup_*.tar.gz" -mtime +$RETENTION_DAYS -delete
DELETED=$(find "$BACKUP_ROOT" -name "woragis_backup_*.tar.gz" -mtime +$RETENTION_DAYS 2>/dev/null | wc -l)
if [ "$DELETED" -gt 0 ]; then
    log "Deleted $DELETED old backup(s)"
fi

log "=========================================="
log "Backup completed successfully!"
log "Backup file: $BACKUP_ROOT/woragis_backup_$DATE.tar.gz"
log "=========================================="
