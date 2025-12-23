#!/bin/bash
# Setup automated backups with cron
# Usage: ./scripts/setup-cron-backups.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_SCRIPT="$SCRIPT_DIR/backup-all.sh"
CRON_SCHEDULE="${CRON_SCHEDULE:-0 2 * * *}"  # Default: 2 AM daily

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}Setting up automated backups...${NC}"

# Check if backup script exists
if [ ! -f "$BACKUP_SCRIPT" ]; then
    echo "Error: Backup script not found: $BACKUP_SCRIPT"
    exit 1
fi

# Make script executable
chmod +x "$BACKUP_SCRIPT"

# Create cron entry
CRON_ENTRY="$CRON_SCHEDULE $BACKUP_SCRIPT >> /var/log/woragis-backup.log 2>&1"

# Check if cron entry already exists
if crontab -l 2>/dev/null | grep -q "$BACKUP_SCRIPT"; then
    echo -e "${YELLOW}Cron entry already exists${NC}"
    echo "Current cron entries:"
    crontab -l | grep "$BACKUP_SCRIPT"
else
    # Add to crontab
    (crontab -l 2>/dev/null; echo "$CRON_ENTRY") | crontab -
    echo -e "${GREEN}Cron entry added successfully${NC}"
    echo "Schedule: $CRON_SCHEDULE"
    echo "Script: $BACKUP_SCRIPT"
fi

echo ""
echo "To view cron entries:"
echo "  crontab -l"
echo ""
echo "To remove cron entry:"
echo "  crontab -e"
echo "  (then delete the line with backup-all.sh)"
