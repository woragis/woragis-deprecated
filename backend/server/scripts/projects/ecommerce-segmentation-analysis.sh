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
# 8. ADDITIONAL TECHNICAL WRITINGS (DEEP DIVES)
# ==========================================
echo "8. Creating Additional Technical Writings (Deep Dives)..."

# Technical Writing 4: RFM Implementation Deep Dive
TECH_WRITING_4='{
    "title": "Implementing RFM Analysis from Scratch: Quintile Scoring and Segment Mapping",
    "description": "A comprehensive deep dive into implementing RFM (Recency, Frequency, Monetary) analysis using Python and Pandas. Covers quintile-based scoring, custom threshold methods, RFM cell creation, segment mapping strategies, and validation techniques. Includes code examples and best practices for production implementations.",
    "type": "article",
    "platform": "github",
    "url": "https://github.com/example/rfm-analysis-implementation",
    "excerpt": "Master RFM analysis implementation with detailed code walkthroughs, scoring strategies, and segment mapping techniques.",
    "publishedAt": "2024-06-18T10:00:00Z",
    "readingTime": 22,
    "topics": ["RFM Analysis", "Python", "Pandas", "Customer Segmentation", "Data Analysis", "Implementation"],
    "technologies": ["Python", "Pandas", "NumPy"],
    "projectId": "'$PROJECT_ID'",
    "featured": true,
    "displayOrder": 4
}'

api_call "POST" "/technical-writings" "$TECH_WRITING_4" "Technical Writing: RFM Implementation Deep Dive"
echo ""

# Technical Writing 5: K-means Optimization Guide
TECH_WRITING_5='{
    "title": "Finding the Optimal Number of Clusters: Elbow Method, Silhouette Analysis, and Beyond",
    "description": "Complete guide to determining optimal cluster count for K-means clustering. Covers elbow method implementation, silhouette score analysis, gap statistic, Davies-Bouldin index, and Calinski-Harabasz index. Includes Python implementations, visualization techniques, and decision frameworks for choosing the right number of clusters.",
    "type": "guide",
    "platform": "personal_blog",
    "url": "https://blog.example.com/kmeans-optimization",
    "excerpt": "Learn multiple methods for finding the optimal number of clusters in K-means clustering with practical Python examples.",
    "publishedAt": "2024-06-22T10:00:00Z",
    "readingTime": 20,
    "topics": ["K-means Clustering", "Machine Learning", "Optimization", "Data Science", "Python"],
    "technologies": ["Python", "scikit-learn", "NumPy", "Matplotlib"],
    "projectId": "'$PROJECT_ID'",
    "featured": true,
    "displayOrder": 5
}'

api_call "POST" "/technical-writings" "$TECH_WRITING_5" "Technical Writing: K-means Optimization Guide"
echo ""

# Technical Writing 6: Data Preprocessing Best Practices
TECH_WRITING_6='{
    "title": "Data Preprocessing for Customer Segmentation: Cleaning, Validation, and Feature Engineering",
    "description": "Comprehensive guide to preprocessing e-commerce transaction data for customer segmentation. Covers data validation, missing value handling, outlier detection using IQR and Z-score methods, duplicate removal, data type conversion, and feature engineering techniques. Includes production-ready Python code and quality check frameworks.",
    "type": "tutorial",
    "platform": "dev_to",
    "url": "https://dev.to/example/data-preprocessing-segmentation",
    "excerpt": "Master data preprocessing techniques for customer segmentation with real-world examples and production code.",
    "publishedAt": "2024-06-25T10:00:00Z",
    "readingTime": 16,
    "topics": ["Data Preprocessing", "Data Cleaning", "Feature Engineering", "Python", "Data Quality"],
    "technologies": ["Python", "Pandas", "NumPy"],
    "projectId": "'$PROJECT_ID'",
    "featured": false,
    "displayOrder": 6
}'

api_call "POST" "/technical-writings" "$TECH_WRITING_6" "Technical Writing: Data Preprocessing Guide"
echo ""

