# Quick Test Guide

## Run Tests Inside Docker Container

Since the worker runs in Docker, the easiest way to test is inside the container:

### Step 1: Get into the container

```bash
docker exec -it woragis-job-application-worker bash
```

### Step 2: Set environment variables

```bash
export DATABASE_URL="postgres://postgres:postgres@database:5432/woragis?sslmode=disable"
export REDIS_URL="redis://redis:6379/0"
export AI_SERVICE_URL="http://ai-service:8000"
export PLAYWRIGHT_HEADLESS="false"  # See the browser
export PLAYWRIGHT_SLOW_MO="500"
```

### Step 3: Run the test

```bash
cd /app
node test-all.js
```

The script will:
1. ✅ Find your user ID from email: `masteringthecode.woragis@gmail.com`
2. ✅ Fetch your profile (projects, posts, writings, etc.)
3. ✅ Scrape the LinkedIn job page (extract title, company, description)
4. ✅ Test AI selector finding
5. ✅ Generate English cover letter using your profile
6. ✅ Generate Portuguese cover letter

## What You'll See

### Browser Window
- If `PLAYWRIGHT_HEADLESS=false`, a browser window will open
- You can manually log in to LinkedIn if needed
- The script will extract job information

### Console Output
```
🧪 Job Application Worker - Comprehensive Test

📧 Step 1: Finding User ID...
✅ Found user: abc-123-def-456

📊 Step 2: Fetching User Profile
  ✅ Projects: 5
  ✅ Posts: 10
  ✅ Skills: 25

🔍 Step 3: Testing LinkedIn Scraping
  ✅ Job Title: Senior Software Engineer
  ✅ Company: Google

📝 Step 4: Generating Cover Letter
[Your personalized cover letter here]

📝 Step 5: Generating Portuguese Cover Letter
[Cover letter em português]
```

## Alternative: Run Locally (Outside Docker)

### Prerequisites
- Node.js 20+
- PostgreSQL running
- Redis running
- AI service running

### Setup

```bash
cd job-application-worker
npm install

# Create .env file
cat > .env << EOF
DATABASE_URL=postgres://postgres:postgres@localhost:5432/woragis?sslmode=disable
REDIS_URL=redis://localhost:6379/0
AI_SERVICE_URL=http://localhost:8000
PLAYWRIGHT_HEADLESS=false
PLAYWRIGHT_SLOW_MO=500
EOF

# Run test
npm run test
```

## Expected Results

### ✅ Profile Fetching
- Should find your user by email
- Should fetch projects, posts, writings, etc.
- Should show counts for each data type

### ✅ LinkedIn Scraping
- Should navigate to job page
- Should extract job title, company, description
- May require login (browser will show login page)
- Screenshots saved for debugging

### ✅ Cover Letter Generation
- Should use your actual profile data
- Should mention your projects/skills
- Should be personalized to the job
- Portuguese version should be in Portuguese

## Troubleshooting

### "User not found"
- Check if email exists in database
- Run: `SELECT id, email FROM users WHERE email = 'masteringthecode.woragis@gmail.com';`

### "LinkedIn login required"
- Set `PLAYWRIGHT_HEADLESS=false`
- Manually log in when browser opens
- Script will continue automatically

### "AI service unavailable"
- Check if `ai-service` container is running: `docker ps | grep ai-service`
- Check AI service logs: `docker logs woragis-ai-service`

### "No profile data"
- Check if you have projects/posts in database
- Run: `SELECT COUNT(*) FROM projects WHERE user_id = 'your-uuid';`

## Notes

- ⚠️ **No application will be made** - Tests are read-only
- 🔒 **Credentials are not stored** - Login is manual for testing
- 📸 **Screenshots help debug** - Check `/tmp/` directory
- 🤖 **AI learns** - Selectors improve over time

