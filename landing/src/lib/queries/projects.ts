import { queryOptions, createQuery } from '@tanstack/svelte-query';
import { derived, get } from 'svelte/store';
import { language } from '$lib/i18n';
import { listProjects, getProjectBySlug, type ListProjectsParams, type ProjectWithTechnologies } from '$lib/api/projects';

// Query keys factory - includes language for proper cache separation
export const projectKeys = {
	all: ['projects'] as const,
	lists: () => [...projectKeys.all, 'list'] as const,
	list: (filters: ListProjectsParams, lang: string) => [...projectKeys.lists(), filters, lang] as const,
	details: () => [...projectKeys.all, 'detail'] as const,
	detail: (slug: string, lang?: string) => [...projectKeys.details(), slug, lang] as const
};

// Query options for listing projects
export function getProjectsQueryOptions(filters: ListProjectsParams, lang: string) {
	return queryOptions({
		queryKey: projectKeys.list(filters, lang),
		queryFn: () => listProjects(filters),
		enabled: true
	});
}

// Query options for getting a project by slug
export function getProjectBySlugQueryOptions(slug: string, lang: string) {
	return queryOptions({
		queryKey: projectKeys.detail(slug, lang),
		queryFn: () => getProjectBySlug(slug),
		enabled: !!slug
	});
}

// Hook for listing projects - reactive to language changes
export function useProjectsQuery(filters: ListProjectsParams, lang: string) {
	// In v6, createQuery expects an Accessor function that's reactive
	// Access the language store directly - TanStack Query will track changes
	return createQuery(() => {
		// Read language from store - this should be tracked by TanStack Query
		const currentLang = get(language);
		return {
			queryKey: projectKeys.list(filters, currentLang),
			queryFn: () => listProjects(filters)
		};
	});
}

// Hook for getting a project by slug - reactive to language changes
export function useProjectBySlugQuery(slug: string, lang: string) {
	// In v6, createQuery expects an Accessor function that's reactive
	// Access the language store directly - TanStack Query will track changes
	return createQuery(() => {
		// Read language from store - this should be tracked by TanStack Query
		const currentLang = get(language);
		return {
			queryKey: projectKeys.detail(slug, currentLang),
			queryFn: () => getProjectBySlug(slug),
			enabled: !!slug
		};
	});
}

