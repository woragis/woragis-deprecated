import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { get } from 'svelte/store';
import { language } from '$lib/i18n';
import {
	listProblemSolutions,
	getFeaturedProblemSolutions,
	getProblemSolution,
	getProblemSolutionMatrix
} from '$lib/api/problem-solutions';
import type { ProblemSolution, ProblemSolutionMatrixEntry } from '$lib/types/problem-solution';

// Query keys factory - includes language for proper cache separation
export const problemSolutionKeys = {
	all: ['problem-solutions'] as const,
	lists: () => [...problemSolutionKeys.all, 'list'] as const,
	list: (params?: { featured?: boolean }, lang?: string) => [...problemSolutionKeys.lists(), params, lang] as const,
	featured: (lang?: string) => [...problemSolutionKeys.all, 'featured', lang] as const,
	details: () => [...problemSolutionKeys.all, 'detail'] as const,
	detail: (id: string, lang?: string) => [...problemSolutionKeys.details(), id, lang] as const,
	matrix: (lang?: string) => [...problemSolutionKeys.all, 'matrix', lang] as const
};

// Query options for listing problem solutions
export function getProblemSolutionsQueryOptions(params?: { featured?: boolean }) {
	return queryOptions({
		queryKey: problemSolutionKeys.list(params),
		queryFn: () => listProblemSolutions(params)
	});
}

// Query options for featured problem solutions
export function getFeaturedProblemSolutionsQueryOptions() {
	return queryOptions({
		queryKey: problemSolutionKeys.featured(),
		queryFn: () => getFeaturedProblemSolutions()
	});
}

// Query options for getting a problem solution by ID
export function getProblemSolutionQueryOptions(id: string) {
	return queryOptions({
		queryKey: problemSolutionKeys.detail(id),
		queryFn: () => getProblemSolution(id),
		enabled: !!id
	});
}

// Hook for listing problem solutions
export function useProblemSolutionsQuery(params?: { featured?: boolean }) {
	return createQuery(() => ({
		queryKey: problemSolutionKeys.list(params),
		queryFn: () => listProblemSolutions(params)
	}));
}

// Hook for featured problem solutions - reactive to language changes
export function useFeaturedProblemSolutionsQuery(lang?: string) {
	return createQuery(() => {
		const currentLang = lang ?? get(language);
		return {
			queryKey: problemSolutionKeys.featured(currentLang),
			queryFn: () => getFeaturedProblemSolutions()
		};
	});
}

// Hook for getting a problem solution by ID
export function useProblemSolutionQuery(id: string) {
	return createQuery(() => ({
		queryKey: problemSolutionKeys.detail(id),
		queryFn: () => getProblemSolution(id),
		enabled: !!id
	}));
}

// Query options for getting problem-solution matrix
export function getProblemSolutionMatrixQueryOptions() {
	return queryOptions({
		queryKey: problemSolutionKeys.matrix(),
		queryFn: () => getProblemSolutionMatrix()
	});
}

// Hook for getting problem-solution matrix
export function useProblemSolutionMatrixQuery() {
	return createQuery(() => ({
		queryKey: problemSolutionKeys.matrix(),
		queryFn: () => getProblemSolutionMatrix()
	}));
}

