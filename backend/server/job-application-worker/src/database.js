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
    const profile = {
      projects: [],
      posts: [],
      technicalWritings: [],
      caseStudies: [],
      certifications: [],
      skills: [],
      interests: [],
      projectCaseStudies: [],
    };

    // Fetch Projects with Technologies
    try {
      const projectsResult = await this.pool.query(
        `SELECT p.id, p.name, p.description, p.status, p.health_score, p.mrr, p.cac, p.ltv, p.churn_rate
         FROM projects p
         WHERE p.user_id = $1
         ORDER BY p.created_at DESC
         LIMIT 20`,
        [userId]
      );

      for (const project of projectsResult.rows) {
        // Fetch technologies for each project
        const techResult = await this.pool.query(
          `SELECT name, version, category, purpose
           FROM project_technologies
           WHERE project_id = $1
           ORDER BY category, name`,
          [project.id]
        );

        profile.projects.push({
          name: project.name,
          description: project.description,
          status: project.status,
          healthScore: project.health_score,
          metrics: {
            mrr: project.mrr,
            cac: project.cac,
            ltv: project.ltv,
            churnRate: project.churn_rate,
          },
          techStack: techResult.rows.map(tech => ({
            name: tech.name,
            version: tech.version,
            category: tech.category,
            purpose: tech.purpose,
          })),
        });
      }
    } catch (error) {
      logger.warn('Failed to fetch projects', { error: error.message });
    }

    // Fetch Posts with Skills/Tags
    try {
      const postsResult = await this.pool.query(
        `SELECT id, title, content, excerpt, status, published_at, featured
         FROM posts
         WHERE user_id = $1 AND status = 'published'
         ORDER BY published_at DESC
         LIMIT 15`,
        [userId]
      );

      for (const post of postsResult.rows) {
        // Fetch skills associated with post
        const postSkillsResult = await this.pool.query(
          `SELECT s.name, s.category, s.description
           FROM skills s
           INNER JOIN post_skills ps ON s.id = ps.skill_id
           WHERE ps.post_id = $1`,
          [post.id]
        );

        profile.posts.push({
          title: post.title,
          content: post.content?.substring(0, 1000) || post.excerpt || '', // Limit content length
          excerpt: post.excerpt,
          publishedAt: post.published_at,
          featured: post.featured,
          skills: postSkillsResult.rows.map(skill => skill.name),
        });
      }
    } catch (error) {
      logger.warn('Failed to fetch posts', { error: error.message });
    }

    // Fetch Technical Writings
    try {
      const writingsResult = await this.pool.query(
        `SELECT title, description, content, type, platform, url, published_at, featured
         FROM technical_writings
         WHERE user_id = $1
         ORDER BY published_at DESC
         LIMIT 15`,
        [userId]
      );

      profile.technicalWritings = writingsResult.rows.map(writing => ({
        title: writing.title,
        content: writing.content?.substring(0, 1000) || writing.description || '',
        description: writing.description,
        type: writing.type,
        platform: writing.platform,
        url: writing.url,
        publishedAt: writing.published_at,
        featured: writing.featured,
      }));
    } catch (error) {
      logger.warn('Failed to fetch technical writings', { error: error.message });
    }

    // Fetch Case Studies
    try {
      const caseStudiesResult = await this.pool.query(
        `SELECT title, problem, context, solution, approach, technologies, lessons_learned, featured
         FROM case_studies
         WHERE user_id = $1
         ORDER BY created_at DESC
         LIMIT 10`,
        [userId]
      );

      profile.caseStudies = caseStudiesResult.rows.map(cs => ({
        title: cs.title,
        problem: cs.problem?.substring(0, 500) || '',
        context: cs.context?.substring(0, 500) || '',
        solution: cs.solution?.substring(0, 1000) || '',
        approach: cs.approach || [],
        technologies: cs.technologies || [],
        lessonsLearned: cs.lessons_learned || [],
        featured: cs.featured,
      }));
    } catch (error) {
      logger.warn('Failed to fetch case studies', { error: error.message });
    }

    // Fetch Project Case Studies
    try {
      const projectCaseStudiesResult = await this.pool.query(
        `SELECT pcs.title, pcs.problem, pcs.context, pcs.solution, pcs.approach, 
                pcs.technologies, pcs.lessons_learned, p.name as project_name
         FROM project_case_studies pcs
         INNER JOIN projects p ON p.id = pcs.project_id
         WHERE p.user_id = $1
         ORDER BY pcs.created_at DESC
         LIMIT 10`,
        [userId]
      );

      profile.projectCaseStudies = projectCaseStudiesResult.rows.map(pcs => ({
        title: pcs.title,
        projectName: pcs.project_name,
        problem: pcs.problem?.substring(0, 500) || '',
        context: pcs.context?.substring(0, 500) || '',
        solution: pcs.solution?.substring(0, 1000) || '',
        approach: pcs.approach || [],
        technologies: pcs.technologies || [],
        lessonsLearned: pcs.lessons_learned || [],
      }));
    } catch (error) {
      logger.warn('Failed to fetch project case studies', { error: error.message });
    }

    // Fetch Certifications
    try {
      const certsResult = await this.pool.query(
        `SELECT name, issuer, issue_date, expiry_date, description, category, status, featured
         FROM certifications
         WHERE user_id = $1 AND status = 'active'
         ORDER BY issue_date DESC`,
        [userId]
      );

      profile.certifications = certsResult.rows.map(cert => ({
        name: cert.name,
        issuer: cert.issuer,
        issueDate: cert.issue_date,
        expiryDate: cert.expiry_date,
        description: cert.description,
        category: cert.category,
        featured: cert.featured,
      }));
    } catch (error) {
      logger.warn('Failed to fetch certifications', { error: error.message });
    }

    // Fetch Skills (all skills, not just user-specific)
    // Skills are global, but we can fetch featured ones or all
    try {
      const skillsResult = await this.pool.query(
        `SELECT DISTINCT s.name, s.category, s.description
         FROM skills s
         INNER JOIN project_technologies pt ON LOWER(pt.name) = LOWER(s.name)
         INNER JOIN projects p ON p.id = pt.project_id
         WHERE p.user_id = $1
         UNION
         SELECT DISTINCT s.name, s.category, s.description
         FROM skills s
         INNER JOIN post_skills ps ON s.id = ps.skill_id
         INNER JOIN posts p ON p.id = ps.post_id
         WHERE p.user_id = $1
         ORDER BY category, name
         LIMIT 50`,
        [userId]
      );

      profile.skills = skillsResult.rows.map(skill => ({
        name: skill.name,
        category: skill.category,
        description: skill.description,
      }));
    } catch (error) {
      logger.warn('Failed to fetch skills', { error: error.message });
    }

    // Fetch Interests (interests are global, not user-specific)
    try {
      const interestsResult = await this.pool.query(
        `SELECT title, description, featured
         FROM interests
         WHERE featured = true
         ORDER BY created_at DESC
         LIMIT 10`
      );

      profile.interests = interestsResult.rows.map(interest => ({
        title: interest.title,
        description: interest.description,
      }));
    } catch (error) {
      logger.warn('Failed to fetch interests', { error: error.message });
    }

    logger.info('User profile fetched', {
      userId,
      projectsCount: profile.projects.length,
      postsCount: profile.posts.length,
      writingsCount: profile.technicalWritings.length,
      caseStudiesCount: profile.caseStudies.length,
      certificationsCount: profile.certifications.length,
      skillsCount: profile.skills.length,
    });

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

