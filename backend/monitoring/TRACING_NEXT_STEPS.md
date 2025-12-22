# Distributed Tracing - Next Steps

## What's Done ✅

1. **Infrastructure**: Jaeger service added to docker-compose.yml
2. **Go Services**: Tracing package created, middleware added to main server
3. **Python Services**: Tracing modules created for all Python services
4. **Node.js Service**: Tracing module created
5. **Grafana Integration**: Jaeger data source configured
6. **Documentation**: Complete guides created

## What Needs to Be Done

### 1. Add Dependencies

#### Go Services (server)
Run in `server/` directory:
```bash
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/trace
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp
go get go.opentelemetry.io/otel/propagation
go get go.opentelemetry.io/otel/sdk/resource
go get go.opentelemetry.io/otel/sdk/trace
go get go.opentelemetry.io/otel/semconv/v1.24.0
```

#### Python Services
Add to `requirements.txt` for each Python service:
```python
opentelemetry-api
opentelemetry-sdk
opentelemetry-exporter-otlp-proto-http
opentelemetry-instrumentation-fastapi
opentelemetry-instrumentation-requests
```

Then install:
```bash
pip install -r requirements.txt
```

#### Node.js Service
Add to `package.json`:
```json
"@opentelemetry/api": "^1.8.0",
"@opentelemetry/sdk-node": "^0.51.0",
"@opentelemetry/exporter-otlp-http": "^0.51.0",
"@opentelemetry/instrumentation-http": "^0.51.0",
"@opentelemetry/instrumentation-amqplib": "^0.40.0"
```

Then install:
```bash
npm install
```

### 2. Test Tracing

1. **Start Jaeger:**
   ```bash
   docker-compose up -d jaeger
   ```

2. **Restart services:**
   ```bash
   docker-compose restart app ai-service creative-service docs-service
   ```

3. **Make a request:**
   ```bash
   curl http://localhost:8080/api/health
   ```

4. **Check Jaeger:**
   - Go to http://localhost:16686
   - Select service: "woragis-server"
   - Click "Find Traces"
   - You should see the trace

### 3. Verify Trace Propagation

1. **Make a request that calls multiple services:**
   ```bash
   curl http://localhost:8080/api/projects
   ```

2. **Check trace in Jaeger:**
   - Should see spans from multiple services
   - Trace should span across service boundaries

### 4. Test Trace-to-Log Correlation

1. **Get trace ID from Jaeger**
2. **Search logs in Grafana:**
   ```
   {job="docker", trace_id="YOUR-TRACE-ID"}
   ```
3. **Verify logs match trace**

## Known Issues to Fix

1. **Queue Trace Propagation**: Workers need to extract trace context from queue message headers
2. **Database Instrumentation**: Add GORM and SQLAlchemy instrumentation
3. **Worker Instrumentation**: Complete queue consumer trace context extraction

## Testing Checklist

- [ ] Jaeger UI accessible
- [ ] Traces appear in Jaeger
- [ ] Trace spans multiple services
- [ ] Trace IDs match in logs
- [ ] Grafana can query Jaeger
- [ ] Trace-to-log links work

## After Testing

Once tracing is verified:
1. Document any issues found
2. Complete queue trace propagation
3. Add database instrumentation
4. Create trace visualization dashboard
5. Set up trace retention policies
