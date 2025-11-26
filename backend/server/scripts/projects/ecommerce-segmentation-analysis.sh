#!/bin/bash

# Seed script for E-commerce Customer Segmentation Analysis Project
# This script creates a complete project entry with related entities:
# - Project
# - Technical Writings
# - Problem Solutions
# - System Designs
# - Posts
# - Case Study

BASE_URL="${API_BASE_URL:-http://localhost:8080/api}"
AUTH_TOKEN="${AUTH_TOKEN:-}"

if [ -z "$AUTH_TOKEN" ]; then
    echo "Error: AUTH_TOKEN environment variable is not set"
    echo "Please set AUTH_TOKEN to a valid JWT token"
    exit 1
fi

echo "=========================================="
echo "Registering E-commerce Segmentation Analysis Project"
echo "=========================================="
echo ""

# Store created IDs for linking
PROJECT_ID=""
CASE_STUDY_ID=""

# Helper function to make API calls (returns response, prints status to stderr)
api_call() {
    local method=$1
    local endpoint=$2
    local payload=$3
    local description=$4
    
    if [ -z "$payload" ]; then
        response=$(curl -s -X "$method" "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $AUTH_TOKEN")
    else
        response=$(curl -s -X "$method" "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $AUTH_TOKEN" \
            -d "$payload")
    fi
    
    if echo "$response" | grep -q '"id"'; then
        echo "✓ $description" >&2
        echo "$response"
        return 0
    else
        echo "✗ Failed: $description" >&2
        echo "Response: $response" >&2
        echo "" >&2
        return 1
    fi
}

# Helper function to extract ID from response
extract_id() {
    local response=$1
    echo "$response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4
}

# ==========================================
# 1. CREATE PROJECT
# ==========================================
echo "1. Creating Project..."

PROJECT_PAYLOAD='{
    "name": "E-commerce Customer Segmentation Analysis",
    "description": "Machine learning project analyzing 100K+ customer transactions to identify distinct purchasing behavior segments using K-means clustering and RFM analysis. Delivers actionable insights for targeted marketing campaigns.",
    "status": "completed",
    "healthScore": 95,
    "mrr": 0,
    "cac": 0,
    "ltv": 0,
    "churnRate": 0
}'

PROJECT_RESPONSE=$(api_call "POST" "/projects" "$PROJECT_PAYLOAD" "Project created")
PROJECT_ID=$(extract_id "$PROJECT_RESPONSE")

if [ -z "$PROJECT_ID" ]; then
    echo "Failed to create project. Exiting."
    exit 1
fi

echo "Project ID: $PROJECT_ID"
echo ""

# ==========================================
# 2. CREATE TECHNICAL WRITINGS
# ==========================================
echo "2. Creating Technical Writings..."

# Technical Writing 1: Methodology Article
TECH_WRITING_1='{
    "title": "RFM Analysis and K-means Clustering for Customer Segmentation",
    "description": "A comprehensive guide to implementing RFM (Recency, Frequency, Monetary) analysis combined with K-means clustering for e-commerce customer segmentation. Covers data preprocessing, feature engineering, cluster optimization, and interpretation strategies.",
    "type": "article",
    "platform": "personal_blog",
    "url": "https://blog.example.com/rfm-kmeans-segmentation",
    "excerpt": "Learn how to combine traditional RFM analysis with modern machine learning techniques to create actionable customer segments.",
    "publishedAt": "2024-06-15T10:00:00Z",
    "readingTime": 18,
    "topics": ["Machine Learning", "Customer Segmentation", "RFM Analysis", "K-means Clustering", "Data Science"],
    "technologies": ["Python", "Pandas", "scikit-learn", "NumPy", "Matplotlib", "Seaborn"],
    "projectId": "'$PROJECT_ID'",
    "featured": true,
    "displayOrder": 1
}'

api_call "POST" "/technical-writings" "$TECH_WRITING_1" "Technical Writing: RFM Analysis Article"
echo ""

