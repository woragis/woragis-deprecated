# Secrets Management Guide

**Last Updated:** 2025-12-22  
**Purpose:** Guide for managing secrets and sensitive configuration in production

---

## Overview

This guide covers how to securely manage secrets (API keys, passwords, tokens) in the Woragis backend. Never commit secrets to version control.

---

## Current State

**Current Approach:** Environment variables via `.env` files (development only)

**Production Requirements:**
- Secrets must not be in code or config files
- Secrets must be encrypted at rest
- Secrets must be encrypted in transit
- Access to secrets must be audited
- Secrets must be rotatable

---

## Recommended Solutions

### Option 1: Docker Secrets (Recommended for Docker Swarm)

**Best for:** Docker Swarm deployments

**Setup:**
```bash
# Create a secret
echo "my-secret-password" | docker secret create db_password -

# Use in docker-compose.yml
services:
  app:
    secrets:
      - db_password
    environment:
      - DATABASE_PASSWORD_FILE=/run/secrets/db_password

secrets:
  db_password:
    external: true
```

**Pros:**
- Built into Docker
- Encrypted at rest
- Only accessible to services that need them
- Simple to use

**Cons:**
- Requires Docker Swarm (not available in Docker Compose)
- Limited to Docker environments

---

### Option 2: Environment Variables with Encryption (Current + Enhancement)

**Best for:** Docker Compose, Kubernetes, cloud platforms

**Current Setup:**
```bash
# .env file (never commit)
DATABASE_PASSWORD=encrypted_value_here
```

**Enhanced Approach:**
1. Encrypt secrets before storing
2. Decrypt at runtime
3. Use tools like `sops` or `ansible-vault`

**Example with SOPS:**
```bash
# Install sops
brew install sops  # macOS
# or download from https://github.com/mozilla/sops

# Encrypt a secret
sops -e -i .env.production

# Decrypt at runtime (in CI/CD or deployment script)
sops -d .env.production > .env
```

**Pros:**
- Works with Docker Compose
- Secrets can be version controlled (encrypted)
- Simple to implement
- Works with any deployment method

**Cons:**
- Requires encryption key management
- Need to decrypt before use

---

### Option 3: HashiCorp Vault (Advanced)

**Best for:** Large-scale deployments, multiple environments

**Setup:**
```bash
# Start Vault (development)
docker run -d --name vault -p 8200:8200 vault:latest

# Write a secret
vault kv put secret/woragis/database password=my-secret-password

# Read in application
vault kv get -field=password secret/woragis/database
```

**Pros:**
- Enterprise-grade security
- Dynamic secrets
- Audit logging
- Secret rotation
- Fine-grained access control

**Cons:**
- Complex setup
- Requires Vault infrastructure
- Learning curve

---

### Option 4: Cloud Secrets Managers

**Best for:** Cloud deployments (AWS, Azure, GCP)

**AWS Secrets Manager:**
```bash
# Store secret
aws secretsmanager create-secret \
  --name woragis/database \
  --secret-string '{"password":"my-secret-password"}'

# Retrieve in application
aws secretsmanager get-secret-value --secret-id woragis/database
```

**Azure Key Vault:**
```bash
# Store secret
az keyvault secret set --vault-name woragis-vault --name database-password --value my-secret-password

# Retrieve in application
az keyvault secret show --vault-name woragis-vault --name database-password
```

**Pros:**
- Managed service
- Automatic rotation
- Integration with cloud services
- Audit logging

**Cons:**
- Cloud-specific
- Cost per secret
- Vendor lock-in

---

## Recommended Implementation (Phase 1)

**For Current Setup:** Use **Option 2 (Environment Variables with SOPS)**

### Step 1: Install SOPS

```bash
# macOS
brew install sops

# Linux
wget https://github.com/mozilla/sops/releases/download/v3.8.0/sops-v3.8.0.linux
chmod +x sops-v3.8.0.linux
sudo mv sops-v3.8.0.linux /usr/local/bin/sops

# Windows
# Download from https://github.com/mozilla/sops/releases
```

### Step 2: Generate Encryption Key

```bash
# Generate age key (recommended)
age-keygen -o secrets-key.txt

# Or use PGP key
gpg --generate-key
```

### Step 3: Create Encrypted Secrets File

```bash
# Create .sops.yaml
cat > .sops.yaml <<EOF
creation_rules:
  - path_regex: \.env\.production$
    age: >-
      age1your-public-key-here
EOF

# Encrypt .env.production
sops -e -i .env.production
```

