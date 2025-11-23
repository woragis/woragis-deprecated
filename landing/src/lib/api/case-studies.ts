import { apiClient } from './client';
import type { CaseStudy } from '$lib/types/case-study';

export interface ListCaseStudiesParams {
	projectId?: string;
	projectSlug?: string;
	featured?: boolean;
	limit?: number;
	offset?: number;
	orderBy?: string;
	order?: 'asc' | 'desc';
}

// List all case studies
export async function listCaseStudies(
	params?: ListCaseStudiesParams
): Promise<CaseStudy[]> {
	try {
		const queryParams = new URLSearchParams();

		if (params?.projectId) {
			queryParams.append('projectId', params.projectId);
		}
		if (params?.projectSlug) {
			queryParams.append('projectSlug', params.projectSlug);
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
		const url = queryString ? `/case-studies?${queryString}` : '/case-studies';

		const response = await apiClient.get<{ success: boolean; data: CaseStudy[] }>(url);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching case studies:', error);
		throw error;
	}
}

// Get case study by ID
export async function getCaseStudy(id: string): Promise<CaseStudy | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: CaseStudy }>(
			`/case-studies/${id}`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching case study ${id}:`, error);
		return null;
	}
}

// Get case study by project slug
export async function getCaseStudyByProjectSlug(projectSlug: string): Promise<CaseStudy | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: CaseStudy }>(
			`/case-studies/project-slug/${projectSlug}`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching case study for project ${projectSlug}:`, error);
		return null;
	}
}

