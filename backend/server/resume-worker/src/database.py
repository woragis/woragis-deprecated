"""
Database operations for resume generation worker.
"""

import logging
import psycopg2
from psycopg2.extras import RealDictCursor
from typing import Dict, List, Optional

logger = logging.getLogger(__name__)


class Database:
    """Database connection handler"""
    
    def __init__(self, database_url: str):
        self.database_url = database_url
        self.conn = None
    
    def connect(self):
        """Establish database connection"""
        try:
            self.conn = psycopg2.connect(self.database_url)
            logger.info("Database connection established")
        except Exception as e:
            logger.error(f"Failed to connect to database: {e}")
            raise
    
    def close(self):
        """Close database connection"""
        if self.conn:
            self.conn.close()
            logger.info("Database connection closed")
    
    def get_user_projects(self, user_id: str, tech_categories: List[str] = None) -> List[Dict]:
        """Get projects for a user, optionally filtered by technology categories"""
        try:
            with self.conn.cursor(cursor_factory=RealDictCursor) as cur:
                query = """
                    SELECT 
                        p.id,
                        p.name,
                        p.description,
                        p.status,
                        p.slug,
                        COALESCE(
                            json_agg(
                                DISTINCT jsonb_build_object(
                                    'name', s.name,
                                    'category', COALESCE(s.category, 'other')
                                )
                            ) FILTER (WHERE s.id IS NOT NULL),
                            '[]'::json
                        ) as technologies
                    FROM projects p
                    LEFT JOIN project_skills ps ON ps.project_id = p.id
                    LEFT JOIN skills s ON s.id = ps.skill_id
                    WHERE p.user_id = %s
                """
                
                params = [user_id]
                
                if tech_categories:
                    query += """
                        AND s.category = ANY(%s)
                    """
                    params.append(tech_categories)
                
                query += """
                    GROUP BY p.id, p.name, p.description, p.status, p.slug
                    ORDER BY p.created_at DESC
                """
                
                cur.execute(query, params)
                projects = cur.fetchall()
                
                # Convert to list of dicts
                result = []
                for project in projects:
                    project_dict = dict(project)
                    # Parse technologies JSON
                    if isinstance(project_dict['technologies'], str):
                        import json
                        project_dict['technologies'] = json.loads(project_dict['technologies'])
                    # Filter out None values
                    if project_dict['technologies']:
                        project_dict['technologies'] = [
                            t for t in project_dict['technologies'] 
                            if t.get('name') is not None
                        ]
                    result.append(project_dict)
                
                return result
        except Exception as e:
            logger.error(f"Error fetching projects: {e}")
            raise
    
    def get_user_certifications(self, user_id: str, categories: List[str] = None) -> List[Dict]:
        """Get certifications for a user, optionally filtered by categories"""
        try:
            with self.conn.cursor(cursor_factory=RealDictCursor) as cur:
                query = """
                    SELECT 
                        id,
                        name,
                        issuer,
                        issue_date,
                        expiry_date,
                        credential_id,
                        verification_url,
                        description,
                        status,
                        category
                    FROM certifications
                    WHERE user_id = %s AND status = 'active'
                """
                
                params = [user_id]
                
                if categories:
                    query += " AND category = ANY(%s)"
                    params.append(categories)
                
                query += " ORDER BY issue_date DESC"
                
                cur.execute(query, params)
                certs = cur.fetchall()
                
                return [dict(cert) for cert in certs]
        except Exception as e:
            logger.error(f"Error fetching certifications: {e}")
            raise
    
    def get_user_info(self, user_id: str) -> Optional[Dict]:
        """Get basic user information"""
        try:
            with self.conn.cursor(cursor_factory=RealDictCursor) as cur:
                query = """
                    SELECT id, email
                    FROM users
                    WHERE id = %s
                """
                cur.execute(query, [user_id])
                user = cur.fetchone()
                if user:
                    user_dict = dict(user)
                    # Extract name from email if name column doesn't exist
                    email = user_dict.get('email', '')
                    if email:
                        name = email.split('@')[0].replace('.', ' ').title()
                        user_dict['name'] = name
                    return user_dict
                return None
        except Exception as e:
            logger.error(f"Error fetching user info: {e}")
            return None