# Technical Writing 7: Visualization Strategies
TECH_WRITING_7='{
    "title": "Data Visualization for Customer Segmentation: RFM Heatmaps, Cluster Plots, and Business Dashboards",
    "description": "Complete guide to creating effective visualizations for customer segmentation analysis. Covers RFM heatmaps, cluster scatter plots with PCA, silhouette plots, elbow curves, segment comparison charts, LTV visualizations, and interactive dashboards. Includes Matplotlib, Seaborn, and Plotly examples with design best practices.",
    "type": "article",
    "platform": "medium",
    "url": "https://medium.com/@example/segmentation-visualization",
    "excerpt": "Create compelling visualizations that communicate customer segmentation insights effectively to stakeholders.",
    "publishedAt": "2024-06-28T10:00:00Z",
    "readingTime": 18,
    "topics": ["Data Visualization", "Matplotlib", "Seaborn", "Customer Segmentation", "Business Intelligence"],
    "technologies": ["Python", "Matplotlib", "Seaborn", "Plotly"],
    "projectId": "'$PROJECT_ID'",
    "featured": true,
    "displayOrder": 7
}'

api_call "POST" "/technical-writings" "$TECH_WRITING_7" "Technical Writing: Visualization Strategies"
echo ""

# Technical Writing 8: Customer Persona Deep Dive - Champions
TECH_WRITING_8='{
    "title": "Understanding the Champions Segment: High-Value Customer Analysis and Retention Strategies",
    "description": "Deep analysis of the Champions customer segment - high value, frequent, recent customers. Covers segment characteristics, behavioral patterns, lifetime value calculation, purchase frequency analysis, and targeted retention strategies. Includes case studies and marketing recommendations for maximizing value from this critical segment.",
    "type": "case_study",
    "platform": "personal_blog",
    "url": "https://blog.example.com/champions-segment-analysis",
    "excerpt": "Deep dive into the Champions customer segment with actionable insights for retention and value maximization.",
    "publishedAt": "2024-07-05T10:00:00Z",
    "readingTime": 14,
    "topics": ["Customer Segmentation", "Customer Retention", "Marketing Strategy", "Business Analysis"],
    "technologies": ["Python", "Pandas"],
    "projectId": "'$PROJECT_ID'",
    "featured": false,
    "displayOrder": 8
}'

api_call "POST" "/technical-writings" "$TECH_WRITING_8" "Technical Writing: Champions Segment Analysis"
echo ""

# ==========================================
# 9. ADDITIONAL POSTS (DETAILED PHASE COVERAGE)
# ==========================================
echo "9. Creating Additional Posts (Detailed Phase Coverage)..."

# Post 3: Phase 1 - Data Collection & Preparation
POST_3='{
    "title": "Phase 1: Data Collection and Preparation for Customer Segmentation",
    "content": "# Phase 1: Data Collection and Preparation\n\n## Overview\n\nThis post covers the first phase of our customer segmentation project: collecting and preparing 100K+ transaction records for analysis.\n\n## Data Collection\n\n### Data Sources\n- E-commerce transaction database\n- 100K+ customer transaction records\n- Time period: Last 12-24 months\n\n### Required Fields\n- `customer_id`: Unique customer identifier\n- `transaction_date`: Date and time of purchase\n- `transaction_amount`: Monetary value of transaction\n- `order_id`: Unique order identifier\n\n## Data Quality Assessment\n\n### Checklist Completed\n- ✅ Total record count: 100,000+ transactions\n- ✅ Date range coverage: 12 months\n- ✅ Unique customer count: 15,000+ customers\n- ✅ Missing value percentage: < 2%\n- ✅ Duplicate record count: Identified and removed\n- ✅ Data type consistency: Validated\n- ✅ Value range validation: Completed\n- ✅ Outlier identification: Using IQR method\n\n## Data Cleaning Process\n\n### Steps Taken\n1. **Missing Value Handling**: Removed records with critical missing values (< 2% of data)\n2. **Duplicate Removal**: Identified and removed 1,234 duplicate transactions\n3. **Outlier Detection**: Used IQR method to identify and cap extreme values\n4. **Data Type Conversion**: Converted dates to datetime, amounts to float\n5. **Data Validation**: Removed negative amounts, future dates, invalid customer IDs\n\n### Results\n- Initial records: 102,345\n- After cleaning: 100,234\n- Data quality score: 98.2%\n\n## Feature Engineering\n\n### Base Features Created\n- Transaction count per customer\n- Total spending per customer\n- Average transaction value\n- Days since first purchase\n- Days since last purchase\n\n### Derived Features\n- Customer lifetime (days between first and last purchase)\n- Purchase frequency (transactions per month)\n- Average days between purchases\n\n## Data Dictionary\n\nCreated comprehensive data dictionary documenting:\n- Field names and descriptions\n- Data types\n- Value ranges\n- Business rules\n- Quality metrics\n\n## Next Steps\n\nWith clean, validated data in hand, we proceed to Phase 2: Exploratory Data Analysis.",
    "excerpt": "Detailed walkthrough of data collection and preparation phase, including quality assessment, cleaning procedures, and feature engineering.",
    "status": "published",
    "featured": false,
    "metaTitle": "Data Collection and Preparation for Customer Segmentation | Phase 1",
    "metaDescription": "Complete guide to collecting and preparing 100K+ transaction records for customer segmentation analysis.",
    "metaKeywords": "data collection, data preparation, data cleaning, customer segmentation, data quality"
}'

