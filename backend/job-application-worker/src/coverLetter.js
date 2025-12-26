import axios from 'axios';
import { logger } from './utils/logger.js';
import { validateString, validateNoSQLInjection, validateNoXSS } from './validation.js';

// Note: Import logger if not already imported

export class CoverLetterService {
  constructor() {
    this.aiServiceUrl = process.env.AI_SERVICE_URL || 'http://ai-service:8000';
  }

  async generateCoverLetter(profile, jobInfo) {
    // Validate inputs
    if (!jobInfo || !jobInfo.companyName || !jobInfo.jobTitle) {
      throw new Error('jobInfo must contain companyName and jobTitle');
    }
    validateString(jobInfo.companyName, 1, 200, 'companyName');
    validateString(jobInfo.jobTitle, 1, 200, 'jobTitle');
    validateNoSQLInjection(jobInfo.companyName, 'companyName');
    validateNoSQLInjection(jobInfo.jobTitle, 'jobTitle');
    validateNoXSS(jobInfo.companyName, 'companyName');
    validateNoXSS(jobInfo.jobTitle, 'jobTitle');

    if (jobInfo.location) {
      validateString(jobInfo.location, 1, 200, 'location');
      validateNoSQLInjection(jobInfo.location, 'location');
      validateNoXSS(jobInfo.location, 'location');
    }

    if (jobInfo.jobDescription) {
      validateString(jobInfo.jobDescription, 1, 50000, 'jobDescription');
      validateNoSQLInjection(jobInfo.jobDescription, 'jobDescription');
      validateNoXSS(jobInfo.jobDescription, 'jobDescription');
    }

    logger.info('Generating cover letter', {
      company: jobInfo.companyName,
      jobTitle: jobInfo.jobTitle,
    });

    const prompt = this.buildPrompt(profile, jobInfo);

    try {
      const response = await axios.post(
        `${this.aiServiceUrl}/api/chat/completions`,
        {
          provider: 'openai',
          model: 'gpt-4o-mini',
          temperature: 0.7,
          messages: [
            {
              role: 'user',
              content: prompt,
            },
          ],
          max_tokens: 1500,
        },
        {
          timeout: 30000,
        }
      );

      const coverLetter = response.data.message?.content || response.data.choices?.[0]?.message?.content;
      
      if (!coverLetter) {
        throw new Error('No cover letter content in AI response');
      }

      logger.info('Cover letter generated', { length: coverLetter.length });
      return coverLetter.trim();
    } catch (error) {
      logger.error('Failed to generate cover letter', {
        error: error.message,
        response: error.response?.data,
      });
      throw new Error(`Failed to generate cover letter: ${error.message}`);
    }
  }

