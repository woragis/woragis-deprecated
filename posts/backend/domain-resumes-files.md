# Resumes Domain - File Management

## Overview
How resume files are managed (upload, storage, download, deletion).

## Key Points

### File Storage Structure
- Base file path: Configurable (baseFilePath)
- Upload directory: `uploads/`
- File naming: Safe filename generation (sanitized)
- Full path: `{baseFilePath}/uploads/{safeFilename}`

### File Operations

#### Upload
- User uploads resume file via API
- File validated (size, type)
- Safe filename generated
- File saved to uploads directory
- Database record created with file path

#### Download
- User requests download by resume ID
- File path retrieved from database
- File served from file system
- Content-Type: application/pdf (or detected)
- Content-Disposition: attachment

#### Deletion
- User deletes resume
- Database record deleted
- File deletion from file system (if exists)
- Error handling for missing files

### File Generation
- Generated resumes saved to worker output directory
- Separate from uploaded resumes
- Output path stored in job result
- Accessible via download endpoint

### Security Considerations
- User-scoped access (only owner can download)
- Filename sanitization (prevent path traversal)
- File size limits
- File type validation

## Potential Improvements
- Add file storage abstraction (S3, cloud storage)
- Implement file versioning
- Add file compression
- Support multiple file formats (PDF, DOCX, TXT)
- Add file preview generation (thumbnails)
- Implement file streaming for large files
- Add file caching headers
- Support file upload resume (chunked uploads)
- Add virus scanning
- Implement file encryption at rest
- Add file backup/replication
- Support file sharing (temporary URLs)

