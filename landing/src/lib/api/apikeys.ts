import { apiClient } from './client';

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

export interface CreateAPIKeyRequest {
	name: string;
	expiresAt?: string; // ISO 8601 format
}

export interface UpdateAPIKeyRequest {
	name: string;
}

// List all API keys
export async function listAPIKeys(): Promise<APIKey[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: APIKey[] }>('/api-keys');
		return response.data || [];
	} catch (error) {
		console.error('Error fetching API keys:', error);
		throw error;
	}
}

// Get API key by ID
export async function getAPIKey(id: string): Promise<APIKey | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: APIKey }>(`/api-keys/${id}`);
		return response.data || null;
	} catch (error) {
		console.error('Error fetching API key:', error);
		return null;
	}
}

// Create API key
export async function createAPIKey(
	apiKey: CreateAPIKeyRequest
): Promise<APIKeyWithToken> {
	try {
		const response = await apiClient.post<{ success: boolean; data: APIKeyWithToken }>(
			'/api-keys',
			apiKey
		);
		return response.data;
	} catch (error) {
		console.error('Error creating API key:', error);
		throw error;
	}
}

// Update API key
export async function updateAPIKey(id: string, apiKey: UpdateAPIKeyRequest): Promise<APIKey> {
	try {
		const response = await apiClient.patch<{ success: boolean; data: APIKey }>(
			`/api-keys/${id}`,
			apiKey
		);
		return response.data;
	} catch (error) {
		console.error('Error updating API key:', error);
		throw error;
	}
}

// Delete API key
export async function deleteAPIKey(id: string): Promise<void> {
	try {
		await apiClient.delete(`/api-keys/${id}`);
	} catch (error) {
		console.error('Error deleting API key:', error);
		throw error;
	}
}