api_call "POST" "/posts" "$POST_3" "Post: Phase 1 - Data Collection" > /dev/null
echo ""

# Post 4: Phase 2 - Exploratory Data Analysis
POST_4='{
    "title": "Phase 2: Exploratory Data Analysis - Uncovering Customer Behavior Patterns",
    "content": "# Phase 2: Exploratory Data Analysis\n\n## Overview\n\nThis post details our comprehensive exploratory data analysis (EDA) that uncovered key patterns in customer purchasing behavior.\n\n## Univariate Analysis\n\n### Transaction Amount Distribution\n- **Mean**: $45.32\n- **Median**: $32.15\n- **Standard Deviation**: $28.90\n- **Skewness**: 2.3 (right-skewed)\n- **Key Insight**: Most transactions are small, but a few large transactions significantly impact average\n\n### Transaction Date Patterns\n- **Peak Day**: Friday (18% of transactions)\n- **Peak Month**: November (holiday season)\n- **Seasonality**: Strong seasonal patterns with 40% increase in Q4\n- **Key Insight**: Clear temporal patterns suggest time-based segmentation opportunities\n\n### Customer Metrics\n- **Average transactions per customer**: 6.7\n- **Average total spending**: $303.64\n- **Average transaction value**: $45.32\n- **Customer lifetime range**: 1-365 days\n\n## Bivariate Analysis\n\n### Key Relationships Discovered\n1. **Transaction Count vs. Total Spending**: Strong positive correlation (r=0.78)\n2. **Customer Lifetime vs. Total Spending**: Moderate correlation (r=0.52)\n3. **Frequency vs. Recency**: Negative correlation (r=-0.34)\n4. **Monetary vs. Frequency**: Strong positive correlation (r=0.71)\n\n### Visualizations Created\n- Scatter plots showing relationships\n- Correlation matrices\n- Heatmaps for RFM dimensions\n\n## Multivariate Analysis\n\n### RFM Dimensions Correlation\n- R-F correlation: -0.34\n- R-M correlation: -0.28\n- F-M correlation: 0.71\n\n### Customer Behavior Patterns\n- High frequency customers tend to have higher monetary value\n- Recent customers show varied frequency patterns\n- Long-term customers (high lifetime) show consistent spending\n\n## Key Findings\n\n1. **Pareto Principle**: 20% of customers account for 60% of revenue\n2. **Recency Matters**: Recent customers show higher engagement\n3. **Frequency-Value Link**: More frequent buyers spend more overall\n4. **Seasonal Impact**: Q4 shows 40% transaction increase\n\n## Insights for Segmentation\n\n- RFM analysis will be effective given clear patterns\n- Clustering should consider temporal patterns\n- Segment size will likely be imbalanced (Pareto effect)\n- Need to account for seasonality in analysis\n\n## Next Steps\n\nWith EDA complete, we proceed to Phase 3: RFM Analysis to quantify customer value.",
    "excerpt": "Comprehensive exploratory data analysis revealing customer behavior patterns, correlations, and insights that inform segmentation strategy.",
    "status": "published",
    "featured": true,
    "metaTitle": "Exploratory Data Analysis for Customer Segmentation | Phase 2",
    "metaDescription": "Detailed EDA findings including transaction patterns, customer metrics, and behavior insights for segmentation.",
    "metaKeywords": "exploratory data analysis, EDA, customer behavior, data patterns, segmentation"
}'

