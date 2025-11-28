import pg from 'pg';
const { Pool } = pg;
import { logger } from './utils/logger.js';

export class Database {
  constructor() {
    this.pool = null;
  }

  async connect() {
    const connectionString = process.env.DATABASE_URL;
    if (!connectionString) {
      throw new Error('DATABASE_URL environment variable is required');
    }

    this.pool = new Pool({
      connectionString,
      ssl: process.env.DATABASE_SSL === 'true' ? { rejectUnauthorized: false } : false,
    });

    // Test connection
    await this.pool.query('SELECT NOW()');
    logger.info('Connected to database');
  }

  async findApplicationByUrl(jobUrl, website) {
    const result = await this.pool.query(
      'SELECT * FROM job_applications WHERE job_url = $1 AND website = $2 LIMIT 1',
      [jobUrl, website]
    );
    return result.rows[0] || null;
  }

  async createApplication(data) {
    const result = await this.pool.query(
      `INSERT INTO job_applications 
       (id, user_id, company_name, location, job_title, job_url, website, status, created_at, updated_at)
       VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
       RETURNING *`,
      [
        data.userId,
        data.companyName,
        data.location,
        data.jobTitle,
        data.jobUrl,
        data.website,
        data.status || 'pending',
      ]
    );
    return result.rows[0];
  }

  async updateApplicationStatus(id, status) {
    const result = await this.pool.query(
      'UPDATE job_applications SET status = $1, updated_at = NOW() WHERE id = $2 RETURNING *',
      [status, id]
    );
    return result.rows[0];
  }

  async updateApplication(id, data) {
    const updates = [];
    const values = [];
    let paramCount = 1;

    if (data.status) {
      updates.push(`status = $${paramCount++}`);
      values.push(data.status);
    }
    if (data.coverLetter) {
      updates.push(`cover_letter = $${paramCount++}`);
      values.push(data.coverLetter);
    }
    if (data.appliedAt) {
      updates.push(`applied_at = $${paramCount++}`);
      values.push(data.appliedAt);
    }
    if (data.errorMessage) {
      updates.push(`error_message = $${paramCount++}`);
      values.push(data.errorMessage);
    }

    updates.push(`updated_at = NOW()`);
    values.push(id);

    const query = `UPDATE job_applications SET ${updates.join(', ')} WHERE id = $${paramCount} RETURNING *`;
    const result = await this.pool.query(query, values);
    return result.rows[0];
  }

  async fetchUserProfile(userId) {
    // TODO: Fetch projects, posts, technical writings, skills, etc.
    // This is a placeholder - implement based on your schema
    const profile = {
      projects: [],
      posts: [],
      technicalWritings: [],
      skills: [],
      interests: [],
      certifications: [],
    };

    // Example: Fetch projects
    try {
      const projectsResult = await this.pool.query(
        'SELECT name, description FROM projects WHERE user_id = $1',
        [userId]
      );
      profile.projects = projectsResult.rows.map(row => ({
        name: row.name,
        description: row.description,
        techStack: [], // TODO: Fetch from project_technologies
      }));
    } catch (error) {
      logger.warn('Failed to fetch projects', { error: error.message });
    }

    // TODO: Fetch other profile data similarly

    return profile;
  }

  async getWebsiteByName(name) {
    const result = await this.pool.query(
      'SELECT * FROM job_websites WHERE name = $1 LIMIT 1',
      [name]
    );
    return result.rows[0] || null;
  }

  async updateWebsiteCount(name, count) {
    await this.pool.query(
      'UPDATE job_websites SET current_count = $1, updated_at = NOW() WHERE name = $2',
      [count, name]
    );
  }

  async resetWebsiteCount(id) {
    await this.pool.query(
      'UPDATE job_websites SET current_count = 0, last_reset = NOW(), updated_at = NOW() WHERE id = $1',
      [id]
    );
  }

  async disconnect() {
    if (this.pool) {
      await this.pool.end();
      logger.info('Disconnected from database');
    }
  }
}