  buildPrompt(profile, jobInfo) {
    let prompt = `You are a professional cover letter writer. Write a personalized cover letter for the following job application.

Job Information:
- Company: ${jobInfo.companyName}
- Position: ${jobInfo.jobTitle}
- Location: ${jobInfo.location}
- Job Description: ${jobInfo.jobDescription || 'Not provided'}

Applicant Profile:
`;

    // Add Projects with Technologies
    if (profile.projects && profile.projects.length > 0) {
      prompt += '\n## Projects & Experience:\n';
      profile.projects.slice(0, 5).forEach(project => {
        prompt += `- **${project.name}**: ${project.description || 'No description'}`;
        if (project.techStack && project.techStack.length > 0) {
          const techNames = project.techStack.map(t => t.name).join(', ');
          prompt += `\n  Technologies: ${techNames}`;
        }
        if (project.metrics && (project.metrics.mrr > 0 || project.metrics.healthScore > 0)) {
          prompt += `\n  Status: ${project.status}, Health Score: ${project.metrics.healthScore || 'N/A'}`;
        }
        prompt += '\n';
      });
    }

    // Add Case Studies
    if (profile.caseStudies && profile.caseStudies.length > 0) {
      prompt += '\n## Case Studies:\n';
      profile.caseStudies.slice(0, 3).forEach(cs => {
        prompt += `- **${cs.title}**: ${cs.problem || 'Problem solved'}\n`;
        prompt += `  Solution: ${cs.solution?.substring(0, 300) || 'Complex technical solution'}\n`;
        if (cs.technologies && cs.technologies.length > 0) {
          prompt += `  Technologies: ${cs.technologies.join(', ')}\n`;
        }
        if (cs.lessonsLearned && cs.lessonsLearned.length > 0) {
          prompt += `  Key Learnings: ${cs.lessonsLearned.slice(0, 2).join('; ')}\n`;
        }
      });
    }

    // Add Project Case Studies
    if (profile.projectCaseStudies && profile.projectCaseStudies.length > 0) {
      prompt += '\n## Project Case Studies:\n';
      profile.projectCaseStudies.slice(0, 3).forEach(pcs => {
        prompt += `- **${pcs.title}** (Project: ${pcs.projectName}): ${pcs.solution?.substring(0, 300) || 'Technical solution'}\n`;
        if (pcs.technologies && pcs.technologies.length > 0) {
          prompt += `  Technologies: ${pcs.technologies.join(', ')}\n`;
        }
      });
    }

    // Add Technical Writings
    if (profile.technicalWritings && profile.technicalWritings.length > 0) {
      prompt += '\n## Technical Writings & Publications:\n';
      profile.technicalWritings.slice(0, 5).forEach(writing => {
        prompt += `- **${writing.title}** (${writing.platform}): ${writing.description || writing.content?.substring(0, 200) || ''}\n`;
        if (writing.url) {
          prompt += `  URL: ${writing.url}\n`;
        }
      });
    }

    // Add Posts
    if (profile.posts && profile.posts.length > 0) {
      prompt += '\n## Blog Posts & Articles:\n';
      profile.posts.slice(0, 5).forEach(post => {
        prompt += `- **${post.title}**: ${post.excerpt || post.content?.substring(0, 200) || ''}\n`;
        if (post.skills && post.skills.length > 0) {
          prompt += `  Topics: ${post.skills.join(', ')}\n`;
        }
      });
    }

    // Add Skills (comprehensive list)
    if (profile.skills && profile.skills.length > 0) {
      const skillsByCategory = {};
      profile.skills.forEach(skill => {
        if (!skillsByCategory[skill.category]) {
          skillsByCategory[skill.category] = [];
        }
        skillsByCategory[skill.category].push(skill.name);
      });

      prompt += '\n## Technical Skills:\n';
      Object.entries(skillsByCategory).forEach(([category, skills]) => {
        prompt += `- **${category}**: ${skills.join(', ')}\n`;
      });
    }

    // Add Certifications
    if (profile.certifications && profile.certifications.length > 0) {
      prompt += '\n## Certifications:\n';
      profile.certifications.forEach(cert => {
        prompt += `- **${cert.name}** from ${cert.issuer}`;
        if (cert.issueDate) {
          prompt += ` (Issued: ${new Date(cert.issueDate).getFullYear()})`;
        }
        if (cert.description) {
          prompt += `: ${cert.description}`;
        }
        prompt += '\n';
      });
    }

    // Add Interests
    if (profile.interests && profile.interests.length > 0) {
      prompt += '\n## Interests & Focus Areas:\n';
      profile.interests.forEach(interest => {
        prompt += `- **${interest.title}**: ${interest.description}\n`;
      });
    }

    prompt += `
## Instructions:
1. Write a professional, engaging cover letter (3-4 paragraphs)
2. **Highlight relevant experience** from the projects, case studies, and technical writings
3. **Mention specific technologies** that match the job requirements
4. **Reference concrete achievements** from case studies and projects
5. **Show enthusiasm** for the specific role and company
6. **Use a professional but personable tone**
7. **Do NOT include placeholders or generic statements** - use real information from the profile
8. **Tailor the letter** to the specific job requirements mentioned
9. If the job description mentions specific technologies, prioritize mentioning those from the profile

Write the cover letter now:`;

    return prompt;
  }
}

