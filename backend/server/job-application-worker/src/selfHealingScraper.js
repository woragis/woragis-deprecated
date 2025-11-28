import { chromium } from 'playwright';
import { logger } from './utils/logger.js';
import { SelectorCache } from './selectorCache.js';
import { AISelectorFinder } from './aiSelectorFinder.js';

/**
 * Self-healing scraper that uses AI to find selectors when cached ones fail
 */
export class SelfHealingScraper {
  constructor() {
    this.selectorCache = new SelectorCache();
    this.aiFinder = new AISelectorFinder();
    this.browser = null;
  }

  async initialize() {
    await this.selectorCache.connect();
  }

  /**
   * Find an element using cached selectors, with AI fallback
   * @param {Page} page - Playwright page
   * @param {string} website - Website name
   * @param {string} action - Action name (e.g., 'easy-apply-button')
   * @param {string} description - Human-readable description
   * @returns {Promise<ElementHandle|null>} Found element or null
   */
  async findElement(page, website, action, description) {
    // Try cached selectors first
    const cached = await this.selectorCache.getSelectors(website, action);
    
    if (cached) {
      logger.info('Using cached selectors', { website, action });
      const element = await this.trySelectors(page, cached);
      
      if (element) {
        return element;
      }

      // Cached selectors failed, invalidate and try AI
      logger.warn('Cached selectors failed, invalidating', { website, action });
      await this.selectorCache.invalidateSelectors(website, action);
    }

    // Use AI to find new selectors
    logger.info('Finding new selectors using AI', { website, action, description });
    
    let selectors;
    try {
      // Try HTML analysis first (faster, cheaper)
      const html = await page.content();
      selectors = await this.aiFinder.findSelectorsFromHTML(html, description, website);
    } catch (error) {
      logger.warn('HTML analysis failed, trying vision', { error: error.message });
      
      // Fallback to vision if HTML analysis fails
      try {
        const screenshot = await page.screenshot();
        selectors = await this.aiFinder.findSelectorsFromScreenshot(screenshot, description, website);
      } catch (visionError) {
        logger.error('Both HTML and vision analysis failed', {
          htmlError: error.message,
          visionError: visionError.message,
        });
        throw new Error('AI selector finding failed');
      }
    }

    // Try the new selectors
    const element = await this.trySelectors(page, selectors);
    
    if (!element) {
      throw new Error(`Could not find element: ${description}`);
    }

    // Cache the successful selectors
    await this.selectorCache.setSelectors(website, action, selectors);
    logger.info('Cached new selectors', { website, action });

    return element;
  }

  /**
   * Try multiple selector strategies
   * @param {Page} page - Playwright page
   * @param {Object} selectors - Selector data with primary, alternatives, etc.
   * @returns {Promise<ElementHandle|null>} Found element or null
   */
  async trySelectors(page, selectors) {
    const strategies = [
      selectors.primary,
      ...(selectors.alternatives || []),
      selectors.xpath,
    ].filter(Boolean);

    for (const selector of strategies) {
      try {
        // Try CSS selector
        if (selector.startsWith('//') || selector.startsWith('(//')) {
          // XPath
          const element = await page.locator(selector).first();
          if (await element.count() > 0) {
            return element;
          }
        } else {
          // CSS selector
          const element = await page.locator(selector).first();
          if (await element.count() > 0) {
            return element;
          }
        }
      } catch (error) {
        // Try next strategy
        continue;
      }
    }

    // Try text-based search if available
    if (selectors.text) {
      try {
        const element = await page.getByText(selectors.text, { exact: false }).first();
        if (await element.count() > 0) {
          return element;
        }
      } catch (error) {
        // Ignore
      }
    }

    return null;
  }

  /**
   * Click an element with self-healing
   */
  async clickElement(page, website, action, description) {
    const element = await this.findElement(page, website, action, description);
    
    if (!element) {
      throw new Error(`Could not find element to click: ${description}`);
    }

    await element.click();
    logger.info('Clicked element', { website, action, description });
  }

  /**
   * Fill an input field with self-healing
   */
  async fillField(page, website, action, description, value) {
    const element = await this.findElement(page, website, action, description);
    
    if (!element) {
      throw new Error(`Could not find field to fill: ${description}`);
    }

    await element.fill(value);
    logger.info('Filled field', { website, action, description });
  }

  /**
   * Wait for an element with self-healing
   */
  async waitForElement(page, website, action, description, timeout = 10000) {
    const element = await this.findElement(page, website, action, description);
    
    if (!element) {
      throw new Error(`Could not find element to wait for: ${description}`);
    }

    await element.waitFor({ state: 'visible', timeout });
    logger.info('Element appeared', { website, action, description });
  }

  async cleanup() {
    await this.selectorCache.disconnect();
    if (this.browser) {
      await this.browser.close();
    }
  }
}

