# Environment Variables Guide

## How Environment Variables Work

### Development vs Production

**Development Mode (`ENV != "production"`):**
- `.env` file is **automatically loaded** from the server directory
- If `.env` file is missing, it logs a warning but continues
- Environment variables can be set in `.env` file OR via OS environment
- OS environment variables take precedence over `.env` file

**Production Mode (`ENV=production`):**
- `.env` file is **NOT loaded**
- All environment variables **MUST** be set via OS environment (Docker, Kubernetes, system env, etc.)
- The application expects all variables to be set externally
- Missing required variables will cause the application to fail at startup

### Code Logic

```go
// From config/config.go
if os.Getenv("ENV") != "production" {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
}
```

## Required Environment Variables

### Both Auth and Jobs Services

These variables **MUST** be set (using `getEnvRequired()` - will panic if missing):

1. **DATABASE_URL** - PostgreSQL connection string
   - Example: `postgres://user:password@host:5432/dbname?sslmode=disable`
   
2. **REDIS_URL** - Redis connection string
   - Example: `redis://localhost:6379/0`
   
3. **AES_KEY** - AES encryption key (must be exactly 32 bytes)
   - Used for encrypting sensitive data
   - Example: `65b41ee0f9c63c4b0f2b584e165d0439` (32 hex characters = 32 bytes)
   
4. **HASH_SALT** - Salt for password hashing
   - Used for additional security in password hashing
   - Example: `65b41ee0f9c63c4b0f2b584e165d0439`

### Auth Service Specific

5. **AUTH_JWT_SECRET** (or **JWT_SECRET**) - Required in production
   - JWT secret key for signing tokens
   - In development: defaults to `"dev-secret-change-me"` if not set
   - In production: **MUST** be set, will fail if using default
   - Should be a long, random, secure string

## Optional Environment Variables

These have defaults and are not required:

### Application Configuration
- **APP_ENV** / **ENV** - Environment mode (default: `"development"`)
- **APP_NAME** - Application name (default: `"woragis-auth-service"` or `"woragis-jobs-service"`)
- **APP_PORT** / **PORT** - Port to listen on (default: `"3000"`)
- **APP_PUBLIC_URL** - Public URL of the service (default: `"http://localhost:3000"`)

### Database Configuration
- **DATABASE_MAX_OPEN_CONNS** - Max open connections (default: `25`)
- **DATABASE_MAX_IDLE_CONNS** - Max idle connections (default: `25`)
- **DATABASE_MAX_IDLE_TIME** - Max idle time (default: `15m`)
- **DATABASE_CONN_MAX_LIFETIME** - Connection max lifetime (default: `60m`)

### Redis Configuration
- **REDIS_PASSWORD** - Redis password (default: `""`)
- **REDIS_DB** - Redis database number (default: `0`)

### Auth Configuration
- **AUTH_JWT_TTL** / **JWT_EXPIRE_HOURS** - JWT expiration in hours (default: `24`)
- **JWT_REFRESH_EXPIRE_HOURS** - Refresh token expiration in hours (default: `168` = 7 days)
- **BCRYPT_COST** - Bcrypt cost factor (default: `12`)

### CORS Configuration
- **CORS_ENABLED** - Enable CORS (default: `true`)
- **CORS_ALLOWED_ORIGINS** - Allowed origins (comma-separated)
- **CORS_ALLOWED_METHODS** - Allowed methods (default: `GET,POST,PUT,PATCH,DELETE,OPTIONS`)
- **CORS_ALLOWED_HEADERS** - Allowed headers
- **CORS_EXPOSED_HEADERS** - Exposed headers
- **CORS_ALLOW_CREDENTIALS** - Allow credentials (default: `true`)
- **CORS_MAX_AGE** - Max age in seconds (default: `86400`)

### Email Configuration (Optional)
- **SMTP_HOST** - SMTP server host
- **SMTP_PORT** - SMTP server port (default: `587`)
- **SMTP_USERNAME** - SMTP username
- **SMTP_PASSWORD** - SMTP password
- **SMTP_FROM** / **EMAIL_FROM** - From email address
- **SMTP_TLS** - Use TLS (default: `true`)

### Monitoring & Tracing (Optional)
- **MONITORING_ENABLED** - Enable monitoring (default: `true`)
- **METRICS_NAMESPACE** - Metrics namespace
- **OTLP_ENDPOINT** - OpenTelemetry endpoint
- **JAEGER_ENDPOINT** - Jaeger endpoint

### RabbitMQ Configuration (Jobs Service)
- **RABBITMQ_URL** - RabbitMQ connection URL
- **RABBITMQ_USER** - RabbitMQ user (default: `woragis`)
- **RABBITMQ_PASSWORD** - RabbitMQ password (default: `woragis`)
- **RABBITMQ_HOST** - RabbitMQ host (default: `rabbitmq`)
- **RABBITMQ_PORT** - RabbitMQ port (default: `5672`)
- **RABBITMQ_VHOST** - RabbitMQ virtual host (default: `/`)

## Setting Environment Variables

### Development (with .env file)

Create a `.env` file in the `server/` directory:

```bash
# Required
DATABASE_URL=postgres://woragis:password@localhost:5442/auth_service?sslmode=disable
REDIS_URL=redis://localhost:6389/0
AES_KEY=65b41ee0f9c63c4b0f2b584e165d0439
HASH_SALT=65b41ee0f9c63c4b0f2b584e165d0439
AUTH_JWT_SECRET=your-secret-key-here

# Optional
ENV=development
APP_PORT=3000
APP_PUBLIC_URL=http://localhost:3010
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
```

### Production (Docker Compose)

Set variables in `docker-compose.yml`:

```yaml
environment:
  ENV: production
  DATABASE_URL: ${DATABASE_URL}
  REDIS_URL: ${REDIS_URL}
  AES_KEY: ${AES_KEY}
  HASH_SALT: ${HASH_SALT}
  AUTH_JWT_SECRET: ${AUTH_JWT_SECRET}
  # ... other variables
```

Or use environment file:

```yaml
env_file:
  - .env.production
```

### Production (Kubernetes)

Set via ConfigMap or Secrets:

```yaml
env:
  - name: DATABASE_URL
    valueFrom:
      secretKeyRef:
        name: auth-secrets
        key: database-url
  - name: AES_KEY
    valueFrom:
      secretKeyRef:
        name: auth-secrets
        key: aes-key
```

## Summary

1. **Development**: `.env` file is loaded automatically (if `ENV != "production"`)
2. **Production**: `.env` file is **NOT loaded** - set variables via OS environment
3. **Required variables**: `DATABASE_URL`, `REDIS_URL`, `AES_KEY`, `HASH_SALT`, `AUTH_JWT_SECRET` (production)
4. **Missing required variables**: Application will fail to start with a fatal error
5. **Optional variables**: Have sensible defaults, can be overridden as needed

