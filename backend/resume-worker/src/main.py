#!/usr/bin/env python3
"""
Resume Generation Worker - Main Entry Point

This worker generates customized resumes based on job requirements.
It searches through projects (by technology tags) and certifications,
uses AI to generate resume sections, and creates PDFs using WeasyPrint.

Supports two modes:
1. CLI mode: Direct invocation with command-line arguments
2. Queue mode: Consumes jobs from RabbitMQ queue (default when no args provided)
"""

import os
import sys
import json
import logging
from datetime import datetime
from dotenv import load_dotenv

import sys
import os

# Add src directory to path for imports
sys.path.insert(0, os.path.dirname(__file__))

from database import Database
from ai_service import AIService
from resume_generator import ResumeGenerator
from translation_helper import TranslationHelper
from queue_consumer import create_consumer_from_env
from metrics import record_job_processed, record_job_failed
import time

# Load environment variables
load_dotenv()

# Configure structured logging
from logger import configure_logging, get_logger

env = os.getenv("ENV", "development")
log_to_file = os.getenv("LOG_TO_FILE", "false").lower() == "true"
log_dir = os.getenv("LOG_DIR", "logs")
configure_logging(env=env, log_to_file=log_to_file, log_dir=log_dir)

logger = get_logger()

# Configuration
DATABASE_URL = os.getenv('DATABASE_URL', 'postgres://postgres:postgres@database:5432/woragis?sslmode=disable')
AI_SERVICE_URL = os.getenv('AI_SERVICE_URL', 'http://ai-service:8000')
OUTPUT_DIR = os.getenv('RESUME_OUTPUT_DIR', '/app/output')
RESULTS_LOG_DIR = os.getenv('RESULTS_LOG_DIR', '/app/results')


def save_result(result: dict, results_dir: str = RESULTS_LOG_DIR):
    """Save generation result to a JSON file for tracking"""
    os.makedirs(results_dir, exist_ok=True)
    
    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    result_filename = f"resume_result_{result['user_id']}_{timestamp}.json"
    result_path = os.path.join(results_dir, result_filename)
    
    try:
        with open(result_path, 'w') as f:
            json.dump(result, f, indent=2, default=str)
        logger.info(f"Result saved to: {result_path}")
    except Exception as e:
        logger.error(f"Failed to save result: {e}")


def process_resume_job(message: dict) -> bool:
    """
    Process a resume generation job from the queue.
    
    Args:
        message: Job message containing:
            - user_id: User ID
            - job_description: Job description text
            - job_title: Job title (optional, default: "Software Engineer")
            - output_filename: Output filename (optional)
            - language: Language code (optional, default: "en")
    
    Returns:
        True if job was processed successfully, False otherwise
    """
    # Extract job parameters
    user_id = message.get('user_id')
    job_description = message.get('job_description')
    job_title = message.get('job_title', 'Software Engineer')
    output_filename = message.get('output_filename')
    language = message.get('language', 'en')
    
    if not user_id or not job_description:
        logger.error("Missing required fields: user_id and job_description are required")
        return False
    
    logger.info(f"Processing resume generation job for user: {user_id}, language: {language}")
    
    start_time = time.time()
    worker_name = "resume-worker"
    
    db = None
    translation_helper = None
    
    try:
        # Initialize components
        db = Database(DATABASE_URL)
        ai_service = AIService(AI_SERVICE_URL)
        translation_helper = TranslationHelper(DATABASE_URL)
        generator = ResumeGenerator(db, ai_service, OUTPUT_DIR, translation_helper)
        
        db.connect()
        translation_helper.connect()
        
        # Generate resume
        result = generator.generate_resume(
            user_id=user_id,
            job_description=job_description,
            job_title=job_title,
            output_filename=output_filename,
            language=language
        )
        
        # Save result metadata
        save_result(result, RESULTS_LOG_DIR)
        
        duration = time.time() - start_time
        record_job_processed(worker_name, "success", duration)
        logger.info(f"Resume successfully generated: {result['output_path']}")
        logger.info(f"File size: {result['file_size']} bytes")
        logger.info(f"Projects included: {result['projects_count']}")
        logger.info(f"Certifications included: {result['certifications_count']}")
        
        return True
        
    except Exception as e:
        duration = time.time() - start_time
        record_job_processed(worker_name, "failed", duration)
        record_job_failed(worker_name, "processing_error")
        logger.error(f"Error processing resume job: {e}", exc_info=True)
        return False
    finally:
        if db:
            db.close()
        if translation_helper:
            translation_helper.close()


