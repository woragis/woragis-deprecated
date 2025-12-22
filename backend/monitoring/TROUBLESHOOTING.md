# Troubleshooting Guide

**Version:** 1.0  
**Last Updated:** 2025-12-22

## Common Issues and Solutions

### No Logs Appearing in Grafana

#### Symptoms
- Grafana shows no logs
- Queries return empty results
- "No data" message in panels

#### Diagnosis

1. **Check if services are running:**
   ```bash
   docker-compose ps
   ```

2. **Check if Loki is running:**
   ```bash
   docker-compose ps loki
   docker-compose logs loki
   ```

3. **Check if Promtail is running:**
   ```bash
   docker-compose ps promtail
   docker-compose logs promtail
   ```

4. **Verify Loki is receiving logs:**
   ```bash
   curl http://localhost:3100/ready
   curl http://localhost:3100/metrics
   ```

5. **Check Promtail logs for errors:**
   ```bash
   docker-compose logs promtail | grep -i error
   ```

#### Solutions

**Issue: Promtail can't access Docker socket**
```bash
# On Linux, ensure Docker socket is accessible
sudo chmod 666 /var/run/docker.sock

# Or add user to docker group
sudo usermod -aG docker $USER
```

**Issue: Promtail not discovering containers**
- Check Docker Compose project name matches filter in `promtail-config.yml`
- Verify containers have correct labels
- Check Promtail logs for discovery messages

**Issue: Loki not accessible**
- Verify Loki is healthy: `docker-compose ps loki`
- Check network connectivity: `docker network inspect woragis_default`
- Review Loki logs: `docker-compose logs loki`

### Promtail Can't Access Docker Socket

#### Symptoms
- Promtail logs show "permission denied"
- No containers discovered
- Errors about Docker socket

#### Solutions

**Linux:**
```bash
# Check socket permissions
ls -la /var/run/docker.sock

# Fix permissions (temporary)
sudo chmod 666 /var/run/docker.sock

# Fix permissions (permanent - add user to docker group)
sudo usermod -aG docker $USER
# Log out and back in
```

**Windows/Mac:**
- Docker Desktop should handle this automatically
- If issues persist, restart Docker Desktop

**Docker Compose:**
- Ensure Promtail has access to Docker socket in `docker-compose.yml`:
  ```yaml
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock:ro
  ```

### High Memory Usage

#### Symptoms
- Loki using excessive memory
- System running out of memory
- Slow queries

#### Solutions

1. **Reduce retention period:**
   Edit `monitoring/loki-config.yml`:
   ```yaml
   limits_config:
     retention_period: 168h  # 7 days instead of 30
   ```

2. **Limit log volume:**
   - Reduce log verbosity in services
   - Filter out health check logs
   - Use appropriate log levels

3. **Adjust Loki limits:**
   ```yaml
   limits_config:
     max_streams_per_user: 5000  # Reduce from 10000
     max_line_size: 128KB  # Reduce from 256KB
   ```

4. **Restart Loki:**
   ```bash
   docker-compose restart loki
   ```

### Slow Queries

#### Symptoms
- Queries take a long time
- Timeouts in Grafana
- High CPU usage

#### Solutions

1. **Use label filters:**
   ```
   ✅ {job="docker", service="app", level="error"}
   ❌ {job="docker"} |= "app" |= "error"
   ```

2. **Specify time ranges:**
   ```
   ✅ {job="docker"} [1h]
   ❌ {job="docker"}  # May query all time
   ```

3. **Use shorter time ranges:**
   - Use 1h instead of 24h when possible
   - Use 5m for recent logs

4. **Avoid broad regex:**
   ```
   ✅ {job="docker", level="error"} |= "database"
   ❌ {job="docker"} |~ ".*database.*"
   ```

5. **Use aggregations efficiently:**
   ```
   ✅ sum(rate({job="docker"}[5m]))
   ❌ count_over_time({job="docker"}[1h])  # More expensive
   ```

### Grafana Can't Connect to Loki

#### Symptoms
- "Data source error" in Grafana
- "Connection refused" errors
- Loki data source shows as offline

#### Solutions

1. **Verify Loki is running:**
   ```bash
   docker-compose ps loki
   ```

2. **Check Loki URL in Grafana:**
   - Go to Configuration → Data Sources → Loki
   - Verify URL is: `http://loki:3100`
   - Test connection

3. **Check network:**
   ```bash
   docker network inspect woragis_default | grep -A 5 loki
   ```

4. **Restart Grafana:**
   ```bash
   docker-compose restart grafana
   ```

5. **Check Grafana logs:**
   ```bash
   docker-compose logs grafana | grep -i loki
   ```

### Logs Not Parsing Correctly