# Technical Writing 2: Tutorial
TECH_WRITING_2='{
    "title": "Complete Guide to Customer Segmentation with Python",
    "description": "Step-by-step tutorial on building a customer segmentation system from scratch. Includes data preprocessing, RFM score calculation, K-means implementation, and visualization techniques.",
    "type": "tutorial",
    "platform": "github",
    "url": "https://github.com/example/customer-segmentation-tutorial",
    "excerpt": "Build a production-ready customer segmentation system using Python and scikit-learn.",
    "publishedAt": "2024-06-20T10:00:00Z",
    "readingTime": 25,
    "topics": ["Tutorial", "Python", "Machine Learning", "Customer Segmentation", "Data Analysis"],
    "technologies": ["Python", "Pandas", "scikit-learn", "Jupyter Notebooks"],
    "projectId": "'$PROJECT_ID'",
    "featured": true,
    "displayOrder": 2
}'

api_call "POST" "/technical-writings" "$TECH_WRITING_2" "Technical Writing: Segmentation Tutorial"
echo ""

# Technical Writing 3: Case Study
TECH_WRITING_3='{
    "title": "Case Study: Improving Email Marketing Conversion by 18% Through Customer Segmentation",
    "description": "Real-world case study demonstrating how customer segmentation analysis led to a 18% improvement in email marketing conversion rates through targeted campaigns.",
    "type": "case_study",
    "platform": "personal_blog",
    "url": "https://blog.example.com/segmentation-case-study",
    "excerpt": "Discover how data-driven customer segmentation transformed marketing effectiveness.",
    "publishedAt": "2024-07-01T10:00:00Z",
    "readingTime": 12,
    "topics": ["Case Study", "Marketing", "Customer Segmentation", "Business Impact", "Data Science"],
    "technologies": ["Python", "Pandas", "scikit-learn"],
    "projectId": "'$PROJECT_ID'",
    "featured": true,
    "displayOrder": 3
}'

api_call "POST" "/technical-writings" "$TECH_WRITING_3" "Technical Writing: Marketing Case Study"
echo ""

# ==========================================
# 3. CREATE PROBLEM SOLUTIONS
# ==========================================
echo "3. Creating Problem Solutions..."

# Problem Solution 1: Handling Large Datasets
PROBLEM_SOLUTION_1='{
    "problem": "Processing 100K+ customer transaction records efficiently for segmentation analysis without running into memory constraints or performance bottlenecks.",
    "context": "Traditional pandas operations on large datasets can be slow and memory-intensive. Need to calculate RFM scores, perform clustering, and generate visualizations for 100K+ records efficiently.",
    "solution": "Implemented chunked processing for data loading, vectorized operations for RFM calculations using pandas groupby and aggregation functions, optimized K-means clustering with feature scaling, and used efficient visualization libraries. Created modular functions for each step to enable parallel processing and caching of intermediate results.",
    "technologies": ["Python", "Pandas", "NumPy", "scikit-learn"],
    "impact": "Reduced processing time from 45 minutes to 8 minutes. Memory usage decreased by 60%. Enabled real-time analysis capabilities.",
    "metrics": {
        "before": "Processing time: 45 minutes, Memory: 8GB",
        "after": "Processing time: 8 minutes, Memory: 3.2GB",
        "improvement": "82% faster, 60% less memory"
    },
    "featured": true
}'

api_call "POST" "/problem-solutions" "$PROBLEM_SOLUTION_1" "Problem Solution: Large Dataset Processing"
echo ""

