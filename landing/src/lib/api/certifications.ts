import { apiClient } from './client';
import type { Certification, ListCertificationsParams } from '$lib/types/certification';

// List all certifications
export async function listCertifications(
	params?: ListCertificationsParams
): Promise<Certification[]> {
	try {
		const queryParams = new URLSearchParams();

		if (params?.status) {
			queryParams.append('status', params.status);
		}
		if (params?.category) {
			queryParams.append('category', params.category);
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
		const url = queryString ? `/certifications?${queryString}` : '/certifications';

		const response = await apiClient.get<{ success: boolean; data: Certification[] }>(url);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching certifications:', error);
		throw error;
	}
}

// Get featured certifications (public endpoint)
export async function getFeaturedCertifications(): Promise<Certification[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Certification[] }>(
			'/certifications/featured'
		);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching featured certifications:', error);
		throw error;
	}
}

// Get certification by ID
export async function getCertification(id: string): Promise<Certification | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Certification }>(
			`/certifications/${id}/public`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching certification ${id}:`, error);
		return null;
	}
}

// Get certifications by skill
export async function getCertificationsBySkill(skillId: string): Promise<Certification[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Certification[] }>(
			`/certifications/skill/${skillId}`
		);
		return response.data || [];
	} catch (error) {
		console.error(`Error fetching certifications for skill ${skillId}:`, error);
		return [];
	}
}