def run_cli_mode():
    """Run in CLI mode (direct invocation with command-line arguments)"""
    logger.info("Starting Resume Generation Worker (CLI Mode)")
    
    # Initialize components
    db = Database(DATABASE_URL)
    ai_service = AIService(AI_SERVICE_URL)
    translation_helper = TranslationHelper(DATABASE_URL)
    generator = ResumeGenerator(db, ai_service, OUTPUT_DIR, translation_helper)
    
    try:
        db.connect()
        translation_helper.connect()
        
        # Check command line arguments
        if len(sys.argv) < 3:
            logger.error("Usage: python main.py <user_id> <job_description> [job_title] [output_filename] [language]")
            sys.exit(1)
        
        user_id = sys.argv[1]
        job_description = sys.argv[2]
        job_title = sys.argv[3] if len(sys.argv) > 3 else "Software Engineer"
        output_filename = sys.argv[4] if len(sys.argv) > 4 else None
        language = sys.argv[5] if len(sys.argv) > 5 else "en"
        
        # Generate resume
        result = generator.generate_resume(
            user_id=user_id,
            job_description=job_description,
            job_title=job_title,
            output_filename=output_filename,
            language=language
        )
        
        # Save result metadata
        save_result(result, RESULTS_LOG_DIR)
        
        logger.info(f"Resume successfully generated: {result['output_path']}")
        print(f"Resume generated: {result['output_path']}")
        print(f"File size: {result['file_size']} bytes")
        print(f"Projects included: {result['projects_count']}")
        print(f"Certifications included: {result['certifications_count']}")
        
        # Return success
        sys.exit(0)
        
    except Exception as e:
        logger.error(f"Error in main: {e}", exc_info=True)
        sys.exit(1)
    finally:
        db.close()
        translation_helper.close()


def run_queue_mode():
    """Run in queue mode (consume jobs from RabbitMQ)"""
    logger.info("Starting Resume Generation Worker (Queue Mode)")
    
    # Start health check HTTP server
    from http.server import HTTPServer, BaseHTTPRequestHandler
    import json as json_lib
    from health import check_health
    from prometheus_client import generate_latest, CONTENT_TYPE_LATEST
    
    class HealthHandler(BaseHTTPRequestHandler):
        def do_GET(self):
            if self.path == '/healthz':
                result = check_health()
                status_code = 200 if result["status"] != "unhealthy" else 503
                
                self.send_response(status_code)
                self.send_header('Content-Type', 'application/json')
                self.end_headers()
                self.wfile.write(json_lib.dumps(result).encode())
            elif self.path == '/metrics':
                # Prometheus metrics endpoint
                self.send_response(200)
                self.send_header('Content-Type', CONTENT_TYPE_LATEST)
                self.end_headers()
                self.wfile.write(generate_latest())
            else:
                self.send_response(404)
                self.end_headers()
        
        def log_message(self, format, *args):
            # Suppress default HTTP server logs
            pass
    
    health_server = HTTPServer(('0.0.0.0', 8080), HealthHandler)
    import threading
    health_thread = threading.Thread(target=health_server.serve_forever, daemon=True)
    health_thread.start()
    logger.info("Health check server started on port 8080 (includes /metrics endpoint)")
    
    try:
        # Create queue consumer
        consumer = create_consumer_from_env()
        
        # Connect to RabbitMQ
        if not consumer.connect():
            logger.error("Failed to connect to RabbitMQ, exiting...")
            sys.exit(1)
        
        # Start consuming messages
        consumer.start_consuming(process_resume_job)
        
    except KeyboardInterrupt:
        logger.info("Received keyboard interrupt, shutting down...")
    except Exception as e:
        logger.error(f"Error in queue mode: {e}", exc_info=True)
        sys.exit(1)


def main():
    """Main entry point - determines mode based on command-line arguments"""
    # If command-line arguments are provided, run in CLI mode
    # Otherwise, run in queue mode
    if len(sys.argv) > 1:
        # CLI mode: requires at least user_id and job_description
        if len(sys.argv) >= 3:
            run_cli_mode()
        else:
            logger.error("CLI mode requires at least <user_id> and <job_description>")
            logger.error("Usage: python main.py <user_id> <job_description> [job_title] [output_filename] [language]")
            sys.exit(1)
    else:
        # Queue mode: no arguments provided
        run_queue_mode()


if __name__ == '__main__':
    main()