api_call "POST" "/posts" "$POST_4" "Post: Phase 2 - EDA" > /dev/null
echo ""

# Post 5: Phase 3 - RFM Analysis Implementation
POST_5='{
    "title": "Phase 3: RFM Analysis Implementation - Quantifying Customer Value",
    "content": "# Phase 3: RFM Analysis Implementation\n\n## Overview\n\nThis post covers the implementation of RFM (Recency, Frequency, Monetary) analysis to quantify customer value and create initial segments.\n\n## RFM Calculation Process\n\n### Step 1: Define Analysis Date\n```python\nanalysis_date = max(transaction_date)  # Most recent transaction date\n```\n\n### Step 2: Calculate RFM Metrics\n\n**Recency Calculation:**\n```python\nrecency = analysis_date - last_transaction_date\ndays_since_last_purchase = recency.days\n```\n\n**Frequency Calculation:**\n```python\nfrequency = transactions.groupby(\"customer_id\")[\"order_id\"].count()\n```\n\n**Monetary Calculation:**\n```python\nmonetary = transactions.groupby(\"customer_id\")[\"transaction_amount\"].sum()\n```\n\n## RFM Scoring Strategy\n\n### Quintile-Based Scoring (Selected Method)\n\nWe chose quintile-based scoring for its statistical robustness:\n\n- **Recency**: Lower values = Higher score (1 = most recent, 5 = least recent)\n- **Frequency**: Higher values = Higher score (1 = least frequent, 5 = most frequent)\n- **Monetary**: Higher values = Higher score (1 = lowest spending, 5 = highest spending)\n\n### Implementation\n```python\nrfm_df[\"R_Score\"] = pd.qcut(rfm_df[\"recency\"].rank(method=\"first\"), 5, labels=[5,4,3,2,1])\nrfm_df[\"F_Score\"] = pd.qcut(rfm_df[\"frequency\"].rank(method=\"first\"), 5, labels=[1,2,3,4,5])\nrfm_df[\"M_Score\"] = pd.qcut(rfm_df[\"monetary\"].rank(method=\"first\"), 5, labels=[1,2,3,4,5])\n```\n\n## RFM Segment Assignment\n\n### Segment Definitions\n\n**Champions (555, 554, 544, 545, 454, 455, 445):**\n- Best customers, high value, frequent, recent\n- 12% of customer base\n- 35% of total revenue\n\n**Loyal Customers (543, 444, 435, 355, 354, 345, 344, 335):**\n- Regular buyers with good recency\n- 18% of customer base\n- 28% of total revenue\n\n**Potential Loyalists (512, 511, 422, 421, 412, 411, 311):**\n- Recent customers with growing frequency\n- 15% of customer base\n- 12% of total revenue\n\n**At Risk (344, 343, 334, 333, 323, 322, 233, 232, 223, 222):**\n- Previously active but declining engagement\n- 20% of customer base\n- 15% of total revenue\n\n**Lost Customers (111, 112, 113, 114, 115):**\n- Inactive for extended periods\n- 35% of customer base\n- 10% of total revenue\n\n## Results Summary\n\n- **Total customers analyzed**: 15,234\n- **RFM segments created**: 125 unique RFM cells\n- **Primary segments identified**: 5 major groups\n- **Revenue concentration**: Top 2 segments account for 63% of revenue\n\n## Validation\n\n- All customers assigned to segments\n- No empty segments\n- Segment sizes reasonable (5-35%)\n- Business validation: Segments align with known customer behaviors\n\n## Next Steps\n\nWith RFM segments established, we proceed to Phase 4: K-means Clustering to validate and refine segments using machine learning.",
    "excerpt": "Complete implementation of RFM analysis including scoring methodology, segment definitions, and results validation.",
    "status": "published",
    "featured": true,
    "metaTitle": "RFM Analysis Implementation | Phase 3 | Customer Segmentation",
    "metaDescription": "Detailed RFM analysis implementation with quintile scoring, segment mapping, and validation results.",
    "metaKeywords": "RFM analysis, customer segmentation, RFM scoring, customer value, segmentation"
}'

api_call "POST" "/posts" "$POST_5" "Post: Phase 3 - RFM Analysis" > /dev/null
echo ""

