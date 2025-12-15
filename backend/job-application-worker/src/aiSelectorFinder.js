import axios from 'axios';
import { logger } from './utils/logger.js';

export class AISelectorFinder {
  constructor() {
    this.aiServiceUrl = process.env.AI_SERVICE_URL || 'http://ai-service:8000';
  }

  /**
   * Find selectors using AI by analyzing HTML
   * @param {string} html - HTML content of the page
   * @param {string} description - What we're looking for (e.g., "Easy Apply button", "email input field")
   * @param {string} website - Website name for context
   * @returns {Promise<Object>} Selector data with multiple strategies
   */
  async findSelectorsFromHTML(html, description, website) {
    logger.info('Finding selectors using AI (HTML analysis)', {
      website,
      description,
      htmlLength: html.length,
    });

    const prompt = this.buildHTMLPrompt(html, description, website);

    try {
      const response = await axios.post(
        `${this.aiServiceUrl}/api/chat/completions`,
        {
          provider: 'openai',
          model: 'gpt-4o-mini',
          temperature: 0.3, // Lower temperature for more consistent selectors
          messages: [
            {
              role: 'user',
              content: prompt,
            },
          ],
          max_tokens: 1000,
        },
        {
          timeout: 30000,
        }
      );

      const content = response.data.message?.content || response.data.choices?.[0]?.message?.content;
      
      if (!content) {
        throw new Error('No content in AI response');
      }

      // Parse AI response to extract selectors
      const selectors = this.parseSelectorResponse(content);
      
      logger.info('AI found selectors', {
        website,
        description,
        selectors: Object.keys(selectors),
      });

      return selectors;
    } catch (error) {
      logger.error('Failed to find selectors using AI', {
        error: error.message,
        website,
        description,
      });
      throw new Error(`AI selector finding failed: ${error.message}`);
    }
  }

  /**
   * Find selectors using AI by analyzing screenshot (vision model)
   * @param {Buffer} screenshot - Screenshot image buffer
   * @param {string} description - What we're looking for
   * @param {string} website - Website name
   * @returns {Promise<Object>} Selector data with coordinates/selectors
   */
  async findSelectorsFromScreenshot(screenshot, description, website) {
    logger.info('Finding selectors using AI (vision)', {
      website,
      description,
    });

    // Convert screenshot to base64
    const base64Image = screenshot.toString('base64');

    try {
      const response = await axios.post(
        `${this.aiServiceUrl}/api/chat/completions`,
        {
          provider: 'openai',
          model: 'gpt-4o', // Use vision-capable model
          temperature: 0.3,
          messages: [
            {
              role: 'user',
              content: [
                {
                  type: 'text',
                  text: this.buildVisionPrompt(description, website),
                },
                {
                  type: 'image_url',
                  image_url: {
                    url: `data:image/png;base64,${base64Image}`,
                  },
                },
              ],
            },
          ],
          max_tokens: 1000,
        },
        {
          timeout: 60000, // Vision models can be slower
        }
      );

      const content = response.data.message?.content || response.data.choices?.[0]?.message?.content;
      
      if (!content) {
        throw new Error('No content in AI response');
      }

      const selectors = this.parseSelectorResponse(content);
      
      logger.info('AI found selectors from screenshot', {
        website,
        description,
        selectors: Object.keys(selectors),
      });

      return selectors;
    } catch (error) {
      logger.error('Failed to find selectors using vision AI', {
        error: error.message,
        website,
        description,
      });
      throw new Error(`AI vision selector finding failed: ${error.message}`);
    }
  }

  buildHTMLPrompt(html, description, website) {
    // Truncate HTML if too long (keep structure)
    const truncatedHTML = html.length > 50000 
      ? html.substring(0, 50000) + '... [truncated]'
      : html;

    return `You are a web scraping expert. Analyze the HTML below and find the best selectors for: "${description}"

Website: ${website}
Looking for: ${description}

HTML:
${truncatedHTML}

Instructions:
1. Find the element that matches the description
2. Provide multiple selector strategies (CSS selector, XPath, text content, etc.)
3. Return a JSON object with this structure:
{
  "primary": "best CSS selector",
  "alternatives": ["alternative selector 1", "alternative selector 2"],
  "xpath": "XPath selector if useful",
  "text": "text content to match if applicable",
  "attributes": {"data-testid": "value", "aria-label": "value"},
  "explanation": "why this selector was chosen"
}

Return ONLY the JSON object, no other text.`;
  }

  buildVisionPrompt(description, website) {
    return `You are a web scraping expert. Analyze this screenshot and find the best way to locate: "${description}"

Website: ${website}
Looking for: ${description}

Instructions:
1. Identify the element in the screenshot
2. Provide selectors and coordinates
3. Return a JSON object with this structure:
{
  "primary": "best CSS selector",
  "alternatives": ["alternative selector 1", "alternative selector 2"],
  "coordinates": {"x": 100, "y": 200} if you can see it,
  "text": "visible text if applicable",
  "explanation": "description of the element"
}

Return ONLY the JSON object, no other text.`;
  }

  parseSelectorResponse(content) {
    try {
      // Try to extract JSON from the response
      const jsonMatch = content.match(/\{[\s\S]*\}/);
      if (jsonMatch) {
        return JSON.parse(jsonMatch[0]);
      }
      
      // Fallback: try parsing the whole content
      return JSON.parse(content);
    } catch (error) {
      logger.warn('Failed to parse AI selector response', {
        error: error.message,
        content: content.substring(0, 200),
      });
      
      // Fallback: return a basic structure
      return {
        primary: content.trim(),
        alternatives: [],
        explanation: 'AI response could not be parsed, using raw content',
      };
    }
  }
}

