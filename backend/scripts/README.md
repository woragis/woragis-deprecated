# Backup Scripts

**Purpose:** Automated backup and restore scripts for Woragis backend

---

## Scripts

### `backup-database.sh`

Backs up PostgreSQL database only.

**Usage:**
```bash
./scripts/backup-database.sh
```

**Environment Variables:**
- `BACKUP_DIR` - Backup directory (default: `./backups/postgres`)
- `DB_NAME` - Database name (default: `woragis`)
- `DB_USER` - Database user (default: `postgres`)
- `CONTAINER_NAME` - Container name (default: `woragis-database`)
- `RETENTION_DAYS` - Days to keep backups (default: `30`)

**Output:**
- Creates compressed backup: `woragis_YYYYMMDD_HHMMSS.dump.gz`
- Removes backups older than retention period

---

### `backup-all.sh`

Backs up database, Redis, files, and configuration.

**Usage:**
```bash
./scripts/backup-all.sh
```

**Environment Variables:**
- `BACKUP_ROOT` - Root backup directory (default: `./backups`)
- `LOG_FILE` - Log file path (default: `./backups/backup.log`)
- `RETENTION_DAYS` - Days to keep backups (default: `30`)

**Output:**
- Creates compressed archive: `woragis_backup_YYYYMMDD_HHMMSS.tar.gz`
- Includes: database, Redis, files, configuration
- Creates manifest.json with backup metadata

---

### `restore-backup.sh`

Restores from a backup archive.

**Usage:**
```bash
./scripts/restore-backup.sh backups/woragis_backup_20241222_120000.tar.gz
```

**Features:**
- Extracts backup archive
- Restores database
- Restores Redis (if backup exists)
- Restores files
- Verifies restore

**⚠️ Warning:** This will overwrite existing data!

---

## Setup

### Make Scripts Executable

```bash
chmod +x scripts/*.sh
```

### Schedule Automated Backups

**Add to crontab:**
```bash
# Daily backup at 2 AM
0 2 * * * /path/to/woragis/backend/scripts/backup-all.sh

# Or weekly
0 2 * * 0 /path/to/woragis/backend/scripts/backup-all.sh
```

---

## Backup Storage

### Local Storage

Backups are stored in `./backups/` by default.

### Remote Storage (Optional)

**Upload to S3:**
```bash
aws s3 cp backups/woragis_backup_*.tar.gz s3://woragis-backups/
```

**Upload to Azure:**
```bash
az storage blob upload --container-name backups --file backups/woragis_backup_*.tar.gz
```

---

## Testing

### Test Backup

```bash
# Run backup
./scripts/backup-all.sh

# Verify backup exists
ls -lh backups/woragis_backup_*.tar.gz

# Test restore (on test environment)
./scripts/restore-backup.sh backups/woragis_backup_YYYYMMDD_HHMMSS.tar.gz
```

---

## Troubleshooting

### Backup Fails

**Check:**
1. Docker containers are running
2. Sufficient disk space
3. Write permissions on backup directory
4. Database container name is correct

### Restore Fails

**Check:**
1. Backup file is not corrupted
2. Database container is running
3. Sufficient disk space
4. Database is not in use

---

## Related Documentation

- [Backup & Disaster Recovery Guide](../docs/operations/backup-restore.md)
- [Docker Setup Guide](../docs/deployment/docker-setup.md)

---

**Last Updated:** 2025-12-22
