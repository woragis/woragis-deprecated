import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface Component {
	name: string;
	description: string;
	technology: string;
}

export interface ComponentsData {
	components?: Component[];
}

export interface SystemDesign {
	id: string;
	userId: string;
	title: string;
	description: string;
	components?: ComponentsData;
	dataFlow?: string;
	scalability?: string;
	reliability?: string;
	diagram?: string;
	featured: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface CreateSystemDesignInput {
	title: string;
	description: string;
	components?: ComponentsData;
	dataFlow?: string;
	scalability?: string;
	reliability?: string;
	diagram?: string;
	featured?: boolean;
}

export interface UpdateSystemDesignInput {
	title?: string;
	description?: string;
	components?: ComponentsData;
	dataFlow?: string;
	scalability?: string;
	reliability?: string;
	diagram?: string;
	featured?: boolean;
}

export async function listSystemDesigns(): Promise<SystemDesign[]> {
	const response = await apiClient.get<ApiResponse<SystemDesign[]>>('/system-designs');
	return response.data.data ?? [];
}

export async function getSystemDesign(id: string): Promise<SystemDesign> {
	const response = await apiClient.get<ApiResponse<SystemDesign>>(`/system-designs/${id}`);
	return response.data.data;
}

export async function createSystemDesign(input: CreateSystemDesignInput): Promise<SystemDesign> {
	const response = await apiClient.post<ApiResponse<SystemDesign>>('/system-designs', input);
	return response.data.data;
}

export async function updateSystemDesign(id: string, input: UpdateSystemDesignInput): Promise<SystemDesign> {
	const response = await apiClient.patch<ApiResponse<SystemDesign>>(`/system-designs/${id}`, input);
	return response.data.data;
}

export async function deleteSystemDesign(id: string): Promise<void> {
	await apiClient.delete(`/system-designs/${id}`);
}

