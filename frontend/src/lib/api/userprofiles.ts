import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface UserProfile {
	id: string;
	userId: string;
	aboutMe: string;
	createdAt: string;
	updatedAt: string;
}

export interface UpdateUserProfileInput {
	aboutMe: string;
}

// Get current user's profile
export async function getUserProfile(): Promise<UserProfile> {
	const response = await apiClient.get<ApiResponse<UserProfile>>('/profile');
	return response.data.data;
}

// Create or update user profile
export async function upsertUserProfile(input: UpdateUserProfileInput): Promise<UserProfile> {
	const response = await apiClient.put<ApiResponse<UserProfile>>('/profile', input);
	return response.data.data;
}

