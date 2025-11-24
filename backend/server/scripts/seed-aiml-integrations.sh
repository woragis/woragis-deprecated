#!/bin/bash

# Seed script for AI/ML Integrations
# This script creates various AI/ML integration showcases

BASE_URL="${API_BASE_URL:-http://localhost:8080/api}"
AUTH_TOKEN="${AUTH_TOKEN:-}"

if [ -z "$AUTH_TOKEN" ]; then
    echo "Error: AUTH_TOKEN environment variable is not set"
    echo "Please set AUTH_TOKEN to a valid JWT token"
    exit 1
fi

echo "Seeding AI/ML Integrations..."

# Helper function to create an integration
create_integration() {
    local title=$1
    local description=$2
    local type=$3
    local framework=$4
    local model_name=$5
    local model_version=$6
    local use_case=$7
    local impact=$8
    local technologies=$9
    local architecture=${10}
    local metrics=${11}
    local featured=${12}
    local display_order=${13}
    local demo_url=${14}
    local github_url=${15}

    local payload="{
        \"title\": \"$title\",
        \"description\": \"$description\",
        \"type\": \"$type\",
        \"framework\": \"$framework\""

    if [ -n "$model_name" ]; then
        payload="$payload,
        \"modelName\": \"$model_name\""
    fi

    if [ -n "$model_version" ]; then
        payload="$payload,
        \"modelVersion\": \"$model_version\""
    fi

    if [ -n "$use_case" ]; then
        payload="$payload,
        \"useCase\": \"$use_case\""
    fi

    if [ -n "$impact" ]; then
        payload="$payload,
        \"impact\": \"$impact\""
    fi

    if [ -n "$technologies" ]; then
        payload="$payload,
        \"technologies\": [$technologies]"
    fi

    if [ -n "$architecture" ]; then
        payload="$payload,
        \"architecture\": \"$architecture\""
    fi

    if [ -n "$metrics" ]; then
        payload="$payload,
        \"metrics\": \"$metrics\""
    fi

    if [ "$featured" = "true" ]; then
        payload="$payload,
        \"featured\": true"
    fi

    if [ -n "$display_order" ]; then
        payload="$payload,
        \"displayOrder\": $display_order"
    fi

    if [ -n "$demo_url" ]; then
        payload="$payload,
        \"demoUrl\": \"$demo_url\""
    fi

    if [ -n "$github_url" ]; then
        payload="$payload,
        \"githubUrl\": \"$github_url\""
    fi

    payload="$payload
    }"

    response=$(curl -s -X POST "$BASE_URL/aiml-integrations" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $AUTH_TOKEN" \
        -d "$payload")

    if echo "$response" | grep -q '"id"'; then
        echo "✓ Created $type integration: $title"
    else
        echo "✗ Failed to create $type integration: $response"
    fi
}

# RAG Systems
echo ""
echo "Creating RAG (Retrieval-Augmented Generation) integrations..."
create_integration \
    "Enterprise Knowledge Base RAG System" \
    "Built a RAG system that enables semantic search across 50,000+ documents with 95% accuracy. Uses vector embeddings and hybrid search to provide context-aware responses." \
    "rag" \
    "langchain" \
    "gpt-4-turbo" \
    "latest" \
    "Enterprise knowledge management and customer support automation" \
    "Reduced support ticket resolution time by 70% and improved customer satisfaction scores by 40%" \
    "\"Python\", \"LangChain\", \"OpenAI\", \"Pinecone\", \"FastAPI\", \"PostgreSQL\"" \
    "Vector database (Pinecone) → Embedding service → LLM (GPT-4) → Response generation with source citations" \
    "{\"accuracy\": \"95%\", \"latency\": \"<2s\", \"cost_reduction\": \"60%\"}" \
    "true" \
    "1" \
    "https://demo.example.com/rag" \
    "https://github.com/example/rag-system"

create_integration \
    "Document Q&A System with RAG" \
    "Implemented a RAG-powered Q&A system for technical documentation. Uses chunking strategies and metadata filtering for precise answers." \
    "rag" \
    "llamaindex" \
    "claude-3-opus" \
    "20240229" \
    "Technical documentation search and Q&A" \
    "Enabled developers to find answers 10x faster than traditional search" \
    "\"Python\", \"LlamaIndex\", \"Anthropic Claude\", \"ChromaDB\", \"Next.js\"" \
    "Document ingestion → Chunking → Embedding → Vector store → Query processing → LLM synthesis" \
    "{\"retrieval_accuracy\": \"92%\", \"response_time\": \"1.5s\", \"user_satisfaction\": \"4.8/5\"}" \
    "true" \
    "2" \
    "" \
    "https://github.com/example/doc-qa-rag"

