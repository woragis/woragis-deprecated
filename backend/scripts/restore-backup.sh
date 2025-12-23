#!/bin/bash
# Restore backup script for Woragis backend
# Usage: ./scripts/restore-backup.sh <backup-file.tar.gz>

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <backup-file.tar.gz>"
    echo "Example: $0 backups/woragis_backup_20241222_120000.tar.gz"
    exit 1
fi

BACKUP_FILE="$1"
RESTORE_DIR="./restore_temp"
CONTAINER_NAME="${CONTAINER_NAME:-woragis-database}"
DB_NAME="${DB_NAME:-woragis}"
DB_USER="${DB_USER:-postgres}"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    error "Backup file not found: $BACKUP_FILE"
    exit 1
fi

# Check if Docker container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    error "Database container '$CONTAINER_NAME' is not running"
    exit 1
fi

log "=========================================="
log "Starting restore process"
log "Backup file: $BACKUP_FILE"
log "=========================================="

# Extract backup
log "Extracting backup..."
mkdir -p "$RESTORE_DIR"
tar -xzf "$BACKUP_FILE" -C "$RESTORE_DIR"
EXTRACTED_DIR=$(ls -td "$RESTORE_DIR"/*/ | head -1)
log "✓ Backup extracted"

# Confirm restore
warn "WARNING: This will overwrite existing data!"
read -p "Are you sure you want to continue? (yes/no): " CONFIRM
if [ "$CONFIRM" != "yes" ]; then
    log "Restore cancelled"
    rm -rf "$RESTORE_DIR"
    exit 0
fi

# Stop application services (optional, recommended)
warn "Consider stopping application services before restore"
read -p "Stop application services? (yes/no): " STOP_SERVICES
if [ "$STOP_SERVICES" = "yes" ]; then
    log "Stopping application services..."
    docker-compose stop app ai-service creative-service docs-service 2>/dev/null || true
    log "✓ Services stopped"
fi

# Restore database
if [ -f "$EXTRACTED_DIR/database.sql.gz" ]; then
    log "Restoring database..."
    gunzip -c "$EXTRACTED_DIR/database.sql.gz" | docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" > /dev/null
    log "✓ Database restored"
else
    warn "Database backup not found in archive"
fi

# Restore Redis (if backup exists)
if [ -f "$EXTRACTED_DIR/redis.rdb" ]; then
    log "Restoring Redis..."
    docker-compose stop redis 2>/dev/null || true
    docker cp "$EXTRACTED_DIR/redis.rdb" woragis-redis:/data/dump.rdb 2>/dev/null || warn "Redis restore skipped (container not accessible)"
    docker-compose start redis 2>/dev/null || true
    log "✓ Redis restored"
fi

# Restore files
if [ -f "$EXTRACTED_DIR/files.tar.gz" ]; then
    log "Restoring files..."
    tar -xzf "$EXTRACTED_DIR/files.tar.gz" -C . 2>/dev/null || warn "Some files may not restore correctly"
    log "✓ Files restored"
fi

# Cleanup
log "Cleaning up temporary files..."
rm -rf "$RESTORE_DIR"
log "✓ Cleanup completed"

# Restart services if stopped
if [ "$STOP_SERVICES" = "yes" ]; then
    log "Restarting application services..."
    docker-compose start app ai-service creative-service docs-service 2>/dev/null || true
    log "✓ Services restarted"
fi

log "=========================================="
log "Restore completed successfully!"
log "=========================================="

# Verify restore
log "Verifying restore..."
if docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "SELECT COUNT(*) FROM information_schema.tables;" > /dev/null 2>&1; then
    log "✓ Database is accessible"
else
    error "✗ Database verification failed"
    exit 1
fi
