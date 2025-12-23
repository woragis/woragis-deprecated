#!/bin/bash
# Encrypt secrets files with SOPS
# Usage: ./scripts/encrypt-secrets.sh [file1] [file2] ...

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

if [ $# -eq 0 ]; then
    # Default: encrypt common secret files
    FILES=(
        ".env.production"
        ".env.staging"
    )
else
    FILES=("$@")
fi

# Check if sops is installed
if ! command -v sops &> /dev/null; then
    echo -e "${RED}SOPS is not installed. Run ./scripts/setup-sops.sh first${NC}"
    exit 1
fi

# Check if .sops.yaml exists
if [ ! -f ".sops.yaml" ]; then
    echo -e "${RED}.sops.yaml not found. Run ./scripts/setup-sops.sh first${NC}"
    exit 1
fi

for file in "${FILES[@]}"; do
    if [ ! -f "$file" ]; then
        echo -e "${YELLOW}File not found: $file (skipping)${NC}"
        continue
    fi
    
    # Check if already encrypted
    if grep -q "sops:" "$file" 2>/dev/null; then
        echo -e "${YELLOW}File already encrypted: $file (skipping)${NC}"
        continue
    fi
    
    echo -e "${GREEN}Encrypting: $file${NC}"
    if sops -e -i "$file"; then
        echo -e "${GREEN}✓ Encrypted: $file${NC}"
    else
        echo -e "${RED}✗ Failed to encrypt: $file${NC}"
        exit 1
    fi
done

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Encryption complete!${NC}"
echo -e "${GREEN}========================================${NC}"
