import { apiClient } from './client';
import type { ProblemSolution, ListProblemSolutionsParams } from '$lib/types/problem-solution';

// List all problem solutions
export async function listProblemSolutions(
	params?: ListProblemSolutionsParams
): Promise<ProblemSolution[]> {
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
		const url = queryString ? `/problem-solutions?${queryString}` : '/problem-solutions';

		const response = await apiClient.get<{ success: boolean; data: ProblemSolution[] }>(url);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching problem solutions:', error);
		throw error;
	}
}

// Get featured problem solutions (public endpoint)
export async function getFeaturedProblemSolutions(): Promise<ProblemSolution[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: ProblemSolution[] }>(
			'/problem-solutions/featured'
		);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching featured problem solutions:', error);
		throw error;
	}
}

// Get problem solution by ID
export async function getProblemSolution(id: string): Promise<ProblemSolution | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: ProblemSolution }>(
			`/problem-solutions/${id}/public`
		);
		return response.data || null;
	} catch (error) {
		console.error(`Error fetching problem solution ${id}:`, error);
		return null;
	}
}

