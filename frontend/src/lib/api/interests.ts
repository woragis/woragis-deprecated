import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface Interest {
	id: string;
	title: string;
	slug: string;
	description: string;
	icon?: string;
	color?: string;
	bgGradient?: string;
	borderColor?: string;
	hoverBorderColor?: string;
	shadowColor?: string;
	fullWidth: boolean;
	featured: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface CreateInterestInput {
	title: string;
	description: string;
	icon?: string;
	color?: string;
	bgGradient?: string;
	borderColor?: string;
	hoverBorderColor?: string;
	shadowColor?: string;
	fullWidth?: boolean;
	featured?: boolean;
}

export interface UpdateInterestInput {
	title?: string;
	description?: string;
	icon?: string;
	color?: string;
	bgGradient?: string;
	borderColor?: string;
	hoverBorderColor?: string;
	shadowColor?: string;
	fullWidth?: boolean;
	featured?: boolean;
}

export async function listInterests(): Promise<Interest[]> {
	const response = await apiClient.get<ApiResponse<Interest[]>>('/interests');
	return response.data.data ?? [];
}

export async function getInterest(id: string): Promise<Interest> {
	const response = await apiClient.get<ApiResponse<Interest>>(`/interests/${id}`);
	return response.data.data;
}

export async function createInterest(input: CreateInterestInput): Promise<Interest> {
	const response = await apiClient.post<ApiResponse<Interest>>('/interests', input);
	return response.data.data;
}

export async function updateInterest(id: string, input: UpdateInterestInput): Promise<Interest> {
	const response = await apiClient.patch<ApiResponse<Interest>>(`/interests/${id}`, input);
	return response.data.data;
}

export async function deleteInterest(id: string): Promise<void> {
	await apiClient.delete(`/interests/${id}`);
}

