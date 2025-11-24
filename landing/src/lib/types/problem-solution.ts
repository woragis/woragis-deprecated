export interface MetricsData {
	before: string;
	after: string;
	improvement: string;
}

export interface ProblemSolution {
	id: string;
	userId: string;
	problem: string;
	context: string;
	solution: string;
	technologies: string[];
	impact: string;
	metrics?: MetricsData;
	featured: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface ListProblemSolutionsParams {
	featured?: boolean;
	limit?: number;
	offset?: number;
	orderBy?: string;
	order?: 'asc' | 'desc';
}

export interface ProblemSolutionMatrixEntry {
	technology: string;
	problems: string[];
	count: number;
}

