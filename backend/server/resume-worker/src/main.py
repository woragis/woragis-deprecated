#!/usr/bin/env python3
"""
Resume Generation Worker - Main Entry Point

This worker generates customized resumes based on job requirements.
It searches through projects (by technology tags) and certifications,
uses AI to generate resume sections, and creates PDFs using WeasyPrint.
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

# Load environment variables
load_dotenv()

# Configure logging
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


def main():
    """Main entry point"""
    logger.info("Starting Resume Generation Worker")
    
    # Initialize components
    db = Database(DATABASE_URL)
    ai_service = AIService(AI_SERVICE_URL)
    generator = ResumeGenerator(db, ai_service, OUTPUT_DIR)
    
    try:
        db.connect()
        
        # Check command line arguments
        if len(sys.argv) < 3:
            logger.error("Usage: python main.py <user_id> <job_description> [job_title] [output_filename]")
            sys.exit(1)
        
        user_id = sys.argv[1]
        job_description = sys.argv[2]
        job_title = sys.argv[3] if len(sys.argv) > 3 else "Software Engineer"
        output_filename = sys.argv[4] if len(sys.argv) > 4 else None
        
        # Generate resume
        result = generator.generate_resume(
            user_id=user_id,
            job_description=job_description,
            job_title=job_title,
            output_filename=output_filename
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


if __name__ == '__main__':
    main()
