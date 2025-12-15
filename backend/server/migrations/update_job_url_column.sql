-- Migration: Update job_url column to support longer URLs
-- Date: 2025-12-13
-- Description: Changes job_url column from varchar(512) to text to accommodate longer URLs
--              (e.g., LinkedIn URLs with tracking parameters)

-- Alter the job_url column to text type
ALTER TABLE job_applications
ALTER COLUMN job_url TYPE TEXT;