# Post 6: Phase 4 - K-means Clustering
POST_6='{
    "title": "Phase 4: K-means Clustering - Machine Learning Validation of Customer Segments",
    "content": "# Phase 4: K-means Clustering\n\n## Overview\n\nThis post details the K-means clustering implementation that validated and refined our RFM-based customer segments using machine learning.\n\n## Feature Preparation\n\n### Selected Features\n- Recency (days)\n- Frequency (count)\n- Monetary (total amount)\n\n### Feature Scaling\n\n**Why Scaling is Critical:**\nK-means uses Euclidean distance, so features with larger scales dominate. We used StandardScaler to normalize all features.\n\n```python\nfrom sklearn.preprocessing import StandardScaler\n\nscaler = StandardScaler()\nfeatures_scaled = scaler.fit_transform(rfm_features)\n```\n\n## Optimal Cluster Number Determination\n\n### Elbow Method\n\nTested k values from 2 to 10:\n- **K=2**: Inertia = 45,234\n- **K=3**: Inertia = 32,156\n- **K=4**: Inertia = 24,891\n- **K=5**: Inertia = 19,234 (elbow point)\n- **K=6**: Inertia = 16,789\n- **K=7**: Inertia = 15,123\n\n**Elbow identified at K=5**\n\n### Silhouette Analysis\n\n- **K=2**: Silhouette = 0.42\n- **K=3**: Silhouette = 0.48\n- **K=4**: Silhouette = 0.51\n- **K=5**: Silhouette = 0.54 (optimal)\n- **K=6**: Silhouette = 0.49\n\n**Optimal K = 5** (highest silhouette score)\n\n### Additional Metrics\n\n**Davies-Bouldin Index** (lower is better):\n- K=5: 0.89 (best)\n\n**Calinski-Harabasz Index** (higher is better):\n- K=5: 1,234 (best)\n\n## K-means Implementation\n\n```python\nfrom sklearn.cluster import KMeans\n\nkmeans = KMeans(n_clusters=5, random_state=42, n_init=10)\nclusters = kmeans.fit_predict(features_scaled)\n```\n\n## Cluster Analysis\n\n### Cluster Characteristics\n\n**Cluster 0 - Champions:**\n- Size: 1,829 customers (12%)\n- Avg Recency: 15 days\n- Avg Frequency: 12.3 transactions\n- Avg Monetary: $1,234\n- Avg LTV: $2,456\n\n**Cluster 1 - Loyal Customers:**\n- Size: 2,742 customers (18%)\n- Avg Recency: 45 days\n- Avg Frequency: 8.7 transactions\n- Avg Monetary: $856\n- Avg LTV: $1,789\n\n**Cluster 2 - Potential Loyalists:**\n- Size: 2,285 customers (15%)\n- Avg Recency: 25 days\n- Avg Frequency: 4.2 transactions\n- Avg Monetary: $423\n- Avg LTV: $987\n\n**Cluster 3 - At Risk:**\n- Size: 3,047 customers (20%)\n- Avg Recency: 120 days\n- Avg Frequency: 3.1 transactions\n- Avg Monetary: $234\n- Avg LTV: $456\n\n**Cluster 4 - Lost Customers:**\n- Size: 5,331 customers (35%)\n- Avg Recency: 280 days\n- Avg Frequency: 1.2 transactions\n- Avg Monetary: $89\n- Avg LTV: $123\n\n## Clustering Quality Metrics\n\n- **Silhouette Score**: 0.54 (good)\n- **Davies-Bouldin Index**: 0.89 (excellent)\n- **Calinski-Harabasz Index**: 1,234 (excellent)\n\n## Comparison with RFM Segments\n\n- **Alignment**: 78% of customers in same segment category\n- **Refinement**: Clustering identified 22% of customers that RFM misclassified\n- **Validation**: Clusters confirm RFM segment validity\n\n## Lifetime Value Analysis\n\n- **Highest LTV**: Champions - $2,456\n- **Lowest LTV**: Lost Customers - $123\n- **LTV Difference**: 1,897% (35% when normalized)\n\n## Next Steps\n\nWith validated clusters, we proceed to Phase 5: Customer Persona Development to create actionable marketing personas.",
    "excerpt": "Complete K-means clustering implementation including optimal cluster selection, quality metrics, and detailed cluster analysis.",
    "status": "published",
    "featured": true,
    "metaTitle": "K-means Clustering for Customer Segmentation | Phase 4",
    "metaDescription": "Detailed K-means clustering implementation with optimal cluster selection, quality metrics, and cluster analysis results.",
    "metaKeywords": "K-means clustering, machine learning, customer segmentation, cluster analysis, optimal clusters"
}'

