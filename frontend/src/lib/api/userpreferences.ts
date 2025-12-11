import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface UserPreferences {
	id: string;
	userId: string;
	defaultLanguage: string;
	defaultCurrency: string;
	defaultWebsite?: string;
	createdAt: string;
	updatedAt: string;
}

export interface UpdateUserPreferencesInput {
	defaultLanguage?: string;
	defaultCurrency?: string;
}

export async function getUserPreferences(): Promise<UserPreferences> {
	const response = await apiClient.get<ApiResponse<UserPreferences>>('/user-preferences');
	return response.data.data;
}

export async function updateUserPreferences(
	input: UpdateUserPreferencesInput
): Promise<UserPreferences> {
	const response = await apiClient.patch<ApiResponse<UserPreferences>>('/user-preferences', input);
	return response.data.data;
}

