"""
AI Service client for generating resume content.
"""

import logging
import requests
from typing import Dict, List

logger = logging.getLogger(__name__)

# Profile template - base structure for AI to follow
PROFILE_TEMPLATE = """Passionate Full-Stack Developer specializing in clean architecture and domain-driven design principles.

Currently pursuing a Computer Science degree with a strong focus on building scalable, maintainable applications across backend, frontend, and mobile platforms. Committed to writing elegant code that reflects real-world business domains and delivers measurable business value.

Experienced in designing and implementing distributed systems, microservices architectures, and high-performance APIs. Continuously learning and applying modern software engineering principles to create solutions that balance technical excellence with practical business outcomes."""

# About Me template - base structure for AI to follow
ABOUT_ME_TEMPLATE = """Passionate Full-Stack Developer specializing in clean architecture, domain-driven design, and scalable system development.

Currently pursuing a Computer Science degree with a strong focus on building robust, maintainable applications across backend, frontend, and mobile platforms. Committed to writing elegant, production-ready code that reflects real-world business domains and solves complex technical challenges.

Experienced in designing and implementing distributed systems, microservices architectures, and high-performance APIs. Continuously learning and applying modern software engineering principles to deliver solutions that balance technical excellence with business value."""


class AIService:
    """AI Service client for generating resume content"""
    
    def __init__(self, base_url: str):
        self.base_url = base_url.rstrip('/')
    
    def generate_resume_section(
        self,
        section_type: str,
        job_description: str,
        projects: List[Dict],
        certifications: List[Dict] = None,
        context: Dict = None
    ) -> str:
        """Generate a resume section using AI"""
        try:
            # Prepare prompt based on section type
            if section_type == 'profile':
                prompt = self._build_profile_prompt(job_description, projects)
            elif section_type == 'about':
                prompt = self._build_about_prompt(job_description, projects)
            elif section_type == 'experience':
                prompt = self._build_experience_prompt(job_description, projects)
            elif section_type == 'skills':
                prompt = self._build_skills_prompt(job_description, projects)
            else:
                prompt = f"Generate a {section_type} section for a resume based on this job description: {job_description}"
            
            # Call AI service
            response = requests.post(
                f"{self.base_url}/api/chat/completions",
                json={
                    "provider": "anthropic",  # Use Anthropic for better resume writing
                    "model": "claude-3-5-sonnet-latest",
                    "messages": [
                        {
                            "role": "user",
                            "content": f"You are an expert resume writer. Generate professional, ATS-friendly resume content in English. Be concise, use strong action verbs, and quantify achievements when possible.\n\n{prompt}"
                        }
                    ],
                    "temperature": 0.3,
                    "max_tokens": 1000
                },
                timeout=30
            )
            
            response.raise_for_status()
            data = response.json()
            
            # Extract content from response (handle both response formats)
            content = ""
            if 'message' in data and 'content' in data['message']:
                content = data['message']['content']
            elif 'choices' in data and len(data['choices']) > 0:
                content = data['choices'][0]['message']['content']
            else:
                logger.warning(f"Unexpected AI response format: {data}")
                return ""
            
            # Clean up content: remove extra whitespace but preserve normal spacing
            # Replace any non-breaking spaces or special Unicode spaces with regular spaces
            content = content.replace('\u00A0', ' ')  # Non-breaking space
            content = content.replace('\u2009', ' ')   # Thin space
            content = content.replace('\u2006', ' ')   # Six-per-em space
            content = content.replace('\u2007', ' ')   # Figure space
            content = content.replace('\u2008', ' ')   # Punctuation space
            # Normalize multiple spaces to single space, but preserve line breaks
            import re
            content = re.sub(r' +', ' ', content)  # Multiple spaces to single
            content = re.sub(r'\n\s*\n', '\n\n', content)  # Multiple newlines to double
            return content.strip()
                
        except Exception as e:
            logger.error(f"Error generating {section_type} section: {e}")
            return ""
    
    def _build_profile_prompt(self, job_description: str, projects: List[Dict]) -> str:
        """Build prompt for professional summary/profile section - uses base template"""
        projects_text = "\n".join([
            f"- {p['name']}: {p.get('description', '')}"
            for p in projects[:10]  # Limit to top 10
        ])
        
        return f"""Write a personalized "Profile" section (3-4 sentences) for a resume based on this job description. Follow this structure and style:

BASE TEMPLATE STRUCTURE:
{PROFILE_TEMPLATE}

Job Description:
{job_description}

Relevant Projects:
{projects_text}

Requirements:
- Follow the same 3-paragraph structure as the template
- First paragraph: Passion + specialization + key principles (match job requirements)
- Second paragraph: Education/background + focus areas + commitment to quality
- Third paragraph: Experience highlights + continuous learning + value delivery
- Adapt the technical focus to match the job (e.g., if job mentions Golang, emphasize backend/distributed systems; if frontend, emphasize UI/UX)
- Keep the same professional tone and length
- Use strong action words and technical terminology relevant to the job
- Maintain the structure but personalize content based on job requirements
- Write in English
- Optimize for ATS (Applicant Tracking Systems)"""
    
    def _build_about_prompt(self, job_description: str, projects: List[Dict]) -> str:
        """Build prompt for about me section - follows a specific template structure"""
        projects_summary = ", ".join([p['name'] for p in projects[:5]])
        
        return f"""Write a personalized "About Me" section (3-4 sentences) for a resume based on this job description. Follow this structure and style:

BASE TEMPLATE STRUCTURE:
{ABOUT_ME_TEMPLATE}

Job Description:
{job_description}

Relevant Projects: {projects_summary}

Requirements:
- Follow the same 3-paragraph structure as the template
- First paragraph: Passion + specialization + key principles (match job requirements)
- Second paragraph: Education/background + focus areas + commitment to quality
- Third paragraph: Experience highlights + continuous learning + value delivery
- Adapt the technical focus to match the job (e.g., if job mentions Golang, emphasize backend/distributed systems; if frontend, emphasize UI/UX)
- Keep the same professional tone and length
- Use strong action words and technical terminology relevant to the job
- Maintain the structure but personalize content based on job requirements"""
    
    def _build_experience_prompt(self, job_description: str, projects: List[Dict]) -> str:
        """Build prompt for experience/projects section"""
        projects_text = "\n\n".join([
            f"Project: {p['name']}\nDescription: {p.get('description', '')}\nTechnologies: {', '.join([t.get('name', '') for t in p.get('technologies', [])])}"
            for p in projects[:8]  # Top 8 projects
        ])
        
        return f"""Format these projects as resume experience entries. For each project, write 2-3 bullet points highlighting achievements and impact.

Job Description:
{job_description}

Projects:
{projects_text}

Requirements:
- Use past tense action verbs
- Quantify results when possible
- Focus on impact and achievements
- Keep each bullet point concise (one line)
- Match the job requirements
- Format as plain text with project names as headers followed by bullet points"""
    
    def _build_skills_prompt(self, job_description: str, projects: List[Dict]) -> str:
        """Build prompt for technical skills section - formatted for 2-column display"""
        all_techs = []
        for p in projects:
            all_techs.extend([t.get('name', '') for t in p.get('technologies', [])])
        
        unique_techs = list(set(all_techs))
        
        return f"""Based on the job description and these technologies, create a categorized technical skills list formatted for a 2-column resume layout.

Job Description:
{job_description}

Technologies Used:
{', '.join(unique_techs[:30])}

Requirements:
- Group skills into 4-6 categories (e.g., "Backend Development", "Frontend Development", "Architecture & Design", "Tools & Technologies", "Cloud & DevOps", "Databases")
- Format each category as: "Category Name" (bold) followed by skills separated by " • " (bullet points)
- Only include skills relevant to the job
- Use standard naming conventions
- Keep categories balanced for 2-column display
- Format example:
Backend Development
Go • JavaScript/TypeScript • Python • Java

Frontend Development
React.js • Next.js • SvelteKit

Architecture & Design
Domain-Driven Design (DDD) • Clean Architecture • Microservices • RESTful APIs

Tools & Technologies
Docker • Kubernetes • PostgreSQL • Redis"""