# Problem Solution 2: Optimal Cluster Number Selection
PROBLEM_SOLUTION_2='{
    "problem": "Determining the optimal number of customer segments (clusters) that provides meaningful business insights without over-segmentation or under-segmentation.",
    "context": "K-means clustering requires specifying the number of clusters (k) beforehand. Too few clusters may miss important customer distinctions, while too many may create segments that are not actionable for marketing teams.",
    "solution": "Implemented a comprehensive cluster evaluation approach using elbow method, silhouette analysis, and gap statistic. Combined quantitative metrics with business domain knowledge to select 5 clusters that represent distinct customer personas: Champions, Loyal Customers, Potential Loyalists, At Risk, and Lost Customers.",
    "technologies": ["Python", "scikit-learn", "NumPy", "Matplotlib"],
    "impact": "Selected optimal cluster count that resulted in 35% difference in lifetime value across segments, enabling targeted marketing strategies with clear ROI.",
    "metrics": {
        "before": "Arbitrary cluster selection, unclear segment differentiation",
        "after": "5 validated segments with 35% LTV difference, clear personas",
        "improvement": "Data-driven segmentation with measurable business impact"
    },
    "featured": true
}'

api_call "POST" "/problem-solutions" "$PROBLEM_SOLUTION_2" "Problem Solution: Cluster Optimization"
echo ""

# ==========================================
# 4. CREATE SYSTEM DESIGNS
# ==========================================
echo "4. Creating System Designs..."

SYSTEM_DESIGN='{
    "title": "Customer Segmentation Analysis Pipeline Architecture",
    "description": "System design for an end-to-end customer segmentation analysis pipeline that processes transaction data, performs RFM analysis, applies machine learning clustering, and generates actionable insights.",
    "components": {
        "components": [
            {
                "name": "Data Ingestion Layer",
                "description": "Loads and validates customer transaction data from various sources (CSV, databases, APIs)",
                "technology": "Python, Pandas"
            },
            {
                "name": "Data Preprocessing Module",
                "description": "Cleans data, handles missing values, removes duplicates, and performs feature engineering",
                "technology": "Python, Pandas, NumPy"
            },
            {
                "name": "RFM Analysis Engine",
                "description": "Calculates Recency, Frequency, and Monetary scores for each customer using optimized aggregation functions",
                "technology": "Python, Pandas"
            },
            {
                "name": "Clustering Service",
                "description": "Implements K-means clustering with feature scaling, cluster optimization, and evaluation metrics",
                "technology": "Python, scikit-learn, NumPy"
            },
            {
                "name": "Visualization Generator",
                "description": "Creates charts, graphs, and dashboards for segment analysis and business reporting",
                "technology": "Python, Matplotlib, Seaborn"
            },
            {
                "name": "Insights Engine",
                "description": "Generates customer personas, calculates lifetime value, and provides marketing recommendations",
                "technology": "Python, Pandas"
            }
        ]
    },
    "dataFlow": "Transaction data → Data Ingestion → Preprocessing → RFM Calculation → Feature Scaling → Clustering → Visualization → Insights Generation → Marketing Recommendations",
    "scalability": "Designed to handle datasets from 10K to 1M+ records. Uses chunked processing for large datasets. Can be parallelized across multiple cores. Supports incremental updates for new transaction data.",
    "reliability": "Includes data validation at each stage, error handling for edge cases, and logging for debugging. Implements reproducible analysis with fixed random seeds. Version control for analysis code and results.",
    "diagram": "https://example.com/diagrams/customer-segmentation-pipeline.png",
    "featured": true
}'

api_call "POST" "/system-designs" "$SYSTEM_DESIGN" "System Design: Segmentation Pipeline"
echo ""

# ==========================================
# 5. CREATE POSTS
# ==========================================
echo "5. Creating Posts..."

