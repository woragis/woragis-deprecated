# Production Readiness Implementation Guide

**Last Updated:** 2025-12-22  
**Purpose:** Step-by-step guide to implement production readiness features

---

## Quick Start

### 1. Install Pre-commit Hooks (5 minutes)

```bash
# Install pre-commit
pip install pre-commit

# Install hooks
pre-commit install

# Test hooks
pre-commit run --all-files
```

---

### 2. Set Up Secrets Management (15 minutes)

```bash
# Run setup script
./scripts/setup-sops.sh

# This will:
# - Install SOPS and age
# - Generate encryption keys
# - Update .sops.yaml
# - Add secrets-key.txt to .gitignore

# Encrypt your production secrets
./scripts/encrypt-secrets.sh .env.production

# Decrypt when needed (for deployment)
./scripts/decrypt-secrets.sh .env.production .env
```

---

### 3. Set Up Automated Backups (10 minutes)

```bash
# Make scripts executable
chmod +x scripts/*.sh

# Test backup
./scripts/backup-all.sh

# Set up automated backups (Linux/Mac)
./scripts/setup-cron-backups.sh

# Or manually add to crontab:
# 0 2 * * * /path/to/scripts/backup-all.sh
```

---

### 4. Test Security Middleware (5 minutes)

```bash
# Build server
cd server/app
go build ./cmd/server

# Start server
./server

# Test security headers
curl -I http://localhost:8080/healthz

# Should see security headers:
# X-Content-Type-Options: nosniff
# X-Frame-Options: DENY
# etc.
```

---

## Implementation Checklist

### ✅ Completed
- [x] Security middleware created and integrated
- [x] Validation utilities created
- [x] Backup scripts created
- [x] SOPS configuration created
- [x] Pre-commit hooks configured
- [x] CI/CD code quality workflow created

### ⏳ Next Steps

#### Week 1: Security & Backups
- [ ] Install and test pre-commit hooks
- [ ] Set up SOPS and encrypt secrets
- [ ] Test backup scripts
- [ ] Schedule automated backups
- [ ] Test security middleware

#### Week 2: SSL/TLS & Validation
- [ ] Obtain SSL certificates (Let's Encrypt)
- [ ] Configure HTTPS
- [ ] Add validation to API endpoints
- [ ] Test input validation
- [ ] Audit security measures

#### Week 3: Monitoring & Performance
- [ ] Configure alerting channels
- [ ] Set up alert rules
- [ ] Run performance tests
- [ ] Optimize based on results
- [ ] Document procedures

---

## Testing

### Test Security Middleware

```bash
# Start server
cd server/app
go run ./cmd/server/main.go

# Test rate limiting (should fail after 100 requests)
for i in {1..105}; do curl http://localhost:8080/healthz; done

# Test request size limit
curl -X POST http://localhost:8080/api/endpoint \
  -H "Content-Type: application/json" \
  -d "$(python -c "print('x' * 11 * 1024 * 1024)")"
# Should return 413 Request Entity Too Large
```

### Test Backup

```bash
# Run backup
./scripts/backup-all.sh

# Verify backup exists
ls -lh backups/woragis_backup_*.tar.gz

# Test restore (on test environment)
./scripts/restore-backup.sh backups/woragis_backup_YYYYMMDD_HHMMSS.tar.gz
```

### Test SOPS

```bash
# Setup
./scripts/setup-sops.sh

# Create test file
echo "SECRET_KEY=test123" > .env.test

# Encrypt
sops -e -i .env.test

# Verify encrypted
cat .env.test  # Should show encrypted content

# Decrypt
sops -d .env.test
```

---

## Troubleshooting

### Security Middleware Not Working

**Issue:** Headers not appearing in responses

**Solution:**
1. Check middleware is added before routes
2. Verify helmet package is installed: `go get github.com/gofiber/fiber/v2/middleware/helmet`
3. Check middleware order in main.go

### Backup Scripts Fail

**Issue:** Permission denied or container not found

**Solution:**
1. Make scripts executable: `chmod +x scripts/*.sh`
2. Check Docker containers are running: `docker-compose ps`
3. Verify container names match script defaults

### SOPS Encryption Fails

**Issue:** "age key not found" or encryption fails

**Solution:**
1. Run setup script: `./scripts/setup-sops.sh`
2. Check .sops.yaml has correct public key
3. Verify secrets-key.txt exists (never commit this!)

---

## Related Documentation

- [Secrets Management](./secrets-management.md)
- [SSL/TLS Configuration](./ssl-tls-configuration.md)
- [Backup & Disaster Recovery](../operations/backup-restore.md)
- [Security Headers](../operations/security.md) (when created)

---

**Last Updated:** 2025-12-22
