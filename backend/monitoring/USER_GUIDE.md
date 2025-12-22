# Logging Aggregation User Guide

**Version:** 1.0  
**Last Updated:** 2025-12-22

## Table of Contents

1. [Introduction](#introduction)
2. [Accessing Grafana](#accessing-grafana)
3. [Viewing Logs](#viewing-logs)
4. [Using Dashboards](#using-dashboards)
5. [Searching Logs](#searching-logs)
6. [Creating Custom Dashboards](#creating-custom-dashboards)
7. [Setting Up Alerts](#setting-up-alerts)
8. [Best Practices](#best-practices)

## Introduction

The Woragis logging aggregation system uses Grafana Loki to collect, store, and analyze logs from all services. This guide will help you use the system effectively.

### What You Can Do

- **Search logs** across all services
- **View dashboards** for system health and error analysis
- **Set up alerts** for critical issues
- **Trace requests** across services using trace IDs
- **Analyze patterns** in log data

## Accessing Grafana

### Initial Access

1. Open your browser and navigate to: `http://localhost:3000`
2. Login with:
   - **Username**: `admin`
   - **Password**: `admin` (change this immediately!)

### First-Time Setup

1. **Change Password**: Click your profile icon → Preferences → Change Password
2. **Explore Interface**: Familiarize yourself with the left sidebar
3. **Check Data Source**: Go to Configuration → Data Sources → Verify "Loki" is configured

## Viewing Logs

### Using Explore

1. Click **Explore** (compass icon) in the left sidebar
2. Select **Loki** as the data source
3. Enter a query (see [Searching Logs](#searching-logs))
4. Click **Run query**

### Basic Queries

**View all logs:**
```
{job="docker"}
```

**View logs from a specific service:**
```
{job="docker", service="app"}
```

**View only errors:**
```
{job="docker", level="error"}
```

**View errors from a specific service:**
```
{job="docker", service="app", level="error"}
```

### Viewing Log Details

- Click on any log line to see full details
- Expand JSON fields to see structured data
- Use the time range selector to filter by time

## Using Dashboards

### Available Dashboards

1. **Woragis Logs Overview**
   - All logs stream
   - Log volume by service
   - Error rates
   - Recent errors

2. **Service Health Overview**
   - Service status
   - Log volume trends
   - Error and warning rates
   - Service health table

3. **Error Analysis**
   - Recent errors
   - Errors by service
   - Error rate trends
   - Top error messages
   - Error details table

### Accessing Dashboards

1. Click **Dashboards** in the left sidebar
2. Browse or search for dashboards
3. Click on a dashboard to open it

### Dashboard Features

- **Time Range**: Use the time picker to change the time range
- **Refresh**: Click the refresh icon to update data
- **Panel Details**: Click on panels to see more details
- **Export**: Use the share icon to export or link to dashboards

## Searching Logs

### Label Filters

Use label filters for fast, efficient searches:

```
{job="docker", service="app", level="error"}
```

Available labels:
- `job`: Always "docker"
- `service`: Service name (app, ai-service, etc.)
- `level`: Log level (error, warn, info, debug)
- `trace_id`: Distributed tracing ID
- `request_id`: HTTP request ID

### Pattern Matching

Search for text in log messages:

```
{job="docker"} |= "database"
{job="docker"} |= "connection timeout"
```

Use regex for complex patterns:

```
{job="docker"} |~ "error|exception|failure"
```

### Time Ranges

Specify time ranges for queries:

```
{job="docker"} [5m]    # Last 5 minutes
{job="docker"} [1h]    # Last hour
{job="docker"} [24h]   # Last 24 hours
```

### Combining Filters

Combine multiple filters:

```
{job="docker", service="app", level="error"} |= "database" [1h]
```

## Creating Custom Dashboards

### Step 1: Create New Dashboard

1. Click **Dashboards** → **New Dashboard**
2. Click **Add visualization**
3. Select **Loki** as the data source

### Step 2: Add Panels

**Logs Panel:**
1. Select **Logs** visualization type
2. Enter LogQL query
3. Configure options (show labels, time, etc.)

**Time Series Panel:**
1. Select **Time series** visualization type
2. Enter LogQL query with aggregation
3. Configure axes and legend

**Stat Panel:**
1. Select **Stat** visualization type
2. Enter LogQL query
3. Configure thresholds and colors

### Step 3: Save Dashboard

1. Click **Save dashboard**
2. Enter dashboard name
3. Add tags (optional)
4. Choose folder
5. Click **Save**

### Example: Service-Specific Dashboard

Create a dashboard for a specific service:

1. **Panel 1**: Logs from service
   ```
   {job="docker", service="resume-worker"}
   ```

2. **Panel 2**: Error rate
   ```
   sum(rate({job="docker", service="resume-worker", level="error"}[5m]))
   ```

3. **Panel 3**: Recent errors
   ```
   {job="docker", service="resume-worker", level="error"}
   ```

## Setting Up Alerts

### Using Pre-configured Alerts

Alerts are configured in `monitoring/grafana/provisioning/alerting/rules.yml`:

- **High Error Rate**: Triggers when error rate > 10 errors/sec
- **Service Down**: Triggers when service stops logging
- **Critical Error Pattern**: Triggers on database/connection errors

### Creating Custom Alerts

1. Go to **Alerting** → **Alert rules**
2. Click **New alert rule**
3. Configure:
   - **Name**: Alert name
   - **Condition**: LogQL query
   - **Threshold**: When to trigger
   - **Evaluation**: How often to check
4. Add **Notification channels**
5. Save alert rule

### Alert Examples

**High Error Rate:**
```
sum(rate({job="docker", level="error"}[5m])) > 10
```

**Service Not Logging:**
```
count(count_over_time({job="docker", service="app"}[5m]) > 0) == 0
```

**Specific Error Pattern:**
```
sum(count_over_time({job="docker", level="error"} |~ "database.*connection" [5m])) > 5
```

## Best Practices

### 1. Use Label Filters First

Label filters are faster than pattern matching:
```
✅ {job="docker", service="app", level="error"}
❌ {job="docker"} |= "app" |= "error"
```

### 2. Specify Time Ranges

Always specify appropriate time ranges:
```
✅ {job="docker"} [1h]
❌ {job="docker"}  # May be too slow
```

### 3. Use Aggregations for Metrics

Use rate() and count_over_time() for metrics:
```
✅ sum(rate({job="docker", level="error"}[5m]))
❌ count({job="docker", level="error"})  # Less efficient
```

### 4. Trace Requests Across Services

Use trace_id to follow requests:
```
{job="docker", trace_id="550e8400-e29b-41d4-a716-446655440000"}
```

### 5. Save Common Queries

Save frequently used queries:
1. In Explore, enter your query
2. Click **Add to dashboard** or save as favorite

### 6. Use Dashboards for Monitoring

Create dashboards for:
- Service health
- Error monitoring
- Performance metrics
- Business metrics

## Common Tasks

### Find All Errors in Last Hour

```
{job="docker", level="error"} [1h]
```

### Check Service Health

```
sum(rate({job="docker", service="app"}[5m]))
```

### Find Slow Requests

```
{job="docker", service="app"} | json | duration_ms > 1000
```

### Track a Request

```
{job="docker", trace_id="YOUR-TRACE-ID"}
```

### Analyze Error Patterns

```
topk(10, sum(count_over_time({job="docker", level="error"}[1h])) by (message))
```

## Tips and Tricks

1. **Use Log Details**: Click on log lines to see full structured data
2. **Copy Query**: Right-click on panels to copy LogQL queries
3. **Share Dashboards**: Use the share icon to share dashboards with team
4. **Export Data**: Use Explore to export log data for analysis
5. **Bookmark Queries**: Save common queries for quick access

## Getting Help

- **Documentation**: See `monitoring/README.md` and `monitoring/LOGQL_QUERY_LIBRARY.md`
- **LogQL Reference**: [Grafana LogQL Documentation](https://grafana.com/docs/loki/latest/logql/)
- **Grafana Docs**: [Grafana Documentation](https://grafana.com/docs/grafana/latest/)

## Next Steps

1. Explore the pre-configured dashboards
2. Try some queries in Explore
3. Create a custom dashboard for your service
4. Set up alerts for critical issues
5. Share useful queries with your team
