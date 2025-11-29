"""
Main resume generator that orchestrates the resume creation process.
"""

import os
import logging
from datetime import datetime
from typing import Dict, List, Optional
from jinja2 import Environment, FileSystemLoader
from weasyprint import HTML, CSS

import sys
import os

# Add src directory to path for imports
sys.path.insert(0, os.path.dirname(__file__))

from database import Database
from ai_service import AIService
from keyword_extractor import extract_keywords

logger = logging.getLogger(__name__)


class ResumeGenerator:
    """Main resume generator"""
    
    def __init__(self, db: Database, ai_service: AIService, output_dir: str = "/app/output"):
        self.db = db
        self.ai_service = ai_service
        self.output_dir = output_dir
        
        # Try multiple possible template locations
        possible_dirs = [
            os.path.join(os.path.dirname(__file__), '..', 'templates'),
            os.path.join(os.path.dirname(__file__), '..', '..', 'templates'),
            '/app/templates',
            './templates'
        ]
        template_dir = None
        for dir_path in possible_dirs:
            abs_path = os.path.abspath(dir_path)
            if os.path.exists(abs_path) and os.path.isdir(abs_path):
                template_dir = abs_path
                break
        
        if not template_dir:
            raise FileNotFoundError(f"Could not find templates directory. Tried: {possible_dirs}")
        
        self.template_dir = template_dir
        self.template_env = Environment(
            loader=FileSystemLoader(template_dir),
            autoescape=True
        )
        
        # Ensure output directory exists
        os.makedirs(self.output_dir, exist_ok=True)
    
    def generate_resume(
        self,
        user_id: str,
        job_description: str,
        job_title: str = "Software Engineer",
        output_filename: str = None
    ) -> Dict[str, any]:
        """Generate a customized resume PDF and return metadata"""
        logger.info(f"Generating resume for user {user_id}, job: {job_title}")
        
        # Get user info
        user_info = self.db.get_user_info(user_id)
        if not user_info:
            raise ValueError(f"User {user_id} not found")
        
        # Extract keywords and determine relevant categories
        keywords = extract_keywords(job_description)
        logger.info(f"Extracted keywords: {keywords}")
        
        # Fetch relevant projects
        projects = self.db.get_user_projects(
            user_id,
            tech_categories=keywords['tech_categories'],
            skill_names=keywords.get('skill_names')
        )
        logger.info(f"Found {len(projects)} relevant projects")
        
        # Fetch relevant certifications
        certifications = self.db.get_user_certifications(
            user_id,
            categories=keywords['cert_categories']
        )
        logger.info(f"Found {len(certifications)} relevant certifications")
        
        # Generate resume sections using AI
        logger.info("Generating resume sections with AI...")
        profile = self.ai_service.generate_resume_section('profile', job_description, projects)
        about = self.ai_service.generate_resume_section('about', job_description, projects)
        experience = self.ai_service.generate_resume_section('experience', job_description, projects)
        skills = self.ai_service.generate_resume_section('skills', job_description, projects)
        
        # Prepare template context
        # Set default social links for this user
        website = None
        github = None
        linkedin = None
        
        user_email = user_info.get('email', '')
        if user_email == 'masteringthecode.woragis@gmail.com':
            website = 'www.woragis.me'
            github = 'github.com/woragis'
            linkedin = 'linkedin.com/in/jezreel-andrade'
        
        # Set the name - use hardcoded name for this user
        display_name = "Jezreel de Andrade Galvao Veloso"
        if user_email == 'masteringthecode.woragis@gmail.com':
            display_name = "Jezreel de Andrade Galvao Veloso"
        
        # Project limit logic: up to 4 projects, but if exactly 3, show only 2
        if projects:
            if len(projects) == 3:
                projects_to_show = projects[:2]
            else:
                projects_to_show = projects[:4]
        else:
            projects_to_show = []
        
        context = {
            'name': display_name,
            'email': user_info.get('email', ''),
            'website': website,
            'github': github,
            'linkedin': linkedin,
            'profile': profile,
            'about': about,
            'experience': experience,
            'skills': skills,
            'projects': projects_to_show,
            'certifications': certifications if certifications else [],
            'job_title': job_title,
            'generated_date': datetime.now().strftime('%B %Y')
        }
        
        # Render HTML template
        template = self.template_env.get_template('resume.html')
        html_content = template.render(**context)
        
        # Generate PDF filename
        if not output_filename:
            safe_job_title = "".join(c for c in job_title if c.isalnum() or c in (' ', '-', '_')).strip()
            safe_job_title = safe_job_title.replace(' ', '_')
            timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
            output_filename = f"resume_{safe_job_title}_{timestamp}.pdf"
        
        output_path = os.path.join(self.output_dir, output_filename)
        
        # Generate PDF with WeasyPrint
        css_path = os.path.join(self.template_dir, 'style.css')
        try:
            HTML(string=html_content, base_url=self.template_dir).write_pdf(
                output_path,
                stylesheets=[CSS(filename=css_path)]
            )
        except Exception as e:
            # Fallback: try without base_url
            logger.warning(f"PDF generation with base_url failed: {e}, trying without base_url")
            HTML(string=html_content).write_pdf(
                output_path,
                stylesheets=[CSS(filename=css_path)]
            )
        
        # Get file size
        file_size = os.path.getsize(output_path)
        
        logger.info(f"Resume generated: {output_path} ({file_size} bytes)")
        
        # Return metadata about the generated resume
        return {
            'output_path': output_path,
            'filename': output_filename,
            'file_size': file_size,
            'user_id': user_id,
            'job_title': job_title,
            'projects_count': len(projects),
            'certifications_count': len(certifications),
            'generated_at': datetime.now().isoformat(),
            'keywords': keywords
        }