# Post 1: Project Overview
POST_1='{
    "title": "E-commerce Customer Segmentation: From Data to Actionable Insights",
    "content": "# E-commerce Customer Segmentation Analysis\n\n## Overview\n\nThis project analyzes 100K+ customer transactions to identify distinct purchasing behavior segments using K-means clustering. Through comprehensive RFM (Recency, Frequency, Monetary) analysis, we reveal 5 key customer personas and deliver actionable insights for targeted marketing campaigns.\n\n## Key Achievements\n\n- ✅ Analyzed **100K+ customer transactions** to identify distinct purchasing behavior segments\n- ✅ Implemented **K-means clustering** for customer segmentation\n- ✅ Performed **RFM (Recency, Frequency, Monetary) analysis** revealing **5 key customer personas**\n- ✅ Created data visualizations demonstrating **35% difference in lifetime value** across segments\n- ✅ Delivered actionable recommendations adopted for email marketing campaigns, improving conversion by **18%**\n\n## Customer Segments Identified\n\n1. **Champions** - High value, frequent, recent customers\n2. **Loyal Customers** - Regular purchasers with good recency\n3. **Potential Loyalists** - Recent customers with growing frequency\n4. **At Risk** - Previously active but declining engagement\n5. **Lost Customers** - Inactive for extended periods\n\n## Technologies Used\n\n- **Python 3.8+**\n- **Pandas** - Data manipulation and analysis\n- **NumPy** - Numerical computing\n- **scikit-learn** - Machine learning (K-means clustering)\n- **Matplotlib/Seaborn** - Data visualization\n- **Jupyter Notebooks** - Interactive analysis\n\n## Business Impact\n\n- **35% difference in lifetime value** across customer segments\n- **18% improvement in email marketing conversion** through targeted campaigns\n- Actionable recommendations for personalized marketing strategies",
    "excerpt": "Discover how machine learning and RFM analysis transformed customer segmentation, leading to 18% improvement in marketing conversion rates.",
    "status": "published",
    "featured": true,
    "metaTitle": "E-commerce Customer Segmentation Analysis | Data Science Project",
    "metaDescription": "Machine learning project analyzing 100K+ customer transactions using K-means clustering and RFM analysis. Improved marketing conversion by 18%.",
    "metaKeywords": "customer segmentation, RFM analysis, K-means clustering, machine learning, data science, e-commerce"
}'

api_call "POST" "/posts" "$POST_1" "Post: Project Overview" > /dev/null
echo ""

# Post 2: Methodology Deep Dive
POST_2='{
    "title": "RFM Analysis and K-means Clustering: A Complete Methodology",
    "content": "# RFM Analysis and K-means Clustering Methodology\n\n## Introduction\n\nThis document outlines the technical methodology for performing customer segmentation analysis using RFM analysis and K-means clustering on e-commerce transaction data.\n\n## 1. Data Preprocessing\n\n### Data Collection\n- Source: E-commerce transaction database\n- Minimum required fields:\n  - `customer_id`: Unique customer identifier\n  - `transaction_date`: Date of purchase\n  - `transaction_amount`: Monetary value of transaction\n  - `order_id`: Unique order identifier\n\n### Data Cleaning\n1. **Handle Missing Values** - Identify and handle missing data\n2. **Remove Duplicates** - Ensure data integrity\n3. **Outlier Detection** - Identify and handle statistical outliers\n4. **Data Type Conversion** - Ensure proper data types\n5. **Data Validation** - Validate data ranges and consistency\n\n## 2. RFM Analysis\n\n### RFM Score Calculation\n\n**Recency (R):** Days since customer's last purchase\n- Lower values = higher recency (more recent)\n- Typical scale: 1-5\n\n**Frequency (F):** Number of transactions in analysis period\n- Higher values = higher frequency\n- Typical scale: 1-5\n\n**Monetary (M):** Total amount spent by customer\n- Higher values = higher monetary value\n- Typical scale: 1-5\n\n## 3. K-means Clustering\n\n### Feature Selection\n- Recency (normalized)\n- Frequency (normalized)\n- Monetary value (normalized)\n- Average transaction value\n- Customer lifetime\n- Purchase frequency rate\n\n### Optimal Cluster Number\n- Elbow Method\n- Silhouette Analysis\n- Gap Statistic\n\n## 4. Customer Persona Development\n\nFor each cluster, we analyze:\n- Behavioral Patterns\n- Value Metrics\n- Purchase Preferences\n- Lifetime Value\n\n## 5. Visualization Strategy\n\n- RFM score distributions\n- Cluster visualizations\n- Business insights charts\n- Segment comparisons",
    "excerpt": "A comprehensive guide to implementing RFM analysis and K-means clustering for customer segmentation.",
    "status": "published",
    "featured": false,
    "metaTitle": "RFM Analysis and K-means Clustering Methodology | Customer Segmentation",
    "metaDescription": "Complete methodology for customer segmentation using RFM analysis and K-means clustering with Python.",
    "metaKeywords": "RFM analysis, K-means clustering, customer segmentation, methodology, data science"
}'

