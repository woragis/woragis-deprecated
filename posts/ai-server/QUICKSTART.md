# Quick Start: Posts AI Service

## 5-Minute Setup

### Prerequisites

- Go 1.21+
- PostgreSQL running locally
- AI Service running on port 8000

### Step 1: Clone & Navigate

```bash
cd c:\Users\Jezreel de Andrade\dev\Projects\woragis\posts\ai-server
```

### Step 2: Setup Environment

```bash
cp .env.example .env
# Edit .env if needed (defaults should work for local dev)
```

### Step 3: Create Database

```bash
# Using psql directly
psql -U postgres -c "CREATE DATABASE posts_ai;"

# Or use Makefile
make db-create
```

### Step 4: Start Service

```bash
# Option A: Direct run (migrations run automatically)
go run cmd/server/main.go

# Option B: With hot reload
go install github.com/cosmtrek/air@latest
air

# Option C: Docker
docker-compose up
```

Service will be available at `http://localhost:3014`

### Step 5: Verify Health

```bash
curl http://localhost:3014/healthz
# Should return: {"status":"healthy","service":"posts-ai-service"}
```

## Testing Endpoints

### Generate Draft

```bash
curl -X POST http://localhost:3014/api/v1/chats/generate \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "prompt": "Write a beginner guide to Go microservices",
    "agent": "auto"
  }'
```

### List Chats

```bash
curl http://localhost:3014/api/v1/chats?user_id=550e8400-e29b-41d4-a716-446655440000&limit=10
```

### Get Chat Details

```bash
curl http://localhost:3014/api/v1/chats/{chat_id}
```

## Integration with Frontend

See `INTEGRATION.md` for detailed frontend setup:

1. Update `.env` with `PUBLIC_POSTS_AI_SERVICE_URL`
2. Create Posts-AI client service
3. Update components to route through Posts-AI Service
4. Benefits: persistence, audit trail, cost tracking

## Common Commands

```bash
# Build binary
make build

# Run tests
make test

# Format code
make fmt

# Lint code
make lint

# Database operations
make db-create    # Create database
make db-reset     # Drop and recreate
make psql         # Open psql shell

# Docker operations
make docker-build  # Build image
make docker-up     # Start stack
make docker-down   # Stop stack
```

## Architecture Overview

```
Posts-AI Service (Go)
├── cmd/server/main.go          Entry point
├── internal/
│   ├── config/                 Configuration
│   ├── db/                      PostgreSQL
│   ├── services/
│   │   ├── ai.go               AI service client
│   │   └── chat.go             Chat business logic
│   ├── handlers/
│   │   ├── chat.go             HTTP handlers
│   │   └── websocket.go        WebSocket support
│   └── models/                 Data models
└── Database (PostgreSQL)
    ├── chats
    ├── chat_turns
    └── ai_usage_stats
```

## Endpoints at a Glance

| Method | Endpoint                       | Purpose          |
| ------ | ------------------------------ | ---------------- |
| GET    | `/healthz`                     | Health check     |
| POST   | `/api/v1/chats/generate`       | Generate draft   |
| POST   | `/api/v1/posts/:id/ai/improve` | Improve content  |
| GET    | `/api/v1/chats/:id`            | Get chat details |
| GET    | `/api/v1/chats`                | List chats       |
| GET    | `/api/v1/usage/stats`          | Get usage stats  |
| WS     | `/ws/chats/:id`                | Real-time chat   |

## Environment Variables

```env
PORT=3014                                    # Server port
DATABASE_URL=postgres://...                  # PostgreSQL connection
AI_SERVICE_URL=http://localhost:8000         # AI service
CORS_ORIGINS=http://localhost:5173,...       # Allowed CORS origins
ENV=development                              # development|production
LOG_LEVEL=info                               # debug|info|warn|error
```

## Troubleshooting

### "connection refused" to database

```bash
# Check PostgreSQL is running
psql -U woragis -d posts_ai

# Or restart in Docker
docker-compose down
docker-compose up -d
```

### "connection refused" to AI service

```bash
# Verify AI service is running
curl http://localhost:8000/healthz

# Check service URL in .env
```

### Build errors

```bash
# Clean and rebuild
make clean
go mod download
go mod tidy
make build
```

## Next Steps

1. ✅ Service running locally
2. ✅ Database with tables and indexes
3. ✅ Health check passing
4. Next → Integrate with Posts Frontend (see `INTEGRATION.md`)
5. Next → Deploy to staging/production

## Documentation

- `README.md` - Full documentation
- `INTEGRATION.md` - Frontend integration guide
- `Makefile` - Common commands
- `.env.example` - Configuration template

---

**Need help?**

- Check logs: Look at terminal output or `docker-compose logs`
- Test endpoint: Use `curl` or Postman
- See integration guide: `INTEGRATION.md`
