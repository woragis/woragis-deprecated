import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { listProjects, getProjectBySlug, type ListProjectsParams, type ProjectWithTechnologies } from '$lib/api/projects';

// Query keys factory
export const projectKeys = {
	all: ['projects'] as const,
	lists: () => [...projectKeys.all, 'list'] as const,
	list: (filters?: ListProjectsParams) => [...projectKeys.lists(), filters] as const,
	details: () => [...projectKeys.all, 'detail'] as const,
	detail: (slug: string) => [...projectKeys.details(), slug] as const
};

// Query options for listing projects
export function getProjectsQueryOptions(filters?: ListProjectsParams) {
	return queryOptions({
		queryKey: projectKeys.list(filters),
		queryFn: () => listProjects(filters),
		enabled: true
	});
}

// Query options for getting a project by slug
export function getProjectBySlugQueryOptions(slug: string) {
	return queryOptions({
		queryKey: projectKeys.detail(slug),
		queryFn: () => getProjectBySlug(slug),
		enabled: !!slug
	});
}

// Hook for listing projects
export function useProjectsQuery(filters?: ListProjectsParams) {
	return createQuery(() => ({
		queryKey: projectKeys.list(filters),
		queryFn: () => listProjects(filters)
	}));
}

// Hook for getting a project by slug
export function useProjectBySlugQuery(slug: string) {
	return createQuery(() => ({
		queryKey: projectKeys.detail(slug),
		queryFn: () => getProjectBySlug(slug),
		enabled: !!slug
	}));
}

