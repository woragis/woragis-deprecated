#!/bin/bash
# Test script for backup functionality
# Usage: ./scripts/test-backups.sh

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[TEST]${NC} $1"
}

error() {
    echo -e "${RED}[FAIL]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Check if Docker is running
log "Checking Docker..."
if ! docker ps > /dev/null 2>&1; then
    error "Docker is not running"
    exit 1
fi
log "✓ Docker is running"

# Check if database container exists
CONTAINER_NAME="${CONTAINER_NAME:-woragis-database}"
log "Checking database container: $CONTAINER_NAME"
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    warn "Database container '$CONTAINER_NAME' is not running"
    warn "Starting containers..."
    docker-compose up -d database 2>/dev/null || {
        error "Failed to start database container"
        exit 1
    }
    sleep 5
fi
log "✓ Database container is running"

# Test 1: Database Backup
log "Test 1: Testing database backup..."
if ./scripts/backup-database.sh; then
    BACKUP_FILE=$(ls -t backups/postgres/woragis_*.dump.gz 2>/dev/null | head -1)
    if [ -f "$BACKUP_FILE" ]; then
        BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
        log "✓ Database backup created: $BACKUP_FILE ($BACKUP_SIZE)"
    else
        error "✗ Backup file not found"
        exit 1
    fi
else
    error "✗ Database backup failed"
    exit 1
fi

# Test 2: Complete Backup
log "Test 2: Testing complete backup..."
if ./scripts/backup-all.sh; then
    BACKUP_FILE=$(ls -t backups/woragis_backup_*.tar.gz 2>/dev/null | head -1)
    if [ -f "$BACKUP_FILE" ]; then
        BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
        log "✓ Complete backup created: $BACKUP_FILE ($BACKUP_SIZE)"
    else
        error "✗ Backup file not found"
        exit 1
    fi
else
    error "✗ Complete backup failed"
    exit 1
fi

# Test 3: Backup Integrity
log "Test 3: Testing backup integrity..."
if [ -f "$BACKUP_FILE" ]; then
    if tar -tzf "$BACKUP_FILE" > /dev/null 2>&1; then
        log "✓ Backup archive is valid"
        
        # Check for required files
        MANIFEST=$(tar -tzf "$BACKUP_FILE" | grep manifest.json || echo "")
        DATABASE=$(tar -tzf "$BACKUP_FILE" | grep database.sql.gz || echo "")
        
        if [ -n "$MANIFEST" ]; then
            log "✓ Manifest file present"
        else
            warn "⚠ Manifest file not found"
        fi
        
        if [ -n "$DATABASE" ]; then
            log "✓ Database backup present"
        else
            warn "⚠ Database backup not found"
        fi
    else
        error "✗ Backup archive is corrupted"
        exit 1
    fi
else
    error "✗ Backup file not found"
    exit 1
fi

# Summary
echo ""
log "=========================================="
log "Backup Test Summary"
log "=========================================="
log "Database Backup: ✓ PASS"
log "Complete Backup: ✓ PASS"
log "Backup Integrity: ✓ PASS"
log "=========================================="
log ""
warn "Note: Restore testing should be done on a test environment only!"
log "To test restore: ./scripts/restore-backup.sh $BACKUP_FILE"
