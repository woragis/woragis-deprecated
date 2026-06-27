# Translations Domain - Translation Queue & Processing

## Overview
How translation jobs are queued and processed asynchronously.

## Key Points

### Queue Integration
- Translation service integrates with Redis queue
- Jobs created on translation request
- Worker processes jobs asynchronously
- Queue pattern: LPush/BRPop

### Translation Job Flow
1. User/System requests translation
2. Check if translation already exists (completed)
3. Create TranslationJob (entity_type, entity_id, language, fields, source_text)
4. Enqueue job to Redis
5. Create/update Translation record (status: pending)
6. Worker dequeues and processes
7. AI service called for translation
8. Translation saved to database
9. Status updated to completed

### Job Processing
- ProcessTranslationJob handles worker processing
- Fetches entity from database if source_text not provided
- Calls AI service for each field translation
- Saves translations to database
- Handles errors and retries

### Status Management
- Pending: Job queued, not processed
- Processing: Worker is processing
- Completed: Translation successful
- Failed: Translation failed

### Entity Support
- Projects, Skills, Interests, Posts, Case Studies, System Designs
- Field-level translations (title, description, content)
- Source text provided or fetched from entity

## Potential Improvements
- Add translation priority levels
- Implement translation batching
- Add translation caching (avoid re-translating)
- Support translation versioning
- Add translation quality scoring
- Implement translation retry with backoff
- Add translation deduplication
- Support bulk translation requests
- Add translation progress tracking
- Implement translation cancellation
- Add translation result validation
- Support translation preview

