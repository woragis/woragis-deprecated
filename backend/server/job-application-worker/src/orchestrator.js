import { logger } from './utils/logger.js';

export class Orchestrator {
  constructor(db) {
    this.db = db;
  }

  async shouldProcessWebsite(websiteName) {
    const website = await this.db.getWebsiteByName(websiteName);
    
    if (!website) {
      logger.warn('Website not found', { website: websiteName });
      return false;
    }

    if (!website.enabled) {
      return false;
    }

    // Check if should reset (new day)
    const lastReset = new Date(website.last_reset);
    const now = new Date();
    if (this.isNewDay(lastReset, now)) {
      await this.db.resetWebsiteCount(website.id);
      // Reload website
      const updated = await this.db.getWebsiteByName(websiteName);
      return updated.current_count < updated.daily_limit;
    }

    return website.current_count < website.daily_limit;
  }

  async incrementWebsiteCount(websiteName) {
    const website = await this.db.getWebsiteByName(websiteName);
    if (!website) {
      return;
    }

    const newCount = website.current_count + 1;
    await this.db.updateWebsiteCount(websiteName, newCount);
  }

  isNewDay(date1, date2) {
    return (
      date1.getFullYear() !== date2.getFullYear() ||
      date1.getMonth() !== date2.getMonth() ||
      date1.getDate() !== date2.getDate()
    );
  }
}

