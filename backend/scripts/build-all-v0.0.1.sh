#!/bin/bash
# Build and tag all services as v0.0.1
# This script builds all Docker images locally for testing

set -e

VERSION="v0.0.1"
DOCKER_HUB_USERNAME="${DOCKER_HUB_USERNAME:-woragis}"

echo "=========================================="
echo "Building all services as $VERSION"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to build a service
build_service() {
    local service_name=$1
    local image_name=$2
    local dockerfile_path=$3
    local context_path=$4
    
    echo -e "${BLUE}Building $service_name...${NC}"
    echo "  Image: $image_name:$VERSION"
    echo "  Dockerfile: $dockerfile_path"
    echo "  Context: $context_path"
    
    docker build \
        -f "$dockerfile_path" \
        -t "$DOCKER_HUB_USERNAME/$image_name:$VERSION" \
        -t "$DOCKER_HUB_USERNAME/$image_name:latest" \
        "$context_path"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ $service_name built successfully${NC}"
    else
        echo -e "✗ $service_name build failed"
        exit 1
    fi
    echo ""
}

# Build all services
echo "Building server..."
build_service "server" "woragis-server" "./backend/server/Dockerfile" "./backend/server"

echo "Building email-worker..."
build_service "email-worker" "woragis-email-worker" "./backend/Dockerfile.email-worker" "./backend/email-worker"

echo "Building translation-worker..."
build_service "translation-worker" "woragis-translation-worker" "./backend/Dockerfile.translation-worker" "./backend/translation-worker"

echo "Building whatsapp-worker..."
build_service "whatsapp-worker" "woragis-whatsapp-worker" "./backend/Dockerfile.whatsapp-worker" "./backend/whatsapp-worker"

echo "Building job-application-worker..."
build_service "job-application-worker" "woragis-job-application-worker" "./backend/Dockerfile.job-application-worker" "./backend/job-application-worker"

echo "Building resume-worker..."
build_service "resume-worker" "woragis-resume-worker" "./backend/Dockerfile.resume-worker" "./backend/resume-worker"

echo "Building ai-service..."
build_service "ai-service" "woragis-ai-service" "./backend/ai-service/Dockerfile" "./backend/ai-service"

echo "Building creative-service..."
build_service "creative-service" "woragis-creative-service" "./backend/creative-service/Dockerfile" "./backend/creative-service"

echo "Building docs-service..."
build_service "docs-service" "woragis-docs-service" "./backend/Dockerfile.docs-service" "./backend"

echo "=========================================="
echo -e "${GREEN}All services built successfully!${NC}"
echo "=========================================="
echo ""
echo "Built images:"
echo "  - $DOCKER_HUB_USERNAME/woragis-server:$VERSION"
echo "  - $DOCKER_HUB_USERNAME/woragis-email-worker:$VERSION"
echo "  - $DOCKER_HUB_USERNAME/woragis-translation-worker:$VERSION"
echo "  - $DOCKER_HUB_USERNAME/woragis-whatsapp-worker:$VERSION"
echo "  - $DOCKER_HUB_USERNAME/woragis-job-application-worker:$VERSION"
echo "  - $DOCKER_HUB_USERNAME/woragis-resume-worker:$VERSION"
echo "  - $DOCKER_HUB_USERNAME/woragis-ai-service:$VERSION"
echo "  - $DOCKER_HUB_USERNAME/woragis-creative-service:$VERSION"
echo "  - $DOCKER_HUB_USERNAME/woragis-docs-service:$VERSION"
echo ""
echo "To push to Docker Hub, run:"
echo "  docker push $DOCKER_HUB_USERNAME/woragis-server:$VERSION"
echo "  docker push $DOCKER_HUB_USERNAME/woragis-email-worker:$VERSION"
echo "  docker push $DOCKER_HUB_USERNAME/woragis-translation-worker:$VERSION"
echo "  docker push $DOCKER_HUB_USERNAME/woragis-whatsapp-worker:$VERSION"
echo "  docker push $DOCKER_HUB_USERNAME/woragis-job-application-worker:$VERSION"
echo "  docker push $DOCKER_HUB_USERNAME/woragis-resume-worker:$VERSION"
echo "  docker push $DOCKER_HUB_USERNAME/woragis-ai-service:$VERSION"
echo "  docker push $DOCKER_HUB_USERNAME/woragis-creative-service:$VERSION"
echo "  docker push $DOCKER_HUB_USERNAME/woragis-docs-service:$VERSION"


















