# Configuration Reference

**Last Updated:** 2025-12-22  
**Purpose:** Complete reference for all environment variables and configuration options

---

## Overview

The Woragis backend uses environment variables for configuration. Variables can be set in:
- `.env` file (for local development)
- Docker Compose environment variables
- System environment variables
- Container environment variables

**Priority:** System env > Docker Compose > `.env` file

---

## Application Configuration

### Core Application Settings

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `APP_NAME` | Application name | `woragis-server` | No |
| `APP_ENV` | Environment (development, staging, production) | `development` | No |
| `APP_PORT` | Main server port | `8080` | No |
| `APP_PUBLIC_URL` | Public URL for the application | `http://localhost:5173` | Yes (production) |

**Example:**
```bash
APP_NAME=woragis-server
APP_ENV=production
APP_PORT=8080
APP_PUBLIC_URL=https://api.woragis.com
```

---

## Database Configuration

### PostgreSQL

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DATABASE_URL` | PostgreSQL connection string | - | Yes |
| `DB_HOST` | Database host (if not using DATABASE_URL) | `database` | No |
| `DB_PORT` | Database port | `5432` | No |
| `DB_USER` | Database user | `postgres` | No |
| `DB_PASSWORD` | Database password | `postgres` | No |
| `DB_NAME` | Database name | `woragis` | No |
| `DB_SSLMODE` | SSL mode (disable, require, etc.) | `disable` | No |

**Example:**
```bash
DATABASE_URL=postgres://postgres:postgres@database:5432/woragis?sslmode=disable
```

**Production Example:**
```bash
DATABASE_URL=postgres://user:password@db.example.com:5432/woragis?sslmode=require
```

---

## Redis Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `REDIS_URL` | Redis connection string | `redis://redis:6379/0` | No |
| `REDIS_HOST` | Redis host (if not using REDIS_URL) | `redis` | No |
| `REDIS_PORT` | Redis port | `6379` | No |
| `REDIS_DB` | Redis database number | `0` | No |
| `REDIS_PASSWORD` | Redis password (if required) | - | No |

**Example:**
```bash
REDIS_URL=redis://redis:6379/0
# Or with password:
REDIS_URL=redis://:password@redis:6379/0
```

---

## RabbitMQ Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `RABBITMQ_URL` | RabbitMQ AMQP connection URL | `amqp://woragis:woragis@rabbitmq:5672/` | Yes |
| `RABBITMQ_USER` | RabbitMQ username | `woragis` | No |
| `RABBITMQ_PASSWORD` | RabbitMQ password | `woragis` | No |
| `RABBITMQ_VHOST` | RabbitMQ virtual host | `/` | No |
| `RABBITMQ_HOST` | RabbitMQ host | `rabbitmq` | No |
| `RABBITMQ_PORT` | RabbitMQ port | `5672` | No |

**Example:**
```bash
RABBITMQ_URL=amqp://woragis:woragis@rabbitmq:5672/
RABBITMQ_VHOST=/
```

**Note:** Use `/` (root vhost) for default setup. Custom vhosts require RabbitMQ configuration.

---

## CORS Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `CORS_ENABLED` | Enable CORS | `true` | No |
| `CORS_ALLOWED_ORIGINS` | Comma-separated list of allowed origins | - | Yes (if CORS enabled) |
| `CORS_ALLOWED_METHODS` | Allowed HTTP methods | `GET,POST,PUT,PATCH,DELETE,OPTIONS` | No |
| `CORS_ALLOWED_HEADERS` | Allowed headers | `Authorization,Content-Type,X-Requested-With` | No |
| `CORS_EXPOSED_HEADERS` | Headers to expose | - | No |
| `CORS_ALLOW_CREDENTIALS` | Allow credentials | `true` | No |
| `CORS_MAX_AGE` | Max age for preflight requests (seconds) | `86400` | No |

**Example:**
```bash
CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=http://localhost:5173,https://app.woragis.com
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
CORS_ALLOW_CREDENTIALS=true
```

---

## AI Service Configuration

### Chat Provider Selection

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `CHAT_PROVIDER` | AI provider (openai, anthropic, xai, manus, cipher) | `openai` | No |
| `DEFAULT_TEMPERATURE` | Default temperature for AI requests | `0.3` | No |

