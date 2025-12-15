"""
Keyword extraction and category matching for job descriptions.
"""

from typing import Dict, List

# Technology category mapping to job keywords
# Maps job description keywords to skill categories in the database
TECH_CATEGORY_MAPPING = {
    'backend': ['backend', 'back-end', 'server', 'api', 'rest', 'graphql', 'microservices', 'golang', 'go programming', 'go language'],
    'language': ['golang', 'go programming', 'go language', 'python', 'java', 'javascript', 'typescript', 'rust', 'programming language'],
    'devops': ['devops', 'dev-ops', 'ci/cd', 'docker', 'kubernetes', 'k8s', 'terraform', 'ansible', 'jenkins'],
    'infrastructure': ['infrastructure', 'cloud', 'aws', 'azure', 'gcp', 'infrastructure'],
    'database': ['database', 'sql', 'nosql', 'postgresql', 'mysql', 'mongodb', 'redis'],
    'frontend': ['frontend', 'front-end', 'react', 'vue', 'angular', 'javascript', 'typescript'],
    'monitoring': ['monitoring', 'observability', 'prometheus', 'grafana', 'datadog'],
    'testing': ['testing', 'qa', 'test', 'tdd', 'bdd'],
    'other': []
}

# Certification category mapping
CERT_CATEGORY_MAPPING = {
    'cloud': ['cloud', 'aws', 'azure', 'gcp', 'google cloud'],
    'devops': ['devops', 'kubernetes', 'docker', 'terraform'],
    'security': ['security', 'pentesting', 'penetration testing', 'cybersecurity', 'ethical hacking'],
    'programming': ['programming', 'python', 'golang', 'java', 'spring boot'],
    'database': ['database', 'sql'],
    'architecture': ['architecture', 'solution architect'],
    'other': []
}


def extract_keywords(job_description: str) -> Dict[str, List[str]]:
    """Extract relevant keywords from job description"""
    job_lower = job_description.lower()
    
    # Find matching tech categories
    matching_tech_categories = []
    for category, keywords in TECH_CATEGORY_MAPPING.items():
        if any(keyword in job_lower for keyword in keywords):
            matching_tech_categories.append(category)
    
    # Find matching cert categories
    matching_cert_categories = []
    for category, keywords in CERT_CATEGORY_MAPPING.items():
        if any(keyword in job_lower for keyword in keywords):
            matching_cert_categories.append(category)
    
    # Extract specific skill names mentioned in job description
    skill_names = []
    skill_keywords = ['golang', 'go programming', 'python', 'java', 'javascript', 'typescript', 'rust', 
                      'kubernetes', 'docker', 'postgresql', 'redis', 'grpc', 'websocket']
    for skill in skill_keywords:
        if skill in job_lower:
            # Map variations to actual skill names
            if skill in ['golang', 'go programming', 'go language']:
                skill_names.append('Golang')
            elif skill == 'python':
                skill_names.append('Python')
            elif skill == 'java':
                skill_names.append('Java')
            elif skill == 'javascript':
                skill_names.append('JavaScript')
            elif skill == 'typescript':
                skill_names.append('TypeScript')
            elif skill == 'rust':
                skill_names.append('Rust')
            elif skill == 'kubernetes':
                skill_names.append('Kubernetes')
            elif skill == 'docker':
                skill_names.append('Docker')
            elif skill == 'postgresql':
                skill_names.append('PostgreSQL')
            elif skill == 'redis':
                skill_names.append('Redis')
            elif skill == 'grpc':
                skill_names.append('gRPC')
            elif skill == 'websocket':
                skill_names.append('WebSocket')
    
    return {
        'tech_categories': matching_tech_categories if matching_tech_categories else None,
        'cert_categories': matching_cert_categories if matching_cert_categories else None,
        'skill_names': skill_names if skill_names else None
    }

