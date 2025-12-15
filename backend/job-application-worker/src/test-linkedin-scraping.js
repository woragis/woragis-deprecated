import 'dotenv/config';
import { chromium } from 'playwright';
import { SelfHealingScraper } from './selfHealingScraper.js';
import { logger } from './utils/logger.js';

/**
 * Test LinkedIn scraping capabilities
 * This script will:
 * 1. Navigate to LinkedIn job page
 * 2. Try to extract job information
 * 3. Test AI selector finding
 * 4. Take screenshots for debugging
 * 
 * IMPORTANT: This does NOT apply to the job
 */
async function testLinkedInScraping() {
  console.log('🔍 Testing LinkedIn Scraping Capabilities\n');

  const jobUrl = 'https://www.linkedin.com/jobs/search/?currentJobId=4347824058&f_WT=2&origin=JOB_COLLECTION_PAGE_LOCATION_SUGGESTION&refresh=true';
  
  const headless = process.env.PLAYWRIGHT_HEADLESS !== 'false';
  const slowMo = parseInt(process.env.PLAYWRIGHT_SLOW_MO || '500');
  
  console.log('Configuration:');
  console.log(`  Headless: ${headless}`);
  console.log(`  Slow Motion: ${slowMo}ms`);
  console.log(`  Job URL: ${jobUrl}\n`);

  const browser = await chromium.launch({
    headless,
    slowMo,
  });

  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
  });

  const page = await context.newPage();
  const selfHealing = new SelfHealingScraper();
  await selfHealing.initialize();

  try {
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('STEP 1: Navigate to Job Page');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

    console.log('🌐 Navigating to LinkedIn...');
    await page.goto(jobUrl, { waitUntil: 'networkidle', timeout: 30000 });
    await page.waitForTimeout(3000); // Wait for page to fully load
    console.log('✅ Page loaded\n');

    // Take initial screenshot
    await page.screenshot({ path: '/tmp/linkedin-initial.png', fullPage: false });
    console.log('📸 Screenshot saved: /tmp/linkedin-initial.png\n');

    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('STEP 2: Extract Job Information');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

    const jobInfo = {};

    // Try multiple selectors for job title
    const titleSelectors = [
      'h1.job-details-jobs-unified-top-card__job-title',
      'h1[data-test-id="job-title"]',
      'h1.jobs-details-top-card__job-title',
      'h1',
    ];

    for (const selector of titleSelectors) {
      try {
        const element = await page.locator(selector).first();
        if (await element.count() > 0) {
          jobInfo.title = (await element.textContent())?.trim();
          console.log(`✅ Job Title: ${jobInfo.title}`);
          break;
        }
      } catch (error) {
        continue;
      }
    }

    // Try multiple selectors for company
    const companySelectors = [
      'a.job-details-jobs-unified-top-card__company-name',
      'a[data-test-id="company-name"]',
      'a.jobs-details-top-card__company-name',
      '.jobs-company__box a',
    ];

    for (const selector of companySelectors) {
      try {
        const element = await page.locator(selector).first();
        if (await element.count() > 0) {
          jobInfo.company = (await element.textContent())?.trim();
          console.log(`✅ Company: ${jobInfo.company}`);
          break;
        }
      } catch (error) {
        continue;
      }
    }

    // Try to find location
    const locationSelectors = [
      '.job-details-jobs-unified-top-card__primary-description-without-tagline',
      '.jobs-details-top-card__bullet',
      '[data-test-id="job-location"]',
    ];

    for (const selector of locationSelectors) {
      try {
        const element = await page.locator(selector).first();
        if (await element.count() > 0) {
          jobInfo.location = (await element.textContent())?.trim();
          console.log(`✅ Location: ${jobInfo.location}`);
          break;
        }
      } catch (error) {
        continue;
      }
    }

    // Try to extract job description
    const descSelectors = [
      '.jobs-description-content__text',
      '[data-test-id="job-description"]',
      '.jobs-description__text',
      '#job-details',
    ];

    for (const selector of descSelectors) {
      try {
        const element = await page.locator(selector).first();
        if (await element.count() > 0) {
          const desc = await element.textContent();
          jobInfo.description = desc?.trim() || '';
          console.log(`✅ Description: ${jobInfo.description.substring(0, 200)}...`);
          console.log(`   Full length: ${jobInfo.description.length} characters\n`);
          break;
        }
      } catch (error) {
        continue;
      }
    }

    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('STEP 3: Test AI Selector Finding');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

    // Test finding "Easy Apply" button using AI
    console.log('🤖 Testing AI-powered selector finding for "Easy Apply" button...\n');

    try {
      // First, try with cached selectors
      const element = await selfHealing.findElement(
        page,
        'linkedin',
        'easy-apply-button',
        'Easy Apply button or Apply button'
      );

      if (element) {
        console.log('✅ Found "Easy Apply" button using AI selectors\n');
        
        // Get button text and position
        const buttonText = await element.textContent();
        const boundingBox = await element.boundingBox();
        console.log(`  Button Text: ${buttonText?.trim()}`);
        console.log(`  Position: x=${boundingBox?.x}, y=${boundingBox?.y}\n`);
      }
    } catch (error) {
      console.log(`⚠️  Could not find "Easy Apply" button: ${error.message}`);
      console.log('  This might be because:');
      console.log('    1. You need to log in first');
      console.log('    2. The job doesn\'t have Easy Apply');
      console.log('    3. The page structure is different\n');
    }

    // Test finding other common elements
    const testElements = [
      { action: 'job-title', description: 'Job title heading' },
      { action: 'company-name', description: 'Company name link' },
      { action: 'job-description', description: 'Job description section' },
    ];

    for (const test of testElements) {
      try {
        const element = await selfHealing.findElement(
          page,
          'linkedin',
          test.action,
          test.description
        );
        if (element) {
          console.log(`✅ Found ${test.description}\n`);
        }
      } catch (error) {
        console.log(`⚠️  Could not find ${test.description}: ${error.message}\n`);
      }
    }

    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('STEP 4: Page Analysis');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

    // Get page title
    const pageTitle = await page.title();
    console.log(`Page Title: ${pageTitle}\n`);

    // Get URL (might redirect)
    const currentUrl = page.url();
    console.log(`Current URL: ${currentUrl}\n`);

    // Check if we're on a login page
    const isLoginPage = currentUrl.includes('/login') || 
                        currentUrl.includes('/checkpoint') ||
                        await page.locator('input[type="password"]').count() > 0;

    if (isLoginPage) {
      console.log('⚠️  Detected login page - authentication required\n');
      console.log('To test full scraping, you would need to:');
      console.log('  1. Log in with credentials');
      console.log('  2. Handle 2FA if required');
      console.log('  3. Save session cookies\n');
    }

    // Take final screenshot
    await page.screenshot({ path: '/tmp/linkedin-final.png', fullPage: true });
    console.log('📸 Full page screenshot saved: /tmp/linkedin-final.png\n');

    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('✅ LinkedIn Scraping Test Completed');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

    console.log('📊 Extracted Job Information:');
    console.log(JSON.stringify(jobInfo, null, 2));

  } catch (error) {
    console.error('❌ Error during scraping test:', error.message);
    console.error('Stack:', error.stack);
    
    // Take error screenshot
    await page.screenshot({ path: '/tmp/linkedin-error.png', fullPage: true });
    console.log('\n📸 Error screenshot saved: /tmp/linkedin-error.png');
  } finally {
    await browser.close();
    await selfHealing.cleanup();
  }
}

// Run test
testLinkedInScraping().catch(console.error);

