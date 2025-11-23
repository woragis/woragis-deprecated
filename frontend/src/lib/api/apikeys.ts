import { apiClient } from '@clients/apiClient';

export interface APIKey {
	id: string;
	name: string;
	prefix: string;
	userId: string;
	lastUsedAt?: string;
	expiresAt?: string;
	createdAt: string;
	updatedAt: string;
}

export interface APIKeyWithToken extends APIKey {
	token: string; // Only included on creation
}

export interface CreateAPIKeyPayload {
	name: string;
	expiresAt?: string; // ISO 8601 format
}

export interface UpdateAPIKeyPayload {
	name: string;
}

export interface APIKeyResponse {
	success: boolean;
	data: APIKey | APIKey[] | APIKeyWithToken;
	message?: string;
	error?: string;
}

// List all API keys
export async function listAPIKeys(): Promise<APIKey[]> {
	const response = await apiClient.get<APIKeyResponse>('/api-keys');
	if (response.data.success && Array.isArray(response.data.data)) {
		return response.data.data;
	}
	return [];
}

// Get API key by ID
export async function getAPIKey(id: string): Promise<APIKey | null> {
	try {
		const response = await apiClient.get<APIKeyResponse>(`/api-keys/${id}`);
		if (response.data.success && !Array.isArray(response.data.data)) {
			return response.data.data as APIKey;
		}
		return null;
	} catch (error) {
		console.error('Error fetching API key:', error);
		return null;
	}
}

// Create API key
export async function createAPIKey(payload: CreateAPIKeyPayload): Promise<APIKeyWithToken> {
	const response = await apiClient.post<APIKeyResponse>('/api-keys', payload);
	if (response.data.success && !Array.isArray(response.data.data)) {
		return response.data.data as APIKeyWithToken;
	}
	throw new Error('Failed to create API key');
}

// Update API key
export async function updateAPIKey(id: string, payload: UpdateAPIKeyPayload): Promise<APIKey> {
	const response = await apiClient.patch<APIKeyResponse>(`/api-keys/${id}`, payload);
	if (response.data.success && !Array.isArray(response.data.data)) {
		return response.data.data as APIKey;
	}
	throw new Error('Failed to update API key');
}

// Delete API key
export async function deleteAPIKey(id: string): Promise<void> {
	await apiClient.delete(`/api-keys/${id}`);
}

