# Prometheus Metrics Quick Start

## Start Services

```bash
docker-compose up -d prometheus grafana
```

## Verify Prometheus

1. **Check Prometheus is running:**
   ```bash
   docker-compose ps prometheus
   ```

2. **Access Prometheus UI:**
   - http://localhost:9090

3. **Check targets:**
   - http://localhost:9090/targets
   - All services should show as "UP"

4. **Test a query:**
   - Go to http://localhost:9090/graph
   - Try: `up`

## Verify Grafana

1. **Access Grafana:**
   - http://localhost:3000
   - Login: admin/admin

2. **Check Prometheus data source:**
   - Configuration → Data Sources → Prometheus
   - Should show as "Connected"

3. **View dashboards:**
   - Dashboards → System Overview
   - Dashboards → Queue Monitoring

## Test Metrics

1. **Query a service metric:**
   ```bash
   curl http://localhost:8080/metrics | grep http_requests
   ```

2. **Check in Prometheus:**
   - Query: `http_requests_total`
   - Should see metrics from services

3. **Check in Grafana:**
   - Go to Explore
   - Select Prometheus
   - Query: `up`

## Common Issues

### Services Not Scraping

1. Check service is running
2. Verify `/metrics` endpoint works
3. Check Prometheus targets page
4. Review Prometheus logs

### No Data in Grafana

1. Verify Prometheus data source is connected
2. Check Prometheus has data (query in Prometheus UI)
3. Verify time range in Grafana
4. Check dashboard queries are correct

## Next Steps

- Review dashboards
- Set up custom alerts
- Create service-specific dashboards
- Monitor metrics over time