# LLM Integrations
echo ""
echo "Creating LLM integrations..."
create_integration \
    "Multi-LLM Chat Application" \
    "Built a chat application that intelligently routes queries to different LLMs (GPT-4, Claude, Gemini) based on task complexity and cost optimization." \
    "llm" \
    "langchain" \
    "gpt-4-turbo, claude-3-opus, gemini-pro" \
    "multi-model" \
    "Intelligent task routing and cost optimization for LLM usage" \
    "Reduced LLM costs by 45% while maintaining 98% response quality" \
    "\"Python\", \"LangChain\", \"OpenAI\", \"Anthropic\", \"Google AI\", \"Redis\", \"FastAPI\"" \
    "Request router → Model selector → LLM execution → Response caching → Cost tracking" \
    "{\"cost_reduction\": \"45%\", \"avg_latency\": \"1.2s\", \"quality_score\": \"98%\"}" \
    "true" \
    "3" \
    "https://demo.example.com/multi-llm" \
    ""

create_integration \
    "Code Generation Assistant" \
    "Developed an AI-powered code generation tool that understands context and generates production-ready code with tests." \
    "llm" \
    "openai" \
    "gpt-4-turbo" \
    "0125-preview" \
    "Automated code generation and refactoring" \
    "Increased developer productivity by 3x for boilerplate code generation" \
    "\"TypeScript\", \"OpenAI API\", \"VS Code Extension\", \"AST Parsing\"" \
    "Code context extraction → LLM prompt engineering → Code generation → Validation → Test generation" \
    "{\"code_quality\": \"95%\", \"test_coverage\": \"90%\", \"time_saved\": \"65%\"}" \
    "true" \
    "4" \
    "" \
    "https://github.com/example/code-gen-assistant"

# ML Models
echo ""
echo "Creating ML Model integrations..."
create_integration \
    "Customer Churn Prediction Model" \
    "Built a machine learning model that predicts customer churn with 89% accuracy using ensemble methods and feature engineering." \
    "ml_model" \
    "pytorch" \
    "XGBoost + Neural Network Ensemble" \
    "v2.1" \
    "Predictive analytics for customer retention" \
    "Reduced churn rate by 25% through proactive intervention" \
    "\"Python\", \"PyTorch\", \"XGBoost\", \"Scikit-learn\", \"Pandas\", \"MLflow\"" \
    "Data pipeline → Feature engineering → Model training → Model serving → Real-time inference" \
    "{\"accuracy\": \"89%\", \"precision\": \"87%\", \"recall\": \"91%\", \"f1_score\": \"89%\"}" \
    "false" \
    "0" \
    "" \
    ""

create_integration \
    "Image Classification API" \
    "Deployed a production-ready image classification service using transfer learning with ResNet50, achieving 94% accuracy." \
    "computer_vision" \
    "tensorflow" \
    "ResNet50" \
    "v1.0" \
    "Automated image categorization and tagging" \
    "Processed 1M+ images with 94% accuracy, reducing manual tagging by 90%" \
    "\"Python\", \"TensorFlow\", \"Keras\", \"FastAPI\", \"Docker\", \"Kubernetes\"" \
    "Image upload → Preprocessing → Model inference → Post-processing → Results API" \
    "{\"accuracy\": \"94%\", \"throughput\": \"1000 req/s\", \"latency\": \"50ms\"}" \
    "false" \
    "0" \
    "https://api.example.com/image-classify" \
    ""

# NLP Integrations
echo ""
echo "Creating NLP integrations..."
create_integration \
    "Sentiment Analysis Pipeline" \
    "Built a real-time sentiment analysis system that processes social media feeds and customer reviews with 92% accuracy." \
    "nlp" \
    "huggingface" \
    "distilbert-base-uncased-finetuned-sst-2-english" \
    "latest" \
    "Real-time sentiment monitoring and analysis" \
    "Enabled real-time brand monitoring and customer feedback analysis" \
    "\"Python\", \"Hugging Face Transformers\", \"Kafka\", \"Elasticsearch\", \"FastAPI\"" \
    "Data ingestion → Text preprocessing → Model inference → Sentiment scoring → Dashboard updates" \
    "{\"accuracy\": \"92%\", \"throughput\": \"5000 req/s\", \"latency\": \"30ms\"}" \
    "false" \
    "0" \
    "" \
    ""