#### Symptoms
- Logs appear as plain text
- JSON fields not extracted
- Labels not showing

#### Solutions

1. **Verify log format:**
   - Check if services output JSON in production
   - Verify log format matches specification

2. **Check Promtail pipeline:**
   - Review `promtail-config.yml` pipeline stages
   - Verify JSON parsing is configured

3. **Test log format:**
   ```bash
   # Check a service's log output
   docker-compose logs app | head -5
   ```

4. **Update Promtail config:**
   - Adjust pipeline stages in `promtail-config.yml`
   - Restart Promtail: `docker-compose restart promtail`

### Disk Space Issues

#### Symptoms
- Disk full errors
- Loki failing to write
- Docker volume full

#### Solutions

1. **Check disk usage:**
   ```bash
   docker system df
   docker volume inspect woragis_loki-data
   ```

2. **Reduce retention:**
   Edit `monitoring/loki-config.yml`:
   ```yaml
   limits_config:
     retention_period: 168h  # 7 days
   ```
   Restart Loki: `docker-compose restart loki`

3. **Clean up old data:**
   ```bash
   # Remove old Docker logs
   docker system prune -a --volumes
   ```

4. **Increase disk space:**
   - Add more disk space
   - Move volumes to larger disk
   - Use external storage

### Alerts Not Firing

#### Symptoms
- Alerts configured but not triggering
- No notifications received

#### Solutions

1. **Check alert rules:**
   - Go to Alerting → Alert rules
   - Verify rules are enabled
   - Check evaluation interval

2. **Test query:**
   - Run the alert query in Explore
   - Verify it returns data
   - Check if threshold is correct

3. **Check notification channels:**
   - Go to Alerting → Notification channels
   - Verify channels are configured
   - Test channel

4. **Review alert logs:**
   ```bash
   docker-compose logs grafana | grep -i alert
   ```

### Services Not Appearing in Logs

#### Symptoms
- Some services missing from logs
- No logs from specific service

#### Solutions

1. **Check service is running:**
   ```bash
   docker-compose ps <service-name>
   ```

2. **Check service logs directly:**
   ```bash
   docker-compose logs <service-name>
   ```

3. **Verify service outputs to stdout:**
   - Services must log to stdout/stderr
   - Check service logging configuration

4. **Check Promtail discovery:**
   ```bash
   docker-compose logs promtail | grep <service-name>
   ```

5. **Verify Docker labels:**
   - Check if service has correct Docker Compose labels
   - Review `promtail-config.yml` filters

## Diagnostic Commands

### Check Service Status
```bash
docker-compose ps
```

### Check Logs
```bash
# Loki logs
docker-compose logs loki

# Promtail logs
docker-compose logs promtail

# Grafana logs
docker-compose logs grafana
```

### Test Loki
```bash
# Health check
curl http://localhost:3100/ready

# Metrics
curl http://localhost:3100/metrics

# Query test
curl -G -s "http://localhost:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={job="docker"}' \
  --data-urlencode 'limit=10' | jq
```

### Test Promtail
```bash
# Check if Promtail is discovering containers
docker-compose logs promtail | grep "discovered"

# Check Promtail metrics
curl http://localhost:9080/metrics
```

### Check Network
```bash
# Inspect Docker network
docker network inspect woragis_default

# Test connectivity between containers
docker exec woragis-grafana ping loki
```

### Check Disk Usage
```bash
# Docker system usage
docker system df

# Volume usage
docker volume inspect woragis_loki-data
docker volume inspect woragis_grafana-data
```

## Getting More Help

1. **Check Documentation:**
   - `monitoring/README.md` - Overview
   - `monitoring/QUICK_START.md` - Quick start
   - `monitoring/LOGQL_QUERY_LIBRARY.md` - Query reference

2. **Review Logs:**
   - Check service logs for errors
   - Review Promtail discovery logs
   - Check Loki error logs

3. **Grafana Resources:**
   - [Grafana Loki Documentation](https://grafana.com/docs/loki/latest/)
   - [LogQL Documentation](https://grafana.com/docs/loki/latest/logql/)
   - [Grafana Troubleshooting](https://grafana.com/docs/grafana/latest/troubleshooting/)

4. **Community:**
   - Grafana Community Forums
   - Loki GitHub Issues

## Prevention

### Regular Maintenance

1. **Monitor disk usage:**
   - Set up alerts for disk usage
   - Clean up old logs regularly

2. **Review configurations:**
   - Check retention policies
   - Verify log levels
   - Review alert rules

3. **Test backups:**
   - Test log export procedures
   - Verify restoration works

4. **Update documentation:**
   - Document custom configurations
   - Update troubleshooting guide with new issues
