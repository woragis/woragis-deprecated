import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface Client {
	id: string;
	userId: string;
	name: string;
	email?: string;
	phoneNumber: string;
	company?: string;
	notes?: string;
	isArchived: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface CreateClientInput {
	name: string;
	email?: string;
	phoneNumber: string;
	company?: string;
	notes?: string;
}

export interface UpdateClientInput {
	name?: string;
	email?: string;
	phoneNumber?: string;
	company?: string;
	notes?: string;
}

export async function listClients(): Promise<Client[]> {
	const response = await apiClient.get<ApiResponse<Client[]>>('/clients');
	return response.data.data ?? [];
}

export async function getClient(id: string): Promise<Client> {
	const response = await apiClient.get<ApiResponse<Client>>(`/clients/${id}`);
	return response.data.data;
}

export async function createClient(input: CreateClientInput): Promise<Client> {
	const response = await apiClient.post<ApiResponse<Client>>('/clients', input);
	return response.data.data;
}

export async function updateClient(id: string, input: UpdateClientInput): Promise<Client> {
	const response = await apiClient.patch<ApiResponse<Client>>(`/clients/${id}`, input);
	return response.data.data;
}

export async function toggleArchiveClient(id: string): Promise<Client> {
	const response = await apiClient.patch<ApiResponse<Client>>(`/clients/${id}/archive`);
	return response.data.data;
}

export async function deleteClient(id: string): Promise<void> {
	await apiClient.delete(`/clients/${id}`);
}

// Aliases for hooks compatibility
export async function fetchClients(includeArchived: boolean = false): Promise<Client[]> {
	const response = await apiClient.get<ApiResponse<Client[]>>(
		`/clients${includeArchived ? '?includeArchived=true' : ''}`
	);
	return response.data.data ?? [];
}

export async function toggleClientArchived(id: string, archived: boolean): Promise<Client> {
	if (archived) {
		return toggleArchiveClient(id);
	} else {
		const response = await apiClient.patch<ApiResponse<Client>>(`/clients/${id}/restore`);
		return response.data.data;
	}
}