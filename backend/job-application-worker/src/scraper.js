import { chromium } from 'playwright';
import { logger } from './utils/logger.js';
import { SelfHealingScraper } from './selfHealingScraper.js';

export class Scraper {
  constructor() {
    this.browser = null;
    this.selfHealing = new SelfHealingScraper();
  }

  async initialize() {
    await this.selfHealing.initialize();
  }

  /**
   * Test method: Extract job information without applying
   * Useful for testing scraping capabilities
   */
  async testJobExtraction(jobUrl) {
    logger.info('Testing job extraction', { jobUrl });

    const headless = process.env.PLAYWRIGHT_HEADLESS !== 'false';
    const slowMo = parseInt(process.env.PLAYWRIGHT_SLOW_MO || '100');
    
    this.browser = await chromium.launch({
      headless,
      slowMo,
    });

    const context = await this.browser.newContext({
      viewport: { width: 1920, height: 1080 },
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
    });

    const page = await context.newPage();

    try {
      await page.goto(jobUrl, { waitUntil: 'networkidle', timeout: 30000 });
      await page.waitForTimeout(2000);

      const jobInfo = {};

      // Extract job title
      const titleSelectors = [
        'h1.job-details-jobs-unified-top-card__job-title',
        'h1[data-test-id="job-title"]',
        'h1.jobs-details-top-card__job-title',
      ];

      for (const selector of titleSelectors) {
        try {
          const element = await page.locator(selector).first();
          if (await element.count() > 0) {
            jobInfo.title = (await element.textContent())?.trim();
            break;
          }
        } catch (error) {
          continue;
        }
      }

      // Extract company
      const companySelectors = [
        'a.job-details-jobs-unified-top-card__company-name',
        'a[data-test-id="company-name"]',
      ];

      for (const selector of companySelectors) {
        try {
          const element = await page.locator(selector).first();
          if (await element.count() > 0) {
            jobInfo.company = (await element.textContent())?.trim();
            break;
          }
        } catch (error) {
          continue;
        }
      }

      // Extract description
      try {
        const descElement = await page.locator('.jobs-description-content__text').first();
        if (await descElement.count() > 0) {
          jobInfo.description = (await descElement.textContent())?.trim();
        }
      } catch (error) {
        // Ignore
      }

      // Take screenshot
      await page.screenshot({ path: '/tmp/linkedin-job-test.png', fullPage: false });

      await context.close();
      await this.browser.close();

      return jobInfo;
    } catch (error) {
      await context.close();
      await this.browser.close();
      throw error;
    }
  }

  async applyToJob(job, coverLetter) {
    logger.info('Applying to job', {
      company: job.companyName,
      website: job.website,
      jobUrl: job.jobUrl,
    });

    // Launch browser
    const headless = process.env.PLAYWRIGHT_HEADLESS !== 'false';
    const slowMo = parseInt(process.env.PLAYWRIGHT_SLOW_MO || '100');
    
    this.browser = await chromium.launch({
      headless,
      slowMo,
    });

    const context = await this.browser.newContext({
      viewport: { width: 1920, height: 1080 },
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
    });

    const page = await context.newPage();

    try {
      // Navigate to job URL
      await page.goto(job.jobUrl, { waitUntil: 'networkidle' });

      // Apply based on website type
      switch (job.website.toLowerCase()) {
        case 'linkedin':
          await this.applyLinkedIn(page, job, coverLetter);
          break;
        case 'glassdoor':
          await this.applyGlassdoor(page, job, coverLetter);
          break;
        case 'weworkremotely':
          await this.applyWeWorkRemotely(page, job, coverLetter);
          break;
        default:
          await this.applyGeneric(page, job, coverLetter);
      }

      logger.info('Successfully applied to job', {
        company: job.companyName,
        website: job.website,
      });
    } catch (error) {
      logger.error('Failed to apply to job', {
        error: error.message,
        company: job.companyName,
        website: job.website,
        stack: error.stack,
      });
      throw error;
    } finally {
      await context.close();
      await this.browser.close();
    }
  }

