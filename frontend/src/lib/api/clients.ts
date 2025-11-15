import { apiClient } from '@clients/apiClient';

export interface Client {
	id: string;
	user_id: string;
	name: string;
	email?: string;
	phone_number: string;
	company?: string;
	notes?: string;
	is_archived: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateClientInput {
	name: string;
	email?: string;
	phone_number: string;
	company?: string;
	notes?: string;
}

export interface UpdateClientInput {
	name?: string;
	email?: string;
	phone_number?: string;
	company?: string;
	notes?: string;
}

export interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export async function fetchClients(includeArchived = false): Promise<Client[]> {
	const response = await apiClient.get<ApiResponse<Client[]>>('/clients', {
		params: {
			include_archived: includeArchived
		}
	});
	return response.data.data ?? [];
}

export async function fetchClient(id: string): Promise<Client> {
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

export async function deleteClient(id: string): Promise<void> {
	await apiClient.delete<ApiResponse<void>>(`/clients/${id}`);
}

export async function toggleClientArchived(id: string, archived: boolean): Promise<void> {
	await apiClient.patch<ApiResponse<void>>(`/clients/${id}/archive`, { archived });
}

