# Testing Scripts

**Purpose:** Automated testing scripts for production readiness features

---

## Available Test Scripts

### 1. `test-security-middleware.sh`

Tests security middleware functionality.

**Usage:**
```bash
# Start server first
cd server/app
go run ./cmd/server/main.go

# In another terminal, run tests
./scripts/test-security-middleware.sh

# Or test against different URL
./scripts/test-security-middleware.sh http://localhost:3000
```

**Tests:**
- ✅ Security headers presence
- ✅ Rate limiting (100 req/min)
- ✅ Request size limits (10MB)
- ✅ Input sanitization

**Expected Output:**
```
[TEST] Checking if server is running...
[TEST] Server is running
[TEST] Test 1: Checking security headers...
  ✓ X-Content-Type-Options present
  ✓ X-Frame-Options present
  ✓ X-XSS-Protection present
✓ All security headers present
[TEST] Test 2: Testing rate limiting...
  Making 105 requests...
  Successful requests: 100
  Rate limited requests: 5
✓ Rate limiting is working
...
```

---

### 2. `test-backups.sh`

Tests backup functionality.

**Usage:**
```bash
# Make sure Docker containers are running
docker-compose up -d

# Run tests
./scripts/test-backups.sh
```

**Tests:**
- ✅ Database backup creation
- ✅ Complete backup creation
- ✅ Backup integrity
- ✅ Required files present

**Expected Output:**
```
[TEST] Checking Docker...
✓ Docker is running
[TEST] Test 1: Testing database backup...
✓ Database backup created: backups/postgres/woragis_20241222_120000.dump.gz (2.5M)
[TEST] Test 2: Testing complete backup...
✓ Complete backup created: backups/woragis_backup_20241222_120000.tar.gz (5.2M)
[TEST] Test 3: Testing backup integrity...
✓ Backup archive is valid
✓ Manifest file present
✓ Database backup present
...
```

---

## Running All Tests

### Quick Test Suite

```bash
# 1. Start services
docker-compose up -d

# 2. Start server (in separate terminal)
cd server/app
go run ./cmd/server/main.go

# 3. Run security tests (in another terminal)
./scripts/test-security-middleware.sh

# 4. Run backup tests
./scripts/test-backups.sh
```

---

## Test Requirements

### Security Middleware Tests
- Server must be running
- Server must have security middleware enabled
- Port 8080 (or specified port) must be accessible

### Backup Tests
- Docker must be running
- Database container must exist
- Sufficient disk space for backups
- Write permissions in `backups/` directory

---

## Troubleshooting

### Security Tests Fail

**Issue:** "Server is not running"
```bash
# Start the server
cd server/app
go run ./cmd/server/main.go
```

**Issue:** "Rate limiting not working"
- Check middleware is enabled in `main.go`
- Verify rate limiter configuration
- Check Redis is running (if using Redis for rate limiting)

### Backup Tests Fail

**Issue:** "Docker is not running"
```bash
# Start Docker
# Linux/Mac: Start Docker Desktop
# Or: sudo systemctl start docker
```

**Issue:** "Database container not found"
```bash
# Start database container
docker-compose up -d database
```

**Issue:** "Permission denied"
```bash
# Make scripts executable
chmod +x scripts/*.sh
```

---

## Continuous Testing

### Add to CI/CD

You can add these tests to your CI/CD pipeline:

```yaml
# .github/workflows/test-production-readiness.yml
name: Test Production Readiness

on:
  pull_request:
    branches: [main, develop]

jobs:
  test-security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Start server
        run: |
          cd server/app
          go run ./cmd/server/main.go &
          sleep 5
      - name: Run security tests
        run: ./scripts/test-security-middleware.sh

  test-backups:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_PASSWORD: test
    steps:
      - uses: actions/checkout@v4
      - name: Run backup tests
        run: ./scripts/test-backups.sh
```

---

## Manual Testing

### Test Security Headers Manually

```bash
curl -I http://localhost:8080/healthz | grep -i "x-"
```

### Test Rate Limiting Manually

```bash
# Make 105 requests
for i in {1..105}; do
  curl -w "%{http_code}\n" -o /dev/null -s http://localhost:8080/healthz
done
# Should see 429 after 100 requests
```

### Test Request Size Limit Manually

```bash
# Create 11MB file
dd if=/dev/zero of=large.txt bs=1M count=11

# Try to POST it
curl -X POST http://localhost:8080/api/test \
  -H "Content-Type: application/json" \
  -d @large.txt
# Should return 413 Request Entity Too Large
```

---

**Last Updated:** 2025-12-22