api_call "POST" "/posts" "$POST_6" "Post: Phase 4 - K-means Clustering" > /dev/null
echo ""

# ==========================================
# 10. ADDITIONAL CASE STUDIES
# ==========================================
echo "10. Creating Additional Case Studies..."

# Case Study 2: Data Preprocessing Challenges
CASE_STUDY_2='{
    "title": "Case Study: Handling 100K+ Transaction Records - Data Quality and Preprocessing Challenges",
    "description": "Detailed case study on the data preprocessing challenges encountered when processing 100K+ e-commerce transaction records. Covers missing value strategies, outlier detection methods, duplicate handling, and performance optimization techniques.",
    "challenge": "Processing 100K+ transaction records with inconsistent data quality, missing values, duplicates, outliers, and performance bottlenecks. Initial processing took 45 minutes and consumed 8GB of memory, making iterative analysis impractical.",
    "solution": "Implemented a comprehensive data preprocessing pipeline with: 1) Chunked data loading for memory efficiency, 2) Vectorized pandas operations for speed, 3) IQR-based outlier detection with business logic validation, 4) Intelligent missing value handling (removal for critical fields, imputation for optional), 5) Deduplication with transaction fingerprinting, 6) Data validation rules with automated reporting, 7) Caching of intermediate results for iterative analysis.",
    "technologies": ["Python", "Pandas", "NumPy"],
    "architecture": "Data Ingestion → Validation → Cleaning → Outlier Detection → Feature Engineering → Quality Reporting → Cached Output",
    "metrics": {
        "metrics": [
            {
                "label": "Processing Time",
                "value": "8 minutes (from 45 minutes)",
                "improvement": "82% faster"
            },
            {
                "label": "Memory Usage",
                "value": "3.2GB (from 8GB)",
                "improvement": "60% reduction"
            },
            {
                "label": "Data Quality Score",
                "value": "98.2%",
                "improvement": "From 87% initial quality"
            },
            {
                "label": "Records Processed",
                "value": "100,234 (from 102,345)",
                "improvement": "2.1% removed as invalid"
            }
        ]
    },
    "tradeoffs": {
        "tradeoffs": [
            {
                "decision": "Removing vs. Imputing Missing Values",
                "pros": ["Removal ensures data integrity", "No artificial data introduced", "Simpler validation"],
                "cons": ["Loss of 2% of records", "Potential bias if missing is non-random"]
            },
            {
                "decision": "IQR vs. Z-score for Outlier Detection",
                "pros": ["IQR is robust to extreme values", "Business-validated outliers kept", "Better for skewed distributions"],
                "cons": ["More complex implementation", "Requires domain expertise"]
            }
        ]
    },
    "lessonsLearned": [
        "Chunked processing is essential for large datasets",
        "Vectorized operations provide 10x speed improvement",
        "Business validation of outliers prevents data loss",
        "Caching intermediate results enables iterative analysis",
        "Data quality reporting is crucial for stakeholder confidence"
    ]
}'

CASE_STUDY_2_RESPONSE=$(api_call "POST" "/projects/$PROJECT_ID/case-studies" "$CASE_STUDY_2" "Case Study: Data Preprocessing")
echo ""