api_call "POST" "/posts" "$POST_2" "Post: Methodology Deep Dive"
echo ""

# ==========================================
# 6. CREATE CASE STUDY
# ==========================================
echo "6. Creating Case Study..."

CASE_STUDY='{
    "title": "E-commerce Customer Segmentation: Driving 18% Marketing Conversion Improvement",
    "description": "A detailed case study documenting the complete customer segmentation analysis project, from data collection to actionable marketing insights that improved email campaign conversion by 18%.",
    "challenge": "An e-commerce business needed to understand their customer base better to create targeted marketing campaigns. With 100K+ transactions and no clear segmentation strategy, marketing efforts were generic and had low conversion rates. The challenge was to identify distinct customer segments that would enable personalized marketing strategies.",
    "solution": "Implemented a comprehensive customer segmentation analysis combining traditional RFM (Recency, Frequency, Monetary) analysis with modern K-means clustering. The solution involved: 1) Data preprocessing and quality assessment, 2) RFM score calculation using quintile-based scoring, 3) Feature engineering and normalization, 4) K-means clustering with optimal cluster selection (5 segments), 5) Customer persona development, 6) Lifetime value analysis, and 7) Actionable marketing recommendations for each segment.",
    "technologies": ["Python", "Pandas", "NumPy", "scikit-learn", "Matplotlib", "Seaborn", "Jupyter Notebooks"],
    "architecture": "The analysis pipeline consists of: Data Ingestion → Preprocessing → RFM Calculation → Feature Scaling → K-means Clustering → Persona Development → Visualization → Insights Generation. Each stage is modular and can be run independently or as part of the complete pipeline.",
    "metrics": {
        "metrics": [
            {
                "label": "Dataset Size",
                "value": "100K+ transactions",
                "improvement": "Comprehensive analysis coverage"
            },
            {
                "label": "Customer Segments Identified",
                "value": "5 distinct personas",
                "improvement": "Clear, actionable segments"
            },
            {
                "label": "Lifetime Value Difference",
                "value": "35% across segments",
                "improvement": "Significant segment differentiation"
            },
            {
                "label": "Marketing Conversion Improvement",
                "value": "18% increase",
                "improvement": "Targeted campaigns effectiveness"
            },
            {
                "label": "Processing Time",
                "value": "8 minutes",
                "improvement": "82% faster than initial implementation"
            }
        ]
    },
    "tradeoffs": {
        "tradeoffs": [
            {
                "decision": "Using 5 clusters instead of 3 or 7",
                "pros": ["More granular segmentation", "Better captures customer diversity", "More actionable for marketing"],
                "cons": ["Slightly more complex to manage", "Requires more marketing resources"]
            },
            {
                "decision": "Combining RFM with K-means instead of using only one",
                "pros": ["Leverages domain knowledge (RFM)", "Adds machine learning insights", "More robust segmentation"],
                "cons": ["More complex implementation", "Requires understanding both approaches"]
            }
        ]
    },
    "lessonsLearned": [
        "Feature scaling is critical for K-means clustering performance",
        "Business validation is as important as statistical validation for cluster selection",
        "Visualization is key for communicating insights to non-technical stakeholders",
        "Modular code design enables iterative analysis and experimentation",
        "Combining traditional analysis (RFM) with ML (K-means) provides best results"
    ]
}'

