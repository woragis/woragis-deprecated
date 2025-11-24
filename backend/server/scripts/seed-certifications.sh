#!/bin/bash

# Script to seed certifications using curl
# Usage: ./seed-certifications.sh [email] [password]

API_BASE="http://localhost:8080/api"
EMAIL="${1:-masteringthecode.woragis@gmail.com}"
PASSWORD="${2}"

if [ -z "$PASSWORD" ]; then
    echo "Usage: $0 [email] [password]"
    echo "Please provide your password to authenticate"
    exit 1
fi

echo "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

# Extract access_token from JSON response
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
if [ -z "$ACCESS_TOKEN" ]; then
    ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"accessToken":"[^"]*' | cut -d'"' -f4)
fi

if [ -z "$ACCESS_TOKEN" ]; then
    echo "Failed to login. Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "Login successful! Token: ${ACCESS_TOKEN:0:20}..."

# Create certifications
echo "Creating certifications..."

# Certification 1: AWS Solutions Architect
curl -s -X POST "$API_BASE/certifications" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "AWS Certified Solutions Architect - Associate",
        "issuer": "Amazon Web Services",
        "issueDate": "2023-06-15",
        "expiryDate": "2026-06-15",
        "credentialId": "AWS-CSA-123456789",
        "verificationUrl": "https://www.credly.com/badges/abc123",
        "certificateUrl": "https://example.com/certificates/aws-saa.pdf",
        "description": "Validates expertise in designing distributed systems on AWS. Covers architecture best practices, cost optimization, security, and scalability.",
        "status": "active",
        "category": "cloud",
        "featured": true,
        "displayOrder": 1,
        "skillIds": []
    }'

echo ""
echo "Certification 1 created"

# Certification 2: Google Cloud Professional Cloud Architect
curl -s -X POST "$API_BASE/certifications" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "Google Cloud Professional Cloud Architect",
        "issuer": "Google Cloud",
        "issueDate": "2023-09-20",
        "expiryDate": "2025-09-20",
        "credentialId": "GCP-PCA-987654321",
        "verificationUrl": "https://www.credential.net/xyz789",
        "certificateUrl": "https://example.com/certificates/gcp-pca.pdf",
        "description": "Demonstrates ability to design, develop, and manage robust, secure, scalable, highly available, and dynamic solutions on Google Cloud Platform.",
        "status": "active",
        "category": "cloud",
        "featured": true,
        "displayOrder": 2,
        "skillIds": []
    }'

echo ""
echo "Certification 2 created"

# Certification 3: Kubernetes Administrator (CKA)
curl -s -X POST "$API_BASE/certifications" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "Certified Kubernetes Administrator (CKA)",
        "issuer": "Cloud Native Computing Foundation",
        "issueDate": "2023-03-10",
        "expiryDate": "2026-03-10",
        "credentialId": "CKA-456789012",
        "verificationUrl": "https://www.credly.com/badges/def456",
        "certificateUrl": "https://example.com/certificates/cka.pdf",
        "description": "Validates skills in Kubernetes administration including cluster architecture, installation, configuration, networking, storage, security, and troubleshooting.",
        "status": "active",
        "category": "devops",
        "featured": true,
        "displayOrder": 3,
        "skillIds": []
    }'

echo ""
echo "Certification 3 created"

# Certification 4: HashiCorp Terraform Associate
curl -s -X POST "$API_BASE/certifications" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "HashiCorp Certified: Terraform Associate",
        "issuer": "HashiCorp",
        "issueDate": "2023-11-05",
        "expiryDate": "2025-11-05",
        "credentialId": "HCTA-789012345",
        "verificationUrl": "https://www.credly.com/badges/ghi789",
        "certificateUrl": "https://example.com/certificates/terraform-associate.pdf",
        "description": "Validates knowledge of infrastructure automation using Terraform. Covers core concepts, state management, modules, and best practices.",
        "status": "active",
        "category": "devops",
        "featured": false,
        "displayOrder": 4,
        "skillIds": []
    }'

echo ""
echo "Certification 4 created"

# Certification 5: PostgreSQL Certified Professional
curl -s -X POST "$API_BASE/certifications" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "PostgreSQL Certified Professional",
        "issuer": "PostgreSQL Global Development Group",
        "issueDate": "2022-08-12",
        "expiryDate": null,
        "credentialId": "PGCP-2022-001234",
        "verificationUrl": "https://www.postgresql.org/certification/verify/001234",
        "certificateUrl": "https://example.com/certificates/postgresql-professional.pdf",
        "description": "Validates advanced knowledge of PostgreSQL database administration, optimization, replication, and high availability configurations.",
        "status": "active",
        "category": "database",
        "featured": false,
        "displayOrder": 5,
        "skillIds": []
    }'

echo ""
echo "Certification 5 created"

# Certification 6: Go Programming Language (no expiry)
curl -s -X POST "$API_BASE/certifications" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "Go Programming Language Certification",
        "issuer": "Go Institute",
        "issueDate": "2022-01-20",
        "expiryDate": null,
        "credentialId": "GO-CERT-2022-567890",
        "verificationUrl": "https://www.goinstitute.org/verify/567890",
        "certificateUrl": "https://example.com/certificates/go-certification.pdf",
        "description": "Validates proficiency in Go programming language including concurrency patterns, interfaces, error handling, and best practices.",
        "status": "active",
        "category": "programming",
        "featured": true,
        "displayOrder": 6,
        "skillIds": []
    }'

echo ""
echo "Certification 6 created"

# Certification 7: Certified Information Systems Security Professional (CISSP)
curl -s -X POST "$API_BASE/certifications" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "Certified Information Systems Security Professional (CISSP)",
        "issuer": "(ISC)²",
        "issueDate": "2023-12-01",
        "expiryDate": "2026-12-01",
        "credentialId": "CISSP-123456",
        "verificationUrl": "https://www.isc2.org/verify/123456",
        "certificateUrl": "https://example.com/certificates/cissp.pdf",
        "description": "Validates expertise in information security. Covers security and risk management, asset security, security architecture, and more.",
        "status": "active",
        "category": "security",
        "featured": true,
        "displayOrder": 7,
        "skillIds": []
    }'

echo ""
echo "Certification 7 created"

# Certification 8: System Design Mastery (Custom/Internal)
curl -s -X POST "$API_BASE/certifications" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "name": "System Design Mastery Certification",
        "issuer": "Tech Academy",
        "issueDate": "2023-05-30",
        "expiryDate": null,
        "credentialId": "SDM-2023-001",
        "verificationUrl": "https://techacademy.com/verify/sdm-001",
        "certificateUrl": "https://example.com/certificates/system-design-mastery.pdf",
        "description": "Comprehensive certification covering distributed systems design, scalability patterns, microservices architecture, and performance optimization.",
        "status": "active",
        "category": "architecture",
        "featured": false,
        "displayOrder": 8,
        "skillIds": []
    }'

echo ""
echo "Certification 8 created"

echo ""
echo "All certifications created successfully!"
echo ""
echo "Note: Update the skillIds arrays to link certifications to actual skill IDs in your database"
echo "You can link skills using: POST /api/certifications/{certificationId}/skills/{skillId}"