### OpenAI Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `OPENAI_API_KEY` | OpenAI API key | - | Yes (if using OpenAI) |
| `OPENAI_MODEL` | OpenAI model name | `gpt-4o-mini` | No |
| `OPENAI_TEMPERATURE` | Temperature for OpenAI requests | `0.3` | No |

**Example:**
```bash
CHAT_PROVIDER=openai
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4o-mini
```

### Anthropic Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `ANTHROPIC_API_KEY` | Anthropic API key | - | Yes (if using Anthropic) |
| `ANTHROPIC_MODEL` | Anthropic model name | `claude-3-5-sonnet-latest` | No |

### X.AI Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `XAI_API_KEY` | X.AI API key | - | Yes (if using X.AI) |
| `XAI_BASE_URL` | X.AI API base URL | `https://api.x.ai/v1` | No |
| `XAI_MODEL` | X.AI model name | `grok-beta` | No |

### Manus Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `MANUS_API_KEY` | Manus API key | - | Yes (if using Manus) |
| `MANUS_BASE_URL` | Manus API base URL | - | Yes (if using Manus) |
| `MANUS_MODEL` | Manus model name | `manus-chat` | No |

### Cipher Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `CIPHER_API_KEY` | Cipher API key | - | Yes (if using Cipher) |
| `CIPHER_BASE_URL` | Cipher API base URL | `https://api.nofiltergpt.com/v1/chat/completions` | No |
| `CIPHER_MAX_TOKENS` | Max tokens for Cipher | `800` | No |
| `CIPHER_TOP_P` | Top-p for Cipher | `1.0` | No |
| `CIPHER_IMAGE_URL` | Cipher image generation URL | `https://api.nofiltergpt.com/v1/images/generations` | No |
| `CIPHER_IMAGE_SIZE` | Image size | `1024x1024` | No |
| `CIPHER_IMAGE_N` | Number of images | `1` | No |

---

## Email Configuration (SMTP)

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SMTP_HOST` | SMTP server host | - | Yes (if email enabled) |
| `SMTP_PORT` | SMTP server port | `587` | No |
| `SMTP_USERNAME` | SMTP username | - | Yes (if email enabled) |
| `SMTP_PASSWORD` | SMTP password | - | Yes (if email enabled) |
| `SMTP_FROM` | Default "from" email address | - | Yes (if email enabled) |
| `SMTP_TLS` | Enable TLS | `true` | No |

**Example:**
```bash
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USERNAME=noreply@mail.woragis.me
SMTP_PASSWORD=your-password
SMTP_FROM=noreply@mail.woragis.me
SMTP_TLS=true
```

---

## WhatsApp Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `WHATSAPP_ENABLED` | Enable WhatsApp worker | `true` | No |
| `WHATSAPP_SESSION_PATH` | Path to WhatsApp session files | `./whatsapp-session` | No |

**Example:**
```bash
WHATSAPP_ENABLED=true
WHATSAPP_SESSION_PATH=./whatsapp-worker/session
```

---

## Translation Worker Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `TRANSLATION_PROVIDER` | Translation provider (google, deepl, libre) | `libre` | No |
| `GOOGLE_TRANSLATE_API_KEY` | Google Translate API key | - | No (falls back to LibreTranslate) |
| `GOOGLE_PROJECT_ID` | Google Cloud project ID | - | No |
| `DEEPL_API_KEY` | DeepL API key | - | No (falls back to LibreTranslate) |
| `LIBRE_API_URL` | LibreTranslate API URL | `https://libretranslate.com/translate` | No |
| `LIBRE_API_KEY` | LibreTranslate API key (optional) | - | No |
| `TRANSLATION_TIMEOUT` | Translation request timeout | `30s` | No |
| `TRANSLATION_MAX_RETRIES` | Maximum retry attempts | `3` | No |
| `TRANSLATION_RETRY_DELAY` | Delay between retries | `1s` | No |

**Example:**
```bash
TRANSLATION_PROVIDER=google
GOOGLE_TRANSLATE_API_KEY=your-key-here
GOOGLE_PROJECT_ID=your-project-id
```

---

## Monitoring & Metrics Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `MONITORING_ENABLED` | Enable monitoring | `true` | No |
| `MONITORING_DB_URL` | Monitoring database URL (optional) | - | No |
| `METRICS_NAMESPACE` | Metrics namespace | `woragis` | No |

**Example:**
```bash
MONITORING_ENABLED=true
METRICS_NAMESPACE=woragis
```

---