  async applyLinkedIn(page, job, coverLetter) {
    const website = 'linkedin';
    
    try {
      // Step 1: Find and click "Easy Apply" button
      await this.selfHealing.clickElement(
        page,
        website,
        'easy-apply-button',
        'Easy Apply button or Apply button'
      );

      // Wait for application form to appear
      await page.waitForTimeout(2000);

      // Step 2: Fill phone number (if required)
      try {
        await this.selfHealing.fillField(
          page,
          website,
          'phone-field',
          'Phone number input field',
          '1234567890' // TODO: Get from user profile
        );
      } catch (error) {
        logger.warn('Phone field not found or not required', { error: error.message });
      }

      // Step 3: Fill cover letter
      try {
        await this.selfHealing.fillField(
          page,
          website,
          'cover-letter-field',
          'Cover letter or message textarea',
          coverLetter
        );
      } catch (error) {
        logger.warn('Cover letter field not found', { error: error.message });
      }

      // Step 4: Submit application
      await this.selfHealing.clickElement(
        page,
        website,
        'submit-button',
        'Submit application button'
      );

      // Wait for confirmation
      await page.waitForTimeout(2000);

      logger.info('LinkedIn application completed');
    } catch (error) {
      logger.error('LinkedIn application failed', { error: error.message });
      throw error;
    }
  }

  async applyGlassdoor(page, job, coverLetter) {
    const website = 'glassdoor';
    
    try {
      // Step 1: Find and click "Apply Now" button
      await this.selfHealing.clickElement(
        page,
        website,
        'apply-button',
        'Apply Now or Apply button'
      );

      await page.waitForTimeout(2000);

      // Step 2: Fill application form fields
      // Glassdoor often redirects to external sites, so this may vary
      try {
        await this.selfHealing.fillField(
          page,
          website,
          'cover-letter-field',
          'Cover letter textarea',
          coverLetter
        );
      } catch (error) {
        logger.warn('Cover letter field not found', { error: error.message });
      }

      // Step 3: Submit
      await this.selfHealing.clickElement(
        page,
        website,
        'submit-button',
        'Submit application button'
      );

      await page.waitForTimeout(2000);
      logger.info('Glassdoor application completed');
    } catch (error) {
      logger.error('Glassdoor application failed', { error: error.message });
      throw error;
    }
  }

  async applyWeWorkRemotely(page, job, coverLetter) {
    const website = 'weworkremotely';
    
    try {
      // WeWorkRemotely usually has a simpler form
      await this.selfHealing.fillField(
        page,
        website,
        'email-field',
        'Email input field',
        'user@example.com' // TODO: Get from user profile
      );

      await this.selfHealing.fillField(
        page,
        website,
        'cover-letter-field',
        'Cover letter or message textarea',
        coverLetter
      );

      await this.selfHealing.clickElement(
        page,
        website,
        'submit-button',
        'Submit application button'
      );

      await page.waitForTimeout(2000);
      logger.info('WeWorkRemotely application completed');
    } catch (error) {
      logger.error('WeWorkRemotely application failed', { error: error.message });
      throw error;
    }
  }

  async applyGeneric(page, job, coverLetter) {
    // Generic application: try to find common form elements
    try {
      // Look for common form fields
      const emailField = await page.locator('input[type="email"], input[name*="email" i]').first();
      if (await emailField.count() > 0) {
        await emailField.fill('user@example.com'); // TODO: Get from user profile
      }

      const messageField = await page.locator('textarea[name*="message" i], textarea[name*="cover" i], textarea[name*="letter" i]').first();
      if (await messageField.count() > 0) {
        await messageField.fill(coverLetter);
      }

      const submitButton = await page.locator('button[type="submit"], input[type="submit"], button:has-text("Submit"), button:has-text("Apply")').first();
      if (await submitButton.count() > 0) {
        await submitButton.click();
      }

      await page.waitForTimeout(2000);
      logger.info('Generic application completed');
    } catch (error) {
      logger.error('Generic application failed', { error: error.message });
      throw error;
    }
  }

  async cleanup() {
    await this.selfHealing.cleanup();
  }
}