### Step 4: Update Deployment Scripts

```bash
# In deployment script
sops -d .env.production > .env
docker-compose up -d
```

---

## Secrets Audit Checklist

### Secrets to Identify and Secure:

- [ ] Database passwords (`DATABASE_URL`, `DB_PASSWORD`)
- [ ] Redis passwords (`REDIS_URL` with password)
- [ ] RabbitMQ credentials (`RABBITMQ_USER`, `RABBITMQ_PASSWORD`)
- [ ] API keys:
  - [ ] `OPENAI_API_KEY`
  - [ ] `ANTHROPIC_API_KEY`
  - [ ] `GOOGLE_TRANSLATE_API_KEY`
  - [ ] `DEEPL_API_KEY`
  - [ ] `XAI_API_KEY`
  - [ ] `MANUS_API_KEY`
  - [ ] `CIPHER_API_KEY`
- [ ] SMTP credentials (`SMTP_USERNAME`, `SMTP_PASSWORD`)
- [ ] JWT secrets (`JWT_SECRET`)
- [ ] Session secrets (`SESSION_SECRET`)
- [ ] Public API keys (`PUBLIC_API_KEY`)
- [ ] OAuth client secrets (if applicable)
- [ ] Encryption keys
- [ ] Service-to-service authentication tokens

---

## Best Practices

### 1. Never Commit Secrets

```bash
# Add to .gitignore
.env
.env.*
!.env.example
*.key
*.pem
secrets/
```

### 2. Use Different Secrets per Environment

- Development: `.env.development` (can be unencrypted, local only)
- Staging: `.env.staging` (encrypted)
- Production: `.env.production` (encrypted)

### 3. Rotate Secrets Regularly

- Database passwords: Every 90 days
- API keys: When compromised or quarterly
- JWT secrets: Every 180 days
- Service tokens: Every 90 days

### 4. Limit Access

- Only grant access to secrets to services/users that need them
- Use least privilege principle
- Audit secret access

### 5. Monitor Secret Usage

- Log secret access (without logging the secret itself)
- Alert on unusual access patterns
- Review access logs regularly

---

## Implementation Steps

### Phase 1: Audit (Week 1)
1. [ ] Identify all secrets in codebase
2. [ ] Document where secrets are used
3. [ ] Create secrets inventory
4. [ ] Remove any hardcoded secrets

### Phase 2: Setup (Week 1)
1. [ ] Choose secrets management solution
2. [ ] Install required tools (SOPS, Vault, etc.)
3. [ ] Generate encryption keys
4. [ ] Set up key management

### Phase 3: Migration (Week 2)
1. [ ] Encrypt existing secrets
2. [ ] Update services to read from secrets manager
3. [ ] Update deployment scripts
4. [ ] Test in staging environment

### Phase 4: Documentation (Week 2)
1. [ ] Document secrets management process
2. [ ] Create runbooks for secret rotation
3. [ ] Train team on new process
4. [ ] Update CI/CD pipelines

---

## Emergency Procedures

### If Secret is Compromised:

1. **Immediately rotate the secret:**
   ```bash
   # Generate new secret
   # Update in secrets manager
   # Deploy updated configuration
   ```

2. **Revoke old secret:**
   - Disable old API keys
   - Change passwords
   - Invalidate tokens

3. **Audit access:**
   - Review access logs
   - Identify unauthorized access
   - Take corrective action

4. **Notify team:**
   - Alert security team
   - Document incident
   - Update procedures if needed

---

## Tools Reference

### SOPS (Secrets OPerationS)
- **Website:** https://github.com/mozilla/sops
- **Use case:** Encrypt files for version control
- **Supports:** Age, PGP, AWS KMS, GCP KMS, Azure Key Vault

### HashiCorp Vault
- **Website:** https://www.vaultproject.io/
- **Use case:** Centralized secrets management
- **Features:** Dynamic secrets, secret rotation, audit logging

### Docker Secrets
- **Documentation:** https://docs.docker.com/engine/swarm/secrets/
- **Use case:** Docker Swarm deployments
- **Features:** Built-in encryption, service-scoped access

---

## Related Documentation

- [Configuration Reference](./configuration.md)
- [Docker Setup Guide](./docker-setup.md)
- [Deployment Procedures](./deployment-procedures.md) (when created)
- [Security Guide](../operations/security.md) (when created)

---

**Next Steps:**
1. Audit all secrets in the codebase
2. Choose secrets management solution
3. Set up encryption for production secrets
4. Update deployment processes

---

**Last Updated:** 2025-12-22
