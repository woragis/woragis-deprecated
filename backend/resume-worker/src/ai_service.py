"""
AI Service client for generating resume content.
"""

import logging
import requests
from typing import Dict, List
from circuitbreaker import circuit

logger = logging.getLogger(__name__)

# Circuit breaker configuration for AI Service calls
# Opens after 5 consecutive failures, timeout 30 seconds before half-open
@circuit(failure_threshold=5, recovery_timeout=30, expected_exception=requests.RequestException)
def _call_ai_service(base_url: str, payload: Dict) -> requests.Response:
    """Make HTTP call to AI service with circuit breaker protection"""
    response = requests.post(
        f"{base_url}/v1/chat",
        json=payload,
        timeout=30
    )
    response.raise_for_status()
    return response

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
        context: Dict = None,
        language: str = "en",
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
            elif section_type == 'hard_skills':
                prompt = self._build_hard_skills_prompt(job_description, projects)
            elif section_type == 'tags':
                prompt = self._build_tags_prompt(job_description, projects)
            else:
                prompt = f"Generate a {section_type} section for a resume based on this job description: {job_description}"
            
            # Choose language instruction
            if language.lower().startswith("pt"):
                lang_instruction = "in Brazilian Portuguese (pt-BR). Write naturally, using professional tone in Portuguese."
            else:
                lang_instruction = "in English."

            system_prompt = (
                "You are an expert resume writer. "
                f"Generate professional, ATS-friendly resume content {lang_instruction} "
                "Be concise, use strong action verbs, and quantify achievements when possible. "
                "IMPORTANT: Return ONLY valid HTML code. Do NOT use markdown. Use proper HTML tags like <p>, <strong>, <em>, <ul>, <li>, <h3>, <h4>, <div>, <span> with appropriate CSS classes as specified in the prompt."
            )
            
            # Use circuit breaker protected call
            payload = {
                "agent": "auto",
                "provider": "anthropic",
                "model": "claude-3-5-sonnet-latest",
                "system": system_prompt,
                "input": prompt,
                "temperature": 0.3
            }
            
            try:
                response = _call_ai_service(self.base_url, payload)
            except Exception as e:
                # Circuit breaker will handle failures and open circuit if threshold reached
                logger.error(f"AI service call failed (circuit breaker protected): {e}")
                raise
            data = response.json()
            # Extract content from response
            content = ""
            if 'output' in data:
                content = data['output']
            elif 'message' in data and 'content' in data['message']:
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
- Return ONLY valid HTML code (no markdown, no plain text)
- Use <p> tags for paragraphs
- Use <strong> for emphasis on key terms
- Use <em> for subtle emphasis
- Follow the same 3-paragraph structure as the template
- First paragraph: Passion + specialization + key principles (match job requirements)
- Second paragraph: Education/background + focus areas + commitment to quality
- Third paragraph: Experience highlights + continuous learning + value delivery
- Adapt the technical focus to match the job (e.g., if job mentions Golang, emphasize backend/distributed systems; if frontend, emphasize UI/UX)
- Keep the same professional tone and length
- Use strong action words and technical terminology relevant to the job
- Maintain the structure but personalize content based on job requirements
- Optimize for ATS (Applicant Tracking Systems)

Example HTML format:
<p>Passionate <strong>Full-Stack Developer</strong> specializing in <strong>clean architecture</strong> and <strong>domain-driven design</strong> principles.</p>
<p>Currently pursuing a Computer Science degree with a strong focus on building scalable, maintainable applications across backend, frontend, and mobile platforms.</p>
<p>Experienced in designing and implementing distributed systems, microservices architectures, and high-performance APIs.</p>"""
    
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
- Return ONLY valid HTML code (no markdown, no plain text)
- Use <p> tags for paragraphs
- Use <strong> for emphasis on key terms
- Use <em> for subtle emphasis
- Follow the same 3-paragraph structure as the template
- First paragraph: Passion + specialization + key principles (match job requirements)
- Second paragraph: Education/background + focus areas + commitment to quality
- Third paragraph: Experience highlights + continuous learning + value delivery
- Adapt the technical focus to match the job (e.g., if job mentions Golang, emphasize backend/distributed systems; if frontend, emphasize UI/UX)
- Keep the same professional tone and length
- Use strong action words and technical terminology relevant to the job
- Maintain the structure but personalize content based on job requirements

Example HTML format:
<p>Passionate <strong>Full-Stack Developer</strong> specializing in <strong>clean architecture</strong>, <strong>domain-driven design</strong>, and scalable system development.</p>
<p>Currently pursuing a Computer Science degree with a strong focus on building robust, maintainable applications across backend, frontend, and mobile platforms.</p>
<p>Experienced in designing and implementing distributed systems, microservices architectures, and high-performance APIs.</p>"""
    
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
- Return ONLY valid HTML code (no markdown, no plain text)
- Use <div class="experience-item"> for each project entry
- Use <h3 class="experience-title"> for project names
- Use <ul class="experience-description"> and <li> for bullet points
- Use past tense action verbs
- Quantify results when possible
- Focus on impact and achievements
- Keep each bullet point concise (one line)
- Match the job requirements