# Case Study 3: Cluster Optimization Journey
CASE_STUDY_3='{
    "title": "Case Study: Finding the Optimal Number of Clusters - A Data-Driven Journey",
    "description": "Case study documenting the iterative process of determining the optimal number of customer segments using multiple evaluation methods. Covers elbow method, silhouette analysis, gap statistic, and business validation approaches.",
    "challenge": "Determining the optimal number of customer segments (k) for K-means clustering. Too few clusters miss important customer distinctions, while too many create segments that are not actionable for marketing. Initial attempts used arbitrary k values (3, 7, 10) without systematic evaluation.",
    "solution": "Implemented a comprehensive cluster evaluation framework: 1) Elbow method to identify inflection point (k=5), 2) Silhouette analysis for cluster quality (optimal at k=5 with score 0.54), 3) Davies-Bouldin index for cluster separation (best at k=5 with 0.89), 4) Calinski-Harabasz index for cluster variance (best at k=5 with 1,234), 5) Business validation with marketing team to ensure segments are actionable, 6) Cluster stability testing with multiple random seeds.",
    "technologies": ["Python", "scikit-learn", "NumPy", "Matplotlib"],
    "architecture": "Feature Preparation → Multiple K Testing → Metric Calculation → Visualization → Business Validation → Final Selection",
    "metrics": {
        "metrics": [
            {
                "label": "Optimal K Selected",
                "value": "5 clusters",
                "improvement": "Data-driven decision vs. arbitrary selection"
            },
            {
                "label": "Silhouette Score",
                "value": "0.54",
                "improvement": "Good cluster quality (0.3+ is acceptable)"
            },
            {
                "label": "Cluster Stability",
                "value": "92% consistency",
                "improvement": "High reproducibility across runs"
            },
            {
                "label": "Business Alignment",
                "value": "100% actionable",
                "improvement": "All segments usable for marketing"
            }
        ]
    },
    "tradeoffs": {
        "tradeoffs": [
            {
                "decision": "K=5 vs. K=4 or K=6",
                "pros": ["K=5 provides best statistical metrics", "Segments are actionable", "Good balance of granularity"],
                "cons": ["Slightly more complex than K=4", "Requires more marketing resources than fewer segments"]
            },
            {
                "decision": "Statistical vs. Business Validation",
                "pros": ["Combined approach ensures both quality and usability", "Reduces risk of over-optimization"],
                "cons": ["More time-consuming", "Requires stakeholder involvement"]
            }
        ]
    },
    "lessonsLearned": [
        "Multiple evaluation methods provide confidence in selection",
        "Business validation is as important as statistical metrics",
        "Cluster stability testing prevents overfitting to specific data",
        "Visualization (elbow curve, silhouette plot) aids decision-making",
        "Optimal k balances statistical quality with business practicality"
    ]
}'

api_call "POST" "/projects/$PROJECT_ID/case-studies" "$CASE_STUDY_3" "Case Study: Cluster Optimization" > /dev/null
echo ""

# ==========================================
# 11. ADDITIONAL PROBLEM SOLUTIONS
# ==========================================
echo "11. Creating Additional Problem Solutions..."

# Problem Solution 3: Feature Scaling for K-means
PROBLEM_SOLUTION_3='{
    "problem": "K-means clustering produced poor results with unbalanced cluster sizes and unclear segment differentiation. Features had vastly different scales (Recency: 0-365 days, Frequency: 1-50, Monetary: $10-$10,000), causing the algorithm to be dominated by the Monetary feature.",
    "context": "Initial K-means implementation used raw RFM values without scaling. The Monetary feature (range $10-$10,000) dominated the Euclidean distance calculations, causing clusters to form primarily based on spending amount while ignoring Recency and Frequency patterns. This resulted in poor segmentation that didn't capture customer behavior nuances.",
    "solution": "Implemented StandardScaler from scikit-learn to normalize all features to have mean=0 and standard deviation=1. This ensures all features contribute equally to distance calculations. Applied scaling before clustering and maintained the scaler for future predictions. Also tested MinMaxScaler as alternative but StandardScaler provided better results.",
    "technologies": ["Python", "scikit-learn", "NumPy"],
    "impact": "Clustering quality improved dramatically. Silhouette score increased from 0.28 to 0.54. Cluster sizes became more balanced (12-35% vs. previous 5-60%). Segments now capture all three RFM dimensions effectively, resulting in more meaningful customer personas.",
    "metrics": {
        "before": "Silhouette: 0.28, Unbalanced clusters (5-60% size), Monetary-dominated",
        "after": "Silhouette: 0.54, Balanced clusters (12-35% size), All features contribute",
        "improvement": "93% better silhouette, balanced segmentation, multi-dimensional insights"
    },
    "featured": true
}'

api_call "POST" "/problem-solutions" "$PROBLEM_SOLUTION_3" "Problem Solution: Feature Scaling"
echo ""

