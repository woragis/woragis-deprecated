"""
AI Service client for generating resume content.
"""

import logging
import requests
from typing import Dict, List

logger = logging.getLogger(__name__)


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
            if 'message' in data and 'content' in data['message']:
                return data['message']['content'].strip()
            elif 'choices' in data and len(data['choices']) > 0:
                return data['choices'][0]['message']['content'].strip()
            else:
                logger.warning(f"Unexpected AI response format: {data}")
                return ""
                
        except Exception as e:
            logger.error(f"Error generating {section_type} section: {e}")
            return ""
    
    def _build_profile_prompt(self, job_description: str, projects: List[Dict]) -> str:
        """Build prompt for professional summary/profile section"""
        projects_text = "\n".join([
            f"- {p['name']}: {p.get('description', '')}"
            for p in projects[:10]  # Limit to top 10
        ])
        
        return f"""Write a professional summary (4-5 lines) for a resume optimized for this job description:

Job Description:
{job_description}

Relevant Projects:
{projects_text}

Requirements:
- Write in English
- Use strong action verbs
- Highlight relevant technical skills
- Be concise and impactful
- Optimize for ATS (Applicant Tracking Systems)
- Focus on achievements and impact"""
    
    def _build_about_prompt(self, job_description: str, projects: List[Dict]) -> str:
        """Build prompt for about me section"""
        return f"""Write a brief "About Me" section (3-4 sentences) for a resume based on this job description:

{job_description}

Focus on:
- Professional background
- Key strengths
- Passion for the field
- What makes you unique"""
    
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
        """Build prompt for technical skills section"""
        all_techs = []
        for p in projects:
            all_techs.extend([t.get('name', '') for t in p.get('technologies', [])])
        
        unique_techs = list(set(all_techs))
        
        return f"""Based on the job description and these technologies, create a categorized technical skills list:

Job Description:
{job_description}

Technologies Used:
{', '.join(unique_techs[:30])}

Requirements:
- Group skills into categories (e.g., Languages, Frameworks, Tools, Cloud)
- Only include skills relevant to the job
- Use standard naming conventions
- Format as a clean list"""

