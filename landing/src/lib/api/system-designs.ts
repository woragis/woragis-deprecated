import { apiClient } from './client';
import type { SystemDesign, ListSystemDesignsParams } from '$lib/types/system-design';

// List all system designs
export async function listSystemDesigns(
	params?: ListSystemDesignsParams
): Promise<SystemDesign[]> {
	try {
		const queryParams = new URLSearchParams();

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
		const url = queryString ? `/system-designs?${queryString}` : '/system-designs';

		const response = await apiClient.get<{ success: boolean; data: SystemDesign[] }>(url);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching system designs:', error);
		throw error;
	}
}

// Get featured system designs (public endpoint)
export async function getFeaturedSystemDesigns(): Promise<SystemDesign[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: SystemDesign[] }>(
			'/system-designs/featured'
		);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching featured system designs:', error);
		throw error;
	}
}

// Get system design by ID
export async function getSystemDesign(id: string): Promise<SystemDesign | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: SystemDesign }>(
			`/system-designs/${id}/public`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching system design ${id}:`, error);
		return null;
	}
}

