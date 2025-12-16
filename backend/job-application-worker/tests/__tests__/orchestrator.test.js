import { describe, test, expect, beforeEach, jest } from '@jest/globals';
import { Orchestrator } from '../../src/orchestrator.js';

describe('Orchestrator', () => {
  let orchestrator;
  let mockDb;

  beforeEach(() => {
    mockDb = {
      getWebsiteByName: jest.fn(),
      resetWebsiteCount: jest.fn(),
      updateWebsiteCount: jest.fn(),
    };
    orchestrator = new Orchestrator(mockDb);
  });

  describe('shouldProcessWebsite', () => {
    test('should return false when website is not found', async () => {
      mockDb.getWebsiteByName.mockResolvedValue(null);

      const result = await orchestrator.shouldProcessWebsite('unknown-site');

      expect(result).toBe(false);
      expect(mockDb.getWebsiteByName).toHaveBeenCalledWith('unknown-site');
    });

    test('should return false when website is disabled', async () => {
      mockDb.getWebsiteByName.mockResolvedValue({
        id: 1,
        name: 'test-site',
        enabled: false,
        current_count: 0,
        daily_limit: 10,
        last_reset: new Date(),
      });

      const result = await orchestrator.shouldProcessWebsite('test-site');

      expect(result).toBe(false);
    });

    test('should return true when website is enabled and under limit', async () => {
      mockDb.getWebsiteByName.mockResolvedValue({
        id: 1,
        name: 'test-site',
        enabled: true,
        current_count: 5,
        daily_limit: 10,
        last_reset: new Date(),
      });

      const result = await orchestrator.shouldProcessWebsite('test-site');

      expect(result).toBe(true);
    });

    test('should return false when website is at limit', async () => {
      mockDb.getWebsiteByName.mockResolvedValue({
        id: 1,
        name: 'test-site',
        enabled: true,
        current_count: 10,
        daily_limit: 10,
        last_reset: new Date(),
      });

      const result = await orchestrator.shouldProcessWebsite('test-site');

      expect(result).toBe(false);
    });

    test('should reset count and return true when new day', async () => {
      const yesterday = new Date();
      yesterday.setDate(yesterday.getDate() - 1);

      mockDb.getWebsiteByName
        .mockResolvedValueOnce({
          id: 1,
          name: 'test-site',
          enabled: true,
          current_count: 10,
          daily_limit: 10,
          last_reset: yesterday,
        })
        .mockResolvedValueOnce({
          id: 1,
          name: 'test-site',
          enabled: true,
          current_count: 0,
          daily_limit: 10,
          last_reset: new Date(),
        });

      mockDb.resetWebsiteCount.mockResolvedValue(undefined);

      const result = await orchestrator.shouldProcessWebsite('test-site');

      expect(result).toBe(true);
      expect(mockDb.resetWebsiteCount).toHaveBeenCalledWith(1);
    });
  });

  describe('incrementWebsiteCount', () => {
    test('should increment website count', async () => {
      mockDb.getWebsiteByName.mockResolvedValue({
        id: 1,
        name: 'test-site',
        current_count: 5,
      });
      mockDb.updateWebsiteCount.mockResolvedValue(undefined);

      await orchestrator.incrementWebsiteCount('test-site');

      expect(mockDb.getWebsiteByName).toHaveBeenCalledWith('test-site');
      expect(mockDb.updateWebsiteCount).toHaveBeenCalledWith('test-site', 6);
    });

    test('should not increment when website is not found', async () => {
      mockDb.getWebsiteByName.mockResolvedValue(null);

      await orchestrator.incrementWebsiteCount('unknown-site');

      expect(mockDb.updateWebsiteCount).not.toHaveBeenCalled();
    });
  });

  describe('isNewDay', () => {
    test('should return true for different days', () => {
      const date1 = new Date('2024-01-01');
      const date2 = new Date('2024-01-02');

      const result = orchestrator.isNewDay(date1, date2);

      expect(result).toBe(true);
    });

    test('should return false for same day', () => {
      const date1 = new Date('2024-01-01T10:00:00');
      const date2 = new Date('2024-01-01T15:00:00');

      const result = orchestrator.isNewDay(date1, date2);

      expect(result).toBe(false);
    });

    test('should return true for different months', () => {
      const date1 = new Date('2024-01-01');
      const date2 = new Date('2024-02-01');

      const result = orchestrator.isNewDay(date1, date2);

      expect(result).toBe(true);
    });

    test('should return true for different years', () => {
      const date1 = new Date('2023-12-31');
      const date2 = new Date('2024-01-01');

      const result = orchestrator.isNewDay(date1, date2);

      expect(result).toBe(true);
    });
  });
});
