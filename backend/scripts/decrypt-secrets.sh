#!/bin/bash
# Decrypt secrets files with SOPS
# Usage: ./scripts/decrypt-secrets.sh [encrypted-file] [output-file]

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

if [ $# -lt 1 ]; then
    echo "Usage: $0 <encrypted-file> [output-file]"
    echo "Example: $0 .env.production .env"
    exit 1
fi

ENCRYPTED_FILE="$1"
OUTPUT_FILE="${2:-${ENCRYPTED_FILE%.encrypted}}"

# Check if sops is installed
if ! command -v sops &> /dev/null; then
    echo -e "${RED}SOPS is not installed. Run ./scripts/setup-sops.sh first${NC}"
    exit 1
fi

# Check if encrypted file exists
if [ ! -f "$ENCRYPTED_FILE" ]; then
    echo -e "${RED}File not found: $ENCRYPTED_FILE${NC}"
    exit 1
fi

# Check if output file already exists
if [ -f "$OUTPUT_FILE" ] && [ "$OUTPUT_FILE" != "$ENCRYPTED_FILE" ]; then
    echo -e "${YELLOW}Output file already exists: $OUTPUT_FILE${NC}"
    read -p "Overwrite? (yes/no): " CONFIRM
    if [ "$CONFIRM" != "yes" ]; then
        echo "Cancelled"
        exit 0
    fi
fi

echo -e "${GREEN}Decrypting: $ENCRYPTED_FILE${NC}"
if sops -d "$ENCRYPTED_FILE" > "$OUTPUT_FILE"; then
    echo -e "${GREEN}✓ Decrypted to: $OUTPUT_FILE${NC}"
    echo -e "${YELLOW}⚠️  Remember to delete $OUTPUT_FILE after use!${NC}"
else
    echo -e "${RED}✗ Failed to decrypt: $ENCRYPTED_FILE${NC}"
    exit 1
fi
