#!/bin/bash

# Script to seed testimonials with context using curl
# Usage: ./seed-testimonials.sh [email] [password]

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

# Create testimonials
echo "Creating testimonials..."

# Testimonial 1: General testimonial with context
curl -s -X POST "$API_BASE/testimonials" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "authorName": "Sarah Chen",
        "authorRole": "Senior Engineering Manager",
        "authorCompany": "TechCorp Inc.",
        "authorPhoto": "https://i.pravatar.cc/150?img=1",
        "content": "Working with Jezreel was an absolute pleasure. His technical expertise and problem-solving skills are exceptional. He consistently delivered high-quality solutions and was always willing to go the extra mile to ensure project success.",
        "context": "Worked together on a critical microservices migration project over 6 months. Jezreel led the architecture design and implementation.",
        "type": "general",
        "rating": 5,
        "linkedinUrl": "https://linkedin.com/in/sarahchen",
        "displayOrder": 1
    }'

echo ""
echo "Testimonial 1 created"

# Testimonial 2: Project-specific testimonial
curl -s -X POST "$API_BASE/testimonials" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "authorName": "Michael Rodriguez",
        "authorRole": "CTO",
        "authorCompany": "StartupXYZ",
        "authorPhoto": "https://i.pravatar.cc/150?img=2",
        "content": "Jezreel transformed our backend infrastructure. The system he built handles 10x more traffic with better performance. His attention to detail and understanding of scalability challenges is impressive.",
        "context": "Hired Jezreel to rebuild our backend API. The project involved migrating from a monolithic architecture to microservices using Go and PostgreSQL.",
        "type": "project_specific",
        "rating": 5,
        "linkedinUrl": "https://linkedin.com/in/michaelrodriguez",
        "displayOrder": 2,
        "entityLinks": [
            {
                "entityType": "project",
                "entityId": "00000000-0000-0000-0000-000000000001"
            }
        ]
    }'

echo ""
echo "Testimonial 2 created"

# Testimonial 3: Skill-specific testimonial (Go/Backend)
curl -s -X POST "$API_BASE/testimonials" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "authorName": "David Kim",
        "authorRole": "Lead Backend Engineer",
        "authorCompany": "CloudScale Systems",
        "authorPhoto": "https://i.pravatar.cc/150?img=3",
        "content": "Jezreel is one of the best Go developers I'\''ve worked with. His code is clean, well-structured, and follows best practices. He has a deep understanding of concurrency patterns and system design.",
        "context": "Collaborated on a high-performance Go service that processes millions of requests per day. Jezreel implemented critical optimizations and design patterns.",
        "type": "skill_specific",
        "rating": 5,
        "linkedinUrl": "https://linkedin.com/in/davidkim",
        "displayOrder": 3,
        "entityLinks": [
            {
                "entityType": "skill",
                "entityId": "00000000-0000-0000-0000-000000000002"
            }
        ]
    }'

echo ""
echo "Testimonial 3 created"

# Testimonial 4: General testimonial with video
curl -s -X POST "$API_BASE/testimonials" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "authorName": "Emily Watson",
        "authorRole": "Product Manager",
        "authorCompany": "InnovateLab",
        "authorPhoto": "https://i.pravatar.cc/150?img=4",
        "content": "Jezreel is an exceptional engineer who combines technical excellence with great communication skills. He always explains complex technical concepts clearly and works collaboratively with the team.",
        "context": "Worked together on multiple product features over 2 years. Jezreel was instrumental in delivering key features on time and with high quality.",
        "type": "general",
        "videoUrl": "https://example.com/videos/testimonial-emily-watson.mp4",
        "rating": 5,
        "linkedinUrl": "https://linkedin.com/in/emilywatson",
        "displayOrder": 4
    }'

echo ""
echo "Testimonial 4 created"

# Testimonial 5: Project-specific with detailed context
curl -s -X POST "$API_BASE/testimonials" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "authorName": "James Anderson",
        "authorRole": "Founder & CEO",
        "authorCompany": "DataFlow Analytics",
        "authorPhoto": "https://i.pravatar.cc/150?img=5",
        "content": "Jezreel built our entire data pipeline infrastructure from scratch. The system he designed is robust, scalable, and maintainable. We'\''ve had zero downtime since launch.",
        "context": "Hired Jezreel as a consultant to design and implement a real-time data processing pipeline. The project involved Kafka, PostgreSQL, Redis, and Go microservices. Completed in 3 months with excellent documentation.",
        "type": "project_specific",
        "rating": 5,
        "linkedinUrl": "https://linkedin.com/in/jamesanderson",
        "displayOrder": 5,
        "entityLinks": [
            {
                "entityType": "project",
                "entityId": "00000000-0000-0000-0000-000000000003"
            }
        ]
    }'

echo ""
echo "Testimonial 5 created"

# Testimonial 6: Skill-specific (System Design)
curl -s -X POST "$API_BASE/testimonials" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d '{
        "authorName": "Lisa Park",
        "authorRole": "Principal Architect",
        "authorCompany": "Enterprise Solutions",
        "authorPhoto": "https://i.pravatar.cc/150?img=6",
        "content": "Jezreel'\''s system design skills are outstanding. He can break down complex problems into elegant, scalable solutions. His architecture diagrams and documentation are always top-notch.",
        "context": "Worked together on designing a distributed system for handling millions of concurrent users. Jezreel created comprehensive architecture documentation and led technical discussions.",
        "type": "skill_specific",
        "rating": 5,
        "linkedinUrl": "https://linkedin.com/in/lisapark",
        "displayOrder": 6,
        "entityLinks": [
            {
                "entityType": "skill",
                "entityId": "00000000-0000-0000-0000-000000000004"
            }
        ]
    }'

echo ""
echo "Testimonial 6 created"

echo ""
echo "All testimonials created successfully!"
echo ""
echo "Note: Update the entityIds in entityLinks to match actual project/skill IDs in your database"

