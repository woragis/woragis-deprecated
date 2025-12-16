import { describe, test, expect, beforeEach, afterEach, jest } from '@jest/globals';

// Mock axios before importing the service
const mockAxiosPost = jest.fn();
jest.unstable_mockModule('axios', () => ({
  default: {
    post: mockAxiosPost,
  },
}));

// Import after mocking
const { CoverLetterService } = await import('../../src/coverLetter.js');

describe('CoverLetterService', () => {
  let service;
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
    service = new CoverLetterService();
    mockAxiosPost.mockClear();
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  describe('constructor', () => {
    test('should use default AI service URL', () => {
      delete process.env.AI_SERVICE_URL;
      service = new CoverLetterService();
      expect(service.aiServiceUrl).toBe('http://ai-service:8000');
    });

    test('should use environment AI service URL', () => {
      process.env.AI_SERVICE_URL = 'http://custom-ai:9000';
      service = new CoverLetterService();
      expect(service.aiServiceUrl).toBe('http://custom-ai:9000');
    });
  });

  describe('generateCoverLetter', () => {
    const mockProfile = {
      name: 'John Doe',
      email: 'john@example.com',
      skills: ['JavaScript', 'Node.js'],
      experience: '5 years',
    };

    const mockJobInfo = {
      companyName: 'Test Company',
      jobTitle: 'Senior Developer',
      jobDescription: 'Looking for an experienced developer',
    };

    test('should generate cover letter successfully', async () => {
      const mockResponse = {
        data: {
          message: {
            content: 'Dear Test Company, I am writing to apply...',
          },
        },
      };
      mockAxiosPost.mockResolvedValue(mockResponse);

      const result = await service.generateCoverLetter(mockProfile, mockJobInfo);

      expect(result).toBe('Dear Test Company, I am writing to apply...');
      expect(mockAxiosPost).toHaveBeenCalledWith(
        expect.stringContaining('/api/chat/completions'),
        expect.objectContaining({
          provider: 'openai',
          model: 'gpt-4o-mini',
        }),
        expect.objectContaining({
          timeout: 30000,
        })
      );
    });

    test('should handle response with choices array', async () => {
      const mockResponse = {
        data: {
          choices: [
            {
              message: {
                content: 'Alternative response format',
              },
            },
          ],
        },
      };
      mockAxiosPost.mockResolvedValue(mockResponse);

      const result = await service.generateCoverLetter(mockProfile, mockJobInfo);

      expect(result).toBe('Alternative response format');
    });

    test('should throw error when no content in response', async () => {
      const mockResponse = {
        data: {},
      };
      mockAxiosPost.mockResolvedValue(mockResponse);

      await expect(
        service.generateCoverLetter(mockProfile, mockJobInfo)
      ).rejects.toThrow('No cover letter content in AI response');
    });

    test('should handle API errors', async () => {
      const error = new Error('Network error');
      error.response = { data: { error: 'API error' } };
      mockAxiosPost.mockRejectedValue(error);

      await expect(
        service.generateCoverLetter(mockProfile, mockJobInfo)
      ).rejects.toThrow('Failed to generate cover letter: Network error');
    });

    test('should include profile and job info in prompt', async () => {
      mockAxiosPost.mockResolvedValue({
        data: { message: { content: 'Test' } },
      });

      await service.generateCoverLetter(mockProfile, mockJobInfo);

      const callArgs = mockAxiosPost.mock.calls[0];
      const prompt = callArgs[1].messages[0].content;

      // The prompt should contain job information
      expect(prompt).toContain('Test Company');
      expect(prompt).toContain('Senior Developer');
      // Prompt should be a non-empty string
      expect(typeof prompt).toBe('string');
      expect(prompt.length).toBeGreaterThan(0);
    });
  });

  describe('buildPrompt', () => {
    test('should build prompt with profile and job info', () => {
      const profile = {
        name: 'John Doe',
        email: 'john@example.com',
        skills: ['JavaScript'],
      };
      const jobInfo = {
        companyName: 'Test Company',
        jobTitle: 'Developer',
        jobDescription: 'Test description',
      };

      const prompt = service.buildPrompt(profile, jobInfo);

      // The prompt should contain job information
      expect(prompt).toContain('Test Company');
      expect(prompt).toContain('Developer');
      expect(prompt).toContain('Test description');
      // Prompt should be a string
      expect(typeof prompt).toBe('string');
      expect(prompt.length).toBeGreaterThan(0);
    });
  });
});
