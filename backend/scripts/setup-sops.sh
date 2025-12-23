#!/bin/bash
# Setup SOPS for secrets management
# Usage: ./scripts/setup-sops.sh

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}Setting up SOPS for secrets management...${NC}"

# Check if age-keygen is installed
if ! command -v age-keygen &> /dev/null; then
    echo -e "${YELLOW}age-keygen not found. Installing age...${NC}"
    
    # Detect OS
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        wget https://github.com/FiloSottile/age/releases/latest/download/age-v1.1.1-linux-amd64.tar.gz
        tar -xzf age-v1.1.1-linux-amd64.tar.gz
        sudo mv age/age* /usr/local/bin/
        rm -rf age age-v1.1.1-linux-amd64.tar.gz
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        brew install age
    else
        echo "Please install age manually from https://github.com/FiloSottile/age"
        exit 1
    fi
fi

# Check if sops is installed
if ! command -v sops &> /dev/null; then
    echo -e "${YELLOW}sops not found. Installing SOPS...${NC}"
    
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        wget https://github.com/mozilla/sops/releases/latest/download/sops-v3.8.1.linux
        chmod +x sops-v3.8.1.linux
        sudo mv sops-v3.8.1.linux /usr/local/bin/sops
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        brew install sops
    else
        echo "Please install sops manually from https://github.com/mozilla/sops"
        exit 1
    fi
fi

# Generate age key if it doesn't exist
if [ ! -f "secrets-key.txt" ]; then
    echo -e "${GREEN}Generating age key pair...${NC}"
    age-keygen -o secrets-key.txt
    echo -e "${GREEN}Age key generated: secrets-key.txt${NC}"
    echo -e "${YELLOW}⚠️  IMPORTANT: Keep secrets-key.txt secure and never commit it!${NC}"
    echo -e "${YELLOW}⚠️  Add secrets-key.txt to .gitignore${NC}"
else
    echo -e "${GREEN}Age key already exists: secrets-key.txt${NC}"
fi

# Extract public key
PUBLIC_KEY=$(grep "public key:" secrets-key.txt | cut -d: -f2 | tr -d ' ')
echo -e "${GREEN}Public key: ${PUBLIC_KEY}${NC}"

# Update .sops.yaml with public key
if [ -f ".sops.yaml" ]; then
    echo -e "${GREEN}Updating .sops.yaml with public key...${NC}"
    # Replace placeholder with actual public key
    sed -i.bak "s/age1your-public-key-here-replace-me/$PUBLIC_KEY/g" .sops.yaml
    rm -f .sops.yaml.bak
    echo -e "${GREEN}.sops.yaml updated${NC}"
else
    echo -e "${YELLOW}.sops.yaml not found. Please create it manually.${NC}"
fi

# Add to .gitignore
if ! grep -q "secrets-key.txt" .gitignore 2>/dev/null; then
    echo "" >> .gitignore
    echo "# SOPS secrets key" >> .gitignore
    echo "secrets-key.txt" >> .gitignore
    echo -e "${GREEN}Added secrets-key.txt to .gitignore${NC}"
fi

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}SOPS setup complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Next steps:"
echo "1. Encrypt your .env.production file:"
echo "   sops -e -i .env.production"
echo ""
echo "2. Decrypt when needed:"
echo "   sops -d .env.production > .env"
echo ""
echo "3. Keep secrets-key.txt secure!"