Example HTML format:
<div class="experience-item">
  <h3 class="experience-title">Project Name</h3>
  <ul class="experience-description">
    <li>Developed scalable microservices architecture using Go and Docker, reducing response time by 40%</li>
    <li>Implemented RESTful APIs serving 10,000+ requests per day with 99.9% uptime</li>
    <li>Led team of 3 developers in adopting domain-driven design principles</li>
  </ul>
</div>
<div class="experience-item">
  <h3 class="experience-title">Another Project</h3>
  <ul class="experience-description">
    <li>Built responsive frontend using React.js and TypeScript</li>
    <li>Optimized database queries reducing load time by 50%</li>
  </ul>
</div>"""
    
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
- Return ONLY valid HTML code (no markdown, no plain text)
- Use <div class="skill-category"> for each category
- Use <div class="skill-category-title"> or <strong> for category names
- Use <div class="skill-items"> or <span> for skills, separated by " • " (bullet separator)
- Group skills into 4-6 categories (e.g., "Backend Development", "Frontend Development", "Architecture & Design", "Tools & Technologies", "Cloud & DevOps", "Databases")
- Only include skills relevant to the job
- Use standard naming conventions
- Keep categories balanced for 2-column display

Example HTML format:
<div class="skill-category">
  <div class="skill-category-title">Backend Development</div>
  <div class="skill-items">Go • JavaScript/TypeScript • Python • Java</div>
</div>
<div class="skill-category">
  <div class="skill-category-title">Frontend Development</div>
  <div class="skill-items">React.js • Next.js • SvelteKit</div>
</div>
<div class="skill-category">
  <div class="skill-category-title">Architecture & Design</div>
  <div class="skill-items">Domain-Driven Design (DDD) • Clean Architecture • Microservices • RESTful APIs</div>
</div>
<div class="skill-category">
  <div class="skill-category-title">Tools & Technologies</div>
  <div class="skill-items">Docker • Kubernetes • PostgreSQL • Redis</div>
</div>"""
    
    def _build_hard_skills_prompt(self, job_description: str, projects: List[Dict]) -> str:
        """Build prompt for hard skills section - comprehensive technical skills list"""
        all_techs = []
        for p in projects:
            all_techs.extend([t.get('name', '') for t in p.get('technologies', [])])
        
        unique_techs = list(set(all_techs))
        
        return f"""Based on the job description and these technologies, create a comprehensive "Hard Skills" section formatted for a 2-column resume layout with categorized skills.

Job Description:
{job_description}

Technologies Used:
{', '.join(unique_techs[:40])}

Requirements:
- Return ONLY valid HTML code (no markdown, no plain text)
- Use <div class="skill-category"> for each category
- Use <div class="skill-category-title"> or <strong> for category names
- Use <div class="skill-items"> or <span> for skills, separated by " • " (bullet separator)
- Group skills into 6-8 categories such as:
  * Primary Stack (main programming languages and frameworks)
  * Backend Expertise (patterns, architectures, concepts)
  * Databases & Storage
  * Infrastructure & DevOps
  * Additional Technologies
  * Testing & Quality
  * Cloud Platforms (AWS, GCP, Azure if relevant)
  * Other relevant categories based on job requirements
- Only include skills relevant to the job description
- Use standard naming conventions
- Keep categories balanced for 2-column display
- Be comprehensive but focused on job requirements

Example HTML format:
<div class="skill-category">
  <div class="skill-category-title">Primary Stack</div>
  <div class="skill-items">Go (Golang) • Microservices Architecture • RESTful & gRPC APIs • Event-Driven Systems</div>
</div>
<div class="skill-category">
  <div class="skill-category-title">Backend Expertise</div>
  <div class="skill-items">Concurrent Programming • Design Patterns • Domain-Driven Design (DDD) • Clean Architecture • CQRS • Message Queues</div>
</div>
<div class="skill-category">
  <div class="skill-category-title">Databases & Storage</div>
  <div class="skill-items">PostgreSQL • MySQL • MongoDB • Redis • Database Design & Optimization</div>
</div>
<div class="skill-category">
  <div class="skill-category-title">Infrastructure & DevOps</div>
  <div class="skill-items">Docker • Kubernetes • CI/CD Pipelines • Git • Linux/Unix</div>
</div>
<div class="skill-category">
  <div class="skill-category-title">Additional Technologies</div>
  <div class="skill-items">JavaScript/TypeScript • Python • Java • GraphQL</div>
</div>
<div class="skill-category">
  <div class="skill-category-title">Testing & Quality</div>
  <div class="skill-items">Unit Testing • Integration Testing • Test-Driven Development (TDD) • Performance Testing</div>
</div>"""
    
    def _build_tags_prompt(self, job_description: str, projects: List[Dict]) -> str:
        """Build prompt for generating tags (technologies, languages, tools)"""
        all_techs = []
        for p in projects:
            all_techs.extend([t.get('name', '') for t in p.get('technologies', [])])
        
        unique_techs = list(set(all_techs))
        
        return f"""Based on the job description and projects, extract and generate a list of relevant tags (technologies, programming languages, tools, frameworks, cloud services, etc.) for this resume.

Job Description:
{job_description}

Technologies from Projects:
{', '.join(unique_techs[:50])}

Requirements:
- Extract 5-10 relevant tags from the job description and projects
- Include programming languages (e.g., golang, python, javascript, typescript)
- Include frameworks and libraries (e.g., react, svelte, gin, echo)
- Include cloud services and infrastructure (e.g., aws, docker, kubernetes, eks, ec2, s3)
- Include databases (e.g., postgresql, mongodb, redis)
- Include tools and platforms (e.g., git, github, ci/cd, terraform)
- Use lowercase, no spaces (use hyphens if needed: e.g., "ci-cd" not "CI/CD")
- Return ONLY a comma-separated list of tags, nothing else
- Example format: golang, python, aws, docker, kubernetes, postgresql, react, typescript, redis, eks
- Focus on the most relevant and commonly used technologies"""
    
    def generate_tags(self, job_description: str, projects: List[Dict], language: str = "en") -> List[str]:
        """Generate tags using AI"""
        try:
            prompt = self._build_tags_prompt(job_description, projects)
            
            # Use circuit breaker protected call
            payload = {
                "agent": "auto",
                "provider": "anthropic",
                "model": "claude-3-5-sonnet-latest",
                "system": "You are a technical tag extractor. Extract relevant technology tags from the job description and projects. Return ONLY a comma-separated list of lowercase tags, nothing else.",
                "input": prompt,
                "temperature": 0.2
            }
            
            try:
                response = _call_ai_service(self.base_url, payload)
            except Exception as e:
                # Circuit breaker will handle failures and open circuit if threshold reached
                logger.error(f"AI service call failed (circuit breaker protected): {e}")
                raise
            data = response.json()

            # Extract content
            content = ""
            if 'output' in data:
                content = data['output']
            elif 'message' in data and 'content' in data['message']:
                content = data['message']['content']
            elif 'choices' in data and len(data['choices']) > 0:
                content = data['choices'][0]['message']['content']
            else:
                logger.warning(f"Unexpected AI response format for tags: {data}")
                return []
            
            # Parse comma-separated tags
            tags = [tag.strip().lower() for tag in content.split(",") if tag.strip()]
            # Remove duplicates and limit to 10
            unique_tags = list(dict.fromkeys(tags))[:10]
            return unique_tags
                
        except Exception as e:
            logger.error(f"Error generating tags: {e}")
            # Fallback: extract from technologies
            all_techs = []
            for p in projects:
                all_techs.extend([t.get('name', '').lower() for t in p.get('technologies', [])])
            unique_techs = list(dict.fromkeys(all_techs))[:10]
            return unique_techs

