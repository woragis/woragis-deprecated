import { Queue } from './queue.js';
import { Database } from './database.js';
import { Orchestrator } from './orchestrator.js';
import { Scraper } from './scraper.js';
import { CoverLetterService } from './coverLetter.js';
import { logger } from './utils/logger.js';

export class Worker {
  constructor() {
    this.queue = new Queue();
    this.db = new Database();
    this.orchestrator = new Orchestrator(this.db);
    this.scraper = new Scraper();
    this.coverLetterService = new CoverLetterService();
    this.running = false;
  }

  async start() {
    logger.info('Starting job application worker');
    this.running = true;

    // Initialize connections
    await this.queue.connect();
    await this.db.connect();
    await this.scraper.initialize();

    // Start processing loop
    this.processLoop();
  }

  async processLoop() {
    while (this.running) {
      try {
        // Dequeue job with 5 second timeout
        const job = await this.queue.dequeueJob(5000);
        
        if (!job) {
          // No job available, continue polling
          continue;
        }

        logger.info('Processing job application', {
          jobId: job.id,
          company: job.companyName,
          website: job.website,
        });

        // Check if we should process this website (rate limit check)
        const shouldProcess = await this.orchestrator.shouldProcessWebsite(job.website);
        
        if (!shouldProcess) {
          logger.info('Website limit reached, re-enqueuing job', {
            website: job.website,
          });
          // Re-enqueue for later
          await this.queue.enqueueJob(job);
          // Wait before checking again
          await this.sleep(3600000); // 1 hour
          continue;
        }

        // Process the job
        await this.processApplication(job);

        // Increment website count
        await this.orchestrator.incrementWebsiteCount(job.website);

        // Mark job as complete
        await this.queue.markJobComplete(job.id);

        logger.info('Job application completed', { jobId: job.id });
      } catch (error) {
        logger.error('Error processing job', { error: error.message, stack: error.stack });
        // Continue processing other jobs
      }
    }
  }

  async processApplication(job) {
    // Find or create application record
    let application = await this.db.findApplicationByUrl(job.jobUrl, job.website);
    
    if (!application) {
      application = await this.db.createApplication({
        userId: job.userId,
        companyName: job.companyName,
        location: job.location,
        jobTitle: job.jobTitle,
        jobUrl: job.jobUrl,
        website: job.website,
        status: 'processing',
      });
    } else {
      await this.db.updateApplicationStatus(application.id, 'processing');
    }

    // Fetch user profile data
    const profile = await this.db.fetchUserProfile(job.userId);

    // Generate cover letter
    const jobInfo = {
      companyName: job.companyName,
      jobTitle: job.jobTitle,
      jobDescription: '', // TODO: Fetch from job URL if possible
      location: job.location,
      requirements: [], // TODO: Extract from job description
    };

    const coverLetter = await this.coverLetterService.generateCoverLetter(profile, jobInfo);

    // Apply to job using Playwright
    try {
      await this.scraper.applyToJob(job, coverLetter);
      
      // Mark as applied
      await this.db.updateApplication(application.id, {
        status: 'applied',
        coverLetter: coverLetter,
        appliedAt: new Date(),
      });
    } catch (error) {
      // Mark as failed
      await this.db.updateApplication(application.id, {
        status: 'failed',
        errorMessage: error.message,
      });
      throw error;
    }
  }

  stop() {
    logger.info('Stopping job application worker');
    this.running = false;
    this.queue.disconnect();
    this.db.disconnect();
    this.scraper.cleanup();
  }

  sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}