## Grafana Configuration (Monitoring Stack)

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `GRAFANA_ADMIN_USER` | Grafana admin username | `admin` | No |
| `GRAFANA_ADMIN_PASSWORD` | Grafana admin password | `admin` | No |
| `GRAFANA_ROOT_URL` | Grafana root URL | `http://localhost:3000` | No |

**Example:**
```bash
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=secure-password
GRAFANA_ROOT_URL=http://localhost:3000
```

---

## Security Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PUBLIC_API_KEY` | Public API key for authentication | - | Yes (production) |
| `JWT_SECRET` | JWT secret key (if using JWT) | - | Yes (if JWT enabled) |
| `SESSION_SECRET` | Session secret key | - | Yes (if sessions enabled) |

**Example:**
```bash
PUBLIC_API_KEY=your-secure-api-key-here
JWT_SECRET=your-jwt-secret-here
```

**⚠️ Security Warning:** Never commit secrets to version control. Use secrets management in production.

---

## Service-Specific Configuration

### AI Service

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `AI_SERVICE_URL` | AI service URL (for workers) | `http://ai-service:8000` | No |
| `AI_SERVICE_PORT` | AI service port | `8000` | No |

### Creative Service

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `CREATIVE_SERVICE_URL` | Creative service URL | `http://creative-service:8001` | No |
| `CREATIVE_SERVICE_PORT` | Creative service port | `8001` | No |

### Docs Service

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DOCS_SERVICE_URL` | Docs service URL | `http://docs-service:8002` | No |
| `DOCS_SERVICE_PORT` | Docs service port | `8002` | No |

---

## Environment-Specific Configuration

### Development

```bash
APP_ENV=development
DATABASE_URL=postgres://postgres:postgres@localhost:5432/woragis?sslmode=disable
REDIS_URL=redis://localhost:6379/0
CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
```

### Production

```bash
APP_ENV=production
DATABASE_URL=postgres://user:password@db.example.com:5432/woragis?sslmode=require
REDIS_URL=redis://:password@redis.example.com:6379/0
CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=https://app.woragis.com
MONITORING_ENABLED=true
```

---

## Configuration Validation

### Required Variables Checklist

**For Basic Setup:**
- [ ] `DATABASE_URL`
- [ ] `REDIS_URL`
- [ ] `RABBITMQ_URL`
- [ ] `CORS_ALLOWED_ORIGINS` (if CORS enabled)

**For AI Features:**
- [ ] `OPENAI_API_KEY` (or other provider key)
- [ ] `CHAT_PROVIDER`

**For Email:**
- [ ] `SMTP_HOST`
- [ ] `SMTP_USERNAME`
- [ ] `SMTP_PASSWORD`
- [ ] `SMTP_FROM`

**For Production:**
- [ ] `APP_PUBLIC_URL`
- [ ] `PUBLIC_API_KEY`
- [ ] Secure database credentials
- [ ] Secure Redis credentials
- [ ] Secure RabbitMQ credentials

---

## Configuration Best Practices

1. **Use `.env` files for local development** - Never commit `.env` files
2. **Use secrets management in production** - Use Docker secrets, Kubernetes secrets, or cloud secret managers
3. **Validate configuration on startup** - Services should validate required variables
4. **Use environment-specific configs** - Different values for dev/staging/prod
5. **Document defaults** - Always document default values
6. **Use connection strings** - Prefer `DATABASE_URL` over individual variables
7. **Rotate secrets regularly** - Change API keys and passwords periodically

---

## Troubleshooting Configuration

### Issue: Service won't start

**Check:**
1. All required variables are set
2. Connection strings are valid
3. Credentials are correct
4. Ports are not in use

### Issue: Connection errors

**Check:**
1. Database/Redis/RabbitMQ are running
2. Connection strings use correct hostnames
3. Network connectivity (Docker network)
4. Credentials are correct

### Issue: Environment variables not loading

**Check:**
1. `.env` file is in correct location
2. Docker Compose is reading `.env` file
3. Variable names are correct (case-sensitive)
4. No typos in variable names

---

## Related Documentation

- [Development Setup Guide](../development/setup-guide.md)
- [Docker Setup Guide](./docker-setup.md)
- [Deployment Procedures](./deployment-procedures.md) (when created)
- [Environment Variables Reference](./environment-variables.md) (detailed reference)

---

**Last Updated:** 2025-12-22
