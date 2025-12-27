/**
 * Validation utilities for job application worker.
 * Validates inputs, job messages, and scraped data.
 */

const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
const URL_REGEX = /^https?:\/\/[^\s/$.?#].[^\s]*$/;
const SQL_INJECTION_PATTERNS = [
  /\b(SELECT|INSERT|UPDATE|DELETE|DROP|CREATE|ALTER|EXEC|EXECUTE)\b/i,
  /(--|#|\/\*|\*\/)/,
  /\b(UNION|OR|AND)\s+\d+\s*=\s*\d+/i,
  /('|(\\')|(--)|(;)|(\|)|(\*))/,
];
const XSS_PATTERNS = [
  /<script[^>]*>/i,
  /javascript:/i,
  /on\w+\s*=/i,
  /<iframe[^>]*>/i,
  /<object[^>]*>/i,
  /<embed[^>]*>/i,
];

/**
 * Validate UUID format
 */
function validateUUID(value, fieldName = 'id') {
  if (!value) {
    throw new Error(`${fieldName} is required`);
  }
  if (!UUID_REGEX.test(value)) {
    throw new Error(`${fieldName} must be a valid UUID`);
  }
}

/**
 * Validate string length
 */
function validateString(value, minLength, maxLength, fieldName) {
  if (!value) {
    if (minLength > 0) {
      throw new Error(`${fieldName} is required`);
    }
    return;
  }
  if (value.length < minLength) {
    throw new Error(`${fieldName} is too short (minimum ${minLength} characters)`);
  }
  if (value.length > maxLength) {
    throw new Error(`${fieldName} is too long (maximum ${maxLength} characters)`);
  }
}

/**
 * Validate URL format
 */
function validateURL(url, fieldName = 'url') {
  if (!url) {
    throw new Error(`${fieldName} is required`);
  }
  if (!URL_REGEX.test(url)) {
    throw new Error(`${fieldName} must be a valid HTTP/HTTPS URL`);
  }
}

/**
 * Check for potential SQL injection patterns
 */
function validateNoSQLInjection(value, fieldName) {
  if (!value) {
    return;
  }
  for (const pattern of SQL_INJECTION_PATTERNS) {
    if (pattern.test(value)) {
      throw new Error(`${fieldName} contains potentially dangerous content`);
    }
  }
}

/**
 * Check for potential XSS patterns
 */
function validateNoXSS(value, fieldName) {
  if (!value) {
    return;
  }
  for (const pattern of XSS_PATTERNS) {
    if (pattern.test(value)) {
      throw new Error(`${fieldName} contains potentially dangerous content`);
    }
  }
}

/**
 * Validate job application job message
 */
function validateJobApplicationJob(job) {
  // Validate id
  if (!job.id) {
    throw new Error('job id is required');
  }
  validateUUID(job.id, 'id');

  // Validate userId
  if (!job.userId) {
    throw new Error('user ID is required');
  }
  validateUUID(job.userId, 'userId');

  // Validate companyName
  if (!job.companyName) {
    throw new Error('company name is required');
  }
  validateString(job.companyName, 1, 200, 'companyName');
  validateNoSQLInjection(job.companyName, 'companyName');
  validateNoXSS(job.companyName, 'companyName');

  // Validate jobTitle
  if (!job.jobTitle) {
    throw new Error('job title is required');
  }
  validateString(job.jobTitle, 1, 200, 'jobTitle');
  validateNoSQLInjection(job.jobTitle, 'jobTitle');
  validateNoXSS(job.jobTitle, 'jobTitle');

  // Validate jobUrl
  if (!job.jobUrl) {
    throw new Error('job URL is required');
  }
  validateURL(job.jobUrl, 'jobUrl');

  // Validate website
  if (!job.website) {
    throw new Error('website is required');
  }
  validateString(job.website, 1, 255, 'website');
  // Website should not contain spaces
  if (job.website.includes(' ')) {
    throw new Error('website contains invalid characters');
  }

  // Validate location (optional)
  if (job.location) {
    validateString(job.location, 1, 200, 'location');
    validateNoSQLInjection(job.location, 'location');
    validateNoXSS(job.location, 'location');
  }
}

/**
 * Validate scraped job information
 */
function validateScrapedJobInfo(jobInfo) {
  if (!jobInfo) {
    throw new Error('scraped job info is required');
  }

  // Validate title if present
  if (jobInfo.title) {
    validateString(jobInfo.title, 1, 200, 'title');
    validateNoSQLInjection(jobInfo.title, 'title');
    validateNoXSS(jobInfo.title, 'title');
  }

  // Validate company if present
  if (jobInfo.company) {
    validateString(jobInfo.company, 1, 200, 'company');
    validateNoSQLInjection(jobInfo.company, 'company');
    validateNoXSS(jobInfo.company, 'company');
  }

  // Validate description if present
  if (jobInfo.description) {
    validateString(jobInfo.description, 1, 50000, 'description');
    validateNoSQLInjection(jobInfo.description, 'description');
    validateNoXSS(jobInfo.description, 'description');
  }
}

/**
 * Validate cover letter content
 */
function validateCoverLetter(coverLetter) {
  if (!coverLetter) {
    throw new Error('cover letter is required');
  }
  validateString(coverLetter, 100, 10000, 'coverLetter');
  validateNoSQLInjection(coverLetter, 'coverLetter');
  // Note: XSS check might be too strict for generated content, so we skip it
}

/**
 * Sanitize string by removing dangerous characters
 */
function sanitizeString(value) {
  if (!value) {
    return '';
  }
  // Remove null bytes
  value = value.replace(/\x00/g, '');
  // Remove control characters except newlines and tabs
  value = value.replace(/[\x00-\x08\x0B-\x0C\x0E-\x1F\x7F]/g, '');
  return value.trim();
}

export {
  validateUUID,
  validateString,
  validateURL,
  validateNoSQLInjection,
  validateNoXSS,
  validateJobApplicationJob,
  validateScrapedJobInfo,
  validateCoverLetter,
  sanitizeString,
};

