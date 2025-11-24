import { apiClient } from './client';
import type { TechnicalWriting, ListTechnicalWritingsParams } from '$lib/types/technical-writing';

// List all technical writings
export async function listTechnicalWritings(
	params?: ListTechnicalWritingsParams
): Promise<TechnicalWriting[]> {
	try {
		const queryParams = new URLSearchParams();

		if (params?.type) {
			queryParams.append('type', params.type);
		}
		if (params?.platform) {
			queryParams.append('platform', params.platform);
		}
		if (params?.projectId) {
			queryParams.append('projectId', params.projectId);
		}
		if (params?.featured !== undefined) {
			queryParams.append('featured', params.featured.toString());
		}
		if (params?.limit) {
			queryParams.append('limit', params.limit.toString());
		}
		if (params?.offset) {
			queryParams.append('offset', params.offset.toString());
		}
		if (params?.orderBy) {
			queryParams.append('orderBy', params.orderBy);
		}
		if (params?.order) {
			queryParams.append('order', params.order);
		}

		const queryString = queryParams.toString();
		const url = queryString ? `/technical-writings?${queryString}` : '/technical-writings';

		const response = await apiClient.get<{ success: boolean; data: TechnicalWriting[] }>(url);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching technical writings:', error);
		throw error;
	}
}

// Get featured technical writings (public endpoint)
export async function getFeaturedTechnicalWritings(): Promise<TechnicalWriting[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: TechnicalWriting[] }>(
			'/technical-writings/featured'
		);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching featured technical writings:', error);
		throw error;
	}
}

// Get technical writing by ID
export async function getTechnicalWriting(id: string): Promise<TechnicalWriting | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: TechnicalWriting }>(
			`/technical-writings/${id}/public`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching technical writing ${id}:`, error);
		return null;
	}
}

// Search technical writings
export async function searchTechnicalWritings(query: string): Promise<TechnicalWriting[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: TechnicalWriting[] }>(
			`/technical-writings/search?q=${encodeURIComponent(query)}`
		);
		return response.data || [];
	} catch (error) {
		console.error(`Error searching technical writings:`, error);
		return [];
	}
}

// Get writings by type
export async function getWritingsByType(type: string): Promise<TechnicalWriting[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: TechnicalWriting[] }>(
			`/technical-writings/type/${type}`
		);
		return response.data || [];
	} catch (error) {
		console.error(`Error fetching technical writings by type ${type}:`, error);
		return [];
	}
}

// Get writings by platform
export async function getWritingsByPlatform(platform: string): Promise<TechnicalWriting[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: TechnicalWriting[] }>(
			`/technical-writings/platform/${platform}`
		);
		return response.data || [];
	} catch (error) {
		console.error(`Error fetching technical writings by platform ${platform}:`, error);
		return [];
	}
}

