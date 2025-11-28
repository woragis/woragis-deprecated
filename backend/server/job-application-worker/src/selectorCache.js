import Redis from 'ioredis';
import { logger } from './utils/logger.js';

const SELECTOR_PREFIX = 'scraper:selectors:';
const DEFAULT_TTL = 7 * 24 * 60 * 60; // 7 days in seconds

export class SelectorCache {
  constructor() {
    this.client = null;
  }

  async connect() {
    const redisUrl = process.env.REDIS_URL || 'redis://localhost:6379/0';
    this.client = new Redis(redisUrl);
    
    this.client.on('error', (err) => {
      logger.error('SelectorCache Redis error', { error: err.message });
    });

    await this.client.ping();
  }

  /**
   * Get cached selectors for a website and action
   * @param {string} website - Website name (e.g., 'linkedin')
   * @param {string} action - Action name (e.g., 'easy-apply-button', 'email-field')
   * @returns {Promise<Object|null>} Cached selector data or null
   */
  async getSelectors(website, action) {
    const key = this.getKey(website, action);
    const data = await this.client.get(key);
    
    if (!data) {
      return null;
    }

    try {
      return JSON.parse(data);
    } catch (error) {
      logger.warn('Failed to parse cached selectors', { website, action, error: error.message });
      return null;
    }
  }

  /**
   * Store selectors for a website and action
   * @param {string} website - Website name
   * @param {string} action - Action name
   * @param {Object} selectors - Selector data (can include multiple strategies)
   * @param {number} ttl - Time to live in seconds (default: 7 days)
   */
  async setSelectors(website, action, selectors, ttl = DEFAULT_TTL) {
    const key = this.getKey(website, action);
    const data = JSON.stringify({
      ...selectors,
      cachedAt: new Date().toISOString(),
      website,
      action,
    });

    await this.client.setex(key, ttl, data);
    logger.info('Cached selectors', { website, action, ttl });
  }

  /**
   * Invalidate selectors (delete from cache)
   * Called when selectors fail
   */
  async invalidateSelectors(website, action) {
    const key = this.getKey(website, action);
    await this.client.del(key);
    logger.info('Invalidated selectors', { website, action });
  }

  /**
   * Get all cached selectors for a website (for debugging)
   */
  async getAllSelectors(website) {
    const pattern = this.getKey(website, '*');
    const keys = await this.client.keys(pattern);
    
    if (keys.length === 0) {
      return {};
    }

    const values = await this.client.mget(keys);
    const result = {};

    keys.forEach((key, index) => {
      const action = key.replace(SELECTOR_PREFIX + website + ':', '');
      try {
        result[action] = JSON.parse(values[index]);
      } catch (error) {
        logger.warn('Failed to parse selector', { key, error: error.message });
      }
    });

    return result;
  }

  getKey(website, action) {
    return `${SELECTOR_PREFIX}${website}:${action}`;
  }

  async disconnect() {
    if (this.client) {
      await this.client.quit();
    }
  }
}

