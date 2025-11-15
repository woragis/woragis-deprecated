/**
 * Utility functions for generating cron expressions
 * Cron format: minute hour day-of-month month day-of-week
 * Day of week: 0 = Sunday, 1 = Monday, ..., 6 = Saturday
 */

export interface CronPreset {
	name: string;
	description: string;
	cron: string;
	frequency: string;
}

export const WEEKDAYS = [
	{ value: '0', label: 'Sunday' },
	{ value: '1', label: 'Monday' },
	{ value: '2', label: 'Tuesday' },
	{ value: '3', label: 'Wednesday' },
	{ value: '4', label: 'Thursday' },
	{ value: '5', label: 'Friday' },
	{ value: '6', label: 'Saturday' }
] as const;

/**
 * Generate cron expression for daily schedule
 */
export function dailyCron(hour: number = 8, minute: number = 0): string {
	return `${minute} ${hour} * * *`;
}

/**
 * Generate cron expression for weekly schedule on specific weekday
 */
export function weeklyCron(weekday: number, hour: number = 8, minute: number = 0): string {
	return `${minute} ${hour} * * ${weekday}`;
}

/**
 * Generate cron expression for monthly schedule on specific day
 */
export function monthlyCron(dayOfMonth: number, hour: number = 8, minute: number = 0): string {
	return `${minute} ${hour} ${dayOfMonth} * *`;
}

/**
 * Generate cron expression for "every N days" schedule
 * Note: Standard cron doesn't support "every N days" perfectly.
 * This uses day-of-month ranges as an approximation.
 */
export function everyNDaysCron(days: number, hour: number = 8, minute: number = 0): string {
	if (days === 1) {
		return dailyCron(hour, minute);
	}
	// For every N days, we'll use a pattern in day-of-month
	// This is an approximation - for exact "every N days", you'd need a more complex scheduler
	// Using */N pattern for day-of-month (runs on days 1, N+1, 2N+1, etc.)
	return `${minute} ${hour} */${days} * *`;
}

/**
 * Generate cron expression for "every 7 days" (weekly)
 */
export function every7DaysCron(hour: number = 8, minute: number = 0): string {
	// Every 7 days starting from day 1 of month
	return `${minute} ${hour} */7 * *`;
}

/**
 * Generate cron expression for "every 14 days" (bi-weekly)
 */
export function every14DaysCron(hour: number = 8, minute: number = 0): string {
	return `${minute} ${hour} */14 * *`;
}

/**
 * Generate cron expression for "every 30 days" (monthly)
 */
export function every30DaysCron(hour: number = 8, minute: number = 0): string {
	return monthlyCron(1, hour, minute);
}

/**
 * Get common presets for quick scheduling
 */
export function getCronPresets(hour: number = 8, minute: number = 0): CronPreset[] {
	return [
		{
			name: 'Daily',
			description: 'Every day at the same time',
			cron: dailyCron(hour, minute),
			frequency: 'daily'
		},
		{
			name: 'Every 7 Days',
			description: 'Every 7 days',
			cron: every7DaysCron(hour, minute),
			frequency: 'weekly'
		},
		{
			name: 'Every 14 Days',
			description: 'Every 14 days (bi-weekly)',
			cron: every14DaysCron(hour, minute),
			frequency: 'bi-weekly'
		},
		{
			name: 'Every 30 Days',
			description: 'Every 30 days (monthly)',
			cron: every30DaysCron(hour, minute),
			frequency: 'monthly'
		}
	];
}

/**
 * Parse cron expression to extract components
 */
export function parseCron(cron: string): {
	minute: string;
	hour: string;
	dayOfMonth: string;
	month: string;
	dayOfWeek: string;
} | null {
	const parts = cron.trim().split(/\s+/);
	if (parts.length !== 5) {
		return null;
	}
	return {
		minute: parts[0],
		hour: parts[1],
		dayOfMonth: parts[2],
		month: parts[3],
		dayOfWeek: parts[4]
	};
}

