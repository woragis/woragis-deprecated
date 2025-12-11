-- Migration: Add resume metrics columns and indexes
-- Date: 2024-01-XX
-- Description: Adds applications_used, interview_rate, and offer_rate columns to resumes table
--              and adds indexes for performance optimization

-- Add new columns to resumes table
ALTER TABLE resumes
ADD COLUMN IF NOT EXISTS applications_used INTEGER NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS interview_rate DOUBLE PRECISION NOT NULL DEFAULT 0.0,
ADD COLUMN IF NOT EXISTS offer_rate DOUBLE PRECISION NOT NULL DEFAULT 0.0;

-- Add indexes for performance
-- Index on job_applications.resume_id for faster metric calculations
CREATE INDEX IF NOT EXISTS idx_job_applications_resume_id ON job_applications(resume_id);

-- Index on job_application_stages.job_application_id for faster interview queries
CREATE INDEX IF NOT EXISTS idx_job_application_stages_application_id ON job_application_stages(job_application_id);

-- Index on job_application_stages.completed_date for faster completed interview queries
CREATE INDEX IF NOT EXISTS idx_job_application_stages_completed_date ON job_application_stages(completed_date) WHERE completed_date IS NOT NULL;

-- Index on job_application_responses.job_application_id for faster offer queries
CREATE INDEX IF NOT EXISTS idx_job_application_responses_application_id ON job_application_responses(job_application_id);

-- Index on job_application_responses.response_type for faster offer filtering
CREATE INDEX IF NOT EXISTS idx_job_application_responses_type ON job_application_responses(response_type);

-- Update existing resumes to have default values (already set by DEFAULT, but ensure consistency)
UPDATE resumes
SET applications_used = 0,
    interview_rate = 0.0,
    offer_rate = 0.0
WHERE applications_used IS NULL OR interview_rate IS NULL OR offer_rate IS NULL;

