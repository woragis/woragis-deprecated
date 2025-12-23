# SSL/TLS Configuration Guide

**Last Updated:** 2025-12-22  
**Purpose:** Guide for configuring SSL/TLS certificates for secure communication

---

## Overview

This guide covers SSL/TLS configuration for:
- HTTPS for web services
- TLS for database connections
- TLS for Redis connections
- TLS for RabbitMQ connections
- Certificate management and renewal

---

## Current State

**Development:** HTTP only (no TLS)  
**Production:** TLS required for all external-facing services

---

## Certificate Options

### Option 1: Let's Encrypt (Recommended - Free)

**Best for:** Public-facing services, automatic renewal

**Features:**
- Free SSL certificates
- Automatic renewal
- 90-day validity (auto-renewed)
- Trusted by all browsers

**Setup:**
```bash
# Install certbot
sudo apt-get install certbot  # Ubuntu/Debian
brew install certbot  # macOS

# Generate certificate
sudo certbot certonly --standalone -d api.woragis.com

# Certificates stored in:
# /etc/letsencrypt/live/api.woragis.com/
#   - fullchain.pem (certificate + chain)
#   - privkey.pem (private key)
```

**Auto-renewal:**
```bash
# Add to crontab
0 0 * * * certbot renew --quiet --deploy-hook "docker-compose restart nginx"
```

---

### Option 2: Self-Signed Certificates (Development Only)

**Best for:** Local development, internal services

**Generate:**
```bash
# Generate private key
openssl genrsa -out server.key 2048

# Generate certificate signing request
openssl req -new -key server.key -out server.csr

# Generate self-signed certificate
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
```

**⚠️ Warning:** Self-signed certificates are not trusted by browsers. Only use for development.

---

### Option 3: Cloud Provider Certificates

**AWS Certificate Manager (ACM):**
- Free certificates
- Automatic renewal
- Integrated with AWS services

**Azure App Service Certificates:**
- Managed certificates
- Auto-renewal
- Integrated with Azure services

**Google Cloud SSL Certificates:**
- Managed certificates
- Auto-renewal
- Integrated with GCP services

---

## Service Configuration

### 1. Main API Server (Go/Fiber)

**With Reverse Proxy (Nginx/Traefik):**
```nginx
# nginx.conf
server {
    listen 443 ssl http2;
    server_name api.woragis.com;

    ssl_certificate /etc/letsencrypt/live/api.woragis.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.woragis.com/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    location / {
        proxy_pass http://app:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**Direct TLS in Go:**
```go
// main.go
import (
    "crypto/tls"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    // TLS configuration
    cfg := &tls.Config{
        MinVersion: tls.VersionTLS12,
        CurvePreferences: []tls.CurveID{
            tls.CurveP521,
            tls.CurveP384,
            tls.CurveP256,
        },
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        },
    }

    ln, err := tls.Listen("tcp", ":8443", cfg)
    if err != nil {
        log.Fatal(err)
    }

    log.Fatal(app.Listener(ln))
}
```

---

### 2. Database (PostgreSQL) TLS

**Enable TLS in PostgreSQL:**
```bash
# postgresql.conf
ssl = on
ssl_cert_file = '/var/lib/postgresql/server.crt'
ssl_key_file = '/var/lib/postgresql/server.key'
ssl_ca_file = '/var/lib/postgresql/ca.crt'
```

**Connection String:**
```bash
# Require TLS
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require

# Verify certificate
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=verify-full&sslcert=/path/to/client.crt&sslkey=/path/to/client.key&sslrootcert=/path/to/ca.crt
```

**Generate Certificates:**
```bash
# Create CA
openssl req -new -x509 -days 3650 -nodes -out ca.crt -keyout ca.key

# Create server certificate
openssl req -new -nodes -out server.csr -keyout server.key
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365
```

---

### 3. Redis TLS

**Redis Configuration:**
```bash
# redis.conf
port 0
tls-port 6380
tls-cert-file /etc/redis/redis.crt
tls-key-file /etc/redis/redis.key
tls-ca-cert-file /etc/redis/ca.crt
tls-protocols "TLSv1.2 TLSv1.3"
```

**Connection:**
```bash
# Go (go-redis)
rdb := redis.NewClient(&redis.Options{
    Addr: "localhost:6380",
    TLSConfig: &tls.Config{
        ServerName: "redis.woragis.com",
    },
})

# Python (redis-py)
import redis
r = redis.Redis(
    host='localhost',
    port=6380,
    ssl=True,
    ssl_cert_reqs='required',
    ssl_ca_certs='/path/to/ca.crt'
)
```

---

### 4. RabbitMQ TLS

**RabbitMQ Configuration:**
```bash
# rabbitmq.conf
listeners.ssl.default = 5671
ssl_options.cacertfile = /etc/rabbitmq/ca.crt
ssl_options.certfile = /etc/rabbitmq/server.crt
ssl_options.keyfile = /etc/rabbitmq/server.key
ssl_options.verify = verify_peer
ssl_options.fail_if_no_peer_cert = false
```

**Connection URL:**
```bash
# AMQPS (TLS)
RABBITMQ_URL=amqps://user:pass@rabbitmq:5671/vhost

