import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface Schedule {
	id: string;
	userId: string;
	reportType: string;
	agentAlias: string;
	frequency: string;
	weekday?: string;
	timeOfDay: string;
	timezone: string;
	rrule?: string;
	priority: number;
	email?: string;
	phoneNumber?: string;
	channels?: Record<string, boolean>;
	active: boolean;
	paused: boolean;
	nextRun: string;
	lastRun?: string;
	createdAt: string;
	updatedAt: string;
}

export interface CreateScheduleInput {
	reportType: string;
	agentAlias: string;
	frequency: string;
	weekday?: string;
	timeOfDay: string;
	timezone: string;
	rrule?: string;
	priority?: number;
	email?: string;
	phoneNumber?: string;
	channels?: Record<string, boolean>;
	active?: boolean;
}

export interface UpdateScheduleInput {
	reportType?: string;
	agentAlias?: string;
	frequency?: string;
	weekday?: string;
	timeOfDay?: string;
	timezone?: string;
	rrule?: string;
	priority?: number;
	email?: string;
	phoneNumber?: string;
	channels?: Record<string, boolean>;
	active?: boolean;
	paused?: boolean;
}

export async function listSchedules(): Promise<Schedule[]> {
	const response = await apiClient.get<ApiResponse<Schedule[]>>('/scheduler');
	return response.data.data ?? [];
}

export async function createSchedule(input: CreateScheduleInput): Promise<Schedule> {
	const response = await apiClient.post<ApiResponse<Schedule>>('/scheduler', input);
	return response.data.data;
}

export async function updateSchedule(id: string, input: UpdateScheduleInput): Promise<Schedule> {
	const response = await apiClient.patch<ApiResponse<Schedule>>(`/scheduler/${id}`, input);
	return response.data.data;
}

export async function bulkActivate(ids: string[]): Promise<void> {
	await apiClient.post('/scheduler/bulk/activate', { ids });
}

export async function bulkDeactivate(ids: string[]): Promise<void> {
	await apiClient.post('/scheduler/bulk/deactivate', { ids });
}

export async function bulkPause(ids: string[]): Promise<void> {
	await apiClient.post('/scheduler/bulk/pause', { ids });
}

export async function bulkResume(ids: string[]): Promise<void> {
	await apiClient.post('/scheduler/bulk/resume', { ids });
}

export async function deleteSchedule(id: string): Promise<void> {
	await apiClient.delete(`/scheduler/${id}`);
}