CASE_STUDY_RESPONSE=$(api_call "POST" "/projects/$PROJECT_ID/case-studies" "$CASE_STUDY" "Case Study: Project Case Study")
CASE_STUDY_ID=$(extract_id "$CASE_STUDY_RESPONSE")
echo ""

# ==========================================
# 7. ADD PROJECT TECHNOLOGIES
# ==========================================
echo "7. Adding Project Technologies..."

# Technology 1: Python
TECH_1='{
    "name": "Python",
    "version": "3.8+",
    "category": "backend",
    "purpose": "Primary programming language for data analysis and machine learning",
    "link": "https://www.python.org/"
}'

api_call "POST" "/projects/$PROJECT_ID/technologies" "$TECH_1" "Technology: Python"
echo ""

# Technology 2: Pandas
TECH_2='{
    "name": "Pandas",
    "version": "1.3+",
    "category": "backend",
    "purpose": "Data manipulation and analysis, RFM calculations, data preprocessing",
    "link": "https://pandas.pydata.org/"
}'

api_call "POST" "/projects/$PROJECT_ID/technologies" "$TECH_2" "Technology: Pandas"
echo ""

# Technology 3: scikit-learn
TECH_3='{
    "name": "scikit-learn",
    "version": "0.24+",
    "category": "backend",
    "purpose": "K-means clustering implementation and machine learning utilities",
    "link": "https://scikit-learn.org/"
}'

api_call "POST" "/projects/$PROJECT_ID/technologies" "$TECH_3" "Technology: scikit-learn"
echo ""

# Technology 4: NumPy
TECH_4='{
    "name": "NumPy",
    "version": "1.20+",
    "category": "backend",
    "purpose": "Numerical computing, array operations, feature scaling",
    "link": "https://numpy.org/"
}'

api_call "POST" "/projects/$PROJECT_ID/technologies" "$TECH_4" "Technology: NumPy"
echo ""

# Technology 5: Matplotlib
TECH_5='{
    "name": "Matplotlib",
    "version": "3.3+",
    "category": "backend",
    "purpose": "Data visualization and chart generation",
    "link": "https://matplotlib.org/"
}'

api_call "POST" "/projects/$PROJECT_ID/technologies" "$TECH_5" "Technology: Matplotlib"
echo ""

# Technology 6: Seaborn
TECH_6='{
    "name": "Seaborn",
    "version": "0.11+",
    "category": "backend",
    "purpose": "Statistical data visualization and advanced plotting",
    "link": "https://seaborn.pydata.org/"
}'

api_call "POST" "/projects/$PROJECT_ID/technologies" "$TECH_6" "Technology: Seaborn"
echo ""

# Technology 7: Jupyter Notebooks
TECH_7='{
    "name": "Jupyter Notebooks",
    "version": "6.0+",
    "category": "other",
    "purpose": "Interactive analysis, documentation, and reproducible research",
    "link": "https://jupyter.org/"
}'

api_call "POST" "/projects/$PROJECT_ID/technologies" "$TECH_7" "Technology: Jupyter Notebooks"
echo ""

# ==========================================
# SUMMARY
# ==========================================
echo ""
echo "=========================================="
echo "Registration Complete!"
echo "=========================================="
echo ""
echo "Project ID: $PROJECT_ID"
echo "Case Study ID: $CASE_STUDY_ID"
echo ""
echo "Created entities:"
echo "  ✓ 1 Project"
echo "  ✓ 3 Technical Writings"
echo "  ✓ 2 Problem Solutions"
echo "  ✓ 1 System Design"
echo "  ✓ 2 Posts"
echo "  ✓ 1 Case Study"
echo "  ✓ 7 Technologies"
echo ""
echo "View project at: GET $BASE_URL/projects/$PROJECT_ID"
echo "View case study at: GET $BASE_URL/projects/$PROJECT_ID/case-studies"
echo ""

