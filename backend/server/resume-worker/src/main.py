#!/usr/bin/env python3
"""
Resume Generation Worker - Main Entry Point

This worker generates customized resumes based on job requirements.
It can run in two modes:
1. Worker mode: Long-running process that processes jobs from Redis queue
2. CLI mode: One-time execution with command-line arguments (for testing)

Usage:
    Worker mode: python main.py --worker
    CLI mode: python main.py <user_id> <job_description> [job_title] [output_filename] [language]
"""

import os
import sys
import json
import logging
import argparse
from datetime import datetime
from dotenv import load_dotenv

# Add src directory to path for imports
sys.path.insert(0, os.path.dirname(__file__))

from database import Database
from ai_service import AIService
from resume_generator import ResumeGenerator
from translation_helper import TranslationHelper
from worker import Worker

# Load environment variables
load_dotenv()

# Configure basic logging for CLI mode
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

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


def run_cli_mode():
    """Run in CLI mode (one-time execution for testing)."""
    logger.info("Starting Resume Generation Worker (CLI mode)")
    
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
        logger.error(f"Error in CLI mode: {e}", exc_info=True)
        sys.exit(1)
    finally:
        db.close()
        translation_helper.close()


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Resume Generation Worker")
    parser.add_argument(
        "--worker",
        action="store_true",
        help="Run in worker mode (processes jobs from Redis queue)"
    )
    
    args = parser.parse_args()
    
    if args.worker:
        # Run in worker mode
        worker = Worker()
        try:
            worker.start()
        except KeyboardInterrupt:
            logger.info("Received interrupt signal")
        except Exception as e:
            logger.error("Worker failed to start", exc_info=True)
            sys.exit(1)
        finally:
            worker.stop()
    else:
        # Run in CLI mode (backward compatibility)
        run_cli_mode()


if __name__ == '__main__':
    main()
