#!/bin/bash

# Seed script for Technical Writings
# This script creates various technical writing portfolio entries

BASE_URL="${API_BASE_URL:-http://localhost:8080/api}"
AUTH_TOKEN="${AUTH_TOKEN:-}"

if [ -z "$AUTH_TOKEN" ]; then
    echo "Error: AUTH_TOKEN environment variable is not set"
    echo "Please set AUTH_TOKEN to a valid JWT token"
    exit 1
fi

echo "Seeding Technical Writings..."

# Helper function to create a writing
create_writing() {
    local title=$1
    local description=$2
    local type=$3
    local platform=$4
    local url=$5
    local excerpt=$6
    local published_at=$7
    local reading_time=$8
    local topics=$9
    local technologies=${10}
    local views=${11}
    local likes=${12}
    local featured=${13}
    local display_order=${14}

    local payload="{
        \"title\": \"$title\",
        \"description\": \"$description\",
        \"type\": \"$type\",
        \"platform\": \"$platform\",
        \"url\": \"$url\""

    if [ -n "$excerpt" ]; then
        payload="$payload,
        \"excerpt\": \"$excerpt\""
    fi

    if [ -n "$published_at" ]; then
        payload="$payload,
        \"publishedAt\": \"$published_at\""
    fi

    if [ -n "$reading_time" ]; then
        payload="$payload,
        \"readingTime\": $reading_time"
    fi

    if [ -n "$topics" ]; then
        payload="$payload,
        \"topics\": [$topics]"
    fi

    if [ -n "$technologies" ]; then
        payload="$payload,
        \"technologies\": [$technologies]"
    fi

    if [ -n "$views" ]; then
        payload="$payload,
        \"views\": $views"
    fi

    if [ -n "$likes" ]; then
        payload="$payload,
        \"likes\": $likes"
    fi

    if [ "$featured" = "true" ]; then
        payload="$payload,
        \"featured\": true"
    fi

    if [ -n "$display_order" ]; then
        payload="$payload,
        \"displayOrder\": $display_order"
    fi

    payload="$payload
    }"

    response=$(curl -s -X POST "$BASE_URL/technical-writings" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $AUTH_TOKEN" \
        -d "$payload")

    if echo "$response" | grep -q '"id"'; then
        echo "✓ Created $type: $title"
    else
        echo "✗ Failed to create $type: $response"
    fi
}

# Technical Articles
echo ""
echo "Creating Technical Articles..."
create_writing \
    "Building Scalable Microservices with Go and Kubernetes" \
    "A comprehensive guide to designing and deploying microservices architecture using Go and Kubernetes. Covers service discovery, load balancing, and monitoring strategies." \
    "article" \
    "medium" \
    "https://medium.com/@example/microservices-go-kubernetes" \
    "Learn how to build production-ready microservices with Go, covering everything from architecture design to deployment on Kubernetes." \
    "2024-03-15" \
    "12" \
    "\"Microservices\", \"Go\", \"Kubernetes\", \"Architecture\", \"DevOps\"" \
    "\"Go\", \"Kubernetes\", \"Docker\", \"gRPC\", \"Prometheus\"" \
    "15000" \
    "450" \
    "true" \
    "1"

create_writing \
    "Mastering TypeScript: Advanced Patterns and Best Practices" \
    "Deep dive into TypeScript's advanced features including conditional types, mapped types, and utility types. Real-world examples and patterns for enterprise applications." \
    "article" \
    "dev_to" \
    "https://dev.to/example/typescript-advanced-patterns" \
    "Explore TypeScript's most powerful features and learn how to write type-safe, maintainable code at scale." \
    "2024-02-20" \
    "15" \
    "\"TypeScript\", \"Programming\", \"Best Practices\", \"Type Safety\"" \
    "\"TypeScript\", \"JavaScript\", \"Node.js\"" \
    "22000" \
    "680" \
    "true" \
    "2"

# Documentation
echo ""
echo "Creating Documentation..."
create_writing \
    "API Documentation: RESTful Design Principles" \
    "Complete guide to writing comprehensive API documentation. Includes OpenAPI specifications, code examples, and best practices for developer experience." \
    "documentation" \
    "github" \
    "https://github.com/example/api-docs" \
    "Learn how to create documentation that developers love to use." \
    "2024-01-10" \
    "8" \
    "\"API\", \"Documentation\", \"REST\", \"OpenAPI\", \"Developer Experience\"" \
    "\"OpenAPI\", \"Swagger\", \"Markdown\", \"Postman\"" \
    "8500" \
    "230" \
    "true" \
    "3"

create_writing \
    "System Architecture Documentation Template" \
    "A comprehensive template for documenting system architecture, including diagrams, component descriptions, and deployment strategies." \
    "documentation" \
    "personal_blog" \
    "https://blog.example.com/architecture-docs" \
    "Standardize your architecture documentation with this proven template." \
    "2024-04-05" \
    "10" \
    "\"Architecture\", \"Documentation\", \"System Design\", \"Best Practices\"" \
    "\"C4 Model\", \"PlantUML\", \"Mermaid\"" \
    "6200" \
    "180" \
    "false" \
    "0"

