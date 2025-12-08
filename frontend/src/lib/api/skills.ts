import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

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

export type ProficiencyLevel = 'expert' | 'advanced' | 'proficient' | 'learning';

export interface Skill {
	id: string;
	name: string;
	slug: string;
	category: SkillCategory;
	description?: string;
	icon?: string;
	color?: string;
	bgGradient?: string;
	borderColor?: string;
	hoverBorderColor?: string;
	shadowColor?: string;
	proficiencyLevel?: ProficiencyLevel;
	yearsOfExperience?: number;
	firstUsedDate?: string;
	lastUsedDate?: string;
	createdAt: string;
	updatedAt: string;
}

export interface CreateSkillInput {
	name: string;
	description?: string;
	icon?: string;
	color?: string;
	bgGradient?: string;
	borderColor?: string;
	hoverBorderColor?: string;
	shadowColor?: string;
	category: SkillCategory;
	proficiencyLevel?: ProficiencyLevel;
	yearsOfExperience?: number;
	firstUsedDate?: string;
	lastUsedDate?: string;
}

export interface UpdateSkillInput {
	name?: string;
	description?: string;
	icon?: string;
	color?: string;
	bgGradient?: string;
	borderColor?: string;
	hoverBorderColor?: string;
	shadowColor?: string;
	category?: SkillCategory;
	proficiencyLevel?: ProficiencyLevel;
	yearsOfExperience?: number;
	firstUsedDate?: string;
	lastUsedDate?: string;
}

export async function listSkills(): Promise<Skill[]> {
	const response = await apiClient.get<ApiResponse<Skill[]>>('/skills');
	return response.data.data ?? [];
}

export async function getSkill(id: string): Promise<Skill> {
	const response = await apiClient.get<ApiResponse<Skill>>(`/skills/${id}`);
	return response.data.data;
}

export async function getSkillBySlug(slug: string): Promise<Skill> {
	const response = await apiClient.get<ApiResponse<Skill>>(`/skills/slug/${slug}`);
	return response.data.data;
}

export async function createSkill(input: CreateSkillInput): Promise<Skill> {
	const response = await apiClient.post<ApiResponse<Skill>>('/skills', input);
	return response.data.data;
}

export async function updateSkill(id: string, input: UpdateSkillInput): Promise<Skill> {
	const response = await apiClient.patch<ApiResponse<Skill>>(`/skills/${id}`, input);
	return response.data.data;
}

export async function searchSkills(query: string): Promise<Skill[]> {
	const response = await apiClient.get<ApiResponse<Skill[]>>(`/skills/search?q=${encodeURIComponent(query)}`);
	return response.data.data ?? [];
}

export async function listSkillsByCategory(category: SkillCategory): Promise<Skill[]> {
	const response = await apiClient.get<ApiResponse<Skill[]>>(`/skills/category?category=${category}`);
	return response.data.data ?? [];
}

export async function deleteSkill(id: string): Promise<void> {
	await apiClient.delete(`/skills/${id}`);
}

