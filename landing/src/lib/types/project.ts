export type ProjectStatus = 'idea' | 'planning' | 'executing' | 'monitoring' | 'completed';

export type TechnologyCategory =
	| 'language'
	| 'framework'
	| 'database'
	| 'tool'
	| 'service'
	| 'library'
	| 'other';

export interface Project {
	id: string;
	userId: string;
	name: string;
	description: string;
	slug: string;
	status: ProjectStatus;
	healthScore: number;
	mrr: number;
	cac: number;
	ltv: number;
	churnRate: number;
	createdAt: string;
	updatedAt: string;
	technologies?: ProjectTechnology[];
	skills?: import('../api/skills').Skill[];
}

export interface ProjectTechnology {
	id: string;
	projectId: string;
	name: string;
	version?: string;
	category: TechnologyCategory;
	purpose?: string;
	link?: string;
	createdAt: string;
	updatedAt: string;
}

export interface ProjectFilters {
	status?: ProjectStatus | 'all';
	search?: string;
	technology?: string;
	sortBy?: 'name' | 'createdAt' | 'updatedAt' | 'status' | 'healthScore';
	sortOrder?: 'asc' | 'desc';
}

export interface ApiResponse<T> {
	data?: T;
	message?: string;
	error?: string;
}