# Tutorials
echo ""
echo "Creating Tutorials..."
create_writing \
    "Complete Guide to Building a Full-Stack Application with Next.js and PostgreSQL" \
    "Step-by-step tutorial on building a modern full-stack application from scratch. Covers authentication, database design, API development, and deployment." \
    "tutorial" \
    "hashnode" \
    "https://example.hashnode.dev/nextjs-postgresql-tutorial" \
    "Build a production-ready full-stack application with Next.js 14, PostgreSQL, and Prisma." \
    "2024-05-12" \
    "25" \
    "\"Tutorial\", \"Next.js\", \"PostgreSQL\", \"Full-Stack\", \"Web Development\"" \
    "\"Next.js\", \"React\", \"PostgreSQL\", \"Prisma\", \"TypeScript\", \"Tailwind CSS\"" \
    "18000" \
    "520" \
    "true" \
    "4"

create_writing \
    "Getting Started with Docker: From Zero to Production" \
    "A beginner-friendly tutorial on Docker, covering containerization basics, Docker Compose, multi-stage builds, and production deployment strategies." \
    "tutorial" \
    "medium" \
    "https://medium.com/@example/docker-tutorial" \
    "Master Docker and containerization with this comprehensive tutorial for developers." \
    "2024-03-28" \
    "18" \
    "\"Docker\", \"DevOps\", \"Containers\", \"Tutorial\", \"Deployment\"" \
    "\"Docker\", \"Docker Compose\", \"Kubernetes\"" \
    "12500" \
    "380" \
    "true" \
    "5"

# Guides
echo ""
echo "Creating Guides..."
create_writing \
    "The Complete Guide to CI/CD Pipelines with GitHub Actions" \
    "Comprehensive guide to setting up CI/CD pipelines using GitHub Actions. Includes testing, building, deploying, and monitoring workflows." \
    "guide" \
    "dev_to" \
    "https://dev.to/example/github-actions-cicd" \
    "Automate your development workflow with GitHub Actions and deploy with confidence." \
    "2024-04-18" \
    "20" \
    "\"CI/CD\", \"GitHub Actions\", \"DevOps\", \"Automation\", \"Deployment\"" \
    "\"GitHub Actions\", \"Docker\", \"Kubernetes\", \"AWS\", \"Terraform\"" \
    "14500" \
    "420" \
    "true" \
    "6"

create_writing \
    "Security Best Practices for Web Applications" \
    "A practical guide to securing web applications, covering authentication, authorization, data encryption, and common vulnerabilities." \
    "guide" \
    "company_blog" \
    "https://blog.company.com/security-best-practices" \
    "Protect your web applications from common security threats with these proven practices." \
    "2024-02-14" \
    "14" \
    "\"Security\", \"Web Development\", \"Best Practices\", \"Authentication\", \"Encryption\"" \
    "\"JWT\", \"OAuth\", \"HTTPS\", \"OWASP\"" \
    "9800" \
    "290" \
    "false" \
    "0"

# Blog Posts
echo ""
echo "Creating Blog Posts..."
create_writing \
    "My Journey from Junior to Senior Developer: Lessons Learned" \
    "Personal reflections on career growth, technical skills development, and the mindset shifts needed to advance as a software engineer." \
    "blog_post" \
    "medium" \
    "https://medium.com/@example/junior-to-senior" \
    "Insights and lessons from my journey to becoming a senior developer." \
    "2024-06-01" \
    "8" \
    "\"Career\", \"Personal Development\", \"Software Engineering\", \"Growth\"" \
    "" \
    "11200" \
    "340" \
    "false" \
    "0"

create_writing \
    "Why I Switched from REST to GraphQL (And When You Should Too)" \
    "A detailed comparison of REST and GraphQL, discussing when to use each, migration strategies, and real-world performance implications." \
    "blog_post" \
    "hashnode" \
    "https://example.hashnode.dev/rest-vs-graphql" \
    "Understanding when GraphQL makes sense and when REST is still the better choice." \
    "2024-05-22" \
    "11" \
    "\"GraphQL\", \"REST\", \"API Design\", \"Backend Development\"" \
    "\"GraphQL\", \"Apollo\", \"Node.js\", \"Express\"" \
    "8900" \
    "260" \
    "false" \
    "0"

# Case Studies
echo ""
echo "Creating Case Studies..."
create_writing \
    "Case Study: Scaling a SaaS Platform from 0 to 100K Users" \
    "Detailed case study on the technical challenges and solutions involved in scaling a SaaS platform, including database optimization, caching strategies, and infrastructure decisions." \
    "case_study" \
    "personal_blog" \
    "https://blog.example.com/scaling-saas-case-study" \
    "Real-world insights from scaling a SaaS platform to handle 100K+ users." \
    "2024-04-30" \
    "22" \
    "\"Case Study\", \"Scaling\", \"SaaS\", \"Performance\", \"Infrastructure\"" \
    "\"PostgreSQL\", \"Redis\", \"AWS\", \"Kubernetes\", \"Node.js\"" \
    "16800" \
    "490" \
    "true" \
    "7"

echo ""
echo "Technical Writings seeding completed!"
echo ""
echo "You can now view the portfolio at: GET $BASE_URL/technical-writings/featured"
echo "Filter by type: GET $BASE_URL/technical-writings/type/article"
echo "Filter by platform: GET $BASE_URL/technical-writings/platform/medium"
echo "Search: GET $BASE_URL/technical-writings/search?q=typescript"