# Recommendation Systems
echo ""
echo "Creating Recommendation System integrations..."
create_integration \
    "Personalized Content Recommendation Engine" \
    "Developed a hybrid recommendation system combining collaborative filtering and content-based approaches, improving engagement by 35%." \
    "recommendation" \
    "custom" \
    "Hybrid Collaborative + Content-Based" \
    "v3.0" \
    "Personalized content recommendations" \
    "Increased user engagement by 35% and click-through rates by 42%" \
    "\"Python\", \"TensorFlow\", \"Pandas\", \"Redis\", \"PostgreSQL\", \"Apache Spark\"" \
    "User behavior tracking → Feature extraction → Model training → Real-time recommendations → A/B testing" \
    "{\"engagement_increase\": \"35%\", \"ctr_improvement\": \"42%\", \"diversity_score\": \"0.78\"}" \
    "true" \
    "5" \
    "" \
    ""

# Chatbots
echo ""
echo "Creating Chatbot integrations..."
create_integration \
    "AI Customer Support Chatbot" \
    "Built an intelligent chatbot that handles 80% of customer inquiries without human intervention, using RAG for knowledge retrieval." \
    "chatbot" \
    "langchain" \
    "gpt-4-turbo" \
    "latest" \
    "24/7 customer support automation" \
    "Reduced support costs by 60% and improved response time from 2 hours to 30 seconds" \
    "\"Python\", \"LangChain\", \"OpenAI\", \"FastAPI\", \"WebSocket\", \"Redis\"" \
    "User query → Intent classification → RAG retrieval → LLM response → Human handoff (if needed)" \
    "{\"resolution_rate\": \"80%\", \"avg_response_time\": \"30s\", \"cost_reduction\": \"60%\"}" \
    "true" \
    "6" \
    "https://chat.example.com/support" \
    ""

# Generative AI
echo ""
echo "Creating Generative AI integrations..."
create_integration \
    "AI Content Generation Platform" \
    "Developed a platform that generates marketing content, blog posts, and social media copy using fine-tuned LLMs." \
    "generative_ai" \
    "openai" \
    "gpt-4-turbo" \
    "fine-tuned-v1" \
    "Automated content generation for marketing" \
    "Reduced content creation time by 75% while maintaining quality standards" \
    "\"Python\", \"OpenAI API\", \"FastAPI\", \"PostgreSQL\", \"React\"" \
    "Content request → Prompt engineering → LLM generation → Quality check → Human review → Publishing" \
    "{\"time_saved\": \"75%\", \"quality_score\": \"4.5/5\", \"cost_per_article\": \"\$2.50\"}" \
    "false" \
    "0" \
    "" \
    ""

# Anomaly Detection
echo ""
echo "Creating Anomaly Detection integrations..."
create_integration \
    "Real-time Fraud Detection System" \
    "Implemented an anomaly detection system for financial transactions using autoencoders and isolation forests." \
    "anomaly_detection" \
    "pytorch" \
    "Autoencoder + Isolation Forest Ensemble" \
    "v2.3" \
    "Real-time fraud detection in financial transactions" \
    "Detected 95% of fraudulent transactions with only 0.1% false positive rate" \
    "\"Python\", \"PyTorch\", \"Scikit-learn\", \"Kafka\", \"Redis\", \"PostgreSQL\"" \
    "Transaction stream → Feature extraction → Anomaly scoring → Alert system → Manual review queue" \
    "{\"detection_rate\": \"95%\", \"false_positive_rate\": \"0.1%\", \"latency\": \"<100ms\"}" \
    "false" \
    "0" \
    "" \
    ""

echo ""
echo "AI/ML Integrations seeding completed!"
echo ""
echo "You can now view featured integrations at: GET $BASE_URL/aiml-integrations/featured"
echo "Filter by type: GET $BASE_URL/aiml-integrations/type/rag"
echo "Filter by framework: GET $BASE_URL/aiml-integrations/framework/langchain"

