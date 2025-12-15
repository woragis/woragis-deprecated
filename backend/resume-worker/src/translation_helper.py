"""
Translation helper for resume generation.
Fetches translations for entities and merges them with original data.
"""

import logging
import psycopg2
from psycopg2.extras import RealDictCursor
from typing import Dict, List, Optional, Any
import json

logger = logging.getLogger(__name__)


class TranslationHelper:
    """Helper class for fetching and applying translations to resume data"""
    
    def __init__(self, database_url: str):
        self.database_url = database_url
        self.conn = None
    
    def connect(self):
        """Establish database connection"""
        try:
            self.conn = psycopg2.connect(self.database_url)
            logger.info("Translation helper database connection established")
        except Exception as e:
            logger.error(f"Failed to connect to database: {e}")
            raise
    
    def close(self):
        """Close database connection"""
        if self.conn:
            self.conn.close()
    
    def get_translation(self, entity_type: str, entity_id: str, language: str) -> Optional[Dict[str, str]]:
        """
        Get translation fields for an entity.
        Returns a dict of field_name -> translated_value mapping.
        """
        if not self.conn:
            return None
        
        try:
            with self.conn.cursor(cursor_factory=RealDictCursor) as cur:
                # Map language codes (e.g., "pt" -> "pt-BR")
                lang_map = {
                    "pt": "pt-BR",
                    "en": "en",
                    "es": "es",
                    "fr": "fr",
                    "de": "de",
                    "ru": "ru",
                    "ja": "ja",
                    "ko": "ko",
                    "zh": "zh-CN",
                    "el": "el",
                    "la": "la"
                }
                mapped_lang = lang_map.get(language.lower(), language)
                
                query = """
                    SELECT fields, status
                    FROM translations
                    WHERE entity_type = %s 
                      AND entity_id = %s 
                      AND language = %s
                      AND status = 'completed'
                """
                cur.execute(query, (entity_type, entity_id, mapped_lang))
                result = cur.fetchone()
                
                if result:
                    fields = result['fields']
                    if isinstance(fields, str):
                        return json.loads(fields)
                    return fields
                return None
        except Exception as e:
            logger.warning(f"Error fetching translation for {entity_type} {entity_id}: {e}")
            return None
    
    def apply_translation_to_project(self, project: Dict, language: str) -> Dict:
        """Apply translation to a project if available"""
        if not project or not project.get('id'):
            return project
        
        translation = self.get_translation('project', str(project['id']), language)
        if translation:
            # Translatable fields: name, description
            if 'name' in translation:
                project['name'] = translation['name']
            if 'description' in translation:
                project['description'] = translation['description']
        
        return project
    
    def apply_translation_to_certification(self, cert: Dict, language: str) -> Dict:
        """Apply translation to a certification if available"""
        if not cert or not cert.get('id'):
            return cert
        
        translation = self.get_translation('certification', str(cert['id']), language)
        if translation:
            # Translatable fields: name, category (as text), description
            if 'name' in translation:
                cert['name'] = translation['name']
            if 'category' in translation:
                cert['category'] = translation['category']
            if 'description' in translation:
                cert['description'] = translation['description']
        
        return cert
    
    def apply_translation_to_experience(self, exp: Dict, language: str) -> Dict:
        """Apply translation to an experience if available"""
        if not exp or not exp.get('id'):
            return exp
        
        translation = self.get_translation('experience', str(exp['id']), language)
        if translation:
            # Translatable fields: position, description
            if 'position' in translation:
                exp['position'] = translation['position']
            if 'description' in translation:
                exp['description'] = translation['description']
        
        return exp
    
    def apply_translation_to_post(self, post: Dict, language: str) -> Dict:
        """Apply translation to a post/publication if available"""
        if not post or not post.get('id'):
            return post
        
        translation = self.get_translation('post', str(post['id']), language)
        if translation:
            # Translatable fields: title, excerpt, content
            if 'title' in translation:
                post['title'] = translation['title']
            if 'excerpt' in translation:
                post['excerpt'] = translation['excerpt']
            if 'content' in translation:
                post['content'] = translation['content']
        
        return post
    
    def translate_projects(self, projects: List[Dict], language: str) -> List[Dict]:
        """Apply translations to a list of projects"""
        return [self.apply_translation_to_project(p, language) for p in projects]
    
    def translate_certifications(self, certifications: List[Dict], language: str) -> List[Dict]:
        """Apply translations to a list of certifications"""
        return [self.apply_translation_to_certification(c, language) for c in certifications]
    
    def translate_experiences(self, experiences: List[Dict], language: str) -> List[Dict]:
        """Apply translations to a list of experiences"""
        return [self.apply_translation_to_experience(e, language) for e in experiences]
    
    def translate_posts(self, posts: List[Dict], language: str) -> List[Dict]:
        """Apply translations to a list of posts"""
        return [self.apply_translation_to_post(p, language) for p in posts]
    
    def translate_education_status(self, status: str, language: str) -> str:
        """Translate education status values"""
        if language.lower().startswith('pt'):
            status_map = {
                'In Progress': 'Em Andamento',
                'Completed': 'Concluído',
                'On Hold': 'Em Pausa',
                'Dropped': 'Abandonado'
            }
            return status_map.get(status, status)
        return status
    
    def translate_section_header(self, header: str, language: str) -> str:
        """Translate section headers"""
        if language.lower().startswith('pt'):
            header_map = {
                'Profile': 'Perfil',
                'Education': 'Educação',
                'Work Experience': 'Experiência Profissional',
                'Volunteer Work': 'Trabalho Voluntário',
                'Projects': 'Projetos',
                'Experience': 'Experiência',
                'Certifications': 'Certificações',
                'Publications': 'Publicações',
                'Languages': 'Idiomas',
                'Relevant Coursework': 'Disciplinas Relevantes',
                'Status': 'Status'
            }
            return header_map.get(header, header)
        return header

