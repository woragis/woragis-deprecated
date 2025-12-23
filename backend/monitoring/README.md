# Monitoring & Observability

This directory contains the configuration for centralized logging and metrics collection using Loki, Prometheus, and Grafana.

## Architecture

### Logging
- **Loki**: Log aggregation system that collects and stores logs
- **Promtail**: Log shipper that collects logs from Docker containers and sends them to Loki

### Metrics
- **Prometheus**: Metrics collection and storage system
- **Grafana**: Visualization and dashboard tool for logs and metrics

## Services

The following services are added to `docker-compose.yml`:

- `loki`: Log aggregation backend (port 3100)
- `promtail`: Log collector (port 9080)
- `grafana`: Visualization UI (port 3000)

## Quick Start

1. **Start all monitoring services:**
   ```bash
   docker-compose up -d loki promtail prometheus grafana
   ```

   Or start individually:
   ```bash
   # Logging only
   docker-compose up -d loki promtail grafana
   
   # Metrics only
   docker-compose up -d prometheus grafana
   
   # Everything
   docker-compose up -d loki promtail prometheus grafana
   ```

2. **Access Grafana:**
   - URL: http://localhost:3000
   - Default username: `admin`
   - Default password: `admin` (change via `GRAFANA_ADMIN_PASSWORD` env var)

3. **View logs:**
   - Navigate to "Explore" in Grafana
   - Select "Loki" as the data source
   - Use LogQL queries to search logs

## LogQL Queries

### View all logs
```
{job="docker"}
```

### Filter by service
```
{job="docker", service="app"}
{job="docker", service="ai-service"}
{job="docker", service="resume-worker"}
```

### Filter by log level
```
{job="docker", level="error"}
{job="docker", level="warn"}
{job="docker", level="info"}
```

### Combine filters
```
{job="docker", service="app", level="error"}
```

### Count logs by service
```
sum(count_over_time({job="docker"}[1m])) by (service)
```

### Error rate by service
```
sum(rate({job="docker", level="error"}[5m])) by (service)
```

## Dashboards

A default dashboard is automatically provisioned:
- **Woragis Logs Overview**: Shows log volume, error rates, and recent logs

Access dashboards via: Grafana → Dashboards → Woragis Logs Overview

## Configuration Files

- `loki-config.yml`: Loki server configuration (retention, storage, etc.)
- `promtail-config.yml`: Promtail configuration (log collection rules)
- `grafana/provisioning/`: Grafana auto-provisioning configs
- `grafana/dashboards/`: Dashboard JSON definitions

## Log Retention

- **Default retention**: 30 days (720 hours)
- Configured in `loki-config.yml` under `limits_config.retention_period`

## Environment Variables

Set these in your `.env` file or `docker-compose.yml`:

- `GRAFANA_ADMIN_USER`: Grafana admin username (default: `admin`)
- `GRAFANA_ADMIN_PASSWORD`: Grafana admin password (default: `admin`)
- `GRAFANA_ROOT_URL`: Grafana root URL (default: `http://localhost:3000`)

## Troubleshooting

### Logs not appearing in Grafana

1. Check if Promtail is running:
   ```bash
   docker-compose ps promtail
   ```

2. Check Promtail logs:
   ```bash
   docker-compose logs promtail
   ```

3. Verify Loki is accessible:
   ```bash
   curl http://localhost:3100/ready
   ```

4. Check if containers are being discovered:
   - Promtail uses Docker service discovery
   - Ensure containers have the `com.docker.compose.project=woragis` label

### Promtail can't access Docker socket

On Linux, ensure the Docker socket is accessible:
```bash
sudo chmod 666 /var/run/docker.sock
```

On Windows/Mac with Docker Desktop, this should work automatically.

### High memory usage

- Reduce log retention period in `loki-config.yml`
- Adjust `max_streams_per_user` in `loki-config.yml`
- Consider using Loki's object storage backend for production

## Next Steps

1. Create custom dashboards for specific services
2. Set up alerts for critical errors
3. Configure log retention policies per service
4. Add log parsing rules for specific log formats
5. Integrate with Prometheus metrics (if using Prometheus)

## Documentation

### Logging
- **Quick Start**: `QUICK_START.md`
- **User Guide**: `USER_GUIDE.md`
- **Query Library**: `LOGQL_QUERY_LIBRARY.md`
- **Troubleshooting**: `TROUBLESHOOTING.md`
- **Retention Policy**: `RETENTION_POLICY.md`
- **Security**: `SECURITY.md`

### Metrics
- **Prometheus Guide**: `PROMETHEUS_GUIDE.md`
- **Quick Start**: `METRICS_QUICK_START.md`

### Tracing
- **Tracing Guide**: `TRACING_GUIDE.md`
- **Quick Start**: `TRACING_QUICK_START.md`

### Planning & Improvements
- **TODO**: `TODO.md` - Monitoring improvements and next steps

## References

### Logging
- [Grafana Loki Documentation](https://grafana.com/docs/loki/latest/)
- [LogQL Query Language](https://grafana.com/docs/loki/latest/logql/)
- [Promtail Configuration](https://grafana.com/docs/loki/latest/clients/promtail/configuration/)

### Metrics
- [Prometheus Documentation](https://prometheus.io/docs/)
- [PromQL Guide](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Grafana Prometheus Data Source](https://grafana.com/docs/grafana/latest/datasources/prometheus/)
