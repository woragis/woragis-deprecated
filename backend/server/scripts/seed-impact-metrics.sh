#!/bin/bash

# Seed script for Impact Metrics
# This script creates various impact metrics to populate the dashboard

BASE_URL="${API_BASE_URL:-http://localhost:8080/api}"
AUTH_TOKEN="${AUTH_TOKEN:-}"

if [ -z "$AUTH_TOKEN" ]; then
    echo "Error: AUTH_TOKEN environment variable is not set"
    echo "Please set AUTH_TOKEN to a valid JWT token"
    exit 1
fi

echo "Seeding Impact Metrics..."

# Helper function to create a metric
create_metric() {
    local type=$1
    local value=$2
    local unit=$3
    local description=$4
    local entity_type=$5
    local entity_id=$6
    local featured=$7
    local display_order=$8
    local period_start=$9
    local period_end=${10}

    local payload="{
        \"type\": \"$type\",
        \"value\": $value,
        \"unit\": \"$unit\",
        \"description\": \"$description\""

    if [ -n "$entity_type" ] && [ -n "$entity_id" ]; then
        payload="$payload,
        \"entityType\": \"$entity_type\",
        \"entityId\": \"$entity_id\""
    fi

    if [ "$featured" = "true" ]; then
        payload="$payload,
        \"featured\": true"
    fi

    if [ -n "$display_order" ]; then
        payload="$payload,
        \"displayOrder\": $display_order"
    fi

    if [ -n "$period_start" ]; then
        payload="$payload,
        \"periodStart\": \"$period_start\""
    fi

    if [ -n "$period_end" ]; then
        payload="$payload,
        \"periodEnd\": \"$period_end\""
    fi

    payload="$payload
    }"

    response=$(curl -s -X POST "$BASE_URL/impact-metrics" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $AUTH_TOKEN" \
        -d "$payload")

    if echo "$response" | grep -q '"id"'; then
        echo "✓ Created $type metric: $description"
    else
        echo "✗ Failed to create $type metric: $response"
    fi
}

# Projects Delivered Metrics
echo ""
echo "Creating Projects Delivered metrics..."
create_metric "projects_delivered" 15 "count" "Total projects successfully delivered in 2024" "" "" "true" 1 "2024-01-01" "2024-12-31"
create_metric "projects_delivered" 8 "count" "Projects delivered in Q1 2024" "" "" "false" 0 "2024-01-01" "2024-03-31"
create_metric "projects_delivered" 7 "count" "Projects delivered in Q2 2024" "" "" "false" 0 "2024-04-01" "2024-06-30"

# Users Impacted Metrics
echo ""
echo "Creating Users Impacted metrics..."
create_metric "users_impacted" 50000 "count" "Total users impacted by all projects" "" "" "true" 2 "" ""
create_metric "users_impacted" 25000 "count" "Users impacted by e-commerce platform migration" "project" "" "false" 0 "" ""
create_metric "users_impacted" 15000 "count" "Users impacted by mobile app launch" "project" "" "false" 0 "" ""
create_metric "users_impacted" 10000 "count" "Users impacted by API optimization project" "project" "" "false" 0 "" ""

# Performance Improvements Metrics
echo ""
echo "Creating Performance Improvements metrics..."
create_metric "performance_improvement" 75 "percentage" "API response time improvement (from 800ms to 200ms)" "problem_solution" "" "true" 3 "" ""
create_metric "performance_improvement" 60 "percentage" "Database query optimization (reduced query time by 60%)" "problem_solution" "" "false" 0 "" ""
create_metric "performance_improvement" 85 "percentage" "Frontend load time reduction (from 5s to 0.75s)" "problem_solution" "" "false" 0 "" ""
create_metric "performance_improvement" 50 "percentage" "Image processing pipeline optimization" "problem_solution" "" "false" 0 "" ""
create_metric "performance_improvement" 90 "percentage" "Caching layer implementation (90% cache hit rate)" "system_design" "" "false" 0 "" ""

# Cost Savings Metrics
echo ""
echo "Creating Cost Savings metrics..."
create_metric "cost_savings" 45000 "currency" "Annual cloud infrastructure cost savings through optimization" "" "" "true" 4 "2024-01-01" "2024-12-31"
create_metric "cost_savings" 12000 "currency" "Q1 2024 cost savings from database optimization" "" "" "false" 0 "2024-01-01" "2024-03-31"
create_metric "cost_savings" 15000 "currency" "Q2 2024 cost savings from auto-scaling implementation" "" "" "false" 0 "2024-04-01" "2024-06-30"
create_metric "cost_savings" 18000 "currency" "Q3 2024 cost savings from CDN optimization" "" "" "false" 0 "2024-07-01" "2024-09-30"
create_metric "cost_savings" 25000 "currency" "Cost savings from migrating to serverless architecture" "project" "" "false" 0 "" ""

# Time Saved Through Automation Metrics
echo ""
echo "Creating Time Saved metrics..."
create_metric "time_saved" 1200 "hours" "Total time saved through automation in 2024" "" "" "true" 5 "2024-01-01" "2024-12-31"
create_metric "time_saved" 300 "hours" "Time saved through CI/CD pipeline automation" "" "" "false" 0 "" ""
create_metric "time_saved" 450 "hours" "Time saved through automated testing and deployment" "" "" "false" 0 "" ""
create_metric "time_saved" 250 "hours" "Time saved through automated data processing workflows" "" "" "false" 0 "" ""
create_metric "time_saved" 200 "hours" "Time saved through automated monitoring and alerting" "" "" "false" 0 "" ""
create_metric "time_saved" 8 "hours" "Daily time saved through automated report generation" "" "" "false" 0 "" ""

# Additional featured metrics for dashboard
echo ""
echo "Creating additional featured metrics..."
create_metric "users_impacted" 5000 "count" "Active users on new mobile platform" "project" "" "true" 6 "" ""
create_metric "performance_improvement" 40 "percentage" "Overall system performance improvement" "" "" "true" 7 "" ""
create_metric "cost_savings" 8000 "currency" "Monthly recurring cost savings" "" "" "true" 8 "" ""
create_metric "time_saved" 40 "hours" "Weekly time saved through automation" "" "" "true" 9 "" ""

echo ""
echo "Impact Metrics seeding completed!"
echo ""
echo "You can now view the dashboard at: GET $BASE_URL/impact-metrics/dashboard"
echo "Featured metrics: GET $BASE_URL/impact-metrics/featured"

