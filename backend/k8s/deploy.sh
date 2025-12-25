#!/bin/bash
# Kubernetes Deployment Script
# 
# This script deploys all Woragis services to Kubernetes in the correct order.

set -e

NAMESPACE="${NAMESPACE:-default}"

echo "🚀 Deploying Woragis services to Kubernetes (namespace: $NAMESPACE)"

# Check if namespace exists, create if not
if ! kubectl get namespace "$NAMESPACE" &> /dev/null; then
    echo "📦 Creating namespace: $NAMESPACE"
    kubectl create namespace "$NAMESPACE"
fi

# Function to deploy a service
deploy_service() {
    local service=$1
    echo "📦 Deploying $service..."
    kubectl apply -f "$service/" -n "$NAMESPACE"
    echo "✅ $service deployed"
}

# Deploy services first (dependencies)
echo ""
echo "=== Deploying Services ==="
deploy_service "ai-service"
deploy_service "creative-service"
deploy_service "docs-service"

# Deploy workers
echo ""
echo "=== Deploying Workers ==="
deploy_service "email-worker"
deploy_service "translation-worker"
deploy_service "job-application-worker"
deploy_service "resume-worker"
deploy_service "whatsapp-worker"

# Deploy server last (depends on services)
echo ""
echo "=== Deploying Server ==="
deploy_service "server"

echo ""
echo "✅ All services deployed!"
echo ""
echo "📊 Check status with:"
echo "   kubectl get pods -n $NAMESPACE"
echo "   kubectl get services -n $NAMESPACE"
echo "   kubectl get deployments -n $NAMESPACE"