# Problem Solution 4: Handling Imbalanced Customer Segments
PROBLEM_SOLUTION_4='{
    "problem": "Customer segments are highly imbalanced with Lost Customers representing 35% of the base but only 10% of revenue, while Champions represent 12% of base but 35% of revenue. Marketing teams need strategies that account for this imbalance without ignoring valuable smaller segments.",
    "context": "Natural customer behavior follows Pareto principle (80/20 rule), resulting in imbalanced segments. Traditional equal-weight strategies would over-invest in low-value segments and under-invest in high-value ones. Need to develop marketing strategies that account for segment value and size.",
    "solution": "Developed value-weighted marketing strategies: 1) Revenue-based prioritization (focus on Champions and Loyal Customers), 2) Growth potential assessment (invest in Potential Loyalists), 3) Retention campaigns for At Risk (prevent churn), 4) Re-engagement for Lost (low-cost win-back), 5) Segment-specific budget allocation based on LTV contribution, 6) Different campaign frequencies and channels per segment.",
    "technologies": ["Business Strategy", "Data Analysis"],
    "impact": "Marketing ROI improved by 25%. Campaign conversion rates increased: Champions (45%), Loyal Customers (32%), Potential Loyalists (28%), At Risk (18%), Lost (8%). Budget allocation optimized to focus 60% on top 30% of customers (Champions + Loyal).",
    "metrics": {
        "before": "Equal treatment, 18% average conversion, inefficient budget allocation",
        "after": "Value-weighted strategies, 26% average conversion, optimized budget",
        "improvement": "44% conversion increase, 25% ROI improvement, strategic budget allocation"
    },
    "featured": false
}'

api_call "POST" "/problem-solutions" "$PROBLEM_SOLUTION_4" "Problem Solution: Imbalanced Segments"
echo ""

# Problem Solution 5: Interpreting Cluster Results for Business
PROBLEM_SOLUTION_5='{
    "problem": "Machine learning clustering results were statistically valid but difficult for marketing teams to understand and act upon. Technical metrics (silhouette scores, cluster centers) didn't translate to actionable business insights.",
    "context": "Data science team delivered clustering results with technical metrics, but marketing team couldn't understand how to use them. Cluster numbers (0-4) and statistical measures didn't provide business context. Need to bridge the gap between technical analysis and business application.",
    "solution": "Created comprehensive cluster interpretation framework: 1) Persona naming (Champions, Loyal, Potential Loyalists, At Risk, Lost), 2) Business-friendly descriptions with characteristics, 3) Visualizations (radar charts, comparison bar charts), 4) Marketing recommendations per persona, 5) LTV calculations and revenue contribution, 6) Behavioral pattern summaries, 7) Actionable next steps for each segment.",
    "technologies": ["Data Visualization", "Business Communication", "Python", "Matplotlib", "Seaborn"],
    "impact": "Marketing team adoption increased from 30% to 95%. Campaign development time reduced by 40%. All 5 personas now have dedicated marketing strategies. Stakeholder presentations are more effective with business-focused insights.",
    "metrics": {
        "before": "30% team adoption, unclear actions, long campaign development",
        "after": "95% team adoption, clear personas and strategies, 40% faster campaigns",
        "improvement": "217% adoption increase, actionable insights, significant time savings"
    },
    "featured": true
}'

api_call "POST" "/problem-solutions" "$PROBLEM_SOLUTION_5" "Problem Solution: Cluster Interpretation"
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
echo "  ✓ 8 Technical Writings (including deep dives)"
echo "  ✓ 5 Problem Solutions"
echo "  ✓ 1 System Design"
echo "  ✓ 6 Posts (covering all project phases)"
echo "  ✓ 3 Case Studies (comprehensive coverage)"
echo "  ✓ 7 Technologies"
echo ""
echo "Content Coverage:"
echo "  • RFM Analysis implementation details"
echo "  • K-means clustering optimization"
echo "  • Data preprocessing best practices"
echo "  • Visualization strategies"
echo "  • Customer persona deep dives"
echo "  • All 7 project phases documented"
echo "  • Technical challenges and solutions"
echo "  • Business impact case studies"
echo ""
echo "View project at: GET $BASE_URL/projects/$PROJECT_ID"
echo "View case study at: GET $BASE_URL/projects/$PROJECT_ID/case-studies"
echo ""

