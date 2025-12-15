# How to Test the Job Application Worker

## Quick Test (Recommended)

### Option 1: Run Test Inside Docker Container

```bash
# 1. Get into the container
docker exec -it woragis-job-application-worker bash

# 2. Inside the container, set environment and run test
export DATABASE_URL="postgres://postgres:postgres@database:5432/woragis?sslmode=disable"
export REDIS_URL="redis://redis:6379/0"
export AI_SERVICE_URL="http://ai-service:8000"
export PLAYWRIGHT_HEADLESS="false"  # See browser
export PLAYWRIGHT_SLOW_MO="500"

cd /app
node test-all.js
```

### Option 2: Run Test from Host (Windows PowerShell)

```powershell
# Copy test file into container and run
docker cp job-application-worker/test-all.js woragis-job-application-worker:/app/test-all.js

# Run test
docker exec -it -e DATABASE_URL="postgres://postgres:postgres@database:5432/woragis?sslmode=disable" -e REDIS_URL="redis://redis:6379/0" -e AI_SERVICE_URL="http://ai-service:8000" -e PLAYWRIGHT_HEADLESS="false" woragis-job-application-worker node test-all.js
```

## What the Test Does

1. **Finds Your User ID**
   - Searches database for: `masteringthecode.woragis@gmail.com`
   - Gets your user UUID

2. **Fetches Your Profile**
   - Projects with technologies
   - Posts with skills
   - Technical writings
   - Case studies
   - Certifications
   - Skills
   - Interests

3. **Scrapes LinkedIn Job Page**
   - Navigates to: https://www.linkedin.com/jobs/search/?currentJobId=4347824058...
   - Extracts: Job title, company, location, description
   - Tests AI selector finding
   - Takes screenshots
   - **Does NOT apply**

4. **Generates English Cover Letter**
   - Uses your profile data
   - Personalized to the job
   - Shows in console

5. **Generates Portuguese Cover Letter**
   - Same profile data
   - Written in Portuguese (pt-BR)
   - Shows in console

## Expected Output

```
🧪 Job Application Worker - Comprehensive Test
======================================================================

📧 Step 1: Finding User ID...
✅ Found user: abc-123-def-456

📊 Step 2: Fetching User Profile
  ✅ Projects: 5
  ✅ Posts: 10
  ✅ Technical Writings: 5
  ✅ Skills: 25

🔍 Step 3: Testing LinkedIn Scraping
  ✅ Job Title: [extracted]
  ✅ Company: [extracted]
  ✅ Description: [extracted]

📝 Step 4: Generating Cover Letter
[Your personalized cover letter]

📝 Step 5: Generating Portuguese Cover Letter
[Cover letter em português]
```

## Troubleshooting

### Browser Doesn't Open
- Set `PLAYWRIGHT_HEADLESS=false`
- Make sure you're running in a GUI environment
- For headless servers, check screenshots in `/tmp/`

### "User not found"
- Check database: `SELECT id, email FROM users WHERE email = 'masteringthecode.woragis@gmail.com';`
- Make sure database is accessible from container

### "LinkedIn login required"
- Browser will open (if headless=false)
- Manually log in with your credentials
- Script will continue automatically

### "AI service unavailable"
- Check: `docker ps | grep ai-service`
- Check logs: `docker logs woragis-ai-service`
- Verify `AI_SERVICE_URL` is correct

## Screenshots

Screenshots are saved to `/tmp/` inside the container:
- `linkedin-test.png` - Job page
- `linkedin-error.png` - If error occurs

To view:
```bash
docker exec woragis-job-application-worker ls -la /tmp/linkedin-*.png
docker cp woragis-job-application-worker:/tmp/linkedin-test.png ./
```

