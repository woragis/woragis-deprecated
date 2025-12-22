# Distributed Tracing Quick Start

## Start Jaeger

```bash
docker-compose up -d jaeger
```

## Access Jaeger UI

- URL: http://localhost:16686
- Search for traces by service, operation, or trace ID

## Verify Tracing

1. **Check Jaeger is running:**
   ```bash
   docker-compose ps jaeger
   ```

2. **Make a request to your API:**
   ```bash
   curl http://localhost:8080/api/health
   ```

3. **View trace in Jaeger:**
   - Go to http://localhost:16686
   - Select service: "woragis-server"
   - Click "Find Traces"
   - You should see the request trace

## Trace Context Propagation

Traces automatically propagate via:
- HTTP headers (`traceparent`, `X-Trace-ID`)
- Queue message headers (when implemented)
- Context propagation in code

## Viewing Traces

### In Jaeger UI
1. Open http://localhost:16686
2. Select service from dropdown
3. Click "Find Traces"
4. Click on a trace to see details

### In Grafana (when configured)
1. Open Grafana
2. Go to Explore
3. Select Jaeger data source
4. Search for traces

## Common Issues

### No Traces Appearing

1. **Check Jaeger is running:**
   ```bash
   docker-compose logs jaeger
   ```

2. **Check service is sending traces:**
   - Verify OTLP endpoint is correct
   - Check service logs for tracing errors

3. **Check sampling:**
   - Development: 100% sampling
   - Production: 10% sampling (may miss some traces)

### Trace ID Mismatch

- Ensure trace_id from OpenTelemetry matches log trace_id
- Check middleware order (tracing before logging)

## Next Steps

- Configure Grafana Jaeger data source
- Create trace visualization dashboard
- Link traces to logs and metrics
