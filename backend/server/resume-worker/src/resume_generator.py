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
from ai_service import AIService, PROFILE_TEMPLATE
from keyword_extractor import extract_keywords
from translation_helper import TranslationHelper

logger = logging.getLogger(__name__)


class ResumeGenerator:
    """Main resume generator"""
    
    def __init__(self, db: Database, ai_service: AIService, output_dir: str = "/app/output", translation_helper: TranslationHelper = None):
        self.db = db
        self.ai_service = ai_service
        self.output_dir = output_dir
        self.translation_helper = translation_helper
        
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
    
    def _parse_hard_skills(self, skills_text: str) -> List[Dict]:
        """Parse hard skills text into structured format for template rendering"""
        if not skills_text or not skills_text.strip():
            return []
        
        categories = []
        current_category = None
        
        lines = skills_text.split('\n')
        i = 0
        while i < len(lines):
            line = lines[i].strip()
            if not line:
                i += 1
                continue
            
            # Check if line is a category title (no bullet points)
            if '•' not in line:
                # Save previous category if it has skills
                if current_category and current_category.get('skills'):
                    categories.append(current_category)
                
                # Start new category
                current_category = {
                    'title': line,
                    'skills': []
                }
                i += 1
                
                # Check if next line has skills
                if i < len(lines):
                    next_line = lines[i].strip()
                    if '•' in next_line:
                        # This is a skills line - split by bullet points
                        skills = [s.strip() for s in next_line.split('•') if s.strip()]
                        current_category['skills'].extend(skills)
                        i += 1
            else:
                # This is a skills line without a category - skip or handle
                i += 1
        
        # Add last category
        if current_category and current_category.get('skills'):
            categories.append(current_category)
        
        return categories
    
    def generate_resume(
        self,
        user_id: str,
        job_description: str,
        job_title: str = "Software Engineer",
        output_filename: str = None,
        language: str = "en"
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
        
        # Apply translations to projects if available
        if self.translation_helper:
            projects = self.translation_helper.translate_projects(projects, language)
        
        # Fetch relevant certifications
        certifications = self.db.get_user_certifications(
            user_id,
            categories=keywords['cert_categories']
        )
        logger.info(f"Found {len(certifications)} relevant certifications")
        
        # Apply translations to certifications if available
        if self.translation_helper:
            certifications = self.translation_helper.translate_certifications(certifications, language)
        
        # Fetch publications (posts)
        publications = self.db.get_user_publications(user_id, limit=4)
        logger.info(f"Found {len(publications)} publications")
        
        # Apply translations to publications if available
        if self.translation_helper:
            publications = self.translation_helper.translate_posts(publications, language)
        
        # Format publication dates for template
        for pub in publications:
            if pub.get('published_at'):
                try:
                    if hasattr(pub['published_at'], 'strftime'):
                        pub['formatted_date'] = pub['published_at'].strftime('%B %Y')
                    else:
                        pub['formatted_date'] = str(pub['published_at'])
                except Exception as e:
                    logger.warning(f"Error formatting publication date: {e}")
                    pub['formatted_date'] = str(pub.get('published_at', ''))
            else:
                pub['formatted_date'] = None
        
        # Fetch work experiences from database
        work_experiences_db = self.db.get_user_experiences(user_id)
        
        # Apply translations to experiences if available
        if self.translation_helper and work_experiences_db:
            work_experiences_db = self.translation_helper.translate_experiences(work_experiences_db, language)
        
        # Convert database experiences to template format
        work_experience = []
        for exp in work_experiences_db:
            work_exp = {
                'title': exp.get('position', ''),
                'company': exp.get('company', ''),
                'location': exp.get('location', ''),
                'period': exp.get('period', ''),
                'description': exp.get('description', [])
            }
            work_experience.append(work_exp)
        
        # Hardcoded Education section (can be moved to database later)
        education = [
            {
                'degree': 'Bachelor of Science in Computer Science',
                'institution': 'Federal University of Paraíba (UFPB)',
                'location': 'João Pessoa, PB, Brazil',
                'period': '2021 - Present',
                'status': 'In Progress',
                'relevant_coursework': [
                    'Data Structures and Algorithms',
                    'Software Engineering',
                    'Database Systems',
                    'Distributed Systems',
                    'Computer Networks'
                ]
            }
        ]
        
        # Translate education status if translation helper is available
        if self.translation_helper:
            for edu in education:
                if 'status' in edu:
                    edu['status'] = self.translation_helper.translate_education_status(edu['status'], language)
        
        # Hardcoded Volunteer Work section (can be moved to database later)
        volunteer_work = [
            {
                'title': 'Open Source Contributor',
                'organization': 'Various Open Source Projects',
                'location': 'Remote',
                'period': '2021 - Present',
                'description': [
                    'Contributed to open source projects on GitHub',
                    'Helped improve documentation and fix bugs in community projects',
                    'Shared knowledge through technical blog posts and tutorials'
                ]
            },
            {
                'title': 'Technical Mentor',
                'organization': 'Programming Communities',
                'location': 'Online',
                'period': '2022 - Present',
                'description': [
                    'Mentored junior developers in software engineering best practices',
                    'Conducted code reviews and provided constructive feedback',
                    'Organized technical workshops and knowledge-sharing sessions'
                ]
            }
        ]
        
        # Generate resume sections using AI
        logger.info("Generating resume sections with AI...")
        profile = self.ai_service.generate_resume_section('profile', job_description, projects, language=language)
        about = self.ai_service.generate_resume_section('about', job_description, projects, language=language)
        experience = self.ai_service.generate_resume_section('experience', job_description, projects, language=language)
        skills = self.ai_service.generate_resume_section('skills', job_description, projects, language=language)
        hard_skills = self.ai_service.generate_resume_section('hard_skills', job_description, projects, language=language)
        logger.info(f"Hard skills generated: {len(hard_skills) if hard_skills else 0} characters")
        
        # Generate tags using AI
        logger.info("Generating tags...")
        tags = self.ai_service.generate_tags(job_description, projects, language=language)
        logger.info(f"Generated {len(tags)} tags: {', '.join(tags)}")

        # Ensure profile is never empty - use fallback if AI returns empty
        if not profile or not profile.strip():
            logger.warning("Profile generation returned empty, using fallback")
            profile = PROFILE_TEMPLATE

        # Parse hard_skills text into structured format for 2-column grid display,
        # with a deterministic fallback when the AI section is empty or on error.
        if not hard_skills or not hard_skills.strip():
            logger.warning("Hard skills generation returned empty, using fallback hard skills template")
            fallback_skills = (
                "Primary Stack\n"
                "Go (Golang) • Microservices Architecture • RESTful & gRPC APIs • Event-Driven Systems\n\n"
                "Backend Expertise\n"
                "Concurrent Programming • Design Patterns • Domain-Driven Design (DDD) • Clean Architecture • CQRS • Message Queues\n\n"
                "Databases & Storage\n"
                "PostgreSQL • MySQL • MongoDB • Redis • Database Design & Optimization\n\n"
                "Infrastructure & DevOps\n"
                "Docker • Kubernetes • CI/CD Pipelines • Git • Linux/Unix\n\n"
                "Additional Technologies\n"
                "JavaScript/TypeScript • Python • Java • GraphQL\n\n"
                "Testing & Quality\n"
                "Unit Testing • Integration Testing • Test-Driven Development (TDD) • Performance Testing"
            )
            hard_skills_parsed = self._parse_hard_skills(fallback_skills)
        else:
            hard_skills_parsed = self._parse_hard_skills(hard_skills)
        
        # Choose template based on language (English vs Brazilian Portuguese)
        if language.lower().startswith("pt"):
            template_name = 'resume_pt.html'
        else:
            template_name = 'resume.html'
        
        # Translate section headers if translation helper is available
        section_headers = {
            'profile': 'Profile',
            'education': 'Education',
            'work_experience': 'Work Experience',
            'volunteer_work': 'Volunteer Work',
            'projects': 'Projects',
            'experience': 'Experience',
            'certifications': 'Certifications',
            'publications': 'Publications',
            'languages': 'Languages',
            'status_label': 'Status',
            'relevant_coursework': 'Relevant Coursework'
        }
        if self.translation_helper:
            for key, header in section_headers.items():
                section_headers[key] = self.translation_helper.translate_section_header(header, language)

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
        
        # Limit certifications to 6
        certifications_to_show = certifications[:6] if certifications else []
        
        # Limit publications to 4
        publications_to_show = publications[:4] if publications else []
        
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
            'hard_skills': hard_skills_parsed,
            'projects': projects_to_show,
            'certifications': certifications_to_show,
            'education': education,
            'work_experience': work_experience,
            'volunteer_work': volunteer_work,
            'publications': publications_to_show,
            'job_title': job_title,
            'section_headers': section_headers
        }
        
        # Render HTML template
        template = self.template_env.get_template(template_name)
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
            'keywords': keywords,
            'tags': tags
        }

