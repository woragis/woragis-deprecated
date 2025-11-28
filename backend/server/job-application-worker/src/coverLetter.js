import axios from 'axios';
import { logger } from './utils/logger.js';

export class CoverLetterService {
  constructor() {
    this.aiServiceUrl = process.env.AI_SERVICE_URL || 'http://ai-service:8000';
  }

  async generateCoverLetter(profile, jobInfo) {
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

    // Add projects
    if (profile.projects && profile.projects.length > 0) {
      prompt += '\nProjects:\n';
      profile.projects.forEach(project => {
        prompt += `- ${project.name}: ${project.description}`;
        if (project.techStack && project.techStack.length > 0) {
          prompt += ` (Tech: ${project.techStack.join(', ')})`;
        }
        prompt += '\n';
      });
    }

    // Add skills
    if (profile.skills && profile.skills.length > 0) {
      prompt += `\nSkills: ${profile.skills.join(', ')}\n`;
    }

    // Add technical writings
    if (profile.technicalWritings && profile.technicalWritings.length > 0) {
      prompt += '\nTechnical Writings:\n';
      profile.technicalWritings.forEach(writing => {
        const preview = writing.content?.substring(0, 200) || '';
        prompt += `- ${writing.title}: ${preview}\n`;
      });
    }

    prompt += `
Instructions:
1. Write a professional, engaging cover letter
2. Highlight relevant experience and skills from the applicant's profile
3. Show enthusiasm for the specific role and company
4. Keep it concise (3-4 paragraphs)
5. Use a professional but personable tone
6. Do not include placeholders or generic statements

Write the cover letter now:`;

    return prompt;
  }
}

