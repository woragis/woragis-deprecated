import 'dotenv/config';
import { Database } from './src/database.js';
import { CoverLetterService } from './src/coverLetter.js';
import { chromium } from 'playwright';
import { SelfHealingScraper } from './src/selfHealingScraper.js';
import axios from 'axios';

const JOB_URL = 'https://www.linkedin.com/jobs/search/?currentJobId=4347824058&f_WT=2&origin=JOB_COLLECTION_PAGE_LOCATION_SUGGESTION&refresh=true';

async function main() {
  console.log('🧪 Job Application Worker - Comprehensive Test\n');
  console.log('='.repeat(70));

  // Step 1: Get User ID from email
  console.log('\n📧 Step 1: Finding User ID...\n');
  const db = new Database();
  await db.connect();

  const email = 'masteringthecode.woragis@gmail.com';
  const userResult = await db.pool.query(
    'SELECT id FROM users WHERE email = $1 LIMIT 1',
    [email]
  );

  if (userResult.rows.length === 0) {
    console.error(`❌ User not found with email: ${email}`);
    console.log('\nPlease check your database or use a different email.\n');
    process.exit(1);
  }

  const userId = userResult.rows[0].id;
  console.log(`✅ Found user: ${userId}\n`);

  // Step 2: Fetch Profile
  console.log('='.repeat(70));
  console.log('📊 Step 2: Fetching User Profile from Database\n');

  const profile = await db.fetchUserProfile(userId);
  
  console.log('Profile Summary:');
  console.log(`  ✅ Projects: ${profile.projects.length}`);
  console.log(`  ✅ Posts: ${profile.posts.length}`);
  console.log(`  ✅ Technical Writings: ${profile.technicalWritings.length}`);
  console.log(`  ✅ Case Studies: ${profile.caseStudies.length}`);
  console.log(`  ✅ Project Case Studies: ${profile.projectCaseStudies.length}`);
  console.log(`  ✅ Certifications: ${profile.certifications.length}`);
  console.log(`  ✅ Skills: ${profile.skills.length}`);
  console.log(`  ✅ Interests: ${profile.interests.length}\n`);

  if (profile.projects.length > 0) {
    console.log('Sample Project:');
    const p = profile.projects[0];
    console.log(`  Name: ${p.name}`);
    console.log(`  Description: ${p.description?.substring(0, 100) || 'N/A'}...`);
    console.log(`  Tech Stack: ${p.techStack?.map(t => t.name).join(', ') || 'None'}\n`);
  }

  // Step 3: Test LinkedIn Scraping
  console.log('='.repeat(70));
  console.log('🔍 Step 3: Testing LinkedIn Job Page Scraping\n');
  console.log(`Job URL: ${JOB_URL}\n`);

  const headless = process.env.PLAYWRIGHT_HEADLESS !== 'false';
  const browser = await chromium.launch({ 
    headless, 
    slowMo: parseInt(process.env.PLAYWRIGHT_SLOW_MO || '500') 
  });

  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
  });

  const page = await context.newPage();
  const selfHealing = new SelfHealingScraper();
  await selfHealing.initialize();

  let extractedJobInfo = {};

  try {
    console.log('🌐 Navigating to LinkedIn...');
    await page.goto(JOB_URL, { waitUntil: 'networkidle', timeout: 30000 });
    await page.waitForTimeout(3000);
    console.log('✅ Page loaded\n');

    // Extract job information
    console.log('📋 Extracting job information...\n');

    // Try to find job title
    const titleSelectors = [
      'h1.job-details-jobs-unified-top-card__job-title',
      'h1[data-test-id="job-title"]',
      'h1.jobs-details-top-card__job-title',
      'h1',
    ];

    for (const selector of titleSelectors) {
      try {
        const el = await page.locator(selector).first();
        if (await el.count() > 0) {
          extractedJobInfo.title = (await el.textContent())?.trim();
          console.log(`  ✅ Job Title: ${extractedJobInfo.title}`);
          break;
        }
      } catch (e) {
        continue;
      }
    }

    // Try to find company
    const companySelectors = [
      'a.job-details-jobs-unified-top-card__company-name',
      'a[data-test-id="company-name"]',
      '.jobs-company__box a',
    ];

    for (const selector of companySelectors) {
      try {
        const el = await page.locator(selector).first();
        if (await el.count() > 0) {
          extractedJobInfo.company = (await el.textContent())?.trim();
          console.log(`  ✅ Company: ${extractedJobInfo.company}`);
          break;
        }
      } catch (e) {
        continue;
      }
    }

    // Try to find location
    try {
      const loc = await page.locator('.job-details-jobs-unified-top-card__primary-description-without-tagline, .jobs-details-top-card__bullet').first();
      if (await loc.count() > 0) {
        extractedJobInfo.location = (await loc.textContent())?.trim();
        console.log(`  ✅ Location: ${extractedJobInfo.location}`);
      }
    } catch (e) {
      // Ignore
    }

    // Try to find description
    try {
      const desc = await page.locator('.jobs-description-content__text, [data-test-id="job-description"]').first();
      if (await desc.count() > 0) {
        const descText = await desc.textContent();
        extractedJobInfo.description = descText?.trim() || '';
        console.log(`  ✅ Description: ${extractedJobInfo.description.substring(0, 200)}...`);
        console.log(`     Full length: ${extractedJobInfo.description.length} characters\n`);
      }
    } catch (e) {
      console.log('  ⚠️  Could not extract description\n');
    }

    // Test AI selector finding
    console.log('🤖 Testing AI Selector Finding...\n');
    try {
      const element = await selfHealing.findElement(
        page,
        'linkedin',
        'easy-apply-button',
        'Easy Apply button or Apply button'
      );
      if (element) {
        const text = await element.textContent();
        console.log(`  ✅ Found "Easy Apply" button: "${text?.trim()}"\n`);
      }
    } catch (error) {
      console.log(`  ⚠️  Could not find Easy Apply button: ${error.message}`);
      console.log('     (This is expected if you need to log in first)\n');
    }

    // Take screenshot
    await page.screenshot({ path: '/tmp/linkedin-test.png', fullPage: false });
    console.log('📸 Screenshot saved: /tmp/linkedin-test.png\n');

  } catch (error) {
    console.error('❌ Scraping error:', error.message);
    await page.screenshot({ path: '/tmp/linkedin-error.png', fullPage: true });
    console.log('📸 Error screenshot saved: /tmp/linkedin-error.png\n');
  } finally {
    await browser.close();
    await selfHealing.cleanup();
  }

  // Step 4: Generate Cover Letter (English)
  console.log('='.repeat(70));
  console.log('📝 Step 4: Generating Cover Letter in English\n');

  const coverLetterService = new CoverLetterService();
  
  const jobInfo = {
    companyName: extractedJobInfo.company || 'Company',
    jobTitle: extractedJobInfo.title || 'Software Engineer',
    jobDescription: extractedJobInfo.description || 'Software engineering position',
    location: extractedJobInfo.location || 'Remote',
    requirements: [],
  };

  console.log('Job Information:');
  console.log(`  Company: ${jobInfo.companyName}`);
  console.log(`  Position: ${jobInfo.jobTitle}`);
  console.log(`  Location: ${jobInfo.location}\n`);

  try {
    const coverLetter = await coverLetterService.generateCoverLetter(profile, jobInfo);
    
    console.log('Generated Cover Letter:');
    console.log('─'.repeat(70));
    console.log(coverLetter);
    console.log('─'.repeat(70));
    console.log(`\n📊 Length: ${coverLetter.length} characters\n`);

    // Check if it uses profile data
    const usesProfile = profile.projects.some(p => 
      coverLetter.toLowerCase().includes(p.name.toLowerCase())
    ) || profile.skills.some(s => 
      coverLetter.toLowerCase().includes(s.name.toLowerCase())
    );

    if (usesProfile) {
      console.log('✅ Cover letter uses profile data!\n');
    } else {
      console.log('⚠️  Cover letter may not be using profile data effectively\n');
    }

  } catch (error) {
    console.error('❌ Failed to generate cover letter:', error.message);
  }

  // Step 5: Generate Cover Letter (Portuguese)
  console.log('='.repeat(70));
  console.log('📝 Step 5: Generating Cover Letter in Portuguese\n');

  try {
    const portuguesePrompt = coverLetterService.buildPrompt(profile, jobInfo);
    const portuguesePromptWithLanguage = portuguesePrompt.replace(
      'Write the cover letter now:',
      'IMPORTANT: Write the cover letter in Portuguese (Brazilian Portuguese - pt-BR). Use professional Brazilian business language. Write the cover letter now:'
    );

    const response = await axios.post(
      `${process.env.AI_SERVICE_URL || 'http://localhost:8000'}/api/chat/completions`,
      {
        provider: 'openai',
        model: 'gpt-4o-mini',
        temperature: 0.7,
        messages: [{ role: 'user', content: portuguesePromptWithLanguage }],
        max_tokens: 1500,
      },
      { timeout: 30000 }
    );

    const portugueseCoverLetter = response.data.message?.content || response.data.choices?.[0]?.message?.content;

    console.log('Cover Letter em Português:');
    console.log('─'.repeat(70));
    console.log(portugueseCoverLetter.trim());
    console.log('─'.repeat(70));
    console.log(`\n📊 Length: ${portugueseCoverLetter.trim().length} characters\n`);

    // Check if it's in Portuguese
    const isPortuguese = /(experiência|oportunidade|empresa|desenvolvimento|tecnologia|carreira)/i.test(portugueseCoverLetter);
    if (isPortuguese) {
      console.log('✅ Cover letter is in Portuguese!\n');
    }

  } catch (error) {
    console.error('❌ Failed to generate Portuguese cover letter:', error.message);
  }

  console.log('='.repeat(70));
  console.log('✅ All Tests Completed!');
  console.log('='.repeat(70));

  await db.disconnect();
}

main().catch(console.error);

