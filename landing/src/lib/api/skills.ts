import { apiClient } from './client';

export type SkillCategory =
	| 'backend'
	| 'frontend'
	| 'database'
	| 'infrastructure'
	| 'devops'
	| 'language'
	| 'framework'
	| 'tool'
	| 'service'
	| 'library'
	| 'other';

export interface Skill {
	id: string;
	name: string;
	slug: string;
	category: SkillCategory;
	description?: string;
	icon?: string;
	createdAt: string;
	updatedAt: string;
}

export interface SkillWithCount extends Skill {
	projectCount: number;
}

export interface CreateSkillRequest {
	name: string;
	description?: string;
	icon?: string;
	category: SkillCategory;
}

export interface UpdateSkillRequest {
	name?: string;
	description?: string;
	icon?: string;
	category?: SkillCategory;
}

// List all skills
export async function listSkills(): Promise<Skill[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Skill[] }>('/skills');
		return response.data || [];
	} catch (error) {
		console.error('Error fetching skills:', error);
		throw error;
	}
}

// Get all skills with project counts
export async function listSkillsWithCounts(): Promise<SkillWithCount[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: SkillWithCount[] }>(
			'/skills/with-counts'
		);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching skills with counts:', error);
		throw error;
	}
}

// Get skill by ID
export async function getSkill(id: string): Promise<Skill | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Skill }>(`/skills/${id}`);
		return response.data || null;
	} catch (error) {
		console.error('Error fetching skill:', error);
		return null;
	}
}

// Get skill by slug
export async function getSkillBySlug(slug: string): Promise<Skill | null> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Skill }>(
			`/skills/slug/${slug}`
		);
		return response.data || null;
	} catch (error) {
		console.error('Error fetching skill by slug:', error);
		return null;
	}
}

// Search skills
export async function searchSkills(query: string): Promise<Skill[]> {
	try {
		const params = new URLSearchParams({ q: query });
		const response = await apiClient.get<{ success: boolean; data: Skill[] }>(
			`/skills/search?${params.toString()}`
		);
		return response.data || [];
	} catch (error) {
		console.error('Error searching skills:', error);
		throw error;
	}
}

// List skills by category
export async function listSkillsByCategory(category: SkillCategory): Promise<Skill[]> {
	try {
		const params = new URLSearchParams({ category });
		const response = await apiClient.get<{ success: boolean; data: Skill[] }>(
			`/skills/category?${params.toString()}`
		);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching skills by category:', error);
		throw error;
	}
}

// Create skill
export async function createSkill(skill: CreateSkillRequest): Promise<Skill> {
	try {
		const response = await apiClient.post<{ success: boolean; data: Skill }>('/skills', skill);
		return response.data;
	} catch (error) {
		console.error('Error creating skill:', error);
		throw error;
	}
}

// Update skill
export async function updateSkill(id: string, skill: UpdateSkillRequest): Promise<Skill> {
	try {
		const response = await apiClient.patch<{ success: boolean; data: Skill }>(
			`/skills/${id}`,
			skill
		);
		return response.data;
	} catch (error) {
		console.error('Error updating skill:', error);
		throw error;
	}
}

// Project-Skill relationship operations

// Get skills for a project
export async function getProjectSkills(projectId: string): Promise<Skill[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: Skill[] }>(
			`/skills/projects/${projectId}/skills`
		);
		return response.data || [];
	} catch (error) {
		console.error('Error fetching project skills:', error);
		return [];
	}
}

// Attach skill to project
export async function attachSkillToProject(
	projectId: string,
	skillId: string
): Promise<void> {
	try {
		await apiClient.post(`/skills/projects/${projectId}/skills/${skillId}`);
	} catch (error) {
		console.error('Error attaching skill to project:', error);
		throw error;
	}
}

// Detach skill from project
export async function detachSkillFromProject(
	projectId: string,
	skillId: string
): Promise<void> {
	try {
		await apiClient.delete(`/skills/projects/${projectId}/skills/${skillId}`);
	} catch (error) {
		console.error('Error detaching skill from project:', error);
		throw error;
	}
}

// Get projects using a skill
export async function getProjectsBySkill(skillId: string): Promise<string[]> {
	try {
		const response = await apiClient.get<{ success: boolean; data: { project_ids: string[] } }>(
			`/skills/${skillId}/projects`
		);
		return response.data?.project_ids || [];
	} catch (error) {
		console.error('Error fetching projects by skill:', error);
		return [];
	}
}