# Go (amqp091-go)
conn, err := amqp.DialTLS("amqps://user:pass@rabbitmq:5671/vhost", &tls.Config{
    ServerName: "rabbitmq.woragis.com",
})
```

---

## Docker Compose TLS Setup

### Using Let's Encrypt with Traefik

```yaml
services:
  traefik:
    image: traefik:v2.10
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.letsencrypt.acme.tlschallenge=true"
      - "--certificatesresolvers.letsencrypt.acme.email=admin@woragis.com"
      - "--certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./letsencrypt:/letsencrypt
    labels:
      - "traefik.enable=true"

  app:
    image: woragis-app:latest
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.app.rule=Host(`api.woragis.com`)"
      - "traefik.http.routers.app.entrypoints=websecure"
      - "traefik.http.routers.app.tls.certresolver=letsencrypt"
```

---

## Certificate Renewal

### Let's Encrypt Auto-Renewal

```bash
# Create renewal script
cat > /usr/local/bin/renew-certificates.sh <<'EOF'
#!/bin/bash
certbot renew --quiet --deploy-hook "docker-compose -f /path/to/docker-compose.yml restart nginx"
EOF

chmod +x /usr/local/bin/renew-certificates.sh

# Add to crontab (runs twice daily)
0 0,12 * * * /usr/local/bin/renew-certificates.sh
```

### Manual Renewal

```bash
# Renew certificate
sudo certbot renew

# Test renewal (dry run)
sudo certbot renew --dry-run
```

---

## Security Best Practices

### 1. Use Strong Cipher Suites

```nginx
ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';
ssl_prefer_server_ciphers on;
```

### 2. Enable HSTS (HTTP Strict Transport Security)

```nginx
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
```

### 3. Disable Weak Protocols

```nginx
ssl_protocols TLSv1.2 TLSv1.3;
```

### 4. Use Perfect Forward Secrecy

- Use ECDHE cipher suites
- Disable static RSA key exchange

### 5. Certificate Pinning (Optional)

For mobile apps, pin certificates to prevent MITM attacks.

---

## Testing TLS Configuration

### Test Certificate

```bash
# Check certificate validity
openssl s_client -connect api.woragis.com:443 -showcerts

# Check certificate expiration
echo | openssl s_client -connect api.woragis.com:443 2>/dev/null | openssl x509 -noout -dates

# Test TLS version
openssl s_client -connect api.woragis.com:443 -tls1_2
```

### Online Tools

- **SSL Labs:** https://www.ssllabs.com/ssltest/
- **SSL Test:** https://www.ssllabs.com/ssltest/analyze.html

---

## Troubleshooting

### Certificate Not Trusted

**Issue:** Browser shows "Not Secure" warning

**Solutions:**
- Ensure certificate is from trusted CA (Let's Encrypt, etc.)
- Check certificate chain is complete
- Verify domain matches certificate

### Certificate Expired

**Issue:** Certificate expired, connections fail

**Solutions:**
- Renew certificate immediately
- Set up auto-renewal
- Monitor certificate expiration

### TLS Handshake Failed

**Issue:** Connection fails with TLS error

**Solutions:**
- Check TLS version compatibility
- Verify cipher suite support
- Check certificate validity
- Review server logs

---

## Implementation Checklist

### Phase 1: External Services (Week 1)
- [ ] Obtain SSL certificates (Let's Encrypt)
- [ ] Configure reverse proxy (Nginx/Traefik)
- [ ] Enable HTTPS for main API
- [ ] Test HTTPS connections
- [ ] Set up certificate auto-renewal

### Phase 2: Internal Services (Week 2)
- [ ] Generate internal CA
- [ ] Create certificates for database
- [ ] Create certificates for Redis
- [ ] Create certificates for RabbitMQ
- [ ] Configure services to use TLS
- [ ] Test internal TLS connections

### Phase 3: Hardening (Week 2)
- [ ] Configure strong cipher suites
- [ ] Enable HSTS
- [ ] Disable weak protocols
- [ ] Set up monitoring for certificate expiration
- [ ] Document procedures

---

## Related Documentation

- [Secrets Management](./secrets-management.md)
- [Docker Setup Guide](./docker-setup.md)
- [Deployment Procedures](./deployment-procedures.md) (when created)
- [Security Guide](../operations/security.md) (when created)

---

**Next Steps:**
1. Set up Let's Encrypt for external services
2. Generate internal CA for service-to-service communication
3. Configure TLS for all services
4. Set up certificate monitoring and auto-renewal

---

**Last Updated:** 2025-12-22
