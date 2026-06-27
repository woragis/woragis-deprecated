import { apiClient } from '$lib/clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface PlatformConfig {
	id: string;
	name: string;
	displayName: string;
	postingFrequency?: number;
	bestDays?: string[];
	bestTimes?: string[];
	supportedFormats: string[];
	isActive: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface OptimalTimesResponse {
	platform: string;
	bestDays?: string[];
	bestTimes?: string[];
	postingFrequency?: number;
}

export interface UpdatePlatformConfigRequest {
	displayName?: string;
	postingFrequency?: number;
	bestDays?: string[];
	bestTimes?: string[];
	supportedFormats?: string[];
	isActive?: boolean;
}

export async function listPlatforms(activeOnly: boolean = false): Promise<PlatformConfig[]> {
	const response = await apiClient.get<ApiResponse<PlatformConfig[]>>(
		`/social-media-posts/platforms?activeOnly=${activeOnly}`
	);
	return response.data.data ?? [];
}

export async function getPlatformConfig(id: string): Promise<PlatformConfig> {
	const response = await apiClient.get<ApiResponse<PlatformConfig>>(
		`/social-media-posts/platforms/${id}`
	);
	return response.data.data;
}

export async function getPlatformConfigByName(name: string): Promise<PlatformConfig> {
	const response = await apiClient.get<ApiResponse<PlatformConfig>>(
		`/social-media-posts/platforms/by-name/${name}`
	);
	return response.data.data;
}

export async function getOptimalTimes(name: string): Promise<OptimalTimesResponse> {
	const response = await apiClient.get<ApiResponse<OptimalTimesResponse>>(
		`/social-media-posts/platforms/${name}/optimal-times`
	);
	return response.data.data;
}

export async function updatePlatformConfig(
	id: string,
	input: UpdatePlatformConfigRequest
): Promise<PlatformConfig> {
	const response = await apiClient.patch<ApiResponse<PlatformConfig>>(
		`/social-media-posts/platforms/${id}`,
		input
	);
	return response.data.data;
}
