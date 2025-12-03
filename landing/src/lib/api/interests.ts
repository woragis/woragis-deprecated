import { apiClient } from './client';
import type { Interest } from '$lib/types/interest';

// List all interests
export async function listInterests(): Promise<Interest[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Interest[] }>('/interests');
		return response.data || [];
	} catch (error) {
		console.error('Error fetching interests:', error);
		throw error;
	}
}

// Get featured interests
export async function getFeaturedInterests(): Promise<Interest[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Interest[] }>(
			'/interests/featured'
		);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching featured interests:', error);
		throw error;
	}
}

// Get interest by ID
export async function getInterest(id: string): Promise<Interest | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Interest }>(
			`/interests/${id}`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching interest ${id}:`, error);
		return null;
	}
}

// Get interest by slug
export async function getInterestBySlug(slug: string): Promise<Interest | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Interest }>(
			`/interests/slug/${slug}`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching interest by slug ${slug}:`, error);
		return null;
	}
}

