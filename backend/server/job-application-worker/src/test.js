import 'dotenv/config';
import { Database } from './database.js';
import { CoverLetterService } from './coverLetter.js';
import { Scraper } from './scraper.js';
import { logger } from './utils/logger.js';

/**
 * Test script for job application worker
 * Tests: Profile fetching, AI cover letter generation, LinkedIn scraping
 */
async function test() {
  console.log('🧪 Starting Job Application Worker Tests\n');

  // Initialize services
  const db = new Database();
  const coverLetterService = new CoverLetterService();
  const scraper = new Scraper();

  try {
    // Connect to database
    console.log('📊 Connecting to database...');
    await db.connect();
    console.log('✅ Database connected\n');

    // Test 1: Fetch User Profile
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('TEST 1: Fetch User Profile');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');
    
    // You'll need to provide a valid user ID from your database
    const testUserId = process.env.TEST_USER_ID || '00000000-0000-0000-0000-000000000000';
    console.log(`Fetching profile for user: ${testUserId}\n`);
    
    const profile = await db.fetchUserProfile(testUserId);
    
    console.log('📋 Profile Summary:');
    console.log(`  - Projects: ${profile.projects.length}`);
    console.log(`  - Posts: ${profile.posts.length}`);
    console.log(`  - Technical Writings: ${profile.technicalWritings.length}`);
    console.log(`  - Case Studies: ${profile.caseStudies.length}`);
    console.log(`  - Certifications: ${profile.certifications.length}`);
    console.log(`  - Skills: ${profile.skills.length}`);
    console.log(`  - Interests: ${profile.interests.length}\n`);

    if (profile.projects.length > 0) {
      console.log('📁 Sample Project:');
      const sampleProject = profile.projects[0];
      console.log(`  Name: ${sampleProject.name}`);
      console.log(`  Description: ${sampleProject.description?.substring(0, 100)}...`);
      console.log(`  Technologies: ${sampleProject.techStack?.map(t => t.name).join(', ') || 'None'}\n`);
    }

    // Test 2: Generate Cover Letter
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('TEST 2: Generate Cover Letter with AI');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

    const jobInfo = {
      companyName: 'Google',
      jobTitle: 'Senior Software Engineer',
      jobDescription: 'We are looking for a Senior Software Engineer with experience in distributed systems, microservices architecture, and cloud technologies. The ideal candidate should have strong backend development skills, experience with Go or Python, and knowledge of Kubernetes and Docker.',
      location: 'Remote, United States',
      requirements: ['Go', 'Python', 'Kubernetes', 'Docker', 'Microservices', 'Distributed Systems'],
    };

    console.log('Job Information:');
    console.log(`  Company: ${jobInfo.companyName}`);
    console.log(`  Position: ${jobInfo.jobTitle}`);
    console.log(`  Location: ${jobInfo.location}\n`);

    console.log('🤖 Generating cover letter with AI...\n');
    const coverLetter = await coverLetterService.generateCoverLetter(profile, jobInfo);
    
    console.log('📝 Generated Cover Letter:');
    console.log('─'.repeat(70));
    console.log(coverLetter);
    console.log('─'.repeat(70));
    console.log(`\nLength: ${coverLetter.length} characters\n`);

    // Test 3: Test LinkedIn Scraping (without applying)
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('TEST 3: LinkedIn Scraping Test (No Application)');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

    const testJob = {
      id: 'test-job-1',
      userId: testUserId,
      companyName: 'Test Company',
      location: 'Remote',
      jobTitle: 'Software Engineer',
      jobUrl: 'https://www.linkedin.com/jobs/search/?currentJobId=4347824058&f_WT=2&origin=JOB_COLLECTION_PAGE_LOCATION_SUGGESTION&refresh=true',
      website: 'linkedin',
    };

    console.log('🔍 Testing LinkedIn page scraping...');
    console.log(`Job URL: ${testJob.jobUrl}\n`);

    await scraper.initialize();

    // Launch browser and navigate (but don't apply)
    const { chromium } = await import('playwright');
    const headless = process.env.PLAYWRIGHT_HEADLESS !== 'false';
    const browser = await chromium.launch({ headless, slowMo: 500 });
    const context = await browser.newContext({
      viewport: { width: 1920, height: 1080 },
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    });
    const page = await context.newPage();

    try {
      console.log('🌐 Navigating to LinkedIn job page...');
      await page.goto(testJob.jobUrl, { waitUntil: 'networkidle', timeout: 30000 });
      console.log('✅ Page loaded\n');

      // Take screenshot for debugging
      await page.screenshot({ path: '/tmp/linkedin-job-page.png', fullPage: false });
      console.log('📸 Screenshot saved to /tmp/linkedin-job-page.png\n');

      // Try to find job details
      console.log('🔍 Extracting job information...\n');
      
      // Try to find job title
      try {
        const jobTitle = await page.locator('h1.job-details-jobs-unified-top-card__job-title, h1[data-test-id="job-title"]').first();
        if (await jobTitle.count() > 0) {
          const title = await jobTitle.textContent();
          console.log(`  Job Title: ${title?.trim()}`);
        }
      } catch (error) {
        console.log('  ⚠️  Could not find job title');
      }

      // Try to find company name
      try {
        const company = await page.locator('a.job-details-jobs-unified-top-card__company-name, a[data-test-id="company-name"]').first();
        if (await company.count() > 0) {
          const companyName = await company.textContent();
          console.log(`  Company: ${companyName?.trim()}`);
        }
      } catch (error) {
        console.log('  ⚠️  Could not find company name');
      }

      // Try to find job description
      try {
        const description = await page.locator('.jobs-description-content__text, [data-test-id="job-description"]').first();
        if (await description.count() > 0) {
          const descText = await description.textContent();
          console.log(`  Description Preview: ${descText?.substring(0, 200).trim()}...\n`);
        }
      } catch (error) {
        console.log('  ⚠️  Could not find job description\n');
      }

      // Test selector finding with AI (if needed)
      console.log('🤖 Testing AI Selector Finding...\n');
      console.log('  Looking for "Easy Apply" button...');
      
      // This will test the self-healing scraper
      try {
        await scraper.selfHealing.findElement(page, 'linkedin', 'easy-apply-button', 'Easy Apply button');
        console.log('  ✅ Found element using AI selector finder\n');
      } catch (error) {
        console.log(`  ⚠️  Element not found: ${error.message}\n`);
        console.log('  This is expected if you need to log in first\n');
      }

      // Get page HTML for analysis
      const html = await page.content();
      console.log(`📄 Page HTML length: ${html.length} characters\n`);

      console.log('✅ LinkedIn scraping test completed (no application made)\n');

    } catch (error) {
      console.error('❌ LinkedIn scraping error:', error.message);
      console.error('Stack:', error.stack);
    } finally {
      await browser.close();
    }

    // Test 4: Test with different language
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('TEST 4: Generate Cover Letter in Portuguese');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

    // Modify the prompt to request Portuguese
    const portugueseJobInfo = {
      ...jobInfo,
      language: 'pt-BR', // Add language hint
    };

    // We'll modify the prompt to request Portuguese
    const portuguesePrompt = coverLetterService.buildPrompt(profile, portugueseJobInfo);
    const portuguesePromptWithLanguage = portuguesePrompt.replace(
      'Write the cover letter now:',
      'IMPORTANT: Write the cover letter in Portuguese (Brazil). Write the cover letter now:'
    );

    console.log('🤖 Generating Portuguese cover letter...\n');
    
    try {
      const axios = (await import('axios')).default;
      const response = await axios.post(
        `${process.env.AI_SERVICE_URL || 'http://ai-service:8000'}/api/chat/completions`,
        {
          provider: 'openai',
          model: 'gpt-4o-mini',
          temperature: 0.7,
          messages: [
            {
              role: 'user',
              content: portuguesePromptWithLanguage,
            },
          ],
          max_tokens: 1500,
        },
        { timeout: 30000 }
      );

      const portugueseCoverLetter = response.data.message?.content || response.data.choices?.[0]?.message?.content;
      
      console.log('📝 Cover Letter em Português:');
      console.log('─'.repeat(70));
      console.log(portugueseCoverLetter.trim());
      console.log('─'.repeat(70));
      console.log(`\nLength: ${portugueseCoverLetter.trim().length} characters\n`);
    } catch (error) {
      console.error('❌ Failed to generate Portuguese cover letter:', error.message);
    }

    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('✅ All tests completed!');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

  } catch (error) {
    console.error('❌ Test failed:', error);
    console.error('Stack:', error.stack);
    process.exit(1);
  } finally {
    await db.disconnect();
    await scraper.cleanup();
  }
}

// Run tests
test().catch(console.error);

