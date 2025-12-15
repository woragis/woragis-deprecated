# Testing Guide

## Quick Start

### 1. Set Environment Variables

Create a `.env` file in the `job-application-worker` directory:

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/woragis?sslmode=disable
REDIS_URL=redis://localhost:6379/0
AI_SERVICE_URL=http://localhost:8000
PLAYWRIGHT_HEADLESS=false  # Set to false to see browser
PLAYWRIGHT_SLOW_MO=500     # Slow down for visibility
TEST_USER_ID=your-user-uuid-here  # Get from database
```

### 2. Get Your User ID

```sql
-- Run in PostgreSQL
SELECT id, email FROM users WHERE email = 'masteringthecode.woragis@gmail.com';
```

Copy the UUID and set it as `TEST_USER_ID` in `.env`.

### 3. Run Tests

```bash
# Test everything (profile fetching + cover letter + scraping)
npm run test

# Test only LinkedIn scraping
npm run test:scraping

# Test only cover letter generation
npm run test:cover-letter
```

## Test Scripts

### 1. Full Test (`test.js`)

Tests:
- ✅ User profile fetching from database
- ✅ AI cover letter generation (English)
- ✅ LinkedIn page scraping (no application)
- ✅ AI selector finding
- ✅ Portuguese cover letter generation

**Usage:**
```bash
npm run test
```

### 2. LinkedIn Scraping Test (`test-linkedin-scraping.js`)

Tests:
- ✅ Navigate to LinkedIn job page
- ✅ Extract job information (title, company, description)
- ✅ Test AI selector finding
- ✅ Take screenshots for debugging
- ❌ Does NOT apply to job

**Usage:**
```bash
npm run test:scraping
```

**Output:**
- Screenshots saved to `/tmp/linkedin-*.png`
- Extracted job information printed to console

### 3. Cover Letter Test (`test-cover-letter.js`)

Tests:
- ✅ Fetch user profile from database
- ✅ Generate English cover letter
- ✅ Generate Portuguese cover letter
- ✅ Verify profile data is used

**Usage:**
```bash
npm run test:cover-letter
```

## Testing LinkedIn Scraping

### Prerequisites

1. **LinkedIn Login**: You may need to log in first
   - The script will detect if you're on a login page
   - You can manually log in if `PLAYWRIGHT_HEADLESS=false`

2. **Job URL**: Use the provided URL or any LinkedIn job URL

### What Gets Tested

1. **Page Navigation**: Can we load the LinkedIn job page?
2. **Element Extraction**: Can we find job title, company, description?
3. **AI Selector Finding**: Does AI find selectors when cached ones fail?
4. **Screenshot Capture**: For visual debugging

### Expected Output

```
🔍 Testing LinkedIn Scraping Capabilities

✅ Page loaded
✅ Job Title: Senior Software Engineer
✅ Company: Google
✅ Location: Remote, United States
✅ Description: We are looking for...

🤖 Testing AI Selector Finding...
✅ Found "Easy Apply" button using AI selectors
```

## Testing Cover Letter Generation

### What Gets Tested

1. **Profile Fetching**: Can we get all user data?
2. **AI Integration**: Does AI service respond?
3. **Personalization**: Does cover letter use profile data?
4. **Language Support**: Can we generate in Portuguese?

### Expected Output

```
📋 Profile Data:
  Projects: 5
  Posts: 10
  Technical Writings: 5
  Skills: 25

📝 Generated Cover Letter:
[Personalized cover letter using your profile data]

✅ Mentions projects from profile
✅ Mentions skills from profile
```

## Testing with Your Credentials

### Option 1: Manual Login (Recommended for Testing)

1. Set `PLAYWRIGHT_HEADLESS=false` in `.env`
2. Run the scraping test
3. When browser opens, manually log in to LinkedIn
4. Script will continue automatically

### Option 2: Automated Login (Future Enhancement)

We can add automated login, but it requires:
- Storing credentials securely
- Handling 2FA
- Managing session cookies

For now, manual login is safer for testing.

## Debugging

### Screenshots

All tests save screenshots to `/tmp/`:
- `linkedin-initial.png` - Initial page load
- `linkedin-final.png` - Final state
- `linkedin-error.png` - If error occurs

### Logs

Check worker logs:
```bash
docker logs woragis-job-application-worker
```

### Common Issues

1. **"Could not find element"**
   - LinkedIn may have changed their HTML
   - AI will try to find new selectors automatically
   - Check screenshots to see page state

2. **"Login required"**
   - LinkedIn detected automation
   - Manually log in with `PLAYWRIGHT_HEADLESS=false`
   - Or implement session cookie storage

3. **"AI service unavailable"**
   - Check if `ai-service` container is running
   - Verify `AI_SERVICE_URL` is correct
   - Check AI service logs

4. **"No profile data"**
   - Verify `TEST_USER_ID` is correct
   - Check if user has projects/posts in database
   - Run: `SELECT COUNT(*) FROM projects WHERE user_id = 'your-uuid';`

## Next Steps After Testing

Once tests pass:

1. **Implement Login**: Add automated LinkedIn login
2. **Refine Selectors**: Update selectors based on test results
3. **Add More Websites**: Implement Glassdoor, WeWorkRemotely
4. **Error Handling**: Improve error recovery
5. **Monitoring**: Add metrics and alerts

## Notes

- ⚠️ **Do NOT apply to jobs during testing** - Tests are read-only
- 🔒 **Keep credentials secure** - Don't commit `.env to git
- 📸 **Screenshots help** - Review them to understand page state
- 🤖 **AI is learning** - Selectors improve over time

