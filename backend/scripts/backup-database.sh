#!/bin/bash
# Database backup script for Woragis backend
# Usage: ./scripts/backup-database.sh

set -e

BACKUP_DIR="${BACKUP_DIR:-./backups/postgres}"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="${DB_NAME:-woragis}"
DB_USER="${DB_USER:-postgres}"
CONTAINER_NAME="${CONTAINER_NAME:-woragis-database}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Check if Docker container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    error "Database container '$CONTAINER_NAME' is not running"
    exit 1
fi

# Create backup directory
mkdir -p "$BACKUP_DIR"
log "Backup directory: $BACKUP_DIR"

# Perform backup
log "Starting database backup..."
BACKUP_FILE="$BACKUP_DIR/woragis_$DATE.dump"

if docker exec "$CONTAINER_NAME" pg_dump -U "$DB_USER" -d "$DB_NAME" -F c > "$BACKUP_FILE"; then
    log "Database backup created: $BACKUP_FILE"
    
    # Compress backup
    log "Compressing backup..."
    gzip "$BACKUP_FILE"
    BACKUP_FILE="${BACKUP_FILE}.gz"
    log "Backup compressed: $BACKUP_FILE"
    
    # Get backup size
    BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    log "Backup size: $BACKUP_SIZE"
    
    # Remove old backups
    log "Cleaning up backups older than $RETENTION_DAYS days..."
    find "$BACKUP_DIR" -name "woragis_*.dump.gz" -mtime +$RETENTION_DAYS -delete
    DELETED_COUNT=$(find "$BACKUP_DIR" -name "woragis_*.dump.gz" -mtime +$RETENTION_DAYS 2>/dev/null | wc -l)
    if [ "$DELETED_COUNT" -gt 0 ]; then
        log "Deleted $DELETED_COUNT old backup(s)"
    fi
    
    log "Backup completed successfully: $BACKUP_FILE"
    exit 0
else
    error "Database backup failed"
    exit 1
fi
