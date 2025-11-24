import { queryOptions, createQuery } from '@tanstack/svelte-query';
import {
	listProblemSolutions,
	getFeaturedProblemSolutions,
	getProblemSolution,
	getProblemSolutionMatrix,
	type ProblemSolution
} from '$lib/api/problem-solutions';
import type { ProblemSolutionMatrixEntry } from '$lib/types/problem-solution';

// Query keys factory
export const problemSolutionKeys = {
	all: ['problem-solutions'] as const,
	lists: () => [...problemSolutionKeys.all, 'list'] as const,
	list: (params?: { featured?: boolean }) => [...problemSolutionKeys.lists(), params] as const,
	featured: () => [...problemSolutionKeys.all, 'featured'] as const,
	details: () => [...problemSolutionKeys.all, 'detail'] as const,
	detail: (id: string) => [...problemSolutionKeys.details(), id] as const,
	matrix: () => [...problemSolutionKeys.all, 'matrix'] as const
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

// Hook for featured problem solutions
export function useFeaturedProblemSolutionsQuery() {
	return createQuery(() => ({
		queryKey: problemSolutionKeys.featured(),
		queryFn: () => getFeaturedProblemSolutions()
	}));
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

