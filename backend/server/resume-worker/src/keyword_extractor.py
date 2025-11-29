"""
Keyword extraction and category matching for job descriptions.
"""

from typing import Dict, List

# Technology category mapping to job keywords
TECH_CATEGORY_MAPPING = {
    'backend': ['backend', 'back-end', 'server', 'api', 'rest', 'graphql', 'microservices'],
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
    
    return {
        'tech_categories': matching_tech_categories if matching_tech_categories else None,
        'cert_categories': matching_cert_categories if matching_cert_categories else None
    }

