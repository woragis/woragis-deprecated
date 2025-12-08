import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export type ExperienceType = 'full-time' | 'freelance' | 'contract' | 'internship';

export interface Experience {
	id: string;
	userId: string;
	company: string;
	position: string;
	periodStart?: string;
	periodEnd?: string;
	periodText?: string;
	location?: string;
	description?: string;
	type: ExperienceType;
	companyUrl?: string;
	linkedinUrl?: string;
	displayOrder: number;
	isCurrent: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface CreateExperienceInput {
	company: string;
	position: string;
	periodStart?: string;
	periodEnd?: string;
	periodText?: string;
	location?: string;
	description?: string;
	type?: ExperienceType;
	companyUrl?: string;
	linkedinUrl?: string;
	displayOrder?: number;
	isCurrent?: boolean;
}

export interface UpdateExperienceInput {
	company?: string;
	position?: string;
	periodStart?: string;
	periodEnd?: string;
	periodText?: string;
	location?: string;
	description?: string;
	type?: ExperienceType;
	companyUrl?: string;
	linkedinUrl?: string;
	displayOrder?: number;
	isCurrent?: boolean;
}

export async function listExperiences(): Promise<Experience[]> {
	const response = await apiClient.get<ApiResponse<Experience[]>>('/experiences');
	return response.data.data ?? [];
}

export async function getExperience(id: string): Promise<Experience> {
	const response = await apiClient.get<ApiResponse<Experience>>(`/experiences/${id}`);
	return response.data.data;
}

export async function createExperience(input: CreateExperienceInput): Promise<Experience> {
	const response = await apiClient.post<ApiResponse<Experience>>('/experiences', input);
	return response.data.data;
}

export async function updateExperience(id: string, input: UpdateExperienceInput): Promise<Experience> {
	const response = await apiClient.patch<ApiResponse<Experience>>(`/experiences/${id}`, input);
	return response.data.data;
}

export async function deleteExperience(id: string): Promise<void> {
	await apiClient.delete(`/experiences/${id}`);
}

