# Quick Start: Logging Aggregation

## What Was Implemented

✅ **Phase 3: Infrastructure Setup** (Complete)
- Loki service added to docker-compose.yml
- Promtail service added for log collection
- Grafana service added with auto-provisioning
- Persistent volumes configured
- Basic dashboard created

## Starting the Services

1. **Start all logging services:**
   ```bash
   docker-compose up -d loki promtail grafana
   ```

2. **Or start with all services:**
   ```bash
   docker-compose up -d
   ```

3. **Verify services are running:**
   ```bash
   docker-compose ps loki promtail grafana
   ```

## Accessing Grafana

1. Open your browser: http://localhost:3000
2. Login with:
   - Username: `admin`
   - Password: `admin` (change this in production!)

3. The Loki data source is automatically configured
4. The "Woragis Logs Overview" dashboard is available

## Testing Log Collection

1. **Generate some logs** by using your services:
   ```bash
   # Make a request to your API
   curl http://localhost:8080/health
   ```

2. **View logs in Grafana:**
   - Go to "Explore" (compass icon in left menu)
   - Select "Loki" data source
   - Try query: `{job="docker"}`

3. **View the dashboard:**
   - Go to "Dashboards" → "Woragis Logs Overview"

## Common LogQL Queries

```logql
# All logs
{job="docker"}

# By service
{job="docker", service="app"}
{job="docker", service="ai-service"}

# Errors only
{job="docker", level="error"}

# Errors from specific service
{job="docker", service="app", level="error"}

# Log volume by service
sum(count_over_time({job="docker"}[1m])) by (service)
```

## Next Steps

1. **Customize retention** in `monitoring/loki-config.yml`
2. **Create custom dashboards** for specific services
3. **Set up alerts** for critical errors
4. **Review log formats** to ensure structured logging is working

## Troubleshooting

### Services won't start
```bash
# Check logs
docker-compose logs loki
docker-compose logs promtail
docker-compose logs grafana
```

### No logs appearing
1. Check Promtail is discovering containers:
   ```bash
   docker-compose logs promtail | grep "discovered"
   ```

2. Verify Docker socket access (Linux):
   ```bash
   ls -la /var/run/docker.sock
   ```

3. Check Loki is receiving logs:
   ```bash
   curl http://localhost:3100/ready
   ```

### Grafana can't connect to Loki
- Verify Loki is healthy: `docker-compose ps loki`
- Check network connectivity between containers
- Review Grafana logs: `docker-compose logs grafana`

## Configuration Files

- `loki-config.yml` - Loki server settings
- `promtail-config.yml` - Log collection rules
- `grafana/provisioning/` - Auto-configuration
- `grafana/dashboards/` - Dashboard definitions

## Environment Variables

Add to `.env` file:
```bash
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=your-secure-password
GRAFANA_ROOT_URL=http://localhost:3000
```
