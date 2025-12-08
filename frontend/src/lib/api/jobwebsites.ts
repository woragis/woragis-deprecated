import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface JobWebsite {
	id: string;
	name: string;
	displayName: string;
	dailyLimit: number;
	currentCount: number;
	lastReset: string;
	enabled: boolean;
	baseUrl?: string;
	loginUrl?: string;
	createdAt: string;
	updatedAt: string;
}

export interface CreateJobWebsiteInput {
	name: string;
	displayName: string;
	dailyLimit: number;
	baseUrl?: string;
	loginUrl?: string;
	enabled?: boolean;
}

export interface UpdateJobWebsiteInput {
	displayName?: string;
	dailyLimit?: number;
	baseUrl?: string;
	loginUrl?: string;
	enabled?: boolean;
}

export async function listJobWebsites(): Promise<JobWebsite[]> {
	const response = await apiClient.get<ApiResponse<JobWebsite[]>>('/job-websites');
	return response.data.data ?? [];
}

export async function getJobWebsite(id: string): Promise<JobWebsite> {
	const response = await apiClient.get<ApiResponse<JobWebsite>>(`/job-websites/${id}`);
	return response.data.data;
}

export async function createJobWebsite(input: CreateJobWebsiteInput): Promise<JobWebsite> {
	const response = await apiClient.post<ApiResponse<JobWebsite>>('/job-websites', input);
	return response.data.data;
}

export async function updateJobWebsite(id: string, input: UpdateJobWebsiteInput): Promise<JobWebsite> {
	const response = await apiClient.patch<ApiResponse<JobWebsite>>(`/job-websites/${id}`, input);
	return response.data.data;
}

export async function deleteJobWebsite(id: string): Promise<void> {
	await apiClient.delete(`/job-websites/${id}`);
}

export async function resetJobWebsiteCounter(id: string): Promise<void> {
	await apiClient.post(`/job-websites/${id}/reset-counter`);
}

