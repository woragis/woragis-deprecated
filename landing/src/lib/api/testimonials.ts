import { apiClient } from './client';
import type { Testimonial, ListTestimonialsParams } from '$lib/types/testimonial';

// List all testimonials
export async function listTestimonials(
	params?: ListTestimonialsParams
): Promise<Testimonial[]> {
	try {
		const queryParams = new URLSearchParams();

		if (params?.status) {
			queryParams.append('status', params.status);
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
		const url = queryString ? `/testimonials?${queryString}` : '/testimonials';

		const response = await apiClient.get<{ success: boolean; data: Testimonial[] }>(url);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching testimonials:', error);
		throw error;
	}
}

// Get testimonial by ID
export async function getTestimonial(id: string): Promise<Testimonial | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Testimonial }>(
			`/testimonials/${id}`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching testimonial ${id}:`, error);
		return null;
	}
}

