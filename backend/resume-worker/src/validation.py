"""
Validation utilities for resume worker.
Validates inputs, job messages, and generated content.
"""

import re
import uuid
from typing import Dict, Any, Optional


def validate_uuid(value: str, field_name: str = "id") -> None:
    """Validate UUID format"""
    if not value:
        raise ValueError(f"{field_name} is required")
    try:
        uuid.UUID(value)
    except ValueError:
        raise ValueError(f"{field_name} must be a valid UUID")


def validate_string(value: str, min_length: int, max_length: int, field_name: str) -> None:
    """Validate string length"""
    if not value:
        if min_length > 0:
            raise ValueError(f"{field_name} is required")
        return
    if len(value) < min_length:
        raise ValueError(f"{field_name} is too short (minimum {min_length} characters)")
    if len(value) > max_length:
        raise ValueError(f"{field_name} is too long (maximum {max_length} characters)")


def validate_email(email: str) -> None:
    """Validate email format"""
    if not email:
        raise ValueError("email is required")
    pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    if not re.match(pattern, email):
        raise ValueError("invalid email format")


def validate_url(url: str, field_name: str = "url") -> None:
    """Validate URL format"""
    if not url:
        raise ValueError(f"{field_name} is required")
    pattern = r'^https?://[^\s/$.?#].[^\s]*$'
    if not re.match(pattern, url):
        raise ValueError(f"{field_name} must be a valid HTTP/HTTPS URL")


def validate_language_code(language: str) -> None:
    """Validate ISO 639-1 language code"""
    if not language:
        raise ValueError("language is required")
    if len(language) != 2:
        raise ValueError("language must be exactly 2 characters (ISO 639-1 code)")
    if not language.isalpha():
        raise ValueError("language must contain only letters")


def validate_no_sql_injection(value: str, field_name: str) -> None:
    """Check for potential SQL injection patterns"""
    if not value:
        return
    dangerous_patterns = [
        r"(\b(SELECT|INSERT|UPDATE|DELETE|DROP|CREATE|ALTER|EXEC|EXECUTE)\b)",
        r"(--|#|/\*|\*/)",
        r"(\b(UNION|OR|AND)\s+\d+\s*=\s*\d+)",
        r"('|(\\')|(--)|(;)|(\|)|(\*))",
    ]
    for pattern in dangerous_patterns:
        if re.search(pattern, value, re.IGNORECASE):
            raise ValueError(f"{field_name} contains potentially dangerous content")


def validate_no_xss(value: str, field_name: str) -> None:
    """Check for potential XSS patterns"""
    if not value:
        return
    dangerous_patterns = [
        r"<script[^>]*>",
        r"javascript:",
        r"on\w+\s*=",
        r"<iframe[^>]*>",
        r"<object[^>]*>",
        r"<embed[^>]*>",
    ]
    for pattern in dangerous_patterns:
        if re.search(pattern, value, re.IGNORECASE):
            raise ValueError(f"{field_name} contains potentially dangerous content")


def validate_resume_job_message(message: Dict[str, Any]) -> None:
    """Validate resume generation job message from queue"""
    # Validate user_id (required, UUID)
    user_id = message.get('user_id')
    if not user_id:
        raise ValueError("user_id is required")
    validate_uuid(user_id, "user_id")
    
    # Validate job_description (required, 10-50000 chars)
    job_description = message.get('job_description')
    if not job_description:
        raise ValueError("job_description is required")
    validate_string(job_description, 10, 50000, "job_description")
    validate_no_sql_injection(job_description, "job_description")
    validate_no_xss(job_description, "job_description")
    
    # Validate job_title (optional, but if provided, validate)
    job_title = message.get('job_title')
    if job_title:
        validate_string(job_title, 1, 200, "job_title")
        validate_no_sql_injection(job_title, "job_title")
        validate_no_xss(job_title, "job_title")
    
    # Validate output_filename (optional, but if provided, validate)
    output_filename = message.get('output_filename')
    if output_filename:
        validate_string(output_filename, 1, 255, "output_filename")
        # Check for path traversal
        if '..' in output_filename or '/' in output_filename or '\\' in output_filename:
            raise ValueError("output_filename contains invalid characters")
    
    # Validate language (optional, but if provided, validate)
    language = message.get('language', 'en')
    validate_language_code(language)
    
    # Validate job_id if present
    job_id = message.get('id')
    if job_id:
        validate_uuid(job_id, "id")


def validate_cli_arguments(args: list) -> None:
    """Validate command-line arguments"""
    if len(args) < 3:
        raise ValueError("At least user_id and job_description are required")
    
    user_id = args[1]
    validate_uuid(user_id, "user_id")
    
    job_description = args[2]
    validate_string(job_description, 10, 50000, "job_description")
    validate_no_sql_injection(job_description, "job_description")
    validate_no_xss(job_description, "job_description")
    
    # Validate optional arguments
    if len(args) > 3:
        job_title = args[3]
        if job_title:
            validate_string(job_title, 1, 200, "job_title")
            validate_no_sql_injection(job_title, "job_title")
            validate_no_xss(job_title, "job_title")
    
    if len(args) > 4:
        output_filename = args[4]
        if output_filename:
            validate_string(output_filename, 1, 255, "output_filename")
            if '..' in output_filename or '/' in output_filename or '\\' in output_filename:
                raise ValueError("output_filename contains invalid characters")
    
    if len(args) > 5:
        language = args[5]
        validate_language_code(language)


def validate_generated_content(content: str, min_length: int = 100, max_length: int = 10000) -> None:
    """Validate generated resume content"""
    if not content:
        raise ValueError("generated content is empty")
    validate_string(content, min_length, max_length, "generated_content")
    # Check for suspicious patterns
    validate_no_sql_injection(content, "generated_content")
    # Note: XSS check might be too strict for generated HTML, so we skip it here


def sanitize_string(value: str) -> str:
    """Sanitize string by removing dangerous characters"""
    if not value:
        return ""
    # Remove null bytes
    value = value.replace('\x00', '')
    # Remove control characters except newlines and tabs
    value = re.sub(r'[\x00-\x08\x0B-\x0C\x0E-\x1F\x7F]', '', value)
    return value.strip()

