import { apiClient } from './client';
import type { Project, ProjectFilters, ProjectTechnology } from '$lib/types/project';

export interface ListProjectsParams extends ProjectFilters {
	limit?: number;
	offset?: number;
}

export interface ProjectWithTechnologies extends Project {
	technologies: ProjectTechnology[];
}

// List all projects
export async function listProjects(
	filters?: ListProjectsParams
): Promise<ProjectWithTechnologies[]> {
	try {
		const params = new URLSearchParams();

		if (filters?.status && filters.status !== 'all') {
			params.append('status', filters.status);
		}
		if (filters?.search) {
			params.append('search', filters.search);
		}
		if (filters?.technology) {
			params.append('technology', filters.technology);
		}
		if (filters?.sortBy) {
			params.append('sortBy', filters.sortBy);
		}
		if (filters?.sortOrder) {
			params.append('sortOrder', filters.sortOrder);
		}
		if (filters?.limit) {
			params.append('limit', filters.limit.toString());
		}
		if (filters?.offset) {
			params.append('offset', filters.offset.toString());
		}

		const queryString = params.toString();
		const url = queryString ? `/projects?${queryString}` : '/projects';

		const response = await apiClient.get<{ success: boolean; data: ProjectWithTechnologies[] }>(url);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching projects:', error);
		throw error;
	}
}

// Get project by slug
export async function getProjectBySlug(slug: string): Promise<ProjectWithTechnologies | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: ProjectWithTechnologies }>(
			`/projects/slug/${slug}`
		);
		return response.data || null;
	} catch (error) {
		console.error('Error fetching project by slug:', error);
		return null;
	}
}

// Search projects by slug
export async function searchProjectsBySlug(query: string): Promise<ProjectWithTechnologies[]> {
	try {
		const params = new URLSearchParams({ slug: query });
		const response = await apiClient.get<{ success: boolean; data: ProjectWithTechnologies[] }>(
			`/projects/slug?${params.toString()}`
		);
		return response.data || [];
	} catch (error) {
		console.error('Error searching projects:', error);
		return [];
	}
}

// Get project technologies
export async function getProjectTechnologies(projectId: string): Promise<ProjectTechnology[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: ProjectTechnology[] }>(
			`/projects/${projectId}/technologies`
		);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching project technologies:', error);
		return [];
	}
}

